package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	elasticsearch6 "github.com/elastic/go-elasticsearch/v6"
	esapi6 "github.com/elastic/go-elasticsearch/v6/esapi"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	esapi7 "github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"regexp"
	"strings"
	"time"
)

var (
	client6 *elasticsearch6.Client
	client7 *elasticsearch7.Client

	buf  bytes.Buffer
	raw  map[string]interface{}

	numItems int
)

type ES struct {
	index    []string
	shards   string
	replicas string
	version  string
	Content

	detailEnable bool
	reg          string
	template     string
}

type Content struct {
	Content   string `json:"content"`
}

func NewEs(addresses,index []string,shards,replicas,version string,detailEnable bool,reg,template string) (es *ES,err error) {
	if version == "6" {
		client6, err = elasticsearch6.NewClient(elasticsearch6.Config{
			Addresses: addresses,
		})
	}
	if version == "7" {
		client7, err = elasticsearch7.NewClient(elasticsearch7.Config{
			Addresses: addresses,
		})
	}
	if err != nil {
		log.Println(err)
		return nil,err
	}

	es = &ES{
		index:index,
		shards:shards,
		replicas:replicas,
		version:version,
		detailEnable:detailEnable,
		reg:reg,
		template:template,
	}
	return es,nil
}

func (c *ES) ToEs(content chan string) {
	if !c.indexIsExist() {
		c.indexCreate()
	}
	if c.version == "6" {
		c.toEs6(content)
	}
	if c.version == "7" {
		if c.detailEnable == true {
			c.toEs7Detail(content)
		}
		if c.detailEnable == false {
			c.toEs7(content)
		}
	}
}

func (c *ES) indexIsExist() bool {
	if c.version == "6" {
		res, err := client6.Indices.Get(c.index)
		if err != nil {
			log.Fatal(err)
			return false
		}
		if res.IsError() {
			log.Println("No index found: ",c.index)
			return false
		}
		defer res.Body.Close()

		log.Print("Find the index: ",c.index)
		return true
	}

	if c.version == "7" {
		res, err := client7.Indices.Get(c.index)
		if err != nil {
			log.Fatal(err)
			return false
		}
		if res.IsError() {
			log.Println("No index found: ",c.index)
			return false
		}
		defer res.Body.Close()

		log.Print("Find the index: ",c.index)
		return true
	}
	return false
}

func (c *ES) indexCreate() {
	var b strings.Builder
	b.WriteString(`{"settings" : {"number_of_shards" :`)
	b.WriteString(c.shards)
	b.WriteString(`,"number_of_replicas" :`)
	b.WriteString(c.replicas)
	b.WriteString(`}}`)

	if c.version == "6" {
		res, err := client6.Indices.Create(c.index[0], func(request *esapi6.IndicesCreateRequest) {
			request.Body = strings.NewReader(b.String())
		})
		if err != nil {
			log.Printf("Malformed response to create index: %s", err)
			return
		}
		if res.IsError() {
			log.Printf("Cannot create index: %s", res)
			return
		}
		defer res.Body.Close()

		log.Print("Successfully created index: ",c.index)
		return
	}
	if c.version == "7" {
		res, err := client7.Indices.Create(c.index[0], func(request *esapi7.IndicesCreateRequest) {
			request.Body = strings.NewReader(b.String())
		})
		if err != nil {
			log.Printf("Malformed response to create index: %s", err)
			return
		}
		if res.IsError() {
			log.Printf("Cannot create index: %s", res)
			return
		}
		defer res.Body.Close()

		log.Print("Successfully created index: ",c.index)
		return
	}
}

func (c *ES) toEs7(content chan string)  {
	for {
		select {
		case v := <-content:
			numItems++
			meta := []byte(fmt.Sprintf(`{ "index":{}}%s`, "\n"))

			c.Content.Content = v

			data, err := json.Marshal(c.Content)
			if err != nil {
				log.Fatalf("Cannot encode article  %s",  err)
			}
			data = append(data, "\n"...) // <-- Comment out to trigger failure for batch

			buf.Grow(len(meta) + len(data))
			buf.Write(meta)
			buf.Write(data)

			// When a threshold is reached, execute the Bulk() request with body from buffer
			if numItems == 2000  {
				res, err := client7.Bulk(bytes.NewReader(buf.Bytes()), client7.Bulk.WithIndex(c.index[0]))
				if err != nil {
					log.Printf("Failure indexing batch %s", err)
					continue
				}

				// If the whole request failed, print error and mark all documents as failed
				if res.IsError() {
					if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
						log.Fatalf("Failure to to parse response body: %s", err)
					} else {
						log.Printf("  Error: [%d] %s: %s",
							res.StatusCode,
							raw["error"].(map[string]interface{})["type"],
							raw["error"].(map[string]interface{})["reason"],
						)
					}
				}

				res.Body.Close()
				buf.Reset()

				numItems = 0
			}
		case <-time.After(time.Second * 1):
			if len(buf.Bytes()) != 0 {
				res, err := client7.Bulk(bytes.NewReader(buf.Bytes()), client7.Bulk.WithIndex(c.index[0]))
				if err != nil {
					log.Printf("Failure indexing batch %s", err)
					continue
				}

				// If the whole request failed, print error and mark all documents as failed
				if res.IsError() {
					if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
						log.Fatalf("Failure to to parse response body: %s", err)
					} else {
						log.Printf("  Error: [%d] %s: %s",
							res.StatusCode,
							raw["error"].(map[string]interface{})["type"],
							raw["error"].(map[string]interface{})["reason"],
						)
					}
				}

				res.Body.Close()
				buf.Reset()

				numItems = 0
			}
		}
	}

}

func (c *ES) toEs7Detail(content chan string)  {
	for {
		select {
		case v := <-content:
			numItems++
			meta := []byte(fmt.Sprintf(`{ "index":{}}%s`, "\n"))

			data := c.detail(v)
			data = append(data, "\n"...) // <-- Comment out to trigger failure for batch

			buf.Grow(len(meta) + len(data))
			buf.Write(meta)
			buf.Write(data)

			// When a threshold is reached, execute the Bulk() request with body from buffer
			if numItems == 2000  {
				res, err := client7.Bulk(bytes.NewReader(buf.Bytes()), client7.Bulk.WithIndex(c.index[0]))
				if err != nil {
					log.Printf("Failure indexing batch %s", err)
					continue
				}

				// If the whole request failed, print error and mark all documents as failed
				if res.IsError() {
					if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
						log.Fatalf("Failure to to parse response body: %s", err)
					} else {
						log.Printf("  Error: [%d] %s: %s",
							res.StatusCode,
							raw["error"].(map[string]interface{})["type"],
							raw["error"].(map[string]interface{})["reason"],
						)
					}
				}

				res.Body.Close()
				buf.Reset()

				numItems = 0
			}
		case <-time.After(time.Second * 1):
			if len(buf.Bytes()) != 0 {
				res, err := client7.Bulk(bytes.NewReader(buf.Bytes()), client7.Bulk.WithIndex(c.index[0]))
				if err != nil {
					log.Printf("Failure indexing batch %s", err)
					continue
				}

				// If the whole request failed, print error and mark all documents as failed
				if res.IsError() {
					if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
						log.Fatalf("Failure to to parse response body: %s", err)
					} else {
						log.Printf("  Error: [%d] %s: %s",
							res.StatusCode,
							raw["error"].(map[string]interface{})["type"],
							raw["error"].(map[string]interface{})["reason"],
						)
					}
				}

				res.Body.Close()
				buf.Reset()

				numItems = 0
			}
		}
	}

}

func (c *ES) toEs6(content chan string)  {
	for {
		select {
		case v := <-content:
			numItems++
			meta := []byte(fmt.Sprintf(`{ "index":{}}%s`, "\n"))

			c.Content.Content = v
			data, err := json.Marshal(c.Content)
			if err != nil {
				log.Fatalf("Cannot encode article  %s",  err)
			}
			data = append(data, "\n"...) // <-- Comment out to trigger failure for batch

			buf.Grow(len(meta) + len(data))
			buf.Write(meta)
			buf.Write(data)

			// When a threshold is reached, execute the Bulk() request with body from buffer
			if numItems == 2000  {
				res, err := client6.Bulk(bytes.NewReader(buf.Bytes()), client6.Bulk.WithIndex(c.index[0]))
				if err != nil {
					log.Printf("Failure indexing batch %s", err)
					continue
				}

				// If the whole request failed, print error and mark all documents as failed
				if res.IsError() {
					if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
						log.Fatalf("Failure to to parse response body: %s", err)
					} else {
						log.Printf("  Error: [%d] %s: %s",
							res.StatusCode,
							raw["error"].(map[string]interface{})["type"],
							raw["error"].(map[string]interface{})["reason"],
						)
					}
				}

				res.Body.Close()
				buf.Reset()

				numItems = 0
			}
		case <-time.After(time.Second * 1):
			if len(buf.Bytes()) != 0 {
				res, err := client6.Bulk(bytes.NewReader(buf.Bytes()), client6.Bulk.WithIndex(c.index[0]))
				if err != nil {
					log.Printf("Failure indexing batch %s", err)
					continue
				}

				// If the whole request failed, print error and mark all documents as failed
				if res.IsError() {
					if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
						log.Fatalf("Failure to to parse response body: %s", err)
					} else {
						log.Printf("  Error: [%d] %s: %s",
							res.StatusCode,
							raw["error"].(map[string]interface{})["type"],
							raw["error"].(map[string]interface{})["reason"],
						)
					}
				}

				res.Body.Close()
				buf.Reset()

				numItems = 0
			}
		}
	}

}

func (c *ES) detail(src string) []byte  {
	Regexp := regexp.MustCompile(c.reg)
	match := Regexp.FindSubmatchIndex([]byte(src))
	return Regexp.Expand(nil, []byte(c.template), []byte(src), match)
}
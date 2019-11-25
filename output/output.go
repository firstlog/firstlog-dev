package output

import (
	"firstlog/initialize"
	"log"
)

var Storage = make(chan string)

type Output struct {
	OutputES *initialize.Elasticsearch
}

func (c *Output) ToEs() {
	var indexName []string
	indexName = append(indexName,c.OutputES.Index)

	toEs, err := NewEs(c.OutputES.Hosts,indexName,c.OutputES.Shards,c.OutputES.Replicas,c.OutputES.Version,
		               c.OutputES.Detail.Enable,c.OutputES.Detail.Regex,c.OutputES.Detail.Template)
	if err != nil {
		log.Println(err)
		return
	}
	toEs.ToEs(Storage)
}

func (c *Output) ToKafka() {

}

func NewOutput(outputES *initialize.Elasticsearch) *Output {
	res := &Output{
		OutputES:outputES,
	}
	return res
}
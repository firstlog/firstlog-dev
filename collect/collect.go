package collect

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Metadata = NewMetadata()

type Collect struct {}

func New() (*Collect,error) {
	res  := &Collect{}

	file, err := os.Open("data/registry")
	if err != nil {
		return res, nil
	}
	result, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return res, nil
	}

	var data1 map[string]data
	err = json.Unmarshal(result,&data1)
	if err != nil {
		log.Println(err)
		return res, nil
	}
	Metadata.metadataChild = data1

	return res, nil
}

func (c *Collect) Run(recursive bool,directory,ignore,match string,ctx context.Context)  {

	watcher, err := NewRecursiveWatcher(recursive,directory,ignore,match)
	if err != nil {
		log.Fatal("\n",directory,err)
		return
	}
	watcher.Run()
	defer watcher.Close()

	for {
		select {
		case file := <-watcher.Files:
			go func(file string,ctx context.Context) {
				Tail(file,ctx)
			}(file,ctx)
		}
	}
}

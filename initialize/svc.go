package initialize

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func init(){
	log.SetFlags(log.Ldate|log.Lshortfile)
}

type svc struct {
	ConfigFile string
}

func (c *svc) InputTask() []Task  {
	var main Main

	config, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = yaml.Unmarshal(config,&main)
	if err != nil {
		log.Println(err)
		return nil
	}
	return main.Input.Tasks
}

func (c *svc) OutputES() *Elasticsearch  {
	var main Main

	config, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = yaml.Unmarshal(config,&main)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &main.Output.Elasticsearch
}

func NewInitSvc(configFile string) *svc {
	res := &svc{
		ConfigFile:configFile,
	}
	return res
}
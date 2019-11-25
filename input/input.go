package input

import (
	"log"
	"context"

	"firstlog/collect"
	"firstlog/initialize"
)

var tasks = make(map[string]initialize.Task)

type Input struct {
	ctx  context.Context
	tasks []initialize.Task
}

func (c *Input) InputStart() {
	cl, err := collect.New()
	if err != nil {
		log.Println(err)
		return
	}
	for _, task := range c.tasks {
		tasks[task.Directory] = task  //存入全局map

		go func(v initialize.Task,ctx context.Context) {
			cl.Run(v.Recursive,v.Directory,v.Ignore,v.Match,ctx)
		}(task,c.ctx)
	}
}

func (c *Input) SaveMetadata2Registry() {
	collect.Metadata.Del(tasks)
	collect.Metadata.Save2Registry()
}

func NewInput(ctx context.Context,tasks []initialize.Task) *Input {
	res := &Input{
		ctx:ctx,
		tasks:tasks,
	}
	return res
}
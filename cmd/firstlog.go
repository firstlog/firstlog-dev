package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"firstlog/collect"
	"firstlog/initialize"
	"firstlog/input"
	"firstlog/output"
)

type cmd struct {
	*input.Input
	*output.Output
}

var (
	ctx, cancel = context.WithCancel(context.Background())
	signalChan	= make(chan os.Signal, 1)

	svc       = initialize.NewInitSvc("conf/firstlog.yaml")
	inputTask = svc.InputTask()
	outputES  = svc.OutputES()
)

func (c *cmd) Run() {
	go func() {
		c.ToEs()
	}()
	c.InputStart()

	go func() {
		for {
			time.Sleep(time.Second * 1)
			c.SaveMetadata2Registry()
		}
	}()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <- signalChan
		fmt.Printf("signal: %v\n", sig)
		switch sig {
		case syscall.SIGINT:
			shutdown(c)
			return
		case syscall.SIGTERM:
			shutdown(c)
			return
		}
	}
}

func shutdown(c *cmd) {
	cancel()
	for {
		time.Sleep(time.Second * 1)
		if collect.GoTailRum != 0 {
			fmt.Println("wait...",collect.GoTailRum)
			continue
		}else if collect.GoTailRum == 0 {
			c.SaveMetadata2Registry()
			fmt.Println("Shutdown to complete")
			break
		}
	}
}

func NewCmd() *cmd {
	Input  := input.NewInput(ctx,inputTask)
	Output := output.NewOutput(outputES)

	res := &cmd{Input, Output,}
	return res
}
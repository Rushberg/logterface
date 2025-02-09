package layouts

import (
	"fmt"
	"logterface/handlers"
)

type PipeLayout struct {
	pipeHandler *handlers.PipeHandler
}

func NewPipeLayout(pipeHandler *handlers.PipeHandler) *PipeLayout {
	return &PipeLayout{
		pipeHandler: pipeHandler,
	}
}

func (ll *PipeLayout) AddHandler(handler handlers.LogHandler) {
}

func (ll *PipeLayout) Print() int {
	logs := ll.pipeHandler.GetValue()
	if len(logs) > 0 {
		fmt.Println(logs)
	}
	return 0
}

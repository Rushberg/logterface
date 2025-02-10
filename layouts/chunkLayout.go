package layouts

import (
	"fmt"
	"logterface/handlers"
	"strings"
)

type ChunkLayout struct {
	handlers []handlers.LogHandler
}

func NewChunkLayout() *ChunkLayout {
	return &ChunkLayout{
		handlers: []handlers.LogHandler{},
	}
}

func (ll *ChunkLayout) AddHandler(handler handlers.LogHandler) {
	ll.handlers = append(ll.handlers, handler)
}

func (ll *ChunkLayout) Print() int {
	totalLines := 0

	for _, handler := range ll.handlers {
		value := handler.GetValue()
		totalLines += strings.Count(value, "\n")
		fmt.Print(value)
	}

	return totalLines
}

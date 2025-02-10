package layouts

import (
	"fmt"
	"logterface/handlers"
	"strings"
)

type Layout interface {
	AddHandler(handler handlers.LogHandler)
	Print() int
}

type LayoutManager struct {
	layouts    []Layout
	piped      bool
	linesCount int
}

func NewLayoutManager() LayoutManager {
	return LayoutManager{
		layouts: []Layout{},
	}
}

func (lm *LayoutManager) AddLayout(layout Layout) {
	lm.layouts = append(lm.layouts, layout)
}

func (lm *LayoutManager) AddPipe(handlerManager *handlers.HandlerManager) {
	if lm.piped {
		return
	}
	pipeHandler := handlers.NewPipeHandler()
	handlerManager.AddHandler(pipeHandler)
	pipeLayout := NewPipeLayout(pipeHandler)
	lm.layouts = append([]Layout{pipeLayout}, lm.layouts...)
	lm.piped = true
}

func (lm *LayoutManager) Print() {
	fmt.Print(strings.Repeat("\033[F\033[K", lm.linesCount))
	lm.linesCount = 0
	for _, layout := range lm.layouts {
		lm.linesCount += layout.Print()
	}
}

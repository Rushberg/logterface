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
	linesCount int
}

func NewLayoutManager() LayoutManager {
	return LayoutManager{
		layouts: []Layout{},
	}
}

func (hm *LayoutManager) AddLayout(handler Layout) {
	hm.layouts = append(hm.layouts, handler)
}

func (lm *LayoutManager) Print() {
	fmt.Print(strings.Repeat("\033[F\033[K", lm.linesCount))
	lm.linesCount = 0
	for _, layout := range lm.layouts {
		lm.linesCount += layout.Print()
	}
}

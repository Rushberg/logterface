package layouts

import (
	"fmt"
	"logterface/handlers"
)

type LineLayout struct {
	handlers []handlers.LogHandler
	width    int
}

func NewLineLayout(width int) *LineLayout {
	return &LineLayout{
		handlers: []handlers.LogHandler{},
		width:    width,
	}
}

func (ll *LineLayout) AddHandler(handler handlers.LogHandler) {
	ll.handlers = append(ll.handlers, handler)
}

func (ll *LineLayout) Print() int {
	totalLines := 0
	currRowChars := 0
	for _, handler := range ll.handlers {
		value := handler.GetValue()
		output := fmt.Sprintf("| %s: %s ", handler.GetName(), value)
		if currRowChars+len(output) > ll.width {
			fmt.Println("|")
			currRowChars = 0
			totalLines += 1
		}
		currRowChars += len(output)
		fmt.Print(output)
	}
	fmt.Println("|")
	totalLines += 1
	return totalLines
}

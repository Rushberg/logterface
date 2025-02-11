package layouts

import (
	"fmt"
	"logterface/handlers"
	"logterface/utils"
	"math"
	"strings"
	"unicode/utf8"
)

type LineLayout struct {
	handlers []handlers.LogHandler
}

func NewLineLayout(width int) *LineLayout {
	return &LineLayout{
		handlers: []handlers.LogHandler{},
	}
}

func (ll *LineLayout) AddHandler(handler handlers.LogHandler) {
	ll.handlers = append(ll.handlers, handler)
}

// "┌"
// "┘"
// "└"
// "┐"
// "─"
// "│"

func (ll *LineLayout) Print() int {
	currLength := 0
	rows := []string{}
	for j, handler := range ll.handlers {
		if j == 0 {
			rows = append(rows, "┌")
		} else {
			rows[0] += "┬"
		}
		value := handler.GetValue()
		maxRowLength := maxRowLength(value)
		borderLength := utils.Max(maxRowLength+2, len(handler.GetName())+4)   // +2 spaces and +2 ─ to be prettier
		topBorderLength := float64(borderLength-len(handler.GetName())-2) / 2 //excluding spaces arond name
		rows[0] += strings.Repeat("─", int(math.Floor(topBorderLength)))
		rows[0] += fmt.Sprintf(" %s ", handler.GetName())
		rows[0] += strings.Repeat("─", int(math.Ceil(topBorderLength)))
		currRows := strings.Split(value, "\n")
		for i, rowContent := range currRows {
			if len(rows) < i+2 {
				rows = append(rows, fmt.Sprintf("%-*s", currLength, ""))
				rows[i+1] += "│"
			} else if len(rows) == i+2 {
				rows[i+1] = strings.TrimRight(rows[i+1], "┘")
				rows[i+1] += "┤"
			}
			rows[i+1] += ""
			indent := float64(borderLength-maxRowLength) / 2
			front := int(math.Floor(indent))
			back := int(math.Ceil(indent))
			rows[i+1] += fmt.Sprintf("%-*s%s%-*s", front, "", rowContent, back, "")
			rows[i+1] += "│"
		}
		if len(rows) < len(currRows)+2 {
			rows = append(rows, fmt.Sprintf("%-*s└", currLength, ""))
		} else if len(rows) > len(currRows)+2 {
			rows[len(currRows)+1] = strings.TrimRight(rows[len(currRows)+1], "│")
			rows[len(currRows)+1] += "├"
		} else {
			rows[len(currRows)+1] = strings.TrimRight(rows[len(currRows)+1], "┘")
			rows[len(currRows)+1] += "┴"
		}
		rows[len(currRows)+1] += fmt.Sprintf("%s┘", strings.Repeat("─", borderLength))
		currLength += borderLength + 1
	}
	rows[0] += "┐"

	// rows[len(rows)-1] += "┘"
	fmt.Println(strings.Join(rows, "\n"))
	return len(rows)
}

func maxRowLength(input string) int {
	// Split the input string into lines
	lines := strings.Split(input, "\n")

	// Initialize maxLength to 0
	maxLength := 0

	// Iterate through each line to find the maximum length
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = utf8.RuneCountInString(line)
		}
	}
	return maxLength
}

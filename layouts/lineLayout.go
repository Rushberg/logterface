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
	var rows []string
	var prevRowCount int
	var currLength int

	for j, handler := range ll.handlers {
		// Build the top border
		if j == 0 {
			rows = append(rows, "┌")
		} else {
			rows[0] += "┬"
		}

		value := handler.GetValue()
		maxRowLength := maxRowLength(value)
		nameLength := len(handler.GetName())
		borderLength := utils.Max(maxRowLength+2, nameLength+4) // +2 spaces and +2 ─ for aesthetics

		// Calculate padding for the handler name
		topBorderPadding := float64(borderLength-nameLength-2) / 2
		rows[0] += strings.Repeat("─", int(math.Floor(topBorderPadding)))
		rows[0] += fmt.Sprintf(" %s ", handler.GetName())
		rows[0] += strings.Repeat("─", int(math.Ceil(topBorderPadding)))

		// Split the value into rows
		currRows := strings.Split(value, "\n")

		// Process each row of the value
		for i, rowContent := range currRows {
			if len(rows) < i+2 {
				rows = append(rows, fmt.Sprintf("%-*s", currLength, ""))
			}

			// Add the left border
			if prevRowCount > 0 && i == prevRowCount {
				rows[i+1] += "┤"
			} else {
				rows[i+1] += "│"
			}

			// Center-align the row content
			indent := float64(borderLength-maxRowLength) / 2
			front := int(math.Floor(indent))
			back := int(math.Ceil(indent))
			rows[i+1] += fmt.Sprintf("%-*s%s%-*s", front, "", rowContent, back, "")
		}

		// Handle the bottom border
		if len(currRows) < prevRowCount {
			rows[prevRowCount+1] += "┘"
			rows[len(currRows)+1] += "├"
		} else if len(currRows) == prevRowCount {
			rows[len(currRows)+1] += "┴"
		} else {
			if len(rows) < len(currRows)+2 {
				rows = append(rows, fmt.Sprintf("%-*s", currLength, ""))
			}
			rows[len(currRows)+1] += "└"
		}

		// Fill in gaps for rows with fewer lines
		for i := len(currRows) + 2; i < len(rows); i++ {
			if i <= prevRowCount {
				rows[i] += "│"
			} else if i > prevRowCount+1 {
				rows[i] += " "
			}
			rows[i] += fmt.Sprintf("%-*s", borderLength, "")
		}

		// Add the bottom border
		rows[len(currRows)+1] += strings.Repeat("─", borderLength)

		prevRowCount = len(currRows)
		currLength += borderLength + 1
	}

	// Finalize the borders
	rows[0] += "┐"
	for i := 1; i <= prevRowCount; i++ {
		rows[i] += "│"
	}
	rows[prevRowCount+1] += "┘"

	// Print the rows
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
			maxLength = utf8.RuneCountInString(utils.StripANSI(line))
		}
	}
	return maxLength
}

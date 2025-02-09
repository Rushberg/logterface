package main

import (
	"bufio"
	"fmt"
	"logterface/handlers"
	"logterface/layouts"
	"os"
	"time"
)

func main() {
	pipeHandler := handlers.NewPipeHandler()
	numHandlers := []handlers.LogHandler{
		handlers.NewNumbersHandler("Min", ".*: (\\d+)", handlers.Min),
		handlers.NewNumbersHandler("Max", ".*: (\\d+)", handlers.Max),
		handlers.NewNumbersHandler("Avg", ".*: (\\d+)", handlers.Avg),
		handlers.NewNumbersHandler("Sum", ".*: (\\d+)", handlers.Sum),
		handlers.NewNumbersHandler("Last", ".*: (\\d+)", handlers.Latest),
	}

	lineLayout := layouts.NewLineLayout(40)
	pipeLayout := layouts.NewPipeLayout(pipeHandler)

	hm := handlers.NewHandlerManager()
	lm := layouts.NewLayoutManager()

	lm.AddLayout(pipeLayout)
	lm.AddLayout(lineLayout)
	hm.AddHandler(pipeHandler)

	for _, nh := range numHandlers {
		hm.AddHandler(nh)
		lineLayout.AddHandler(nh)
	}

	go Printer(&lm)
	// Create a scanner to read from os.Stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Read input line by line
	for scanner.Scan() {
		line := scanner.Text() // Get the current line
		hm.ProcessLog(line)
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}

func Printer(lm *layouts.LayoutManager) {
	for {
		lm.Print()
		time.Sleep(1000 * time.Millisecond)
	}
}

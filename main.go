package main

import (
	"bufio"
	"fmt"
	"logterface/layouts"
	"logterface/utils"
	"os"
	"time"
)

func main() {
	hm, lm, refresh := utils.ParseConfig("./example/config.json")
	go Printer(lm, refresh)
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

func Printer(lm *layouts.LayoutManager, refresh int) {
	for {
		lm.Print()
		time.Sleep(time.Duration(refresh) * time.Millisecond)
	}
}

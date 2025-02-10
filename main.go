package main

import (
	"bufio"
	"fmt"
	"logterface/config"
	"logterface/layouts"
	"os"
	"time"
)

func main() {
	hm, lm, refresh := config.ParseConfig("./example/config.json")
	// hm.ProcessLog("number: 1")
	// hm.ProcessLog("number: 2")
	// hm.ProcessLog("number: 3")
	// hm.ProcessLog("number: 4")
	// hm.ProcessLog("number: 5")
	// hm.ProcessLog("number: 6")
	// hm.ProcessLog("number: 5")
	// hm.ProcessLog("number: 4")
	// hm.ProcessLog("number: 3")
	// hm.ProcessLog("number: 3")
	// hm.ProcessLog("number: 5")
	// hm.ProcessLog("number: 6")
	// lm.Print()
	// return

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

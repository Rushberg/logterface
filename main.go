package main

import (
	"bufio"
	"fmt"
	"logterface/config"
	"logterface/layouts"
	"os"
	"strings"
	"time"
)

func main() {
	hm, lm, refresh := config.ParseConfig(os.Args[1])

	go Printer(lm, refresh)
	// Create a scanner to read from os.Stdin
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	for {
		line, err := reader.ReadString('\n') // Read input until newline
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}
		hm.ProcessLog(strings.TrimSuffix(line, "\n"))
	}
	// print last batch
	time.Sleep(time.Duration(refresh+5000) * time.Millisecond)
}

func Printer(lm *layouts.LayoutManager, refresh int) {
	for {
		lm.Print()
		time.Sleep(time.Duration(refresh) * time.Millisecond)
	}
}

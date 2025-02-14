package main

import (
	"fmt"
	"logterface/utils"
	"math/rand"
	"time"
)

func main() {
	progress := 0
	usage := 50
	fmt.Println("target: 1000")
	for progress < 1001 {
		switch rand.Intn(5) {
		case 1:
			fmt.Printf("random number: %d\n", rand.Intn(20000))
		case 2:
			fmt.Printf("progress: %d\n", progress%1001)
			progress++
		case 3:
			fmt.Printf("error: Uh oh! Something bad happened. Code:%d\n", rand.Intn(10))
			progress++
		case 4:
			usage += -5 + rand.Intn(10)
			usage = utils.Max(10, utils.Min(100, usage))
			fmt.Printf("resources usage: %d/100\n", usage)
		}

		time.Sleep(50 * time.Millisecond)
	}
}

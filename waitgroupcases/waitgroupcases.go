package waitgroupcases

import (
	"fmt"
	"sync"
)

var words []string = []string{
	"alpha",
	"beta",
	"delta",
}

func WaitgroupRun() {
	var wg sync.WaitGroup
	wg.Add(len(words))

	for i, x := range words {
		go printString(fmt.Sprintf("%d: %s \n", i, x), &wg)
	}

	wg.Wait()

	wg.Add(1)
	printString("This function has finished", &wg)
}

func printString(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf(s)
}

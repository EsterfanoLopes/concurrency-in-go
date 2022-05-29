package mutex

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func MutexRun() {
	msg = "Hello, world!"

	wg.Add(2)
	updateMessage("Hello, universe!")
	updateMessage("Hello, cosmos!")
	wg.Wait()

	fmt.Println(msg)
}

func updateMessage(s string) {
	defer wg.Done()

	msg = s
}

// func MutexRun() {
// 	msg = "Hello, world!"

// 	var mutex sync.Mutex

// 	wg.Add(2)
// 	updateMessage("Hello, universe!", &mutex)
// 	updateMessage("Hello, cosmos!", &mutex)
// 	wg.Wait()

// 	fmt.Println(msg)
// }

// func updateMessage(s string, m *sync.Mutex) {
// 	defer wg.Done()

// 	m.Lock()
// 	msg = s
// 	m.Unlock()
// }

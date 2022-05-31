package dining_philosophers

import (
	"fmt"
	"sync"
	"time"
)

var (
	philosophers = []string{"Plato", "Socrates", "Aristotle", "Pascal", "Locke"}
	wg           sync.WaitGroup
	sleepTime    = 1 * time.Second
	eatTime      = 2 * time.Second
	thinkTime    = 1 * time.Second
)

const (
	hunger = 3
)

func Run() {
	// print intro
	fmt.Println("The Dining Philosohpers Problem")
	fmt.Println("-------------------------------")
	wg.Add(len(philosophers))

	forkLeft := &sync.Mutex{}

	// spawn go routine for each philosopher
	for i := 0; i < len(philosophers); i++ {
		// create a mutex for the right fork
		forkRight := &sync.Mutex{}
		//	call a goroutine
		go diningProblem(philosophers[i], forkLeft, forkRight)

		forkLeft = forkRight
	}

	wg.Wait()

	fmt.Println("The table is empty.")
}

func diningProblem(philosopher string, leftFork, rightFork *sync.Mutex) {
	defer wg.Done()

	// print a message
	fmt.Println(philosopher, "is seated.")
	time.Sleep(sleepTime)

	for i := hunger; i > 0; i-- {
		fmt.Println(philosopher, "is hungry.")
		time.Sleep(sleepTime)

		// lock both forks
		leftFork.Lock()
		fmt.Printf("\t%s picked up the fork to his left.\n", philosopher)
		rightFork.Lock()
		fmt.Printf("\t%s picked up the fork to his right.\n", philosopher)

		// philospher has both forks
		fmt.Println(philosopher, "has both forks and is eating.")
		time.Sleep(eatTime)

		// give the philosohper a time to think
		fmt.Println(philosopher, "is thinking.")
		time.Sleep(thinkTime)

		// unlock the mutexes
		rightFork.Unlock()
		fmt.Printf("\t%s put down the fork on his right\n", philosopher)
		leftFork.Unlock()
		fmt.Printf("\t%s put down the fork on his left\n", philosopher)
		time.Sleep(sleepTime)
	}

	// print out message
	fmt.Println(philosopher, "is satisfied.")
	time.Sleep(sleepTime)
	fmt.Println(philosopher, "has left the table.")
}

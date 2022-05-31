package dining_philosophers

import "sync"

// variables
// philosophers
var philosophers = []string{"Plato", "Socrates", "Aristotle", "Pascal", "Locke"}
var wg sync.WaitGroup

func Run() {
	// print intro

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

	return
}

func diningProblem(philosopher string, dominantHand, otherHand *sync.Mutex) {
	defer wg.Done()

	// print a message

	// lock both forks

	// philospher has both forks

	// unlock the mutexes
}

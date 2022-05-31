package dining_philosophers

import "sync"

// variables
// philosophers
var philosophers = []string{"Plato", "Socrates", "Aristotle", "Pascal", "Locke"}
var wg sync.WaitGroup

func Run() {
	// print intro

	wg.Add(len(philosophers))

	// spawn go routine for each philosopher
	for i := 0; i < len(philosophers); i++ {
		//	call a goroutine
		go diningProblem()
	}

	wg.Wait()

	return
}

func diningProblem() {
	defer wg.Done()
}

package ongoing_worker_pool

import (
	"concurrency-in-go/ongoing_worker_pool/event"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type job struct{}

func New() job {
	return job{}
}

func (j *job) Run(ctx context.Context) {
	// Number of workers in the pool
	numWorkers := 5

	// Create a channel to receive events to process (events)
	eventsChan := make(chan event.Event)

	// Section 1: Events generation
	// Can be a queue reading, a file reading, a database reading, etc.
	go generateEvents(eventsChan)

	// Section 2: Workers
	// Initiate the workers who will process the events.
	// the number of workers determine the number of concurrent process running.

	// Create a channel to receive events to execute (jobs)
	jobsChan := make(chan string)

	// Wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Start the workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobsChan, &wg)
	}

	// Section 3: Event processing
	// With workers ready to process the events, it is only necessary to send the events to the workers through the jobs channel.

	// generate infinite events
	go generateEvents(eventsChan)

	// Send events to workers to process
	go func() {
		for e := range eventsChan {
			// Parse the event to the worker expectation
			jobsChan <- e.Data
		}
	}()

	// Wait for them to finish

	// Wait for all workers to finish processing
	wg.Wait()

	// Close the jobs channel when done (this will exit the worker goroutines)
	close(jobsChan)
	close(eventsChan)
}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Process the event (job) here
		fmt.Printf("Worker %d processing: message %s\n", id, job)
		random := rand.Intn(5)
		time.Sleep(time.Duration(random) * time.Second)
	}
}

func generateEvents(sendEvents chan event.Event) {
	// infinite loop to simulate an infinite number of events
	for {
		i := 0

		random := rand.Intn(999)
		eventID := uint(random + i)
		e := event.New(eventID, fmt.Sprintf("Event number %d", eventID), uint(i))

		sendEvents <- e
	}
}

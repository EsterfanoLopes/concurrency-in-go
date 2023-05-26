package channel_wait_worker_pool

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Event struct {
	ID    int
	Data  string
	Cycle int
}

type Job struct {
	numberOfWorkers        int
	numberOfEventsPerCycle int
	waitBeforeNextCycle    time.Duration
	shutdownChannel        chan any
}

func (j *Job) gracefulShutdown(chan any) {
	<-j.shutdownChannel
	// TODO: Control the close of channels on the shutdown premisse
}

func New(numberOfWorkers, numberOfEventsPerCycle int, waitBeforeNextCycle time.Duration) Job {
	return Job{
		numberOfWorkers:        numberOfWorkers,
		numberOfEventsPerCycle: numberOfEventsPerCycle,
		waitBeforeNextCycle:    waitBeforeNextCycle,
	}
}

func (j *Job) Run(ctx context.Context) {
	// Create channels for events and control
	events := make(chan Event)
	control := make(chan struct{})
	shutdown := make(chan bool)

	// Notify external context shutdown
	go func() {
		for {
			a := ctx.Done()
			if a != nil {
				fmt.Printf("Context done received\n")
				shutdown <- true
				return
			}
		}
	}()

	go func() {
		select {
		case <-shutdown:
			fmt.Println("Shutdown received. Closing all channels")
			j.shouldClose = true
		}
	}()

	cycle := 0
	for {
		fmt.Printf("Starting cycle %d\n", cycle)
		// Create the worker pool
		for i := 0; i < j.numberOfWorkers; i++ {
			go worker(i, events, control)
		}

		// Generate events and send them to the channel
		go func() {
			for i := 1; i <= j.numberOfEventsPerCycle; i++ {
				event := Event{
					ID:    i,
					Data:  fmt.Sprintf("Event %d", i),
					Cycle: cycle,
				}
				events <- event
			}
		}()

		// Wait for all events to be processed
		for i := 0; i < j.numberOfEventsPerCycle; i++ {
			<-control
		}

		if j.shouldClose {
			break
		}

		fmt.Printf("Finished cycle %d\n", cycle)
		cycle++
	}

	close(events)

	fmt.Println("All events processed")
}

func worker(id int, events <-chan Event, control chan<- struct{}) {
	for event := range events {
		processEvent(id, event)
		control <- struct{}{}
	}
}

func processEvent(workerID int, event Event) {
	fmt.Printf("Worker %d processing event %d from cycle %d: %s\n", workerID, event.ID, event.Cycle, event.Data)
	// Simulate some processing time
	time.Sleep(time.Duration(rand.Intn(9)) * time.Second)
	fmt.Printf("Worker %d finished processing event %d from cycle %d\n", workerID, event.ID, event.Cycle)
}

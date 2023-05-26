package channel_wait_worker_pool

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Event struct {
	ID   int
	Data string
}

type Job struct {
	numberOfWorkers        int
	numberOfEventsPerCycle int
	waitBeforeNextCycle    time.Duration
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
	shutdown := make(chan bool) // TODO: Use this channel to notify graceful shutdown

	// Notify external context shutdown
	// TODO: Create listener to close all channels after shutdown is received
	go func() {
		for {
			<-ctx.Done()
			shutdown <- true
			return
		}
	}()

	// TODO: Create infinite loop to control the cycles

	// Create the worker pool
	for i := 0; i < j.numberOfWorkers; i++ {
		go worker(i, events, control)
	}

	// Generate events and send them to the channel
	go func() {
		for i := 1; i <= j.numberOfEventsPerCycle; i++ {
			event := Event{
				ID:   i,
				Data: fmt.Sprintf("Event %d", i),
			}
			events <- event
		}
		close(events)
	}()

	// Wait for all events to be processed
	for i := 0; i < j.numberOfEventsPerCycle; i++ {
		<-control
	}

	fmt.Println("All events processed")
}

func worker(id int, events <-chan Event, control chan<- struct{}) {
	for event := range events {
		processEvent(id, event)
		control <- struct{}{}
	}
}

func processEvent(workerID int, event Event) {
	fmt.Printf("Worker %d processing event %d: %s\n", workerID, event.ID, event.Data)
	// Simulate some processing time
	time.Sleep(time.Duration(rand.Intn(9)) * time.Second)
	fmt.Printf("Worker %d finished processing event %d\n", workerID, event.ID)
}

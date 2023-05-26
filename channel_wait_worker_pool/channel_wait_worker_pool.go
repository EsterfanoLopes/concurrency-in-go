package channel_wait_worker_pool

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Event represents the data to be processed
type Event struct {
	ID    int    // event ID
	Data  string // data to be processed
	Cycle uint   // cycle number where this event was generated
}

// Job represents the setup for the job and variables who modifies its behavior
type Job struct {
	numberOfWorkers        int           // size of the worker pool
	numberOfEventsPerCycle int           // number of events to be processed per cycle
	waitBeforeNextCycle    time.Duration // time to wait before starting the next cycle
	shutdownChannel        chan bool     // channel to notify shutdown to all goroutines
	stopCycle              bool          // flag to stop the cycle for loop
}

// gracefulShutdown controls the shutdown of the job and its related channels from multiple sources
// - context: external parameter controled by the job caller
// - channel: channel to receive any shutdown signal, from inside or outside the job. E.g.: Sigterm from OS defined on main.go
func (j *Job) gracefulShutdown(ctx context.Context, channelsToClose []any) {
	go func() { // goroutine to listen shutdown coming from the context
		deadlineReceived := ctx.Done()
		for {
			select {
			case <-deadlineReceived:
				fmt.Println("Shutdown received through context")
				j.stopCycle = true
				break
			}
		}
	}()

	go func() { // goroutine to listen shutdown coming from the channel
		<-j.shutdownChannel
		fmt.Println("Shutdown received through channel")
		fmt.Println("gracefully shutting down...")

		j.closeChannels(channelsToClose)
	}()
}

func (j *Job) closeChannels(channelsToClose []any) {
	fmt.Println("Closing Shutdown channel")
	close(j.shutdownChannel)

	for _, channel := range channelsToClose {
		switch channel.(type) {
		case chan Event:
			fmt.Println("Closing Event channel")
			c := channel.(chan Event)
			close(c)
		case chan struct{}:
			fmt.Println("Closing Control channel")
			c := channel.(chan struct{})
			close(c)
		default:
			fmt.Println("Channel type not supported")
			// FIXME: What to do with unexpected channel type? Panic? external channel to errors?
		}
	}
}

func New(numberOfWorkers, numberOfEventsPerCycle int, waitBeforeNextCycle time.Duration, shutdownChannel chan bool) Job {
	return Job{
		numberOfWorkers:        numberOfWorkers,
		numberOfEventsPerCycle: numberOfEventsPerCycle,
		waitBeforeNextCycle:    waitBeforeNextCycle,
		shutdownChannel:        shutdownChannel,
	}
}

func (j *Job) Run(ctx context.Context) {
	// Create channels for events and control
	events := make(chan Event)     // event channel controls the received events to be processed
	control := make(chan struct{}) // control channel controls if the workers are ready to receive a new event to process or not

	go j.gracefulShutdown(ctx, []any{events, control})

	// cycle controls in which cycle the job is
	var cycle uint
	for {
		fmt.Printf("Starting cycle %d\n", cycle)
		// Instantiate the worker pool
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

		if j.stopCycle {
			fmt.Println("Stopping cycle run")
			j.shutdownChannel <- true
			break
		}

		fmt.Printf("Finished cycle %d\n", cycle)
		cycle++
		// This wait is more relevant in a production environment, to avoid overloading the number of requests to the empty queue
		// It can be enhanced to instead of having a waiting time, to wait only if the last fetch for messages was empty
		time.Sleep(j.waitBeforeNextCycle)
	}

	fmt.Println("All events processed")
}

// worker controls the distribution of events to be processed by a available worker
func worker(id int, events <-chan Event, control chan<- struct{}) {
	for event := range events {
		processEvent(id, event)
		control <- struct{}{} // this channel controls if the workers are ready to receive a new event to process or not
	}
}

// processEvent simulates the processing of an event. This would contain the real processing for the job
func processEvent(workerID int, event Event) {
	fmt.Printf("Worker %d processing event %d from cycle %d: %s\n", workerID, event.ID, event.Cycle, event.Data)
	// Simulate some random processing time
	time.Sleep(time.Duration(rand.Intn(9)) * time.Second)
	fmt.Printf("Worker %d finished processing event %d from cycle %d\n", workerID, event.ID, event.Cycle)
}

package channel_wait_worker_pool

import (
	"concurrency-in-go/channel_wait_worker_pool/event"
	"concurrency-in-go/channel_wait_worker_pool/mutex_map"
	"context"
	"fmt"
	"time"
)

// job represents the setup for the job and variables who modifies its behavior
type job struct {
	numberOfWorkers        int                 // size of the worker pool
	numberOfEventsPerCycle uint                // number of events to be processed per cycle
	waitBeforeNextCycle    time.Duration       // time to wait before starting the next cycle
	controlIDMutexMap      *mutex_map.MutexMap // map to store the IDs being operated and the mutex for each worker
	shutdownChannel        chan bool           // channel to notify shutdown to all goroutines
	stopCycle              bool                // flag to stop the cycle for loop
}

func New(numberOfWorkers int, numberOfEventsPerCycle uint, waitBeforeNextCycle time.Duration, shutdownChannel chan bool) job {
	controlIDMutexMap := mutex_map.New()

	return job{
		numberOfWorkers:        numberOfWorkers,
		numberOfEventsPerCycle: numberOfEventsPerCycle,
		waitBeforeNextCycle:    waitBeforeNextCycle,
		shutdownChannel:        shutdownChannel,
		controlIDMutexMap:      controlIDMutexMap,
	}
}

func (j *job) Run(ctx context.Context) {
	// Create channels for events and control
	events := make(chan event.Event) // event channel controls the received events to be processed
	control := make(chan uint)       // control channel controls if the workers are ready to receive a new event to process or not

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
			for i := uint(1); i <= j.numberOfEventsPerCycle; i++ {
				event := event.New(i, fmt.Sprintf("Event %d", i), cycle)

				_, ok := j.controlIDMutexMap.Get(event.ID)
				if ok {
					fmt.Printf("Event %d with message %s already exists in the map. Retry it next cycle\n", event.ID, event.Data)
					i-- // to have another event on the same cycle
					continue
				}

				// Set id to the control mutex
				j.controlIDMutexMap.Set(event.ID, struct{}{})

				// Send event to the channel to be processed
				events <- event
			}
		}()

		// Wait for all events to be processed
		for i := uint(0); i < j.numberOfEventsPerCycle; i++ {
			idToRelease := <-control
			j.controlIDMutexMap.Delete(idToRelease) // release the mutex for the worker to receive a new event
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
func worker(id int, events <-chan event.Event, control chan<- uint) {
	for event := range events {
		event.ProcessEvent(id)
		control <- event.ID // this channel controls if the workers are ready to receive a new event to process or not
	}
}

// gracefulShutdown controls the shutdown of the job and its related channels from multiple sources
// - context: external parameter controled by the job caller
// - channel: channel to receive any shutdown signal, from inside or outside the job. E.g.: Sigterm from OS defined on main.go
func (j *job) gracefulShutdown(ctx context.Context, channelsToClose []any) {
	go func() { // goroutine to listen shutdown coming from the context
		deadlineReceived := ctx.Done()

		<-deadlineReceived
		fmt.Println("Shutdown received through context")
		j.stopCycle = true
	}()

	go func() { // goroutine to listen shutdown coming from the channel
		<-j.shutdownChannel
		fmt.Println("Shutdown received through channel")
		fmt.Println("gracefully shutting down...")

		j.closeChannels(channelsToClose)
	}()
}

func (j *job) closeChannels(channelsToClose []any) {
	fmt.Println("Closing Shutdown channel")
	close(j.shutdownChannel)

	for _, channel := range channelsToClose {
		switch channel.(type) {
		case chan event.Event:
			fmt.Println("Closing Event channel")
			c := channel.(chan event.Event)
			close(c)
		case chan uint:
			fmt.Println("Closing Control channel")
			c := channel.(chan uint)
			close(c)
		default:
			fmt.Println("Channel type not supported")
			// FIXME: What to do with unexpected channel type? Panic? external channel to errors?
		}
	}
}

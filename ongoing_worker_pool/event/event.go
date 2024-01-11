package event

import (
	"fmt"
	"math/rand"
	"time"
)

// Event represents the data to be processed
type Event struct {
	ID   uint   // event ID
	Data string // data to be processed
}

// New creates a new Event
func New(ID uint, Data string, Cycle uint) Event {
	return Event{
		ID:   ID,
		Data: Data,
	}
}

// ProcessEvent simulates the processing of an event. This would contain the real processing for the job
func (e *Event) ProcessEvent(workerID int) {
	fmt.Printf("Worker %d processing event %d from: %s\n", workerID, e.ID, e.Data)
	// Simulate some random processing time
	time.Sleep(time.Duration(rand.Intn(9)) * time.Second)
	fmt.Printf("Worker %d finished processing event %d from cycle %d\n", workerID, e.ID)
}

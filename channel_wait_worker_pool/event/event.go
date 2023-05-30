package event

import (
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

// New creates a new Event
func New(ID int, Data string, Cycle uint) Event {
	return Event{
		ID:    ID,
		Data:  Data,
		Cycle: Cycle,
	}
}

// ProcessEvent simulates the processing of an event. This would contain the real processing for the job
func (e *Event) ProcessEvent(workerID int) {
	fmt.Printf("Worker %d processing event %d from cycle %d: %s\n", workerID, e.ID, e.Cycle, e.Data)
	// Simulate some random processing time
	time.Sleep(time.Duration(rand.Intn(9)) * time.Second)
	fmt.Printf("Worker %d finished processing event %d from cycle %d\n", workerID, e.ID, e.Cycle)
}

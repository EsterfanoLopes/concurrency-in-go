package dining_philosophers

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	sleepTime = 0 * time.Second
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second

	for i := 0; i < 1000; i++ {
		Run()
		if len(finishingOrder) != 5 {
			t.Error("wrong numbe of entries in slice")
		}

		finishingOrder = []string{}
	}
}

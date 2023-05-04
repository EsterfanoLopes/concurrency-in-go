package main

<<<<<<< Updated upstream
import (
	"concurrency-in-go/dining_philosophers"
)
=======
import "concurrency-in-go/producer_consumer"
>>>>>>> Stashed changes

func main() {
	// waitgroupcases.WaitgroupRun()

	// waitgroupcases.Challenge()

	producer_consumer.ProducerConsumerRun()

<<<<<<< Updated upstream
	dining_philosophers.Run()
=======
	// dining_philosophers.Run()

	// simple_channels.Run()

	// channel_select.Run()

	// buffered_channels.Run()

	// sleeping_barber.Run()
>>>>>>> Stashed changes
}

package main

import (
	"concurrency-in-go/channel_wait_worker_pool"
	"context"
	"time"
)

func main() {
	// waitgroupcases.WaitgroupRun()

	// waitgroupcases.Challenge()

	// producer_consumer.ProducerConsumerRun()

	// dining_philosophers.Run()

	var (
		numberOfWorkers        = 10
		numberOfEventsPerCicle = 50
		waitUntilNextCicle     = 5 * time.Second

		ctxWithCancel, cancel = context.WithCancel(context.Background())
	)

	// new job with settings
	job := channel_wait_worker_pool.New(
		numberOfWorkers,
		numberOfEventsPerCicle,
		waitUntilNextCicle,
	)

	job.Run(ctxWithCancel)

	// default time to call cancel to the context
	time.Sleep(1 * time.Second)
	cancel()
}

package channel_select

import (
	"fmt"
	"time"
)

func Run() {
	fmt.Println("Select with channels")
	fmt.Println("--------------------")

	channelOne := make(chan string)
	channelTwo := make(chan string)

	go serverOne(channelOne)
	go serverTwo(channelTwo)

	for {
		select {
		case s1 := <-channelOne:
			fmt.Println("Case one: ", s1)
		case s2 := <-channelOne:
			fmt.Println("Case two: ", s2)
			// If there's more ten one possible receiver, it will be delivered in one of those randomly
		case s3 := <-channelTwo:
			fmt.Println("Case three: ", s3)
		case s4 := <-channelTwo:
			fmt.Println("Case three: ", s4)
			// default:
			// avoid deadlock
		}
	}
}

func serverOne(ch chan string) {
	for {
		time.Sleep(time.Second * 6)
		ch <- "This is from server 1"
	}
}

func serverTwo(ch chan string) {
	for {
		time.Sleep(time.Second * 3)
		ch <- "This is from server 2"
	}
}

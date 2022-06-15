package simple_channels

import (
	"fmt"
	"strings"
)

func Run() {
	// create two channels
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and press ENTER (enter Q to quit)")

	for {
		// print a prompt
		fmt.Print("-> ")

		// get user input
		var userInput string
		_, _ = fmt.Scanln(&userInput)
		if userInput == strings.ToLower("q") {
			break
		}

		ping <- userInput
		// wait for a response
		response := <-pong
		fmt.Println("Response: ", response)
	}

	fmt.Println("All done. Closing channels.")
	close(ping)
	close(pong)
}

func shout(
	ping <-chan string, // receive only
	pong chan<- string, // send only
) {
	for {
		value := <-ping
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(value))
	}
}

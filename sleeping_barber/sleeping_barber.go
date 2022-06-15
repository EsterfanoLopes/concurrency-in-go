package sleeping_barber

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// variables
var (
	seatingCapacity = 10
	arrivalRate     = 100 // miliseconds
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10 * time.Second
)

func Run() {
	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// print welcome message
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	// create channels
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create data structure for the barbershop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day!")

	// add barbers
	shop.addBarber("John Doe")

	// start the barbershop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients

	// block until the barbershop is closed

	time.Sleep(5 * time.Second)
}

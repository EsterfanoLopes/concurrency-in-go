package sleeping_barber

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (bs *BarberShop) addBarber(barberName string) {
	bs.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients", barberName)

		for {
			// if there are no clients, the barber goes to sleep
			if len(bs.ClientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barberName)
				isSleeping = true
			}

			client, shopOpen := <-bs.ClientsChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barberName)
					isSleeping = false
				}

				// cut hair
				bs.cutHair(barberName, client)
			} else {
				// shop is closed. Send the barber home and close this goroutine
				bs.sendBarberHome(barberName)
				return
			}
		}
	}()
}

func (bs *BarberShop) cutHair(barberName, client string) {
	color.Green("%s is cutting %s's hair.", barberName, client)
	time.Sleep(bs.HairCutDuration)
	color.Green("%s is finished cutting %s hair", barberName, client)
}

func (bs *BarberShop) sendBarberHome(barberName string) {
	color.Cyan("%s is going home.", barberName)
	bs.BarbersDoneChan <- true
}

func (bs *BarberShop) closeShopForDay() {
	color.Cyan("Closing shop for the day.")

	close(bs.ClientsChan)
	bs.Open = false

	for a := 1; a <= bs.NumberOfBarbers; a++ {
		<-bs.BarbersDoneChan
	}

	close(bs.BarbersDoneChan)
	color.Green("-------------------------------------------------------------------")
	color.Green("The barbershop is now closed for th day and everyone has gone home.")
}

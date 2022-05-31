package mutex

import (
	"fmt"
	"sync"
)

var wgcm sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func ComplexMutexRun() {
	// variable for bank balance
	var bankBalance int
	var balance sync.Mutex

	// print out starting values

	fmt.Printf("Initial account balance: $%d.00", bankBalance)
	fmt.Println()

	// define weekly revenuw
	incomes := []Income{
		{Source: "Main job", Amount: 500},
		{Source: "Gifts", Amount: 10},
		{Source: "Part time job", Amount: 50},
		{Source: "Investments", Amount: 100},
	}

	wgcm.Add(len(incomes))

	// loop trhough 52  weeks and print out how much is made; keep a running total
	for i, income := range incomes {
		go func(i int, income Income) {
			defer wgcm.Done()
			for week := 1; week <= 52; week++ {
				balance.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				balance.Unlock()

				fmt.Printf("On week %d, you earned $%d.00 from %s\n", week, income.Amount, income.Source)
			}
		}(i, income)
	}

	wgcm.Wait()

	// print out final balance
	fmt.Printf("Final bank balace: %d.00", bankBalance)
	fmt.Println()
}

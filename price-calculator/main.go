package main

import (
	"fmt"

	"example.com/price-calculator/cmdmanager"
	"example.com/price-calculator/prices"
)

type IOManager interface {
	ReadLines() ([]string, error)
	WriteResults(data interface{}) error
}

func main() {
	var taxRates []float64 = []float64{0, 0.07, 0.1, 0.15}

	for _, taxRate := range taxRates {
		// fm := filemanager.New("prices.txt", fmt.Sprintf("prices_%.0f.json", taxRate*100))
		cmdm := cmdmanager.New()
		pricesJob := prices.NewTaxIncludedPriceJob(cmdm, taxRate)
		err := pricesJob.Process()
		if err != nil {
			fmt.Println("Could not precess prices ", err)
		}
	}
}

package main

import (
	"fmt"
	"math"
)

const inflation float64 = 2.5

func main() {

	var investmentAmount, years float64
	var expectedReturnRate float64

	var revenue, expenses, taxRate float64

	valuesInput(&investmentAmount, &years, &expectedReturnRate, &revenue, &expenses)

	futureValue := investmentAmount * math.Pow((1+expectedReturnRate/100), float64(years))
	futureRealValue := futureValue / math.Pow(1+inflation/100, years)

	fmt.Printf("Future value: %f \n", futureValue)
	fmt.Printf("Future real value: %f", futureRealValue)

	ebt, profit := calculateEbtAndProfit(revenue, expenses, taxRate)
	netProfit := ebt - profit
	ratio := ebt / profit
	fmt.Printf("\nEBT: %.2f \nTax: %.2f \nNet profit: %.2f \nRatio: %.2f", ebt, profit, netProfit, ratio)
}

func calculateEbtAndProfit(revenue, expenses, taxRate float64) (float64, float64) {
	return revenue - expenses, revenue - expenses*(1-taxRate/100)
}

func valuesInput(investmentAmount *float64, years *float64, expectedReturnRate *float64, revenue *float64, expenses *float64) {
	fmt.Print("Enter investment amount: ")
	fmt.Scan(investmentAmount)
	fmt.Print("Enter years: ")
	fmt.Scan(years)
	fmt.Print("Enter expected return rate: ")
	fmt.Scan(expectedReturnRate)
	fmt.Print("Enter expected revenue: ")
	fmt.Scan(revenue)
	fmt.Print("Enter expected expenses: ")
	fmt.Scan(expenses)
	fmt.Print("Enter expected tax rate: ")
	fmt.Scan(expenses)
}

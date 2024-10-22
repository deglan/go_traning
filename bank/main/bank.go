package main

import (
	"fmt"

	"example.com/bank/fileutils"
	"github.com/Pallinder/go-randomdata"
)

const accountBalanceFile = "balance.txt"

func main() {
	var isExit = false
	var balance, err = fileutils.GetFloatFromFile(accountBalanceFile)

	if err != nil {
		fmt.Println("ERROR:")
		fmt.Println(err)
		fmt.Println("-----------------------")
		panic("Can't continue")
	}

	fmt.Println("Welcome to Go Bank!")
	fmt.Println("Reach us 27/7", randomdata.PhoneNumber())

	for !isExit {
		presentOptions()

		var choice int
		fmt.Scan(&choice)

		fmt.Println("Your choice: ", choice)
		if choice == 1 {
			fmt.Println("Your balance is ", balance)
		} else if choice == 2 {
			fmt.Print("Your deposit: ")
			var depositAmount float64
			fmt.Scan(&depositAmount)
			if depositAmount <= 0 {
				fmt.Print("Deposit smaller then 0. Invalid input")
				return
			}
			balance += depositAmount
			fmt.Println("Your current balance: ", balance)
			fileutils.WriteToFile(balance, accountBalanceFile)
		} else if choice == 3 {
			fmt.Print("How much do you want to withdraw: ")
			var withdraw float64
			fmt.Scan(&withdraw)
			if withdraw <= 0 {
				fmt.Print("withdraw smaller then 0. Invalid input")
				continue
			}
			if withdraw > balance {
				fmt.Print("Invalid input. Withdraw greater then balance")
				continue
			}
			balance -= withdraw
			fmt.Println("Your current balance: ", balance)
			fileutils.WriteToFile(balance, accountBalanceFile)
		} else if choice == 4 {
			fmt.Println("Goodbye!")
			isExit = true
		} else {
			fmt.Println("Invalid input")
		}
	}
}

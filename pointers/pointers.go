package main

import "fmt"

func main() {
	age := 32

	agePointer := &age
	fmt.Println(age)
	fmt.Println(*agePointer)
	fmt.Println(getAdultAge(agePointer))
}

func getAdultAge(age *int) int {
	return *age - 18
}

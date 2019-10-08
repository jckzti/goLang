package main

import (
	"fmt"
)

func main() {
	x := 5

	if x > 6 {
		fmt.Println("Print more than 6")
	} else if x < 2 {

	} else {
		a := [5]int{5, 4, 3, 2, 1}

		b := []int{5, 4, 3, 2, 1}
		b = append(b, 13)
		a[2] = 7
		fmt.Println(a[1] + b[5])
	}
}

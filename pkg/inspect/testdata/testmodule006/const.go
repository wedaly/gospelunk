package main

import "fmt"

const intConst = 123

func main() {
	var x *int
	if x == nil {
		fmt.Printf("x is nil\n")
	}

	fmt.Printf("%v\n", intConst)
}

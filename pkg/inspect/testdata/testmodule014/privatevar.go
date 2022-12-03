package testmodule014

import "fmt"

var privateVar int

func funcUsingPrivateVar() {
	fmt.Printf("%d", privateVar)
}

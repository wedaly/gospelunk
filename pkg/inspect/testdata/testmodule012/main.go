package main

import "fmt"

func main() {
	s := FirstStruct{A: "foo"}.ReturnSecondStruct().B
	fmt.Printf("%s", s)
}

type FirstStruct struct {
	A string
}

type SecondStruct struct {
	B string
}

func (f FirstStruct) ReturnSecondStruct() SecondStruct {
	return SecondStruct{B: f.A}
}

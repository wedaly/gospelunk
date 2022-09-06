package main

import "fmt"

type TestStruct struct {
	StringField string
	IntField    int
}

func main() {
	s := TestStruct{
		StringField: "foo",
		IntField:    123,
	}
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.IntField)
}

package main

import "fmt"

type TestInterface interface {
	fmt.Stringer
	MyString() string
}

type TestStruct struct{}

func (s *TestStruct) String() string {
	return "Hello world!"
}

func (s *TestStruct) Sub() *TestSubStruct {
	return &TestSubStruct{}
}

type TestSubStruct struct{}

func (ss *TestSubStruct) String() string {
	return "Hello from sub struct"
}

func main() {
	s := TestStruct{}
	fmt.Printf("%s\n", s.String())
	fmt.Printf("%s\n", s.Sub().String())

	var iface TestInterface
	iface = &s
	fmt.Printf("%s\n", iface.MyString())
	fmt.Printf("%s\n", iface.String())
}

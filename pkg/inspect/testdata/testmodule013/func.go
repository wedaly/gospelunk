package testmodule013

import "fmt"

type MyArg string

type MyReturn int

type MyFunc func(s MyArg) MyReturn

type MyStructWithFunc struct {
	FuncField MyFunc
}

func ConstructStructWithFuncField() {
	x := MyStructWithFunc{
		FuncField: func(s MyArg) MyReturn {
			return MyReturn(len(s))
		},
	}
	fmt.Printf("%v", x)
}

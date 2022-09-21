package testmodule002

// #include <stdio.h>
import "C"

type MyStruct struct{}

func Random() int {
	return int(C.random())
}

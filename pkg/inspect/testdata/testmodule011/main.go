package main

import (
	"github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule011/subpkg"
)

func main() {
	x := subpkg.MyStruct{Foo: "foo"}
	subpkg.PrintMyStruct(x)
}

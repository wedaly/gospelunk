package main

import (
	"fmt"

	"github.com/wedaly/gospelunk/pkg/list/testdata/testmodule003/subpkg"
)

func MyFunc(x int) {
	fmt.Printf("x = %d\n", x)
}

func main() {
	MyFunc(subpkg.SubFunc())
}

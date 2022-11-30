package subpkgWithIface

import "fmt"

type MyInterface interface {
	FooString() string
}

func Print(x MyInterface) {
	fmt.Printf("x.FooString() = %q\n", x.FooString())
}

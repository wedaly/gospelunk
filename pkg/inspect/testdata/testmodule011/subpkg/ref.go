package subpkg

import (
	"fmt"
)

func PrintMyStruct(x MyStruct) {
	fmt.Printf("%v", x)
	fmt.Printf("%s", x.MyMethod())
}

func CallMyFunc() {
	x := MyFunc()
	fmt.Printf("%v", x)
}

func PrintMyInterface(x MyInterface) {
	fmt.Printf("%s", x.MyMethod())
}

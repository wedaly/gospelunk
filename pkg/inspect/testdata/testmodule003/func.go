package testmodule003

import (
	"fmt"

	"github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule003/subpkg"
)

func FuncInSameFile() int {
	return 0
}

func CallFuncs() {
	x := FuncInSameFile()
	y := FuncInOtherFile()
	z := subpkg.FuncInOtherPkg(x + y)
	fmt.Printf("x = %d, y = %d, z = %d\n", x, y, z)
}

package main

import (
	"github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule010/subpkgWithIface"
	"github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule010/subpkgWithImpl"
)

func main() {
	impl := subpkgWithImpl.MyInterfaceImpl{Foo: "foo"}
	subpkgWithIface.Print(&impl)
}

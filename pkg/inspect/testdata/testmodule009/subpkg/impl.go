package subpkg

import (
	p "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule009"
)

type MyInterfaceSubpkgImpl struct{}

func (x MyInterfaceSubpkgImpl) IfaceMethodOne() string {
	return "test"
}

func (x MyInterfaceSubpkgImpl) IfaceMethodTwo(param p.ParamType) int {
	return int(param) + 1
}

var _ p.MyInterface = &MyInterfaceSubpkgImpl{}

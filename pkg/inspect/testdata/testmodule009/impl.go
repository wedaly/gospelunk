package testmodule009

type MyInterfaceImpl struct{}

func (x MyInterfaceImpl) IfaceMethodOne() string {
	return "test"
}

func (x MyInterfaceImpl) IfaceMethodTwo(param ParamType) int {
	return int(param) + 1
}

type MyInterfacePointerImpl struct{}

func (y *MyInterfacePointerImpl) IfaceMethodOne() string {
	return "test"
}

func (y *MyInterfacePointerImpl) IfaceMethodTwo(param ParamType) int {
	return int(param) + 1
}

var _ MyInterface = MyInterfaceImpl{}
var _ MyInterface = &MyInterfacePointerImpl{}

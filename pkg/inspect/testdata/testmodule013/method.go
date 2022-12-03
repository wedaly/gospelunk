package testmodule013

type MyMethodReturn int

type MyStructWithPointerRecv struct{}

func (x *MyStructWithPointerRecv) MyMethod() MyMethodReturn {
	return MyMethodReturn(1)
}

package subpkg

type MyStruct struct {
	Foo string
}

func MyFunc() string {
	return "test"
}

type MyInterface interface {
	MyMethod() string
}

func (x MyStruct) MyMethod() string {
	return "method impl"
}

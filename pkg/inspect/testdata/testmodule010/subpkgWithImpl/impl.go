package subpkgWithImpl

type MyInterfaceImpl struct {
	Foo string
}

func (x *MyInterfaceImpl) FooString() string {
	return x.Foo
}

package testmodule001

var MyVar string

const MyConst = 123

type MyStruct struct {
	MyField      string
	privateField int
}

type privateStruct struct {
	PublicField  int
	privateField int
}

type MyInterface interface {
	String() string
}

func MyFunc() string {
	return "Test"
}

func privateFunc() string {
	return "private"
}

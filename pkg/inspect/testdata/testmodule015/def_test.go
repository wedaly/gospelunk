package testmodule015

import "testing"

type MyTestStruct struct{}

var testVar MyTestStruct

func TestMyFunc(t *testing.T) {
	s := MyStruct{MyField: "test"}
	f := MyFunc(s)
	if f != "test" {
		t.Fail()
	}
}

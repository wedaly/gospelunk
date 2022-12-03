package testmodule013

type MyEmbeddedStruct struct{}

type MyParentStruct struct {
	MyEmbeddedStruct
	MyField int
}

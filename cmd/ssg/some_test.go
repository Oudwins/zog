package main

// example with most basic usage
//
//go:generate go run main.go ssg
type MyType struct {
	Field1 string `json:"field1"`
	Field2 int
}

// example with custom output path
//
//go:generate go run main.go ssg -output=./schema/generated.go
type MyOtherType struct {
	Field1 string
	Field2 int
}

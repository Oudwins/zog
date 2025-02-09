package main

//go:generate go run main.go ssg -output=./schema/generated.go
type MyType struct {
	Field1 string `json:"field1"`
	Field2 int
}

//go:generate go run main.go ssg -output=./schema/generated.go
type MyOtherType struct {
	Field1 string
	Field2 int
}

package main

import (
	"fmt"

	z "github.com/Oudwins/zog"
)

func main() {

	t := z.Time(z.Time.Format("Hello world"))
	fmt.Println("HELLO WORLD")
	fmt.Println(z.ToJson(t))
}

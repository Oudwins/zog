package main

import (
	"fmt"
	"os"

	"github.com/Oudwins/zog/cmd/zog/ssg"
)

func main() {
	switch os.Args[1] {
	case "ssg":
		ssg.Run(os.Args[2:])
	default:
		fmt.Println("Unknown command")
	}
}

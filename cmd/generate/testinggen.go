package main

import (
	z "github.com/Oudwins/zog"
)

/*
I am writing something to generate go code. But I need to do so with information from go:generate command. Basically here is the first step of what I am trying to do.

I need to get code like this
go
// go run zog/cmd
var userSchema = z.Struct(z.Schema{
	"Name": z.String(),
	"Age":  z.Int(),
})


execute on the go generate command and create this structure from the code:

json
{
 name: "userSchema",
type: "struct",
 Shape: {
 "Name": {
 type: "string"
},
"Age": {
   type: "int"
}
}
}


The actual code that I am analyzing creates data structures that could give me some of that information but not all. I'm not sure what the best approach for this is. What do you think?

*/

//go:generate go run ./main.go
var userSchema = z.Struct(z.Schema{
	"Name": z.String(),
	"Age":  z.Int(),
})

// another comment

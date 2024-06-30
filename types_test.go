package zog

// func TestDefaulter(t *testing.T) {
// 	type User struct {
// 		Email string
// 	}
// 	schema := Schema{
// 		"email": Default("foo@bar.com", String().Email()),
// 	}
// 	user := User{
// 		Email: "",
// 	}
// 	errors, ok := Validate(&user, schema)
// 	fmt.Println(user)
// 	fmt.Println(errors)
// 	fmt.Println(ok)
// }

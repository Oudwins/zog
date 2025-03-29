---
sidebar_position: 2
hide_table_of_contents: false
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Getting Started

#### **1 Install**

```bash
go get github.com/Oudwins/zog
```

#### **2 Create a user schema and its struct**

```go
import (
	z "github.com/Oudwins/zog"
)

type User struct {
	Name string
	Age  int
}

var userSchema = z.Struct(z.Schema{
	// its very important that schema keys like "name" match the struct field name NOT the input data
	"name": z.String().Min(3, z.Message("Override default message")).Max(10),
	"age":  z.Int().GT(18),
})
```

#### **3 Validate your schema**

**Using [schema.Parse()](https://zog.dev/core-concepts/parsing)**

```go
func main() {
	u := User{}
	m := map[string]string{
		"name": "Zog",
		"age":  "", // won't return an error because fields are optional by default
	}
	errsMap := userSchema.Parse(m, &u)
	if errsMap != nil {
		// handle errors -> see Errors section
	}
	u.Name // "Zog"
	// note that this might look weird but we didn't say age was required so Zog just skipped the empty string and we are left with the uninitialized int
	// If we need 0 to be a valid value for age we can use a pointer to an int which will be nil if the value was not present in the input data
	u.Age // 0
}
```

**Using [schema.Validate()](https://zog.dev/core-concepts/validate)**

```go
func main() {
	u := User{
		Name: "Zog",
		Age:  0, // wont return an error because fields are optional by default otherwise it will error
	}
	errsMap := userSchema.Validate(&u)
	if errsMap != nil {
		// handle errors -> see Errors section
	}
}
```

#### **4. Its easy to use with http & json**

The [zhttp package](https://zog.dev/packages/zhttp) has you covered for JSON, Forms and Query Params, just do:

```go
import (
	zhttp "github.com/Oudwins/zog/zhttp"
)

err := userSchema.Parse(zhttp.Request(r), &user)
```

If you are receiving json some other way you can use the [zjson package](https://zog.dev/packages/zjson)

```go
import (
	zjson "github.com/Oudwins/zog/zjson"
)

err := userSchema.Parse(zjson.Decode(bytes.NewReader(jsonBytes)), &user)
```

#### **5. Or to validate your environment variables**

The [zenv package](https://zog.dev/packages/zenv) has you covered, just do:

```go
import (
	zenv "github.com/Oudwins/zog/zenv"
)

err := envSchema.Parse(zenv.NewDataProvider(), &envs)
```

#### **6. You can also parse individual fields**

```go
var t = time.Time
errsList := Time().Required().Parse("2020-01-01T00:00:00Z", &t)
```

#### **7 And do stuff before and after parsing**

```go
var dest []string
Slice(String().Email().Required()).PreTransform(func(data any, ctx z.Ctx) (any, error) {
	s := data.(string)
	return strings.Split(s, ","), nil
}).PostTransform(func(destPtr any, ctx z.Ctx) error {
	s := destPtr.(*[]string)
	for i, v := range *s {
		(*s)[i] = strings.TrimSpace(v)
	}
	return nil
}).Parse("foo@bar.com,bar@foo.com", &dest) // dest = [foo@bar.com bar@foo.com]
```

<div align="center">
    <br />
    <a href="https://github.com/Oudwins/zog">
     <img src="https://raw.githubusercontent.com/Oudwins/zog/master/assets/zog-banner.png" alt="Zog, a Zod-like schema parser & validator" />
    </a>
</div>

# ZOG - A Zod & Yup like Schema Parser & Validator for GO

[![Coverage Status](https://coveralls.io/repos/github/Oudwins/zog/badge.svg?branch=master)](https://coveralls.io/github/Oudwins/zog?branch=master)
[![Go Report Card](https://goreportcard.com/badge/Oudwins/zog)](https://goreportcard.com/report/github.com/Oudwins/zog)
[![GitHub tag](https://img.shields.io/github/tag/Oudwins/zog?include_prereleases=&sort=semver&color=blue)](https://github.com/Oudwins/zog/releases/)
<a href="https://pkg.go.dev/github.com/Oudwins/zog"><img src="https://pkg.go.dev/badge/github.com//github.com/Oudwins/tailwind-merge-go.svg" alt="Go Reference" /></a>
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Mentioned in Awesome Templ](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/templ-go/awesome-templ)
[![stars - zog](https://img.shields.io/github/stars/Oudwins/zog?style=social)](https://github.com/Oudwins/zog)

[![view - Documentation](https://img.shields.io/badge/view-Documentation-blue?style=for-the-badge)](https://zog.dev)

Zog is a schema builder for runtime value parsing and validation. Define a schema, transform a value to match, assert the shape of an existing value, or both. Zog schemas are extremely expressive and allow modeling complex, interdependent validations, or value transformations.

Killer Features:

- Concise yet expressive schema interface, equipped to model simple to complex data models
- **[Zod](https://github.com/colinhacks/zod)-like API**, use method chaining to build schemas in a typesafe manner
- **Extensible**: add your own Tests and Schemas
- **Rich errors** with detailed context, make debugging a breeze
- **Fast**: Zog is one of the fastest Go validation libraries. We are just behind the goplayground/validator for most of the [govalidbench](https://github.com/Oudwins/govalidbench/tree/master) benchmarks.
- **Built-in coercion** support for most types
- Zero dependencies!
- **Four Helper Packages**
  - **zenv**: parse environment variables
  - **zhttp**: parse http forms & query params
  - **zjson**: parse json
  - **i18n**: Opinionated solution to good i18n zog errors

> **API Stability:**
>
> - I will consider the API stable when we reach v1.0.0
> - However, I believe very little API changes will happen from the current implementation. The APIs most likely to change are the **data providers** (please don't make your own if possible use the helpers whose APIs will not change meaningfully) and the z.Ctx most other APIs should remain the same. I could be wrong but I don't expect many breaking changes.
> - Zog will not respect semver until v1.0.0 is released. Consider each minor version to potentially have breaking changes until then.

## Introduction

#### **0. Read the docs at [zog.dev](https://zog.dev)**

Or don't, below is the quickstart guide

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
	Name string `zog:"firstname"` // tag is optional. If not set zog will check for "name" field in the input data
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
		"firstname": "Zog", // Note we are using "firstname" here as specified in the struct tag
		"age":       "",    // won't return an error because fields are optional by default
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

## Roadmap

These are some of the things I want to add to zog before v1.0.0

- Support for schema.Clone()
- support for catch & default for structs & slices
- Struct generation from the schemas

## Support

The damm domain costs me some outrageous amount like 100$ a year, so if any one wants to help cover that cost through github sponsors that is more than welcome.

## Acknowledgments

- Big thank you to @AlexanderArvidsson for being there to talk about architecture and design decisions. It helped a lot to have someone to bounce ideas off of
- Credit for all the inspiration goes to /colinhacks/zod & /jquense/yup
- Credit for the initial idea goes to anthony (@anthonyGG) -> /anthdm/superkit he made a hacky version of this idea that I used as a starting point, I was never happy with it so I inspired me to rewrite it from scratch. I owe him a lot
- Credit for the zod logo goes to /colinhacks/zod

## License

This project is licensed under the MIT License -
see the [LICENSE](LICENSE) file for details.

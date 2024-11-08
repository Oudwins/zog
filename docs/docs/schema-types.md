---
sidebar_position: 5
---

# Schema types

```go
// Primtives. Calling .Parse() on these will return []ZogError
z.String()
z.Int()
z.Float()
z.Bool()
z.Time()

// Complex Types. Calling .Parse() on these will return map[string][]ZogError. Where the key is the field path ("user.email") & $root is the list of complex type level errors not the specific field errors
z.Struct(z.Schema{
  "name": z.String(),
})
z.Slice(z.String())

```

## Primtive Types

### String

```go
// PreTransforms
z.String().Trim() // trims the input data of whitespace if it is a string (does nothing otherwise)

// Tests / Validations
z.String().Test() // custom test
z.String().Min(5) // validates min length
z.String().Max(10) // validates max length
z.String().Len(5) // validates length
z.String().Email() // validates email
z.String().URL() // validates url
z.String().UUID() // validates uuid v4
z.String().Match(regex) // matches a regex
z.String().Contains(substring) // validates string contains substring
z.String().ContainsUpper() // validates string contains uppercase letter
z.String().ContainsDigit() // validates string contains digit
z.String().ContainsSpecial() // validates string contains special character
z.String().HasPrefix(prefix) // validates string has prefix
z.String().HasSuffix(suffix) // validates string has suffix
z.String().OneOf([]string{"a", "b", "c"}) // validates string is one of the values
```

### Numbers / Ints & Floats

```go
// Tests / Validators
z.Int().GT(n) // validates int is greater than n
z.Float().GTE(n) // validates float is greater than or equal to n
z.Int().LT(n) // validates int is less than n
z.Float().LTE(n) // validates float is less than or equal to n
z.Int().EQ(n) // validates int is equal to n
z.Float().OneOf([]float64{1.0, 2.0, 3.0}) // validates float is one of the values
```

### Booleans

```go
// Tests / Validators
z.Bool().True() // validates bool is true
z.Bool().False() // validates bool is false
```

### Times & Dates

Use Time to validate `time.Time` instances

```go
// Tests / Validators
z.Time().After(time.Now()) // validates time is after now
z.Time().Before(time.Now()) // validates time is before now
z.Time().Is(time.Now()) // validates time is equal to now
```

## Complex Types

### Structs

```go
// usage
s := z.Struct(z.Schema{
  "name": String().Required(),
  "age": Int().Required(),
})
// Tests / Validators
// None right now
```

### Slices

```go
// usage
schema := z.Slice(String())
// Tests / Validators
z.Slice(Int()).Min(5) // validates slice has at least 5 elements
z.Slice(Float()).Max(5) // validates slice has at most 5 elements
z.Slice(Bool()).Length(5) // validates slice has exactly 5 elements
z.Slice(String()).Contains("foo") // validates slice contains the element "foo"
```

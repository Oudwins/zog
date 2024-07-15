##

### Types

```go
// Primtives
String()
Int()
Float()
Bool()
Time()

// nil
Nil()
```

#### Strings

```go
// Validations
String().Min(5)
String().Max(10)
String().Len(5)
String().Email()
String().URL()
String().Emoji() // TODO
String().UUID() // TODO
String().NanoID() // TODO
String().CuID() // TODO
String().CuID2() // TODO
String().Regex(regex) // TODO
String().Contains(string)
String().ContainsUpper()
String().ContainsDigit()
String().ContainsSpecial()
String().StartsWith(string)
String().EndsWith(string)
String().Datetime(); // TODO ->  ISO 8601; by default only `Z` timezone allowed -> https://zod.dev/?id=datetimes
String().Date() // TODO -> 2020-01-01 -> https://zod.dev/?id=dates
String().Time() // TODO -> 12:00:00 -> https://zod.dev/?id=times
String().IP(); // TODO -> defaults to allow both IPv4 and IPv6 -> https://zod.dev/?id=ip-addresses

```

#### Numbers

```go
// Validators
Int().Min(5)
Int().Max(10)
Int().GT(5)
Int().GTE(5)
Int().LT(5)
Int().LTE(5)
Int().EQ(5)
Float().Min(5)
...
```

#### Booleans

```go
Bool().True()
Bool().False()
```

#### Times & Dates

Use Time to validate `time.Time` instances

```go
Time().After(time.Now())
Time().Before(time.Now())
Time().Is(time.Now())
// https://zod.dev/?id=dates-1
Time().Min() // TODO -> Should we have a Min & Max for times like zod?
Time().Max() // TODO -> Should we have a Min & Max for times like zod?
```

#### Enums

```go
e := Enum([]string{"a", "b", "c"})

e.Extract([]string{"a"}) // TODO -> https://zod.dev/?id=zod-enums
e.Exclude([]string{"a"}) // TODO -> https://zod.dev/?id=zod-enums
```

#### Structs

```go
s := Struct{
}

s.Merge(s2)
s.Pick() // TODO -> https://zod.dev/?id=pickomit
s.Omit() // TODO -> https://zod.dev/?id=pickomit
s.Partial() // TODO -> https://zod.dev/?id=partial
s.Required() // TODO -> https://zod.dev/?id=required
```

#### Slices

```go
s := Slice(schema)

s.Element // TODO -> access the schema inside the slice


Slice(schema).Min(5)
Slice(schema).Max(5)
Slice(schema).Length(5)
```

#### Unions

Should we do this. Will it be useful? TODO?

- https://zod.dev/?id=unions

#### Maps

TODO

- https://zod.dev/?id=maps

#### Sets

TODO

- https://zod.dev/?id=sets

#### Effects

Optionals

```go
op := Optional(schema)

op.Unwrap() // TODO -> https://zod.dev/?id=optionals
```

Catch

```go
c := Catch(schema, valueOnError)
```

transform(func)

refine(func, params) // refine validation, message

preprocess(func, schema)

pipe(schema)

#### IMPORTANT NOTES

So embeded structs are NOT NIL by default

```go
// z.WithMessage(), z.String.IP.WithVersion(4|6)
// each just returns a func(rule *Rule) -> transforms what ever
z.string({
required_error: "Name is required",
invalid_type_error: "Name must be a string",
});
```

### Common Methods

```go
.In(values []T) // checks that the value is equal to one of the values
.Optional() // marks the field as optional. If the value of the field is the type's zero value, the field will be ignored. Otherwise, it will be validated as normal
.Refine(ruleName string, errorMsg string, validateFunc p.RuleValidateFunc) *T // adds a new rule to the validator. Use this to add custom rules. It takes a ruleName, an errorMessage and a validateFunc func(rule Rule) bool
```

## Planned Improvements

1. Option to override the default message for a rule
2. Type coersion

## Notes

TODO

How to handle coersion

- Add ParseValue(val any) any -> by default returns the same value but can be overriden
  - to the validationFieldSchema
- Make the coercion a wrapper. But will take a schema which won't be type safe -> Coerce().String().Schema(IntSchema)

z.Struct()
z.Map()
z.Slice()

// for customizing errors provide a map of rule name to error function
// make your own errors?
// pass the error fnc as second optional arg?

Questions for anthony about the validation library:

2. Why does it not take the fieldValue as an argument instead of holding the fieldValue and Fieldname in the ruleset? FieldName is never used
3.

Default values:
Conceptually, this is how Zod processes default values:

1. If the input is undefined, the default value is returned
2. Otherwise, the data is parsed using the base schema

Optionals: if value is zero value then errors are ignored

catch: provides a value to be used if the validation fails

We could return nil from parse function is value did not change.

!!
preprocess instead of coerce

How do we build the error map...

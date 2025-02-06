---
sidebar_position: 2
---

# Parsing vs Validation


Zog supports two main ways of processing data, both of which support the exact same schemas and can be used interchangeably without modifying the schema:
- [`Schema.Parse(data, &dest, ...options)`](/core-concepts/parsing) - Parses the data into the destination pointer and returns a list of errors if any.
- [`Schema.Validate(&data, ...options)`](/core-concepts/validate) - Validates the data and returns a list of errors if any.


> **You are probably wondering**
> [What is the difference?](#What is the difference?)
> [Which should I use?](#Which should I use?)


## What is the difference?
There is only one difference between the two:
- For [Schema.Parse(data, &dest, ...options)](/core-concepts/parsing) you must provide data that Zog will parse into the destination structure. For example, if you use one of the helper packages like [zog-json](/packages/zjson) zog will unmarshal the json into the destination structure.
- For [Schema.Validate(&data, ...options)](/core-concepts/validate) you are expected to have already parsed the data into the final structure you want and are now just validating that it is correct.


**Okay, but what does this mean in practice?**
It means that Parse will handle things like type coercion, zero value checking, etc... for you. Whereas Validate will not. For example:

```go
data := "2024-01-01"
var dest time.Time
z.Time().Parse(data, &dest) // dest will be 2024-01-01 00:00:00 +0000 UTC
z.Time().Validate(&data) // Error: string is not a valid time
```


It also means that Validate cannot know if a value was provided by the user or if it was set by a default value. Therefore it will consider zero values as invalid when the schema is required. For example:

```go
var dest int
z.Int().Required().Parse(0, &dest) // dest will be 0
val := 0
z.Int().Required().Validate(&val) // will return a required error

// To fix this you can use a pointer:
var valPtr *int
z.Ptr(z.Int()).NotNil().Validate(&valPtr) // Error pointer is nil
*valPtr = 0
z.Ptr(z.Int()).NotNil().Validate(&valPtr) // No error
```


## Which should I use?
If you can, use [`Schema.Validate(&data, ...options)`](/core-concepts/validate) as it is more efficient. But if you need the type coercion or zero value checking without pointers feel free to use [`Schema.Parse(data, &dest, ...options)`](/core-concepts/parsing).



## Validation Execution Structure
Generally speaking when executing `schema.Validate()` Zog will follow a very similar execution structure to the one described in [Parsing Execution Structure](/core-concepts/parsing/#parsing-execution-structure).




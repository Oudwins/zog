```go
z.String()
z.Struct(z.Schema{})
z.Preprocess()
```

```json
{
  "type": "string"
}
```

We cannot really get the error messages right? Thats a little complicated to do since they depend on the input maybe

We need to give each processor an ID also. This way we can identify them in the json schema. Maybe we merge issueCode and transformId in the template

- What should be the APi for adding an ID to a transform? z.String().Transform(fn, z.ID())

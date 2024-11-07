---
sidebar_position: 1
---

# Using Zog with HTML templates

(In this example I'll use go templ, but you can use any template engine)

**Example use case: simplified Signup form validation**
Imagine our handler looks like this:

```go
type SignupFormData struct {
  Email string
  Password string
}

schema := z.Struct(z.Schema{
  "email": z.String().Email().Required(),
  "password": z.String().Min(8).Required(),
})

func handleSignup(w http.ResponseWriter, r *http.Request) {
  var signupFormData = SignupFormData{}
  errs := schema.Parse(zhttp.NewRequestDataProvider(r), &signupFormData)

  if errs != nil {
    www.Render(signupFormTempl(&signupFormData, errs))
  }
  // handle successful signup
}

templ signupFormTempl(data *SignupFormData, errs z.ZogErrMap) {
  <input type="text" name="email" value={data.Email}>
  // display only the first error
  if e, ok := errs["email"]; ok {
    <p class="error">{e[0].Message()}</p>
  }
  <input type="text" name="password" value={data.Password}>
  // display only the first error
  if e, ok := errs["password"]; ok {
    <p class="error">{e[0].Message()}</p>
  }
}
```

**PS:** If you are using go html templates & tailwindcss you might be interesting in my port of [tailwind-merge to go.](https://github.com/Oudwins/tailwind-merge-go)

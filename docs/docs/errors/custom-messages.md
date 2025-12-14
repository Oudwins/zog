---
sidebar_position: 2
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Custom Messages

Zog has multiple ways of customizing issue messages as well as support for [i18n](/packages/i18n). Here is a list of the ways you can customize issue messages:

#### **1. Using the z.Message() function**

This is a function available for all tests, it allows you to set a custom message for the test.

```go
err := z.String().Min(5, z.Message("string must be at least 5 characters long")).Parse("bad", &dest)
// err = []ZogIssue{{Message: "string must be at least 5 characters long"}}
```

#### **2. Using the z.MessageFunc() function**

This is a function available for all tests, it allows you to set a custom message for the test.

This function takes in an `IssueFmtFunc` which is the function used to format issue messages in Zog. It has the following signature:

```go
type IssueFmtFunc = func(e *ZogIssue, ctx z.Ctx)
```

```go
err := z.String().Min(5, z.MessageFunc(func(e *z.ZogIssue, ctx z.Ctx) {
	e.SetMessage("string must be at least 5 characters long")
})).Parse("bad", &dest)
// err = []ZogIssue{{Message: "string must be at least 5 characters long"}}
```

#### **3. Using the WithIssueFormatter() ExecOption**

This allows you to set a custom `ZogIssue` formatter for the entire parsing operation. Beware you must handle all `ZogIssue` codes & types or you may get unexpected messages.

```go
err := z.String().Min(5).Email().Parse("zog", &dest, z.WithIssueFormatter(func(e *z.ZogIssue, ctx z.Ctx) {
	e.SetMessage("override message")
}))
// err = []ZogIssue{{Code: min_length_issue, Message: "override message"}, {Code: email_issue, Message: "override message"}}
```

See how our issue messages were overridden? Be careful when using this!

#### **4. Iterate over the returned issues and create custom messages**

```go
errs := userSchema.Parse(data, &user)
msgs := FlattenZogIssues(errs)

// This is basically the implementation of the z.Issues.Flatten() method
func FlattenZogIssues(errs z.ZogIssueList) map[string][]string {
	// iterate over issues and create custom messages based on the issue code, the params and destination type
	result := make(map[string][]string)
	for _, issue := range errs {
		path := issue.PathString()
		if path == "" {
			path = "$root"
		}
		result[path] = append(result[path], issue.Message)
	}
	return result
}
```

#### **5. Configure issue messages globally**

Zog provides a `conf` package where you can override the issue messages for specific issue codes. You will have to do a little digging to be able to do this. But here is an example:

```go
import (
	conf "github.com/Oudwins/zog/zconf"
	zconst "github.com/Oudwins/zog/zconst"
)

// override specific issue messages
// For this I recommend you import `zog/zconst` which contains zog constants but you can just use strings if you prefer
conf.DefaultIssueMessageMap[zconst.TypeString]["my_custom_issue_code"] = "my custom issue message"
conf.DefaultIssueMessageMap[zconst.TypeString][zconst.IssueCodeRequired] = "Now all required issues will get this message"
```

But you can also outright override the issue formatter and ignore the issues map completely:

```go
// override the issue formatter function - CAREFUL with this you can set every issue message to the same thing!
conf.IssueFormatter = func(e *p.ZogIssue, ctx z.Ctx) {
	// do something with the issue
	...
	// fallback to the default issue formatter
	conf.DefaultIssueFormatter(e, p) // this uses the DefaultErrMsgMap to format the issue messages
}
```

#### **6. Use the [i18n](/packages/i18n) package**

Really this only makes sense if you are doing i18n. Please please check out the [i18n](/packages/i18n) section for more information.

---

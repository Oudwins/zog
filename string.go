package zog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/primitives"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	urlRegex   = regexp.MustCompile(`^(http(s)?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?$`)
)

type stringProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *string
	required       *p.Test
	catch          *string
}

func String() *stringProcessor {
	return &stringProcessor{
		tests: []p.Test{},
	}
}

func (v *stringProcessor) Parse(val any, dest *string) p.ZogErrorList {
	// TODO create context -> but for single field
	var ctx = p.NewParseCtx()
	errs := p.NewErrsList()
	path := p.PathBuilder("")
	// TODO handle options

	v.process(val, dest, errs, path, ctx)

	return errs.List
}

func (v *stringProcessor) process(val any, dest any, errs p.ZogErrors, path p.PathBuilder, ctx *p.ParseCtx) {
	primitiveProcessor(val, dest, errs, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, conf.Coercers["string"])
}

// Adds pretransform function to schema
func (v *stringProcessor) PreTransform(transform p.PreTransform) *stringProcessor {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *stringProcessor) PostTransform(transform p.PostTransform) *stringProcessor {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *stringProcessor) Required(options ...TestOption) *stringProcessor {
	r := p.Required(p.DErrorFunc("is a required field"))
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *stringProcessor) Optional() *stringProcessor {
	v.required = nil
	return v
}

// sets the default value
func (v *stringProcessor) Default(val string) *stringProcessor {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *stringProcessor) Catch(val string) *stringProcessor {
	v.catch = &val
	return v
}

// ! VALIDATORS
// custom test function call it -> schema.Test("test_name", z.Message(""), func(val any, ctx *p.ParseCtx) bool {return true})
func (v *stringProcessor) Test(ruleName string, errorMsg TestOption, validateFunc p.TestFunc) *stringProcessor {
	t := p.Test{
		Name:         ruleName,
		ErrorFunc:    nil,
		ValidateFunc: validateFunc,
	}
	errorMsg(&t)
	v.tests = append(v.tests, t)

	return v
}

// checks that the value is one of the enum values
func (v *stringProcessor) OneOf(enum []string, options ...TestOption) *stringProcessor {
	t := p.In(enum, fmt.Sprintf("should be one of %v", enum))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is at least n characters long
func (v *stringProcessor) Min(n int, options ...TestOption) *stringProcessor {
	t := p.LenMin[string](n, p.DErrorFunc(fmt.Sprintf("should be at least %d characters long", n)))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is at most n characters long
func (v *stringProcessor) Max(n int, options ...TestOption) *stringProcessor {
	t := p.LenMax[string](n, p.DErrorFunc(fmt.Sprintf("should be at most %d characters long", n)))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is exactly n characters long
func (v *stringProcessor) Len(n int, options ...TestOption) *stringProcessor {
	t := p.Len[string](n, p.DErrorFunc(fmt.Sprintf("should be exactly %d characters long", n)))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks that the value is a valid email address
func (v *stringProcessor) Email(options ...TestOption) *stringProcessor {
	t := p.Test{
		Name:      "email",
		ErrorFunc: p.DErrorFunc("is not a valid email address"),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			email, ok := v.(string)
			if !ok {
				return false
			}
			return emailRegex.MatchString(email)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) URL(options ...TestOption) *stringProcessor {
	t := p.Test{
		Name:      "url",
		ErrorFunc: p.DErrorFunc("is not a valid url"),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			u, ok := v.(string)
			if !ok {
				return false
			}
			isOk := urlRegex.MatchString(u)
			return isOk
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) HasPrefix(s string, options ...TestOption) *stringProcessor {
	t := p.Test{
		Name:      "startsWith",
		ErrorFunc: p.DErrorFunc(fmt.Sprintf("should start with %s", s)),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.HasPrefix(val, s)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) HasSuffix(s string, options ...TestOption) *stringProcessor {
	t := p.Test{
		Name:      "endsWith",
		ErrorFunc: p.DErrorFunc(fmt.Sprintf("should end with %s", s)),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.HasSuffix(val, s)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) Contains(sub string, options ...TestOption) *stringProcessor {
	t := p.Test{
		Name:      "contains",
		ErrorFunc: p.DErrorFunc(fmt.Sprintf("should contain %s", sub)),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.Contains(val, sub)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) ContainsUpper(options ...TestOption) *stringProcessor {
	t := p.Test{
		Name:      "containsUpper",
		ErrorFunc: p.DErrorFunc("should contain at least one uppercase letter"),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			for _, r := range val {
				if r >= 'A' && r <= 'Z' {
					return true
				}
			}
			return false
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) ContainsDigit(options ...TestOption) *stringProcessor {
	t := p.Test{
		Name: "containsDigit",
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			for _, r := range val {
				if r >= '0' && r <= '9' {
					return true
				}
			}
			return false
		},
		ErrorFunc: p.DErrorFunc("should contain at least one digit"),
	}

	for _, opt := range options {
		opt(&t)
	}

	v.tests = append(v.tests, t)
	return v
}

func (v *stringProcessor) ContainsSpecial(options ...TestOption) *stringProcessor {
	t :=
		p.Test{
			Name: "containsSpecial",
			ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
				val, ok := v.(string)
				if !ok {
					return false
				}
				for _, r := range val {
					if (r >= '!' && r <= '/') ||
						(r >= ':' && r <= '@') ||
						(r >= '[' && r <= '`') ||
						(r >= '{' && r <= '~') {
						return true
					}
				}
				return false
			},
			ErrorFunc: p.DErrorFunc("should contain at least one special character"),
		}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

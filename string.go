package zog

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var (
	_ PrimitiveZogSchema[string] = (*StringSchema)(nil)
	_ NotStringSchema            = (*StringSchema)(nil)
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	uuidRegex  = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
)

type NotStringSchema interface {
	Test(t p.Test, opts ...TestOption) *StringSchema
	OneOf(enum []string, options ...TestOption) *StringSchema
	Min(n int, options ...TestOption) *StringSchema
	Max(n int, options ...TestOption) *StringSchema
	Len(n int, options ...TestOption) *StringSchema
	Email(options ...TestOption) *StringSchema
	URL(options ...TestOption) *StringSchema
	HasPrefix(s string, options ...TestOption) *StringSchema
	HasSuffix(s string, options ...TestOption) *StringSchema
	Contains(sub string, options ...TestOption) *StringSchema
	ContainsUpper(options ...TestOption) *StringSchema
	ContainsDigit(options ...TestOption) *StringSchema
	ContainsSpecial(options ...TestOption) *StringSchema
	UUID(options ...TestOption) *StringSchema
	Match(regex *regexp.Regexp, options ...TestOption) *StringSchema

	// `Not` method is missing here as we do not want the user to do `Not` chaining.
}

type StringSchema struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *string
	required       *p.Test
	catch          *string
	coercer        conf.CoercerFunc
	isNot          bool
}

// ! INTERNALS

// Returns the type of the schema
func (v *StringSchema) getType() zconst.ZogType {
	return zconst.TypeString
}

// Sets the coercer for the schema
func (v *StringSchema) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

// Returns a new String Schema
func String(opts ...SchemaOption) *StringSchema {
	s := &StringSchema{
		coercer: conf.Coercers.String, // default coercer
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Parses the data into the destination string. Returns a list of ZogIssues
func (v *StringSchema) Parse(data any, dest *string, options ...ExecOption) p.ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()

	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}

	path := p.NewPathBuilder()
	defer path.Free()
	v.process(ctx.NewSchemaCtx(data, dest, path, v.getType()))

	return errs.List
}

// Internal function to process the data
func (v *StringSchema) process(ctx *p.SchemaCtx) {
	primitiveProcessor(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Validate Given string
func (v *StringSchema) Validate(data *string, options ...ExecOption) p.ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}

	path := p.NewPathBuilder()
	defer path.Free()
	v.validate(ctx.NewSchemaCtx(data, data, path, v.getType()))
	return errs.List
}

// Internal function to validate the data
func (v *StringSchema) validate(ctx *p.SchemaCtx) {
	primitiveValidator(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch)
}

// Adds pretransform function to schema
func (v *StringSchema) PreTransform(transform p.PreTransform) *StringSchema {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// PreTransform: trims the input data of whitespace if it is a string
func (v *StringSchema) Trim() *StringSchema {
	v.preTransforms = append(v.preTransforms, func(val any, ctx Ctx) (any, error) {
		switch v := val.(type) {
		case string:
			return strings.TrimSpace(v), nil
		default:
			return val, nil
		}
	})
	return v
}

// Adds posttransform function to schema
func (v *StringSchema) PostTransform(transform p.PostTransform) *StringSchema {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *StringSchema) Required(options ...TestOption) *StringSchema {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *StringSchema) Optional() *StringSchema {
	v.required = nil
	return v
}

// sets the default value
func (v *StringSchema) Default(val string) *StringSchema {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *StringSchema) Catch(val string) *StringSchema {
	v.catch = &val
	return v
}

// ! PRETRANSFORMS

// ! Tests
// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *StringSchema) Test(t p.Test, opts ...TestOption) *StringSchema {
	for _, opt := range opts {
		opt(&t)
	}

	t.ValidateFunc = customTestBackwardsCompatWrapper(t.ValidateFunc)
	return v.addTest(t)
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *StringSchema) TestFunc(testFunc p.TestFunc, options ...TestOption) *StringSchema {
	test := TestFunc("", testFunc)
	v.Test(test, options...)
	return v
}

// Test: checks that the value is one of the enum values
func (v *StringSchema) OneOf(enum []string, options ...TestOption) *StringSchema {
	t := p.In(enum)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is at least n characters long
func (v *StringSchema) Min(n int, options ...TestOption) *StringSchema {
	t := p.LenMin[string](n)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is at most n characters long
func (v *StringSchema) Max(n int, options ...TestOption) *StringSchema {
	t := p.LenMax[string](n)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is exactly n characters long
func (v *StringSchema) Len(n int, options ...TestOption) *StringSchema {
	t := p.Len[string](n)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is a valid email address
func (v *StringSchema) Email(options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeEmail,
		ValidateFunc: func(v any, ctx Ctx) bool {
			email, ok := v.(*string)
			if !ok {
				return false
			}
			return emailRegex.MatchString(*email)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is a valid URL
func (v *StringSchema) URL(options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeURL,
		ValidateFunc: func(v any, ctx Ctx) bool {
			s, ok := v.(*string)
			if !ok {
				return false
			}
			u, err := url.Parse(*s)
			return err == nil && u.Scheme != "" && u.Host != ""
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value has the prefix
func (v *StringSchema) HasPrefix(s string, options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeHasPrefix,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*string)
			if !ok {
				return false
			}
			return strings.HasPrefix(*val, s)
		},
	}
	t.Params[zconst.IssueCodeHasPrefix] = s
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value has the suffix
func (v *StringSchema) HasSuffix(s string, options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeHasSuffix,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*string)
			if !ok {
				return false
			}
			return strings.HasSuffix(*val, s)
		},
	}
	t.Params[zconst.IssueCodeHasSuffix] = s
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value contains the substring
func (v *StringSchema) Contains(sub string, options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeContains,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*string)
			if !ok {
				return false
			}
			return strings.Contains(*val, sub)
		},
	}
	t.Params[zconst.IssueCodeContains] = sub
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value contains an uppercase letter
func (v *StringSchema) ContainsUpper(options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeContainsUpper,
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*string)
			if !ok {
				return false
			}
			for _, r := range *val {
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
	return v.addTest(t)
}

// Test: checks that the value contains a digit
func (v *StringSchema) ContainsDigit(options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeContainsDigit,
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*string)
			if !ok {
				return false
			}
			for _, r := range *val {
				if r >= '0' && r <= '9' {
					return true
				}
			}
			return false
		},
	}

	for _, opt := range options {
		opt(&t)
	}

	return v.addTest(t)
}

// Test: checks that the value contains a special character
func (v *StringSchema) ContainsSpecial(options ...TestOption) *StringSchema {
	t :=
		p.Test{
			IssueCode: zconst.IssueCodeContainsSpecial,
			ValidateFunc: func(v any, ctx Ctx) bool {
				val, ok := v.(*string)
				if !ok {
					return false
				}
				for _, r := range *val {
					if (r >= '!' && r <= '/') ||
						(r >= ':' && r <= '@') ||
						(r >= '[' && r <= '`') ||
						(r >= '{' && r <= '~') {
						return true
					}
				}
				return false
			},
		}
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is a valid uuid
func (v *StringSchema) UUID(options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeUUID,
		ValidateFunc: func(v any, ctx Ctx) bool {
			uuid, ok := v.(*string)
			if !ok {
				return false
			}
			return uuidRegex.MatchString(*uuid)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that value matches to regex
func (v *StringSchema) Match(regex *regexp.Regexp, options ...TestOption) *StringSchema {
	t := p.Test{
		IssueCode: zconst.IssueCodeMatch,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			s, ok := v.(*string)
			if !ok {
				return false
			}
			return regex.MatchString(*s)
		},
	}
	t.Params[zconst.IssueCodeMatch] = regex.String()
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: nots the next test fn
func (v *StringSchema) Not() NotStringSchema {
	v.isNot = !v.isNot
	return v
}

func (v *StringSchema) addTest(t p.Test) *StringSchema {
	v.tests = p.AddTest(v.tests, t, v.isNot)
	if v.isNot {
		v.isNot = false
	}

	return v
}

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
	_ PrimitiveZogSchema[string] = (*StringSchema[string])(nil)
	_ NotStringSchema[string]    = (*StringSchema[string])(nil)
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	uuidRegex  = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
)

type NotStringSchema[T ~string] interface {
	Test(t p.Test, opts ...TestOption) *StringSchema[T]
	OneOf(enum []T, options ...TestOption) *StringSchema[T]
	Min(n int, options ...TestOption) *StringSchema[T]
	Max(n int, options ...TestOption) *StringSchema[T]
	Len(n int, options ...TestOption) *StringSchema[T]
	Email(options ...TestOption) *StringSchema[T]
	URL(options ...TestOption) *StringSchema[T]
	HasPrefix(s T, options ...TestOption) *StringSchema[T]
	HasSuffix(s T, options ...TestOption) *StringSchema[T]
	Contains(sub T, options ...TestOption) *StringSchema[T]
	ContainsUpper(options ...TestOption) *StringSchema[T]
	ContainsDigit(options ...TestOption) *StringSchema[T]
	ContainsSpecial(options ...TestOption) *StringSchema[T]
	UUID(options ...TestOption) *StringSchema[T]
	Match(regex *regexp.Regexp, options ...TestOption) *StringSchema[T]

	// `Not` method is missing here as we do not want the user to do `Not` chaining.
}

type StringSchema[T ~string] struct {
	preTransforms  []PreTransform
	tests          []Test
	postTransforms []PostTransform
	defaultVal     *T
	required       *Test
	catch          *T
	coercer        conf.CoercerFunc
	isNot          bool
}

// ! INTERNALS

// Returns the type of the schema
func (v *StringSchema[T]) getType() zconst.ZogType {
	return zconst.TypeString
}

// Sets the coercer for the schema
func (v *StringSchema[T]) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

// Returns a new String Schema
func String(opts ...SchemaOption) *StringSchema[string] {
	s := &StringSchema[string]{
		coercer: conf.Coercers.String, // default coercer
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Parses the data into the destination string. Returns a list of ZogIssues
func (v *StringSchema[T]) Parse(data any, dest *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()

	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}

	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, dest, path, v.getType())
	defer sctx.Free()
	v.process(sctx)

	return errs.List
}

// Internal function to process the data
func (v *StringSchema[T]) process(ctx *p.SchemaCtx) {
	primitiveProcessor(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Validate Given string
func (v *StringSchema[T]) Validate(data *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}

	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, data, path, v.getType())
	defer sctx.Free()
	v.validate(sctx)
	return errs.List
}

// Internal function to validate the data
func (v *StringSchema[T]) validate(ctx *p.SchemaCtx) {
	primitiveValidator(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch)
}

// Adds pretransform function to schema
func (v *StringSchema[T]) PreTransform(transform PreTransform) *StringSchema[T] {
	if v.preTransforms == nil {
		v.preTransforms = []PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// PreTransform: trims the input data of whitespace if it is a string
func (v *StringSchema[T]) Trim() *StringSchema[T] {
	v.preTransforms = append(v.preTransforms, func(val any, ctx Ctx) (any, error) {
		switch v := val.(type) {
		case T:
			return T(strings.TrimSpace(string(v))), nil
		default:
			return val, nil
		}
	})
	return v
}

// Adds posttransform function to schema
func (v *StringSchema[T]) PostTransform(transform PostTransform) *StringSchema[T] {
	if v.postTransforms == nil {
		v.postTransforms = []PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *StringSchema[T]) Required(options ...TestOption) *StringSchema[T] {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *StringSchema[T]) Optional() *StringSchema[T] {
	v.required = nil
	return v
}

// sets the default value
func (v *StringSchema[T]) Default(val T) *StringSchema[T] {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *StringSchema[T]) Catch(val T) *StringSchema[T] {
	v.catch = &val
	return v
}

// ! PRETRANSFORMS

// ! Tests
// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *StringSchema[T]) Test(t Test, opts ...TestOption) *StringSchema[T] {
	for _, opt := range opts {
		opt(&t)
	}

	t.ValidateFunc = customTestBackwardsCompatWrapper(t.ValidateFunc)
	return v.addTest(t)
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *StringSchema[T]) TestFunc(testFunc p.TestFunc, options ...TestOption) *StringSchema[T] {
	test := TestFunc("", testFunc)
	v.Test(test, options...)
	return v
}

// Test: checks that the value is one of the enum values
func (v *StringSchema[T]) OneOf(enum []T, options ...TestOption) *StringSchema[T] {
	t := p.In(enum)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is at least n characters long
func (v *StringSchema[T]) Min(n int, options ...TestOption) *StringSchema[T] {
	t := p.LenMin[T](n)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is at most n characters long
func (v *StringSchema[T]) Max(n int, options ...TestOption) *StringSchema[T] {
	t := p.LenMax[T](n)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is exactly n characters long
func (v *StringSchema[T]) Len(n int, options ...TestOption) *StringSchema[T] {
	t := p.Len[T](n)
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is a valid email address
func (v *StringSchema[T]) Email(options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeEmail,
		ValidateFunc: func(v any, ctx Ctx) bool {
			email, ok := v.(*T)
			if !ok {
				return false
			}
			return emailRegex.MatchString(string(*email))
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value is a valid URL
func (v *StringSchema[T]) URL(options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeURL,
		ValidateFunc: func(v any, ctx Ctx) bool {
			s, ok := v.(*T)
			if !ok {
				return false
			}
			u, err := url.Parse(string(*s))
			return err == nil && u.Scheme != "" && u.Host != ""
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value has the prefix
func (v *StringSchema[T]) HasPrefix(s T, options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeHasPrefix,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*T)
			if !ok {
				return false
			}
			return strings.HasPrefix(string(*val), string(s))
		},
	}
	t.Params[zconst.IssueCodeHasPrefix] = s
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value has the suffix
func (v *StringSchema[T]) HasSuffix(s T, options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeHasSuffix,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*T)
			if !ok {
				return false
			}
			return strings.HasSuffix(string(*val), string(s))
		},
	}
	t.Params[zconst.IssueCodeHasSuffix] = s
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value contains the substring
func (v *StringSchema[T]) Contains(sub T, options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeContains,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*T)
			if !ok {
				return false
			}
			return strings.Contains(string(*val), string(sub))
		},
	}
	t.Params[zconst.IssueCodeContains] = sub
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that the value contains an uppercase letter
func (v *StringSchema[T]) ContainsUpper(options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeContainsUpper,
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*T)
			if !ok {
				return false
			}
			for _, r := range string(*val) {
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
func (v *StringSchema[T]) ContainsDigit(options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeContainsDigit,
		ValidateFunc: func(v any, ctx Ctx) bool {
			val, ok := v.(*T)
			if !ok {
				return false
			}
			for _, r := range string(*val) {
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
func (v *StringSchema[T]) ContainsSpecial(options ...TestOption) *StringSchema[T] {
	t :=
		Test{
			IssueCode: zconst.IssueCodeContainsSpecial,
			ValidateFunc: func(v any, ctx Ctx) bool {
				val, ok := v.(*T)
				if !ok {
					return false
				}
				for _, r := range string(*val) {
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
func (v *StringSchema[T]) UUID(options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeUUID,
		ValidateFunc: func(v any, ctx Ctx) bool {
			uuid, ok := v.(*T)
			if !ok {
				return false
			}
			return uuidRegex.MatchString(string(*uuid))
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: checks that value matches to regex
func (v *StringSchema[T]) Match(regex *regexp.Regexp, options ...TestOption) *StringSchema[T] {
	t := Test{
		IssueCode: zconst.IssueCodeMatch,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx Ctx) bool {
			s, ok := v.(*T)
			if !ok {
				return false
			}
			return regex.MatchString(string(*s))
		},
	}
	t.Params[zconst.IssueCodeMatch] = regex.String()
	for _, opt := range options {
		opt(&t)
	}
	return v.addTest(t)
}

// Test: nots the next test fn
func (v *StringSchema[T]) Not() NotStringSchema[T] {
	v.isNot = !v.isNot
	return v
}

func (v *StringSchema[T]) addTest(t p.Test) *StringSchema[T] {
	v.tests = p.AddTest(v.tests, t, v.isNot)
	if v.isNot {
		v.isNot = false
	}

	return v
}

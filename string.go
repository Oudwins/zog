package zog

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/internals/is"
	"github.com/Oudwins/zog/zconst"
)

var (
	_ PrimitiveZogSchema[string] = (*StringSchema[string])(nil)
	_ NotStringSchema[string]    = (*StringSchema[string])(nil)
)

type likeString interface {
	~string
}

type NotStringSchema[T likeString] interface {
	OneOf(enum []T, options ...TestOption) *StringSchema[T]
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

	// `Test` method is missing here as we require the user to define their own test for their use case.
	// `Not` method is missing here as we do not want the user to do `Not` chaining.
	// `NotNil`, `Min`, `Max` methods are not included as they are opposites of each other.
}

type StringSchema[T likeString] struct {
	processors []p.ZProcessor[*T]
	defaultVal *T
	required   *p.Test[*T]
	catch      *T
	coercer    CoercerFunc
	isNot      bool
}

// ! INTERNALS

// Returns the type of the schema
func (v *StringSchema[T]) getType() zconst.ZogType {
	return zconst.TypeString
}

// Sets the coercer for the schema
func (v *StringSchema[T]) setCoercer(c CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

func StringLike[T likeString](opts ...SchemaOption) *StringSchema[T] {
	s := &StringSchema[T]{
		coercer: func(val any) (any, error) {
			v, err := conf.Coercers.String(val)
			if err != nil {
				return nil, err
			}
			return T(v.(string)), nil
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Returns a new String Shape
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
	primitiveParsing(ctx, v.processors, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
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
	primitiveValidation(ctx, v.processors, v.defaultVal, v.required, v.catch)
}

// Transform: trims the input data of whitespace if it is a string
func (v *StringSchema[T]) Trim() *StringSchema[T] {
	v.processors = append(v.processors, &p.TransformProcessor[*T]{
		Transform: func(val *T, ctx Ctx) error {
			*val = T(strings.TrimSpace(string(*val)))
			return nil
		},
	})

	return v
}

// Adds a transform function to the schema. Runs in the order it is called
func (v *StringSchema[T]) Transform(transform p.Transform[*T]) *StringSchema[T] {
	v.processors = append(v.processors, &p.TransformProcessor[*T]{Transform: transform})
	return v
}

// ! MODIFIERS

// marks field as required
func (v *StringSchema[T]) Required(options ...TestOption) *StringSchema[T] {
	r := p.Required[*T]()
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

// ! Tests
// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *StringSchema[T]) Test(t Test[*T]) *StringSchema[T] {
	x := p.Test[*T](t)
	v.processors = append(v.processors, &x)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *StringSchema[T]) TestFunc(testFunc BoolTFunc[*T], options ...TestOption) *StringSchema[T] {
	test := p.NewTestFunc("", p.BoolTFunc[*T](testFunc), options...)
	v.Test(Test[*T](*test))
	return v
}

// Test: checks that the value is one of the enum values
func (v *StringSchema[T]) OneOf(enum []T, options ...TestOption) *StringSchema[T] {
	t, fn := p.In(enum)
	return v.addTest(t, fn, options...)
}

// Test: checks that the value is at least n characters long
func (v *StringSchema[T]) Min(n int, options ...TestOption) *StringSchema[T] {
	t, fn := p.LenMin[T](n)
	return v.addTest(t, fn, options...)
}

// Test: checks that the value is at most n characters long
func (v *StringSchema[T]) Max(n int, options ...TestOption) *StringSchema[T] {
	t, fn := p.LenMax[T](n)
	return v.addTest(t, fn, options...)
}

// Test: checks that the value is exactly n characters long
func (v *StringSchema[T]) Len(n int, options ...TestOption) *StringSchema[T] {
	t, fn := p.Len[T](n)
	return v.addTest(t, fn, options...)
}

// Test: checks that the value is a valid email address
func (v *StringSchema[T]) Email(options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeEmail}
	fn := func(v *T, ctx Ctx) bool {
		return is.Email(string(*v))
	}
	return v.addTest(t, fn, options...)
}

// Test: checks that the value is a valid URL
func (v *StringSchema[T]) URL(options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeURL}
	fn := func(v *T, ctx Ctx) bool {
		u, err := url.Parse(string(*v))
		return err == nil && u.Scheme != "" && u.Host != ""
	}
	return v.addTest(t, fn, options...)
}

// Test: checks that the value has the prefix
func (v *StringSchema[T]) HasPrefix(s T, options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeHasPrefix, Params: make(map[string]any, 1)}
	t.Params[zconst.IssueCodeHasPrefix] = string(s)
	fn := func(v *T, ctx Ctx) bool {
		return strings.HasPrefix(string(*v), string(s))
	}
	return v.addTest(t, fn, options...)
}

// Test: checks that the value has the suffix
func (v *StringSchema[T]) HasSuffix(s T, options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{
		IssueCode: zconst.IssueCodeHasSuffix,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeHasSuffix] = string(s)
	fn := func(v *T, ctx Ctx) bool {
		return strings.HasSuffix(string(*v), string(s))
	}
	return v.addTest(t, fn, options...)
}

// Test: checks that the value contains the substring
func (v *StringSchema[T]) Contains(sub T, options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{
		IssueCode: zconst.IssueCodeContains,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeContains] = string(sub)
	fn := func(v *T, ctx Ctx) bool {
		return strings.Contains(string(*v), string(sub))
	}
	return v.addTest(t, fn, options...)
}

// Test: checks that the value contains an uppercase letter
func (v *StringSchema[T]) ContainsUpper(options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeContainsUpper}
	fn := func(v *T, ctx Ctx) bool {
		for _, r := range string(*v) {
			if r >= 'A' && r <= 'Z' {
				return true
			}
		}
		return false
	}

	return v.addTest(t, fn, options...)
}

// Test: checks that the value contains a digit
func (v *StringSchema[T]) ContainsDigit(options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeContainsDigit}
	fn := func(v *T, ctx Ctx) bool {
		for _, r := range string(*v) {
			if r >= '0' && r <= '9' {
				return true
			}
		}
		return false
	}

	return v.addTest(t, fn, options...)
}

// Test: checks that the value contains a special character
func (v *StringSchema[T]) ContainsSpecial(options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeContainsSpecial}
	fn := func(v *T, ctx Ctx) bool {
		for _, r := range string(*v) {
			if (r >= '!' && r <= '/') ||
				(r >= ':' && r <= '@') ||
				(r >= '[' && r <= '`') ||
				(r >= '{' && r <= '~') {
				return true
			}
		}
		return false
	}

	return v.addTest(t, fn, options...)
}

// Test: checks that the value is a valid uuid
func (v *StringSchema[T]) UUID(options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeUUID}
	fn := func(v *T, ctx Ctx) bool {
		return is.UUIDv4(string(*v))
	}

	return v.addTest(t, fn, options...)
}

// Test: checks that value matches to regex
func (v *StringSchema[T]) Match(regex *regexp.Regexp, options ...TestOption) *StringSchema[T] {
	t := p.Test[*T]{IssueCode: zconst.IssueCodeMatch, Params: make(map[string]any, 1)}
	t.Params[zconst.IssueCodeMatch] = regex.String()
	fn := func(v *T, ctx Ctx) bool {
		return regex.MatchString(string(*v))
	}

	return v.addTest(t, fn, options...)
}

// Not returns a schema that negates the next validation test.
// For example, `z.String().Not().Email()` validates that the string is NOT a valid email.
// Note: The negation only applies to the next validation test and is reset afterward.
func (v *StringSchema[T]) Not() NotStringSchema[T] {
	v.isNot = true
	return v
}

func (v *StringSchema[T]) addTest(t p.Test[*T], fn p.BoolTFunc[*T], options ...TestOption) *StringSchema[T] {
	if v.isNot {
		p.TestNotFuncFromBool(fn, &t)
		t.IssueCode = zconst.NotIssueCode(t.IssueCode)
		v.isNot = false
	} else {
		p.TestFuncFromBool(fn, &t)
	}

	for _, opt := range options {
		opt(&t)
	}

	v.processors = append(v.processors, &t)
	return v
}

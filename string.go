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
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	uuidRegex  = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
)

type stringProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *string
	required       *p.Test
	catch          *string
	coercer        conf.CoercerFunc
}

// ! INTERNALS

// Internal function to process the data
func (v *stringProcessor) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {
	primitiveProcessor(val, dest, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Returns the type of the schema
func (v *stringProcessor) getType() zconst.ZogType {
	return zconst.TypeString
}

// Sets the coercer for the schema
func (v *stringProcessor) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

// Returns a new String Schema
func String(opts ...SchemaOption) *stringProcessor {
	s := &stringProcessor{
		coercer: conf.Coercers.String, // default coercer
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Parses the data into the destination string. Returns a list of errors
func (v *stringProcessor) Parse(data any, dest *string, options ...ParsingOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.process(data, dest, path, ctx)

	return errs.List
}

// Adds pretransform function to schema
func (v *stringProcessor) PreTransform(transform p.PreTransform) *stringProcessor {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// PreTransform: trims the input data of whitespace if it is a string
func (v *stringProcessor) Trim() *stringProcessor {
	v.preTransforms = append(v.preTransforms, func(val any, ctx ParseCtx) (any, error) {
		s, ok := val.(string)
		if !ok {
			return val, nil
		}
		return strings.TrimSpace(s), nil
	})
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
	r := p.Required()
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

// ! PRETRANSFORMS

// ! Tests
// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *stringProcessor) Test(t p.Test, opts ...TestOption) *stringProcessor {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value is one of the enum values
func (v *stringProcessor) OneOf(enum []string, options ...TestOption) *stringProcessor {
	t := p.In(enum)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value is at least n characters long
func (v *stringProcessor) Min(n int, options ...TestOption) *stringProcessor {
	t := p.LenMin[string](n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value is at most n characters long
func (v *stringProcessor) Max(n int, options ...TestOption) *stringProcessor {
	t := p.LenMax[string](n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value is exactly n characters long
func (v *stringProcessor) Len(n int, options ...TestOption) *stringProcessor {
	t := p.Len[string](n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value is a valid email address
func (v *stringProcessor) Email(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeEmail,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
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

// Test: checks that the value is a valid URL
func (v *stringProcessor) URL(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeURL,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			u, err := url.Parse(s)
			return err == nil && u.Scheme != "" && u.Host != ""
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value has the prefix
func (v *stringProcessor) HasPrefix(s string, options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeHasPrefix,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.HasPrefix(val, s)
		},
	}
	t.Params[zconst.ErrCodeHasPrefix] = s
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value has the suffix
func (v *stringProcessor) HasSuffix(s string, options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeHasSuffix,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.HasSuffix(val, s)
		},
	}
	t.Params[zconst.ErrCodeHasSuffix] = s
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value contains the substring
func (v *stringProcessor) Contains(sub string, options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeContains,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(string)
			if !ok {
				return false
			}
			return strings.Contains(val, sub)
		},
	}
	t.Params[zconst.ErrCodeContains] = sub
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value contains an uppercase letter
func (v *stringProcessor) ContainsUpper(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeContainsUpper,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
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

// Test: checks that the value contains a digit
func (v *stringProcessor) ContainsDigit(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeContainsDigit,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
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
	}

	for _, opt := range options {
		opt(&t)
	}

	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value contains a special character
func (v *stringProcessor) ContainsSpecial(options ...TestOption) *stringProcessor {
	t :=
		p.Test{
			ErrCode: zconst.ErrCodeContainsSpecial,
			ValidateFunc: func(v any, ctx ParseCtx) bool {
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
		}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that the value is a valid uuid
func (v *stringProcessor) UUID(options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeUUID,
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			uuid, ok := v.(string)
			if !ok {
				return false
			}
			return uuidRegex.MatchString(uuid)
		},
	}
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Test: checks that value matches to regex
func (v *stringProcessor) Match(regex *regexp.Regexp, options ...TestOption) *stringProcessor {
	t := p.Test{
		ErrCode: zconst.ErrCodeMatch,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			s, ok := v.(string)
			if !ok {
				return false
			}
			return regex.MatchString(s)
		},
	}
	t.Params[zconst.ErrCodeMatch] = regex.String()
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

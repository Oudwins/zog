package zog

import (
	"fmt"
	"regexp"
	"strings"

	p "github.com/Oudwins/zog/primitives"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	// TODO improve this regex?
	urlRegex = regexp.MustCompile(`^(http(s)?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?$`)
)

type StringValidator struct {
	Rules []p.Rule
}

func String() *StringValidator {
	return &StringValidator{
		Rules: []p.Rule{
			p.IsType[string]("should be a string"),
		},
	}
}

func (v *StringValidator) Parse(val any) (any, []string, bool) {
	errs, ok := p.GenericRulesValidator(val, v.Rules)
	return val, errs, ok
}

func (v *StringValidator) Optional() *optional {
	return Optional(v)
}

func (v *StringValidator) Default(val any) *defaulter {
	return Default(val, v)
}

func (v *StringValidator) Catch(val any) *catcher {
	return Catch(val, v)
}

func (v *StringValidator) Transform(transform func(val any) (any, bool)) *transformer {
	return Transform(v, transform)
}

func (v *StringValidator) Refine(ruleName string, errorMsg string, validateFunc p.RuleValidateFunc) *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         ruleName,
			ErrorMessage: errorMsg,
			ValidateFunc: validateFunc,
		},
	)

	return v
}

// METHODS

func (v *StringValidator) Min(n int, msg ...string) *StringValidator {
	defaultMsg := fmt.Sprintf("should be at least %d characters long", n)
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	v.Rules = append(v.Rules,
		p.LenMin[string](n, defaultMsg))
	return v
}

func (v *StringValidator) Max(n int) *StringValidator {
	v.Rules = append(v.Rules,
		p.LenMax[string](n, fmt.Sprintf("should be at most %d characters long", n)))
	return v
}
func (v *StringValidator) Len(n int) *StringValidator {
	v.Rules = append(v.Rules,
		p.Length[string](n, fmt.Sprintf("should be exactly %d characters long", n)),
	)
	return v
}

// THIS IS ONLY HERE FOR CREATING ERROR MSGS FOR FORMS. DOESN'T ACTUALLY PROVIDE ANY VALUE
func (v *StringValidator) Required() *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name: "required",
			ValidateFunc: func(rule p.Rule) bool {
				str, ok := rule.FieldValue.(string)
				if !ok {
					return false
				}
				return str != ""
			},
			ErrorMessage: "is a required field",
		},
	)
	return v
}

func (v *StringValidator) Email() *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         "email",
			ErrorMessage: "is not a valid email address",
			ValidateFunc: func(set p.Rule) bool {
				email, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return emailRegex.MatchString(email)
			},
		},
	)
	return v
}

func (v *StringValidator) URL() *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         "url",
			ErrorMessage: "is not a valid url",
			ValidateFunc: func(set p.Rule) bool {
				u, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				isOk := urlRegex.MatchString(u)
				return isOk
			},
		},
	)
	return v
}

// Should use the go method name for this? HasPrefix & HasSuffix???
// TODO ???
func (v *StringValidator) StartsWith(s string) *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:      "startsWith",
			RuleValue: s,
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return strings.HasPrefix(val, s)
			},
			ErrorMessage: fmt.Sprintf("should start with %s", s),
		},
	)
	return v
}

func (v *StringValidator) EndsWith(s string) *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:      "startsWith",
			RuleValue: s,
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return strings.HasSuffix(val, s)
			},
			ErrorMessage: fmt.Sprintf("should end with %s", s),
		},
	)
	return v
}

func (v *StringValidator) Contains(sub string) *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:      "contains",
			RuleValue: sub,
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return strings.Contains(val, sub)
			},
			ErrorMessage: fmt.Sprintf("should contain %s", sub),
		},
	)
	return v
}

func (v *StringValidator) ContainsUpper() *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name: "containsUpper",
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(string)
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
			ErrorMessage: "should contain at least one uppercase letter",
		},
	)
	return v
}

func (v *StringValidator) ContainsDigit() *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name: "containsDigit",
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(string)
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
			ErrorMessage: "should contain at least one digit",
		},
	)
	return v
}

func (v *StringValidator) ContainsSpecial() *StringValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name: "containsSpecial",
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(string)
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
			ErrorMessage: "should contain at least one special character",
		},
	)
	return v
}

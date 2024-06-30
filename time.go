package zog

import (
	"fmt"
	"time"

	p "github.com/Oudwins/zog/primitives"
)

type timeValidator struct {
	Rules      []p.Rule
	IsOptional bool
}

func Time() *timeValidator {
	return &timeValidator{
		Rules: []p.Rule{
			p.IsType[time.Time]("is not a a valid time"),
		},
	}
}

func (v *timeValidator) Parse(val any) (any, []string, bool) {
	errs, ok := p.GenericRulesValidator(val, v.Rules)
	return val, errs, ok
}

func (v *timeValidator) Optional() *optional {
	return Optional(v)
}

func (v *timeValidator) Default(val any) *defaulter {
	return Default(val, v)
}
func (v *timeValidator) Catch(val any) *catcher {
	return Catch(val, v)
}
func (v *timeValidator) Transform(transform func(val any) (any, bool)) *transformer {
	return Transform(v, transform)
}

// GLOBAL METHODS

func (v *timeValidator) Refine(ruleName string, errorMsg string, validateFunc p.RuleValidateFunc) *timeValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         ruleName,
			ErrorMessage: errorMsg,
			ValidateFunc: validateFunc,
		},
	)

	return v
}

// UNIQUE METHODS

func (v *timeValidator) After(t time.Time) *timeValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         "timeAfter",
			ErrorMessage: fmt.Sprintf("is not after %v", t),
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(time.Time)
				if !ok {
					return false
				}
				return val.After(t)
			},
		},
	)
	return v
}

func (v *timeValidator) Before(t time.Time) *timeValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         "timeBefore",
			ErrorMessage: fmt.Sprintf("is not before %v", t),
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(time.Time)
				if !ok {
					return false
				}
				return val.Before(t)
			},
		},
	)
	return v
}

func (v *timeValidator) Is(t time.Time) *timeValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         "timeIs",
			ErrorMessage: fmt.Sprintf("is not %v", t),
			ValidateFunc: func(set p.Rule) bool {
				val, ok := set.FieldValue.(time.Time)
				if !ok {
					return false
				}
				return val.Equal(t)
			},
		},
	)

	return v
}

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

func (v *timeValidator) In(values []time.Time) *timeValidator {
	v.Rules = append(v.Rules, p.In(values, fmt.Sprintf("should be in %v", values)))
	return v
}

func (v *timeValidator) Optional() *timeValidator {
	v.IsOptional = true
	return v
}

func (v *timeValidator) Validate(fieldValue any) ([]string, bool) {
	return p.GenericValidator(fieldValue, v.Rules, v.IsOptional)
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

package zog

import (
	"fmt"

	p "github.com/Oudwins/zog/primitives"
)

type EnumValidator struct {
	Rules []p.Rule
}

func Enum[T any](vals []T) *EnumValidator {
	return &EnumValidator{
		Rules: []p.Rule{
			p.In(vals, fmt.Sprintf("should be in %v", vals)),
		},
	}
}

func (v *EnumValidator) Parse(fieldValue any) (any, []string, bool) {
	errs, ok := p.GenericRulesValidator(fieldValue, v.Rules)
	return fieldValue, errs, ok
}

// GLOBAL METHODS

func (v *EnumValidator) Optional() *optional {
	return Optional(v)
}

func (v *EnumValidator) Default(val any) *defaulter {
	return Default(val, v)
}

func (v *EnumValidator) Catch(val any) *catcher {
	return Catch(val, v)
}

func (v *EnumValidator) Transform(transform func(val any) (any, bool)) *transformer {
	return Transform(v, transform)
}

func (v *EnumValidator) Refine(ruleName string, errorMsg string, validateFunc p.RuleValidateFunc) *EnumValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         ruleName,
			ErrorMessage: errorMsg,
			ValidateFunc: validateFunc,
		},
	)
	return v
}

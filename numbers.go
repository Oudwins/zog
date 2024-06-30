package zog

import (
	"fmt"

	p "github.com/Oudwins/zog/primitives"
)

type Numeric interface {
	~int | ~float64
}

type numberValidator[T Numeric] struct {
	Rules []p.Rule
}

func Float() *numberValidator[float64] {
	return &numberValidator[float64]{
		Rules: []p.Rule{
			p.IsType[float64]("should be a decimal number"),
		},
	}
}

func Int() *numberValidator[int] {
	return &numberValidator[int]{
		Rules: []p.Rule{
			p.IsType[int]("should be an whole number"),
		},
	}
}

func (v *numberValidator[Numeric]) Parse(fieldValue any) (any, []string, bool) {
	errs, ok := p.GenericRulesValidator(fieldValue, v.Rules)
	return fieldValue, errs, ok
}

// GLOBAL METHODS

func (v *numberValidator[Numeric]) Optional() *optional {
	return Optional(v)
}

func (v *numberValidator[Numeric]) Default(val any) *defaulter {
	return Default(val, v)
}
func (v *numberValidator[Numeric]) Catch(val any) *catcher {
	return Catch(val, v)
}
func (v *numberValidator[Numeric]) Transform(transform func(val any) (any, bool)) *transformer {
	return Transform(v, transform)
}

func (v *numberValidator[T]) Refine(ruleName string, errorMsg string, validateFunc p.RuleValidateFunc) *numberValidator[T] {
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

func (v *numberValidator[Numeric]) EQ(n Numeric) *numberValidator[Numeric] {
	v.Rules = append(v.Rules, p.EQ(n, fmt.Sprintf("should be equal to %v", n)))
	return v
}

func (v *numberValidator[Numeric]) LTE(n Numeric) *numberValidator[Numeric] {
	v.Rules = append(v.Rules, p.LTE(n, fmt.Sprintf("should be lesser or equal than %v", n)))
	return v
}

func (v *numberValidator[Numeric]) GTE(n Numeric) *numberValidator[Numeric] {
	v.Rules = append(v.Rules, p.GTE(n, fmt.Sprintf("should be greater or equal to %v", n)))
	return v
}

func (v *numberValidator[Numeric]) LT(n Numeric) *numberValidator[Numeric] {
	v.Rules = append(v.Rules, p.LT(n, fmt.Sprintf("should be less than %v", n)))
	return v
}

func (v *numberValidator[Numeric]) GT(n Numeric) *numberValidator[Numeric] {
	v.Rules = append(v.Rules, p.GT(n, fmt.Sprintf("should be greater than %v", n)))
	return v
}

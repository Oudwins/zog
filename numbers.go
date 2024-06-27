package zog

import (
	"fmt"

	p "github.com/Oudwins/zog/primitives"
)

type Numeric interface {
	~int | ~float64
}

type numberValidator[T Numeric] struct {
	Rules      []p.Rule
	IsOptional bool
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

// GLOBAL METHODS

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

// is equal to one of the values
func (v *numberValidator[T]) In(values []T) *numberValidator[T] {
	v.Rules = append(v.Rules, p.In(values, fmt.Sprintf("should be in %v", values)))
	return v
}

func (v *numberValidator[Numeric]) Optional() *numberValidator[Numeric] {
	v.IsOptional = true
	return v
}

func (v *numberValidator[Numeric]) Validate(fieldValue any) ([]string, bool) {
	return p.GenericValidator(fieldValue, v.Rules, v.IsOptional)
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

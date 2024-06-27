package zog

import (
	"github.com/Oudwins/zog/primitives"
)

type boolValidator struct {
	Rules      []primitives.Rule
	IsOptional bool
}

func Bool() *boolValidator {
	return &boolValidator{
		Rules: []primitives.Rule{
			primitives.IsType[bool]("is not a valid boolean"),
		},
	}
}

func (v *boolValidator) Validate(fieldValue any) ([]string, bool) {
	return primitives.GenericValidator(fieldValue, v.Rules, v.IsOptional)
}

func (v *boolValidator) Optional() *boolValidator {
	v.IsOptional = true
	return v
}

func (v *boolValidator) True() *boolValidator {
	v.Rules = append(v.Rules, primitives.EQ[bool](true, "should be true"))
	return v
}

func (v *boolValidator) False() *boolValidator {
	v.Rules = append(v.Rules, primitives.EQ[bool](false, "should be false"))
	return v
}

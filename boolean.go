package zog

import (
	"github.com/Oudwins/zog/primitives"
)

type boolValidator struct {
	Rules []primitives.Rule
}

func Bool() *boolValidator {
	return &boolValidator{
		Rules: []primitives.Rule{
			primitives.IsType[bool]("is not a valid boolean"),
		},
	}
}

func (v *boolValidator) Parse(val any) (any, []string, bool) {
	err, ok := primitives.GenericRulesValidator(val, v.Rules)
	return val, err, ok
}

// GLOBAL METHODS

func (v *boolValidator) Optional() *optional {
	return Optional(v)
}

func (v *boolValidator) Default(val bool) *defaulter {
	return Default(val, v)
}

func (v *boolValidator) Catch(val bool) *catcher {
	return Catch(val, v)
}

func (v *boolValidator) Transform(transform func(val any) (any, bool)) *transformer {
	return Transform(v, transform)
}

// UNIQUE METHODS

func (v *boolValidator) True() *boolValidator {
	v.Rules = append(v.Rules, primitives.EQ[bool](true, "should be true"))
	return v
}

func (v *boolValidator) False() *boolValidator {
	v.Rules = append(v.Rules, primitives.EQ[bool](false, "should be false"))
	return v
}

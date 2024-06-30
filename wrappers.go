package zog

import (
	p "github.com/Oudwins/zog/primitives"
)

type optional struct {
	validator fieldParser
}

func Optional(value fieldParser) *optional {
	return &optional{
		validator: value,
	}
}

func (o *optional) Parse(val any) (any, []string, bool) {
	newVal, errs, ok := o.validator.Parse(val)
	if p.IsZeroValue(val) {
		return newVal, nil, true
	}

	return newVal, errs, ok
}

//

type defaulter struct {
	validator  fieldParser
	defaultVal any
}

func Default(val any, validator fieldParser) *defaulter {
	return &defaulter{
		validator:  validator,
		defaultVal: val,
	}
}

func (d *defaulter) Parse(val any) (any, []string, bool) {
	if p.IsZeroValue(val) {
		return d.defaultVal, nil, true
	}
	newVal, errs, ok := d.validator.Parse(val)

	return newVal, errs, ok

}

type catcher struct {
	validator fieldParser
	value     any
}

func Catch(value any, validator fieldParser) *catcher {
	return &catcher{
		validator: validator,
		value:     value,
	}
}

func (c *catcher) Parse(val any) (any, []string, bool) {
	newVal, errs, ok := c.validator.Parse(val)
	if !ok {
		return c.value, nil, true
	}

	return newVal, errs, ok
}

type transformer struct {
	validator fieldParser
	transform func(val any) (any, bool)
}

func Transform(validator fieldParser, transform func(val any) (any, bool)) *transformer {
	return &transformer{
		validator: validator,
		transform: transform,
	}
}

func (t *transformer) Parse(val any) (any, []string, bool) {
	newVal, errs, ok := t.validator.Parse(val)
	if !ok {
		return val, errs, ok
	}
	transformedVal, ok := t.transform(newVal)

	if !ok {
		return val, nil, ok
	}

	return transformedVal, errs, ok
}

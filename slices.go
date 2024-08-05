package zog

import (
	"fmt"
	"reflect"

	p "github.com/Oudwins/zog/primitives"
)

type sliceProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	schema         Processor
	postTransforms []p.PostTransform
	required       *p.Test
	defaultVal     any
	catch          any
}

func Slice(schema Processor) *sliceProcessor {
	return &sliceProcessor{
		schema: schema,
		tests:  []p.Test{},
	}
}

// only supports val = slice[any] & dest = &slice[]
func (v *sliceProcessor) Parse(val any, dest any) p.ZogSchemaErrors {
	var ctx = p.NewParseCtx()
	errs := p.NewErrsMap()
	path := p.Pather("")
	v.process(val, dest, errs, path, ctx)

	return errs.M
}

func (v *sliceProcessor) process(val any, dest any, errs p.ZogErrors, path p.Pather, ctx *p.ParseCtx) {
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(val, ctx)
			// bail if error in preTransform
			if err != nil {
				errs.Add(path, Errors.WrapUnknown(err))
				return
			}
			val = nVal
		}
	}

	// 4. postTransforms
	if v.postTransforms != nil {
		defer func() {
			for _, fn := range v.postTransforms {
				fn(dest, ctx)
			}
		}()

	}

	// 2. cast data to string & handle default/required
	isZeroVal := p.IsZeroValue(val)

	if isZeroVal {
		if v.defaultVal != nil {
			reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v.defaultVal))
		} else if v.required == nil {
			return
		}
	} else {
		// make sure val is a slice if not try to make it one
		valTp := reflect.TypeOf(val)
		if valTp.Kind() != reflect.Slice {
			// TODO coerce to slice
		}
		// fill dest slice with default value for its type up to size of val
		// call processor on each item while iterating over val & dest

		// 1. val is a slice of
		// 1.1 right type -> copy it
		// 1.2 wrong type -> coerce each element & copy it
		// 2. val is not a slice
		// 2.1 its a string -> split by comma & coerce each element & copy it
		// 2.2 ??
	}

	// required
	if v.required != nil && !v.required.ValidateFunc(dest, ctx) {
		errs.Add(path, Errors.New(v.required.ErrorFunc(dest, ctx)))
	}

	// 3. tests for slice
	for _, test := range v.tests {
		if !test.ValidateFunc(dest, ctx) {
			// catching the first error if catch is set
			if v.catch != nil {
				dest = v.catch
				break
			}
			//
			errs.Add(path, Errors.New(test.ErrorFunc(dest, ctx)))
		}
	}

	// 3.1 tests for slice items
	if v.schema == nil {
		rv := reflect.ValueOf(dest).Elem()
		for idx := 0; idx < rv.Len(); idx++ {
			item := rv.Index(idx).Interface()
			ptr := rv.Index(idx).Addr().Interface()
			p := path.Push(fmt.Sprint(idx))
			v.schema.process(item, ptr, errs, p, ctx)
		}
	}
	// 4. postTransforms -> defered see above
}

// type sliceValidator struct {
// 	Rules []p.Rule
// }

// func Slice(schema fieldParser) *sliceValidator {
// 	return &sliceValidator{
// 		Rules: []p.Rule{
// 			{
// 				Name:      "sliceItemsMatchSchema",
// 				RuleValue: schema,
// 				// TODO this should really be improved. Maybe grab the error message from the schema?
// 				ErrorMessage: "all items should match the schema",
// 				ValidateFunc: func(set p.Rule) bool {
// 					rv := reflect.ValueOf(set.FieldValue)
// 					if rv.Kind() != reflect.Slice {
// 						return false
// 					}
// 					for idx := 0; idx < rv.Len(); idx++ {
// 						v := rv.Index(idx).Interface()
// 						newVal, _, ok := schema.Parse(v)
// 						if !ok {
// 							return false
// 						}
// 						if !reflect.DeepEqual(v, newVal) {
// 							rv.Index(idx).Set(reflect.ValueOf(newVal))
// 						}
// 					}
// 					return true
// 				},
// 			},
// 		},
// 	}
// }

// // GLOBAL METHODS

// func (v *sliceValidator) Refine(ruleName string, errorMsg string, validateFunc p.RuleValidateFunc) *sliceValidator {
// 	v.Rules = append(v.Rules,
// 		p.Rule{
// 			Name:         ruleName,
// 			ErrorMessage: errorMsg,
// 			ValidateFunc: validateFunc,
// 		},
// 	)

// 	return v
// }

// func (v *sliceValidator) Optional() *optional {
// 	return Optional(v)
// }

// // Current implementation is not working. Need to fix.
// // func (v *sliceValidator) Default(val any) *defaulter {
// // 	return Default(val, v)
// // }
// // func (v *sliceValidator) Catch(val any) *catcher {
// // 	return Catch(val, v)
// // }
// // func (v *sliceValidator) Transform(transform func(val any) (any, bool)) *transformer {
// // 	return Transform(v, transform)
// // }

// func (v *sliceValidator) Parse(fieldValue any) (any, []string, bool) {
// 	errs, ok := p.GenericRulesValidator(fieldValue, v.Rules)
// 	return nil, errs, ok
// }

// // UNIQUE METHODS

// // TODO
// // some & every -> pass a validator

// // Minimum number of items
// func (v *sliceValidator) Min(n int) *sliceValidator {
// 	v.Rules = append(v.Rules,
// 		sliceMin(n, fmt.Sprintf("should be at least %d items long", n)),
// 	)
// 	return v
// }

// // Maximum number of items
// func (v *sliceValidator) Max(n int) *sliceValidator {
// 	v.Rules = append(v.Rules,
// 		sliceMax(n, fmt.Sprintf("should be at maximum %d items long", n)),
// 	)
// 	return v
// }

// // Exact number of items
// func (v *sliceValidator) Len(n int) *sliceValidator {
// 	v.Rules = append(v.Rules,
// 		sliceLength(n, fmt.Sprintf("should be exactly %d items long", n)),
// 	)
// 	return v
// }

// func (v *sliceValidator) Contains(val any) *sliceValidator {
// 	v.Rules = append(v.Rules,
// 		p.Rule{
// 			Name:         "contains",
// 			RuleValue:    val,
// 			ErrorMessage: fmt.Sprintf("should contain %v", val),
// 			ValidateFunc: func(set p.Rule) bool {
// 				rv := reflect.ValueOf(set.FieldValue)
// 				if rv.Kind() != reflect.Slice {
// 					return false
// 				}
// 				for idx := 0; idx < rv.Len(); idx++ {
// 					v := rv.Index(idx).Interface()

// 					if reflect.DeepEqual(v, val) {
// 						return true
// 					}
// 				}

// 				return false
// 			},
// 		},
// 	)
// 	return v
// }

// func sliceMin(n int, errMsg string) p.Rule {
// 	return p.Rule{
// 		Name:         "sliceMin",
// 		RuleValue:    n,
// 		ErrorMessage: errMsg,
// 		ValidateFunc: func(set p.Rule) bool {
// 			rv := reflect.ValueOf(set.FieldValue)
// 			if rv.Kind() != reflect.Slice {
// 				return false
// 			}
// 			return rv.Len() >= n
// 		},
// 	}
// }
// func sliceMax(n int, errMsg string) p.Rule {
// 	return p.Rule{
// 		Name:         "sliceMax",
// 		RuleValue:    n,
// 		ErrorMessage: errMsg,
// 		ValidateFunc: func(set p.Rule) bool {
// 			rv := reflect.ValueOf(set.FieldValue)
// 			if rv.Kind() != reflect.Slice {
// 				return false
// 			}
// 			return rv.Len() <= n
// 		},
// 	}
// }
// func sliceLength(n int, errMsg string) p.Rule {
// 	return p.Rule{
// 		Name:         "sliceLength",
// 		RuleValue:    n,
// 		ErrorMessage: errMsg,
// 		ValidateFunc: func(set p.Rule) bool {
// 			rv := reflect.ValueOf(set.FieldValue)
// 			if rv.Kind() != reflect.Slice {
// 				return false
// 			}
// 			return rv.Len() == n
// 		},
// 	}
// }

package zog

import (
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// this is an experiment
// import (
// 	"github.com/Oudwins/zog/conf"
// 	p "github.com/Oudwins/zog/p"
// 	"github.com/Oudwins/zog/zconst"
// )

// type CustomInterfaceDefinition struct {
// 	GetValue func(data any) (any, error)
// 	SetValue func(data any, value any) error
// 	TypeName string
// }

// type CustomInterfaceSchema struct {
// 	preTransforms  []p.PreTransform
// 	tests          []p.Test
// 	postTransforms []p.PostTransform
// 	defaultVal     any
// 	required       *p.Test
// 	catch          any
// 	coercer        conf.CoercerFunc
// 	definition     *CustomInterfaceDefinition
// }

// func (v *CustomInterfaceSchema) getType() zconst.ZogType {
// 	return v.definition.TypeName
// }

// func (v *CustomInterfaceSchema) setCoercer(c conf.CoercerFunc) {
// 	v.coercer = c
// }

// func CustomInterface(definition *CustomInterfaceDefinition, opts ...SchemaOption) *CustomInterfaceSchema {
// 	s := &CustomInterfaceSchema{
// 		definition: definition,
// 	}
// 	for _, opt := range opts {
// 		opt(s)
// 	}
// 	return s
// }

// func (v *CustomInterfaceSchema) Parse(data any, destPtr any, options ...ExecOption) ZogIssueList {
// 	errs := p.NewErrsList()
// 	defer errs.Free()

// 	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
// 	defer ctx.Free()
// 	for _, opt := range options {
// 		opt(ctx)
// 	}

// 	path := p.NewPathBuilder()
// 	defer path.Free()
// 	v.process(ctx.NewSchemaCtx(data, destPtr, path, v.getType()))

// 	return errs.List
// }

// func (v *CustomInterfaceSchema) process(ctx *p.SchemaCtx) {
// 	defer ctx.Free()
// 	canCatch := v.catch != nil
// 	// 4. postTransforms
// 	defer func() {
// 		// only run posttransforms on success
// 		if !ctx.HasErrored() {
// 			for _, fn := range v.postTransforms {
// 				err := fn(ctx.Val, ctx)
// 				if err != nil {
// 					ctx.AddIssue(ctx.IssueFromUnknownError(err))
// 					return
// 				}
// 			}
// 		}
// 	}()

// 	// 1. preTransforms
// 	for _, fn := range v.preTransforms {
// 		nVal, err := fn(ctx.Val, ctx)
// 		// bail if error in preTransform
// 		if err != nil {
// 			if canCatch {
// 				v.definition.SetValue(ctx.Val, v.catch)
// 				return
// 			}
// 			ctx.AddIssue(ctx.IssueFromUnknownError(err))
// 			return
// 		}
// 		v.definition.SetValue(ctx.Val, nVal)
// 	}

// 	// 2. cast data to string & handle default/required
// 	// Warning. This uses generic IsZeroValue because for Validate we treat zero values as invalid for required fields. This is different from Parse.
// 	isZeroVal := p.IsZeroValue(ctx.Val)

// 	if isZeroVal {
// 		if v.defaultVal != nil {
// 			v.definition.SetValue(ctx.Val, v.defaultVal)
// 		} else if v.required == nil {
// 			// This handles optional case
// 			return
// 		} else {
// 			// is required & zero value
// 			// required
// 			if v.catch != nil {
// 				v.definition.SetValue(ctx.Val, v.catch)
// 				return
// 			} else {
// 				ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val))
// 				return
// 			}
// 		}
// 	}
// 	// 3. tests
// 	for _, test := range v.tests {
// 		if !test.ValidateFunc(ctx.Val, ctx) {
// 			// catching the first error if catch is set
// 			if canCatch {
// 				v.definition.SetValue(ctx.Val, v.catch)
// 				return
// 			}
// 			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.Val))
// 		}
// 	}

// }

// // Validate Given string
// func (v *CustomInterfaceSchema) Validate(dataPtr any, options ...ExecOption) p.ZogIssueList {
// 	errs := p.NewErrsList()
// 	defer errs.Free()
// 	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
// 	defer ctx.Free()
// 	for _, opt := range options {
// 		opt(ctx)
// 	}

// 	path := p.NewPathBuilder()
// 	defer path.Free()
// 	v.validate(ctx.NewSchemaCtx(dataPtr, dataPtr, path, v.getType()))
// 	return errs.List
// }

// // Internal function to validate the data
// func (v *CustomInterfaceSchema) validate(ctx *p.SchemaCtx) {
// 	defer ctx.Free()
// 	canCatch := v.catch != nil

// 	// 4. postTransforms
// 	defer func() {
// 		// only run posttransforms on success
// 		if !ctx.HasErrored() {
// 			for _, fn := range v.postTransforms {
// 				err := fn(ctx.Val, ctx)
// 				if err != nil {
// 					ctx.AddIssue(ctx.IssueFromUnknownError(err))
// 					return
// 				}
// 			}
// 		}
// 	}()

// 	// 1. preTransforms
// 	for _, fn := range v.preTransforms {
// 		nVal, err := fn(ctx.Val, ctx)
// 		// bail if error in preTransform
// 		if err != nil {
// 			if canCatch {
// 				v.definition.SetValue(ctx.Val, v.catch)
// 				return
// 			}
// 			ctx.AddIssue(ctx.IssueFromUnknownError(err))
// 			return
// 		}
// 		v.definition.SetValue(ctx.Val, nVal)
// 	}

// 	// 2. cast data to string & handle default/required
// 	// Warning. This uses generic IsZeroValue because for Validate we treat zero values as invalid for required fields. This is different from Parse.
// 	isZeroVal := p.IsZeroValue(ctx.Val)

// 	if isZeroVal {
// 		if v.defaultVal != nil {
// 			v.definition.SetValue(ctx.Val, v.defaultVal)
// 		} else if v.required == nil {
// 			// This handles optional case
// 			return
// 		} else {
// 			// is required & zero value
// 			// required
// 			if v.catch != nil {
// 				v.definition.SetValue(ctx.Val, v.catch)
// 				return
// 			} else {
// 				ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val))
// 				return
// 			}
// 		}
// 	}
// 	// 3. tests
// 	for _, test := range v.tests {
// 		if !test.ValidateFunc(ctx.Val, ctx) {
// 			// catching the first error if catch is set
// 			if canCatch {
// 				v.definition.SetValue(ctx.Val, v.catch)
// 				return
// 			}
// 			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.Val))
// 		}
// 	}

// }

/*
type ZogSchema interface {
	process(ctx *p.SchemaCtx)
	validate(ctx *p.SchemaCtx)
	setCoercer(c conf.CoercerFunc)
	getType() zconst.ZogType
}

*/

// either a function or a PublicZogSchema

type CustomSchemaFn = func(val any, ctx Ctx) any // you return the value you want to modify to
var _ ZogSchema = &CustomSchemaFnSchema{}

// We can't make this support withCoercer unless we make TestOption take an interface
func CustomTest(fn BoolTFunc, opts ...TestOption) *CustomTestSchema {
	test := &Test{}
	p.TestFuncFromBool(fn, test)
	for _, opt := range opts {
		opt(test)
	}
	return &CustomTestSchema{test: test}
}

type CustomTestSchema struct {
	test    *Test
	coercer CoercerFunc
}

func (s *CustomTestSchema) process(ctx *p.SchemaCtx) {
	s.test.Func(ctx.Val, ctx)
}

func (s *CustomTestSchema) validate(ctx *p.SchemaCtx) {
	s.test.Func(ctx.Val, ctx)
}

func (s *CustomTestSchema) setCoercer(coercer CoercerFunc) {
	// no op
}

func (s *CustomTestSchema) getType() string {
	return "custom"
}

type CustomSchemaFnSchema struct {
	fn CustomSchemaFn
}

func (s *CustomSchemaFnSchema) process(ctx *p.SchemaCtx) {
	v := s.fn(ctx.Val, ctx)
	if v == nil {
		return
	}
	refVal := reflect.ValueOf(ctx.DestPtr).Elem()
	refVal.Set(reflect.ValueOf(v))
}

func (s *CustomSchemaFnSchema) validate(ctx *p.SchemaCtx) {
	v := s.fn(ctx.Val, ctx)
	if v == nil {
		return
	}
	refVal := reflect.ValueOf(ctx.Val).Elem()
	refVal.Set(reflect.ValueOf(v))
}

func (s *CustomSchemaFnSchema) setCoercer(coercer CoercerFunc) {
	// no op
}

func (s *CustomSchemaFnSchema) getType() string {
	return "custom"
}

func (s *CustomSchemaFnSchema) Parse(data any, destPtr any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	s.process(ctx.NewSchemaCtx(data, destPtr, path, s.getType()))

	return errs.M
}

func (s *CustomSchemaFnSchema) Validate(dataPtr any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	s.validate(ctx.NewValidateSchemaCtx(dataPtr, path, s.getType()))
	return errs.M
}

type customFunc func(fn CustomSchemaFn) *CustomSchemaFnSchema

type PublicZogSchema interface {
	Process(ctx *p.SchemaCtx)
	Validate(ctx *p.SchemaCtx)
	SetCoercer(c conf.CoercerFunc)
	GetType() zconst.ZogType
}

type PublicZogSchemaWrapper struct {
	schema PublicZogSchema
}

func (s *PublicZogSchemaWrapper) process(ctx *p.SchemaCtx) {
	s.schema.Process(ctx)
}

func (s *PublicZogSchemaWrapper) validate(ctx *p.SchemaCtx) {
	s.schema.Validate(ctx)
}

func (s *PublicZogSchemaWrapper) setCoercer(c CoercerFunc) {
	s.schema.SetCoercer(c)
}

func (s *PublicZogSchemaWrapper) getType() zconst.ZogType {
	return s.schema.GetType()
}

func (f customFunc) Schema(s PublicZogSchema) *PublicZogSchemaWrapper {
	return &PublicZogSchemaWrapper{schema: s}
}

var Custom customFunc = func(fn CustomSchemaFn) *CustomSchemaFnSchema {
	return &CustomSchemaFnSchema{fn: fn}
}

// This is an issue bc we have both validate and parse, so either we add all of the inputs or what?
func CustomFn(fn CustomSchemaFn) *CustomSchemaFnSchema {
	return &CustomSchemaFnSchema{fn: fn}
}

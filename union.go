package zog

import (
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/pkgs/internals"
	zss "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/zconst"
)

var _ ZogSchema = &UnionSchema{}

type UnionSchema struct {
	schemas []ZogSchema
}

func EXPERIMENTAL_UNION(schemas []ZogSchema, options ...SchemaOption) *UnionSchema {
	s := &UnionSchema{
		schemas: schemas,
	}
	for _, opt := range options {
		opt(s)
	}
	return s
}

func (u *UnionSchema) Parse(data any, dest any, options ...ExecOption) p.ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, dest, path, u.getType())
	defer sctx.Free()
	u.process(sctx)
	if len(errs.List) == 0 {
		return nil
	}
	return errs.List
}

func (u *UnionSchema) Validate(dest any, options ...ExecOption) p.ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(dest, dest, path, u.getType())
	defer sctx.Free()
	u.validate(sctx)
	if len(errs.List) == 0 {
		return nil
	}
	return errs.List
}

func (u *UnionSchema) process(ctx *p.SchemaCtx) {
	// Wrap the context and only go to the next one on fail. Keeping all the errors and appending at the end
	listStart := len(ctx.Errors.List)
	values := newUnionBranchValues(ctx.Data, ctx.ValPtr, false)
	for _, s := range u.schemas {
		numIssues := len(ctx.Errors.List)
		data, dest := values.next()
		branchCtx := ctx.NewSchemaCtx(data, dest, ctx.Path, s.getType())
		s.process(branchCtx)
		branchCtx.Free()
		if len(ctx.Errors.List) == numIssues {
			values.commit()
			if listStart == 0 {
				ctx.Errors.List = nil
			} else {
				ctx.Errors.List = ctx.Errors.List[:listStart]
			}
			return // success
		}
		snapshotUnionIssueValues(ctx.Errors.List[numIssues:])
		values.rollback()
	}
}

func (u *UnionSchema) validate(ctx *p.SchemaCtx) {
	// Wrap the context and only go to the next one on fail. Keeping all the errors and appending at the end
	listStart := len(ctx.Errors.List)
	values := newUnionBranchValues(nil, ctx.ValPtr, true)
	for _, s := range u.schemas {
		numIssues := len(ctx.Errors.List)
		data, dest := values.next()
		branchCtx := ctx.NewSchemaCtx(data, dest, ctx.Path, s.getType())
		s.validate(branchCtx)
		branchCtx.Free()
		if len(ctx.Errors.List) == numIssues {
			values.commit()
			if listStart == 0 {
				ctx.Errors.List = nil
			} else {
				ctx.Errors.List = ctx.Errors.List[:listStart]
			}
			return // success
		}
		snapshotUnionIssueValues(ctx.Errors.List[numIssues:])
		values.rollback()
	}

}

type unionBranchValues struct {
	input      unionValue
	output     unionValue
	validating bool
}

func newUnionBranchValues(input, output any, validating bool) unionBranchValues {
	return unionBranchValues{
		input:      newUnionValue(input),
		output:     newUnionValue(output),
		validating: validating,
	}
}

func (v *unionBranchValues) next() (any, any) {
	output := v.output.next()
	if v.validating {
		return output, output
	}
	return v.input.next(), output
}

func (v *unionBranchValues) rollback() {
	if !v.validating {
		v.input.rollback()
	}
	v.output.rollback()
}

func (v *unionBranchValues) commit() {
	v.output.commit()
}

type unionValue struct {
	original   reflect.Value
	working    reflect.Value
	hasWorking bool
}

func newUnionValue(value any) unionValue {
	return unionValue{original: reflect.ValueOf(value)}
}

func (v *unionValue) next() any {
	if !v.hasWorking {
		v.working = cloneUnionValue(v.original, make(map[unionCloneVisit]reflect.Value))
		v.hasWorking = true
	}
	return unionValueInterface(v.working)
}

func (v *unionValue) rollback() {
	if reflect.DeepEqual(unionValueInterface(v.original), unionValueInterface(v.working)) {
		return
	}
	v.working = reflect.Value{}
	v.hasWorking = false
}

func (v *unionValue) commit() {
	if !v.original.IsValid() || v.original.Kind() != reflect.Pointer || v.original.IsNil() {
		return
	}
	if !reflect.DeepEqual(v.original.Interface(), v.working.Interface()) {
		v.original.Elem().Set(v.working.Elem())
	}
}

func unionValueInterface(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}
	return value.Interface()
}

func snapshotUnionIssueValues(issues p.ZogIssueList) {
	for _, issue := range issues {
		value := cloneUnionValue(reflect.ValueOf(issue.Value), make(map[unionCloneVisit]reflect.Value))
		issue.Value = unionValueInterface(value)
	}
}

type unionCloneVisit struct {
	typeOf reflect.Type
	ptr    uintptr
	len    int
	cap    int
}

func cloneUnionValue(value reflect.Value, visited map[unionCloneVisit]reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}

	switch value.Kind() {
	case reflect.Interface:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		clone := reflect.New(value.Type()).Elem()
		clone.Set(cloneUnionValue(value.Elem(), visited))
		return clone
	case reflect.Pointer:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		visit := unionCloneVisit{typeOf: value.Type(), ptr: value.Pointer()}
		if clone, ok := visited[visit]; ok {
			return clone
		}
		clone := reflect.New(value.Type().Elem())
		visited[visit] = clone
		clone.Elem().Set(cloneUnionValue(value.Elem(), visited))
		return clone
	case reflect.Slice:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		visit := unionCloneVisit{typeOf: value.Type(), ptr: value.Pointer(), len: value.Len(), cap: value.Cap()}
		if clone, ok := visited[visit]; ok {
			return clone
		}
		clone := reflect.MakeSlice(value.Type(), value.Len(), value.Cap())
		visited[visit] = clone
		for i := 0; i < value.Len(); i++ {
			clone.Index(i).Set(cloneUnionValue(value.Index(i), visited))
		}
		return clone
	case reflect.Map:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		visit := unionCloneVisit{typeOf: value.Type(), ptr: value.Pointer()}
		if clone, ok := visited[visit]; ok {
			return clone
		}
		clone := reflect.MakeMapWithSize(value.Type(), value.Len())
		visited[visit] = clone
		iter := value.MapRange()
		for iter.Next() {
			clone.SetMapIndex(iter.Key(), cloneUnionValue(iter.Value(), visited))
		}
		return clone
	case reflect.Struct:
		clone := reflect.New(value.Type()).Elem()
		clone.Set(value)
		for i := 0; i < value.NumField(); i++ {
			if value.Type().Field(i).PkgPath == "" {
				clone.Field(i).Set(cloneUnionValue(value.Field(i), visited))
			}
		}
		return clone
	case reflect.Array:
		clone := reflect.New(value.Type()).Elem()
		for i := 0; i < value.Len(); i++ {
			clone.Index(i).Set(cloneUnionValue(value.Index(i), visited))
		}
		return clone
	default:
		return value
	}
}

func (u *UnionSchema) getType() zconst.ZogType {
	return zconst.TypeUnion
}
func (u *UnionSchema) setCoercer(c CoercerFunc) {}

func (u *UnionSchema) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	children := make([]*zss.ZSSSchema, 0, len(u.schemas))
	for _, schema := range u.schemas {
		children = append(children, schema.toZSS(ctx))
	}

	return &zss.ZSSSchema{
		Kind:     zconst.TypeUnion,
		Children: children,
	}
}

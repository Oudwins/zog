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

func Union(schemas []ZogSchema, options ...SchemaOption) *UnionSchema {
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
	for _, s := range u.schemas {
		numIssues := len(ctx.Errors.List)
		branchCtx, commit := newUnionBranchCtx(ctx, s, ctx.Data)
		s.process(branchCtx)
		branchCtx.Free()
		if len(ctx.Errors.List) == numIssues {
			commit()
			if listStart == 0 {
				ctx.Errors.List = nil
			} else {
				ctx.Errors.List = ctx.Errors.List[:listStart]
			}
			return // success
		}
	}
}

func (u *UnionSchema) validate(ctx *p.SchemaCtx) {
	// Wrap the context and only go to the next one on fail. Keeping all the errors and appending at the end
	listStart := len(ctx.Errors.List)
	for _, s := range u.schemas {
		numIssues := len(ctx.Errors.List)
		branchCtx, commit := newUnionBranchCtx(ctx, s, nil)
		branchCtx.Data = branchCtx.ValPtr
		s.validate(branchCtx)
		branchCtx.Free()
		if len(ctx.Errors.List) == numIssues {
			commit()
			if listStart == 0 {
				ctx.Errors.List = nil
			} else {
				ctx.Errors.List = ctx.Errors.List[:listStart]
			}
			return // success
		}
	}

}

func newUnionBranchCtx(ctx *p.SchemaCtx, schema ZogSchema, data any) (*p.SchemaCtx, func()) {
	dest := reflect.ValueOf(ctx.ValPtr)
	if !dest.IsValid() || dest.Kind() != reflect.Pointer || dest.IsNil() {
		return ctx.NewSchemaCtx(data, ctx.ValPtr, ctx.Path, schema.getType()), func() {}
	}

	branchDest := cloneUnionValue(dest, make(map[unionCloneVisit]reflect.Value))
	commit := func() {
		if !reflect.DeepEqual(dest.Interface(), branchDest.Interface()) {
			dest.Elem().Set(branchDest.Elem())
		}
	}
	return ctx.NewSchemaCtx(data, branchDest.Interface(), ctx.Path, schema.getType()), commit
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

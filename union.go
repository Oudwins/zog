package zog

import (
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
	u.process(sctx)
	return errs.List
}

func (u *UnionSchema) process(ctx *p.SchemaCtx) {
	// Wrap the context and only go to the next one on fail. Keeping all the errors and appending at the end
	for _, s := range u.schemas {
		numIssues := len(ctx.Errors.List)
		s.process(ctx)
		if len(ctx.Errors.List) == numIssues {
			ctx.Errors.List = ctx.Errors.List[:0]
			return // success
		}
	}
}

func (u *UnionSchema) validate(ctx *p.SchemaCtx) {
	// Wrap the context and only go to the next one on fail. Keeping all the errors and appending at the end
	for _, s := range u.schemas {
		numIssues := len(ctx.Errors.List)
		s.validate(ctx)
		if len(ctx.Errors.List) == numIssues {
			ctx.Errors.List = ctx.Errors.List[:0]
			return // success
		}
	}

}
func (u *UnionSchema) getType() zconst.ZogType {
	return zconst.TypeUnion
}
func (u *UnionSchema) setCoercer(c CoercerFunc) {}
func (u *UnionSchema) toZSS(*ZSSSerializeCtx) *zss.ZSSSchema {
	return &zss.ZSSSchema{}
}

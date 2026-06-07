package tutils

import (
	"github.com/Oudwins/zog/conf"
	"github.com/Oudwins/zog/pkgs/internals"
	"github.com/Oudwins/zog/zconst"
)

func FakeContextFromValue(val any, schemaType zconst.ZogType) (*internals.SchemaCtx, *internals.ErrsList) {
	errs := internals.NewErrsList()
	ctx := internals.NewExecCtx(errs, conf.IssueFormatter)

	path := internals.NewPathBuilder()
	sctx := ctx.NewSchemaCtx(val, val, path, schemaType)

	return sctx, errs
}

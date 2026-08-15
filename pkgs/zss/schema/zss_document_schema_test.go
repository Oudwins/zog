package zssschema

import (
	"testing"

	zss "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestNilFields(t *testing.T) {
	x := zss.ZSSSchema{Kind: zconst.TypeStruct, Fields: map[string]*zss.ZSSSchema{
		"Field1": nil,
	}, FieldMeta: map[string]zss.ZSSFieldMeta{}}

	errs := ZSSSchemaSchema.Validate(&x)
	assert.NotEmpty(t, errs)
}

package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUnionFirstSchemaSucceeds(t *testing.T) {
	validator := Union([]ZogSchema{
		Int().GT(10, Message("must be greater than 10")),
		Int().LT(0, Message("must be less than 0")),
	})
	dest := 15

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 15, dest)
}

func TestValidateUnionLaterSchemaSucceeds(t *testing.T) {
	validator := Union([]ZogSchema{
		Int().GT(10, Message("must be greater than 10")),
		Int().LT(0, Message("must be less than 0")),
	})
	dest := -5

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, -5, dest)
}

func TestValidateUnionAllSchemasFail(t *testing.T) {
	validator := Union([]ZogSchema{
		Int().GT(10, Message("must be greater than 10")),
		Int().LT(0, Message("must be less than 0")),
	})
	dest := 5

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, "must be greater than 10", errs[0].Message)
	assert.Equal(t, "must be less than 0", errs[1].Message)
	assert.Equal(t, 5, dest)
}

func TestValidateUnionStringOrInt(t *testing.T) {
	validator := Union([]ZogSchema{
		String().Required(),
		Int().Required(),
	})
	dest := 15

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 15, dest)
}

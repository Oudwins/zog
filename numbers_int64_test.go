package zog

import (
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestInt64Parse(t *testing.T) {
	dest := int64(0)
	validator := Int64()
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(5), dest)
}

func TestInt64ParseFormatter(t *testing.T) {
	dest := int64(0)
	fmt := WithIssueFormatter(func(e *ZogIssue, ctx Ctx) {
		e.SetMessage("test2")
	})
	validator := Int64().GTE(10, Message("test1"))
	errs := validator.Parse(int64(5), &dest, fmt)
	assert.Equal(t, "test1", errs[0].Message)
	validator2 := Int64().GTE(10)
	errs2 := validator2.Parse(int64(5), &dest, fmt)
	assert.Equal(t, "test2", errs2[0].Message)
}

func TestInt64SchemaOption(t *testing.T) {
	s := Int64(WithCoercer(func(original any) (value any, err error) {
		return int64(42), nil
	}))

	var result int64
	err := s.Parse("123", &result)
	assert.Nil(t, err)
	assert.Equal(t, int64(42), result)
}

func TestInt64Required(t *testing.T) {
	validator := Int64().Required(Message("custom"))
	var dest int64
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(5), dest)
	dest = 0
	errs = validator.Parse("", &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)

	errs = validator.Parse("     ", &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)

	errs = validator.Parse(nil, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
}

func TestInt64Optional(t *testing.T) {
	validator := Int64().Optional()
	dest := int64(0)
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(5), dest)
}

func TestInt64Default(t *testing.T) {
	validator := Int64().Default(10)
	dest := int64(0)
	errs := validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(10), dest)
}

func TestInt64Catch(t *testing.T) {
	validator := Int64().Catch(0)
	dest := int64(0)
	errs := validator.Parse("not a number", &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(0), dest)
}

func TestInt64PreTransform(t *testing.T) {
	preTransform := func(val any, ctx Ctx) (any, error) {
		if v, ok := val.(int64); ok {
			return v * 2, nil
		}
		return val, nil
	}

	validator := Int64().PreTransform(preTransform)
	var dest int64
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(10), dest)
}

func TestInt64PostTransform(t *testing.T) {
	postTransform := func(val any, ctx Ctx) error {
		if v, ok := val.(*int64); ok {
			*v += 1
		}
		return nil
	}

	validator := Int64().PostTransform(postTransform)
	var dest int64
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(6), dest)
}

func TestInt64MultipleTransforms(t *testing.T) {
	preTransform := func(val any, ctx Ctx) (any, error) {
		if v, ok := val.(int64); ok {
			return v * 2, nil
		}
		return val, nil
	}

	postTransform := func(val any, ctx Ctx) error {
		if v, ok := val.(*int64); ok {
			*v += 1
		}
		return nil
	}

	validator := Int64().PreTransform(preTransform).PostTransform(postTransform)
	var dest int64
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(11), dest)
}

func TestInt64OneOf(t *testing.T) {
	validator := Int64().OneOf([]int64{1, 2, 3}, Message("custom"))
	dest := int64(0)
	errs := validator.Parse(int64(1), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, int64(4), dest)
}

func TestInt64Eq(t *testing.T) {
	validator := Int64().EQ(5, Message("custom"))
	dest := int64(0)
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, int64(4), dest)
}

func TestInt64Gt(t *testing.T) {
	validator := Int64().GT(5, Message("custom"))
	dest := int64(0)
	errs := validator.Parse(int64(6), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(5), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	errs = validator.Parse(int64(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, int64(4), dest)
}

func TestInt64Gte(t *testing.T) {
	dest := int64(0)
	validator := Int64().GTE(5, Message("custom"))
	errs := validator.Parse(int64(6), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, int64(4), dest)
}

func TestInt64Lt(t *testing.T) {
	dest := int64(0)
	validator := Int64().LT(5, Message("custom"))
	errs := validator.Parse(int64(4), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(5), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	errs = validator.Parse(int64(6), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, int64(6), dest)
}

func TestInt64Lte(t *testing.T) {
	dest := int64(0)
	validator := Int64().LTE(5, Message("custom"))
	errs := validator.Parse(int64(4), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(int64(6), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, int64(6), dest)
}

func TestInt64CustomTest(t *testing.T) {
	validator := Int64().TestFunc(func(val any, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, int64(5), val)
		return true
	}, Message("custom"))
	dest := int64(0)
	errs := validator.Parse(int64(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, int64(5), dest)
}

func TestInt64GetType(t *testing.T) {
	i := Int64()
	assert.Equal(t, zconst.TypeNumber, i.getType())
}

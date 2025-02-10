package zog

import (
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestNumberParse(t *testing.T) {
	dest := 0
	validator := Int()
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)
}

func TestNumberParseFormatter(t *testing.T) {
	dest := 0
	fmt := WithErrFormatter(func(e ZogError, ctx Ctx) {
		e.SetMessage("test2")
	})
	validator := Int().GTE(10, Message("test1"))
	errs := validator.Parse(5, &dest, fmt)
	assert.Equal(t, "test1", errs[0].Message())
	validator2 := Int().GTE(10)
	errs2 := validator2.Parse(5, &dest, fmt)
	assert.Equal(t, "test2", errs2[0].Message())
}

func TestIntSchemaOption(t *testing.T) {
	s := Int(WithCoercer(func(original any) (value any, err error) {
		return 42, nil
	}))

	var result int
	err := s.Parse("123", &result)
	assert.Nil(t, err)
	assert.Equal(t, 42, result)
}

func TestFloatSchemaOption(t *testing.T) {
	s := Float(WithCoercer(func(original any) (value any, err error) {
		return 3.14, nil
	}))

	var result float64
	err := s.Parse("2.718", &result)
	assert.Nil(t, err)
	assert.Equal(t, 3.14, result)
}

func TestNumberRequired(t *testing.T) {
	validator := Int().Required(Message("custom"))
	var dest int
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)
	dest = 0
	errs = validator.Parse("", &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())

	errs = validator.Parse("     ", &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())

	errs = validator.Parse(nil, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
}

func TestNumberOptional(t *testing.T) {
	validator := Int().Optional()
	dest := 0
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)
}

func TestNumberDefault(t *testing.T) {
	validator := Int().Default(10)
	dest := 0
	errs := validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 10, dest)
}

func TestNumberCatch(t *testing.T) {
	validator := Int().Catch(0)
	dest := 0
	errs := validator.Parse("not a number", &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 0, dest)
}

func TestNumberPreTransform(t *testing.T) {
	preTransform := func(val any, ctx ParseCtx) (any, error) {
		if v, ok := val.(int); ok {
			return v * 2, nil
		}
		return val, nil
	}

	validator := Int().PreTransform(preTransform)
	var dest int
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 10, dest)
}

func TestNumberPostTransform(t *testing.T) {
	postTransform := func(val any, ctx ParseCtx) error {
		if v, ok := val.(*int); ok {
			*v += 1
		}
		return nil
	}

	validator := Int().PostTransform(postTransform)
	var dest int
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 6, dest)
}

func TestNumberMultipleTransforms(t *testing.T) {
	preTransform := func(val any, ctx ParseCtx) (any, error) {
		if v, ok := val.(int); ok {
			return v * 2, nil
		}
		return val, nil
	}

	postTransform := func(val any, ctx ParseCtx) error {
		if v, ok := val.(*int); ok {
			*v += 1
		}
		return nil
	}

	validator := Int().PreTransform(preTransform).PostTransform(postTransform)
	var dest int
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 11, dest)
}

// validators
func TestNumberOneOf(t *testing.T) {
	validator := Int().OneOf([]int{1, 2, 3}, Message("custom"))
	dest := 0
	errs := validator.Parse(1, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	assert.Equal(t, 4, dest)
}

func TestNumberEq(t *testing.T) {
	validator := Int().EQ(5, Message("custom"))
	dest := 0
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	assert.Equal(t, 4, dest)
}

func TestNumberGt(t *testing.T) {
	validator := Int().GT(5, Message("custom"))
	dest := 0
	errs := validator.Parse(6, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	assert.Equal(t, 4, dest)
}

func TestNumberGte(t *testing.T) {
	dest := 0
	validator := Int().GTE(5, Message("custom"))
	errs := validator.Parse(6, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	assert.Equal(t, 4, dest)
}

func TestNumberLt(t *testing.T) {
	dest := 0
	validator := Int().LT(5, Message("custom"))
	errs := validator.Parse(4, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	errs = validator.Parse(6, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	assert.Equal(t, 6, dest)
}

func TestNumberLte(t *testing.T) {
	dest := 0
	validator := Int().LTE(5, Message("custom"))
	errs := validator.Parse(4, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(6, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message())
	assert.Equal(t, 6, dest)
}

func TestNumberCustomTest(t *testing.T) {
	validator := Int().TestFunc(func(val any, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, 5, val)
		return true
	}, Message("custom"))
	dest := 0
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)
}

func TestIntGetType(t *testing.T) {
	i := Int()
	assert.Equal(t, zconst.TypeNumber, i.getType())
}

func TestFloatGetType(t *testing.T) {
	f := Float()
	assert.Equal(t, zconst.TypeNumber, f.getType())
}

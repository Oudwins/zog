package zog

import (
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestNumberValidate(t *testing.T) {
	dest := 5
	validator := Int()
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)
}

func TestNumberValidateFormatter(t *testing.T) {
	dest := 1
	fmt := WithIssueFormatter(func(e *ZogIssue, ctx Ctx) {
		e.SetMessage("test2")
	})
	validator := Int().GTE(10, Message("test1")).Required()
	errs := validator.Validate(&dest, fmt)
	assert.Equal(t, "test1", errs[0].Message)
	validator2 := Int().GTE(10)
	errs2 := validator2.Validate(&dest, fmt)
	assert.Equal(t, "test2", errs2[0].Message)
}

func TestValidateNumberRequired(t *testing.T) {
	validator := Int().Required(Message("custom"))
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)

	dest = 0
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
}

func TestValidateNumberOptional(t *testing.T) {
	validator := Int().Optional()
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 0
	errs = validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 0, dest)
}

func TestValidateNumberDefault(t *testing.T) {
	validator := Int().Default(10)
	dest := 0
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 10, dest)
}

func TestValidateNumberCatch(t *testing.T) {
	validator := Int().Catch(0)
	dest := 42
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 42, dest)
}

func TestValidateNumberPreTransform(t *testing.T) {
	preTransform := func(val any, ctx ParseCtx) (any, error) {
		if v, ok := val.(int); ok {
			out := v * 2
			return out, nil
		}
		return val, nil
	}

	validator := Int().PreTransform(preTransform)
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 10, dest)
}

func TestValidateNumberPostTransform(t *testing.T) {
	postTransform := func(val any, ctx ParseCtx) error {
		if v, ok := val.(*int); ok {
			*v += 1
		}
		return nil
	}

	validator := Int().PostTransform(postTransform)
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 6, dest)
}

func TestValidateNumberMultipleTransforms(t *testing.T) {
	preTransform := func(val any, ctx ParseCtx) (any, error) {
		if v, ok := val.(int); ok {
			out := v * 2
			return out, nil
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
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 11, dest)
}

func TestValidateNumberOneOf(t *testing.T) {
	validator := Int().OneOf([]int{1, 2, 3}, Message("custom"))
	dest := 1
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 4
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, 4, dest)
}

func TestValidateNumberEq(t *testing.T) {
	validator := Int().EQ(5, Message("custom"))
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 4
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, 4, dest)
}

func TestValidateNumberGt(t *testing.T) {
	validator := Int().GT(5, Message("custom"))
	dest := 6
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 5
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	dest = 4
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, 4, dest)
}

func TestValidateNumberGte(t *testing.T) {
	dest := 6
	validator := Int().GTE(5, Message("custom"))
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 5
	errs = validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 4
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, 4, dest)
}

func TestValidateNumberLt(t *testing.T) {
	dest := 4
	validator := Int().LT(5, Message("custom"))
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 5
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	dest = 6
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, 6, dest)
}

func TestValidateNumberLte(t *testing.T) {
	dest := 4
	validator := Int().LTE(5, Message("custom"))
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 5
	errs = validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = 6
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, 6, dest)
}

func TestValidateNumberValidate(t *testing.T) {
	dest := 5
	validator := Int()
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)
}

func TestValidateNumberCustomTest(t *testing.T) {
	validator := Int().TestFunc(func(val any, ctx ParseCtx) bool {
		// Custom test logic here
		assert.Equal(t, 5, val)
		return true
	}, Message("custom"))
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 5, dest)
}

func TestValidateIntGetType(t *testing.T) {
	i := Int()
	assert.Equal(t, zconst.TypeNumber, i.getType())
}

func TestValidateFloatGetType(t *testing.T) {
	f := Float()
	assert.Equal(t, zconst.TypeNumber, f.getType())
}

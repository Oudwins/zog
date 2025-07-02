package zog

import (
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestUintParse(t *testing.T) {
	dest := uint(0)
	validator := Uint()
	errs := validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(5), dest)
}

func TestUintParseFormatter(t *testing.T) {
	dest := uint(0)
	fmt := WithIssueFormatter(func(e *ZogIssue, ctx Ctx) {
		e.SetMessage("test2")
	})
	validator := Uint().GTE(10, Message("test1"))
	errs := validator.Parse(uint(5), &dest, fmt)
	assert.Equal(t, "test1", errs[0].Message)
	validator2 := Uint().GTE(10)
	errs2 := validator2.Parse(uint(5), &dest, fmt)
	assert.Equal(t, "test2", errs2[0].Message)
}

func TestUintSchemaOption(t *testing.T) {
	s := Uint(WithCoercer(func(original any) (value any, err error) {
		return uint(42), nil
	}))

	var result uint
	err := s.Parse("123", &result)
	assert.Nil(t, err)
	assert.Equal(t, uint(42), result)
}

func TestUintRequired(t *testing.T) {
	validator := Uint().Required(Message("custom"))
	var dest uint
	errs := validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(5), dest)
	dest = 0
	errs = validator.Parse("", &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, zconst.IssueCodeCoerce, errs[0].Code)

	errs = validator.Parse("     ", &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, zconst.IssueCodeCoerce, errs[0].Code)

	errs = validator.Parse(nil, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
}

func TestUintOptional(t *testing.T) {
	validator := Uint().Optional()
	dest := uint(0)
	errs := validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(5), dest)
}

func TestUintDefault(t *testing.T) {
	validator := Uint().Default(10)
	dest := uint(0)
	errs := validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(10), dest)
}

func TestUintCatch(t *testing.T) {
	validator := Uint().Catch(0)
	dest := uint(0)
	errs := validator.Parse("not a number", &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(0), dest)
}

func TestUintPostTransform(t *testing.T) {
	postTransform := func(val *uint, ctx Ctx) error {
		*val += 1
		return nil
	}

	validator := Uint().Transform(postTransform)
	var dest uint
	errs := validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(6), dest)
}

func TestUintMultipleTransforms(t *testing.T) {

	postTransform := func(val *uint, ctx Ctx) error {
		*val += 1
		return nil
	}

	validator := Uint().Transform(postTransform).Transform(postTransform)
	var dest uint
	errs := validator.Parse(uint(9), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(11), dest)
}

func TestUintOneOf(t *testing.T) {
	validator := Uint().OneOf([]uint{1, 2, 3}, Message("custom"))
	dest := uint(0)
	errs := validator.Parse(uint(1), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, uint(4), dest)
}

func TestUintEq(t *testing.T) {
	validator := Uint().EQ(5, Message("custom"))
	dest := uint(0)
	errs := validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, uint(4), dest)
}

func TestUintGt(t *testing.T) {
	validator := Uint().GT(5, Message("custom"))
	dest := uint(0)
	errs := validator.Parse(uint(6), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(5), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	errs = validator.Parse(uint(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, uint(4), dest)
}

func TestUintGte(t *testing.T) {
	dest := uint(0)
	validator := Uint().GTE(5, Message("custom"))
	errs := validator.Parse(uint(6), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(4), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, uint(4), dest)
}

func TestUintLt(t *testing.T) {
	dest := uint(0)
	validator := Uint().LT(5, Message("custom"))
	errs := validator.Parse(uint(4), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(5), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	errs = validator.Parse(uint(6), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, uint(6), dest)
}

func TestUintLte(t *testing.T) {
	dest := uint(0)
	validator := Uint().LTE(5, Message("custom"))
	errs := validator.Parse(uint(4), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(uint(6), &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, uint(6), dest)
}

func TestUintCustomTest(t *testing.T) {
	validator := Uint().TestFunc(func(val *uint, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, uint(5), *val)
		return true
	}, Message("custom"))
	dest := uint(0)
	errs := validator.Parse(uint(5), &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, uint(5), dest)
}

func TestUintGetType(t *testing.T) {
	i := Uint()
	assert.Equal(t, zconst.TypeNumber, i.getType())
}

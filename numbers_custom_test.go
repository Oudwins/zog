package zog

import (
	"testing"

	"github.com/Oudwins/zog/conf"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

type MyInt int64

func CustomInt(opts ...SchemaOption) *NumberSchema[MyInt] {
	s := &NumberSchema[MyInt]{}
	opts = append(
		[]SchemaOption{
			WithCoercer(func(x any) (any, error) {
				v, e := conf.DefaultCoercers.Int(x)
				if e != nil {
					return nil, e
				}
				return MyInt(v.(int)), e
			}),
		},
		opts...,
	)
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func TestCustomNumberParse(t *testing.T) {
	dest := MyInt(0)
	validator := CustomInt()
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(5), dest)
}

func TestCustomNumberParseFormatter(t *testing.T) {
	dest := MyInt(0)
	fmt := WithIssueFormatter(func(e *ZogIssue, ctx Ctx) {
		e.SetMessage("test2")
	})
	validator := CustomInt().GTE(10, Message("test1"))
	errs := validator.Parse(5, &dest, fmt)
	assert.Equal(t, "test1", errs[0].Message)
	validator2 := CustomInt().GTE(10)
	errs2 := validator2.Parse(5, &dest, fmt)
	assert.Equal(t, "test2", errs2[0].Message)
}

func TestCustomNumberSchemaOption(t *testing.T) {
	s := CustomInt(WithCoercer(func(original any) (value any, err error) {
		return MyInt(42), nil
	}))

	var result MyInt
	err := s.Parse("123", &result)
	assert.Nil(t, err)
	assert.Equal(t, MyInt(42), result)
}

func TestCustomNumberRequired(t *testing.T) {
	validator := CustomInt().Required(Message("custom"))
	var dest MyInt
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(5), dest)
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

func TestCustomNumberOptional(t *testing.T) {
	validator := CustomInt().Optional()
	dest := MyInt(0)
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(5), dest)
}

func TestCustomNumberDefault(t *testing.T) {
	validator := CustomInt().Default(10)
	dest := MyInt(0)
	errs := validator.Parse(nil, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(10), dest)
}

func TestCustomNumberCatch(t *testing.T) {
	validator := CustomInt().Catch(0)
	dest := MyInt(0)
	errs := validator.Parse("not a number", &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(0), dest)
}

func TestCustomNumberPostTransform(t *testing.T) {
	postTransform := func(val any, ctx Ctx) error {
		if v, ok := val.(*MyInt); ok {
			*v += 1
		}
		return nil
	}

	validator := CustomInt().PostTransform(postTransform)
	var dest MyInt
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(6), dest)
}

func TestCustomNumberMultipleTransforms(t *testing.T) {

	postTransform := func(val any, ctx Ctx) error {
		if v, ok := val.(*MyInt); ok {
			*v += 1
		}
		return nil
	}

	validator := CustomInt().PostTransform(postTransform).PostTransform(postTransform)
	var dest MyInt
	errs := validator.Parse(9, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(11), dest)
}

func TestCustomNumberOneOf(t *testing.T) {
	validator := CustomInt().OneOf([]MyInt{1, 2, 3}, Message("custom"))
	dest := MyInt(0)
	errs := validator.Parse(1, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, MyInt(4), dest)
}

func TestCustomNumberEq(t *testing.T) {
	validator := CustomInt().EQ(5, Message("custom"))
	dest := MyInt(0)
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, MyInt(4), dest)
}

func TestCustomNumberGt(t *testing.T) {
	validator := CustomInt().GT(5, Message("custom"))
	dest := MyInt(0)
	errs := validator.Parse(6, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	errs = validator.Parse(4, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, MyInt(4), dest)
}

func TestCustomNumberGte(t *testing.T) {
	dest := MyInt(0)
	validator := CustomInt().GTE(5, Message("custom"))
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
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, MyInt(4), dest)
}

func TestCustomNumberLt(t *testing.T) {
	dest := MyInt(0)
	validator := CustomInt().LT(5, Message("custom"))
	errs := validator.Parse(4, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	errs = validator.Parse(5, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	errs = validator.Parse(6, &dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, MyInt(6), dest)
}

func TestCustomNumberLte(t *testing.T) {
	dest := MyInt(0)
	validator := CustomInt().LTE(5, Message("custom"))
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
	assert.Equal(t, "custom", errs[0].Message)
	assert.Equal(t, MyInt(6), dest)
}

func TestCustomNumberCustomTest(t *testing.T) {
	validator := CustomInt().TestFunc(func(val any, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, MyInt(5), val)
		return true
	}, Message("custom"))
	dest := MyInt(0)
	errs := validator.Parse(5, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, MyInt(5), dest)
}

func TestCustomNumberGetType(t *testing.T) {
	i := CustomInt()
	assert.Equal(t, zconst.TypeNumber, i.getType())
}

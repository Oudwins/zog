package zog

import (
	"testing"

	"github.com/Oudwins/zog/tutils"
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

func TestValidateNumberPostTransform(t *testing.T) {
	postTransform := func(val *int, ctx Ctx) error {
		*val += 1
		return nil
	}

	validator := Int().Transform(postTransform)
	dest := 5
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, 6, dest)
}

func TestValidateNumberMultipleTransforms(t *testing.T) {

	postTransform := func(val *int, ctx Ctx) error {
		*val += 1
		return nil
	}

	validator := Int().Transform(postTransform).Transform(postTransform)
	dest := 9
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
	validator := Int().TestFunc(func(val *int, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, 5, *val)
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

func TestValidateIntNot(t *testing.T) {
	tests := map[string]struct {
		schema    *NumberSchema[int]
		value     int
		expectErr bool
	}{
		"Not eq true": {
			schema:    Int().Not().EQ(1),
			value:     2,
			expectErr: false,
		},
		"Not eq false": {
			schema:    Int().Not().EQ(1),
			value:     1,
			expectErr: true,
		},
		"Not one of true": {
			schema:    Int().Not().OneOf([]int{1, 2, 3, 4}),
			value:     5,
			expectErr: false,
		},
		"Not one of false": {
			schema:    Int().Not().OneOf([]int{1, 2, 3, 4}),
			value:     2,
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			errs := tc.schema.Validate(&tc.value)
			if tc.expectErr {
				assert.NotEmpty(t, errs)
				tutils.VerifyDefaultIssueMessages(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestValidateFloatNot(t *testing.T) {
	tests := map[string]struct {
		schema    *NumberSchema[float64]
		value     float64
		expectErr bool
	}{
		"Not eq true": {
			schema:    Float64().Not().EQ(1),
			value:     2,
			expectErr: false,
		},
		"Not eq false": {
			schema:    Float64().Not().EQ(1),
			value:     1,
			expectErr: true,
		},
		"Not one of true": {
			schema:    Float64().Not().OneOf([]float64{1, 2, 3, 4}),
			value:     5,
			expectErr: false,
		},
		"Not one of false": {
			schema:    Float64().Not().OneOf([]float64{1, 2, 3, 4}),
			value:     2,
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			errs := tc.schema.Validate(&tc.value)
			if tc.expectErr {
				assert.NotEmpty(t, errs)
				tutils.VerifyDefaultIssueMessages(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

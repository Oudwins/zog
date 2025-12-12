package zog

import (
	"strings"
	"testing"

	"github.com/Oudwins/zog/tutils"
	"github.com/stretchr/testify/assert"
)

func TestValidateSliceRequired(t *testing.T) {
	validator := Slice(String())
	dest := []string{"test"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, []string{"test"}, dest)

	validator = validator.Required()
	dest = []string{}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateSliceOptional(t *testing.T) {
	validator := Slice(String()).Optional()
	dest := []string{"test"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = []string{}
	errs = validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
}

func TestValidateSliceDefault(t *testing.T) {
	validator := Slice(String()).Default([]string{"default"})
	dest := []string{}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, []string{"default"}, dest)
}

func TestValidateSliceTransform(t *testing.T) {
	transform := func(val any, ctx Ctx) error {
		if v, ok := val.(*[]string); ok {
			for i := range *v {
				(*v)[i] = strings.ToUpper((*v)[i])
			}
		}
		return nil
	}

	validator := Slice(String()).Transform(transform)
	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, []string{"TEST", "EXAMPLE"}, dest)
}

func TestValidateSliceLen(t *testing.T) {
	validator := Slice(String()).Len(2)
	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = []string{"test"}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateSliceMin(t *testing.T) {
	validator := Slice(String()).Min(2)
	dest := []string{"test", "example", "extra"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = []string{"test"}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateSliceMax(t *testing.T) {
	validator := Slice(String()).Max(2)
	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = []string{"test", "example", "extra"}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateSliceContains(t *testing.T) {
	validator := Slice(String()).Contains("test")
	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = []string{"example", "sample"}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateSliceCustomTest(t *testing.T) {
	validator := Slice(String()).TestFunc(func(val any, ctx Ctx) bool {
		if v, ok := val.(*[]string); ok {
			return len(*v) > 0 && (*v)[0] == "test"
		}
		return false
	}, Message("custom"))

	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = []string{"wrong", "example"}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
	// assert.Equal(t, "custom_test", rootErrs[0].Code())
}

// TODO not yet supported
// func TestValidateSliceNestedValidation(t *testing.T) {
// 	type User struct {
// 		Name string
// 	}

// 	userSchema := Struct(Shape{
// 		"name": String().Required(),
// 	})

// 	validator := Slice(userSchema)
// 	dest := []User{
// 		{Name: "John"},
// 		{Name: "Jane"},
// 	}

// 	errs := validator.Validate(&dest)
// 	if len(errs) > 0 {
// 		t.Errorf("Expected no errors, got %v", errs)
// 	}

// 	dest = []User{
// 		{Name: ""},
// 		{Name: ""},
// 	}

// 	errs = validator.Validate(&dest)
// 	assert.NotEmpty(t, errs)
// 	assert.NotEmpty(t, errs["[0].name"])
// 	assert.NotEmpty(t, errs["[1].name"])
// }

func TestValidateSliceMultipleValidators(t *testing.T) {
	validator := Slice(String()).
		Min(2, Message("too short")).
		Max(4, Message("too long")).
		Contains("test", Message("must contain test"))

	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = []string{"wrong"}
	errs = validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	m := Issues.Flatten(errs)
	rootErrs := m["zconst.ISSUE_KEY_ROOT"]
	assert.Contains(t, rootErrs, "too short")
	assert.Contains(t, rootErrs, "must contain test")
	assert.Len(t, rootErrs, 2)

	dest = []string{"a", "b", "c", "d", "e"}
	errs = validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	rootErrs = Issues.Flatten(errs)["zconst.ISSUE_KEY_ROOT"]
	assert.Contains(t, rootErrs, "too long")
	assert.Contains(t, rootErrs, "must contain test")
	assert.Len(t, rootErrs, 2)
}

func TestValidateSliceNot(t *testing.T) {
	tests := map[string]struct {
		schema    *SliceSchema
		value     []int
		expectErr bool
	}{
		"not len true": {
			schema:    Slice(Int()).Not().Len(2),
			value:     []int{1},
			expectErr: false,
		},
		"not len false": {
			schema:    Slice(Int()).Not().Len(2),
			value:     []int{1, 2},
			expectErr: true,
		},
		"not contains true": {
			schema:    Slice(Int()).Not().Contains([]int{1, 3}),
			value:     []int{1, 2},
			expectErr: false,
		},
		"not contains false": {
			schema:    Slice(Int()).Not().Contains(1),
			value:     []int{1, 2},
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

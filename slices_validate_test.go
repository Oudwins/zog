package zog

import (
	"strings"
	"testing"

	"github.com/Oudwins/zog/tutils"
	"github.com/stretchr/testify/assert"
)

func TestValidateSliceRequired(t *testing.T) {
	validator := Slice[string](String())
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
}

func TestValidateSliceOptional(t *testing.T) {
	validator := Slice[string](String()).Optional()
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
	validator := Slice[string](String()).Default([]string{"default"})
	dest := []string{}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, []string{"default"}, dest)
}

func TestValidateSliceTransform(t *testing.T) {
	transform := func(val []string, ctx Ctx) error {
		for i := range val {
			val[i] = strings.ToUpper(val[i])
		}
		return nil
	}

	validator := Slice[string](String()).Transform(transform)
	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, []string{"TEST", "EXAMPLE"}, dest)
}

func TestValidateSliceLen(t *testing.T) {
	validator := Slice[string](String()).Len(2)
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
}

func TestValidateSliceMin(t *testing.T) {
	validator := Slice[string](String()).Min(2)
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
}

func TestValidateSliceMax(t *testing.T) {
	validator := Slice[string](String()).Max(2)
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
}

func TestValidateSliceContains(t *testing.T) {
	validator := Slice[string](String()).Contains("test")
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
}

func TestValidateSliceCustomTest(t *testing.T) {
	validator := Slice[string](String()).TestFunc(func(val []string, ctx Ctx) bool {
		return len(val) > 0 && val[0] == "test"
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
	assert.Equal(t, "custom", errs["$root"][0].Message)
	// assert.Equal(t, "custom_test", errs["$root"][0].Code())
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
	validator := Slice[string](String()).
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
	assert.Contains(t, Issues.SanitizeList(errs["$root"]), "too short")
	assert.Contains(t, Issues.SanitizeList(errs["$root"]), "must contain test")
	assert.Len(t, errs["$root"], 2)

	dest = []string{"a", "b", "c", "d", "e"}
	errs = validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Contains(t, Issues.SanitizeList(errs["$root"]), "too long")
	assert.Contains(t, Issues.SanitizeList(errs["$root"]), "must contain test")
	assert.Len(t, errs["$root"], 2)
}

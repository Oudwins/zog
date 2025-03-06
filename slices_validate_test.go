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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
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

func TestValidateSlicePreTransform(t *testing.T) {
	preTransform := func(val any, ctx ParseCtx) (any, error) {
		if v, ok := val.([]string); ok {
			out := make([]string, len(v))
			copy(out, v)
			for i := range out {
				out[i] = strings.ToUpper(out[i])
			}
			return out, nil
		}
		return val, nil
	}

	validator := Slice(String()).PreTransform(preTransform)
	dest := []string{"test", "example"}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, []string{"TEST", "EXAMPLE"}, dest)
}

func TestValidateSlicePostTransform(t *testing.T) {
	postTransform := func(val any, ctx ParseCtx) error {
		if v, ok := val.(*[]string); ok {
			for i := range *v {
				(*v)[i] = strings.ToUpper((*v)[i])
			}
		}
		return nil
	}

	validator := Slice(String()).PostTransform(postTransform)
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
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
	tutils.VerifyDefaultIssueMessagesMap(t, errs)
}

func TestValidateSliceCustomTest(t *testing.T) {
	validator := Slice(String()).TestFunc(func(val any, ctx ParseCtx) bool {
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
	assert.Equal(t, "custom", errs["$root"][0].Message)
	// assert.Equal(t, "custom_test", errs["$root"][0].Code())
}

// TODO not yet supported
// func TestValidateSliceNestedValidation(t *testing.T) {
// 	type User struct {
// 		Name string
// 	}

// 	userSchema := Struct(Schema{
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

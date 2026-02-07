package zog

import (
	"testing"
	"time"

	"github.com/Oudwins/zog/pkgs/internals/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestValidateMapRequired(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int())
	dest := map[string]int{"test": 1}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, map[string]int{"test": 1}, dest)

	validator = validator.Required()
	dest = map[string]int{}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateMapOptional(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).Optional()
	dest := map[string]int{"test": 1}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	dest = map[string]int{}
	errs = validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
}

func TestValidateMapDefault(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).Default(map[string]int{"default": 42})
	dest := map[string]int{}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, map[string]int{"default": 42}, dest)
}

func TestValidateMapDefaultFunc(t *testing.T) {
	defaultCalls := 0
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).DefaultFunc(func() map[string]int {
		defaultCalls++
		return map[string]int{"default": 42}
	})
	dest := map[string]int{}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, map[string]int{"default": 42}, dest)
	assert.Equal(t, 1, defaultCalls)
}

func TestValidateMapTransform(t *testing.T) {
	transform := func(val any, ctx Ctx) error {
		if v, ok := val.(*map[string]int); ok {
			for k := range *v {
				(*v)[k] = (*v)[k] * 2
			}
		}
		return nil
	}

	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).Transform(transform)
	dest := map[string]int{"test": 5, "example": 10}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, map[string]int{"test": 10, "example": 20}, dest)
}

func TestValidateMapLen(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).Len(2)
	dest := map[string]int{"test": 1, "example": 2}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = map[string]int{"test": 1}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateMapMin(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).Min(2)
	dest := map[string]int{"test": 1, "example": 2, "extra": 3}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = map[string]int{"test": 1}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateMapMax(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).Max(2)
	dest := map[string]int{"test": 1, "example": 2}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = map[string]int{"test": 1, "example": 2, "extra": 3}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateMapCustomTest(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).TestFunc(func(val any, ctx Ctx) bool {
		if v, ok := val.(*map[string]int); ok {
			return len(*v) > 0 && (*v)["test"] == 1
		}
		return false
	}, Message("custom"))

	dest := map[string]int{"test": 1, "example": 2}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = map[string]int{"wrong": 1, "example": 2}
	errs = validator.Validate(&dest)
	if len(errs) == 0 {
		t.Errorf("Expected errors, got none")
	}
	assert.Equal(t, "custom", errs[0].Message)
}

func TestValidateMapMultipleValidators(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int()).
		Min(2, Message("too short")).
		Max(4, Message("too long"))

	dest := map[string]int{"test": 1, "example": 2}
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	dest = map[string]int{"wrong": 1}
	errs = validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	m := Issues.Flatten(errs)
	rootErrs := m[zconst.ISSUE_KEY_ROOT]
	assert.Contains(t, rootErrs, "too short")
	assert.Len(t, rootErrs, 1)

	dest = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	errs = validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	rootErrs = Issues.Flatten(errs)[zconst.ISSUE_KEY_ROOT]
	assert.Contains(t, rootErrs, "too long")
	assert.Len(t, rootErrs, 1)
}

func TestValidateMapValueValidation(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String(), Int().GTE(5))
	dest := map[string]int{
		"valid":   10,
		"invalid": 3,
	}
	errs := validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	errsM := Issues.Flatten(errs)
	assert.Len(t, errsM[`["invalid"]`], 1)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

// Test different key types in Validate
func TestValidateMapFloat64Keys(t *testing.T) {
	validator := EXPERIMENTAL_MAP[float64, string](Float64(), String())
	dest := map[float64]string{
		1.5:  "one point five",
		2.7:  "two point seven",
		3.14: "pi",
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

func TestValidateMapBoolKeys(t *testing.T) {
	validator := EXPERIMENTAL_MAP[bool, string](Bool(), String())
	dest := map[bool]string{
		true:  "yes",
		false: "no",
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 2)
}

func TestValidateMapInt64Keys(t *testing.T) {
	validator := EXPERIMENTAL_MAP[int64, string](Int64(), String())
	dest := map[int64]string{
		100: "hundred",
		200: "two hundred",
		300: "three hundred",
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

func TestValidateMapUintKeys(t *testing.T) {
	validator := EXPERIMENTAL_MAP[uint, string](Uint(), String())
	dest := map[uint]string{
		10: "ten",
		20: "twenty",
		30: "thirty",
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

// Test different value types in Validate
func TestValidateMapBoolValues(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, bool](String(), Bool())
	dest := map[string]bool{
		"flag1": true,
		"flag2": false,
		"flag3": true,
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

func TestValidateMapFloatValues(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, float64](String(), Float64())
	dest := map[string]float64{
		"pi":    3.14,
		"e":     2.718,
		"sqrt2": 1.414,
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

func TestValidateMapTimeValues(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, time.Time](String(), Time())
	dest := map[string]time.Time{
		"date1": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"date2": time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 2)
}

func TestValidateMapInt64Values(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int64](String(), Int64())
	dest := map[string]int64{
		"large1": 9223372036854775807,
		"large2": 9223372036854775806,
		"large3": 100,
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

func TestValidateMapUintValues(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, uint](String(), Uint())
	dest := map[string]uint{
		"val1": 100,
		"val2": 200,
		"val3": 300,
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

// Test maps with complex value types in Validate
func TestValidateMapPtrValues(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, *int](String(), Ptr(Int()))
	val1 := 10
	val2 := 20
	dest := map[string]*int{
		"val1": &val1,
		"val2": &val2,
		"val3": nil,
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 3)
}

func TestValidateMapSliceValues(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, []string](String(), Slice(String()))
	dest := map[string][]string{
		"fruits": {"apple", "banana", "cherry"},
		"colors": {"red", "blue"},
		"empty":  {},
		"single": {"one"},
	}
	errs := validator.Validate(&dest)
	assert.Empty(t, errs)
	assert.Len(t, dest, 4)
}

func TestValidateMapSliceValuesWithValidation(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, []string](String(), Slice(String().Min(3)))
	dest := map[string][]string{
		"valid":   {"apple", "banana"},
		"invalid": {"a", "b"}, // should fail
	}
	errs := validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	errsM := Issues.Flatten(errs)
	// Errors should be at ["invalid"][0] and ["invalid"][1]
	assert.True(t, len(errsM[`["invalid"][0]`]) > 0 || len(errsM[`["invalid"][1]`]) > 0)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

// Test key validation errors
func TestValidateMapKeyValidation(t *testing.T) {
	validator := EXPERIMENTAL_MAP[string, int](String().Min(3), Int())
	dest := map[string]int{
		"valid_key":   10,
		"ab":          20, // key too short, should fail
		"another_key": 30,
	}
	errs := validator.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

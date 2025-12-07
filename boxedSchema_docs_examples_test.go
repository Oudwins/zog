package zog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function for Example 1
func createValuer(v string) StringValuer {
	return myStringValuer{v: v}
}

// Example 3: Omittable pattern types and helpers
type Omittable[T any] interface {
	Value() T
	IsSet() bool
}

type omittableString struct {
	value *string
	set   bool
}

func (o omittableString) Value() string {
	if o.value == nil {
		return ""
	}
	return *o.value
}

func (o omittableString) IsSet() bool {
	return o.set
}

func createOmittable(s *string) Omittable[string] {
	if s == nil {
		return omittableString{value: nil, set: false}
	}
	return omittableString{value: s, set: true}
}

// ============================================================================
// Tests for documentation examples
// ============================================================================

// TestExample1DriverValuerPattern tests Example 1 from the docs
func TestExample1DriverValuerPattern(t *testing.T) {
	// Example 1: driver.Valuer pattern
	type StringValuer interface {
		Value() (string, error)
	}

	schema := Boxed(
		String().Min(3),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	// Test Parse
	var valuer StringValuer
	errs := schema.Parse("hello", &valuer)
	assert.Empty(t, errs)
	val, _ := valuer.Value()
	assert.Equal(t, "hello", val)

	// Test Validate
	valuer = createValuer("hello2")
	errs = schema.Validate(&valuer)
	assert.Empty(t, errs)
	val, _ = valuer.Value()
	assert.Equal(t, "hello2", val)
}

// TestExample2NullablePattern tests Example 2 from the docs
func TestExample2NullablePattern(t *testing.T) {
	// Example 2: Nullable pattern (like sql.NullString)
	type NullString struct {
		String string
		Valid  bool
	}

	schema := Boxed(
		String().Min(3),
		func(ns NullString, ctx Ctx) (string, error) {
			if !ns.Valid {
				return "", errors.New("null string is not valid")
			}
			return ns.String, nil
		},
		func(s string, ctx Ctx) (NullString, error) {
			return NullString{String: s, Valid: true}, nil
		},
	)

	// Test Parse with valid string
	var ns NullString
	errs := schema.Parse("hello", &ns)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", ns.String)
	assert.Equal(t, true, ns.Valid)

	// Test Validate with valid NullString
	ns = NullString{String: "world", Valid: true}
	errs = schema.Validate(&ns)
	assert.Empty(t, errs)
	assert.Equal(t, "world", ns.String)
	assert.Equal(t, true, ns.Valid)
}

// TestExample3OmittablePattern tests Example 3 from the docs
func TestExample3OmittablePattern(t *testing.T) {
	// Example 3: Omittable pattern
	type Omittable[T any] interface {
		Value() T
		IsSet() bool
	}

	schema := Boxed(
		Ptr(String().Min(3)),
		func(o Omittable[string], ctx Ctx) (*string, error) {
			if o.IsSet() {
				val := o.Value()
				return &val, nil
			}
			return nil, nil
		},
		func(s *string, ctx Ctx) (Omittable[string], error) {
			return createOmittable(s), nil
		},
	)

	// Test Parse with valid string (not pointer - Ptr schema will handle pointer creation)
	var omittable Omittable[string]
	errs := schema.Parse("hello", &omittable)
	assert.Empty(t, errs)
	assert.True(t, omittable.IsSet())
	val := omittable.Value()
	assert.Equal(t, "hello", val)

	// Test Parse with nil (omitted)
	var omittable2 Omittable[string]
	errs = schema.Parse(nil, &omittable2)
	assert.Empty(t, errs)
	assert.False(t, omittable2.IsSet())

	// Test Validate with set value
	val2 := "world"
	var omittable3 Omittable[string] = createOmittable(&val2)
	errs = schema.Validate(&omittable3)
	assert.Empty(t, errs)
	assert.True(t, omittable3.IsSet())
	assert.Equal(t, "world", omittable3.Value())

	// Test Validate with IsSet = false (omitted)
	var omittable4 Omittable[string] = createOmittable(nil)
	errs = schema.Validate(&omittable4)
	assert.Empty(t, errs)
	assert.False(t, omittable4.IsSet())
}

package zog

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// 1. Primitive Type Parsing (struct box) - Raw Data Input
// ============================================================================

func TestBoxedStringParse(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("hello", &box)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V)
}

func TestBoxedStringParseFailure(t *testing.T) {
	s := Boxed(
		String().Min(5),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("hi", &box) // too short
	assert.NotEmpty(t, errs)
}

func TestBoxedBoolParse(t *testing.T) {
	s := Boxed(
		Bool(),
		func(b BoolBox, ctx Ctx) (bool, error) { return b.Value, nil },
		func(v bool, ctx Ctx) (BoolBox, error) { return BoolBox{Value: v}, nil },
	)

	var box BoolBox
	errs := s.Parse(true, &box)
	assert.Empty(t, errs)
	assert.Equal(t, true, box.Value)
}

func TestBoxedIntParse(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	errs := s.Parse(42, &box)
	assert.Empty(t, errs)
	assert.Equal(t, 42, box.Value)
}

func TestBoxedIntParseFailure(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	errs := s.Parse(-1, &box) // violates GT(0)
	assert.NotEmpty(t, errs)
}

func TestBoxedFloat64Parse(t *testing.T) {
	s := Boxed(
		Float64().GT(0),
		func(b Float64Box, ctx Ctx) (float64, error) { return b.Value, nil },
		func(v float64, ctx Ctx) (Float64Box, error) { return Float64Box{Value: v}, nil },
	)

	var box Float64Box
	errs := s.Parse(3.14, &box)
	assert.Empty(t, errs)
	assert.Equal(t, 3.14, box.Value)
}

func TestBoxedTimeParse(t *testing.T) {
	s := Boxed(
		Time(),
		func(b TimeBox, ctx Ctx) (time.Time, error) { return b.Value, nil },
		func(v time.Time, ctx Ctx) (TimeBox, error) { return TimeBox{Value: v}, nil },
	)

	timestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	var box TimeBox
	errs := s.Parse("2023-01-01T00:00:00Z", &box)
	assert.Empty(t, errs)
	assert.Equal(t, timestamp, box.Value)
}

// ============================================================================
// 2. Primitive Type Parsing (struct box) - Box Input (B)
// ============================================================================

func TestBoxedStringParseFromBox(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	inputBox := StringBox{V: "hello"}
	errs := s.Parse(inputBox, &box)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V)
}

func TestBoxedStringParseFromBoxFailure(t *testing.T) {
	s := Boxed(
		String().Min(5),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	inputBox := StringBox{V: "hi"} // too short
	errs := s.Parse(inputBox, &box)
	assert.NotEmpty(t, errs)
}

func TestBoxedIntParseFromBox(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	inputBox := IntBox{Value: 42}
	errs := s.Parse(inputBox, &box)
	assert.Empty(t, errs)
	assert.Equal(t, 42, box.Value)
}

// ============================================================================
// 3. Primitive Type Parsing (struct box) - Pointer to Box Input (*B)
// ============================================================================

func TestBoxedStringParseFromBoxPtr(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	inputBox := StringBox{V: "hello"}
	errs := s.Parse(&inputBox, &box)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V)
}

func TestBoxedIntParseFromBoxPtr(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	inputBox := IntBox{Value: 42}
	errs := s.Parse(&inputBox, &box)
	assert.Empty(t, errs)
	assert.Equal(t, 42, box.Value)
}

// ============================================================================
// 4. Primitive Type Parsing (interface box) - Raw Data Input
// ============================================================================

func TestBoxedStringValuerParse(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	var valuer StringValuer
	errs := s.Parse("hello", &valuer)
	assert.Empty(t, errs)
	val, _ := valuer.Value()
	assert.Equal(t, "hello", val)
}

func TestBoxedStringValuerParseFailure(t *testing.T) {
	s := Boxed(
		String().Min(5),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	var valuer StringValuer
	errs := s.Parse("hi", &valuer) // too short
	assert.NotEmpty(t, errs)
}

func TestBoxedIntValuerParse(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntValuer, ctx Ctx) (int, error) { return b.Value() },
		func(v int, ctx Ctx) (IntValuer, error) { return myIntValuer{v: v}, nil },
	)

	var valuer IntValuer
	errs := s.Parse(42, &valuer)
	assert.Empty(t, errs)
	val, _ := valuer.Value()
	assert.Equal(t, 42, val)
}

func TestBoxedIntValuerParseFailure(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntValuer, ctx Ctx) (int, error) { return b.Value() },
		func(v int, ctx Ctx) (IntValuer, error) { return myIntValuer{v: v}, nil },
	)

	var valuer IntValuer
	errs := s.Parse(-1, &valuer) // violates GT(0)
	assert.NotEmpty(t, errs)
}

// ============================================================================
// 5. Complex Type Parsing - Raw Data Input
// ============================================================================

func TestBoxedSliceParse(t *testing.T) {
	s := Boxed(
		Slice(String().Min(1)),
		func(b SliceBox, ctx Ctx) ([]string, error) { return b.Value, nil },
		func(v []string, ctx Ctx) (SliceBox, error) { return SliceBox{Value: v}, nil },
	)

	var box SliceBox
	errs := s.Parse([]string{"hello", "world"}, &box)
	assert.Empty(t, errs)
	assert.Equal(t, []string{"hello", "world"}, box.Value)
}

func TestBoxedSliceParseFailure(t *testing.T) {
	s := Boxed(
		Slice(String().Min(5)),
		func(b SliceBox, ctx Ctx) ([]string, error) { return b.Value, nil },
		func(v []string, ctx Ctx) (SliceBox, error) { return SliceBox{Value: v}, nil },
	)

	var box SliceBox
	errs := s.Parse([]string{"hi"}, &box) // too short
	assert.NotEmpty(t, errs)
}

func TestBoxedStructParse(t *testing.T) {
	s := Boxed(
		Struct(Shape{
			"Id":   String().Min(1),
			"Name": String().Min(1),
		}),
		func(b UserBox, ctx Ctx) (BoxedUser, error) { return b.Value, nil },
		func(u BoxedUser, ctx Ctx) (UserBox, error) { return UserBox{Value: u}, nil },
	)

	var box UserBox
	input := map[string]any{
		"Id":   "1",
		"Name": "John Doe",
	}
	errs := s.Parse(input, &box)
	assert.Empty(t, errs)
	assert.Equal(t, BoxedUser{Id: "1", Name: "John Doe"}, box.Value)
}

func TestBoxedStructParseFailure(t *testing.T) {
	s := Boxed(
		Struct(Shape{
			"Id":   String().Min(1),
			"Name": String().Min(5), // Name must be at least 5 chars
		}),
		func(b UserBox, ctx Ctx) (BoxedUser, error) { return b.Value, nil },
		func(u BoxedUser, ctx Ctx) (UserBox, error) { return UserBox{Value: u}, nil },
	)

	var box UserBox
	input := map[string]any{
		"Id":   "1",
		"Name": "Joe", // too short
	}
	errs := s.Parse(input, &box)
	assert.NotEmpty(t, errs)
}

func TestBoxedSliceValuerParse(t *testing.T) {
	s := Boxed(
		Slice(String().Min(1)),
		func(b SliceValuer, ctx Ctx) ([]string, error) { return b.Value() },
		func(v []string, ctx Ctx) (SliceValuer, error) { return mySliceValuer{v: v}, nil },
	)

	var valuer SliceValuer
	errs := s.Parse([]string{"hello", "world"}, &valuer)
	assert.Empty(t, errs)
	val, _ := valuer.Value()
	assert.Equal(t, []string{"hello", "world"}, val)
}

// ============================================================================
// 6. Error Handling
// ============================================================================

func TestBoxedUnboxErrorStructParse(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringBox, ctx Ctx) (string, error) {
			if b.V == "" {
				return "", errors.New("cannot unbox empty string")
			}
			return b.V, nil
		},
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	inputBox := StringBox{V: ""}
	errs := s.Parse(inputBox, &box)
	assert.NotEmpty(t, errs)
}

func TestBoxedBoxError(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) {
			if s == "error" {
				return StringBox{}, errors.New("box error")
			}
			return StringBox{V: s}, nil
		},
	)

	var box StringBox
	errs := s.Parse("error", &box)
	assert.NotEmpty(t, errs)
}

// ============================================================================
// 7. Real-World Patterns
// ============================================================================

func TestBoxedNullablePatternParse(t *testing.T) {
	s := Boxed(
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

	// Valid nullable string from raw data
	var ns NullString
	errs := s.Parse("hello", &ns)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", ns.String)
	assert.Equal(t, true, ns.Valid)

	// Valid nullable string from box
	var ns2 NullString
	inputBox := NullString{String: "hello", Valid: true}
	errs2 := s.Parse(inputBox, &ns2)
	assert.Empty(t, errs2)
	assert.Equal(t, "hello", ns2.String)
	assert.Equal(t, true, ns2.Valid)
}

func TestBoxedNullablePatternParseFailure(t *testing.T) {
	s := Boxed(
		String().Min(5),
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

	// Valid but too short
	var ns NullString
	errs := s.Parse("hi", &ns)
	assert.NotEmpty(t, errs)
}

func TestBoxedValuerLikePatternParse(t *testing.T) {
	// Similar to database/sql driver.Valuer pattern
	s := Boxed(
		String().Min(3),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	var valuer StringValuer
	errs := s.Parse("hello world", &valuer)
	assert.Empty(t, errs)
	val, _ := valuer.Value()
	assert.Equal(t, "hello world", val)
}

// ============================================================================
// 8. Catch Functionality
// ============================================================================

func TestBoxedStringWithCatchParse(t *testing.T) {
	s := Boxed(
		String().Min(3).Catch("caught"),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("x", &box) // too short, should trigger catch
	assert.Empty(t, errs)      // Catch should suppress errors
	assert.Equal(t, "caught", box.V)
}

func TestBoxedStringWithCatchParseSuccess(t *testing.T) {
	s := Boxed(
		String().Min(3).Catch("caught"),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("hello", &box) // valid, should not trigger catch
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V)
}

func TestBoxedIntWithCatchParse(t *testing.T) {
	s := Boxed(
		Int().GT(0).Catch(42),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	errs := s.Parse(-1, &box) // violates GT(0), should trigger catch
	assert.Empty(t, errs)
	assert.Equal(t, 42, box.Value)
}

func TestBoxedStringValuerWithCatchParse(t *testing.T) {
	s := Boxed(
		String().Min(5).Catch("caught"),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	var valuer StringValuer
	errs := s.Parse("hi", &valuer) // too short, should trigger catch
	assert.Empty(t, errs)
	val, _ := valuer.Value()
	assert.Equal(t, "caught", val)
}

func TestBoxedSliceWithInnerCatchParse(t *testing.T) {
	// Test Catch on the inner string schema within a slice
	s := Boxed(
		Slice(String().Min(5).Catch("caught")),
		func(b SliceBox, ctx Ctx) ([]string, error) { return b.Value, nil },
		func(v []string, ctx Ctx) (SliceBox, error) { return SliceBox{Value: v}, nil },
	)

	var box SliceBox
	errs := s.Parse([]string{"hi"}, &box) // element too short, should trigger catch on inner schema
	assert.Empty(t, errs)                 // Catch should suppress errors
	assert.Equal(t, []string{"caught"}, box.Value)
}

// ============================================================================
// 9. Pointer Schema Parsing
// ============================================================================

func TestBoxedPtrStringParse(t *testing.T) {
	s := Boxed(
		Ptr(String().Min(3)),
		func(b StringPtrBox, ctx Ctx) (*string, error) { return b.Value, nil },
		func(s *string, ctx Ctx) (StringPtrBox, error) { return StringPtrBox{Value: s}, nil },
	)

	var box StringPtrBox
	errs := s.Parse("hello", &box)
	assert.Empty(t, errs)
	assert.NotNil(t, box.Value)
	assert.Equal(t, "hello", *box.Value)
}

func TestBoxedPtrStringNilParse(t *testing.T) {
	s := Boxed(
		Ptr(String().Min(3)),
		func(b StringPtrBox, ctx Ctx) (*string, error) { return b.Value, nil },
		func(s *string, ctx Ctx) (StringPtrBox, error) { return StringPtrBox{Value: s}, nil },
	)

	var box StringPtrBox
	errs := s.Parse(nil, &box)
	assert.Empty(t, errs) // nil is valid for optional pointer
}

func TestBoxedPtrStringNotNilParse(t *testing.T) {
	s := Boxed(
		Ptr(String().Min(3)).NotNil(),
		func(b StringPtrBox, ctx Ctx) (*string, error) { return b.Value, nil },
		func(s *string, ctx Ctx) (StringPtrBox, error) { return StringPtrBox{Value: s}, nil },
	)

	var box StringPtrBox
	errs := s.Parse(nil, &box)
	assert.NotEmpty(t, errs) // nil is invalid when NotNil is required
}

func TestBoxedPtrStringParseFailure(t *testing.T) {
	s := Boxed(
		Ptr(String().Min(5)),
		func(b StringPtrBox, ctx Ctx) (*string, error) { return b.Value, nil },
		func(s *string, ctx Ctx) (StringPtrBox, error) { return StringPtrBox{Value: s}, nil },
	)

	var box StringPtrBox
	errs := s.Parse("hi", &box) // too short
	assert.NotEmpty(t, errs)
}

func TestBoxedPtrIntParse(t *testing.T) {
	s := Boxed(
		Ptr(Int().GT(0)),
		func(b IntPtrBox, ctx Ctx) (*int, error) { return b.Value, nil },
		func(v *int, ctx Ctx) (IntPtrBox, error) { return IntPtrBox{Value: v}, nil },
	)

	var box IntPtrBox
	errs := s.Parse(42, &box)
	assert.Empty(t, errs)
	assert.NotNil(t, box.Value)
	assert.Equal(t, 42, *box.Value)
}

func TestBoxedPtrIntNilParse(t *testing.T) {
	s := Boxed(
		Ptr(Int().GT(0)),
		func(b IntPtrBox, ctx Ctx) (*int, error) { return b.Value, nil },
		func(v *int, ctx Ctx) (IntPtrBox, error) { return IntPtrBox{Value: v}, nil },
	)

	var box IntPtrBox
	errs := s.Parse(nil, &box)
	assert.Empty(t, errs) // nil is valid for optional pointer
}

// ============================================================================
// 10. Transform Tests
// ============================================================================

func TestBoxedStringWithTrimParse(t *testing.T) {
	s := Boxed(
		String().Trim().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("  hello  ", &box)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V) // Trim transform should propagate back to box
}

func TestBoxedStringWithTrimParseFromBox(t *testing.T) {
	s := Boxed(
		String().Trim().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	inputBox := StringBox{V: "  hello  "}
	errs := s.Parse(inputBox, &box)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V) // Trim transform should propagate back to box
}

func TestBoxedStringWithTrimValuerParse(t *testing.T) {
	s := Boxed(
		String().Trim().Min(3),
		func(b StringValuerBox, ctx Ctx) (string, error) { return b.Value(), nil },
		func(s string, ctx Ctx) (StringValuerBox, error) { return &StringBox{V: s}, nil },
	)

	var box StringValuerBox
	errs := s.Parse("  hello  ", &box)
	assert.Empty(t, errs)
	x := box.Value()
	assert.Equal(t, "hello", x) // Trim transform should propagate back to box
}

func TestBoxedStringWithTransformParse(t *testing.T) {
	s := Boxed(
		String().Min(3).Transform(func(val *string, ctx Ctx) error {
			*val = strings.ToUpper(*val)
			return nil
		}),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("hello", &box)
	assert.Empty(t, errs)
	assert.Equal(t, "HELLO", box.V) // Transform should propagate back to box
}

func TestBoxedIntWithTransformParse(t *testing.T) {
	s := Boxed(
		Int().GT(0).Transform(func(val *int, ctx Ctx) error {
			*val = *val * 2
			return nil
		}),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	errs := s.Parse(5, &box)
	assert.Empty(t, errs)
	assert.Equal(t, 10, box.Value) // Transform should propagate back to box
}

// ============================================================================
// 11. Default Value Tests
// ============================================================================

func TestBoxedStringWithDefaultParse(t *testing.T) {
	s := Boxed(
		String().Default("default").Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse(nil, &box) // nil value, should use default
	assert.Empty(t, errs)
	assert.Equal(t, "default", box.V) // Default value should propagate back to box
}

func TestBoxedStringWithDefaultParseNonZero(t *testing.T) {
	s := Boxed(
		String().Default("default").Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("hello", &box) // non-zero value, should not use default
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V)
}

func TestBoxedIntWithDefaultParse(t *testing.T) {
	s := Boxed(
		Int().Default(42).GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	errs := s.Parse(nil, &box) // nil value, should use default
	assert.Empty(t, errs)
	assert.Equal(t, 42, box.Value) // Default value should propagate back to box
}

func TestBoxedStringValuerWithDefaultParse(t *testing.T) {
	s := Boxed(
		String().Default("default").Min(3),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	var valuer StringValuer
	errs := s.Parse(nil, &valuer) // nil value, should use default
	assert.Empty(t, errs)
	val, _ := valuer.Value()
	assert.Equal(t, "default", val) // Default value should propagate back to valuer
}

// ============================================================================
// 12. Required Tests
// ============================================================================

func TestBoxedStringRequiredParse(t *testing.T) {
	s := Boxed(
		String().Required().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("", &box) // zero value, should fail required
	assert.NotEmpty(t, errs)
}

func TestBoxedStringRequiredParseValid(t *testing.T) {
	s := Boxed(
		String().Required().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("hello", &box) // valid value
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.V)
}

func TestBoxedStringOptionalParse(t *testing.T) {
	s := Boxed(
		String().Optional().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse(nil, &box) // nil value, should be valid for optional
	assert.Empty(t, errs)
}

func TestBoxedIntRequiredParse(t *testing.T) {
	s := Boxed(
		Int().Required().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	errs := s.Parse(0, &box) // zero value, should fail required
	assert.NotEmpty(t, errs)
}

func TestBoxedIntRequiredParseValid(t *testing.T) {
	s := Boxed(
		Int().Required().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	var box IntBox
	errs := s.Parse(42, &box) // valid value
	assert.Empty(t, errs)
	assert.Equal(t, 42, box.Value)
}

func TestBoxedStringRequiredWithCatchParse(t *testing.T) {
	s := Boxed(
		String().Required().Min(3).Catch("caught"),
		func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)

	var box StringBox
	errs := s.Parse("", &box) // zero value, should trigger catch
	assert.Empty(t, errs)     // Catch should suppress required error
	assert.Equal(t, "caught", box.V)
}

// ============================================================================
// 13. Boxed Schema Inside Other Schemas
// ============================================================================

func TestBoxedSchemaInsideStructParse(t *testing.T) {
	s := Struct(Shape{
		"BoxedField": Boxed(
			String().Min(3),
			func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
			func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
		),
		"OtherField": String().Min(1),
	})

	var container ContainerStruct
	input := map[string]any{
		"BoxedField": "hello",
		"OtherField": "test",
	}
	errs := s.Parse(input, &container)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", container.BoxedField.V)
	assert.Equal(t, "test", container.OtherField)
}

func TestBoxedSchemaInsideStructParseFromBox(t *testing.T) {
	s := Struct(Shape{
		"BoxedField": Boxed(
			String().Min(3),
			func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
			func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
		),
		"OtherField": String().Min(1),
	})

	var container ContainerStruct
	input := map[string]any{
		"BoxedField": StringBox{V: "hello"},
		"OtherField": "test",
	}
	errs := s.Parse(input, &container)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", container.BoxedField.V)
	assert.Equal(t, "test", container.OtherField)
}

func TestBoxedSchemaInsideStructWithCatchParse(t *testing.T) {
	s := Struct(Shape{
		"BoxedField": Boxed(
			String().Min(5).Catch("caught"),
			func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
			func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
		),
		"OtherField": String().Min(1),
	})

	var container ContainerStruct
	input := map[string]any{
		"BoxedField": "hi", // too short, should trigger catch
		"OtherField": "test",
	}
	errs := s.Parse(input, &container)
	assert.Empty(t, errs)                             // Catch should suppress errors
	assert.Equal(t, "caught", container.BoxedField.V) // Catch value should propagate back
	assert.Equal(t, "test", container.OtherField)
}

func TestBoxedSchemaInsideStructWithTransformParse(t *testing.T) {
	s := Struct(Shape{
		"BoxedField": Boxed(
			String().Min(3).Transform(func(val *string, ctx Ctx) error {
				*val = strings.ToUpper(*val)
				return nil
			}),
			func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
			func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
		),
		"OtherField": String().Min(1),
	})

	var container ContainerStruct
	input := map[string]any{
		"BoxedField": "hello",
		"OtherField": "test",
	}
	errs := s.Parse(input, &container)
	assert.Empty(t, errs)
	assert.Equal(t, "HELLO", container.BoxedField.V) // Transform should propagate back
	assert.Equal(t, "test", container.OtherField)
}

func TestBoxedSchemaInsideStructParseFailure(t *testing.T) {
	s := Struct(Shape{
		"BoxedField": Boxed(
			String().Min(5),
			func(b StringBox, ctx Ctx) (string, error) { return b.V, nil },
			func(s string, ctx Ctx) (StringBox, error) { return StringBox{V: s}, nil },
		),
		"OtherField": String().Min(1),
	})

	var container ContainerStruct
	input := map[string]any{
		"BoxedField": "hi", // too short
		"OtherField": "test",
	}
	errs := s.Parse(input, &container)
	assert.NotEmpty(t, errs) // Should have validation error
}

func TestBoxedSchemaInsideStructWithSliceParse(t *testing.T) {
	s := Struct(Shape{
		"BoxedSlice": Boxed(
			Slice(String().Min(10).Catch("xyz")),
			func(b ValuerBox, ctx Ctx) ([]string, error) { return b.Value() },
			func(v []string, ctx Ctx) (ValuerBox, error) {
				var x ValuerBox = &myValuerBox{v: v}
				return x, nil
			},
		),
		"OtherField": String().Min(1),
	})

	var container ContainerWithSlice
	input := map[string]any{
		"BoxedSlice": []string{"hello", "world"},
		"OtherField": "test",
	}
	errs := s.Parse(input, &container)
	assert.Empty(t, errs)
	v, _ := container.BoxedSlice.Value()
	assert.Equal(t, []string{"xyz", "xyz"}, v)
	assert.Equal(t, "test", container.OtherField)
}

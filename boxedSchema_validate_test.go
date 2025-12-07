package zog

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Struct box types for testing
type StringBox struct {
	Value string
}

type BoolBox struct {
	Value bool
}

type IntBox struct {
	Value int
}

type Float64Box struct {
	Value float64
}

type TimeBox struct {
	Value time.Time
}

type SliceBox struct {
	Value []string
}

type BoxedUser struct {
	Id   string
	Name string
}

type UserBox struct {
	Value BoxedUser
}

// Interface box types (Valuer pattern)
type StringValuer interface {
	Value() (string, error)
}

type IntValuer interface {
	Value() (int, error)
}

type SliceValuer interface {
	Value() ([]string, error)
}

// Implementations of Valuer interfaces
type myStringValuer struct {
	v string
}

func (m myStringValuer) Value() (string, error) {
	return m.v, nil
}

type myIntValuer struct {
	v int
}

func (m myIntValuer) Value() (int, error) {
	return m.v, nil
}

type mySliceValuer struct {
	v []string
}

func (m mySliceValuer) Value() ([]string, error) {
	return m.v, nil
}

// Error implementations for testing
type errorStringValuer struct {
	v string
}

func (e errorStringValuer) Value() (string, error) {
	return "", errors.New("unbox error")
}

// Nullable type pattern (like sql.NullString)
type NullString struct {
	String string
	Valid  bool
}

// ============================================================================
// 1. Primitive Type Boxing (struct box)
// ============================================================================

func TestBoxedStringValidate(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringBox, ctx Ctx) (string, error) { return b.Value, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{Value: s}, nil },
	)

	box := StringBox{Value: "hello"}
	errs := s.Validate(box)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", box.Value) // Original unchanged
}

func TestBoxedStringValidateFailure(t *testing.T) {
	s := Boxed(
		String().Min(5),
		func(b StringBox, ctx Ctx) (string, error) { return b.Value, nil },
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{Value: s}, nil },
	)

	box := StringBox{Value: "hi"} // too short
	errs := s.Validate(box)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "hi", box.Value) // Original unchanged
}

func TestBoxedBoolValidate(t *testing.T) {
	s := Boxed(
		Bool(),
		func(b BoolBox, ctx Ctx) (bool, error) { return b.Value, nil },
		func(v bool, ctx Ctx) (BoolBox, error) { return BoolBox{Value: v}, nil },
	)

	box := BoolBox{Value: true}
	errs := s.Validate(box)
	assert.Empty(t, errs)
	assert.Equal(t, true, box.Value) // Original unchanged
}

func TestBoxedIntValidate(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	box := IntBox{Value: 42}
	errs := s.Validate(box)
	assert.Empty(t, errs)
	assert.Equal(t, 42, box.Value) // Original unchanged
}

func TestBoxedIntValidateFailure(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntBox, ctx Ctx) (int, error) { return b.Value, nil },
		func(v int, ctx Ctx) (IntBox, error) { return IntBox{Value: v}, nil },
	)

	box := IntBox{Value: -1} // violates GT(0)
	errs := s.Validate(box)
	assert.NotEmpty(t, errs)
	assert.Equal(t, -1, box.Value) // Original unchanged
}

func TestBoxedFloat64Validate(t *testing.T) {
	s := Boxed(
		Float64().GT(0),
		func(b Float64Box, ctx Ctx) (float64, error) { return b.Value, nil },
		func(v float64, ctx Ctx) (Float64Box, error) { return Float64Box{Value: v}, nil },
	)

	box := Float64Box{Value: 3.14}
	errs := s.Validate(box)
	assert.Empty(t, errs)
	assert.Equal(t, 3.14, box.Value) // Original unchanged
}

func TestBoxedTimeValidate(t *testing.T) {
	s := Boxed(
		Time(),
		func(b TimeBox, ctx Ctx) (time.Time, error) { return b.Value, nil },
		func(v time.Time, ctx Ctx) (TimeBox, error) { return TimeBox{Value: v}, nil },
	)

	timestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	box := TimeBox{Value: timestamp}
	errs := s.Validate(box)
	assert.Empty(t, errs)
	assert.Equal(t, timestamp, box.Value) // Original unchanged
}

// ============================================================================
// 2. Primitive Type Boxing (interface box)
// ============================================================================

func TestBoxedStringValuerValidate(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	valuer := myStringValuer{v: "hello"}
	errs := s.Validate(valuer)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", valuer.v) // Original unchanged
}

func TestBoxedStringValuerValidateFailure(t *testing.T) {
	s := Boxed(
		String().Min(5),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	valuer := myStringValuer{v: "hi"} // too short
	errs := s.Validate(valuer)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "hi", valuer.v) // Original unchanged
}

func TestBoxedIntValuerValidate(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntValuer, ctx Ctx) (int, error) { return b.Value() },
		func(v int, ctx Ctx) (IntValuer, error) { return myIntValuer{v: v}, nil },
	)

	valuer := myIntValuer{v: 42}
	errs := s.Validate(valuer)
	assert.Empty(t, errs)
	assert.Equal(t, 42, valuer.v) // Original unchanged
}

func TestBoxedIntValuerValidateFailure(t *testing.T) {
	s := Boxed(
		Int().GT(0),
		func(b IntValuer, ctx Ctx) (int, error) { return b.Value() },
		func(v int, ctx Ctx) (IntValuer, error) { return myIntValuer{v: v}, nil },
	)

	valuer := myIntValuer{v: -1} // violates GT(0)
	errs := s.Validate(valuer)
	assert.NotEmpty(t, errs)
	assert.Equal(t, -1, valuer.v) // Original unchanged
}

// ============================================================================
// 3. Complex Type Boxing
// ============================================================================

func TestBoxedSliceValidate(t *testing.T) {
	s := Boxed(
		Slice(String().Min(1)),
		func(b SliceBox, ctx Ctx) ([]string, error) { return b.Value, nil },
		func(v []string, ctx Ctx) (SliceBox, error) { return SliceBox{Value: v}, nil },
	)

	box := SliceBox{Value: []string{"hello", "world"}}
	errs := s.Validate(box)
	assert.Empty(t, errs)
	assert.Equal(t, []string{"hello", "world"}, box.Value) // Original unchanged
}

func TestBoxedSliceValidateFailure(t *testing.T) {
	s := Boxed(
		Slice(String().Min(5)),
		func(b SliceBox, ctx Ctx) ([]string, error) { return b.Value, nil },
		func(v []string, ctx Ctx) (SliceBox, error) { return SliceBox{Value: v}, nil },
	)

	box := SliceBox{Value: []string{"hi"}} // too short
	errs := s.Validate(box)
	assert.NotEmpty(t, errs)
	assert.Equal(t, []string{"hi"}, box.Value) // Original unchanged
}

func TestBoxedStructValidate(t *testing.T) {
	s := Boxed(
		Struct(Shape{
			"Id":   String().Min(1),
			"Name": String().Min(1),
		}),
		func(b UserBox, ctx Ctx) (BoxedUser, error) { return b.Value, nil },
		func(u BoxedUser, ctx Ctx) (UserBox, error) { return UserBox{Value: u}, nil },
	)

	box := UserBox{Value: BoxedUser{Id: "1", Name: "John Doe"}}
	errs := s.Validate(box)
	assert.Empty(t, errs)
	assert.Equal(t, BoxedUser{Id: "1", Name: "John Doe"}, box.Value) // Original unchanged
}

func TestBoxedStructValidateFailure(t *testing.T) {
	s := Boxed(
		Struct(Shape{
			"Id":   String().Min(1),
			"Name": String().Min(5), // Name must be at least 5 chars
		}),
		func(b UserBox, ctx Ctx) (BoxedUser, error) { return b.Value, nil },
		func(u BoxedUser, ctx Ctx) (UserBox, error) { return UserBox{Value: u}, nil },
	)

	box := UserBox{Value: BoxedUser{Id: "1", Name: "Joe"}} // Name too short
	errs := s.Validate(box)
	assert.NotEmpty(t, errs)
	assert.Equal(t, BoxedUser{Id: "1", Name: "Joe"}, box.Value) // Original unchanged
}

func TestBoxedSliceValuerValidate(t *testing.T) {
	s := Boxed(
		Slice(String().Min(1)),
		func(b SliceValuer, ctx Ctx) ([]string, error) { return b.Value() },
		func(v []string, ctx Ctx) (SliceValuer, error) { return mySliceValuer{v: v}, nil },
	)

	valuer := mySliceValuer{v: []string{"hello", "world"}}
	errs := s.Validate(valuer)
	assert.Empty(t, errs)
	assert.Equal(t, []string{"hello", "world"}, valuer.v) // Original unchanged
}

// ============================================================================
// 4. Error Handling
// ============================================================================

func TestBoxedUnboxErrorStruct(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringBox, ctx Ctx) (string, error) {
			if b.Value == "" {
				return "", errors.New("cannot unbox empty string")
			}
			return b.Value, nil
		},
		func(s string, ctx Ctx) (StringBox, error) { return StringBox{Value: s}, nil },
	)

	box := StringBox{Value: ""}
	errs := s.Validate(box)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "", box.Value) // Original unchanged
}

func TestBoxedUnboxErrorInterface(t *testing.T) {
	s := Boxed(
		String().Min(3),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return errorStringValuer{v: s}, nil },
	)

	valuer := errorStringValuer{v: "hello"}
	errs := s.Validate(valuer)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "hello", valuer.v) // Original unchanged
}

// ============================================================================
// 5. Real-World Patterns
// ============================================================================

func TestBoxedNullablePattern(t *testing.T) {
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

	// Valid nullable string
	ns := NullString{String: "hello", Valid: true}
	errs := s.Validate(ns)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", ns.String) // Original unchanged
	assert.Equal(t, true, ns.Valid)     // Original unchanged

	// Invalid nullable string (Valid = false)
	ns2 := NullString{String: "hello", Valid: false}
	errs2 := s.Validate(ns2)
	assert.NotEmpty(t, errs2)
	assert.Equal(t, false, ns2.Valid) // Original unchanged
}

func TestBoxedNullablePatternValidationFailure(t *testing.T) {
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
	ns := NullString{String: "hi", Valid: true}
	errs := s.Validate(ns)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "hi", ns.String) // Original unchanged
	assert.Equal(t, true, ns.Valid)  // Original unchanged
}

func TestBoxedValuerLikePattern(t *testing.T) {
	// Similar to database/sql driver.Valuer pattern
	s := Boxed(
		String().Min(3),
		func(b StringValuer, ctx Ctx) (string, error) { return b.Value() },
		func(s string, ctx Ctx) (StringValuer, error) { return myStringValuer{v: s}, nil },
	)

	valuer := myStringValuer{v: "hello world"}
	errs := s.Validate(valuer)
	assert.Empty(t, errs)
	assert.Equal(t, "hello world", valuer.v) // Original unchanged
}

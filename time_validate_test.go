package zog

import (
	"testing"
	"time"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestTimeValidateRequired(t *testing.T) {
	validator := Time().Required(Message("custom"))
	now := time.Now()
	errs := validator.Validate(&now)
	assert.Nil(t, errs)

	now = time.Time{}
	errs = validator.Validate(&now)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message())
}

func TestTimeValidateOptional(t *testing.T) {
	now := time.Now()
	validator := Time().Optional()
	errs := validator.Validate(&now)
	assert.Nil(t, errs)

	var zeroTime time.Time
	errs = validator.Validate(&zeroTime)
	assert.Nil(t, errs)
}

func TestTimeValidateDefault(t *testing.T) {
	defaultVal := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	validator := Time().Default(defaultVal)
	var now time.Time
	errs := validator.Validate(&now)
	assert.Nil(t, errs)
	assert.Equal(t, defaultVal, now)
}

func TestTimeValidateCatch(t *testing.T) {
	catchVal := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	validator := Time().Required().Catch(catchVal)
	var now time.Time
	errs := validator.Validate(&now)
	assert.Nil(t, errs)
	assert.Equal(t, catchVal, now)
}

func TestTimeValidatePreTransform(t *testing.T) {
	validator := Time().PreTransform(func(data any, ctx ParseCtx) (any, error) {
		// Add 1 hour to the input time
		t, ok := data.(*time.Time)
		if !ok {
			return nil, nil
		}
		result := t.Add(time.Hour)
		return &result, nil
	})

	input := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	expected := input.Add(time.Hour)
	errs := validator.Validate(&input)
	assert.Nil(t, errs)
	assert.Equal(t, expected, input)
}

func TestTimeValidatePostTransform(t *testing.T) {
	validator := Time().PostTransform(func(dataPtr any, ctx ParseCtx) error {
		// Set the time to noon
		t := dataPtr.(*time.Time)
		*t = time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
		return nil
	})

	input := time.Date(2023, 1, 1, 15, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	errs := validator.Validate(&input)
	assert.Nil(t, errs)
	assert.Equal(t, expected, input)
}

func TestTimeValidateAfter(t *testing.T) {
	now := time.Now()
	validator := Time().After(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Message("custom"))
	errs := validator.Validate(&now)
	assert.Nil(t, errs)

	past := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = validator.Validate(&past)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message())
}

func TestTimeValidateBefore(t *testing.T) {
	past := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	validator := Time().Before(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Message("custom"))
	errs := validator.Validate(&past)
	assert.Nil(t, errs)

	future := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = validator.Validate(&future)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message())
}

func TestTimeValidateEQ(t *testing.T) {
	target := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	validator := Time().EQ(target, Message("custom"))
	errs := validator.Validate(&target)
	assert.Nil(t, errs)

	different := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = validator.Validate(&different)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message())
}

func TestTimeValidateCustomTest(t *testing.T) {
	now := time.Now()
	validator := Time().Test(TestFunc("custom_test", func(val any, ctx ParseCtx) bool {
		return val != now
	}), Message("custom"))
	errs := validator.Validate(&now)
	assert.NotNil(t, errs)
	assert.Equal(t, "custom", errs[0].Message())
}

func TestTimeValidateGetType(t *testing.T) {
	validator := Time()
	assert.Equal(t, zconst.TypeTime, validator.getType())
}

package zog

import (
	"testing"
	"time"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestTimeRequired(t *testing.T) {
	var now time.Time
	schema := Time().Required(Message("custom"))
	errs := schema.Parse(time.Now(), &now)
	assert.Nil(t, errs)
	now = time.Time{}
	errs = schema.Parse(nil, &now)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message)
}

func TestTimeOptional(t *testing.T) {
	now := time.Now()
	schema := Time().Optional()
	errs := schema.Parse(nil, &now)
	assert.Nil(t, errs)
	errs = schema.Parse(now, &now)
	assert.Nil(t, errs)
}

func TestTimeDefault(t *testing.T) {
	var now time.Time
	defaultVal := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Time().Default(defaultVal)
	errs := schema.Parse(nil, &now)
	assert.Nil(t, errs)
	assert.Equal(t, defaultVal, now)
}

func TestTimeCatch(t *testing.T) {
	var now time.Time
	catchVal := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Time().Required().Catch(catchVal)
	errs := schema.Parse(nil, &now)
	assert.Nil(t, errs)
	assert.Equal(t, catchVal, now)
}

func TestTimePreTransform(t *testing.T) {
	var now time.Time
	schema := Time().PreTransform(func(data any, ctx ParseCtx) (any, error) {
		// Add 1 hour to the input time
		t, ok := data.(time.Time)
		if !ok {
			return nil, nil
		}
		return t.Add(time.Hour), nil
	})

	input := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	expected := input.Add(time.Hour)

	errs := schema.Parse(input, &now)
	assert.Nil(t, errs)
	assert.Equal(t, expected, now)
}

func TestTimePostTransform(t *testing.T) {
	var now time.Time
	schema := Time().PostTransform(func(dataPtr any, ctx ParseCtx) error {
		// Set the time to noon
		t := dataPtr.(*time.Time)
		*t = time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
		return nil
	})

	input := time.Date(2023, 1, 1, 15, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	errs := schema.Parse(input, &now)
	assert.Nil(t, errs)
	assert.Equal(t, expected, now)
}

// VALIDATORS

func TestTimeAfter(t *testing.T) {
	now := time.Now()

	schema := Time().After(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Message("custom"))
	errs := schema.Parse(now, &now)
	assert.Nil(t, errs)

	now = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = schema.Parse(now, &now)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message)
}
func TestTimeBefore(t *testing.T) {
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Time().Before(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Message("custom"))
	errs := schema.Parse(now, &now)
	assert.Nil(t, errs)
	now = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = schema.Parse(now, &now)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message)
}

func TestTimeEQ(t *testing.T) {
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Time().EQ(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Message("custom"))
	errs := schema.Parse(now, &now)
	assert.Nil(t, errs)
	now = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = schema.Parse(now, &now)
	assert.Len(t, errs, 1)
	assert.Equal(t, "custom", errs[0].Message)
}

func TestTimeCustomTest(t *testing.T) {
	now := time.Now()
	schema := Time().TestFunc(func(val any, ctx ParseCtx) bool {
		return val != now
	}, Message("custom"))
	errs := schema.Parse(now, &now)
	assert.NotNil(t, errs)
	assert.Equal(t, "custom", errs[0].Message)
}

func TestTimeSchemaOption(t *testing.T) {
	s := Time(WithCoercer(func(original any) (value any, err error) {
		return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil
	}))

	var result time.Time
	err := s.Parse("invalid-date", &result)
	assert.Nil(t, err)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), result)
}

func TestTimeFormat(t *testing.T) {
	s := Time(Time.Format(time.RFC1123))
	var result time.Time
	err := s.Parse("Mon, 01 Jan 2024 00:00:00 UTC", &result)
	assert.Nil(t, err)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), result)
}

func TestTimeFormatFunc(t *testing.T) {
	s := Time(Time.FormatFunc(func(data string) (time.Time, error) {
		return time.Parse(time.RFC1123, data)
	}))
	var result time.Time
	err := s.Parse("Mon, 01 Jan 2024 00:00:00 UTC", &result)
	assert.Nil(t, err)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), result)
}

func TestTimeGetType(t *testing.T) {
	s := Time()
	assert.Equal(t, zconst.TypeTime, s.getType())
}

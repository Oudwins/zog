package zog

import (
	"fmt"
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

func TestTimePostTransform(t *testing.T) {
	var now time.Time
	schema := Time().Transform(func(dataPtr *time.Time, ctx Ctx) error {
		// Set the time to noon
		*dataPtr = time.Date(dataPtr.Year(), dataPtr.Month(), dataPtr.Day(), 12, 0, 0, 0, dataPtr.Location())
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
	schema := Time().TestFunc(func(val *time.Time, ctx Ctx) bool {
		return !(*val).Equal(now)
	}, Message("custom"))
	errs := schema.Parse(now, &now)
	assert.NotNil(t, errs)
	assert.Equal(t, "custom", errs[0].Message)
}

// FORMAT TESTS

func TestTimeFormatRFC3339(t *testing.T) {
	var parsed time.Time
	schema := Time(Time.Format(time.RFC3339))

	// Test valid RFC3339 string
	errs := schema.Parse("2024-01-01T12:30:45Z", &parsed)
	assert.Nil(t, errs)
	expected := time.Date(2024, 1, 1, 12, 30, 45, 0, time.UTC)
	assert.Equal(t, expected, parsed)

	// Test invalid RFC3339 string
	errs = schema.Parse("invalid-date", &parsed)
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
}

func TestTimeFormatCustom(t *testing.T) {
	var parsed time.Time
	customFormat := "2006-01-02 15:04:05"
	schema := Time(Time.Format(customFormat))

	// Test valid custom format string
	errs := schema.Parse("2024-01-01 12:30:45", &parsed)
	assert.Nil(t, errs)
	expected := time.Date(2024, 1, 1, 12, 30, 45, 0, time.UTC)
	assert.Equal(t, expected, parsed)

	// Test invalid format string
	errs = schema.Parse("2024/01/01 12:30:45", &parsed)
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
}

func TestTimeFormatDateOnly(t *testing.T) {
	var parsed time.Time
	dateFormat := "2006-01-02"
	schema := Time(Time.Format(dateFormat))

	// Test date-only format
	errs := schema.Parse("2024-12-25", &parsed)
	assert.Nil(t, errs)
	expected := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, parsed)
}

func TestTimeFormatTimeOnly(t *testing.T) {
	var parsed time.Time
	timeFormat := "15:04:05"
	schema := Time(Time.Format(timeFormat))

	// Test time-only format (year will be 0000)
	errs := schema.Parse("14:30:00", &parsed)
	assert.Nil(t, errs)
	expected := time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC)
	assert.Equal(t, expected, parsed)
}

func TestTimeFormatWithTimezone(t *testing.T) {
	var parsed time.Time
	formatWithTZ := "2006-01-02 15:04:05 MST"
	schema := Time(Time.Format(formatWithTZ))

	// Test format with timezone
	errs := schema.Parse("2024-01-01 12:30:45 EST", &parsed)
	assert.Nil(t, errs)
	// Verify basic time components were parsed correctly
	assert.Equal(t, 2024, parsed.Year())
	assert.Equal(t, time.January, parsed.Month())
	assert.Equal(t, 1, parsed.Day())
	assert.Equal(t, 12, parsed.Hour())
	assert.Equal(t, 30, parsed.Minute())
	assert.Equal(t, 45, parsed.Second())
	// Verify timezone was parsed
	zoneName, _ := parsed.Zone()
	assert.Equal(t, "EST", zoneName)
}

func TestTimeFormatWithExistingTimeInput(t *testing.T) {
	var parsed time.Time
	schema := Time(Time.Format(time.RFC3339))
	now := time.Now()

	// Test that time.Time input still works with custom format
	errs := schema.Parse(now, &parsed)
	assert.Nil(t, errs)
	assert.Equal(t, now, parsed)
}

func TestTimeFormatWithIntInput(t *testing.T) {
	var parsed time.Time
	schema := Time(Time.Format(time.RFC3339))

	// Test that unix timestamp still works with custom format
	timestamp := int64(1640995200) // 2022-01-01 00:00:00 UTC
	errs := schema.Parse(timestamp, &parsed)
	assert.Nil(t, errs)
	expected := time.Unix(timestamp, 0)
	assert.Equal(t, expected, parsed)
}

func TestTimeFormatWithRequiredValidation(t *testing.T) {
	var parsed time.Time
	schema := Time(Time.Format("2006-01-02")).Required(Message("time required"))

	// Test valid input
	errs := schema.Parse("2024-01-01", &parsed)
	assert.Nil(t, errs)

	// Test nil input
	errs = schema.Parse(nil, &parsed)
	assert.NotNil(t, errs)
	assert.Equal(t, "time required", errs[0].Message)
}

func TestTimeFormatFunc(t *testing.T) {
	var parsed time.Time

	// Custom format function that parses "DD/MM/YYYY HH:mm" format
	formatFunc := func(data string) (time.Time, error) {
		return time.Parse("02/01/2006 15:04", data)
	}

	schema := Time(Time.FormatFunc(formatFunc))

	// Test valid custom format
	errs := schema.Parse("25/12/2024 14:30", &parsed)
	assert.Nil(t, errs)
	expected := time.Date(2024, 12, 25, 14, 30, 0, 0, time.UTC)
	assert.Equal(t, expected, parsed)

	// Test invalid format
	errs = schema.Parse("2024-12-25 14:30", &parsed)
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
}

func TestTimeFormatFuncWithError(t *testing.T) {
	var parsed time.Time

	// Custom format function that always returns an error
	formatFunc := func(data string) (time.Time, error) {
		return time.Time{}, fmt.Errorf("custom parsing error: %s", data)
	}

	schema := Time(Time.FormatFunc(formatFunc))

	// Test that custom error is propagated
	errs := schema.Parse("any-string", &parsed)
	assert.NotNil(t, errs)
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Err.Error(), "custom parsing error")
}

func TestTimeFormatFuncComplexParsing(t *testing.T) {
	var parsed time.Time

	// Custom format function that handles multiple formats
	formatFunc := func(data string) (time.Time, error) {
		formats := []string{
			"2006-01-02",
			"01/02/2006",
			"2006-01-02 15:04:05",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, data); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("unable to parse date: %s", data)
	}

	schema := Time(Time.FormatFunc(formatFunc))

	// Test multiple valid formats
	testCases := []struct {
		input    string
		expected time.Time
	}{
		{"2024-01-01", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"01/15/2024", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{"2024-01-01 12:30:45", time.Date(2024, 1, 1, 12, 30, 45, 0, time.UTC)},
	}

	for _, tc := range testCases {
		errs := schema.Parse(tc.input, &parsed)
		assert.Nil(t, errs, "Failed to parse: %s", tc.input)
		assert.Equal(t, tc.expected, parsed, "Unexpected result for: %s", tc.input)
	}

	// Test invalid format
	errs := schema.Parse("invalid-date-format", &parsed)
	assert.NotNil(t, errs)
	assert.Contains(t, errs[0].Err.Error(), "unable to parse date")
}

func TestTimeFormatWithValidators(t *testing.T) {
	var parsed time.Time
	minDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Time(Time.Format("2006-01-02")).After(minDate, Message("date too early"))

	// Test valid date after minimum
	errs := schema.Parse("2024-06-01", &parsed)
	assert.Nil(t, errs)

	// Test invalid date before minimum
	errs = schema.Parse("2023-12-31", &parsed)
	assert.NotNil(t, errs)
	assert.Equal(t, "date too early", errs[0].Message)
}

func TestTimeGetType(t *testing.T) {
	s := Time()
	assert.Equal(t, zconst.TypeTime, s.getType())
}

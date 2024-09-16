package zog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeAfter(t *testing.T) {
	now := time.Now()

	schema := Time().After(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	errs := schema.Parse(now, &now)
	assert.Nil(t, errs)

	now = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = schema.Parse(now, &now)
	assert.Len(t, errs, 1)
}
func TestTimeBefore(t *testing.T) {
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Time().Before(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
	errs := schema.Parse(now, &now)
	assert.Nil(t, errs)
	now = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = schema.Parse(now, &now)
	assert.Len(t, errs, 1)
}

func TestTimeEQ(t *testing.T) {
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Time().EQ(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	errs := schema.Parse(now, &now)
	assert.Nil(t, errs)
	now = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	errs = schema.Parse(now, &now)
	assert.Len(t, errs, 1)
}

func TestTimeRequired(t *testing.T) {
	var now time.Time
	schema := Time().Required()
	errs := schema.Parse(time.Now(), &now)
	assert.Nil(t, errs)
	now = time.Time{}
	errs = schema.Parse(nil, &now)
	assert.Len(t, errs, 1)
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

func TestTimeCustomTest(t *testing.T) {
	now := time.Now()
	schema := Time().Test(TestFunc("custom_test", func(val any, ctx ParseCtx) bool {
		// Custom test logic here
		return true
	}))
	errs := schema.Parse(now, &now)
	assert.Nil(t, errs)
}

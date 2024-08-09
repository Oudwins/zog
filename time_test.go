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

package zog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeAfter(t *testing.T) {
	type Foo struct {
		CreatedAt time.Time
	}
	now := time.Now()
	foo := Foo{
		CreatedAt: now,
	}
	schema := Schema{
		"createdAt": Time().After(now),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.Len(t, errors["createdAt"], 1)

	foo.CreatedAt = now.Add(time.Second * 10000)
	_, ok = Validate(foo, schema)
	assert.True(t, ok)
}

func TestTimeBefore(t *testing.T) {
	type Foo struct {
		CreatedAt time.Time
	}
	now := time.Now()
	foo := Foo{
		CreatedAt: now,
	}
	schema := Schema{
		"createdAt": Time().Before(now),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.Len(t, errors["createdAt"], 1)

	foo.CreatedAt = now.Add(time.Second * -10000)
	_, ok = Validate(foo, schema)
	assert.True(t, ok)
}

func TestTimeIs(t *testing.T) {
	type Foo struct {
		CreatedAt time.Time
	}
	now := time.Now()
	foo := Foo{
		CreatedAt: now.Add(time.Second * 10000),
	}
	schema := Schema{
		"createdAt": Time().Is(now),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.Len(t, errors["createdAt"], 1)

	foo.CreatedAt = now
	_, ok = Validate(foo, schema)
	assert.True(t, ok)
}

package zjson

import (
	"strings"
	"testing"

	z "github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestDecodeOmitEmpty(t *testing.T) {
	type User struct {
		Email string `json:"email,omitempty"`
	}

	schema := z.Struct(z.Shape{
		"email": z.String().Required(),
	})

	body := strings.NewReader(`{"other":"x"}`)

	var u User
	errs := schema.Parse(Decode(body), &u)
	assert.NotEmpty(t, errs)
	assert.Equal(t, []string{"email"}, errs[0].Path)
}

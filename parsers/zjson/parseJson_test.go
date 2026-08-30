package zjson

import (
	"bytes"
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

func TestDecodeRespectsJSONTagOnInitialismField(t *testing.T) {
	type SessionCreateRequest struct {
		Path      string `json:"path"`
		SessionID string `json:"sessionId,omitempty" zog:"sessionId"`
	}

	schema := z.Struct(z.Shape{
		"Path":      z.String().Required(),
		"SessionID": z.String().Optional(),
	})

	body := bytes.NewReader([]byte(`{"path":"/bin/x","sessionId":"XYZ-123"}`))

	var payload SessionCreateRequest
	errs := schema.Parse(Decode(body), &payload)
	assert.Empty(t, errs)
	assert.Equal(t, "/bin/x", payload.Path)
	assert.Equal(t, "XYZ-123", payload.SessionID)
}

func TestDecodeNonEmptyBodyUsesTaggedIssuePaths(t *testing.T) {
	type Person struct {
		Name  string `zog:"full_name"`
		Email string `json:"email_address"`
		Alias string `json:"public_alias" zog:"internal_alias"`
		Age   string
	}

	tests := []struct {
		name         string
		schemaKey    string
		expectedPath string
	}{
		{name: "zog tag", schemaKey: "Name", expectedPath: "full_name"},
		{name: "json tag", schemaKey: "Email", expectedPath: "email_address"},
		{name: "json tag takes precedence over zog tag", schemaKey: "Alias", expectedPath: "public_alias"},
		{name: "schema key fallback", schemaKey: "Age", expectedPath: "Age"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := z.Struct(z.Shape{
				tt.schemaKey: z.String().Required(),
			})

			var person Person
			errs := schema.Parse(Decode(strings.NewReader(`{"other":true}`)), &person)
			if assert.NotEmpty(t, errs) {
				assert.Equal(t, []string{tt.expectedPath}, errs[0].Path)
			}
		})
	}
}

func TestDecodeEmptyBodyRespectsJSONTag(t *testing.T) {
	type User struct {
		Name string `json:"full_name"`
	}

	schema := z.Struct(z.Shape{
		"name": z.String().Min(2).Required(),
	})

	var u User
	errs := schema.Parse(Decode(strings.NewReader(`{}`)), &u)
	assert.NotEmpty(t, errs)
	assert.Equal(t, []string{"full_name"}, errs[0].Path)
}

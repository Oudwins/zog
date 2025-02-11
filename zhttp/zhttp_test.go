package zhttp

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	formData := "name=JohnDoe&email=john@doe.com&age=30&isMarried=true&lights=on&cash=10.5&swagger=doweird"

	// Create a fake HTTP request with form data
	req, err := http.NewRequest("POST", "/submit", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	type User struct {
		Email     string  `zog:"email"`
		Name      string  `zog:"name"`
		Age       int     `zog:"age"`
		IsMarried bool    `zog:"isMarried"`
		Lights    bool    `zog:"lights"`
		Cash      float64 `zog:"cash"`
		Swagger   string  `zog:"swagger"`
	}
	schema := z.Struct(z.Schema{
		"email":     z.String().Email(),
		"name":      z.String().Min(3).Max(10),
		"age":       z.Int().GT(18),
		"isMarried": z.Bool().True(),
		"lights":    z.Bool().True(),
		"cash":      z.Float().GT(10.0),
		"swagger": z.String().Test(z.TestFunc("swagger", func(val any, ctx z.Ctx) bool {
			return val.(string) == "doweird"
		})),
	})
	u := User{}

	dp := Request(req)
	errs := schema.Parse(dp, &u)

	assert.Equal(t, "john@doe.com", u.Email)
	assert.Equal(t, "JohnDoe", u.Name)
	assert.Equal(t, 30, u.Age)
	assert.True(t, u.IsMarried)
	assert.True(t, u.Lights)
	assert.Equal(t, 10.5, u.Cash)
	assert.Equal(t, u.Swagger, "doweird")
	assert.Empty(t, errs)
}

func TestRequestParams(t *testing.T) {
	formData := "name=JohnDoe&email=john@doe.com&age=30&isMarried=true&lights=on&cash=10.5&swagger=doweird&swagger=swagger"

	// Create a fake HTTP request with form data
	req, err := http.NewRequest("POST", "/submit?"+formData, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	type User struct {
		Email     string   `param:"email"`
		Name      string   `param:"name"`
		Age       int      `param:"age"`
		IsMarried bool     `param:"isMarried"`
		Lights    bool     `param:"lights"`
		Cash      float64  `param:"cash"`
		Swagger   []string `param:"swagger"`
	}

	schema := z.Struct(z.Schema{
		"email":     z.String().Email(),
		"name":      z.String().Min(3).Max(10),
		"age":       z.Int().GT(18),
		"isMarried": z.Bool().True(),
		"lights":    z.Bool().True(),
		"cash":      z.Float().GT(10.0),
		"swagger": z.Slice(
			z.String().Min(1)).Min(2),
	})
	u := User{}
	dp := Request(req)
	assert.Nil(t, err)
	errs := schema.Parse(dp, &u)

	assert.Equal(t, "john@doe.com", u.Email)
	assert.Equal(t, "JohnDoe", u.Name)
	assert.Equal(t, 30, u.Age)
	assert.True(t, u.IsMarried)
	assert.True(t, u.Lights)
	assert.Equal(t, 10.5, u.Cash)
	assert.Equal(t, u.Swagger, []string{"doweird", "swagger"})
	assert.Empty(t, errs)
}

// Unit tests for url data provider
func TestUrlDataProviderGet(t *testing.T) {
	data := url.Values{
		"single":     []string{"value"},
		"multiple":   []string{"value1", "value2"},
		"array[]":    []string{"item1", "item2"},
		"emptyArray": []string{},
	}
	provider := urlDataProvider{Data: data}

	tests := []struct {
		name     string
		key      string
		expected any
	}{
		{"Single value", "single", "value"},
		{"Multiple values", "multiple", []string{"value1", "value2"}},
		{"Array notation", "array[]", []string{"item1", "item2"}},
		{"Empty array", "emptyArray", ""},
		{"Non-existent key", "nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.Get(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUrlDataProviderGetNestedProvider(t *testing.T) {
	data := url.Values{"key": []string{"value"}}
	provider := urlDataProvider{Data: data}

	nestedProvider := provider.GetNestedProvider("any_key")
	assert.Equal(t, provider, nestedProvider)
}

func TestUrlDataProviderGetUnderlying(t *testing.T) {
	data := url.Values{"key": []string{"value"}}
	provider := urlDataProvider{Data: data}

	underlying := provider.GetUnderlying()
	assert.Equal(t, data, underlying)
}

func TestRequestContentTypeJSON(t *testing.T) {
	jsonBody := `{"name":"John","age":30}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	dpFactory := Request(req)
	dp, err := dpFactory()
	assert.Nil(t, err)
	assert.Equal(t, "John", dp.Get("name"))
	assert.Equal(t, float64(30), dp.Get("age"))
}

func TestRequestContentTypeForm(t *testing.T) {
	formData := "name=John&age=30"
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dpFactory := Request(req)
	dp, err := dpFactory()

	assert.Nil(t, err)
	assert.Equal(t, "John", dp.Get("name"))
	assert.Equal(t, "30", dp.Get("age"))
}

func TestRequestContentTypeDefault(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test?name=John&age=30", nil)

	dpFactory := Request(req)
	dp, err := dpFactory()

	assert.Nil(t, err)
	assert.Equal(t, "John", dp.Get("name"))
	assert.Equal(t, "30", dp.Get("age"))
}

func TestParseJsonValid(t *testing.T) {
	jsonData := `{"name":"John","age":30}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	dp, err := Config.Parsers.JSON(req)()
	assert.Nil(t, err)
	assert.Equal(t, "John", dp.Get("name"))
	assert.Equal(t, float64(30), dp.Get("age"))
}

func TestParseJsonInvalid(t *testing.T) {
	invalidJSON := `{"name":"John","age":30`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	dp, err := Config.Parsers.JSON(req)()

	assert.Error(t, err)
	assert.Nil(t, dp)
	assert.Equal(t, zconst.IssueCodeInvalidJSON, err.Code())
}

func TestParseJsonWithNilValue(t *testing.T) {
	jsonData := `null`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	dp, err := Config.Parsers.JSON(req)()
	assert.NotNil(t, err)
	assert.Nil(t, dp)
}

func TestParseJsonWithEmptyObject(t *testing.T) {
	jsonData := `{}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	dp, err := Config.Parsers.JSON(req)()
	assert.Nil(t, err)
	assert.Nil(t, dp)
}

func TestParseDeeplyNestedJson(t *testing.T) {
	schema := z.Struct(z.Schema{
		"name": z.String().Required(),
		"nested1": z.Struct(z.Schema{
			"name": z.String().Required(),
			"nested3": z.Ptr(z.Struct(z.Schema{
				"name": z.String().Required(),
			})),
		}),
	})
	type User struct {
		Name    string `json:"name"`
		Nested1 struct {
			Name    string `json:"name"`
			Nested3 *struct {
				Name string `json:"name"`
			} `json:"nested3"`
		} `json:"nested1"`
	}

	jsonData := `{"name":"John","nested1":{"name":"nested1"}}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	user := User{}
	errs := schema.Parse(Request(req), &user)
	assert.Nil(t, errs)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "nested1", user.Nested1.Name)
	assert.Nil(t, user.Nested1.Nested3)

}

func TestTopLevelOptionalStruct(t *testing.T) {
	schema := z.Ptr(z.Struct(z.Schema{
		"name": z.String().Required(),
	}))

	type User struct {
		Name string `json:"name"`
	}

	jsonData := `{}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	user := &User{}
	errs := schema.Parse(Request(req), &user)
	assert.Nil(t, errs)
}

func TestForm(t *testing.T) {
	data := url.Values{
		"name": []string{"John"},
		"age":  []string{"30"},
	}

	dp := form(data)

	assert.IsType(t, urlDataProvider{}, dp)
	assert.Equal(t, data, dp.(urlDataProvider).Data)
}

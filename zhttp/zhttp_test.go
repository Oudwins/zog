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
	schema := z.Struct(z.Shape{
		"email":     z.String().Email(),
		"name":      z.String().Min(3).Max(10),
		"age":       z.Int().GT(18),
		"isMarried": z.Bool().True(),
		"lights":    z.Bool().True(),
		"cash":      z.Float64().GT(10.0),
		"swagger": z.String().TestFunc(func(val *string, ctx z.Ctx) bool {
			return *val == "doweird"
		}),
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
	formData := "thename=JohnDoe&mail=john@doe.com&theage=30&Married=true&light=on&money=10.5&swagger=doweird&swagger=swagger&q=test"

	// Create a fake HTTP request with query param data
	req, err := http.NewRequest("POST", "/submit?"+formData, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	type User struct {
		Email     string   `query:"mail"`
		Name      string   `query:"thename"`
		Age       int      `query:"theage"`
		IsMarried bool     `query:"Married"`
		Lights    bool     `query:"light"`
		Cash      float64  `query:"money"`
		Swagger   []string `zog:"random" query:"swagger"` // query takes priority over zog
		Q         string   `zog:"q"`
	}

	schema := z.Struct(z.Shape{
		"email":     z.String().Email(),
		"name":      z.String().Min(3).Max(10),
		"age":       z.Int().GT(18),
		"isMarried": z.Bool().True(),
		"lights":    z.Bool().True(),
		"cash":      z.Float64().GT(10.0),
		"swagger": z.Slice[string](
			z.String().Min(1)).Min(2),
		"q": z.String().Required(),
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
	assert.Equal(t, "test", u.Q)
	assert.Empty(t, errs)
}

func TestRequestParamsOnJsonContentType(t *testing.T) {
	formData := "name=JohnDoe&email=john@doe.com&age=30&isMarried=true&lights=on&cash=10.5&swagger=doweird&swagger=swagger&q=test"

	// Create a fake HTTP request with form data
	req, err := http.NewRequest("GET", "/submit?"+formData, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	type User struct {
		Email     string
		Name      string
		Age       int
		IsMarried bool
		Lights    bool
		Cash      float64
		Swagger   []string
		Q         string
	}

	schema := z.Struct(z.Shape{
		"email":     z.String().Email(),
		"name":      z.String().Min(3).Max(10),
		"age":       z.Int().GT(18),
		"isMarried": z.Bool().True(),
		"lights":    z.Bool().True(),
		"cash":      z.Float64().GT(10.0),
		"swagger": z.Slice[string](
			z.String().Min(1)).Min(2),
		"q": z.String().Required(),
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
	assert.Equal(t, "test", u.Q)
	assert.Empty(t, errs)
}

func TestRequestParamsOnDeleteMethodWithJsonContentType(t *testing.T) {
	// Create a fake HTTP request with JSON data
	jsonData := `{
		"userEmail": "john@doe.com",
		"userName": "JohnDoe",
		"userAge": 30,
		"userMarried": true,
		"userLights": true,
		"userMoney": 10.5,
		"userStyles": ["doweird", "swagger"],
		"userQuery": "test"
	}`

	req, err := http.NewRequest("DELETE", "/submit", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	type User struct {
		Email     string   `json:"userEmail"`
		Name      string   `json:"userName"`
		Age       int      `json:"userAge"`
		IsMarried bool     `json:"userMarried"`
		Lights    bool     `json:"userLights"`
		Cash      float64  `json:"userMoney"`
		Swagger   []string `json:"userStyles"`
		Q         string   `json:"userQuery"`
	}

	schema := z.Struct(z.Shape{
		"email":     z.String().Email(),
		"name":      z.String().Min(3).Max(10),
		"age":       z.Int().GT(18),
		"isMarried": z.Bool().True(),
		"lights":    z.Bool().True(),
		"cash":      z.Float64().GT(10.0),
		"swagger": z.Slice[string](
			z.String().Min(1)).Min(2),
		"q": z.String().Required(),
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
	assert.Equal(t, "test", u.Q)
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

func TestParseJsonWithComplexContentType(t *testing.T) {
	jsonData := `{"name":"John","age":30}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

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
	assert.Equal(t, zconst.IssueCodeInvalidJSON, err.Code)
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
	schema := z.Struct(z.Shape{
		"name": z.String().Required(),
		"nested1": z.Struct(z.Shape{
			"name": z.String().Required(),
			"nested3": z.Ptr(z.Struct(z.Shape{
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
	schema := z.Ptr(z.Struct(z.Shape{
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

	dp := form(data, &formTag)

	assert.IsType(t, urlDataProvider{}, dp)
	assert.Equal(t, data, dp.(urlDataProvider).Data)
}

package zog

import (
	"net/http"
	"strings"
	"testing"

	p "github.com/Oudwins/zog/primitives"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	formData := "name=JohnDoe&email=john@doe.com&age=30&isMarried=true&lights=on&cash=10.5&swagger=doweird"

	// Create a fake HTTP request with form data
	req, err := http.NewRequest("POST", "/submit?foo=bar&bar=foo&foo=baz", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	type User struct {
		Email     string  `form:"email"`
		Name      string  `form:"name"`
		Age       int     `form:"age"`
		IsMarried bool    `form:"isMarried"`
		Lights    bool    `form:"lights"`
		Cash      float64 `form:"cash"`
		Swagger   string  `form:"swagger"`
	}
	schema := Schema{
		"email":     String().Email(),
		"name":      String().Min(3).Max(10),
		"age":       Int().GT(18),
		"isMarried": Bool().True(),
		"lights":    Bool().True(),
		"cash":      Float().GT(10.0),
		"swagger": String().Refine("swagger", "should be doweird", func(rule p.Rule) bool {
			return rule.FieldValue.(string) == "doweird"
		}),
	}
	u := User{}

	errs, ok := Request(req, &u, schema)

	assert.Equal(t, "john@doe.com", u.Email)
	assert.Equal(t, "JohnDoe", u.Name)
	assert.Equal(t, 30, u.Age)
	assert.True(t, u.IsMarried)
	assert.True(t, u.Lights)
	assert.Equal(t, 10.5, u.Cash)
	assert.Equal(t, u.Swagger, "doweird")
	assert.Empty(t, errs)
	assert.True(t, ok)
}

func TestRequestParams(t *testing.T) {
	formData := "name=JohnDoe&email=john@doe.com&age=30&age=20&isMarried=true&lights=on&cash=10.5&swagger=doweird&swagger=swagger"

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

	schema := Schema{
		"email":     String().Email(),
		"name":      String().Min(3).Max(10),
		"age":       Int().GT(18),
		"isMarried": Bool().True(),
		"lights":    Bool().True(),
		"cash":      Float().GT(10.0),
		"swagger": Slice(
			String().Min(1)).Min(2),
	}
	u := User{}

	errs, ok := RequestParams(req, &u, schema)

	assert.Equal(t, "john@doe.com", u.Email)
	assert.Equal(t, "JohnDoe", u.Name)
	assert.Equal(t, 30, u.Age)
	assert.True(t, u.IsMarried)
	assert.True(t, u.Lights)
	assert.Equal(t, 10.5, u.Cash)
	assert.Equal(t, u.Swagger, []string{"doweird", "swagger"})
	assert.Empty(t, errs)
	assert.True(t, ok)
}

func TestStringURL(t *testing.T) {
	type Foo struct {
		Url string
	}
	foo := Foo{
		Url: "not an url",
	}
	schema := Schema{
		"url": String().URL(),
	}
	errors, ok := Validate(foo, schema)
	assert.False(t, ok)
	assert.Len(t, errors["url"], 1)

	foo.Url = "https://www.user.com"
	errors, ok = Validate(foo, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
}

func TestStringIn(t *testing.T) {
	type Foo struct {
		Currency string
	}
	foo := Foo{"eur"}
	schema := Schema{
		"currency": String().In([]string{"eur", "usd", "chz"}),
	}
	errors, ok := Validate(foo, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
	foo = Foo{"foo"}
	errors, ok = Validate(foo, schema)
	assert.False(t, ok)
	assert.Len(t, errors["currency"], 1)
}

func TestValidate(t *testing.T) {
	type User struct {
		Email    string
		Username string
	}
	schema := Schema{
		"email": String().Email(),
		// Test both lower and uppercase
		"username": String().Min(3).Max(10),
	}
	user := User{
		Email:    "foo@bar.com",
		Username: "pedropedro",
	}
	errors, ok := Validate(user, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
	assert.Empty(t, errors)
}

func TestOptional(t *testing.T) {
	type User struct {
		Email    string
		Username string
	}
	schema := Schema{
		"email":    String().Email(),
		"username": String().Min(3).Max(10).Optional(),
	}
	user := User{
		Email: "pedro@gmail.com",
	}

	errors, ok := Validate(user, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
	assert.Empty(t, errors)
}

func TestEmpty(t *testing.T) {
	type User struct {
		Email    string
		Username string
	}
	schema := Schema{
		"email":    String(),
		"username": String(),
	}
	user := User{
		Email:    "",
		Username: "",
	}

	errors, ok := Validate(user, schema)
	assert.True(t, ok)
	assert.Empty(t, errors)
	assert.Empty(t, errors)
}

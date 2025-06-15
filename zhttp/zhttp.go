package zhttp

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/parsers/zjson"
	"github.com/Oudwins/zog/zconst"
)

type ParserFunc = func(r *http.Request) p.DpFactory

var (
	formTag    string = "form"
	queryParam string = "query"
)

var Config = struct {
	Parsers struct {
		JSON  ParserFunc
		Form  ParserFunc
		Query ParserFunc
	}
}{
	Parsers: struct {
		JSON  ParserFunc
		Form  ParserFunc
		Query ParserFunc
	}{
		JSON: func(r *http.Request) p.DpFactory {
			return zjson.Decode(r.Body)
		},
		Form: func(r *http.Request) p.DpFactory {
			return func() (p.DataProvider, *p.ZogIssue) {
				err := r.ParseForm()
				if err != nil {
					return nil, &p.ZogIssue{Code: zconst.IssueCodeZHTTPInvalidForm, Err: err}
				}
				return form(r.Form, &formTag), nil
			}
		},
		Query: func(r *http.Request) p.DpFactory {
			return func() (p.DataProvider, *p.ZogIssue) {
				// This handles generic GET request from browser. We treat it as url.Values
				return form(r.URL.Query(), &queryParam), nil
			}
		},
	},
}

type urlDataProvider struct {
	Data url.Values
	tag  *string
}

var _ p.DataProvider = urlDataProvider{}

func (u urlDataProvider) Get(key string) any {
	if !u.Data.Has(key) {
		return nil
	}

	// if query param ends with [] its always a slice
	if len(key) > 2 && key[len(key)-2:] == "[]" {
		return u.Data[key]
	}

	if len(u.Data[key]) > 1 {
		return u.Data[key]
	} else {
		return u.Data.Get(key)
	}
}

func (u urlDataProvider) GetByField(field reflect.StructField, fallback string) (any, string) {
	key := p.GetKeyFromField(field, fallback, u.tag)
	return u.Get(key), key
}

func (u urlDataProvider) GetNestedProvider(key string) p.DataProvider {
	return u
}
func (u urlDataProvider) GetUnderlying() any {
	return u.Data
}

// Parses JSON, Form & Query data from request based on Content-Type header
// Usage:
// schema.Parse(zhttp.Request(r), &dest)
// WARNING: FOR JSON PARSING DOES NOT SUPPORT JSON ARRAYS OR PRIMITIVES
func Request(r *http.Request) p.DpFactory {
	switch r.Method {
	case "GET":
		return Config.Parsers.Query(r)
	case "HEAD":
		return Config.Parsers.Query(r)
	default:
		// Content-Type follows this format: Content-Type: <media-type> [; parameter=value]
		typ, _, _ := strings.Cut(r.Header.Get("Content-Type"), ";")
		switch typ {
		case "application/json":
			return Config.Parsers.JSON(r)
		case "application/x-www-form-urlencoded":
			return Config.Parsers.Form(r)
		default:
			return Config.Parsers.Query(r)
		}
	}
}

func form(data url.Values, tag *string) p.DataProvider {
	return urlDataProvider{Data: data, tag: tag}
}

// func params(data url.Values) p.DataProvider {
// 	return form(data)
// }

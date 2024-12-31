package zhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type ParserFunc = func(r *http.Request) (p.DataProvider, p.ZogError)

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
		JSON: parseJson,
		Form: func(r *http.Request) (p.DataProvider, p.ZogError) {
			err := r.ParseForm()
			if err != nil {
				return nil, &p.ZogErr{C: zconst.ErrCodeZHTTPInvalidForm, Err: err}
			}
			return form(r.Form), nil
		},
		Query: func(r *http.Request) (p.DataProvider, p.ZogError) {
			// This handles generic GET request from browser. We treat it as url.Values
			return form(r.URL.Query()), nil
		},
	},
}

type urlDataProvider struct {
	Data url.Values
}

var _ p.DataProvider = urlDataProvider{}

func (u urlDataProvider) Get(key string) any {
	// if query param ends with [] its always a slice
	if key[len(key)-2:] == "[]" {
		return u.Data[key]
	}

	if len(u.Data[key]) > 1 {
		return u.Data[key]
	} else {
		return u.Data.Get(key)
	}
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
func Request(r *http.Request) p.DpFactory {
	return func() (p.DataProvider, p.ZogError) {
		switch r.Header.Get("Content-Type") {
		case "application/json":
			return Config.Parsers.JSON(r)
		case "application/x-www-form-urlencoded":
			return Config.Parsers.Form(r)
		default:
			return Config.Parsers.Query(r)
		}
	}
}

// Parses JSON data from request body. Does not support json arrays or primitives
/*
- "null" -> nil -> Not accepted by zhttp -> errs["$root"]-> required error
- "{}" -> okay -> map[]{}
- "" -> parsing error -> errs["$root"]-> parsing error
- "1213" -> zhttp -> plain value
  - struct schema -> hey this valid input
  - "string is not an object"
*/
func parseJson(r *http.Request) (p.DataProvider, p.ZogError) {
	var m map[string]any
	decod := json.NewDecoder(r.Body)
	err := decod.Decode(&m)
	if err != nil {
		return nil, &p.ZogErr{C: zconst.ErrCodeInvalidJSON, Err: err}
	}
	if m == nil {
		return nil, &p.ZogErr{C: zconst.ErrCodeInvalidJSON, Err: errors.New("nill json body")}
	}
	return p.NewMapDataProvider(m), nil
}

func form(data url.Values) p.DataProvider {
	return urlDataProvider{Data: data}
}

// func params(data url.Values) p.DataProvider {
// 	return form(data)
// }

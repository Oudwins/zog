package zhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

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
	return func() (p.DataProvider, *p.ZogErr) {
		switch r.Header.Get("Content-Type") {
		case "application/json":
			return parseJson(r.Body)
		case "application/x-www-form-urlencoded":
			err := r.ParseForm()
			if err != nil {
				return nil, &p.ZogErr{C: zconst.ErrCodeZHTTPInvalidForm, Err: err}
			}
			return form(r.Form), nil
		default:
			// This handles generic GET request from browser. We treat it as url.Values
			err := r.ParseForm()
			if err != nil {
				return nil, &p.ZogErr{C: zconst.ErrCodeZHTTPInvalidQuery, Err: err}
			}
			return form(r.Form), nil
		}
	}
}

func parseJson(data io.ReadCloser) (p.DataProvider, *p.ZogErr) {
	var m map[string]any
	decod := json.NewDecoder(data)
	err := decod.Decode(&data)
	if err != nil {
		return nil, &p.ZogErr{C: zconst.ErrCodeZHTTPInvalidJSON, Err: err}
	}
	return p.NewMapDataProvider(m), nil
}

func form(data url.Values) p.DataProvider {
	return urlDataProvider{Data: data}
}

// func params(data url.Values) p.DataProvider {
// 	return form(data)
// }

// DEPRECATED: DO NOT USE WILL BE REMOVED
func NewRequestDataProvider(r *http.Request) (urlDataProvider, error) {
	err := r.ParseForm()
	if err != nil {
		return urlDataProvider{}, err
	}
	return urlDataProvider{Data: r.Form}, nil
}

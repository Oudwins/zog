package zhttp

import (
	"encoding/json"
	"net/http"
	"net/url"

	p "github.com/Oudwins/zog/internals"
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

// Only supports form data & query parms
func NewRequestDataProvider(r *http.Request) (urlDataProvider, error) {
	err := r.ParseForm()
	if err != nil {
		return urlDataProvider{}, err
	}
	return urlDataProvider{Data: r.Form}, nil
}

func NewJsonDataProvider(r *http.Request) (p.DataProvider, error) {
	var data map[string]any
	decod := json.NewDecoder(r.Body)
	err := decod.Decode(&data)
	if err != nil {
		return nil, err
	}
	return p.NewMapDataProvider(data), nil
}

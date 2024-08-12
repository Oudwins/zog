package zhttp

import (
	"net/http"
	"net/url"

	p "github.com/Oudwins/zog/primitives"
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

func NewRequestDataProvider(r *http.Request) (urlDataProvider, error) {
	err := r.ParseForm()
	if err != nil {
		return urlDataProvider{}, err
	}
	return urlDataProvider{Data: r.Form}, nil
}

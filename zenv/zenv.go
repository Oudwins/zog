package zenv

import (
	"os"
	"reflect"
	"strings"

	p "github.com/Oudwins/zog/internals"
)

var _ p.DataProvider = &envDataProvider{}

var (
	envTag string = "env"
)

type envDataProvider struct {
}

func (e *envDataProvider) Get(key string) any {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return nil
	}
	return val
}

func (e *envDataProvider) GetByField(field reflect.StructField, fallback string) (any, string) {
	key := p.GetKeyFromField(field, fallback, &envTag)
	return e.Get(key), key
}

func (e *envDataProvider) GetNestedProvider(key string) p.DataProvider {
	return e
}

func NewDataProvider() *envDataProvider {
	return &envDataProvider{}
}

func (e *envDataProvider) GetUnderlying() any {
	return nil
}

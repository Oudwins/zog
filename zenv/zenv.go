package zenv

import (
	"os"
	"strings"

	p "github.com/Oudwins/zog/internals"
)

var _ p.DataProvider = &envDataProvider{}

type envDataProvider struct {
}

func (e *envDataProvider) Get(key string) any {
	return strings.TrimSpace(os.Getenv(key))
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

package zenv

import (
	"os"

	p "github.com/Oudwins/zog/internals"
)

type envDataProvider struct {
}

func (e *envDataProvider) Get(key string) any {
	return os.Getenv(key)
}

func (e *envDataProvider) GetNestedProvider(key string) p.DataProvider {
	return e
}

func NewDataProvider() *envDataProvider {
	return &envDataProvider{}
}

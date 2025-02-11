package internals

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/zconst"
)

type DpFactory = func() (DataProvider, ZogIssue)

// This is used for parsing structs & maps
type DataProvider interface {
	Get(key string) any
	GetNestedProvider(key string) DataProvider
	GetUnderlying() any // returns the underlying value the dp is wrapping
}

// checks that we implement the interface
var _ DataProvider = &MapDataProvider[string]{}

type MapDataProvider[T any] struct {
	M map[string]T
}

func (m *MapDataProvider[T]) Get(key string) any {
	return any(m.M[key])
}

func (m *MapDataProvider[T]) GetNestedProvider(key string) DataProvider {
	dataProvider, _ := TryNewAnyDataProvider(m.M[key])
	return dataProvider
}

func (m *MapDataProvider[T]) GetUnderlying() any {
	return m.M
}

func NewMapDataProvider[T any](m map[string]T) DataProvider {
	if m == nil {
		return &EmptyDataProvider{}
	}
	return &MapDataProvider[T]{
		M: m,
	}
}

type EmptyDataProvider struct {
	Underlying any
}

func (e *EmptyDataProvider) Get(key string) any {
	return nil
}

func (e *EmptyDataProvider) GetNestedProvider(key string) DataProvider {
	return e
}

func (e *EmptyDataProvider) GetUnderlying() any {
	return e.Underlying
}

func TryNewAnyDataProvider(val any) (DataProvider, ZogError) {
	dp, ok := val.(DataProvider)
	if ok {
		return dp, nil
	}
	factory, ok := val.(DpFactory)
	if ok {
		return factory()
	}
	x := reflect.ValueOf(val)
	switch x.Kind() {
	case reflect.Map:
		keyTyp := x.Type().Key()

		if keyTyp.Kind() != reflect.String {
			return &EmptyDataProvider{Underlying: val}, &ZogErr{
				C:   zconst.IssueCodeCoerce,
				Err: fmt.Errorf("could not convert map[%s]any to a data provider", keyTyp.String()),
			}
		}

		valTyp := x.Type().Elem()

		switch valTyp.Kind() { // TODO: add more types
		case reflect.String:
			return NewMapDataProvider(x.Interface().(map[string]string)), nil
		case reflect.Int:
			return NewMapDataProvider(x.Interface().(map[string]int)), nil
		case reflect.Float64:
			return NewMapDataProvider(x.Interface().(map[string]float64)), nil
		case reflect.Bool:
			return NewMapDataProvider(x.Interface().(map[string]bool)), nil
		case reflect.Interface:
			return NewMapDataProvider(x.Interface().(map[string]any)), nil
		default:
			return &EmptyDataProvider{Underlying: val}, &ZogErr{
				C:   zconst.IssueCodeCoerce,
				Err: fmt.Errorf("could not convert map[string]%s to a data provider", valTyp.String()),
			}
		}

	case reflect.Pointer:
		if x.IsNil() {
			return &EmptyDataProvider{}, &ZogErr{
				C:   zconst.IssueCodeCoerce,
				Err: errors.New("could not convert pointer to a data provider"),
			}
		}
		return TryNewAnyDataProvider(x.Elem().Interface())

	default:
		return &EmptyDataProvider{Underlying: val}, &ZogErr{
			C:   zconst.IssueCodeCoerce,
			Err: fmt.Errorf("could not convert type %s to a data provider. unsupported type", x.Kind().String()),
		}
	}
}

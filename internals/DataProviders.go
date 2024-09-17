package internals

import "reflect"

// This is used for parsing structs & maps
type DataProvider interface {
	Get(key string) any
	GetNestedProvider(key string) DataProvider
}

type MapDataProvider[T any] struct {
	M map[string]T
}

func (m *MapDataProvider[T]) Get(key string) any {
	return any(m.M[key])
}

func (m *MapDataProvider[T]) GetNestedProvider(key string) DataProvider {
	return NewAnyDataProvider(m.M[key])
}

// checks that we implement the interface
var _ DataProvider = &MapDataProvider[string]{}

type EmptyDataProvider struct{}

func (e *EmptyDataProvider) Get(key string) any {
	return nil
}

func (e *EmptyDataProvider) GetNestedProvider(key string) DataProvider {
	return e
}

func NewMapDataProvider[T any](m map[string]T) DataProvider {
	if m == nil {
		return &EmptyDataProvider{}
	}
	return &MapDataProvider[T]{
		M: m,
	}
}

// Tries to create a map data provider from any value if it cannot it will return an empty data provider (which will always return nil)
func NewAnyDataProvider(val any) DataProvider {
	dataProvider, _ := TryNewAnyDataProvider(val)
	return dataProvider
}

func TryNewAnyDataProvider(val any) (DataProvider, bool) {
	x := reflect.ValueOf(val)

	switch x.Kind() {
	case reflect.Map:
		keyTyp := x.Type().Key()

		if keyTyp.Kind() != reflect.String {
			return &EmptyDataProvider{}, false
		}

		valTyp := x.Type().Elem()

		switch valTyp.Kind() { // TODO: add more types
		case reflect.String:
			return NewMapDataProvider(x.Interface().(map[string]string)), true
		case reflect.Int:
			return NewMapDataProvider(x.Interface().(map[string]int)), true
		case reflect.Float64:
			return NewMapDataProvider(x.Interface().(map[string]float64)), true
		case reflect.Bool:
			return NewMapDataProvider(x.Interface().(map[string]bool)), true
		case reflect.Interface:
			return NewMapDataProvider(x.Interface().(map[string]any)), true
		default:
			return &EmptyDataProvider{}, false
		}

	case reflect.Pointer:
		if x.IsNil() {
			return &EmptyDataProvider{}, false
		}
		return TryNewAnyDataProvider(x.Elem().Interface())

	default:
		return &EmptyDataProvider{}, false
	}
}

package internals

import (
	"fmt"
	"reflect"
	"strings"

	zconst "github.com/Oudwins/zog/zconst"
)

func GetKeyFromField(field reflect.StructField, fallback string, tag *string) string {
	if tag != nil {
		fieldTag, ok := field.Tag.Lookup(*tag)
		if ok {
			name, _, _ := strings.Cut(fieldTag, ",")
			if name != "" {
				return name
			}
		}
	}
	fieldTag, ok := field.Tag.Lookup(zconst.ZogTag)
	if ok {
		return fieldTag
	}
	return fallback
}

type DpFactory = func() (DataProvider, *ZogIssue)

// This is used for parsing structs & maps
type DataProvider interface {
	Get(key string) any
	GetByField(field reflect.StructField, fallback string) (any, string)
	GetNestedProvider(key string) DataProvider
	GetUnderlying() any // returns the underlying value the dp is wrapping
}

// checks that we implement the interface
var _ DataProvider = &MapDataProvider[string]{}
var _ DataProvider = &StructDataProvider{}

type StructDataProvider struct {
	value reflect.Value
	tag   *string
}

func (s *StructDataProvider) Get(key string) any {
	field := s.value.FieldByName(key)
	if !field.IsValid() {
		return nil
	}
	return field.Interface()
}

func (s *StructDataProvider) GetByField(field reflect.StructField, fallback string) (any, string) {
	key := GetKeyFromField(field, fallback, s.tag)
	return s.Get(key), key
}

func (s *StructDataProvider) GetNestedProvider(key string) DataProvider {
	field := s.value.FieldByName(key)
	if !field.IsValid() {
		return nil
	}
	dataProvider, _ := TryNewAnyDataProvider(field.Interface())
	return dataProvider
}

func (s *StructDataProvider) GetUnderlying() any {
	return s.value.Interface()
}

type MapDataProvider[T any] struct {
	M   map[string]T
	tag *string
}

func (m *MapDataProvider[T]) Get(key string) any {
	v, ok := m.M[key]
	if !ok {
		return nil
	}
	// Present with nil emits the sentinel so pointer schemas can distinguish it
	// from an absent key.
	if any(v) == nil {
		return ExplicitNullMarker()
	}
	// A typed nil pointer inside an interface like (*string) is not == nil but is
	// still semantically an explicit null. Other nillable kinds are excluded, a
	// nil slice is accepted as empty input by non-pointer schemas today.
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Pointer && rv.IsNil() {
		return ExplicitNullMarker()
	}
	return v
}

// returns value + key used
func (m *MapDataProvider[T]) GetByField(field reflect.StructField, fallback string) (any, string) {
	key := GetKeyFromField(field, fallback, m.tag)
	return m.Get(key), key
}

func (m *MapDataProvider[T]) GetNestedProvider(key string) DataProvider {
	dataProvider, _ := TryNewAnyDataProvider(m.M[key])
	return dataProvider
}

func (m *MapDataProvider[T]) GetUnderlying() any {
	return m.M
}

func NewMapDataProvider[T any](m map[string]T, tag *string) DataProvider {
	if len(m) == 0 {
		return &EmptyDataProvider{tag: tag}
	}
	return &MapDataProvider[T]{
		M:   m,
		tag: tag,
	}
}

func NewSafeMapDataProvider[T any](m map[string]T) DataProvider {
	if len(m) == 0 {
		return &EmptyDataProvider{}
	}
	return NewMapDataProvider(m, nil)
}

type EmptyDataProvider struct {
	Underlying any
	tag        *string
}

func (e *EmptyDataProvider) Get(key string) any {
	return nil
}

func (e *EmptyDataProvider) GetByField(field reflect.StructField, fallback string) (any, string) {
	return nil, GetKeyFromField(field, fallback, e.tag)
}

func (e *EmptyDataProvider) GetNestedProvider(key string) DataProvider {
	return nil
}

func (e *EmptyDataProvider) GetUnderlying() any {
	return e.Underlying
}

func TryNewAnyDataProvider(val any) (DataProvider, error) {
	dp, ok := val.(DataProvider)
	if ok {
		return dp, nil
	}
	// Here we treat the sentinel and a nil (absence) as the same because we are
	// at the top level, for individual fields we do a different behavior.
	if IsExplicitNull(val) || val == nil {
		return &EmptyDataProvider{Underlying: val}, nil
	}
	x := reflect.ValueOf(val)
	switch x.Kind() {
	case reflect.Map:
		keyTyp := x.Type().Key()

		if keyTyp.Kind() != reflect.String {
			return &EmptyDataProvider{Underlying: val}, fmt.Errorf("could not convert map[%s]any to a data provider", keyTyp.String())
		}

		valTyp := x.Type().Elem()

		switch valTyp.Kind() { // TODO: add more types
		case reflect.String:
			return NewSafeMapDataProvider(x.Interface().(map[string]string)), nil
		case reflect.Int:
			return NewSafeMapDataProvider(x.Interface().(map[string]int)), nil
		case reflect.Float64:
			return NewSafeMapDataProvider(x.Interface().(map[string]float64)), nil
		case reflect.Bool:
			return NewSafeMapDataProvider(x.Interface().(map[string]bool)), nil
		case reflect.Interface:
			return NewSafeMapDataProvider(x.Interface().(map[string]any)), nil
		default:
			return &EmptyDataProvider{Underlying: val}, fmt.Errorf("could not convert map[string]%s to a data provider", valTyp.String())
		}

	case reflect.Struct:
		return &StructDataProvider{value: x, tag: nil}, nil

	case reflect.Pointer:
		if x.IsNil() {
			return &EmptyDataProvider{}, nil
		}
		return TryNewAnyDataProvider(x.Elem().Interface())

	default:
		return &EmptyDataProvider{Underlying: val}, fmt.Errorf("could not convert type %s to a data provider. unsupported type", x.Kind().String())
	}
}

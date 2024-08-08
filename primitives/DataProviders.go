package primitives

// This is used for parsing structs & maps
type DataProvider interface {
	Get(key string) any
}

type MapDataProvider[T any] struct {
	M map[string]T
}

func (m *MapDataProvider[T]) Get(key string) any {
	return any(m.M[key])
}

// checks that we implement the interface
var _ DataProvider = &MapDataProvider[string]{}

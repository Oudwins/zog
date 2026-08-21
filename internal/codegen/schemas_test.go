package codegen

import (
	"reflect"
	"testing"
)

func TestZogSchemas(t *testing.T) {
	schemas, err := ZogSchemas("../..")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"AnySchema",
		"BoolSchema[T]",
		"BoxedSchema[B, T]",
		"CustomSchema",
		"Custom[T]",
		"MapSchema[K, V]",
		"NumberSchema[T]",
		"PointerSchema",
		"PreprocessSchema[F, T]",
		"SliceSchema",
		"StringSchema[T]",
		"StructSchema",
		"TimeSchema",
		"lazySchema",
	}
	if !reflect.DeepEqual(schemas, want) {
		t.Fatalf("ZogSchemas() = %v, want %v", schemas, want)
	}
}

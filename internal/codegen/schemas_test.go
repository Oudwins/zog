package codegen

import (
	"os"
	"path/filepath"
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

func TestZogSchemasChecksMethodSets(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/schemas\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	source := `package zog

type ZogSchema interface { schema(int) string }

type Valid struct{}
func (*Valid) schema(int) string { return "" }

type WrongSignature struct{}
func (*WrongSignature) schema(string) string { return "" }

type Promoted struct{ ZogSchema }
`
	if err := os.WriteFile(filepath.Join(dir, "schemas.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	schemas, err := ZogSchemas(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Promoted", "Valid"}
	if !reflect.DeepEqual(schemas, want) {
		t.Fatalf("ZogSchemas() = %v, want %v", schemas, want)
	}
}

package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zconst"
)

func main() {
	schema := z.String().Min(1)
	json, err := ToJsonSchema(schema)
	_, _ = json, err
}
func ToJsonSchema(zogSchema z.ZogSchema) (string, error) {
	// schema := map[string]any{}
	Walk(zogSchema, []string{})
	return "", nil
}

func Walk(zogSchema any, path []string) {
	refVal := reflect.ValueOf(zogSchema)
	// typ, _ := GetSchemaType(refVal)
	// fmt.Println("type found: ", typ)
	tests, _ := GetTests(refVal)
	fmt.Println("found tests", tests)
}

// Get the tests
func GetTests(s reflect.Value) ([]z.Test, error) {
	val := s
	for val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	testsValue := val.FieldByIndex([]int{1})
	tests, ok := testsValue.Interface().([]z.Test)
	if !ok {
		return nil, errors.New("Could not get tests from schema")
	}
	return tests, nil
}

// Get the schema type
func GetSchemaType(s reflect.Value) (string, error) {
	return getSchemaTypeViaTypeName(s)
}

const (
	SchemaIdString = "zog.StringSchema["
)

func getSchemaTypeViaTypeName(s reflect.Value) (zconst.ZogType, error) {
	stringType := s.Type().String()

	if strings.Contains(stringType, SchemaIdString) {
		return zconst.TypeString, nil
	}

	panic("unsupported schema")
	// return "", errors.New("Invalid schema")
}

func getSchemaTypeViaInterface(s any) (zconst.ZogType, error) {
	// TODO alternative where we have an interface for each and we cast the any value to it then get the values we need from it
	return "", nil
}

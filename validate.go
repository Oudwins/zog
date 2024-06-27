package zog

import (
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"strconv"
	"unicode"
)

// Errors is a map holding all the possible errors that may
// occur during validation.
type Errors map[string][]string

// Any return true if there is any error.
func (e Errors) Any() bool {
	return len(e) > 0
}

// Add adds an error for a specific field
func (e Errors) Add(field string, msg string) {
	if _, ok := e[field]; !ok {
		e[field] = []string{}
	}
	e[field] = append(e[field], msg)
}

// Get returns all the errors for the given field.
func (e Errors) Get(field string) []string {
	return e[field]
}

// the interface of any set of rules for a field. Eg: String().Min(5) -> this is a fieldValidator
type fieldValidator interface {
	Validate(val any) (errors []string, ok bool)
}

// Schema represents a validation schema.
type Schema map[string]fieldValidator

// Validate validates data based on the given Schema.
func (s Schema) Validate(data any) (Errors, bool) {
	errors := Errors{}
	return validateSchema(data, s, errors)
}
func Validate(data any, fields Schema) (Errors, bool) {
	errors := Errors{}
	return validateSchema(data, fields, errors)
}

func validateSchema(data any, schema Schema, errors Errors) (Errors, bool) {
	globalOk := true
	ok := true
	var fieldErrs []string

	for fieldName, validator := range schema {
		fieldName = string(unicode.ToUpper(rune(fieldName[0]))) + fieldName[1:]
		fieldValue := getFieldValueByName(data, fieldName)
		fieldName = string(unicode.ToLower([]rune(fieldName)[0])) + fieldName[1:]
		fieldErrs, ok = validator.Validate(fieldValue)
		if !ok {
			errors[fieldName] = fieldErrs
			globalOk = false
		}

	}

	return errors, globalOk
}

// Merge merges the two given schemas returning a new Schema. In case of clashing second will take priority
func Merge(first, second Schema, rest ...Schema) Schema {
	newSchema := Schema{}
	maps.Copy(newSchema, first)
	maps.Copy(newSchema, second)

	for _, s := range rest {
		maps.Copy(newSchema, s)
	}

	return newSchema
}

// ! PARSE REQUESTS
// Request parses an http.Request into data and validates it based
// on the given schema.
func Request(r *http.Request, data any, schema Schema) (Errors, bool) {
	errors := Errors{}
	if err := parseRequest(r, data); err != nil {
		errors["_error"] = []string{err.Error()}
	}
	return validateSchema(data, schema, errors)
}

// TODO -> Parse requestQueryParams
func RequestParams(r *http.Request, data any, schema Schema) (Errors, bool) {
	errors := Errors{}
	if err := parseRequestParams(r, data); err != nil {
		errors["_error"] = []string{err.Error()}
	}
	return validateSchema(data, schema, errors)
}

func parseRequestParams(r *http.Request, v any) error {

	params := r.URL.Query()
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		paramTag := field.Tag.Get("param")
		param := params[paramTag]

		if len(param) == 0 {
			continue
		}

		fieldVal := val.Field(i)
		t := fieldVal.Kind()
		switch t {
		case reflect.Slice:
			for idx, v := range param {
				if idx < fieldVal.Len() {
					fieldVal.Index(idx).Set(reflect.ValueOf(v))
				} else {
					newElem := reflect.Append(fieldVal, reflect.ValueOf(v))
					fieldVal.Set(newElem)
				}
			}
		default:
			if err := parsePrimitive(&t, &fieldVal, param[0]); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseRequest(r *http.Request, v any) error {
	contentType := r.Header.Get("Content-Type")
	// TODO support more content types
	if contentType == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %v", err)
		}
		val := reflect.ValueOf(v).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			formTag := field.Tag.Get("form")
			formValue := r.FormValue(formTag)

			if formValue == "" {
				continue
			}

			fieldVal := val.Field(i)
			typ := fieldVal.Kind()
			if err := parsePrimitive(&typ, &fieldVal, formValue); err != nil {
				return err
			}
		}

	}
	return nil
}

func parsePrimitive(typ *reflect.Kind, refObj *reflect.Value, value string) error {
	switch *typ {
	case reflect.Bool:
		// There are cases where frontend libraries use "on" as the bool value
		// think about toggles. Hence, let's try this first.
		if value == "on" {
			refObj.SetBool(true)
		} else if value == "off" {
			refObj.SetBool(false)
			return nil
		} else {
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("failed to parse bool: %v", err)
			}
			refObj.SetBool(boolVal)
		}

	case reflect.String:
		refObj.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to parse int: %v", err)
		}
		refObj.SetInt(int64(intVal))
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float: %v", err)
		}
		refObj.SetFloat(floatVal)
	default:
		return fmt.Errorf("unsupported kind %s", refObj.Kind())
	}

	return nil
}

func getFieldValueByName(v any, name string) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil
	}
	return fieldVal.Interface()
}

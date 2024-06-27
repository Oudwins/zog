package primitives

import "reflect"

func IsZeroValue(x any) bool {
	if x == nil {
		return true
	}

	v := reflect.ValueOf(x)
	if !v.IsValid() {
		return true
	}

	// Check if the value is the zero value for its type
	zeroValue := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zeroValue.Interface())
}

func GenericValidator(fieldValue any, rules []Rule, isOptional bool) ([]string, bool) {

	var errors []string = nil
	ok := true

	// if its optional and the value is zero we can skip the validation
	if isOptional && IsZeroValue(fieldValue) {
		return errors, ok
	}

	for _, set := range rules {

		set.FieldValue = fieldValue
		if !set.ValidateFunc(set) {
			ok = false
			msg := set.ErrorMessage
			if errors == nil {
				errors = []string{}
			}
			errors = append(errors, msg)
		}
	}

	return errors, ok
}

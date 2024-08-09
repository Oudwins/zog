package notes

// Should we have an initialize function on the data provider that returns an error?

// This should parse both form and url params
// func Request(r *http.Request, schema z.StructParser, destPtr any) (error, p.ZogSchemaErrors) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		return err, nil
// 	}
// 	fieldNames := getFieldNamesFromStructPtr(destPtr)

// 	m := make(map[string]any, len(fieldNames))
// 	for _, fieldName := range fieldNames {
// 		if len(r.Form[fieldName]) > 1 {
// 			m[fieldName] = r.Form[fieldName]
// 		} else {
// 			m[fieldName] = r.FormValue(fieldName)
// 		}
// 	}

// 	return nil, schema.Parse(z.NewMapDataProvider(m), destPtr)
// }

// func getFieldNamesFromStructPtr(v any) []string {
// 	val := reflect.ValueOf(v).Elem()
// 	fieldNames := make([]string, val.NumField())
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Type().Field(i)
// 		paramTag := field.Tag.Get(p.ZogTag)
// 		if paramTag == "" {
// 			fieldNames[i] = field.Name
// 		} else {
// 			fieldNames[i] = paramTag
// 		}
// 	}
// 	return fieldNames
// }

// func getStructFieldPointerByName(v any, name string) (any, error) {
// 	val := reflect.ValueOf(v)
// 	if val.Kind() != reflect.Ptr {
// 		return nil, fmt.Errorf("expected dest to be a pointer, got %s", val.Kind())
// 	}
// 	val = val.Elem()
// 	if val.Kind() != reflect.Struct {
// 		return nil, fmt.Errorf("expected dest to point to a struct, got %s", val.Kind())
// 	}

// 	fieldVal := val.FieldByName(name)
// 	if !fieldVal.IsValid() {
// 		return nil, fmt.Errorf("field %s not found", name)
// 	}
// 	// TODO HERE WE PROBABLY NEED TO CREATE structs, maps, slices, etc if they don't exist and that is the kind of field
// 	// if its an optional value
// 	fmt.Println("Field vals")
// 	fmt.Println(fieldVal.Kind(), fieldVal.Type())

// 	// HERE WE NEED TO CHECK IF THE VALUE IS A POINTER. If it is, we should check the underlying value for nil
// 	// if its nil we need to create it based on the type

// 	return fieldVal.Addr().Interface(), nil
// }

// func RequestParams(r *http.Request, data any, schema SchemaOld) (Errors, bool) {
// 	errors := Errors{}
// 	if err := parseRequestParams(r, data); err != nil {
// 		errors["_error"] = []string{err.Error()}
// 	}
// 	return parseSchema(data, schema, errors)
// }

// func parseRequestParams(r *http.Request, v any) error {

// 	params := r.URL.Query()
// 	val := reflect.ValueOf(v).Elem()
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Type().Field(i)
// 		paramTag := field.Tag.Get("param")
// 		param := params[paramTag]

// 		if len(param) == 0 || param[0] == "" {
// 			continue
// 		}

// 		fieldVal := val.Field(i)
// 		t := fieldVal.Kind()
// 		switch t {
// 		case reflect.Slice:
// 			for idx, v := range param {
// 				if idx < fieldVal.Len() {
// 					fieldVal.Index(idx).Set(reflect.ValueOf(v))
// 				} else {
// 					newElem := reflect.Append(fieldVal, reflect.ValueOf(v))
// 					fieldVal.Set(newElem)
// 				}
// 			}
// 		default:
// 			if err := parsePrimitive(&t, &fieldVal, param[0]); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func parseRequest(r *http.Request, v any) error {
// 	contentType := r.Header.Get("Content-Type")
// 	// TODO support more content types
// 	if contentType == "application/x-www-form-urlencoded" {
// 		if err := r.ParseForm(); err != nil {
// 			return fmt.Errorf("failed to parse form: %v", err)
// 		}
// 		val := reflect.ValueOf(v).Elem()
// 		for i := 0; i < val.NumField(); i++ {
// 			field := val.Type().Field(i)
// 			formTag := field.Tag.Get("form")
// 			formValue := r.FormValue(formTag)

// 			if formValue == "" {
// 				continue
// 			}

// 			fieldVal := val.Field(i)
// 			typ := fieldVal.Kind()
// 			if err := parsePrimitive(&typ, &fieldVal, formValue); err != nil {
// 				return err
// 			}
// 		}

// 	}
// 	return nil
// }

// func parsePrimitive(typ *reflect.Kind, refObj *reflect.Value, value string) error {
// 	switch *typ {
// 	case reflect.Bool:
// 		// There are cases where frontend libraries use "on" as the bool value
// 		// think about toggles. Hence, let's try this first.
// 		if value == "on" {
// 			refObj.SetBool(true)
// 		} else if value == "off" {
// 			refObj.SetBool(false)
// 			return nil
// 		} else {
// 			boolVal, err := strconv.ParseBool(value)
// 			if err != nil {
// 				return fmt.Errorf("failed to parse bool: %v", err)
// 			}
// 			refObj.SetBool(boolVal)
// 		}

// 	case reflect.String:
// 		refObj.SetString(value)
// 	case reflect.Int:
// 		intVal, err := strconv.Atoi(value)
// 		if err != nil {
// 			return fmt.Errorf("failed to parse int: %v", err)
// 		}
// 		refObj.SetInt(int64(intVal))
// 	case reflect.Float64:
// 		floatVal, err := strconv.ParseFloat(value, 64)
// 		if err != nil {
// 			return fmt.Errorf("failed to parse float: %v", err)
// 		}
// 		refObj.SetFloat(floatVal)
// 	default:
// 		return fmt.Errorf("unsupported kind %s", refObj.Kind())
// 	}

// 	return nil
// }

// func setPrimitiveValue(obj any, newVal any, fieldName string) {

// 	val := reflect.ValueOf(obj).Elem()
// 	fieldVal := val.FieldByName(fieldName)
// 	if !fieldVal.IsValid() {
// 		return
// 	}
// 	fieldVal.Set(reflect.ValueOf(newVal))
// }

// func getFieldValueByName(v any, name string) any {
// 	val := reflect.ValueOf(v)
// 	if val.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 	}
// 	if val.Kind() != reflect.Struct {
// 		return nil
// 	}
// 	fieldVal := val.FieldByName(name)
// 	if !fieldVal.IsValid() {
// 		return nil
// 	}

// 	return fieldVal.Interface()
// }

package zog

import (
	"fmt"
	"reflect"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// Helper test schema implementation for testing EXPERIMENTAL_PUBLIC_ZOG_SCHEMA
type testExperimentalSchema struct {
	validator func(val any) bool
	dtype     zconst.ZogType
	coercer   CoercerFunc
	errorMsg  string
}

func (s *testExperimentalSchema) Process(ctx *p.SchemaCtx) {
	// Try to coerce if coercer is set
	if s.coercer != nil {
		coerced, err := s.coercer(ctx.Data)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromCoerce(err))
			return
		}
		ctx.Data = coerced
	}

	// Type assertion and assignment
	switch ptr := ctx.ValPtr.(type) {
	case *string:
		val, ok := ctx.Data.(string)
		if !ok {
			ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("expected string, got %T", ctx.Data)))
			return
		}
		*ptr = val
		if s.validator != nil && !s.validator(*ptr) {
			issue := ctx.Issue().SetMessage(s.errorMsg)
			if s.errorMsg == "" {
				issue.SetCode(zconst.IssueCodeCustom)
			}
			ctx.AddIssue(issue)
			return
		}
	case *int:
		val, ok := ctx.Data.(int)
		if !ok {
			// Try to coerce from float64 (common case)
			if floatVal, ok := ctx.Data.(float64); ok {
				val = int(floatVal)
			} else {
				ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("expected int, got %T", ctx.Data)))
				return
			}
		}
		*ptr = val
		if s.validator != nil && !s.validator(*ptr) {
			issue := ctx.Issue().SetMessage(s.errorMsg)
			if s.errorMsg == "" {
				issue.SetCode(zconst.IssueCodeCustom)
			}
			ctx.AddIssue(issue)
			return
		}
	case *testCustomType:
		var val testCustomType
		// Try direct type assertion first (fastest path)
		if v, ok := ctx.Data.(testCustomType); ok {
			val = v
		} else if v, ok := ctx.Data.(*testCustomType); ok {
			val = *v
		} else {
			// Use reflection to extract Value field
			dataVal := reflect.ValueOf(ctx.Data)
			if !dataVal.IsValid() {
				ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("expected testCustomType, got nil")))
				return
			}
			if dataVal.Kind() == reflect.Struct {
				valueField := dataVal.FieldByName("Value")
				if valueField.IsValid() && valueField.Kind() == reflect.String {
					val = testCustomType{Value: valueField.String()}
				} else {
					// Try interface conversion as last resort
					if dataVal.Type() == reflect.TypeOf(val) {
						val = dataVal.Interface().(testCustomType)
					} else {
						ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("expected testCustomType with Value field, got %T", ctx.Data)))
						return
					}
				}
			} else {
				ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("expected testCustomType, got %T", ctx.Data)))
				return
			}
		}
		*ptr = val
		if s.validator != nil && !s.validator(*ptr) {
			issue := ctx.Issue().SetMessage(s.errorMsg)
			if s.errorMsg == "" {
				issue.SetCode(zconst.IssueCodeCustom)
			}
			ctx.AddIssue(issue)
			return
		}
	default:
		ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("unsupported type %T", ctx.ValPtr)))
	}
}

func (s *testExperimentalSchema) Validate(ctx *p.SchemaCtx) {
	if s.validator != nil {
		// ctx.ValPtr is already a pointer to the value
		// Pass it directly to validator which expects *T for validation
		if !s.validator(ctx.ValPtr) {
			issue := ctx.Issue().SetMessage(s.errorMsg)
			if s.errorMsg == "" {
				issue.SetCode(zconst.IssueCodeCustom)
			}
			ctx.AddIssue(issue)
		}
	}
}

func (s *testExperimentalSchema) GetType() zconst.ZogType {
	return s.dtype
}

func (s *testExperimentalSchema) SetCoercer(c CoercerFunc) {
	s.coercer = c
}

// Test type for custom schema tests
type testCustomType struct {
	Value string
}

// ============================================================================
// 1. Basic Use() and CustomSchema Tests
// ============================================================================

func TestUseWithCustomSchema_Parse(t *testing.T) {
	type TestStruct struct {
		Value string
	}

	validator := func(val any) bool {
		if str, ok := val.(string); ok {
			return len(str) > 3
		}
		return false
	}

	testSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeString,
		errorMsg:  "string too short",
	}

	schema := Struct(Shape{
		"value": Use(testSchema),
	})

	var result TestStruct
	errs := schema.Parse(map[string]any{"value": "test"}, &result)
	assert.Empty(t, errs)
	assert.Equal(t, "test", result.Value)

	// Test validation failure
	errs = schema.Parse(map[string]any{"value": "ab"}, &result)
	assert.NotNil(t, errs)
	valueErrs := tutils.FindByPath(errs, "value")
	assert.NotEmpty(t, valueErrs)
	assert.Equal(t, "string too short", valueErrs[0].Message)
}

func TestUseWithCustomSchema_Validate(t *testing.T) {
	type TestStruct struct {
		Value int
	}

	validator := func(val any) bool {
		if ptr, ok := val.(*int); ok {
			return *ptr > 10
		}
		return false
	}

	testSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeNumber,
		errorMsg:  "value too small",
	}

	schema := Struct(Shape{
		"value": Use(testSchema),
	})

	// Test validation success
	testVal := TestStruct{Value: 15}
	errs := schema.Validate(&testVal)
	assert.Empty(t, errs)

	// Test validation failure
	testVal.Value = 5
	errs = schema.Validate(&testVal)
	assert.NotNil(t, errs)
	valueErrs := tutils.FindByPath(errs, "value")
	assert.NotEmpty(t, valueErrs)
	assert.Equal(t, "value too small", valueErrs[0].Message)
}

func TestUseSchemaGetType(t *testing.T) {
	testSchema := &testExperimentalSchema{
		dtype: zconst.TypeString,
	}

	wrapped := Use(testSchema)
	assert.Equal(t, zconst.TypeString, wrapped.getType())

	testSchema2 := &testExperimentalSchema{
		dtype: zconst.TypeNumber,
	}
	wrapped2 := Use(testSchema2)
	assert.Equal(t, zconst.TypeNumber, wrapped2.getType())
}

func TestUseSchemaSetCoercer(t *testing.T) {
	type TestStruct struct {
		Value string
	}

	testSchema := &testExperimentalSchema{
		dtype: zconst.TypeString,
	}

	customCoercer := func(original any) (value any, err error) {
		return "coerced", nil
	}

	testSchema.SetCoercer(customCoercer)
	wrapped := Use(testSchema)

	schema := Struct(Shape{
		"value": wrapped,
	})

	var result TestStruct
	errs := schema.Parse(map[string]any{"value": 123}, &result)
	// Coercer should be called and value should be coerced
	assert.Empty(t, errs)
	assert.Equal(t, "coerced", result.Value)
}

// ============================================================================
// 2. Integration with Struct Schema
// ============================================================================

func TestCustomSchemaInStruct_Parse(t *testing.T) {
	type User struct {
		Name string
		ID   testCustomType `zog:"id"`
	}

	validator := func(val any) bool {
		if ct, ok := val.(testCustomType); ok {
			return len(ct.Value) > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "invalid custom type",
	}

	schema := Struct(Shape{
		"name": String().Required(),
		"ID":   Use(customSchema),
	})

	var user User
	data := map[string]any{
		"name": "John",
		"id":   testCustomType{Value: "123"}, // Use "id" because of zog tag
	}

	errs := schema.Parse(data, &user)
	assert.Empty(t, errs)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "123", user.ID.Value)
}

func TestCustomSchemaInStruct_Validate(t *testing.T) {
	type User struct {
		Name string
		ID   testCustomType
	}

	validator := func(val any) bool {
		// Handle both *testCustomType (from Validate) and testCustomType (from Process)
		var ct testCustomType
		if ptr, ok := val.(*testCustomType); ok {
			ct = *ptr
		} else if v, ok := val.(testCustomType); ok {
			ct = v
		} else {
			return false
		}
		return len(ct.Value) > 0
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "invalid custom type",
	}

	schema := Struct(Shape{
		"name": String().Required(),
		"ID":   Use(customSchema),
	})

	// Test validation success
	user := User{
		Name: "John",
		ID:   testCustomType{Value: "123"},
	}
	errs := schema.Validate(&user)
	assert.Empty(t, errs)

	// Test validation failure
	user.ID = testCustomType{Value: ""}
	errs = schema.Validate(&user)
	assert.NotNil(t, errs)
	idErrs := tutils.FindByPath(errs, "ID")
	assert.NotEmpty(t, idErrs)
}

func TestCustomSchemaInStruct_NestedStruct(t *testing.T) {
	type Address struct {
		Street string
		Zip    testCustomType
	}

	type User struct {
		Name    string
		Address Address
	}

	validator := func(val any) bool {
		if ct, ok := val.(testCustomType); ok {
			return len(ct.Value) == 5
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "zip code must be 5 digits",
	}

	schema := Struct(Shape{
		"name": String().Required(),
		"address": Struct(Shape{
			"street": String().Required(),
			"zip":    Use(customSchema),
		}),
	})

	var user User
	data := map[string]any{
		"name": "John",
		"address": map[string]any{
			"street": "123 Main St",
			"zip":    testCustomType{Value: "12345"},
		},
	}

	errs := schema.Parse(data, &user)
	assert.Empty(t, errs)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "123 Main St", user.Address.Street)
	assert.Equal(t, "12345", user.Address.Zip.Value)
}

func TestCustomSchemaInStruct_ValidationError(t *testing.T) {
	type User struct {
		Name string
		ID   testCustomType `zog:"id"`
	}

	validator := func(val any) bool {
		// Handle both *testCustomType (from Validate) and testCustomType (from Process)
		var ct testCustomType
		if ptr, ok := val.(*testCustomType); ok {
			ct = *ptr
		} else if v, ok := val.(testCustomType); ok {
			ct = v
		} else {
			return false
		}
		return len(ct.Value) > 3
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "ID too short",
	}

	schema := Struct(Shape{
		"name": String().Required(),
		"ID":   Use(customSchema),
	})

	var user User
	data := map[string]any{
		"name": "John",
		"id":   testCustomType{Value: "ab"}, // Use "id" because of zog tag
	}

	errs := schema.Parse(data, &user)
	assert.NotNil(t, errs)
	// Check that we have errors - the zog tag "id" means the error path will be "id", not "ID"
	idErrs := tutils.FindByPath(errs, "id")
	if len(idErrs) == 0 {
		idErrs = tutils.FindByPath(errs, "ID")
	}
	assert.NotEmpty(t, idErrs, "Expected error at path 'id' or 'ID'")
	if idErrs[0].Message != "" {
		assert.Equal(t, "ID too short", idErrs[0].Message)
	} else {
		// If message is empty, check that we have a custom error code
		assert.Equal(t, zconst.IssueCodeCustom, idErrs[0].Code)
	}
	// Don't verify default messages for custom types - they don't have i18n support
}

// ============================================================================
// 3. Integration with Slice Schema
// ============================================================================

func TestCustomSchemaInSlice_Parse(t *testing.T) {
	validator := func(val any) bool {
		if str, ok := val.(string); ok {
			return len(str) > 2
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeString,
		errorMsg:  "string too short",
	}

	schema := Slice(Use(customSchema))

	var result []string
	data := []any{"abc", "def", "ghi"}

	errs := schema.Parse(data, &result)
	assert.Empty(t, errs)
	assert.Len(t, result, 3)
	assert.Equal(t, "abc", result[0])
	assert.Equal(t, "def", result[1])
	assert.Equal(t, "ghi", result[2])
}

func TestCustomSchemaInSlice_Validate(t *testing.T) {
	validator := func(val any) bool {
		if ptr, ok := val.(*int); ok {
			return *ptr > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeNumber,
		errorMsg:  "value must be positive",
	}

	schema := Slice(Use(customSchema))

	// Test validation success
	values := []int{1, 2, 3}
	errs := schema.Validate(&values)
	assert.Empty(t, errs)

	// Test validation failure
	values = []int{1, -1, 3}
	errs = schema.Validate(&values)
	assert.NotNil(t, errs)
	indexErrs := tutils.FindByPath(errs, "[1]")
	assert.NotEmpty(t, indexErrs)
}

func TestCustomSchemaInSlice_ErrorPaths(t *testing.T) {
	validator := func(val any) bool {
		if str, ok := val.(string); ok {
			return len(str) > 3
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeString,
		errorMsg:  "string too short",
	}

	schema := Slice(Use(customSchema))

	var result []string
	data := []any{"abc", "ab", "defg"}

	errs := schema.Parse(data, &result)
	assert.NotNil(t, errs)
	indexErrs := tutils.FindByPath(errs, "[1]")
	assert.NotEmpty(t, indexErrs)
	assert.Equal(t, "string too short", indexErrs[0].Message)
	// Don't verify default messages for custom types in slices
}

func TestSliceOfCustomSchemaInStruct(t *testing.T) {
	type Team struct {
		Name  string
		Users []testCustomType
	}

	validator := func(val any) bool {
		if ct, ok := val.(testCustomType); ok {
			return len(ct.Value) > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "invalid user",
	}

	schema := Struct(Shape{
		"name":  String().Required(),
		"users": Slice(Use(customSchema)),
	})

	var team Team
	data := map[string]any{
		"name": "Team A",
		"users": []any{
			testCustomType{Value: "user1"},
			testCustomType{Value: "user2"},
		},
	}

	errs := schema.Parse(data, &team)
	assert.Empty(t, errs)
	assert.Equal(t, "Team A", team.Name)
	assert.Len(t, team.Users, 2)
	assert.Equal(t, "user1", team.Users[0].Value)
	assert.Equal(t, "user2", team.Users[1].Value)
}

// ============================================================================
// 4. Integration with Pointer Schema
// ============================================================================

func TestPtrCustomSchema_Parse_Nil(t *testing.T) {
	validator := func(val any) bool {
		if str, ok := val.(string); ok {
			return len(str) > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeString,
		errorMsg:  "invalid value",
	}

	schema := Ptr(Use(customSchema))

	var result *string
	errs := schema.Parse(nil, &result)
	assert.Empty(t, errs)
	assert.Nil(t, result)
}

func TestPtrCustomSchema_Parse_Value(t *testing.T) {
	validator := func(val any) bool {
		if str, ok := val.(string); ok {
			return len(str) > 3
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeString,
		errorMsg:  "string too short",
	}

	schema := Ptr(Use(customSchema))

	var result *string
	errs := schema.Parse("test", &result)
	assert.Empty(t, errs)
	assert.NotNil(t, result)
	assert.Equal(t, "test", *result)
}

func TestPtrCustomSchema_Validate_Nil(t *testing.T) {
	validator := func(val any) bool {
		if str, ok := val.(string); ok {
			return len(str) > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeString,
		errorMsg:  "invalid value",
	}

	schema := Ptr(Use(customSchema))

	var result *string
	errs := schema.Validate(&result)
	assert.Empty(t, errs)
}

func TestPtrCustomSchema_Validate_Value(t *testing.T) {
	type TestStruct struct {
		Value *int
	}

	validator := func(val any) bool {
		if ptr, ok := val.(*int); ok {
			return *ptr > 10
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeNumber,
		errorMsg:  "value too small",
	}

	schema := Struct(Shape{
		"value": Ptr(Use(customSchema)),
	})

	val := 15
	testVal := TestStruct{Value: &val}
	errs := schema.Validate(&testVal)
	assert.Empty(t, errs)

	val = 5
	errs = schema.Validate(&testVal)
	assert.NotNil(t, errs)
	valueErrs := tutils.FindByPath(errs, "value")
	assert.NotEmpty(t, valueErrs)
}

func TestPtrCustomSchemaInStruct(t *testing.T) {
	type User struct {
		Name string
		ID   *testCustomType
	}

	validator := func(val any) bool {
		if ct, ok := val.(testCustomType); ok {
			return len(ct.Value) > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "invalid ID",
	}

	schema := Struct(Shape{
		"name": String().Required(),
		"ID":   Ptr(Use(customSchema)),
	})

	var user User
	data := map[string]any{
		"name": "John",
		"ID":   testCustomType{Value: "123"}, // Use "ID" because no zog tag
	}

	errs := schema.Parse(data, &user)
	assert.Empty(t, errs)
	assert.Equal(t, "John", user.Name)
	assert.NotNil(t, user.ID)
	assert.Equal(t, "123", user.ID.Value)

	// Test with nil - reset user first
	user = User{}
	data["ID"] = nil
	errs = schema.Parse(data, &user)
	assert.Empty(t, errs)
	assert.Nil(t, user.ID)
}

// ============================================================================
// 5. Complex Nested Scenarios
// ============================================================================

func TestNestedStructWithSliceOfCustomSchema(t *testing.T) {
	type Item struct {
		Tags []testCustomType
	}

	type Container struct {
		Items []Item
	}

	validator := func(val any) bool {
		if ct, ok := val.(testCustomType); ok {
			return len(ct.Value) > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "invalid tag",
	}

	schema := Struct(Shape{
		"items": Slice(Struct(Shape{
			"tags": Slice(Use(customSchema)),
		})),
	})

	var container Container
	data := map[string]any{
		"items": []any{
			map[string]any{
				"tags": []any{
					testCustomType{Value: "tag1"},
					testCustomType{Value: "tag2"},
				},
			},
		},
	}

	errs := schema.Parse(data, &container)
	assert.Empty(t, errs)
	assert.Len(t, container.Items, 1)
	assert.Len(t, container.Items[0].Tags, 2)
	assert.Equal(t, "tag1", container.Items[0].Tags[0].Value)
}

func TestNestedSliceOfStructWithCustomSchema(t *testing.T) {
	type User struct {
		ID testCustomType
	}

	validator := func(val any) bool {
		if ct, ok := val.(testCustomType); ok {
			return len(ct.Value) > 0
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "invalid ID",
	}

	schema := Slice(Struct(Shape{
		"ID": Use(customSchema), // Match struct field name
	}))

	var users []User
	data := []any{
		map[string]any{"ID": testCustomType{Value: "id1"}},
		map[string]any{"ID": testCustomType{Value: "id2"}},
	}

	errs := schema.Parse(data, &users)
	assert.Empty(t, errs)
	assert.Len(t, users, 2)
	assert.Equal(t, "id1", users[0].ID.Value)
	assert.Equal(t, "id2", users[1].ID.Value)
}

func TestNestedPtrToSliceOfCustomSchema(t *testing.T) {
	validator := func(val any) bool {
		if str, ok := val.(string); ok {
			return len(str) > 2
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeString,
		errorMsg:  "string too short",
	}

	schema := Ptr(Slice(Use(customSchema)))

	var result *[]string
	data := []any{"abc", "def"}

	errs := schema.Parse(data, &result)
	assert.Empty(t, errs)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	assert.Equal(t, "abc", (*result)[0])
	assert.Equal(t, "def", (*result)[1])

	// Test with nil - reset result first
	result = nil
	errs = schema.Parse(nil, &result)
	assert.Empty(t, errs)
	assert.Nil(t, result)
}

// ============================================================================
// 6. Error Handling
// ============================================================================

func TestCustomSchemaCoercionError(t *testing.T) {
	type TestStruct struct {
		Value string
	}

	customCoercer := func(original any) (value any, err error) {
		return nil, fmt.Errorf("cannot coerce %T to string", original)
	}

	testSchema := &testExperimentalSchema{
		dtype:   zconst.TypeString,
		coercer: customCoercer,
	}

	schema := Struct(Shape{
		"value": Use(testSchema),
	})

	var result TestStruct
	errs := schema.Parse(map[string]any{"value": 123}, &result)
	assert.NotNil(t, errs)
	valueErrs := tutils.FindByPath(errs, "value")
	assert.NotEmpty(t, valueErrs)
	// The error might be formatted, so check both the message and error field
	assert.True(t,
		valueErrs[0].Message != "" ||
			(valueErrs[0].Err != nil && valueErrs[0].Err.Error() != ""),
		"Expected error message or error field to be set")
}

func TestCustomSchemaValidationError(t *testing.T) {
	type TestStruct struct {
		Value int
	}

	validator := func(val any) bool {
		if ptr, ok := val.(*int); ok {
			return *ptr >= 0
		}
		return false
	}

	testSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     zconst.TypeNumber,
		errorMsg:  "value cannot be negative",
	}

	schema := Struct(Shape{
		"value": Use(testSchema),
	})

	testVal := TestStruct{Value: -5}
	errs := schema.Validate(&testVal)
	assert.NotNil(t, errs)
	valueErrs := tutils.FindByPath(errs, "value")
	assert.NotEmpty(t, valueErrs)
	assert.Equal(t, "value cannot be negative", valueErrs[0].Message)
}

func TestCustomSchemaErrorPath(t *testing.T) {
	type Nested struct {
		Value testCustomType
	}

	type Container struct {
		Nested Nested
	}

	validator := func(val any) bool {
		if ct, ok := val.(testCustomType); ok {
			return len(ct.Value) > 3
		}
		return false
	}

	customSchema := &testExperimentalSchema{
		validator: validator,
		dtype:     "custom",
		errorMsg:  "value too short",
	}

	schema := Struct(Shape{
		"nested": Struct(Shape{
			"value": Use(customSchema),
		}),
	})

	var container Container
	data := map[string]any{
		"nested": map[string]any{
			"value": testCustomType{Value: "ab"},
		},
	}

	errs := schema.Parse(data, &container)
	assert.NotNil(t, errs)
	nestedErrs := tutils.FindByPath(errs, "nested.value")
	assert.NotEmpty(t, nestedErrs)
	assert.Equal(t, "value too short", nestedErrs[0].Message)
	// Don't verify default messages for custom types - they don't have i18n support
}

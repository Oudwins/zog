package zog

import (
	"errors"
	"testing"
	"time"

	"github.com/Oudwins/zog/pkgs/internals/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestValidateUnionFirstSchemaSucceeds(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Int().GT(10, Message("must be greater than 10")),
		Int().LT(0, Message("must be less than 0")),
	})
	dest := 15

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 15, dest)
}

func TestValidateUnionLaterSchemaSucceeds(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Int().GT(10, Message("must be greater than 10")),
		Int().LT(0, Message("must be less than 0")),
	})
	dest := -5

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, -5, dest)
}

func TestValidateUnionAllSchemasFail(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Int().GT(10, Message("must be greater than 10")),
		Int().LT(0, Message("must be less than 0")),
	})
	dest := 5

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, "must be greater than 10", errs[0].Message)
	assert.Equal(t, "must be less than 0", errs[1].Message)
	assert.Equal(t, 5, dest)
}

func TestValidateUnionDefaultBranchSucceeds(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Required(),
		Int().Default(42),
	})
	dest := 0

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 42, dest)
}

func TestValidateUnionCatchBranchSucceeds(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Required(),
		Int().GT(10, Message("must be greater than 10")).Catch(42),
	})
	dest := 5

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 42, dest)
}

func TestValidateUnionInvalidDestinationTypeAllSchemasFail(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Required(),
		Int().Required(),
	})
	dest := true

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, zconst.IssueCodeInvalidType, errs[0].Code)
	assert.Equal(t, zconst.IssueCodeInvalidType, errs[1].Code)
	assert.True(t, dest)
}

func TestValidateUnionStringOrInt(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Required(),
		Int().Required(),
	})
	dest := 15

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 15, dest)
}

func TestValidateUnionPrimitiveSchemaMatrix(t *testing.T) {
	tests := []struct {
		name string
		dest any
	}{
		{name: "string", dest: tutils.PtrOf("hello")},
		{name: "int", dest: tutils.PtrOf(42)},
		{name: "float", dest: tutils.PtrOf(42.5)},
		{name: "bool", dest: tutils.PtrOf(true)},
		{name: "time", dest: tutils.PtrOf(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))},
	}

	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Required(),
		Int().Required(),
		Float().Required(),
		Bool().Required(),
		Time().Required(),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validator.Validate(tt.dest)

			assert.Empty(t, errs)
		})
	}
}

func TestValidateUnionContainerSchemaMatrix(t *testing.T) {
	type unionContainerStruct struct {
		Name string
	}

	tests := []struct {
		name      string
		validator *UnionSchema
		dest      any
	}{
		{
			name: "struct",
			validator: EXPERIMENTAL_UNION([]ZogSchema{
				Struct(Shape{"name": String().Required()}),
				Slice(String()).Required(),
			}),
			dest: &unionContainerStruct{Name: "zog"},
		},
		{
			name: "slice",
			validator: EXPERIMENTAL_UNION([]ZogSchema{
				Struct(Shape{"name": String().Required()}),
				Slice(String()).Min(1),
			}),
			dest: &[]string{"zog"},
		},
		{
			name: "map",
			validator: EXPERIMENTAL_UNION([]ZogSchema{
				Struct(Shape{"name": String().Required()}),
				EXPERIMENTAL_MAP[string, int](String().Required(), Int()).Min(1),
			}),
			dest: &map[string]int{"zog": 1},
		},
		{
			name: "pointer",
			validator: EXPERIMENTAL_UNION([]ZogSchema{
				String().Required(),
				Ptr(Int()).NotNil(),
			}),
			dest: tutils.PtrOf(tutils.PtrOf(42)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.validator.Validate(tt.dest)

			assert.Empty(t, errs)
		})
	}
}

func TestValidateUnionMixedSchemaTypesAllFail(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Min(3, Message("string too short")),
		Int().GT(10, Message("int too small")),
		Bool().Required(Message("bool required")),
	})
	dest := "no"

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 3)
	assert.Equal(t, "string too short", errs[0].Message)
	assert.Equal(t, zconst.IssueCodeInvalidType, errs[1].Code)
	assert.Equal(t, zconst.IssueCodeInvalidType, errs[2].Code)
	assert.Equal(t, "no", dest)
}

func TestValidateUnionShortCircuitsAfterFirstSuccess(t *testing.T) {
	firstCalls := 0
	secondCalls := 0
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Int().TestFunc(func(val *int, ctx Ctx) bool {
			firstCalls++
			return true
		}),
		Int().TestFunc(func(val *int, ctx Ctx) bool {
			secondCalls++
			return true
		}),
	})
	dest := 1

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 1, firstCalls)
	assert.Equal(t, 0, secondCalls)
}

func TestValidateUnionRunsLaterSchemasAfterFailure(t *testing.T) {
	firstCalls := 0
	secondCalls := 0
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Int().TestFunc(func(val *int, ctx Ctx) bool {
			firstCalls++
			return false
		}, Message("first failed")),
		Int().TestFunc(func(val *int, ctx Ctx) bool {
			secondCalls++
			return true
		}),
	})
	dest := 1

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, 1, firstCalls)
	assert.Equal(t, 1, secondCalls)
}

func TestValidateUnionDoesNotCommitFailedBranchMutation(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Trim().Len(3, Message("trimmed string must have length 3")),
		String().OneOf([]string{"x"}),
	})
	dest := " x "

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, " x ", dest)
}

func TestValidateUnionDoesNotCommitNestedFailedBranchMutation(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Slice(String().Trim().Len(2)),
		Slice(String().OneOf([]string{"x"})),
	})
	dest := []string{" x "}

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, []string{" x "}, dest)
}

func TestValidateUnionPreservesUnchangedDestinationIdentity(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{Slice(String().Required())})
	dest := []string{"zog"}
	originalItem := &dest[0]

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Same(t, originalItem, &dest[0])
}

func TestValidateUnionUsesFreshContextForEachBranch(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Transform(func(val *string, ctx Ctx) error {
			return errors.New("first branch failed")
		}),
		String().Min(1).Len(3),
	})
	dest := "ab"

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, zconst.TypeString, errs[1].Dtype)
	assert.Equal(t, "string must be exactly 3 character(s)", errs[1].Message)
	assert.Equal(t, "ab", dest)
}

func TestParseUnionUsesFreshContextForEachBranch(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Transform(func(val *string, ctx Ctx) error {
			return errors.New("first branch failed")
		}),
		String().Min(1).Len(3),
	})
	dest := "original"

	errs := validator.Parse("ab", &dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, zconst.TypeString, errs[1].Dtype)
	assert.Equal(t, "string must be exactly 3 character(s)", errs[1].Message)
	assert.Equal(t, "original", dest)
}

func TestParseUnionIsolatesMutableInputBetweenBranches(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Preprocess(func(data map[string]int, ctx Ctx) (int, error) {
			data["value"] = 2
			return 0, errors.New("first branch failed")
		}, Int()),
		Preprocess(func(data map[string]int, ctx Ctx) (int, error) {
			return data["value"], nil
		}, Int().OneOf([]int{1})),
	})
	input := map[string]int{"value": 1}
	dest := 0

	errs := validator.Parse(input, &dest)

	assert.Empty(t, errs)
	assert.Equal(t, 1, dest)
	assert.Equal(t, map[string]int{"value": 1}, input)
}

func TestParseUnionReusesUnchangedInputClone(t *testing.T) {
	type input struct {
		Value int
	}
	var firstInput *input
	var secondInput *input
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Preprocess(func(data *input, ctx Ctx) (int, error) {
			firstInput = data
			return 0, errors.New("first branch failed")
		}, Int()),
		Preprocess(func(data *input, ctx Ctx) (int, error) {
			secondInput = data
			return data.Value, nil
		}, Int()),
	})
	originalInput := &input{Value: 1}
	dest := 0

	errs := validator.Parse(originalInput, &dest)

	assert.Empty(t, errs)
	assert.NotSame(t, originalInput, firstInput)
	assert.Same(t, firstInput, secondInput)
}

func TestValidateUnionReusesUnchangedOutputClone(t *testing.T) {
	var firstOutput *int
	var secondOutput *int
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		Int().TestFunc(func(val *int, ctx Ctx) bool {
			firstOutput = val
			return false
		}),
		Int().TestFunc(func(val *int, ctx Ctx) bool {
			secondOutput = val
			return true
		}),
	})
	dest := 1

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.NotSame(t, &dest, firstOutput)
	assert.Same(t, firstOutput, secondOutput)
}

func TestValidateUnionPreservesIssueValuesWhenReusingClone(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		CustomFunc(func(val *int, ctx Ctx) bool {
			ctx.AddIssue(ctx.Issue().SetMessage("first branch failed"))
			return true
		}),
		CustomFunc(func(val *int, ctx Ctx) bool {
			*val = 2
			ctx.AddIssue(ctx.Issue().SetMessage("second branch failed"))
			return true
		}),
	})
	dest := 1

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 2)
	assert.Equal(t, 1, *errs[0].Value.(*int))
	assert.Equal(t, 1, dest)
}

func TestValidateUnionInStructField(t *testing.T) {
	type testStruct struct {
		Value string
	}
	validator := Struct(Shape{
		"value": EXPERIMENTAL_UNION([]ZogSchema{
			String().Len(3, Message("must have length 3")),
			String().HasPrefix("z", Message("must start with z")),
		}),
	})

	stringDest := testStruct{Value: "zog"}
	errs := validator.Validate(&stringDest)
	assert.Empty(t, errs)

	laterBranchDest := testStruct{Value: "zod"}
	errs = validator.Validate(&laterBranchDest)
	assert.Empty(t, errs)

	failingDest := testStruct{Value: "go"}
	errs = validator.Validate(&failingDest)
	assert.NotEmpty(t, errs)
}

func TestValidateUnionStructBranchesPreserveNestedErrors(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Required().Min(3),
		Struct(Shape{"name": String().Required(Message("name required"))}),
		Struct(Shape{"age": Int().GT(18, Message("age too low"))}),
	})
	dest := user{Age: 10}

	errs := validator.Validate(&dest)

	assert.Len(t, errs, 3)
	assert.Equal(t, zconst.IssueCodeInvalidType, errs[0].Code)
	assert.Equal(t, "name required", errs[1].Message)
	assert.Equal(t, "age too low", errs[2].Message)
}

func TestValidateUnionAnySchemaBranch(t *testing.T) {
	validator := EXPERIMENTAL_UNION([]ZogSchema{
		String().Len(10),
		EXPERIMENTAL_ANY().Required(),
	})
	var dest any = 42

	errs := validator.Validate(&dest)

	assert.Empty(t, errs)
	assert.Equal(t, any(42), dest)
}

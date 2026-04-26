package zog

import (
	"testing"

	p "github.com/Oudwins/zog/pkgs/internals"
	"github.com/stretchr/testify/assert"
)

func nullableStrPtr(s string) *string { return &s }

func TestPtrNullable_AbsentKeyPreservesExistingValue(t *testing.T) {
	type Req struct {
		Tag *string
	}
	schema := Struct(Shape{
		"tag": Ptr(String()),
	})
	out := Req{Tag: nullableStrPtr("keep-me")}
	errs := schema.Parse(map[string]any{}, &out)
	assert.Empty(t, errs)
	assert.NotNil(t, out.Tag)
	assert.Equal(t, "keep-me", *out.Tag)
}

func TestPtrNullable_ExplicitNullClearsPointer(t *testing.T) {
	type Req struct {
		Tag *string
	}
	schema := Struct(Shape{
		"tag": Ptr(String()),
	})
	out := Req{Tag: nullableStrPtr("clear-me")}
	errs := schema.Parse(map[string]any{"tag": nil}, &out)
	assert.Empty(t, errs)
	assert.Nil(t, out.Tag)
}

func TestPtrNullable_ConcreteValueOverwrites(t *testing.T) {
	type Req struct {
		Tag *string
	}
	schema := Struct(Shape{
		"tag": Ptr(String()),
	})
	var out Req
	errs := schema.Parse(map[string]any{"tag": "v"}, &out)
	assert.Empty(t, errs)
	assert.NotNil(t, out.Tag)
	assert.Equal(t, "v", *out.Tag)
}

func TestPtrNullable_NestedStructPointer_NullClearsOuter(t *testing.T) {
	type Inner struct {
		V int
	}
	type Outer struct {
		Inner *Inner
	}
	schema := Struct(Shape{
		"inner": Ptr(Struct(Shape{
			"v": Int(),
		})),
	})
	out := Outer{Inner: &Inner{V: 42}}
	errs := schema.Parse(map[string]any{"inner": nil}, &out)
	assert.Empty(t, errs)
	assert.Nil(t, out.Inner)
}

func TestPtrNullable_NestedBareStruct_NullBehavesAsEmpty(t *testing.T) {
	type Inner struct {
		V *int
	}
	type Outer struct {
		Inner Inner
	}
	schema := Struct(Shape{
		"inner": Struct(Shape{
			"v": Ptr(Int()),
		}),
	})
	v := 7
	out := Outer{Inner: Inner{V: &v}}
	errs := schema.Parse(map[string]any{"inner": nil}, &out)
	assert.Empty(t, errs)
	assert.NotNil(t, out.Inner.V)
	assert.Equal(t, 7, *out.Inner.V)
}

func TestPtrNullable_DoublePointer_OuterNullable(t *testing.T) {
	type Req struct {
		V **string
	}
	schema := Struct(Shape{
		"v": Ptr(Ptr(String())),
	})
	inner := "keep"
	outer := &inner
	out := Req{V: &outer}
	errs := schema.Parse(map[string]any{"v": nil}, &out)
	assert.Empty(t, errs)
	assert.Nil(t, out.V)
}

func TestPtrNullable_SliceOfNullablePtr(t *testing.T) {
	schema := Slice(Ptr(String()))
	var out []*string
	errs := schema.Parse([]any{"a", nil, "c"}, &out)
	assert.Empty(t, errs)
	assert.Len(t, out, 3)
	assert.NotNil(t, out[0])
	assert.Equal(t, "a", *out[0])
	assert.Nil(t, out[1])
	assert.NotNil(t, out[2])
	assert.Equal(t, "c", *out[2])
}

func TestPtrNullable_StructInput_NilFieldClears(t *testing.T) {
	// Shape key must match the Go field name for struct-to-struct parsing
	// because StructDataProvider.Get uses FieldByName on the resolved key.
	type Req struct {
		Tag *string
	}
	schema := Struct(Shape{
		"Tag": Ptr(String()),
	})
	src := Req{Tag: nil}
	dst := Req{Tag: nullableStrPtr("clear-me")}
	errs := schema.Parse(src, &dst)
	assert.Empty(t, errs)
	assert.Nil(t, dst.Tag)
}

func TestPtrNullable_TopLevelBareNilDoesNotClear(t *testing.T) {
	schema := Ptr(String())
	dest := nullableStrPtr("keep")
	errs := schema.Parse(nil, &dest)
	assert.Empty(t, errs)
	assert.NotNil(t, dest)
	assert.Equal(t, "keep", *dest)
}

func TestPtrNullable_IssueValueIsNotSentinel(t *testing.T) {
	type Req struct {
		Tag *string
	}
	schema := Struct(Shape{
		"tag": Ptr(String()).NotNil(),
	})
	var out Req
	errs := schema.Parse(map[string]any{"tag": nil}, &out)
	assert.NotEmpty(t, errs)
	assert.False(t, p.IsExplicitNull(errs[0].Value), "issue value must not be the internal sentinel")
}

func TestPtrNullable_MapDataProvider_EmitsSentinelForExplicitNull(t *testing.T) {
	m := map[string]any{"present_null": nil, "present_val": "hello"}
	dp := p.NewSafeMapDataProvider(m)
	assert.Nil(t, dp.Get("absent_key"))
	assert.True(t, p.IsExplicitNull(dp.Get("present_null")))
	assert.Equal(t, "hello", dp.Get("present_val"))
}

func TestPtrNullable_MapDataProvider_TypedMap_NoSentinel(t *testing.T) {
	m := map[string]string{"empty": ""}
	dp := p.NewSafeMapDataProvider(m)
	assert.Equal(t, "", dp.Get("empty"))
	assert.Nil(t, dp.Get("missing"))
	assert.False(t, p.IsExplicitNull(dp.Get("empty")))
	assert.False(t, p.IsExplicitNull(dp.Get("missing")))
}

func TestPtrNullable_TypedNilInMapAnyValue_EmitsSentinel(t *testing.T) {
	var typedNil *string
	dp := p.NewSafeMapDataProvider(map[string]any{"tag": typedNil})
	assert.True(t, p.IsExplicitNull(dp.Get("tag")))
}

func TestPtrNullable_TypedNilInMapAnyValue_ClearsPointer(t *testing.T) {
	type Req struct {
		Tag *string
	}
	schema := Struct(Shape{
		"tag": Ptr(String()),
	})
	var typedNil *string
	input := map[string]any{"tag": typedNil}
	out := Req{Tag: nullableStrPtr("preexisting")}
	errs := schema.Parse(input, &out)
	assert.Empty(t, errs)
	assert.Nil(t, out.Tag)
}

func TestPtrNullable_TypedNilInInterfaceStructField_EmitsSentinel(t *testing.T) {
	type Req struct {
		Thing any
	}
	src := Req{Thing: (*string)(nil)}
	dp, err := p.TryNewAnyDataProvider(src)
	assert.NoError(t, err)
	assert.True(t, p.IsExplicitNull(dp.Get("Thing")))
}


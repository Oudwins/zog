package zssschema

import (
	z "github.com/Oudwins/zog"
	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/zconst"
)

// ZSSGoTypeSchema defines the schema for ZSSGoType
var ZSSGoTypeSchema = z.Struct(z.Shape{
	"pkgPath": z.String(),
	"name":    z.String(),
	"display": z.String().Required(),
})

var ZSSFieldMetaSchema = z.Struct(z.Shape{
	"tags": z.String(),
})

// ZSSTransformerSchema defines the schema for ZSSTransformer
var ZSSTransformerSchema = z.Struct(z.Shape{
	"id": z.String().Required(),
})

// ZSSTestSchema defines the schema for ZSSTest
var ZSSTestSchema = z.Struct(z.Shape{
	"id":        z.String().Required(),
	"message":   z.String().Required(),
	"issuePath": z.Slice(z.String()).Required(),
	"params":    z.EXPERIMENTAL_MAP[string, any](z.String(), z.EXPERIMENTAL_ANY()),
})

// ZSSProcessorSchema defines the schema for ZSSProcessor
var ZSSProcessorSchema = z.Struct(z.Shape{
	"kind":        z.String().Required(),
	"test":        z.Ptr(ZSSTestSchema),
	"transformer": z.Ptr(ZSSTransformerSchema),
})

// ZSSExtensionSchema defines the schema for ZSSExtension.
var ZSSExtensionSchema = z.Struct(z.Shape{
	"URI":     URISchema,
	"Content": z.EXPERIMENTAL_ANY(),
})

var zssKind = func(k zconst.ZogType) *z.StringSchema[zconst.ZogType] {
	return z.StringLike[zconst.ZogType]().OneOf([]zconst.ZogType{k})

}

// ZSSSchemaSchema defines the schema for ZSSSchema.
// Note: defaultValue and catchValue are intentionally loose because ZSS allows arbitrary values.
var ZSSSchemaSchema = z.EXPERIMENTAL_RECURSIVE(func(self z.RecursiveSchema[*z.UnionSchema]) *z.UnionSchema {
	ref := z.Struct(z.Shape{
		"Ref": z.String().Required().Min(1),
	})

	str := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeString),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	num := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeNumber),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})
	// bool
	bl := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeBool),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	// time
	tm := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeTime),
		"format":       z.Ptr(z.String()),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	list := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeSlice),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"element":      z.Ptr(self()),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	mp := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeMap),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"key":          z.Ptr(self()),
		"value":        z.Ptr(self()),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	strct := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeStruct),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"fields":       z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(self())),
		"fieldMeta":    z.EXPERIMENTAL_MAP[string, zsscore.ZSSFieldMeta](z.String(), ZSSFieldMetaSchema),
	})

	ptr := z.Struct(z.Shape{
		"Kind":     zssKind(zconst.TypeStruct),
		"element":  z.Ptr(self()),
		"goTypes":  z.Slice(ZSSGoTypeSchema),
		"required": z.Ptr(ZSSTestSchema),
	})

	custom := z.Struct(z.Shape{
		"Kind":       zssKind(zconst.TypeStruct),
		"goTypes":    z.Slice(ZSSGoTypeSchema),
		"processors": z.Slice(ZSSProcessorSchema),
	})

	s := z.Struct(z.Shape{
		"Ref":          z.Ptr(z.String()),
		"kind":         z.StringLike[zconst.ZogType]().OneOf(zconst.ZogTypeValues),
		"Extension":    z.Ptr(ZSSExtensionSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"format":       z.Ptr(z.String()),
		"processors":   z.Slice(ZSSProcessorSchema),
		"fields":       z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(self())),
		"fieldMeta":    z.EXPERIMENTAL_MAP[string, zsscore.ZSSFieldMeta](z.String(), ZSSFieldMetaSchema),
		"element":      z.Ptr(self()),
		"key":          z.Ptr(self()),
		"value":        z.Ptr(self()),
		"children":     z.Slice(z.Ptr(self())),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	return z.Union([]z.ZogSchema{ref, str, num, bl, tm, list, mp, strct, ptr, custom})
})

var URISchema = z.String().Match(zsscore.ZSS_URI_REGEX).Required()

// ZSSDocumentSchema defines the schema for ZSSDocument
var ZSSDocumentSchema = z.Struct(z.Shape{
	"URI":  URISchema, // $Schema
	"Root": z.Ptr(ZSSSchemaSchema).NotNil(),
	"Defs": z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(ZSSSchemaSchema)),
})

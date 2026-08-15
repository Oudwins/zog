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
	"ID": z.String().Required(),
})

// ZSSTestSchema defines the schema for ZSSTest
var ZSSTestSchema = z.Struct(z.Shape{
	"ID":        z.String().Required(),
	"message":   z.String().Required(),
	"issuePath": z.Slice(z.String()),
	"params":    z.EXPERIMENTAL_MAP[string, any](z.String(), z.EXPERIMENTAL_ANY()),
})

// ZSSProcessorSchema defines the schema for ZSSProcessor
var ZSSProcessorSchema = z.Struct(z.Shape{
	"kind":        z.String().Required(),
	"test":        z.Ptr(ZSSTestSchema),
	"transformer": z.Ptr(ZSSTransformerSchema),
})

var ZSSCustomTestSchema = z.Struct(z.Shape{
	"ID":        z.String(),
	"message":   z.String(),
	"issuePath": z.Slice(z.String()),
	"params":    z.EXPERIMENTAL_MAP[string, any](z.String(), z.EXPERIMENTAL_ANY()),
})

var ZSSCustomProcessorSchema = z.Struct(z.Shape{
	"kind": z.StringLike[zconst.ZogProcessor]().OneOf([]zconst.ZogProcessor{zconst.ZogProcessorTest}),
	"test": z.Ptr(ZSSCustomTestSchema).NotNil(),
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
		"Ref": z.Ptr(z.String().Required().Min(1)).NotNil(),
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
		"element":      z.Ptr(self()).NotNil(),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	mp := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeMap),
		"processors":   z.Slice(ZSSProcessorSchema),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"key":          z.Ptr(self()).NotNil(),
		"value":        z.Ptr(self()).NotNil(),
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
		"fields":       z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(self()).NotNil()),
		"fieldMeta":    z.EXPERIMENTAL_MAP[string, zsscore.ZSSFieldMeta](z.String(), ZSSFieldMetaSchema),
	})

	ptr := z.Struct(z.Shape{
		"Kind":     zssKind(zconst.TypePtr),
		"element":  z.Ptr(self()).NotNil(),
		"goTypes":  z.Slice(ZSSGoTypeSchema),
		"required": z.Ptr(ZSSTestSchema),
	})

	custom := z.Struct(z.Shape{
		"Kind":       zssKind(zconst.TypeCustom),
		"goTypes":    z.Slice(ZSSGoTypeSchema),
		"processors": z.Slice(ZSSCustomProcessorSchema),
	})

	preprocess := z.Struct(z.Shape{
		"Kind":    zssKind(zconst.TypePreprocess),
		"element": z.Ptr(self()).NotNil(),
		"goTypes": z.Slice(ZSSGoTypeSchema),
	})

	boxed := z.Struct(z.Shape{
		"Kind":    zssKind(zconst.TypeBoxed),
		"element": z.Ptr(self()).NotNil(),
		"goTypes": z.Slice(ZSSGoTypeSchema),
	})

	anySchema := z.Struct(z.Shape{
		"Kind":         zssKind(zconst.TypeAny),
		"processors":   z.Slice(ZSSProcessorSchema),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})

	union := z.Struct(z.Shape{
		"Kind": zssKind(zconst.TypeUnion),
		// removed min 2 from here because we can't enforce it nicely via api and not an issue (beyond performance if people use it for a 1 or 0 schema union)
		"children": z.Slice(z.Ptr(self()).NotNil()).Required(),
	})

	extended := z.Struct(z.Shape{
		"Kind":      zssKind("extension"),
		"Extension": z.Ptr(ZSSExtensionSchema).NotNil(),
	})

	return z.EXPERIMENTAL_UNION([]z.ZogSchema{
		union, ref, str, num, bl, tm, list, mp, strct, ptr, custom, preprocess, boxed, anySchema, extended,
	})
})

var URISchema = z.String().Match(zsscore.ZSS_URI_REGEX).Required()

// ZSSDocumentSchema defines the schema for ZSSDocument
var ZSSDocumentSchema = z.Struct(z.Shape{
	"URI":  URISchema, // $Schema
	"Root": z.Ptr(ZSSSchemaSchema).NotNil(),
	"Defs": z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(ZSSSchemaSchema)),
})

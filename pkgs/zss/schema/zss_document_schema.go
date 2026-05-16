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

// ZSSSchemaSchema defines the schema for ZSSSchema.
// Note: defaultValue and catchValue are intentionally loose because ZSS allows arbitrary values.
var ZSSSchemaSchema = z.EXPERIMENTAL_RECURSIVE(func(self z.RecursiveSchema[*z.StructSchema]) *z.StructSchema {
	return z.Struct(z.Shape{
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
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})
})

var URISchema = z.String().Match(zsscore.ZSS_URI_REGEX).Required()

// ZSSDocumentSchema defines the schema for ZSSDocument
var ZSSDocumentSchema = z.Struct(z.Shape{
	"URI":  URISchema, // $Schema
	"Root": z.Ptr(ZSSSchemaSchema).NotNil(),
	"Defs": z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(ZSSSchemaSchema)),
})

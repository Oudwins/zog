package zssschema

import (
	z "github.com/Oudwins/zog"
	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
)

// ZSSGoTypeSchema defines the schema for ZSSGoType
var ZSSGoTypeSchema = z.Struct(z.Shape{
	"pkgPath": z.String().Required(),
	"name":    z.String().Required(),
	"display": z.String().Required(),
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

// ZSSSchemaSchema defines the schema for ZSSSchema.
// Note: defaultValue and catchValue are intentionally loose because ZSS allows arbitrary values.
var ZSSSchemaSchema = z.EXPERIMENTAL_RECURSIVE(func(self z.RecursiveSchema[*z.StructSchema]) *z.StructSchema {
	return z.Struct(z.Shape{
		"Ref":          z.Ptr(z.String()),
		"kind":         z.String(),
		"goTypes":      z.Slice(ZSSGoTypeSchema),
		"format":       z.Ptr(z.String()),
		"processors":   z.Slice(ZSSProcessorSchema),
		"fields":       z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(self())),
		"element":      z.Ptr(self()),
		"key":          z.Ptr(self()),
		"value":        z.Ptr(self()),
		"required":     z.Ptr(ZSSTestSchema),
		"defaultValue": z.EXPERIMENTAL_ANY(),
		"catchValue":   z.EXPERIMENTAL_ANY(),
	})
})

// ZSSDocumentSchema defines the schema for ZSSDocument
var ZSSDocumentSchema = z.Struct(z.Shape{
	"Version": z.String().Required(),
	"Root":    z.Ptr(ZSSSchemaSchema).NotNil(),
	"Defs":    z.EXPERIMENTAL_MAP[string, *zsscore.ZSSSchema](z.String(), z.Ptr(ZSSSchemaSchema)),
})

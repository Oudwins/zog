package zssschema

import (
	z "github.com/Oudwins/zog"
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
// Note: params field (map[string]any) cannot be strictly validated, so it's omitted from the schema
var ZSSTestSchema = z.Struct(z.Shape{
	"id":        z.String().Required(),
	"message":   z.String().Required(),
	"issuePath": z.Slice(z.String()).Required(),
	// params is map[string]any - cannot be strictly validated with zog schemas
})

// ZSSProcessorSchema defines the schema for ZSSProcessor
var ZSSProcessorSchema = z.Struct(z.Shape{
	"kind":        z.String().Required(),
	"test":        z.Ptr(ZSSTestSchema),
	"transformer": z.Ptr(ZSSTransformerSchema),
})

// ZSSSchemaSchema defines the schema for ZSSSchema
// Note: recursive child slots, DefaultValue, and CatchValue cannot be strictly validated.
var ZSSSchemaSchema = z.Struct(z.Shape{
	"kind":       z.String().Required(),
	"goTypes":    z.Slice(ZSSGoTypeSchema),
	"format":     z.Ptr(z.String()),
	"processors": z.Slice(ZSSProcessorSchema),
	// fields is map[string]*ZSSSchema - cannot be strictly validated due to recursive nature
	// element is *ZSSSchema - cannot be strictly validated due to recursive nature
	// key is *ZSSSchema - cannot be strictly validated due to recursive nature
	// value is *ZSSSchema - cannot be strictly validated due to recursive nature
	"required": z.Ptr(ZSSTestSchema),
	// defaultValue is any - cannot be strictly validated
	// catchValue is any - cannot be strictly validated
})

// ZSSDocumentSchema defines the schema for ZSSDocument
var ZSSDocumentSchema = z.Struct(z.Shape{
	"$schema": z.String().Required(),
	"root":    z.Ptr(ZSSSchemaSchema).NotNil(),
})

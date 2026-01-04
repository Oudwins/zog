package zss // Zog Schema Specification

import "github.com/Oudwins/zog/zconst" // TODO make zog schemas for all of these to validate them!

type ZSSDocument struct {
	Version ZSSVersion // "1.0.0"
	Schema  *ZSSSchema
}
type ZSSProcessor struct {
	Kind        zconst.ZogProcessor // "transform", "validator"
	Test        *ZSSTest
	Transformer *ZSSTransformer
}

type ZSSTest struct {
	ID        zconst.ZogIssueCode // issue code
	Message   string
	IssuePath *string
	Params    map[string]any
}

type ZSSTransformer struct {
	ID zconst.ZogTransformID
}

type ZSSSchema struct {
	Kind         string  // "string", "number", "bool", "time", "slice", "struct", "ptr"
	GoType       string  // Custom type if available (only if ZSS Exhaustive Metadata is enabled)
	Format       *string // Used for time.Time schemas only right now. (Only if ZSS Exhaustive Metadata is enabled)
	Processors   []ZSSProcessor
	Child        any // *ZSSSchema | map[string]ZSSSchema
	Required     *ZSSTest
	DefaultValue any
	CatchValue   any
}

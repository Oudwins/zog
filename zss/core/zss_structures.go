package zsscore // Zog Schema Specification

import "github.com/Oudwins/zog/zconst" // TODO make zog schemas for all of these to validate them!

type ZSSDocument struct {
	Version ZSSVersion `json:"version"` // "1.0.0"
	Schema  *ZSSSchema `json:"schema"`
}
type ZSSProcessor struct {
	Kind        zconst.ZogProcessor `json:"kind"` // "transform", "validator"
	Test        *ZSSTest            `json:"test"`
	Transformer *ZSSTransformer     `json:"transformer"`
}

type ZSSTest struct {
	ID        zconst.ZogIssueCode `json:"id"` // issue code
	Message   string              `json:"message"`
	IssuePath *string             `json:"issuePath"`
	Params    map[string]any      `json:"params"`
}

type ZSSTransformer struct {
	ID zconst.ZogTransformID `json:"id"`
}

type ZSSSchema struct {
	Kind         string         `json:"kind"`   // "string", "number", "bool", "time", "slice", "struct", "ptr"
	GoType       *string        `json:"goType"` // Custom type if available (only if ZSS Exhaustive Metadata is enabled)
	Format       *string        `json:"format"` // Used for time.Time schemas only right now. (Only if ZSS Exhaustive Metadata is enabled)
	Processors   []ZSSProcessor `json:"processors"`
	Child        any            `json:"child"` // *ZSSSchema | map[string]ZSSSchema
	Required     *ZSSTest       `json:"required"`
	DefaultValue any            `json:"defaultValue"`
	CatchValue   any            `json:"catchValue"`
}


package zss // Zog Schema Specification

import "github.com/Oudwins/zog/zconst" // TODO make zog schemas for all of these to validate them!
type ZSSProcessor struct {
	Type string // "transform", "validator", "required"

	// Validator
	ID        string // issue code or transform ID
	IssuePath *string
	Message   *string
	Params    map[string]any
}

type ZSSRequired struct {
	Type    string // "required"
	ID      zconst.ZogIssueCode
	Message string
	// optional
	IssuePath *string
	Params    map[string]any
}

type ZSSTest struct {
	ID        zconst.ZogIssueCode // issue code
	Message   string
	IssuePath *string
	Params    map[string]any
}

type ZSSTransformer struct {
	Type        string // "transformer"
	TransformId *string
}

type ZSSSchema struct {
	Kind         string  // "string", "number", "bool", "time", "slice", "struct", "ptr"
	Type         string  // Custom type if available (only if ZSS Exhaustive Metadata is enabled)
	Format       *string // Used for time.Time schemas only right now. (Only if ZSS Exhaustive Metadata is enabled)
	Processors   []any   // ZSSTest or ZSSTransformer
	Child        any     // *ZSSSchema | map[string]ZSSSchema
	Required     *ZSSRequired
	DefaultValue any
	CatchValue   any
}

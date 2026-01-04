package zss // Zog Schema Specification

// TODO make zog schemas for all of these to validate them!
type ZSSProcessor struct {
	Type string // "transform", "validator", "required"

	// Validator
	ID        string // issue code or transform ID
	IssuePath *string
	Message   *string
	Params    map[string]any
}

type ZSSTest struct {
	Type string // "test"
	// Validator
	IssueCode *string
	IssuePath *string
	Params    map[string]any
}

type ZSSTransformer struct {
	Type        string // "transformer"
	TransformId *string
}

type ZSSSchema struct {
	Type         string // "string"
	Processors   []any  // ZSSTest or ZSSTransformer
	Format       *string
	Child        any // *ZSSSchema | map[string]ZSSSchema
	Required     *ZSSTest
	DefaultValue any
	CatchValue   any
}

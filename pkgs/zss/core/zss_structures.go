package zsscore // Zog Schema Specification

import "github.com/Oudwins/zog/zconst" // TODO make zog schemas for all of these to validate them!

type ZSSDocument struct {
	Version ZSSVersion `json:"$schema"` // URL to ZSS Json Schema file(e.g., "https://zog.dev/zss/0.0.1/schema.json")
	Root    *ZSSSchema `json:"root"`
}
type ZSSProcessor struct {
	Kind        zconst.ZogProcessor `json:"kind"` // "transform", "validator"
	Test        *ZSSTest            `json:"test"`
	Transformer *ZSSTransformer     `json:"transformer"`
}

type ZSSTest struct {
	ID        zconst.ZogIssueCode `json:"id"` // issue code
	Message   string              `json:"message"`
	IssuePath []string            `json:"issuePath"`
	Params    map[string]any      `json:"params"`
}

type ZSSTransformer struct {
	ID zconst.ZogTransformID `json:"id"`
}

type ZSSGoType struct {
	PkgPath string `json:"pkgPath"` // Package path, empty for builtins
	Name    string `json:"name"`    // Type name, may be empty for unnamed types
	Display string `json:"display"` // Full type string (e.g., "*mypkg.User", "[]string")
}

type ZSSSchemaChildKind = string

const (
	ZSSSchemaChildKindShape  ZSSSchemaChildKind = "shape"
	ZSSSchemaChildKindSchema ZSSSchemaChildKind = "schema"
)

type ZSSSchemaChild struct {
	Kind   ZSSSchemaChildKind   `json:"kind"` // shape or schema
	Schema *ZSSSchema           `json:"schema,omitempty"`
	Shape  map[string]ZSSSchema `json:"shape,omitempty"`
}

type ZSSSchema struct {
	Kind         string           `json:"kind"`              // "string", "number", "bool", "time", "slice", "struct", "ptr"
	GoTypes      []ZSSGoType      `json:"goTypes,omitempty"` // Type metadata (only if ZSS Exhaustive Metadata is enabled)
	Format       *string          `json:"format"`            // Used for time.Time schemas only right now. (Only if ZSS Exhaustive Metadata is enabled)
	Processors   []ZSSProcessor   `json:"processors"`
	Childs       []ZSSSchemaChild `json:"childs"`
	Required     *ZSSTest         `json:"required"`
	DefaultValue any              `json:"defaultValue"`
	CatchValue   any              `json:"catchValue"`
}

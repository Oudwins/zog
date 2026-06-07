package zsscore // Zog Schema Specification

import (
	"strconv"

	"github.com/Oudwins/zog/zconst"
)

func ZSSRefFromKey(key int) string {
	return "#/$defs/" + ZSSDefKeyFromKey(key)
}

func ZSSDefKeyFromKey(key int) string {
	return "schema" + strconv.Itoa(key)
}

type ZSSDocument struct {
	URI  ZSSSchemaVersion      `json:"$schema"` // URL to ZSS Json Schema file(e.g., "https://zog.dev/zss/0.0.1/schema.json")
	Root *ZSSSchema            `json:"root"`
	Defs map[string]*ZSSSchema `json:"$defs,omitempty"`
}
type ZSSProcessor struct {
	Kind        zconst.ZogProcessor `json:"kind"` // "test", "transform"
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

type ZSSFieldMeta struct {
	Tags string `json:"tags,omitempty"`
}

type ZSSExtension struct {
	URI     string `json:"uri"`
	Content any    `json:"content,omitempty"`
}

type ZSSSchema struct {
	Ref          *string                 `json:"$ref,omitempty"`
	Kind         zconst.ZogType          `json:"kind,omitempty"` // one of the first party ZogTypes or a custom one with an extension
	Extension    *ZSSExtension           `json:"extension,omitempty"`
	GoTypes      []ZSSGoType             `json:"goTypes,omitempty"` // Type metadata
	Format       *string                 `json:"format,omitempty"`  // Used for time.Time schemas only right now. (Only if ZSS Exhaustive Metadata is enabled)
	Processors   []ZSSProcessor          `json:"processors,omitempty"`
	Fields       map[string]*ZSSSchema   `json:"fields,omitempty"`    // struct only
	FieldMeta    map[string]ZSSFieldMeta `json:"fieldMeta,omitempty"` // struct only
	Element      *ZSSSchema              `json:"element,omitempty"`   // ptr, slice, preprocess, boxed
	Key          *ZSSSchema              `json:"key,omitempty"`       // map only
	Value        *ZSSSchema              `json:"value,omitempty"`     // map only
	Required     *ZSSTest                `json:"required,omitempty"`
	DefaultValue any                     `json:"defaultValue,omitempty"`
	CatchValue   any                     `json:"catchValue,omitempty"`
}

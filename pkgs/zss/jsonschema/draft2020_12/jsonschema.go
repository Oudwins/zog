package draft2020_12

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/pkgs/zss/jsonschema/shared"
	"github.com/Oudwins/zog/pkgs/zss/regexconv"
	"github.com/Oudwins/zog/zconst"
)

const draft = "https://json-schema.org/draft/2020-12/schema"

type Schema = shared.Schema

type UnknownKindConverter = shared.UnknownKindConverter

type TestConverter = shared.TestConverter

type Options struct {
	UnknownKindConverter UnknownKindConverter
	TestConverter        TestConverter
}

type converter struct {
	opts Options
}

func FromZSS(doc zsscore.ZSSDocument, opts ...Options) (Schema, error) {
	c := converter{opts: options(opts)}
	root, err := c.convertSchema(doc.Root)
	if err != nil {
		return nil, err
	}
	root["$schema"] = draft

	if len(doc.Defs) > 0 {
		defs := Schema{}
		for key, def := range doc.Defs {
			converted, err := c.convertSchema(def)
			if err != nil {
				return nil, fmt.Errorf("convert $defs.%s: %w", key, err)
			}
			defs[key] = converted
		}
		root["$defs"] = defs
	}

	return root, nil
}

func options(opts []Options) Options {
	if len(opts) == 0 {
		return Options{TestConverter: ConvertTest}
	}
	opt := opts[0]
	if opt.TestConverter == nil {
		opt.TestConverter = ConvertTest
	}
	return opt
}

func (c converter) convertSchema(schema *zsscore.ZSSSchema) (Schema, error) {
	if schema == nil {
		return Schema{}, nil
	}
	if schema.Ref != nil {
		return Schema{"$ref": *schema.Ref}, nil
	}

	var out Schema
	var err error
	switch schema.Kind {
	case zconst.TypeString:
		out = Schema{"type": "string"}
	case zconst.TypeNumber:
		out = Schema{"type": "number"}
	case zconst.TypeBool:
		out = Schema{"type": "boolean"}
	case zconst.TypeTime:
		// TODO: map ZSS Format/Go time layouts more precisely than date-time.
		out = Schema{"type": "string", "format": "date-time"}
	case zconst.TypeSlice:
		out, err = c.convertSlice(schema)
	case zconst.TypeMap:
		out, err = c.convertMap(schema)
	case zconst.TypeStruct:
		out, err = c.convertStruct(schema)
	case zconst.TypePtr:
		out, err = c.convertPtr(schema)
	case zconst.TypePreprocess, zconst.TypeBoxed:
		out, err = c.convertSchema(schema.Element)
	case zconst.TypeAny, zconst.TypeCustom:
		out = Schema{}
	case "":
		out = Schema{}
	default:
		if c.opts.UnknownKindConverter != nil {
			out, err = c.opts.UnknownKindConverter(schema)
			break
		}
		return nil, fmt.Errorf("unsupported zss kind %q", schema.Kind)
	}
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = Schema{}
	}

	if err := c.applyProcessors(out, schema); err != nil {
		return nil, err
	}
	return out, nil
}

func (c converter) convertSlice(schema *zsscore.ZSSSchema) (Schema, error) {
	items, err := c.convertSchema(schema.Element)
	if err != nil {
		return nil, err
	}
	return Schema{"type": "array", "items": items}, nil
}

func (c converter) convertMap(schema *zsscore.ZSSSchema) (Schema, error) {
	value, err := c.convertSchema(schema.Value)
	if err != nil {
		return nil, err
	}
	return Schema{"type": "object", "additionalProperties": value}, nil
}

func (c converter) convertStruct(schema *zsscore.ZSSSchema) (Schema, error) {
	properties := Schema{}
	required := []string{}
	for name, field := range schema.Fields {
		propertyName, ok := propertyName(schema.FieldMeta[name], name)
		if !ok {
			continue
		}
		converted, err := c.convertSchema(field)
		if err != nil {
			return nil, fmt.Errorf("convert field %s: %w", name, err)
		}
		properties[propertyName] = converted
		if field != nil && field.Required != nil {
			required = append(required, propertyName)
		}
	}

	out := Schema{"type": "object", "properties": properties}
	if len(required) > 0 {
		sort.Strings(required)
		out["required"] = required
	}
	return out, nil
}

func propertyName(meta zsscore.ZSSFieldMeta, fallback string) (string, bool) {
	tags := reflect.StructTag(meta.Tags)
	if jsonTag, ok := tags.Lookup("json"); ok {
		name, _, _ := strings.Cut(jsonTag, ",")
		if name == "-" {
			return "", false
		}
		if name != "" {
			return name, true
		}
	}
	if zogTag, ok := tags.Lookup(zconst.ZogTag); ok && zogTag != "" {
		return zogTag, true
	}
	return fallback, true
}

func (c converter) convertPtr(schema *zsscore.ZSSSchema) (Schema, error) {
	inner, err := c.convertSchema(schema.Element)
	if err != nil {
		return nil, err
	}
	if schema.Required != nil {
		return inner, nil
	}
	return nullable(inner), nil
}

func nullable(schema Schema) Schema {
	typeValue, ok := schema["type"]
	if !ok {
		return Schema{"anyOf": []any{schema, Schema{"type": "null"}}}
	}
	switch typed := typeValue.(type) {
	case string:
		copy := clone(schema)
		copy["type"] = []string{typed, "null"}
		return copy
	case []string:
		copy := clone(schema)
		copy["type"] = append(append([]string{}, typed...), "null")
		return copy
	default:
		return Schema{"anyOf": []any{schema, Schema{"type": "null"}}}
	}
}

func clone(schema Schema) Schema {
	copy := Schema{}
	for key, value := range schema {
		copy[key] = value
	}
	return copy
}

func (c converter) applyProcessors(out Schema, schema *zsscore.ZSSSchema) error {
	for _, processor := range schema.Processors {
		if processor.Test == nil {
			continue
		}
		if err := c.opts.TestConverter(out, schema.Kind, processor.Test); err != nil {
			return fmt.Errorf("convert test %q: %w", processor.Test.ID, err)
		}
	}
	return nil
}

func ConvertTest(out Schema, kind zconst.ZogType, test *zsscore.ZSSTest) error {
	switch test.ID {
	case zconst.IssueCodeMin:
		applyMin(out, kind, test.Params[zconst.IssueCodeMin])
	case zconst.IssueCodeMax:
		applyMax(out, kind, test.Params[zconst.IssueCodeMax])
	case zconst.IssueCodeLen:
		applyLen(out, kind, test.Params[zconst.IssueCodeLen])
	case zconst.IssueCodeEQ:
		out["const"] = test.Params[zconst.IssueCodeEQ]
	case zconst.IssueCodeOneOf:
		out["enum"] = test.Params[zconst.IssueCodeOneOf]
	case zconst.IssueCodeGT:
		out["exclusiveMinimum"] = test.Params[zconst.IssueCodeGT]
	case zconst.IssueCodeGTE:
		out["minimum"] = test.Params[zconst.IssueCodeGTE]
	case zconst.IssueCodeLT:
		out["exclusiveMaximum"] = test.Params[zconst.IssueCodeLT]
	case zconst.IssueCodeLTE:
		out["maximum"] = test.Params[zconst.IssueCodeLTE]
	case zconst.IssueCodeEmail:
		out["format"] = "email"
	case zconst.IssueCodeUUID:
		out["format"] = "uuid"
	case zconst.IssueCodeURL:
		out["format"] = "uri"
	case zconst.IssueCodeIP:
		out["format"] = "ip"
	case zconst.IssueCodeMatch:
		pattern, ok := test.Params[zconst.IssueCodeMatch].(string)
		if !ok {
			return fmt.Errorf("match test param must be a string")
		}
		out["pattern"] = regexconv.GoRE2ToECMA262(pattern)
	case zconst.IssueCodeTrue:
		out["const"] = true
	case zconst.IssueCodeFalse:
		out["const"] = false
	}
	return nil
}

func applyMin(out Schema, kind zconst.ZogType, value any) {
	switch kind {
	case zconst.TypeString:
		out["minLength"] = value
	case zconst.TypeSlice:
		out["minItems"] = value
	default:
		out["minimum"] = value
	}
}

func applyMax(out Schema, kind zconst.ZogType, value any) {
	switch kind {
	case zconst.TypeString:
		out["maxLength"] = value
	case zconst.TypeSlice:
		out["maxItems"] = value
	default:
		out["maximum"] = value
	}
}

func applyLen(out Schema, kind zconst.ZogType, value any) {
	switch kind {
	case zconst.TypeString:
		out["minLength"] = value
		out["maxLength"] = value
	case zconst.TypeSlice:
		out["minItems"] = value
		out["maxItems"] = value
	}
}

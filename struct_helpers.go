package zog

import (
	"maps"

	p "github.com/Oudwins/zog/internals"
)

// Merges two or more structSchemas into a new structSchema
// PLEASE NOTE: This is a shallow merge. Meaning that:
// - if you merge multiple schemas with the same key, the last one will take precedence
// - preTransforms, postTransforms and tests are merged, but the order is not guaranteed and if you have multiple of the same type they will be duplicated
// - if you modify a nested part of the schema it may effect the other schemas
func (v *structProcessor) Merge(other *structProcessor, others ...*structProcessor) *structProcessor {
	new := &structProcessor{
		preTransforms:  make([]p.PreTransform, 0),
		postTransforms: make([]p.PostTransform, 0),
		tests:          make([]p.Test, 0),
		required:       other.required,
		schema:         Schema{},
	}
	if v.preTransforms != nil {
		new.preTransforms = append(new.preTransforms, v.preTransforms...)
	}
	if other.preTransforms != nil {
		new.preTransforms = append(new.preTransforms, other.preTransforms...)
	}

	if v.postTransforms != nil {
		new.postTransforms = append(new.postTransforms, v.postTransforms...)
	}
	if other.postTransforms != nil {
		new.postTransforms = append(new.postTransforms, other.postTransforms...)
	}

	if v.tests != nil {
		new.tests = append(new.tests, v.tests...)
	}
	if other.tests != nil {
		new.tests = append(new.tests, other.tests...)
	}

	maps.Copy(new.schema, v.schema)
	maps.Copy(new.schema, other.schema)
	for _, o := range others {
		new = new.Merge(o)
	}
	return new
}

func (v *structProcessor) cloneShallow() *structProcessor {
	new := &structProcessor{
		preTransforms:  v.preTransforms,
		postTransforms: v.postTransforms,
		tests:          v.tests,
		required:       v.required,
		schema:         v.schema,
	}
	return new
}

func (v *structProcessor) Omit(vals ...any) *structProcessor {
	new := v.cloneShallow()
	new.schema = Schema{}
	maps.Copy(new.schema, v.schema)
	for _, k := range vals {
		switch k := k.(type) {
		case string:
			delete(new.schema, k)
		case map[string]bool:
			for key, val := range k {
				if val {
					delete(new.schema, key)
				}
			}
		}
	}
	return new
}

func (v *structProcessor) Pick(picks ...any) *structProcessor {
	new := v.cloneShallow()
	new.schema = Schema{}
	for _, pick := range picks {
		switch pick := pick.(type) {
		case string:
			new.schema[pick] = v.schema[pick]
		case map[string]bool:
			for k, pick := range pick {
				if pick {
					new.schema[k] = v.schema[k]
				}
			}
		}
	}
	return new
}

func (v *structProcessor) Extend(schema Schema) *structProcessor {
	new := v.cloneShallow()
	new.schema = Schema{}
	maps.Copy(new.schema, v.schema)
	maps.Copy(new.schema, schema)
	return new
}

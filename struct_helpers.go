package zog

import (
	"maps"

	p "github.com/Oudwins/zog/internals"
)

// Merge combines two or more schemas into a new schema.
// It performs a shallow merge, meaning:
//   - Fields with the same key from later schemas override earlier ones
//   - PreTransforms, PostTransforms and tests are concatenated in order
//   - Modifying nested schemas may affect the original schemas
//
// Parameters:
//   - other: The first schema to merge with
//   - others: Additional schemas to merge
//
// Returns a new schema containing the merged fields and transforms
func (v *StructSchema) Merge(other *StructSchema, others ...*StructSchema) *StructSchema {
	totalProcessors := len(v.processors) + len(other.processors)
	for _, o := range others {
		totalProcessors += len(o.processors)
	}
	new := &StructSchema{
		processors: make([]p.ZProcessor[any], 0, totalProcessors),
		required:   other.required,
		schema:     Shape{},
	}

	// processors
	new.processors = append(new.processors, v.processors...)
	new.processors = append(new.processors, other.processors...)
	for _, s := range others {
		new.processors = append(new.processors, s.processors...)
	}

	maps.Copy(new.schema, v.schema)
	maps.Copy(new.schema, other.schema)
	for _, s := range others {
		maps.Copy(new.schema, s.schema)
	}

	return new
}

// cloneShallow creates a shallow copy of the schema.
// The new schema shares references to the transforms, tests and inner schema.
func (v *StructSchema) cloneShallow() *StructSchema {
	new := &StructSchema{
		processors: v.processors,
		required:   v.required,
		schema:     v.schema,
	}
	return new
}

// Omit creates a new schema with specified fields removed.
// It accepts either strings or map[string]bool as arguments:
//   - Strings directly specify fields to omit
//   - For maps, fields are omitted when their boolean value is true
//
// Returns a new schema with the specified fields removed
func (v *StructSchema) Omit(vals ...any) *StructSchema {
	new := v.cloneShallow()
	new.schema = Shape{}
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

// Pick creates a new schema keeping only the specified fields.
// It accepts either strings or map[string]bool as arguments:
//   - Strings directly specify fields to keep
//   - For maps, fields are kept when their boolean value is true
//
// Returns a new schema containing only the specified fields
func (v *StructSchema) Pick(picks ...any) *StructSchema {
	new := v.cloneShallow()
	new.schema = Shape{}
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

// Extend creates a new schema by adding additional fields from the provided schema.
// Fields in the provided schema override any existing fields with the same key.
//
// Parameters:
//   - schema: The schema containing fields to add
//
// Returns a new schema with the additional fields
func (v *StructSchema) Extend(schema Shape) *StructSchema {
	new := v.cloneShallow()
	new.schema = Shape{}
	maps.Copy(new.schema, v.schema)
	maps.Copy(new.schema, schema)
	return new
}

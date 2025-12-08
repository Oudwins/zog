package ast

import "github.com/Oudwins/zog/zconst"

type Required struct {
	Bool bool
}

type Catch struct {
	Bool  bool
	Value any
}

type Test struct {
	Code   zconst.ZogIssueCode
	Params map[string]any
}

type Schema struct {
	Type    zconst.ZogType
	Shape   *map[string]*Schema
	Element *Schema
	// isRequired -> bool or we could do {Bool: true|false, Condition: something_more_complex}
	// catch -> {Bool: true|false, value: catch_value}
	// coercer -> we could define here the valid coercer values also somehow

	// preTransforms -> {id: "trim" | "custom" , content: string_of_code}
	// tests -> {code: "is_email" | "custom", params: map[string]any, content: string_of_code}
	// postTransform

}

package primitives

type RuleValidateFunc func(Rule) bool

type Rule struct {
	Name         string
	RuleValue    any
	FieldValue   any // TODO I think I can remove this
	FieldName    any // TODO I think I can remove this
	ErrorMessage string
	ValidateFunc RuleValidateFunc
}

// ORIGINAL
// type RuleSet struct {
// 	Name         string
// 	RuleValue    any
// 	FieldValue   any
// 	FieldName    any
// 	ErrorMessage string
// 	MessageFunc  func(RuleSet) string
// 	ValidateFunc func(RuleSet) bool
// }

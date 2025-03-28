package internals

import (
	"fmt"

	zconst "github.com/Oudwins/zog/zconst"
)

// this is the function that formats the error message given a zog error
type IssueFmtFunc = func(e *ZogIssue, p Ctx)

// ZogIssue represents an issue that occurred during parsing or validation.
// When printed it looks like:
// ZogIssue{Code: coercion_issue, Params: map[], Type: number, Value: not_empty, Message: number is invalid, Error: failed to coerce string int: strconv.Atoi: parsing "not_empty": invalid syntax}
type ZogIssue struct {
	// Code is the unique identifier for the issue. Generally also the ID for the Test that caused the issue.
	Code zconst.ZogIssueCode
	// Path is the path to the field that caused the issue
	Path string
	// Value is the data value that caused the issue.
	// If using Schema.Parse(data, dest) then this will be the value of data.
	Value any
	// Dtype is the destination type. i.e The zconst.ZogType of the value that was validated.
	// If using Schema.Parse(data, dest) then this will be the type of dest.
	Dtype string
	// Params is the params map for the issue. Taken from the Test that caused the issue.
	// This may be nil if Test has no params.
	Params map[string]any
	// Message is the human readable, user-friendly message for the issue.
	// This is safe to expose to the user.
	Message string
	// Err is the wrapped error or nil if none
	Err error
}

func NewZogIssue() *ZogIssue {
	e := ZogIssuePool.Get().(*ZogIssue)
	e.Code = ""
	e.Path = ""
	e.Value = nil
	e.Dtype = ""
	e.Params = nil
	e.Message = ""
	e.Err = nil
	return e
}

// SetCode sets the issue code for the issue and returns the issue for chaining
func (i *ZogIssue) SetCode(c zconst.ZogIssueCode) *ZogIssue {
	i.Code = c
	return i
}

// SetPath sets the path for the issue and returns the issue for chaining
func (i *ZogIssue) SetPath(p string) *ZogIssue {
	i.Path = p
	return i
}

// SetValue sets the data value that caused the issue and returns the issue for chaining
func (i *ZogIssue) SetValue(v any) *ZogIssue {
	i.Value = v
	return i
}

// SetDType sets the destination type for the issue and returns the issue for chaining
func (i *ZogIssue) SetDType(t string) *ZogIssue {
	i.Dtype = t
	return i
}

// SetParams sets the params map for the issue and returns the issue for chaining
func (i *ZogIssue) SetParams(p map[string]any) *ZogIssue {
	i.Params = p
	return i
}

// SetMessage sets the human readable, user-friendly message for the issue and returns the issue for chaining
func (i *ZogIssue) SetMessage(m string) *ZogIssue {
	i.Message = m
	return i
}

// SetError sets the wrapped error for the issue and returns the issue for chaining
func (i *ZogIssue) SetError(e error) *ZogIssue {
	i.Err = e
	return i
}

// Unwrap returns the wrapped error or nil if none
func (i *ZogIssue) Unwrap() error {
	return i.Err
}

// Error returns the string representation of the ZogIssue (same as String())
func (i *ZogIssue) Error() string {
	return i.String()
}

// String returns the string representation of the ZogIssue (same as Error())
func (i *ZogIssue) String() string {
	return fmt.Sprintf("ZogIssue{Code: %v, Params: %v, Type: %v, Value: %v, Message: '%v', Error: %v}", SafeString(i.Code), SafeString(i.Params), SafeString(i.Dtype), SafeString(i.Value), SafeString(i.Message), SafeError(i.Err))
}

func FreeIssue(i *ZogIssue) {
	ZogIssuePool.Put(i)
}

// list of errors. This is returned by processors for simple types (e.g. strings, numbers, booleans)
type ZogIssueList = []*ZogIssue

// map of errors. This is returned by processors for complex types (e.g. maps, slices, structs)
type ZogIssueMap = map[string]ZogIssueList

// INTERNAL ONLY: Interface used to add errors during parsing & validation. It represents a group of errors (map or slice)
type ZogIssues interface {
	Add(path string, err *ZogIssue)
	IsEmpty() bool
	Free()
}

// internal only
type ErrsList struct {
	List ZogIssueList
}

// internal only
func NewErrsList() *ErrsList {
	l := InternalIssueListPool.Get().(*ErrsList)
	l.List = nil
	return l
}

func (e *ErrsList) Add(path string, err *ZogIssue) {
	if e.List == nil {
		e.List = make(ZogIssueList, 0, 2)
	}
	e.List = append(e.List, err)
}

func (e *ErrsList) IsEmpty() bool {
	return e.List == nil
}

func (e *ErrsList) Free() {
	InternalIssueListPool.Put(e)
}

// map implementation of Errs
type ErrsMap struct {
	M ZogIssueMap
}

// Factory for errsMap
func NewErrsMap() *ErrsMap {
	m := InternalIssueMapPool.Get().(*ErrsMap)
	m.M = nil
	return m
}

func (s *ErrsMap) Add(p string, err *ZogIssue) {
	// checking if its the first error
	if s.M == nil {
		s.M = ZogIssueMap{}
		s.M[zconst.ISSUE_KEY_FIRST] = []*ZogIssue{err}
	}

	path := p
	if path == "" {
		path = zconst.ISSUE_KEY_ROOT
	}
	if _, ok := s.M[path]; !ok {
		s.M[path] = []*ZogIssue{}
	}
	s.M[path] = append(s.M[path], err)
}

func (s *ErrsMap) IsEmpty() bool {
	return s.M == nil
}

func (s *ErrsMap) Free() {
	InternalIssueMapPool.Put(s)
}

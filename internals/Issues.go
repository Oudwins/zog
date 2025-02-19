package internals

import (
	"fmt"

	zconst "github.com/Oudwins/zog/zconst"
)

type ZogIssue interface {
	// returns the error code for the error. This is a unique identifier for the Issue. Generally also the ID for the Test that caused the error.
	Code() zconst.ZogIssueCode

	// Sets the error code for the Issue. This is a unique identifier for the Issue. Generally also the ID for the Test that caused the error.
	SetCode(zconst.ZogIssueCode) ZogIssue

	// returns the path of the Issue. This is the path of the value that caused the Issue.
	Path() string

	// Sets the path of the Issue. This is the path of the value that caused the Issue.
	SetPath(string) ZogIssue

	// returns the data value that caused the Issue.
	// if using Schema.Parse(data, dest) then this will be the value of data.
	Value() any

	// Deprecated: Use SetValue() instead
	// Sets the data value that caused the Issue.
	// if using Schema.Parse(data, dest) then this will be the value of data.
	SValue(any) ZogIssue

	// Sets the data value that caused the Issue.
	// if using Schema.Parse(data, dest) then this will be the value of data.
	SetValue(any) ZogIssue

	// Returns destination type. i.e The zconst.ZogType of the value that was validated.
	// if Using Schema.Parse(data, dest) then this will be the type of dest.
	Dtype() string

	// Deprecated: Use SetDType() instead
	// Sets destination type. i.e The zconst.ZogType of the value that was validated.
	// if Using Schema.Parse(data, dest) then this will be the type of dest.
	SDType(zconst.ZogType) ZogIssue

	// Sets destination type. i.e The zconst.ZogType of the value that was validated.
	// if Using Schema.Parse(data, dest) then this will be the type of dest.
	SetDType(zconst.ZogType) ZogIssue

	// returns the params map for the error. Taken from the Test that caused the error. This may be nil if Test has no params.
	Params() map[string]any

	// Deprecated: Use SetParams() instead
	// Sets the params map for the error. Taken from the Test that caused the error. This may be nil if Test has no params.
	SParams(map[string]any) ZogIssue
	// Sets the params map for the error. Taken from the Test that caused the error. This may be nil if Test has no params.
	SetParams(map[string]any) ZogIssue
	// returns the human readable, user-friendly message for the error. This is safe to expose to the user.
	Message() string
	// sets the human readable, user-friendly message for the error. This is safe to expose to the user.
	SetMessage(string)
	// returns the string representation of the ZogIssue (same as String())
	Error() string
	// Sets the wrapped error.
	SetError(error) ZogIssue
	// returns the wrapped error or nil if none
	Unwrap() error
	// returns the string representation of the ZogIssue (same as Error())
	String() string
}

// Error interface returned from all processors
type ZogError = ZogIssue

// this is the function that formats the error message given a zog error
type IssueFmtFunc = func(e ZogError, p Ctx)

// INTERNAL ONLY: Error implementation
type ZogErr struct {
	C       zconst.ZogIssueCode // error code
	EPath   string              // path of the value that caused the error
	ParamsM map[string]any      // params for the error (e.g. min, max, len, etc)
	Typ     string              // destination type
	Val     any                 // value that caused the error
	Msg     string              // human readable message
	Err     error               // the underlying error
}

// error code, err uuid
func (e *ZogErr) Code() zconst.ZogIssueCode {
	return e.C
}
func (e *ZogErr) SetCode(c zconst.ZogIssueCode) ZogError {
	e.C = c
	return e
}

func (e *ZogErr) Path() string {
	return e.EPath
}
func (e *ZogErr) SetPath(p string) ZogError {
	e.EPath = p
	return e
}

// value that caused the error
func (e *ZogErr) Value() any {
	return e.Val
}
func (e *ZogErr) SValue(v any) ZogError {
	e.Val = v
	return e
}

func (e *ZogErr) SetValue(v any) ZogError {
	e.Val = v
	return e
}

// destination type TODO
func (e *ZogErr) Dtype() string {
	return e.Typ
}
func (e *ZogErr) SDType(t zconst.ZogType) ZogError {
	e.Typ = t
	return e
}

func (e *ZogErr) SetDType(t zconst.ZogType) ZogError {
	e.Typ = t
	return e
}

func (e *ZogErr) Params() map[string]any {
	return e.ParamsM
}

func (e *ZogErr) SParams(p map[string]any) ZogError {
	e.ParamsM = p
	return e
}

func (e *ZogErr) SetParams(p map[string]any) ZogError {
	e.ParamsM = p
	return e
}
func (e *ZogErr) Message() string {
	return e.Msg
}
func (e *ZogErr) SetMessage(msg string) {
	e.Msg = msg
}
func (e *ZogErr) Error() string {
	return e.String()
}
func (e *ZogErr) SetError(err error) ZogError {
	e.Err = err
	return e
}
func (e *ZogErr) Unwrap() error {
	return e.Err
}

func (e *ZogErr) String() string {
	return fmt.Sprintf("ZogError{Code: %v, Params: %v, Type: %v, Value: %v, Message: '%v', Error: %v}", SafeString(e.C), SafeString(e.ParamsM), SafeString(e.Typ), SafeString(e.Val), SafeString(e.Msg), SafeError(e.Err))
}

func (e *ZogErr) Free() {
	ZogIssuePool.Put(e)
}

// list of errors. This is returned by processors for simple types (e.g. strings, numbers, booleans)
type ZogIssueList = []ZogError

// map of errors. This is returned by processors for complex types (e.g. maps, slices, structs)
type ZogIssueMap = map[string][]ZogError

// INTERNAL ONLY: Interface used to add errors during parsing & validation. It represents a group of errors (map or slice)
type ZogIssues interface {
	Add(path string, err ZogError)
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

func (e *ErrsList) Add(path string, err ZogError) {
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

func (s *ErrsMap) Add(p string, err ZogError) {
	// checking if its the first error
	if s.M == nil {
		s.M = ZogIssueMap{}
		s.M[zconst.ISSUE_KEY_FIRST] = []ZogError{err}
	}

	path := p
	if path == "" {
		path = zconst.ISSUE_KEY_ROOT
	}
	if _, ok := s.M[path]; !ok {
		s.M[path] = []ZogError{}
	}
	s.M[path] = append(s.M[path], err)
}

func (s *ErrsMap) IsEmpty() bool {
	return s.M == nil
}

func (s *ErrsMap) Free() {
	InternalIssueMapPool.Put(s)
}

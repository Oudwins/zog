package internals

import (
	"fmt"

	zconst "github.com/Oudwins/zog/zconst"
)

// Error interface returned from all processors
type ZogError interface {
	// returns the error code for the error. This is a unique identifier for the error. Generally also the ID for the Test that caused the error.
	Code() zconst.ZogErrCode
	// returns the data value that caused the error.
	// if using Schema.Parse(data, dest) then this will be the value of data.
	Value() any
	// Sets the data value that caused the error.
	// if using Schema.Parse(data, dest) then this will be the value of data.
	SValue(any) ZogError
	// Returns destination type. i.e The zconst.ZogType of the value that was validated.
	// if Using Schema.Parse(data, dest) then this will be the type of dest.
	Dtype() string
	// Sets destination type. i.e The zconst.ZogType of the value that was validated.
	// if Using Schema.Parse(data, dest) then this will be the type of dest.
	SDType(zconst.ZogType) ZogError
	// returns the params map for the error. Taken from the Test that caused the error. This may be nil if Test has no params.
	Params() map[string]any
	// Sets the params map for the error. Taken from the Test that caused the error. This may be nil if Test has no params.
	SParams(map[string]any) ZogError
	// returns the human readable, user-friendly message for the error. This is safe to expose to the user.
	Message() string
	// sets the human readable, user-friendly message for the error. This is safe to expose to the user.
	SetMessage(string)
	// returns the string representation of the ZogError (same as String())
	Error() string
	// returns the wrapped error or nil if none
	Unwrap() error
	// returns the string representation of the ZogError (same as Error())
	String() string
}

// this is the function that formats the error message given a zog error
type ErrFmtFunc = func(e ZogError, p ParseCtx)

// INTERNAL ONLY: Error implementation
type ZogErr struct {
	C       zconst.ZogErrCode // error code
	ParamsM map[string]any    // params for the error (e.g. min, max, len, etc)
	Typ     string            // destination type
	Val     any               // value that caused the error
	Msg     string
	Err     error // the underlying error
}

// error code, err uuid
func (e *ZogErr) Code() zconst.ZogErrCode {
	return e.C
}

// value that caused the error
func (e *ZogErr) Value() any {
	return e.Value
}
func (e *ZogErr) SValue(v any) ZogError {
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

func (e *ZogErr) Params() map[string]any {
	return e.ParamsM
}

func (e *ZogErr) SParams(p map[string]any) ZogError {
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
func (e *ZogErr) Unwrap() error {
	return e.Err
}

func (e *ZogErr) String() string {
	return fmt.Sprintf("ZogError{Code: %v, Params: %v, Type: %v, Value: %v, Message: '%v', Error: %v}", SafeString(e.C), SafeString(e.ParamsM), SafeString(e.Typ), SafeString(e.Val), SafeString(e.Msg), SafeError(e.Err))
}

// list of errors. This is returned by processors for simple types (e.g. strings, numbers, booleans)
type ZogErrList = []ZogError

// map of errors. This is returned by processors for complex types (e.g. maps, slices, structs)
type ZogErrMap = map[string][]ZogError

// INTERNAL ONLY: Interface used to add errors during parsing & validation. It represents a group of errors (map or slice)
type ZogErrors interface {
	Add(path PathBuilder, err ZogError)
	IsEmpty() bool
}

// internal only
type ErrsList struct {
	List ZogErrList
}

// internal only
func NewErrsList() *ErrsList {
	return &ErrsList{}
}

func (e *ErrsList) Add(path PathBuilder, err ZogError) {
	if e.List == nil {
		e.List = make(ZogErrList, 0, 2)
	}
	e.List = append(e.List, err)
}

func (e *ErrsList) IsEmpty() bool {
	return e.List == nil
}

// map implementation of Errs
type ErrsMap struct {
	M ZogErrMap
}

const (
	ERROR_KEY_FIRST = "$first"
	ERROR_KEY_ROOT  = "$root"
)

// Factory for errsMap
func NewErrsMap() *ErrsMap {
	return &ErrsMap{}
}

func (s *ErrsMap) Add(p PathBuilder, err ZogError) {
	// checking if its the first error
	if s.M == nil {
		s.M = ZogErrMap{}
		s.M[ERROR_KEY_FIRST] = []ZogError{err}
	}

	path := p.String()
	if path == "" {
		path = ERROR_KEY_ROOT
	}
	if _, ok := s.M[path]; !ok {
		s.M[path] = []ZogError{}
	}
	s.M[path] = append(s.M[path], err)
}

func (s ErrsMap) IsEmpty() bool {
	return s.M == nil
}

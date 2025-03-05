package internals

import (
	"fmt"

	zconst "github.com/Oudwins/zog/zconst"
)

// this is the function that formats the error message given a zog error
type IssueFmtFunc = func(e *ZogIssue, p Ctx)

type ZogIssue struct {
	Code    zconst.ZogIssueCode
	Path    string
	Value   any
	Dtype   string
	Params  map[string]any
	Message string
	Err     error
}

func (i *ZogIssue) SetCode(c zconst.ZogIssueCode) *ZogIssue {
	i.Code = c
	return i
}

func (i *ZogIssue) SetPath(p string) *ZogIssue {
	i.Path = p
	return i
}

func (i *ZogIssue) SetValue(v any) *ZogIssue {
	i.Value = v
	return i
}

func (i *ZogIssue) SetDType(t string) *ZogIssue {
	i.Dtype = t
	return i
}

func (i *ZogIssue) SetParams(p map[string]any) *ZogIssue {
	i.Params = p
	return i
}

func (i *ZogIssue) SetMessage(m string) *ZogIssue {
	i.Message = m
	return i
}

func (i *ZogIssue) SetError(e error) *ZogIssue {
	i.Err = e
	return i
}

func (i *ZogIssue) Unwrap() error {
	return i.Err
}
func (i *ZogIssue) Error() string {
	return i.String()
}

func (i *ZogIssue) String() string {
	return fmt.Sprintf("ZogIssue{Code: %v, Params: %v, Type: %v, Value: %v, Message: '%v', Error: %v}", SafeString(i.Code), SafeString(i.Params), SafeString(i.Dtype), SafeString(i.Value), SafeString(i.Message), SafeError(i.Err))
}

func (i *ZogIssue) Free() {
	ZogIssuePool.Put(i)
}

// list of errors. This is returned by processors for simple types (e.g. strings, numbers, booleans)
type ZogIssueList = []*ZogIssue

// map of errors. This is returned by processors for complex types (e.g. maps, slices, structs)
type ZogIssueMap = map[string][]*ZogIssue

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

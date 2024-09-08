package primitives

// Error interface returned from all processors
type ZogError interface {
	Code() ZogErrCode
	Value() any    // the error value
	Dtype() string // destination type
	Params() map[string]any
	Message() string
	SetMessage(string)
	// returns the string of the wrapped error
	Error() string
	// returns the wrapped error
	Unwrap() error
}

// this is the function that formats the error message given a zog error
type ErrFmtFunc = func(e ZogError, p ParseCtx)

// Error implementation
type ZogErr struct {
	C       ZogErrCode     // error code
	ParamsM map[string]any // params for the error (e.g. min, max, len, etc)
	Typ     string         // destination type
	Val     any            // value that caused the error
	Msg     string
	Err     error // the underlying error
}

// error code, err uuid
func (e *ZogErr) Code() ZogErrCode {
	return e.C
}

// value that caused the error
func (e *ZogErr) Value() any {
	return e.Value
}

// destination type TODO
func (e *ZogErr) Dtype() string {
	return e.Typ
}

func (e *ZogErr) Params() map[string]any {
	return e.ParamsM
}
func (e *ZogErr) Message() string {
	return e.Msg
}
func (e *ZogErr) SetMessage(msg string) {
	e.Msg = msg
}
func (e *ZogErr) Error() string {
	return e.Err.Error()
}
func (e *ZogErr) Unwrap() error {
	return e.Err
}

// list of errors. This is returned for each specific processor
type ZogErrList = []ZogError
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

func (s ErrsMap) First() ZogError {
	if s.IsEmpty() {
		return nil
	}
	return s.M[ERROR_KEY_FIRST][0]
}

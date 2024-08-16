package primitives

// ZogError hides the error and only exposes the message
type ZogError struct {
	Message string
	Err     error
}

func (e ZogError) Error() string {
	return e.Message
}

func (e ZogError) Unwrap() error {
	return e.Err
}

// list of errors. This is returned for each specific processor
type ZogErrList = []ZogError
type ZogErrMap = map[string][]ZogError

// Interface used to add errors during parsing & validation
type ZogErrors interface {
	Add(path PathBuilder, err ZogError)
	IsEmpty() bool
}

type ErrsList struct {
	List ZogErrList
}

func NewErrsList() *ErrsList {
	return &ErrsList{}
}

func (e *ErrsList) Add(path PathBuilder, err ZogError) {
	if e.List == nil {
		e.List = make(ZogErrList, 0, 3)
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

func (s ErrsMap) First() error {
	if s.IsEmpty() {
		return nil
	}
	return s.M[ERROR_KEY_FIRST][0]
}

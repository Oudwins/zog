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
type ZogErrorList = []ZogError
type ZogSchemaErrors = map[string][]ZogError

// Interface used to add errors during parsing & validation
type ZogErrors interface {
	Add(path Pather, err ZogError)
	IsEmpty() bool
}

type ErrsList struct {
	List ZogErrorList
}

func NewErrsList() *ErrsList {
	return &ErrsList{}
}

func (e *ErrsList) Add(path Pather, err ZogError) {
	if e.List == nil {
		e.List = make(ZogErrorList, 0, 3)
	}
	e.List = append(e.List, err)
}

func (e *ErrsList) IsEmpty() bool {
	return e.List == nil
}

// map implementation of Errs
type ErrsMap struct {
	M ZogSchemaErrors
}

const (
	FIRST_ERROR_KEY = "$first"
)

// Factory for errsMap
func NewErrsMap() *ErrsMap {
	return &ErrsMap{}
}

func (s *ErrsMap) Add(p Pather, err ZogError) {
	// checking if its the first error
	if s.M == nil {
		s.M = ZogSchemaErrors{}
		s.M[FIRST_ERROR_KEY] = []ZogError{err}
	}

	path := p.String()
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
	return s.M[FIRST_ERROR_KEY][0]
}

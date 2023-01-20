package executor

type (
	ErrValidateQuery struct {
		err error
	}
	ErrExecuteQuery struct {
		err error
	}
)

func (e ErrValidateQuery) Error() string {
	return "validate query: " + e.err.Error()
}

func (e ErrValidateQuery) Unwrap() error {
	return e.err
}

func (e ErrValidateQuery) Is(err error) bool {
	_, ok := err.(ErrValidateQuery)
	return ok
}

func (e ErrExecuteQuery) Error() string {
	return "execute query: " + e.err.Error()
}

func (e ErrExecuteQuery) Unwrap() error {
	return e.err
}

func (e ErrExecuteQuery) Is(err error) bool {
	_, ok := err.(ErrExecuteQuery)
	return ok
}

func NewErrValidateQuery(err error) error {
	return ErrValidateQuery{err: err}
}

func NewErrExecuteQuery(err error) error {
	return ErrExecuteQuery{err: err}
}

package executor

type (
	ValidateQueryError struct {
		err error
	}
	ExecuteQueryError struct {
		err error
	}
)

func (e ValidateQueryError) Error() string {
	return "validate query: " + e.err.Error()
}

func (e ValidateQueryError) Unwrap() error {
	return e.err
}

func (e ValidateQueryError) Is(err error) bool {
	_, ok := err.(ValidateQueryError) //nolint:errorlint
	return ok
}

func (e ExecuteQueryError) Error() string {
	return "execute query: " + e.err.Error()
}

func (e ExecuteQueryError) Unwrap() error {
	return e.err
}

func (e ExecuteQueryError) Is(err error) bool {
	_, ok := err.(ExecuteQueryError) //nolint:errorlint
	return ok
}

func NewValidateQueryError(err error) error {
	return ValidateQueryError{err: err}
}

func NewExecuteQueryError(err error) error {
	return ExecuteQueryError{err: err}
}

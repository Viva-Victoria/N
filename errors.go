package n

type ErrorWrapper interface {
	Unwrap() error
}

type BadRequestError struct {
	err error
}

func NewBadRequestError(err error) BadRequestError {
	return BadRequestError{
		err: err,
	}
}

func (b BadRequestError) Error() string {
	if b.err == nil {
		return ""
	}

	return b.err.Error()
}

func (b BadRequestError) Unwrap() error {
	return b.err
}

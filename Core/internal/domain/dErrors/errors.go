package dErrors

type Error struct {
	message string
	err     error
}

func (e *Error) Error() string {
	return e.message
}

func NewError(message string) *Error {
	err := Error{
		message: message,
	}

	return &err
}

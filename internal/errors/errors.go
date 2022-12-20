package errors

type ErrorResponse struct {
	Message string `json:"message"`
}

type Error struct {
	Message    string
	StatusCode int
}

func New(message string, statusCode int) *Error {
	return &Error{Message: message, StatusCode: statusCode}
}

func (e *Error) Error() string {
	return e.Message
}

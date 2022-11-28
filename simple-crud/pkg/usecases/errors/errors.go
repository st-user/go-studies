package errors

import "fmt"

type UnexpectedError struct {
	errorSource string
	err         error
}

func NewUnexpectedError(errorSource string, err error) UnexpectedError {
	return UnexpectedError{
		errorSource: errorSource,
		err:         err,
	}
}

func (u UnexpectedError) Error() string {
	return fmt.Sprintf("encountering an unexpected error at %s: %s", u.errorSource, u.err.Error())
}

type DataNotFoundError struct {
	internalMessage string
	userMessage     string
}

func NewDataNotFoundErrorWithLabel(label string) DataNotFoundError {
	msg := fmt.Sprintf("%s Not Found", label)
	return DataNotFoundError{
		internalMessage: msg,
		userMessage:     msg,
	}
}

func (d DataNotFoundError) UserMessage() string {
	return d.userMessage
}

func (d DataNotFoundError) Error() string {
	return d.internalMessage
}

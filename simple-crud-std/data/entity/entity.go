package entity

import "time"

type Employee struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
}

type NotFoundError struct {
	msg string
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{msg}
}

func (e *NotFoundError) Error() string {
	return e.msg
}

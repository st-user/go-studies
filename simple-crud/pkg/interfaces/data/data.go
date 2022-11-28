package data

import (
	"fmt"
	"time"
)

type SimpleDate struct {
	time.Time
}

func (d SimpleDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

func (d *SimpleDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	parsedTime, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.UTC)
	if err != nil {
		return err
	}

	*d = SimpleDate{parsedTime}
	return err
}

type InputDataError struct {
	userMessage string
	detail      error
}

func NewInputDataError(userMessage string, detail error) InputDataError {
	return InputDataError{
		userMessage: userMessage,
		detail:      detail,
	}
}

func (i InputDataError) UserMessage() string {
	return i.userMessage
}

func (i InputDataError) HasDetail() bool {
	return i.detail != nil
}

func (i InputDataError) Error() string {
	return fmt.Sprintf("input data error %s", i.detail)
}

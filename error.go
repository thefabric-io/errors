package errors

import (
	"encoding/json"
)

const (
	CodeSubject = "subject"
	CodeInvalid = "invalid"
	CodeUnknown = "unknown"
)

type Code string

func (c Code) Equal(c2 Code) bool {
	return c == c2
}

type Message string

func (m Message) Equal(m2 Message) bool {
	return m == m2
}

func New(message Message, code Code) *Error {
	if code == "" {
		code = CodeUnknown
	}

	return &Error{message: message, code: code}
}

type Error struct {
	message Message
	code    Code
}

func (e Error) Code() Code {
	return e.code
}

func (e Error) Message() Message {
	return e.message
}

func (e Error) Error() string {
	b, _ := e.MarshalJSON()

	return string(b)
}

type mError struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (e Error) MarshalJSON() ([]byte, error) {
	r := mError{
		Message: string(e.message),
		Code:    string(e.code),
	}

	return json.Marshal(r)
}

func (e *Error) UnmarshalJSON(b []byte) error {
	var mErr mError
	if err := json.Unmarshal(b, &mErr); err != nil {
		return err
	}

	e.code = Code(mErr.Code)
	e.message = Message(mErr.Message)

	return nil
}

package errors

import (
	"encoding/json"
)

func Stack(source error, targets ...error) (result error) {
	switch e := source.(type) {
	case nil:
		if len(targets) != 0 {
			source = targets[0]
			targets = targets[1:]

			return Stack(source, targets...)
		}

	case *Error:
		result = newErrors(e.message, e.code)
		break

	case *Errors:
		result = e
		break

	default:
		result = newErrors(Message(source.Error()), CodeUnknown)
	}

	for _, target := range targets {
		result = result.(*Errors).Stack(target)
	}

	return
}

func IsWithStrategy(err1 error, err2 error, comparable Comparable) bool {
	return comparable.Compare(err1, err2)
}

func Is(err1 error, err2 error) bool {
	return IsWithStrategy(err1, err2, CompareMessageOnlyStrategy())
}

func newErrors(message Message, code Code) *Errors {
	errs := Errors{
		stacks: make([]Error, 1),
	}

	errs.stacks[0] = *New(message, code)

	return &errs
}

type Errors struct {
	stacks []Error
}

func (e *Errors) Stack(target error) error {
	switch err := target.(type) {
	case nil:
		break

	case *Error:
		e.stacks = append(e.stacks, *err)
		break

	case *Errors:
		for _, j := range err.stacks {
			e.stacks = append(e.stacks, j)
		}
		break

	default:
		s := New(Message(target.Error()), CodeUnknown)
		e.stacks = append(e.stacks, *s)
	}

	return e
}

func (e *Errors) Last() Error {
	return e.stacks[len(e.stacks)-1]
}

func (e *Errors) First() Error {
	return e.stacks[0]
}

func (e *Errors) IsWithStrategy(err error, comparable Comparable) bool {
	return comparable.Compare(e, err)
}

func (e *Errors) Is(err error) bool {
	return e.IsWithStrategy(err, CompareMessageOnlyStrategy())
}

func (e *Errors) Error() string {
	b, _ := e.MarshalJSON()

	return string(b)
}

type mErrors struct {
	First Error   `json:"first"`
	Last  Error   `json:"last"`
	Stack []Error `json:"stack,omitempty"`
}

type errorWrapper struct {
	Error mErrors `json:"error"`
}

func (e *Errors) MarshalJSON() ([]byte, error) {
	r := errorWrapper{
		Error: mErrors{
			First: e.First(),
			Last:  e.Last(),
			Stack: e.stacks,
		},
	}

	return json.Marshal(r)
}

func (e *Errors) UnmarshalJSON(b []byte) error {
	var r errorWrapper
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	e.stacks = r.Error.Stack

	return nil
}

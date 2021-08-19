package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Stack(source error, targets ...error) (result *Errors) {
	if source == nil && (len(targets) == 0 || targets == nil) {
		return nil
	}

	if source == nil && len(targets) != 0 {
		source = targets[0]
		targets = targets[1:]
	}

	switch e := source.(type) {
	case *Error:
		result = newErrors(e.message, e.code)
		break
	case *Errors:
		result = e
		break
	default:
		if source == nil {
			break
		}
		result = newErrors(Message(source.Error()), CodeUnknown)
	}

	for _, target := range targets {
		result = result.Stack(target)
	}

	return
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

func (e *Errors) Stack(target error) *Errors {
	if e == nil {
		return Stack(target)
	}

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
	b := strings.Builder{}

	b.WriteString(" ============== Errors ==============\n")

	for i, s := range e.stacks {
		b.WriteString(fmt.Sprintf("\t%d. %s\n", i+1, s.Error()))
	}

	b.WriteString("\n")

	return b.String()
}

func (e *Errors) MarshalJSON() ([]byte, error) {
	type data struct {
		First Error   `json:"first"`
		Last  Error   `json:"last"`
		Stack []Error `json:"stack,omitempty"`
	}

	type result struct {
		Error data `json:"error"`
	}

	r := result{
		Error: data{
			First: e.First(),
			Last:  e.Last(),
			Stack: e.stacks,
		},
	}

	return json.Marshal(r)
}

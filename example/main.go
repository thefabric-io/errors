package main

import (
	"fmt"

	"github.com/thefabric-io/errors"
)

var (
	Err1 = errors.New("error 1", "ERR001")
	Err2 = errors.New("error 2", "ERR002")
	Err3 = errors.New("error 3", "ERR003")
	Err4 = errors.New("error 4", "ERR004")
	Err5 = errors.New("error 5", "ERR005")
	Err6 = errors.New("error 6", "ERR006")
	Err7 = errors.New("error 7", "ERR007")
)

func main() {
	var result error
	result = errors.Stack(result, Err1, Err2, fmt.Errorf("trivial error"))

	// Stacking directly with Errors type
	if e, ok := result.(*errors.Errors); ok {
		_ = e.Stack(Err3)
	}

	result = errors.Stack(
		result,
		Err4,
		errors.New(
			errors.Message(fmt.Sprintf("error with a subject of '%s' and a value of '%.2f' €", "10c39745-c7fe-429f-a0fb-5035dbdc6c47", 12.00)),
			errors.CodeSubject,
		))

	result = errors.Stack(result, Err5, Err6, Err7, fmt.Errorf("trivial error with a value of '%s'", "0.00 €"))

	fmt.Println(result)

	b, err := result.(*errors.Errors).MarshalJSON()
	if err != nil {
		fmt.Println(errors.Stack(errors.New("could not marshal", "MF001"), err))
	}

	fmt.Println(" ============== JSON Representation ==============")
	fmt.Println(string(b))
}

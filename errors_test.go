package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	type args struct {
		source  error
		targets []error
	}
	tests := []struct {
		name       string
		args       args
		wantResult error
	}{
		{
			"err1 error, err2 error",
			args{
				source:  errors.New("first error"),
				targets: []error{errors.New("second error")},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeUnknown,
					},
					{
						message: "second error",
						code:    CodeUnknown,
					},
				},
			},
		},
		{
			"err1 error, err2 Error",
			args{
				source:  errors.New("first error"),
				targets: []error{New("second error", CodeInvalid)},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeUnknown,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
				},
			},
		},
		{
			"err1 Error, err2 error",
			args{
				source:  New("first error", CodeInvalid),
				targets: []error{errors.New("second error")},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeUnknown,
					},
				},
			},
		},
		{
			"err1 Errors, err2 error",
			args{
				source:  newErrors("first error", CodeInvalid),
				targets: []error{errors.New("second error")},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeUnknown,
					},
				},
			},
		},
		{
			"err1 error, err2 Errors",
			args{
				source:  errors.New("first error"),
				targets: []error{newErrors("second error", CodeInvalid)},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeUnknown,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
				},
			},
		},
		{
			"err1 Errors, err2 Errors",
			args{
				source:  Stack(newErrors("first error", CodeInvalid), newErrors("sub first error", CodeInvalid)),
				targets: []error{newErrors("second error", CodeInvalid)},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "sub first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
				},
			},
		},
		{
			"err1 nil, err2 nil",
			args{
				source:  nil,
				targets: nil,
			},
			nil,
		},
		{
			"err1 nil, err2 nil element in array 1",
			args{
				source:  nil,
				targets: []error{nil, nil, nil, nil},
			},
			nil,
		},
		{
			"err1 nil, err2 nil element in array 2",
			args{
				source: nil,
				targets: []error{nil, &Error{
					message: "first error",
					code:    CodeInvalid,
				}, nil, nil, &Error{
					message: "second error",
					code:    CodeInvalid,
				}},
			},
			&Errors{stacks: []Error{
				{
					message: "first error",
					code:    CodeInvalid,
				},
				{
					message: "second error",
					code:    CodeInvalid,
				},
			}},
		},
		{
			"err1 Error, err2 nil element in array 3",
			args{
				source: &Error{
					message: "first error",
					code:    CodeInvalid,
				},
				targets: []error{nil, &Error{
					message: "second error",
					code:    CodeInvalid,
				}, nil, nil, &Error{
					message: "third error",
					code:    CodeInvalid,
				}},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
					{
						message: "third error",
						code:    CodeInvalid,
					},
				},
			},
		},
		{
			"err1 Error, err2 nil element in array 4",
			args{
				source: &Errors{
					stacks: nil,
				},
				targets: []error{nil, &Error{
					message: "first error",
					code:    CodeInvalid,
				}, nil, nil, &Error{
					message: "second error",
					code:    CodeInvalid,
				}},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
				},
			},
		},
		{
			"err1 Error, err2 nil element in array 5",
			args{
				source: &Errors{
					stacks: []Error{
						{
							message: "first error",
							code:    CodeInvalid,
						},
					},
				},
				targets: []error{nil, &Error{
					message: "second error",
					code:    CodeInvalid,
				}, nil, nil, &Error{
					message: "third error",
					code:    CodeInvalid,
				}},
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
					{
						message: "third error",
						code:    CodeInvalid,
					},
				},
			},
		},
		{
			"err1 error, err2 nil",
			args{
				source:  errors.New("first error"),
				targets: nil,
			},
			&Errors{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeUnknown,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Stack(tt.args.source, tt.args.targets...); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Stack() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestStackBis(t *testing.T) {
	var r error

	if got := Stack(r, nil); !reflect.DeepEqual(got, nil) {
		t.Errorf("Stack() Exception = %v, want %v", got, nil)
	}
}

func TestErrors_Stack(t *testing.T) {
	type fields struct {
		stacks []Error
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Errors
	}{
		{
			"trivial error instance",
			fields{stacks: []Error{*New("first error", CodeInvalid)}},
			args{target: errors.New("second error")},
			&Errors{stacks: []Error{
				{
					message: "first error",
					code:    CodeInvalid,
				},
				{
					message: "second error",
					code:    CodeUnknown,
				},
			}},
		},
		{
			"Error instance",
			fields{stacks: []Error{*New("first error", CodeInvalid)}},
			args{target: New("second error", CodeInvalid)},
			&Errors{stacks: []Error{
				{
					message: "first error",
					code:    CodeInvalid,
				},
				{
					message: "second error",
					code:    CodeInvalid,
				},
			}},
		},
		{
			"Errors instance",
			fields{stacks: []Error{*New("first error", CodeInvalid)}},
			args{target: newErrors("second error", CodeInvalid).Stack(errors.New("last error"))},
			&Errors{stacks: []Error{
				{
					message: "first error",
					code:    CodeInvalid,
				},
				{
					message: "second error",
					code:    CodeInvalid,
				},
				{
					message: "last error",
					code:    CodeUnknown,
				},
			}},
		},
	}

	/*var r Errors
	want := r
	if got := r.Stack(nil); !reflect.DeepEqual(got, want) {
		t.Errorf("Stack() Exception = %v, want %v", got, want)
	}*/

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Errors{
				stacks: tt.fields.stacks,
			}
			if got := e.Stack(tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_Last(t *testing.T) {
	type fields struct {
		stacks []Error
	}
	tests := []struct {
		name   string
		fields fields
		want   Error
	}{
		{
			"one error",
			fields{
				stacks: []Error{{
					message: "first error",
					code:    CodeInvalid,
				}},
			},
			Error{
				message: "first error",
				code:    CodeInvalid,
			},
		},
		{
			"multiple error",
			fields{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
					{
						message: "third error",
						code:    CodeInvalid,
					},
				},
			},
			Error{
				message: "third error",
				code:    CodeInvalid,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Errors{
				stacks: tt.fields.stacks,
			}
			if got := e.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_First(t *testing.T) {
	type fields struct {
		stacks []Error
	}
	tests := []struct {
		name   string
		fields fields
		want   Error
	}{
		{
			"one error",
			fields{
				stacks: []Error{{
					message: "first error",
					code:    CodeInvalid,
				}},
			},
			Error{
				message: "first error",
				code:    CodeInvalid,
			},
		},
		{
			"multiple error",
			fields{
				stacks: []Error{
					{
						message: "first error",
						code:    CodeInvalid,
					},
					{
						message: "second error",
						code:    CodeInvalid,
					},
					{
						message: "third error",
						code:    CodeInvalid,
					},
				},
			},
			Error{
				message: "first error",
				code:    CodeInvalid,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Errors{
				stacks: tt.fields.stacks,
			}
			if got := e.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWithStrategy(t *testing.T) {
	type args struct {
		err1       error
		err2       error
		comparable Comparable
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"different trivial errors",
			args{
				err1:       errors.New("first error"),
				err2:       errors.New("second error"),
				comparable: CompareMessageOnlyStrategy(),
			},
			false,
		},
		{
			"different errors types with same message",
			args{
				err1:       errors.New("trivial error"),
				err2:       New("trivial error", CodeInvalid),
				comparable: CompareMessageOnlyStrategy(),
			},
			true,
		},
		{
			"stacked errors",
			args{
				err1:       New("trivial error", CodeInvalid),
				err2:       Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				comparable: CompareMessageOnlyStrategy(),
			},
			true,
		},
		{
			"double stacked errors",
			args{
				err1:       Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				err2:       Stack(New("another different error", CodeInvalid), New("trivial error", CodeInvalid)),
				comparable: CompareMessageOnlyStrategy(),
			},
			true,
		},
		{
			"double stacked errors different",
			args{
				err1:       Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				err2:       Stack(New("another different error", CodeInvalid), New("another trivial error", CodeInvalid)),
				comparable: CompareMessageOnlyStrategy(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWithStrategy(tt.args.err1, tt.args.err2, tt.args.comparable); got != tt.want {
				t.Errorf("IsWithStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIs(t *testing.T) {
	type args struct {
		err1 error
		err2 error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"different trivial errors",
			args{
				err1: errors.New("first error"),
				err2: errors.New("second error"),
			},
			false,
		},
		{
			"different errors types with same message",
			args{
				err1: errors.New("trivial error"),
				err2: New("trivial error", CodeInvalid),
			},
			true,
		},
		{
			"stacked errors",
			args{
				err1: New("trivial error", CodeInvalid),
				err2: Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
			},
			true,
		},
		{
			"double stacked errors",
			args{
				err1: Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				err2: Stack(New("another different error", CodeInvalid), New("trivial error", CodeInvalid)),
			},
			true,
		},
		{
			"double stacked errors different",
			args{
				err1: Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				err2: Stack(New("another different error", CodeInvalid), New("another trivial error", CodeInvalid)),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err1, tt.args.err2); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_IsWithStrategy(t *testing.T) {
	type fields struct {
		stacks []Error
	}
	type args struct {
		err        error
		comparable Comparable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"different trivial errors",
			fields{
				stacks: []Error{*New("first error", CodeUnknown), *New("second error", CodeUnknown)},
			},
			args{
				err:        errors.New("third error"),
				comparable: CompareMessageOnlyStrategy(),
			},
			false,
		},
		{
			"same trivial errors",
			fields{
				stacks: []Error{*New("first error", CodeUnknown), *New("second error", CodeUnknown)},
			},
			args{
				err:        errors.New("first error"),
				comparable: CompareMessageOnlyStrategy(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Errors{
				stacks: tt.fields.stacks,
			}
			if got := e.IsWithStrategy(tt.args.err, tt.args.comparable); got != tt.want {
				t.Errorf("IsWithStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_Is(t *testing.T) {
	type fields struct {
		stacks []Error
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"different trivial errors",
			fields{
				stacks: []Error{*New("first error", CodeUnknown), *New("second error", CodeUnknown)},
			},
			args{
				err: errors.New("third error"),
			},
			false,
		},
		{
			"same trivial errors",
			fields{
				stacks: []Error{*New("first error", CodeUnknown), *New("second error", CodeUnknown)},
			},
			args{
				err: errors.New("first error"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Errors{
				stacks: tt.fields.stacks,
			}
			if got := e.Is(tt.args.err); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

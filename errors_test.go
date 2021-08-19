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
		wantResult *Errors
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
					Error{
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
			"err1 nil, err2 nil element in array",
			args{
				source:  nil,
				targets: []error{nil, nil, nil, nil},
			},
			nil,
		},
		{
			"err1 nil, err2 nil element in array",
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
			"err1 Error, err2 nil element in array",
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
			"err1 Error, err2 nil element in array",
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
			"err1 Error, err2 nil element in array",
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

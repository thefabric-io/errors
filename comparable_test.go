package errors

import (
	"errors"
	"testing"
)

func TestCompareMessageOnly_Compare(t *testing.T) {
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
			c := CompareMessageOnlyStrategy()
			if got := c.Compare(tt.args.err1, tt.args.err2); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareCodeOnly_Compare(t *testing.T) {
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
			true,
		},
		{
			"different errors types with same code",
			args{
				err1: New("first trivial error", CodeInvalid),
				err2: New("trivial error", CodeInvalid),
			},
			true,
		},
		{
			"stacked errors",
			args{
				err1: New("trivial error", CodeInvalid),
				err2: Stack(New("trivial error", CodeInvalid), New("another error", CodeUnknown)),
			},
			true,
		},
		{
			"double stacked errors",
			args{
				err1: Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				err2: Stack(New("another different error", CodeInvalid), New("trivial error", CodeUnknown)),
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
			c := CompareCodeOnlyStrategy()
			if got := c.Compare(tt.args.err1, tt.args.err2); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareStrict_Compare(t *testing.T) {
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
			"different errors types with same message and code",
			args{
				err1: errors.New("trivial error"),
				err2: New("trivial error", CodeUnknown),
			},
			true,
		},
		{
			"stacked errors",
			args{
				err1: Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				err2: Stack(New("trivial error", CodeUnknown)),
			},
			true,
		},
		{
			"double stacked errors",
			args{
				err1: Stack(New("trivial error", CodeUnknown), New("another error", CodeUnknown)),
				err2: Stack(New("another different error", CodeInvalid), New("trivial error", CodeUnknown)),
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
			c := CompareStrictStrategy()
			if got := c.Compare(tt.args.err1, tt.args.err2); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

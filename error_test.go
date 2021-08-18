package errors

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		message Message
		code    Code
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{"nominal", args{
			message: "entity invalid",
			code:    CodeInvalid,
		},
			&Error{
				message: "entity invalid",
				code:    CodeInvalid,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.message, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		message Message
		code    Code
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"nominal",
			fields{
				message: "entity invalid",
				code:    CodeInvalid,
			},
			fmt.Sprintf("[%s] %s", CodeInvalid, "entity invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Error{
				message: tt.fields.message,
				code:    tt.fields.code,
			}
			if got := l.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

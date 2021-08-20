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
		{
			"nominal",
			args{
				message: "entity invalid",
				code:    CodeInvalid,
			},
			&Error{
				message: "entity invalid",
				code:    CodeInvalid,
			},
		},
		{
			"code unknown",
			args{
				message: "entity invalid",
				code:    "",
			},
			&Error{
				message: "entity invalid",
				code:    CodeUnknown,
			},
		},
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

func TestError_MarshalJSON(t *testing.T) {
	type fields struct {
		message Message
		code    Code
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"nominal",
			fields{
				message: "first error",
				code:    CodeInvalid,
			},
			[]byte("{\"message\":\"first error\",\"code\":\"invalid\"}"),
			false,
		},
		{
			"without code",
			fields{
				message: "first error",
			},
			[]byte("{\"message\":\"first error\"}"),
			false,
		},
		{
			"without message",
			fields{
				code: CodeInvalid,
			},
			[]byte("{\"code\":\"invalid\"}"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				message: tt.fields.message,
				code:    tt.fields.code,
			}
			got, err := e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCode_Equal(t *testing.T) {
	type args struct {
		c2 Code
	}
	tests := []struct {
		name string
		c    Code
		args args
		want bool
	}{
		{
			"not equal",
			CodeInvalid,
			args{c2: CodeUnknown},
			false,
		},
		{
			"equal",
			CodeInvalid,
			args{c2: CodeInvalid},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Equal(tt.args.c2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

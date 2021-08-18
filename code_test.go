package errors

import "testing"

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

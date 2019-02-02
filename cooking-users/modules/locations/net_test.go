package locations

import "testing"

func TestResponseMessage_MakeJsonResponse(t *testing.T) {
	type args struct {
		params interface{}
	}
	tests := []struct {
		name string
		s    ResponseMessage
		args args
		want string
	}{
		{
			name: "",
			s:    "FooBar",
			args: args{
				params: []int{1,2},
			},
			want: "[1,2]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.MakeJsonResponse(tt.args.params); got != tt.want {
				t.Errorf("ResponseMessage.MakeJsonResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseMessage_ToString(t *testing.T) {
	type args struct {
		params []interface{}
	}
	tests := []struct {
		name string
		s    ResponseMessage
		args args
		want string
	}{
		{
			name: "",
			s:    "foo %d %d",
			args: args{
				params: []interface{}{1,2},
			},
			want: "foo 1 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ToString(tt.args.params...); got != tt.want {
				t.Errorf("ResponseMessage.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseMessage_AsJsonError(t *testing.T) {
	tests := []struct {
		name string
		s    ResponseMessage
		want string
	}{
		{
			name: "foo",
			s:    "foo",
			want: "{\"error\":\"foo\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AsJsonError(); got != tt.want {
				t.Errorf("ResponseMessage.AsJsonError() = %v, want %v", got, tt.want)
			}
		})
	}
}

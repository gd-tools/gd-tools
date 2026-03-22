package protocol

import "testing"

func TestRequestHasHello(t *testing.T) {
	tests := []struct {
		name string
		req  *Request
		want bool
	}{
		{
			name: "nil request",
			req:  nil,
			want: false,
		},
		{
			name: "empty request",
			req:  &Request{},
			want: false,
		},
		{
			name: "empty greeting",
			req: &Request{
				Hello: Hello{
					Greeting: "",
				},
			},
			want: false,
		},
		{
			name: "hello greeting",
			req: &Request{
				Hello: Hello{
					Greeting: "Hello",
				},
			},
			want: true,
		},
		{
			name: "other greeting still valid",
			req: &Request{
				Hello: Hello{
					Greeting: "Hi there",
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasHello()
			if got != tt.want {
				t.Fatalf("HasHello() = %v, want %v", got, tt.want)
			}
		})
	}
}

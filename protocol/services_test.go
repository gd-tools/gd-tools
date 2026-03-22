package protocol

import "testing"

func TestRequestAddService(t *testing.T) {
	var req Request

	req.AddService("apache2")
	req.AddService("postfix")
	req.AddService("apache2")
	req.AddService("")

	want := []string{"apache2", "postfix"}

	if len(req.Services) != len(want) {
		t.Fatalf("expected %d services, got %d", len(want), len(req.Services))
	}
	for i := range want {
		if req.Services[i] != want[i] {
			t.Fatalf("service %d = %q, want %q", i, req.Services[i], want[i])
		}
	}
}

func TestRequestAddServiceNilReceiver(t *testing.T) {
	var req *Request
	req.AddService("apache2") // should not panic
}

func TestRequestHasServiceList(t *testing.T) {
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
			name: "one service",
			req: &Request{
				ServiceList: ServiceList{
					Services: []string{"apache2"},
				},
			},
			want: true,
		},
		{
			name: "multiple services",
			req: &Request{
				ServiceList: ServiceList{
					Services: []string{"apache2", "postfix"},
				},
			},
			want: true,
		},
		{
			name: "empty slice",
			req: &Request{
				ServiceList: ServiceList{
					Services: []string{},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasServiceList()
			if got != tt.want {
				t.Fatalf("HasServiceList() = %v, want %v", got, tt.want)
			}
		})
	}
}

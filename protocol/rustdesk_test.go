package protocol

import (
	"bytes"
	"testing"
)

func TestRustDeskFQDN(t *testing.T) {
	tests := []struct {
		name string
		rd   *RustDesk
		want string
	}{
		{"nil", nil, ""},
		{"host and domain", &RustDesk{HostName: "rd", DomainName: "example.org"}, "rd.example.org"},
		{"domain only", &RustDesk{DomainName: "example.org"}, "example.org"},
		{"host only", &RustDesk{HostName: "rd"}, "rd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rd.FQDN()
			if got != tt.want {
				t.Fatalf("FQDN() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRustDeskPrivateRoundTrip(t *testing.T) {
	rd := &RustDesk{}
	want := []byte("secret-key-data")

	rd.SetPrivate(want)

	got, err := rd.GetPrivate()
	if err != nil {
		t.Fatalf("GetPrivate() returned error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("GetPrivate() = %q, want %q", got, want)
	}
}

func TestRustDeskSetPrivateEmpty(t *testing.T) {
	rd := &RustDesk{PrivateB64: "abc"}

	rd.SetPrivate(nil)

	if rd.PrivateB64 != "" {
		t.Fatalf("expected empty PrivateB64, got %q", rd.PrivateB64)
	}
}

func TestRustDeskGetPrivateEmpty(t *testing.T) {
	rd := &RustDesk{}

	got, err := rd.GetPrivate()
	if err != nil {
		t.Fatalf("GetPrivate() returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil private data, got %v", got)
	}
}

func TestRustDeskGetPrivateInvalidBase64(t *testing.T) {
	rd := &RustDesk{PrivateB64: "%%%invalid%%%"}
	_, err := rd.GetPrivate()
	if err == nil {
		t.Fatalf("expected base64 decode error")
	}
}

func TestRustDeskPublicNormalize(t *testing.T) {
	rd := &RustDesk{}

	rd.SetPublic("  public-key  ")

	if got := rd.GetPublic(); got != "public-key" {
		t.Fatalf("GetPublic() = %q, want %q", got, "public-key")
	}
}

func TestRequestHasRustDeskApp(t *testing.T) {
	tests := []struct {
		name string
		req  *Request
		want bool
	}{
		{"nil request", nil, false},
		{"empty request", &Request{}, false},
		{
			"with rustdesk",
			&Request{
				RustDeskApp: RustDeskApp{
					RustDesk: &RustDesk{Version: "1.0.0"},
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasRustDeskApp()
			if got != tt.want {
				t.Fatalf("HasRustDeskApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

package protocol

import (
	"path/filepath"
	"testing"
)

func TestDownloadMarkerPath(t *testing.T) {
	tests := []struct {
		name string
		dl   *Download
		root string
		want string
	}{
		{
			name: "empty marker",
			dl: &Download{
				Directory: "downloads",
				Marker:    "",
			},
			root: "/tmp/root",
			want: "",
		},
		{
			name: "marker without directory",
			dl: &Download{
				Marker: "ready.txt",
			},
			root: "/tmp/root",
			want: filepath.Join("/tmp/root", "ready.txt"),
		},
		{
			name: "marker with directory",
			dl: &Download{
				Directory: "downloads/rustdesk",
				Marker:    "ready.txt",
			},
			root: "/tmp/root",
			want: filepath.Join("/tmp/root", "downloads/rustdesk", "ready.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dl.MarkerPath(tt.root)
			if got != tt.want {
				t.Fatalf("MarkerPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequestAddDownload(t *testing.T) {
	req := &Request{}
	dl := &Download{
		DownloadURL: "https://example.com/app.tar.gz",
		Filename:    "app.tar.gz",
	}

	req.AddDownload(dl)

	if len(req.Downloads) != 1 {
		t.Fatalf("expected 1 download, got %d", len(req.Downloads))
	}
	if req.Downloads[0] != dl {
		t.Fatalf("expected appended download pointer to match input")
	}
}

func TestRequestAddDownloadIgnoresNilReceiver(t *testing.T) {
	var req *Request
	dl := &Download{Filename: "app.tar.gz"}

	req.AddDownload(dl)
}

func TestRequestAddDownloadIgnoresNilDownload(t *testing.T) {
	req := &Request{}

	req.AddDownload(nil)

	if len(req.Downloads) != 0 {
		t.Fatalf("expected 0 downloads, got %d", len(req.Downloads))
	}
}

func TestRequestHasDownloadList(t *testing.T) {
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
			name: "with downloads",
			req: &Request{
				DownloadList: DownloadList{
					Downloads: []*Download{
						{Filename: "app.tar.gz"},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasDownloadList()
			if got != tt.want {
				t.Fatalf("HasDownloadList() = %v, want %v", got, tt.want)
			}
		})
	}
}

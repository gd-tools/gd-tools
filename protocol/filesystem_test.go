package protocol

import "testing"

func TestRequestAddMount(t *testing.T) {
	req := &Request{}
	mount := &Mount{
		Provider: "hetzner",
		ID:       "123456789",
		Path:     "/var/gd-tools",
	}

	req.AddMount(mount)

	if len(req.Mounts) != 1 {
		t.Fatalf("expected 1 mount, got %d", len(req.Mounts))
	}
	if req.Mounts[0] != mount {
		t.Fatalf("expected appended mount pointer to match input")
	}
}

func TestRequestAddMountIgnoresNil(t *testing.T) {
	req := &Request{}

	req.AddMount(nil)

	if len(req.Mounts) != 0 {
		t.Fatalf("expected 0 mounts, got %d", len(req.Mounts))
	}
}

func TestRequestAddDirectory(t *testing.T) {
	req := &Request{}
	dir := &Directory{
		Path:  "/var/gd-tools/data",
		Mode:  "0755",
		User:  "gd-tools",
		Group: "gd-tools",
	}

	req.AddDirectory(dir)

	if len(req.Directories) != 1 {
		t.Fatalf("expected 1 directory, got %d", len(req.Directories))
	}
	if req.Directories[0] != dir {
		t.Fatalf("expected appended directory pointer to match input")
	}
}

func TestRequestAddDirectoryIgnoresNil(t *testing.T) {
	req := &Request{}

	req.AddDirectory(nil)

	if len(req.Directories) != 0 {
		t.Fatalf("expected 0 directories, got %d", len(req.Directories))
	}
}

func TestRequestHasFilesystem(t *testing.T) {
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
			name: "mount only",
			req: &Request{
				Filesystem: Filesystem{
					Mounts: []*Mount{
						{Path: "/var/gd-tools"},
					},
				},
			},
			want: true,
		},
		{
			name: "directory only",
			req: &Request{
				Filesystem: Filesystem{
					Directories: []*Directory{
						{Path: "/var/gd-tools/data"},
					},
				},
			},
			want: true,
		},
		{
			name: "both",
			req: &Request{
				Filesystem: Filesystem{
					Mounts: []*Mount{
						{Path: "/var/gd-tools"},
					},
					Directories: []*Directory{
						{Path: "/var/gd-tools/data"},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasFilesystem()
			if got != tt.want {
				t.Fatalf("HasFilesystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

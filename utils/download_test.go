package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadMarkerPath(t *testing.T) {
	tests := []struct {
		name string
		dl   Download
		root string
		want string
	}{
		{
			name: "empty marker",
			dl:   Download{},
			root: "/tmp/root",
			want: "",
		},
		{
			name: "marker without directory",
			dl: Download{
				Marker: "done.marker",
			},
			root: "/tmp/root",
			want: filepath.Join("/tmp/root", "done.marker"),
		},
		{
			name: "marker with directory",
			dl: Download{
				Directory: "downloads",
				Marker:    "done.marker",
			},
			root: "/tmp/root",
			want: filepath.Join("/tmp/root", "downloads", "done.marker"),
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

func TestDownloadMarkerExists(t *testing.T) {
	root := t.TempDir()

	dl := Download{
		Directory: "downloads",
		Marker:    "done.marker",
	}

	exists, err := dl.MarkerExists(root)
	if err != nil {
		t.Fatalf("MarkerExists() error = %v", err)
	}
	if exists {
		t.Fatal("MarkerExists() = true, want false before marker exists")
	}

	path := dl.MarkerPath(root)
	if path == "" {
		t.Fatal("MarkerPath() returned empty path")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte("ok"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	exists, err = dl.MarkerExists(root)
	if err != nil {
		t.Fatalf("MarkerExists() error = %v", err)
	}
	if !exists {
		t.Fatal("MarkerExists() = false, want true after marker exists")
	}
}

func TestDownloadMarkerExistsEmptyMarker(t *testing.T) {
	root := t.TempDir()

	dl := Download{}
	exists, err := dl.MarkerExists(root)
	if err != nil {
		t.Fatalf("MarkerExists() error = %v", err)
	}
	if exists {
		t.Fatal("MarkerExists() = true, want false for empty marker")
	}
}

func TestHashesAndVerify(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.bin")
	data := []byte("hello world\n")

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	md5sum, sha256sum, sha512sum, err := Hashes(path)
	if err != nil {
		t.Fatalf("Hashes() error = %v", err)
	}

	if md5sum == "" || sha256sum == "" || sha512sum == "" {
		t.Fatal("Hashes() returned empty checksum")
	}

	dl := Download{
		Filename: "test.bin",
		MD5:      md5sum,
		SHA256:   sha256sum,
		SHA512:   sha512sum,
	}

	if err := dl.Verify(path); err != nil {
		t.Fatalf("Verify() error = %v, want nil", err)
	}
}

func TestVerifyMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.bin")

	if err := os.WriteFile(path, []byte("hello world\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	dl := Download{
		Filename: "test.bin",
		SHA256:   "definitely-wrong",
	}

	err := dl.Verify(path)
	if err == nil {
		t.Fatal("Verify() error = nil, want mismatch error")
	}
}

func TestHashesMissingFile(t *testing.T) {
	_, _, _, err := Hashes(filepath.Join(t.TempDir(), "missing.bin"))
	if err == nil {
		t.Fatal("Hashes() error = nil, want error for missing file")
	}
}

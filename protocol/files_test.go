package protocol

import "testing"

func TestFileString(t *testing.T) {
	f := &File{
		Task: FileTaskWrite,
		Path: "/etc/test.conf",
	}

	got := f.String()
	want := "task=write path=/etc/test.conf"

	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestFileStringNil(t *testing.T) {
	var f *File

	got := f.String()
	if got != "<nil>" {
		t.Fatalf("expected <nil>, got %q", got)
	}
}

func TestRequestAddFile(t *testing.T) {
	req := &Request{}
	f := &File{Path: "/tmp/test"}

	req.AddFile(f)

	if len(req.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(req.Files))
	}
	if req.Files[0] != f {
		t.Fatalf("expected appended file pointer to match input")
	}
}

func TestRequestAddFileIgnoresNil(t *testing.T) {
	req := &Request{}

	req.AddFile(nil)

	if len(req.Files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(req.Files))
	}
}

func TestRequestHasFileList(t *testing.T) {
	if (&Request{}).HasFileList() {
		t.Fatalf("expected false for empty request")
	}

	req := &Request{}
	req.AddFile(&File{Path: "/tmp/x"})

	if !req.HasFileList() {
		t.Fatalf("expected true after adding file")
	}
}

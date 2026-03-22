package protocol

import (
	"fmt"
)

// FileTask defines the type of file operation.
type FileTask string

const (
	FileTaskWrite   FileTask = "write"
	FileTaskExtract FileTask = "extract"
	FileTaskRead    FileTask = "read"
	FileTaskDelete  FileTask = "delete"
	FileTaskProcess FileTask = "process"
	FileTaskPostmap FileTask = "postmap"
	FileTaskCleanup FileTask = "cleanup"
)

// File describes a file task to be executed on prod.
type File struct {
	Task    FileTask `json:"task"`
	Path    string   `json:"path"`
	Content []byte   `json:"content,omitempty"`
	Target  string   `json:"target,omitempty"`
	Backup  bool     `json:"backup,omitempty"`
	Mode    string   `json:"mode,omitempty"`
	User    string   `json:"user,omitempty"`
	Group   string   `json:"group,omitempty"`
	Service string   `json:"service,omitempty"` // optional service restart/reload
}

// String returns a short human-readable description of the file task.
func (f *File) String() string {
	if f == nil {
		return "<nil>"
	}
	return fmt.Sprintf("task=%s path=%s", f.Task, f.Path)
}

// FileList contains file operations for prod.
type FileList struct {
	Files []*File `json:"files,omitempty"`
}

// AddFile adds a file handling task to the request.
func (req *Request) AddFile(file *File) {
	if req == nil || file == nil {
		return
	}
	req.Files = append(req.Files, file)
}

// HasFileList reports whether the request contains at least one file entry.
func (req *Request) HasFileList() bool {
	if req == nil {
		return false
	}
	return len(req.Files) > 0
}

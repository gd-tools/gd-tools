package protocol

import (
	"fmt"
)

type File struct {
	Task    string `json:"task"`
	Path    string `json:"path"`
	Content []byte `json:"content,omitempty"`
	Target  string `json:"target,omitempty"`
	Backup  bool   `json:"backup,omitempty"`
	Mode    string `json:"mode,omitempty"`
	User    string `json:"user,omitempty"`
	Group   string `json:"group,omitempty"`
	Service string `json:"service,omitempty"`
}

func (f *File) String() string {
	return fmt.Sprintf("Task='%s' Path='%s'", f.Task, f.Path)
}

type FileList struct {
	Files []*File `json:"files"`
}

func (req *Request) AddFile(file *File) {
	if req != nil && file != nil {
		req.Files = append(req.Files, file)
	}
}

func (req *Request) HasFileList() bool {
	if req == nil {
		return false
	}
	return len(req.Files) > 0
}

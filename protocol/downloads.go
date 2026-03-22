package protocol

import (
	"path/filepath"
)

// Download describes a downloadable asset, e.g. zip archive or binary.
type Download struct {
	DownloadURL string `json:"download_url"`
	Filename    string `json:"filename"`
	Directory   string `json:"directory"`
	Binary      string `json:"binary"`
	Marker      string `json:"marker"`
	MD5         string `json:"md5"`
	SHA256      string `json:"sha256"`
	SHA512      string `json:"sha512"`
}

// MarkerPath returns the full marker path below root.
// Marker is interpreted as a relative path inside Directory.
// Return an empty string if Marker is not set.
func (dl *Download) MarkerPath(root string) string {
	if dl.Marker == "" {
		return ""
	}
	if dl.Directory == "" {
		return filepath.Join(root, dl.Marker)
	}
	return filepath.Join(root, dl.Directory, dl.Marker)
}

type DownloadList struct {
	Downloads []*Download `json:"downloads,omitempty"`
}

func (req *Request) HasDownloadList() bool {
	if req == nil {
		return false
	}
	return len(req.Downloads) > 0
}

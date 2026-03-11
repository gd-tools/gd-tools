package releases

import (
	"path/filepath"
)

type Download struct {
	DownloadURL string `json:"download_url"`
	Filename    string `json:"filename"`
	Directory   string `json:"directory"`
	Binary      string `json:"binary"`
	MD5         string `json:"md5"`
	SHA256      string `json:"sha256"`
	SHA512      string `json:"sha512"`
}

var (
	DownloadsDir string
)

func SetDownloadsDir(path string) {
	if path == "" {
		path = "/root/Downloads"
	}
	DownloadsDir = path
}

func GetDownloadsDir(name string) string {
	if DownloadsDir == "" {
		SetDownloadsDir("")
	}
	if name == "" {
		return DownloadsDir
	}
	return filepath.Join(DownloadsDir, name)
}

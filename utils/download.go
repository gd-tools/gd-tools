package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
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
func (dl *Download) MarkerPath(root string) (string, error) {
	if dl.Marker == "" {
		return "", nil
	}
	if dl.Directory == "" {
		return filepath.Join(root, dl.Marker), nil
	}
	return filepath.Join(root, dl.Directory, dl.Marker), nil
}

// MarkerExists reports whether the marker path exists below root.
func (dl *Download) MarkerExists(root string) (bool, error) {
	path, err := dl.MarkerPath(root)
	if err != nil {
		return false, err
	}
	if path == "" {
		return false, nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil // the only positive outcome
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Verify verifies one local file against the configured hashes.
func (dl *Download) Verify(path string) error {
	md5sum, sha256sum, sha512sum, err := Hashes(path)
	if err != nil {
		return err
	}

	if dl.MD5 != "" && dl.MD5 != md5sum {
		return fmt.Errorf("MD5 mismatch for %s", dl.Filename)
	}
	if dl.SHA256 != "" && dl.SHA256 != sha256sum {
		return fmt.Errorf("SHA256 mismatch for %s", dl.Filename)
	}
	if dl.SHA512 != "" && dl.SHA512 != sha512sum {
		return fmt.Errorf("SHA512 mismatch for %s", dl.Filename)
	}

	return nil
}

// Hashes computes MD5, SHA256, and SHA512 for one file.
func Hashes(path string) (string, string, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}
	defer f.Close()

	hMD5 := md5.New()
	hSHA256 := sha256.New()
	hSHA512 := sha512.New()

	if _, err := io.Copy(io.MultiWriter(hMD5, hSHA256, hSHA512), f); err != nil {
		return "", "", "", err
	}

	md5sum := hex.EncodeToString(hMD5.Sum(nil))
	sha256sum := hex.EncodeToString(hSHA256.Sum(nil))
	sha512sum := hex.EncodeToString(hSHA512.Sum(nil))

	return md5sum, sha256sum, sha512sum, nil
}

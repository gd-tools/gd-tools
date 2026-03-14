package platform

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

// Download describes a specific loadable asset, e.g. a zip archive or binary.
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

// LocalPath returns the local cached download path.
func (dl *Download) LocalPath(pf *Platform) string {
	return pf.DownloadsDir(dl.Filename)
}

// TargetPath returns the installation target path for the binary.
// Return an empty string if this download is not installed as a binary.
func (dl *Download) TargetPath(pf *Platform) string {
	if dl.Binary == "" {
		return ""
	}
	return pf.BinDir(dl.Binary)
}

// DirectoryPath returns the base directory for this download below root.
// Return an empty string if Directory is not set.
func (dl *Download) DirectoryPath(root string) string {
	if dl.Directory == "" {
		return ""
	}
	return filepath.Join(root, dl.Directory)
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

// MarkerExists reports whether the marker path exists below root.
func (dl *Download) MarkerExists(root string) (bool, error) {
	path := dl.MarkerPath(root)
	if path == "" {
		return false, nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ExistsLocal reports whether the cached download file exists.
func (dl *Download) ExistsLocal(pf *Platform) (bool, error) {
	path := dl.LocalPath(pf)

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// VerifyLocal verifies the cached local file against the configured hashes.
func (dl *Download) VerifyLocal(pf *Platform) error {
	return dl.VerifyFile(dl.LocalPath(pf))
}

// VerifyFile verifies one local file against the configured hashes.
func (dl *Download) VerifyFile(path string) error {
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

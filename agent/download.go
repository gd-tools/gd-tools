package agent

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gd-tools/gd-tools/releases"
)

// DownloadsTest checks if there is work to be done
func DownloadsTest(req *Request) bool {
	return req != nil && len(req.Downloads) > 0
}

func DownloadsHandler(req *Request, resp *Response) error {
	if req == nil || len(req.Downloads) == 0 || resp == nil {
		return nil
	}
	downloadsRoot := releases.GetDownloadsDir("")

	if err := os.MkdirAll(downloadsRoot, 0755); err != nil {
		return err
	}

	for _, dwn := range req.Downloads {
		path := filepath.Join(downloadsRoot, dwn.Filename)
		status := fmt.Sprintf("✅ download %s", path)
		if _, err := os.Stat(path); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if _, err := RunCommand("curl", "-fsSL", "-o", path, dwn.DownloadURL); err != nil {
				return err
			}
			status = fmt.Sprintf("download %s was successful", dwn.Filename)
		}

		md5sum, sha256sum, sha512sum, err := computeHashes(path)
		if err != nil {
			return err
		}

		if dwn.MD5 != "" && dwn.MD5 != md5sum {
			return fmt.Errorf("MD5 mismatch for %s", dwn.Filename)
		}
		if dwn.SHA256 != "" && dwn.SHA256 != sha256sum {
			return fmt.Errorf("SHA256 mismatch for %s", dwn.Filename)
		}
		if dwn.SHA512 != "" && dwn.SHA512 != sha512sum {
			return fmt.Errorf("SHA512 mismatch for %s", dwn.Filename)
		}
		resp.Say(status)

		if dwn.Binary != "" {
			_, err := RunCommand("install",
				"-m", "0755",
				"-o", "root",
				"-g", "root",
				path,
				releases.GetBinDir(dwn.Binary),
			)
			if err != nil {
				return err
			}
			resp.Sayf("copied %s to %s", path, releases.GetBinDir(dwn.Binary))
		}
	}

	return nil
}

func computeHashes(path string) (string, string, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}
	defer f.Close()

	hMd5 := md5.New()
	hSha256 := sha256.New()
	hSha512 := sha512.New()

	// Write simultaneously in all hashes
	if _, err := io.Copy(io.MultiWriter(hMd5, hSha256, hSha512), f); err != nil {
		return "", "", "", err
	}

	md5sum := hex.EncodeToString(hMd5.Sum(nil))
	sha256sum := hex.EncodeToString(hSha256.Sum(nil))
	sha512sum := hex.EncodeToString(hSha512.Sum(nil))

	return md5sum, sha256sum, sha512sum, nil
}

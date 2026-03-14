package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func testPlatform(t *testing.T) *Platform {
	t.Helper()

	pf, err := LoadPlatformWithPaths([]Path{
		{Name: PathRoot, Value: filepath.Join(t.TempDir(), "root")},
		{Name: PathVar, Value: filepath.Join(t.TempDir(), "var")},
		{Name: PathTools, Value: filepath.Join(t.TempDir(), "var", "gd-tools")},
		{Name: PathEtc, Value: filepath.Join(t.TempDir(), "etc")},
		{Name: PathBin, Value: filepath.Join(t.TempDir(), "bin")},
		{Name: PathRun, Value: filepath.Join(t.TempDir(), "run")},
		{Name: PathDownloads, Value: filepath.Join(t.TempDir(), "downloads")},
	})
	if err != nil {
		t.Fatalf("LoadPlatformWithPaths failed: %v", err)
	}
	return pf
}

func TestDownloadLocalPathAndTargetPath(t *testing.T) {
	pf := testPlatform(t)

	dl := &Download{
		Filename: "nextcloud.tar.bz2",
		Binary:   "occ",
	}

	if got := dl.LocalPath(pf); got != filepath.Join(pf.DownloadsDir(), "nextcloud.tar.bz2") {
		t.Fatalf("LocalPath mismatch: got %q", got)
	}

	if got := dl.TargetPath(pf); got != filepath.Join(pf.BinDir(), "occ") {
		t.Fatalf("TargetPath mismatch: got %q", got)
	}
}

func TestDownloadMarkerPath(t *testing.T) {
	dl := &Download{
		Directory: "nextcloud",
		Marker:    "README",
	}

	got := dl.MarkerPath("/srv/data")
	want := filepath.Join("/srv/data", "nextcloud", "README")
	if got != want {
		t.Fatalf("MarkerPath mismatch: got %q want %q", got, want)
	}
}

func TestDownloadMarkerExists(t *testing.T) {
	root := t.TempDir()
	dl := &Download{
		Directory: "nextcloud",
		Marker:    "README",
	}

	exists, err := dl.MarkerExists(root)
	if err != nil {
		t.Fatalf("MarkerExists failed: %v", err)
	}
	if exists {
		t.Fatal("marker should not exist yet")
	}

	path := filepath.Join(root, "nextcloud")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(path, "README"), []byte("ok"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	exists, err = dl.MarkerExists(root)
	if err != nil {
		t.Fatalf("MarkerExists failed: %v", err)
	}
	if !exists {
		t.Fatal("marker should exist")
	}
}

func TestDownloadExistsLocal(t *testing.T) {
	pf := testPlatform(t)

	if err := os.MkdirAll(pf.DownloadsDir(), 0755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	dl := &Download{
		Filename: "test.bin",
	}

	exists, err := dl.ExistsLocal(pf)
	if err != nil {
		t.Fatalf("ExistsLocal failed: %v", err)
	}
	if exists {
		t.Fatal("download should not exist yet")
	}

	if err := os.WriteFile(dl.LocalPath(pf), []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	exists, err = dl.ExistsLocal(pf)
	if err != nil {
		t.Fatalf("ExistsLocal failed: %v", err)
	}
	if !exists {
		t.Fatal("download should exist")
	}
}

func TestHashesAndVerifyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	md5sum, sha256sum, sha512sum, err := Hashes(path)
	if err != nil {
		t.Fatalf("Hashes failed: %v", err)
	}

	dl := &Download{
		Filename: "file.txt",
		MD5:      md5sum,
		SHA256:   sha256sum,
		SHA512:   sha512sum,
	}

	if err := dl.VerifyFile(path); err != nil {
		t.Fatalf("VerifyFile failed: %v", err)
	}
}

func TestVerifyFileMismatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dl := &Download{
		Filename: "file.txt",
		SHA256:   "deadbeef",
	}

	if err := dl.VerifyFile(path); err == nil {
		t.Fatal("expected checksum mismatch")
	}
}

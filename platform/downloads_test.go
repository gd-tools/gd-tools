package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadLocalPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			rootDir: "/root",
		},
	}

	dl := &Download{
		Filename: "archive.tar.gz",
	}

	got := dl.LocalPath(pf)
	want := "/root/Downloads/archive.tar.gz"

	if got != want {
		t.Fatalf("LocalPath mismatch: got %q want %q", got, want)
	}
}

func TestDownloadTargetPathEmpty(t *testing.T) {
	pf := &Platform{
		options: &Options{
			binDir: "/usr/local/bin",
		},
	}

	dl := &Download{}

	got := dl.TargetPath(pf)
	if got != "" {
		t.Fatalf("TargetPath mismatch: got %q want empty", got)
	}
}

func TestDownloadTargetPath(t *testing.T) {
	pf := &Platform{
		options: &Options{
			binDir: "/usr/local/bin",
		},
	}

	dl := &Download{
		Binary: "ocis",
	}

	got := dl.TargetPath(pf)
	want := "/usr/local/bin/ocis"

	if got != want {
		t.Fatalf("TargetPath mismatch: got %q want %q", got, want)
	}
}

func TestDownloadDirectoryPathEmpty(t *testing.T) {
	dl := &Download{}

	got := dl.DirectoryPath("/srv")
	if got != "" {
		t.Fatalf("DirectoryPath mismatch: got %q want empty", got)
	}
}

func TestDownloadDirectoryPath(t *testing.T) {
	dl := &Download{
		Directory: "nextcloud",
	}

	got := dl.DirectoryPath("/srv")
	want := filepath.Join("/srv", "nextcloud")

	if got != want {
		t.Fatalf("DirectoryPath mismatch: got %q want %q", got, want)
	}
}

func TestDownloadMarkerPathEmpty(t *testing.T) {
	dl := &Download{}

	got := dl.MarkerPath("/srv")
	if got != "" {
		t.Fatalf("MarkerPath mismatch: got %q want empty", got)
	}
}

func TestDownloadMarkerPathWithoutDirectory(t *testing.T) {
	dl := &Download{
		Marker: "installed.ok",
	}

	got := dl.MarkerPath("/srv")
	want := filepath.Join("/srv", "installed.ok")

	if got != want {
		t.Fatalf("MarkerPath mismatch: got %q want %q", got, want)
	}
}

func TestDownloadMarkerPathWithDirectory(t *testing.T) {
	dl := &Download{
		Directory: "nextcloud",
		Marker:    "config/config.php",
	}

	got := dl.MarkerPath("/srv")
	want := filepath.Join("/srv", "nextcloud", "config/config.php")

	if got != want {
		t.Fatalf("MarkerPath mismatch: got %q want %q", got, want)
	}
}

func TestDownloadMarkerExistsEmptyMarker(t *testing.T) {
	dl := &Download{}

	ok, err := dl.MarkerExists("/srv")
	if err != nil {
		t.Fatalf("MarkerExists returned error: %v", err)
	}
	if ok {
		t.Fatal("expected MarkerExists to be false")
	}
}

func TestDownloadMarkerExistsFalse(t *testing.T) {
	root := t.TempDir()

	dl := &Download{
		Directory: "nextcloud",
		Marker:    "config/config.php",
	}

	ok, err := dl.MarkerExists(root)
	if err != nil {
		t.Fatalf("MarkerExists returned error: %v", err)
	}
	if ok {
		t.Fatal("expected MarkerExists to be false")
	}
}

func TestDownloadMarkerExistsTrue(t *testing.T) {
	root := t.TempDir()

	dl := &Download{
		Directory: "nextcloud",
		Marker:    "config/config.php",
	}

	path := dl.MarkerPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(path, []byte("ok"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	ok, err := dl.MarkerExists(root)
	if err != nil {
		t.Fatalf("MarkerExists returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected MarkerExists to be true")
	}
}

func TestDownloadExistsLocalFalse(t *testing.T) {
	root := t.TempDir()

	pf := &Platform{
		options: &Options{
			rootDir: root,
		},
	}

	dl := &Download{
		Filename: "missing.tar.gz",
	}

	ok, err := dl.ExistsLocal(pf)
	if err != nil {
		t.Fatalf("ExistsLocal returned error: %v", err)
	}
	if ok {
		t.Fatal("expected ExistsLocal to be false")
	}
}

func TestDownloadExistsLocalTrue(t *testing.T) {
	root := t.TempDir()

	pf := &Platform{
		options: &Options{
			rootDir: root,
		},
	}

	dl := &Download{
		Filename: "present.tar.gz",
	}

	path := dl.LocalPath(pf)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	ok, err := dl.ExistsLocal(pf)
	if err != nil {
		t.Fatalf("ExistsLocal returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected ExistsLocal to be true")
	}
}

func TestHashes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	md5sum, sha256sum, sha512sum, err := Hashes(path)
	if err != nil {
		t.Fatalf("Hashes returned error: %v", err)
	}

	if md5sum != "900150983cd24fb0d6963f7d28e17f72" {
		t.Fatalf("unexpected MD5: %q", md5sum)
	}
	if sha256sum != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("unexpected SHA256: %q", sha256sum)
	}
	if sha512sum != "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a"+
		"2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f" {
		t.Fatalf("unexpected SHA512: %q", sha512sum)
	}
}

func TestVerifyFileOKWithMD5(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dl := &Download{
		Filename: "test.txt",
		MD5:      "900150983cd24fb0d6963f7d28e17f72",
	}

	if err := dl.VerifyFile(path); err != nil {
		t.Fatalf("VerifyFile returned error: %v", err)
	}
}

func TestVerifyFileOKWithSHA256(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dl := &Download{
		Filename: "test.txt",
		SHA256:   "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
	}

	if err := dl.VerifyFile(path); err != nil {
		t.Fatalf("VerifyFile returned error: %v", err)
	}
}

func TestVerifyFileOKWithSHA512(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dl := &Download{
		Filename: "test.txt",
		SHA512: "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a" +
			"2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f",
	}

	if err := dl.VerifyFile(path); err != nil {
		t.Fatalf("VerifyFile returned error: %v", err)
	}
}

func TestVerifyFileMD5Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dl := &Download{
		Filename: "test.txt",
		MD5:      "deadbeef",
	}

	err := dl.VerifyFile(path)
	if err == nil {
		t.Fatal("expected VerifyFile to fail")
	}
	if got, want := err.Error(), "MD5 mismatch for test.txt"; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}

func TestVerifyFileSHA256Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dl := &Download{
		Filename: "test.txt",
		SHA256:   "deadbeef",
	}

	err := dl.VerifyFile(path)
	if err == nil {
		t.Fatal("expected VerifyFile to fail")
	}
	if got, want := err.Error(), "SHA256 mismatch for test.txt"; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}

func TestVerifyFileSHA512Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dl := &Download{
		Filename: "test.txt",
		SHA512:   "deadbeef",
	}

	err := dl.VerifyFile(path)
	if err == nil {
		t.Fatal("expected VerifyFile to fail")
	}
	if got, want := err.Error(), "SHA512 mismatch for test.txt"; got != want {
		t.Fatalf("unexpected error: got %q want %q", got, want)
	}
}

func TestVerifyLocalOK(t *testing.T) {
	root := t.TempDir()

	pf := &Platform{
		options: &Options{
			rootDir: root,
		},
	}

	dl := &Download{
		Filename: "test.txt",
		SHA256:   "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
	}

	path := dl.LocalPath(pf)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(path, []byte("abc"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	if err := dl.VerifyLocal(pf); err != nil {
		t.Fatalf("VerifyLocal returned error: %v", err)
	}
}

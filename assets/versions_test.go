package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestCatalog(t *testing.T, dir string) {
	t.Helper()

	data := `{
  "products": [
    {
      "name": "wordpress",
      "default": "6.5",
      "versions": [
        { "key": "6.4" },
        { "key": "6.5" },
        { "key": "6.6" }
      ]
    }
  ]
}`

	path := filepath.Join(dir, "templates", "system")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		filepath.Join(path, "releases.json"),
		[]byte(data),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetVersions(t *testing.T) {
	dir := t.TempDir()

	SetTestDirs(Directories{
		rootDir:  dir,
		varDir:   dir,
		toolsDir: dir,
		etcDir:   dir,
	})

	writeTestCatalog(t, dir)

	list, err := GetVersions("wordpress")
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(list))
	}

	if list[0] != "6.4" {
		t.Fatal("wrong order")
	}
}

func TestGetDefaultVersion(t *testing.T) {

	dir := t.TempDir()

	SetTestDirs(Directories{
		rootDir:  dir,
		varDir:   dir,
		toolsDir: dir,
		etcDir:   dir,
	})

	writeTestCatalog(t, dir)

	def, err := GetDefaultVersion("wordpress")
	if err != nil {
		t.Fatal(err)
	}

	if def != "6.5" {
		t.Fatalf("expected 6.5 got %s", def)
	}
}

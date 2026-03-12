package assets

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {

	data := struct {
		Name string
	}{
		Name: "world",
	}

	out, err := Render("test/render.tmpl", data)
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(out))
	want := "hello world"

	if got != want {
		t.Fatalf("expected %q got %q", want, got)
	}
}

func TestSQL(t *testing.T) {

	stmts, err := SQL("test/sql.sql", nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements got %d", len(stmts))
	}

	if !strings.Contains(stmts[0], "a;b") {
		t.Fatalf("semicolon in string lost: %s", stmts[0])
	}
}

func TestLines(t *testing.T) {

	lines, err := Lines("test/lines.txt", "#", nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines got %d", len(lines))
	}

	if lines[0] != "alpha" {
		t.Fatalf("unexpected first line %q", lines[0])
	}
}

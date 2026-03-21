package config

import (
	"reflect"
	"testing"
)

func TestConfigRenderUsesOverride(t *testing.T) {
	cfg := &Config{}

	gotCalled := false
	cfg.render = func(name string, data any) ([]byte, error) {
		gotCalled = true

		if name != "test.txt" {
			t.Fatalf("unexpected name: got %q", name)
		}

		m, ok := data.(map[string]string)
		if !ok {
			t.Fatalf("unexpected data type: %T", data)
		}

		if m["name"] != "volker" {
			t.Fatalf("unexpected data content: %#v", m)
		}

		return []byte("ok"), nil
	}

	got, err := cfg.Render("test.txt", map[string]string{"name": "volker"})
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if !gotCalled {
		t.Fatal("expected render override to be called")
	}

	if string(got) != "ok" {
		t.Fatalf("unexpected result: got %q", string(got))
	}
}

func TestConfigRenderJSONUsesOverride(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	cfg := &Config{}

	gotCalled := false
	cfg.renderJSON = func(name string, data any, v any) error {
		gotCalled = true

		if name != "test.json" {
			t.Fatalf("unexpected name: got %q", name)
		}

		m, ok := data.(map[string]string)
		if !ok {
			t.Fatalf("unexpected data type: %T", data)
		}

		if m["name"] != "volker" {
			t.Fatalf("unexpected data content: %#v", m)
		}

		p, ok := v.(*payload)
		if !ok {
			t.Fatalf("unexpected target type: %T", v)
		}

		p.Name = "done"
		return nil
	}

	var got payload
	err := cfg.RenderJSON("test.json", map[string]string{"name": "volker"}, &got)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	if !gotCalled {
		t.Fatal("expected renderJSON override to be called")
	}

	if got.Name != "done" {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestConfigRenderSQLUsesOverride(t *testing.T) {
	cfg := &Config{}

	gotCalled := false
	cfg.renderSQL = func(name string, data any) ([]string, error) {
		gotCalled = true

		if name != "schema.sql" {
			t.Fatalf("unexpected name: got %q", name)
		}

		if data.(string) != "x" {
			t.Fatalf("unexpected data: %#v", data)
		}

		return []string{
			"CREATE TABLE test (id INTEGER)",
			"INSERT INTO test (id) VALUES (1)",
		}, nil
	}

	got, err := cfg.RenderSQL("schema.sql", "x")
	if err != nil {
		t.Fatalf("RenderSQL returned error: %v", err)
	}

	if !gotCalled {
		t.Fatal("expected renderSQL override to be called")
	}

	want := []string{
		"CREATE TABLE test (id INTEGER)",
		"INSERT INTO test (id) VALUES (1)",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected statements:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestConfigRenderListUsesOverride(t *testing.T) {
	cfg := &Config{}

	gotCalled := false
	cfg.renderList = func(name, comment string, data any) ([]string, error) {
		gotCalled = true

		if name != "packages.txt" {
			t.Fatalf("unexpected name: got %q", name)
		}

		if comment != "#" {
			t.Fatalf("unexpected comment: got %q", comment)
		}

		if data.(string) != "x" {
			t.Fatalf("unexpected data: %#v", data)
		}

		return []string{"curl", "jq", "vim"}, nil
	}

	got, err := cfg.RenderList("packages.txt", "#", "x")
	if err != nil {
		t.Fatalf("RenderList returned error: %v", err)
	}

	if !gotCalled {
		t.Fatal("expected renderList override to be called")
	}

	want := []string{"curl", "jq", "vim"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected lines:\n got: %#v\nwant: %#v", got, want)
	}
}

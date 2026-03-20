package utils

import (
	"reflect"
	"testing"
)

func TestLineBufferAddInsertAppend(t *testing.T) {
	var buf LineBuffer

	buf.Add("b")
	buf.Insert("a")
	buf.Append("c", "d")

	want := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(buf.Lines(), want) {
		t.Fatalf("unexpected lines: got %v want %v", buf.Lines(), want)
	}
}

func TestLineBufferEnsure(t *testing.T) {
	var buf LineBuffer

	buf.Add("a")
	buf.Ensure("b")
	buf.Ensure("a")

	want := []string{"a", "b"}
	if !reflect.DeepEqual(buf.Lines(), want) {
		t.Fatalf("unexpected lines: got %v want %v", buf.Lines(), want)
	}
}

func TestLineBufferAddf(t *testing.T) {
	var buf LineBuffer

	buf.Addf("value=%d", 42)

	want := []string{"value=42"}
	if !reflect.DeepEqual(buf.Lines(), want) {
		t.Fatalf("unexpected lines: got %v want %v", buf.Lines(), want)
	}
}

func TestLineBufferLinesReturnsCopy(t *testing.T) {
	var buf LineBuffer
	buf.Add("a")

	lines := buf.Lines()
	lines[0] = "x"

	want := []string{"a"}
	if !reflect.DeepEqual(buf.Lines(), want) {
		t.Fatalf("internal lines must not be modified: got %v want %v", buf.Lines(), want)
	}
}

func TestLineBufferSort(t *testing.T) {
	var buf LineBuffer
	buf.Append("c", "a", "b")

	buf.Sort()

	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(buf.Lines(), want) {
		t.Fatalf("unexpected lines: got %v want %v", buf.Lines(), want)
	}
}

func TestLineBufferUnique(t *testing.T) {
	var buf LineBuffer
	buf.Append("b", "", "a", "b", "a", "c", "")

	buf.Unique()

	want := []string{"b", "a", "c"}
	if !reflect.DeepEqual(buf.Lines(), want) {
		t.Fatalf("unexpected lines: got %v want %v", buf.Lines(), want)
	}
}

func TestLineBufferNormalize(t *testing.T) {
	var buf LineBuffer
	buf.Append("b", "", "a", "b", "c", "a")

	buf.Normalize()

	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(buf.Lines(), want) {
		t.Fatalf("unexpected lines: got %v want %v", buf.Lines(), want)
	}
}

func TestLineBufferIsEqual(t *testing.T) {
	a := &LineBuffer{}
	b := &LineBuffer{}
	c := &LineBuffer{}

	a.Append("a", "b")
	b.Append("a", "b")
	c.Append("b", "a")

	if !a.IsEqual(b) {
		t.Fatal("expected buffers to be equal")
	}
	if a.IsEqual(c) {
		t.Fatal("expected buffers to be different")
	}
	if a.IsEqual(nil) {
		t.Fatal("expected buffer and nil to be different")
	}
}

func TestLineBufferText(t *testing.T) {
	var buf LineBuffer
	buf.Append("a", "b")

	got := buf.Text()
	want := "a\nb\n"

	if got != want {
		t.Fatalf("unexpected text: got %q want %q", got, want)
	}
}

func TestLineBufferTextEmpty(t *testing.T) {
	var buf LineBuffer

	got := buf.Text()
	want := ""

	if got != want {
		t.Fatalf("unexpected text: got %q want %q", got, want)
	}
}

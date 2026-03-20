package utils

import (
	"fmt"
	"sort"
	"strings"
)

// LineBuffer collects text lines.
type LineBuffer struct {
	lines []string
}

// Add appends a single line.
func (buf *LineBuffer) Add(line string) {
	buf.lines = append(buf.lines, line)
}

// Insert prepends a single line.
func (buf *LineBuffer) Insert(line string) {
	buf.lines = append([]string{line}, buf.lines...)
}

// Append appends multiple lines.
func (buf *LineBuffer) Append(lines ...string) {
	buf.lines = append(buf.lines, lines...)
}

// Ensure appends the line only if it is not already present.
func (buf *LineBuffer) Ensure(line string) {
	for _, cmp := range buf.lines {
		if line == cmp {
			return
		}
	}
	buf.lines = append(buf.lines, line)
}

// Addf formats and appends a single line.
func (buf *LineBuffer) Addf(format string, args ...any) {
	buf.lines = append(buf.lines, fmt.Sprintf(format, args...))
}

// Lines returns a copy of the internal lines.
func (buf *LineBuffer) Lines() []string {
	out := make([]string, len(buf.lines))
	copy(out, buf.lines)
	return out
}

// Sort sorts the lines in place.
func (buf *LineBuffer) Sort() {
	sort.Strings(buf.lines)
}

// Unique removes duplicate and empty lines while preserving first occurrence order.
func (buf *LineBuffer) Unique() {
	if len(buf.lines) == 0 {
		return
	}

	seen := make(map[string]struct{}, len(buf.lines))
	out := make([]string, 0, len(buf.lines))

	for _, line := range buf.lines {
		if line == "" {
			continue
		}
		if _, ok := seen[line]; ok {
			continue
		}
		seen[line] = struct{}{}
		out = append(out, line)
	}

	buf.lines = out
}

// Normalize removes duplicate and empty lines and sorts the result.
func (buf *LineBuffer) Normalize() {
	buf.Unique()
	buf.Sort()
}

// IsEqual reports whether both buffers contain exactly the same lines
// in the same order.
func (buf *LineBuffer) IsEqual(cmp *LineBuffer) bool {
	if cmp == nil {
		return false
	}
	if len(buf.lines) != len(cmp.lines) {
		return false
	}
	for i := range buf.lines {
		if buf.lines[i] != cmp.lines[i] {
			return false
		}
	}
	return true
}

// Text returns all lines joined by "\n" and terminated by a trailing newline.
// For an empty buffer it returns an empty string.
func (buf *LineBuffer) Text() string {
	if len(buf.lines) == 0 {
		return ""
	}
	return strings.Join(buf.lines, "\n") + "\n"
}

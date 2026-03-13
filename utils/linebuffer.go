package utils

import (
	"fmt"
	"strings"
)

// Simple helper function to collect lines
type LineBuffer struct {
	lines []string
}

func (lb *LineBuffer) Add(line string) {
	lb.lines = append(lb.lines, line)
}

func (lb *LineBuffer) Addf(format string, args ...any) {
	lb.lines = append(lb.lines, fmt.Sprintf(format, args...))
}

func (lb *LineBuffer) Lines() []string {
	return lb.lines
}

func (lb *LineBuffer) Text() string {
	return strings.Join(lb.lines, "\n") + "\n"
}

package utils

import (
	"fmt"
	"io"
)

const (
	RunPrefix   = "[run]"
	DebugPrefix = "[dbg] ###################"
)

type Log struct {
	Writer  io.Writer
	Verbose bool
}

func (l *Log) writer() io.Writer {
	if l == nil || l.Writer == nil {
		return io.Discard
	}
	return l.Writer
}

func (l *Log) Say(args ...any) {
	w := l.writer()

	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Fprintln(w, RunPrefix, line)
			}
		default:
			fmt.Fprintln(w, RunPrefix, fmt.Sprint(arg))
		}
	}
}

func (l *Log) Sayf(format string, args ...any) {
	fmt.Fprintln(l.writer(), RunPrefix, fmt.Sprintf(format, args...))
}

func (l *Log) Debug(args ...any) {
	if l == nil || !l.Verbose {
		return
	}

	w := l.writer()

	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Fprintln(w, DebugPrefix, line)
			}
		default:
			fmt.Fprintln(w, DebugPrefix, fmt.Sprint(arg))
		}
	}
}

func (l *Log) Debugf(format string, args ...any) {
	if l == nil || !l.Verbose {
		return
	}

	fmt.Fprintln(l.writer(), DebugPrefix, fmt.Sprintf(format, args...))
}

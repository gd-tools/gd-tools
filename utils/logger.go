package utils

import (
	"fmt"
	"io"
)

const (
	InfoPrefix  = "[info]"
	DebugPrefix = "[dbg]"
	ErrorPrefix = "[err]"
)

type Logger struct {
	Writer  io.Writer
	Verbose bool
}

func (l *Logger) writer() io.Writer {
	if l == nil || l.Writer == nil {
		return io.Discard
	}
	return l.Writer
}

func (l *Logger) Info(args ...any) {
	w := l.writer()

	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Fprintln(w, InfoPrefix, line)
			}
		default:
			fmt.Fprintln(w, InfoPrefix, fmt.Sprint(arg))
		}
	}
}

func (l *Logger) Infof(format string, args ...any) {
	fmt.Fprintln(l.writer(), InfoPrefix, fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(args ...any) {
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

func (l *Logger) Debugf(format string, args ...any) {
	if l == nil || !l.Verbose {
		return
	}

	fmt.Fprintln(l.writer(), DebugPrefix, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(args ...any) {
	w := l.writer()

	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Fprintln(w, ErrorPrefix, line)
			}
		default:
			fmt.Fprintln(w, ErrorPrefix, fmt.Sprint(arg))
		}
	}
}

func (l *Logger) Errorf(format string, args ...any) {
	fmt.Fprintln(l.writer(), ErrorPrefix, fmt.Sprintf(format, args...))
}

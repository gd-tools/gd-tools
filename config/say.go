package config

import (
	"fmt"
)

const (
	RunPrefix   = "[run]"
	DebugPrefix = "[dbg] ##########"
)

func (cfg *Config) Say(args ...any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Println(RunPrefix, line)
			}
		default:
			fmt.Println(RunPrefix, fmt.Sprint(arg))
		}
	}
}

func (cfg *Config) Sayf(format string, args ...any) {
	fmt.Println(RunPrefix, fmt.Sprintf(format, args...))
}

func (cfg *Config) Debug(args ...any) {
	if !cfg.Verbose {
		return
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, line := range v {
				fmt.Println(DebugPrefix, line)
			}
		default:
			fmt.Println(DebugPrefix, fmt.Sprint(arg))
		}
	}
}

func (cfg *Config) Debugf(format string, args ...any) {
	if !cfg.Verbose {
		return
	}

	fmt.Println(DebugPrefix, fmt.Sprintf(format, args...))
}

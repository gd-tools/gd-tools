package config

// Logger defines the logging methods used by Config.
type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
}

// Info logs an informational message.
func (cfg *Config) Info(args ...any) {
	if cfg != nil && cfg.logger != nil {
		cfg.logger.Info(args...)
	}
}

// Infof logs a formatted informational message.
func (cfg *Config) Infof(format string, args ...any) {
	if cfg != nil && cfg.logger != nil {
		cfg.logger.Infof(format, args...)
	}
}

// Debug logs a debug message if verbose mode is enabled.
func (cfg *Config) Debug(args ...any) {
	if cfg != nil && cfg.logger != nil && cfg.Verbose {
		cfg.logger.Debug(args...)
	}
}

// Debugf logs a formatted debug message if verbose mode is enabled.
func (cfg *Config) Debugf(format string, args ...any) {
	if cfg != nil && cfg.logger != nil && cfg.Verbose {
		cfg.logger.Debugf(format, args...)
	}
}

// Error logs an error message.
func (cfg *Config) Error(args ...any) {
	if cfg != nil && cfg.logger != nil {
		cfg.logger.Error(args...)
	}
}

// Errorf logs a formatted error message.
func (cfg *Config) Errorf(format string, args ...any) {
	if cfg != nil && cfg.logger != nil {
		cfg.logger.Errorf(format, args...)
	}
}

// RsyncFlags returns the rsync flags depending on verbosity.
func (cfg *Config) RsyncFlags() string {
	if cfg != nil && cfg.Verbose {
		return "-avz"
	}
	return "-avzq"
}

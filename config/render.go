package config

import (
	"github.com/gd-tools/gd-tools/render"
)

// Render loads a template from the gdt binary.
// In tests, this can be overridden via cfg.render.
func (cfg *Config) Render(name string, data any) ([]byte, error) {
	if cfg != nil {
		if fn := cfg.render; fn != nil {
			return fn(name, data)
		}
	}
	return render.Render(name, data)
}

// RenderJSON loads a JSON template from the gdt binary.
// In tests, this can be overridden via cfg.renderJSON.
func (cfg *Config) RenderJSON(name string, data any, v any) error {
	if cfg != nil {
		if fn := cfg.renderJSON; fn != nil {
			return fn(name, data, v)
		}
	}
	return render.RenderJSON(name, data, v)
}

// RenderSQL loads SQL code from the gdt binary.
// In tests, this can be overridden via cfg.renderSQL.
func (cfg *Config) RenderSQL(name string, data any) ([]string, error) {
	if cfg != nil {
		if fn := cfg.renderSQL; fn != nil {
			return fn(name, data)
		}
	}
	return render.RenderSQL(name, data)
}

// RenderList loads a textual list from the gdt binary.
// In tests, this can be overridden via cfg.renderList.
func (cfg *Config) RenderList(name, comment string, data any) ([]string, error) {
	if cfg != nil {
		if fn := cfg.renderList; fn != nil {
			return fn(name, comment, data)
		}
	}
	return render.RenderList(name, comment, data)
}

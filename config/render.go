package config

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

//go:embed templates
var templateFS embed.FS

func load(name string) ([]byte, error) {
	return templateFS.ReadFile("templates/" + name)
}

// Render loads a template from the gdt binary.
func Render(name string, data any) ([]byte, error) {
	content, err := load(name)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer

	err = tmpl.Execute(&result, data)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

// RenderJSON loads a JSON template from the gdt binary.
func RenderJSON(name string, data any) error {
	content, err := load(name)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", name, err)
	}

	return nil
}

// RenderSQL loads SQL code from the gdt binary.
func RenderSQL(name string, data any) ([]string, error) {
	content, err := Render(name, data)
	if err != nil {
		return nil, err
	}

	var stmts []string
	var buf strings.Builder

	inSingle := false
	inDouble := false

	for _, r := range string(content) {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case ';':
			if !inSingle && !inDouble {
				stmt := strings.TrimSpace(buf.String())
				if stmt != "" {
					stmts = append(stmts, normalizeSQL(stmt))
				}

				buf.Reset()
				continue
			}
		}

		buf.WriteRune(r)
	}

	stmt := strings.TrimSpace(buf.String())
	if stmt != "" {
		stmts = append(stmts, normalizeSQL(stmt))
	}

	return stmts, nil
}

func normalizeSQL(sql string) string {
	words := strings.Fields(sql)
	return strings.Join(words, " ")
}

// RenderList loads a textual list from the gdt binary.
// Where comment is usually something like '#'.
func RenderList(name, comment string, data any) ([]string, error) {
	content, err := Render(name, data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(content)
	scanner := bufio.NewScanner(reader)

	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if comment != "" && strings.HasPrefix(line, comment) {
			continue
		}

		lines = append(lines, line)
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

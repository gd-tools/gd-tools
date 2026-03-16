package platform

import (
	"bufio"
	"bytes"
	"embed"
	"strings"
	"text/template"
)

//go:embed templates
var templateFS embed.FS

func load(name string) ([]byte, error) {
	return templateFS.ReadFile("templates/" + name)
}

func Render(name string, data interface{}) ([]byte, error) {
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

func RenderSQL(name string, data interface{}) ([]string, error) {
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

func RenderLines(name, comment string, data interface{}) ([]string, error) {
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

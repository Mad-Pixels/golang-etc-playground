package utils

import (
	"bytes"
	"strings"
	"text/template"
)

// EscapeString escapes double quotes in strings for JSON embedding.
func EscapeString(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

// Execute template and return result.
func Execute(data string, values any) (string, error) {
	funcMap := template.FuncMap{
		"escape": EscapeString,
	}

	tpl, err := template.New("tpl").Funcs(funcMap).Parse(data)
	if err != nil {
		return "", err
	}
	var body bytes.Buffer
	if err = tpl.Execute(&body, values); err != nil {
		return "", err
	}
	return body.String(), nil
}

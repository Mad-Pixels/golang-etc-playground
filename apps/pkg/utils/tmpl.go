package utils

import (
	"bytes"
	"text/template"
)

// Execute template and return result.
func Execute(data string, values any) (string, error) {
	tpl, err := template.New("tpl").Parse(data)
	if err != nil {
		return "", err
	}
	var body bytes.Buffer
	if err := tpl.Execute(&body, values); err != nil {
		return "", err
	}
	return body.String(), nil
}

package entrypoint

import (
	"encoding/base64"
)

const playgroundNs = "playground"

type requestPlayground struct {
	Version string `json:"version"`
	Source  string `json:"source"`
}

func (r *requestPlayground) SourceDecode() (string, error) {
	b, err := base64.StdEncoding.DecodeString(r.Source)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

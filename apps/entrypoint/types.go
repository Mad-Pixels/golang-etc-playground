package entrypoint

const playgroundNs = "playground"

type requestPlayground struct {
	Version string `json:"version"`
	Source  []byte `json:"source"`
}

type playgroundTmpl struct {
	Name    string
	Version string
}

var playgroundSpec = `
{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "name": "{{ .Name }}"
    },
    "spec": {
        "activeDeadlineSeconds": 45,
		"restartPolicy": "Never",
        "containers": [
            {
                "name": "{{ .Name }}-container",
                "image": "golang:{{ .Version }}-alpine3.18",
				"securityContext": {
                    "runAsNonRoot": true,
                    "runAsUser": 1000
                },
                "resources": {
                    "requests": {
                        "cpu": "100m",
                        "memory": "200Mi"
                    },
                    "limits": {
                        "cpu": "200m",
                        "memory": "400Mi"
                    }
                }
            }
        ]
    }
}
`

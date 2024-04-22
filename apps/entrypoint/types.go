package entrypoint

const playgroundNs = "playground"

type requestPlayground struct {
	Version string `json:"version"`
	Source  string `json:"source"`
	Name    string
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
                "env": [
                    {
                        "name": "GOCACHE",
                        "value": "/workspace/.cache"
                    }
                ],
                "resources": {
                    "requests": {
                        "cpu": "100m",
                        "memory": "200Mi"
                    },
                    "limits": {
                        "cpu": "200m",
                        "memory": "400Mi"
                    }
                },
				"command": ["/bin/sh", "-c"],
				"args": ["echo '{{ .Source | escape }}' > main.go && go build main.go && ./main"],
				"volumeMounts": [
					{
                        "name": "playground-storage",
                        "mountPath": "/workspace"
                    }
				]
            }
        ],
		"volumes": [
            {
                "name": "playground-storage",
                "emptyDir": {}
            }
        ]
    }
}
`

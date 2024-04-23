package entrypoint

const playgroundNs = "playground"

type requestPlayground struct {
	Version string `json:"version"`
	Source  string `json:"source"`
	Name    string
}

var playgroundMapSpec = `
{
  "apiVersion": "v1",
  "kind": "ConfigMap",
  "metadata": {
    "name": "{{ .Name }}"
  },
  "data": {
    "main.go": "{{ .Source | escape }}"
  }
}
`

var playgroundPodSpec = `
{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "name": "{{ .Name }}"
    },
    "spec": {
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
                        "value": "/tmp/.cache"
                    }
                ],
                "resources": {
                    "requests": {
                        "cpu": "1000m",
                        "memory": "400Mi"
                    },
                    "limits": {
                        "cpu": "1600m",
                        "memory": "600Mi"
                    }
                },
				"command": ["go", "run", "/workspace/main.go"],
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
                "configMap": {
					"name": "f1280140-3334-4261-86ef-98fd40f94a70"
				}
            }
        ]
    }
}
`

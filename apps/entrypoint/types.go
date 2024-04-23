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
        "activeDeadlineSeconds": 60,
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
                        "cpu": "100m",
                        "memory": "200Mi"
                    },
                    "limits": {
                        "cpu": "200m",
                        "memory": "400Mi"
                    }
                },
				"command": ["du", "/workspace/main.go"],
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
					"name": "{{ .Name }}",
					"items": [
						{
              				"key": "main.go",
              				"path": "main.go"
            			}
					]
				}
            }
        ]
    }
}
`

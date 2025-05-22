<picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/Mad-Pixels/.github/raw/main/profile/banner.png">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/Mad-Pixels/.github/raw/main/profile/banner.png">
    <img
        alt="MadPixels"
        src="https://github.com/Mad-Pixels/.github/raw/main/profile/banner.png">
</picture>

# 🧪 Golang ETC Playground
Golang ETC Playground is a lightweight web service (for fun) for running Golang code inside Kubernetes. The user submits a base64-encoded Go source, and the backend spawns a Pod using Alpine + Go, mounts the main.go file, runs it, and returns the output.

# Contributing
We're open to any new ideas and contributions. We also have some rules and taboos here, so please read this page and our [Code of Conduct](/CODE_OF_CONDUCT.md) carefully.

## I want to report an issue
If you've found an issue and want to report it, please check our [Issues](https://github.com/Mad-Pixels/golang-etc-playground/issues) page.

## ✨ Features
- Run Go code remotely inside Kubernetes Pods.
- Supports multiple Golang versions (e.g., 1.21, 1.22).
- Dynamic Pod and ConfigMap creation per request.
- Output logs returned to the client.
- Automatic cleanup of all temporary resources.
- Works seamlessly with Minikube.

## 🛠️ Tech Stack
- Language: Go
- Framework: Chi
- Kubernetes Client: client-go
- DevOps: Docker, Minikube, Helm, Taskfile
- Logging: Zap
- Monitoring: Prometheus-compatible metrics

## ⚙️ Installation & Usage
🔧 Requirements
- Docker
- Minikube
- Task (https://taskfile.dev)
- Helm

## 🚀 Quick Start
```bash
# Start Minikube from scratch
task minikube/up

# Build and deploy the service into Minikube
APP=entrypoint TAG=dev ARCH=amd64 task apps/minikube/build
task minikube/deploy/apps

# Open access to the app
task minikube/apps/entrypoint-lookup
```

## 📤 Example Request
POST /api/v1/playground
```json
{
  "version": "1.21",
  "source": "Ly8gSGVsbG8gd29ybGQKbWFpbi5nbyDwn5iA" // base64 of main.go
}
```
Response:
```json
{
  "data": "Hello world\n",
  "message": "",
  "host": "entrypoint-xxxx",
  "id": "req-123",
  "status": 200
}
```

## 🧪 Probes
- GET /api/internal/probe/liveness
- GET /api/internal/probe/readiness

## 🐳 Docker Build
```bash
APP=entrypoint TAG=latest ARCH=amd64 task apps/docker/build
```

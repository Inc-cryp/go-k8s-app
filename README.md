# Go Kubernetes App

A simple production-style Go web service built with Go, containerized using Docker, and deployed with Kubernetes.

## Tech Stack

- Go
- Docker
- Kubernetes

## Features

- Simple HTTP API
- Health check endpoint
- Version endpoint
- Environment-based configuration
- Multi-stage Docker build
- Kubernetes Deployment and Service
- ConfigMap integration
- Secret integration

## Project Structure

```bash
go-k8s-app/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── handler.go
│   └── response/
│       └── response.go
│
├── k8s/
│   ├── configmap.yaml
│   ├── deployment.yaml
│   ├── secret.yaml
│   └── service.yaml
│
├── .dockerignore
├── .gitignore
├── Dockerfile
├── go.mod
└── README.md
```

## API Endpoints

### `GET /`

Returns a simple welcome message.

#### Example Response

```json
{
  "message": "Hello from Production Kubernetes 🚀",
  "app": "go-k8s-app",
  "env": "production"
}
```

---

### `GET /health`

Used for health checking.

#### Example Response

```json
{
  "status": "ok"
}
```

---

### `GET /version`

Returns application version information.

#### Example Response

```json
{
  "app_name": "go-k8s-app",
  "version": "v1.0.0",
  "environment": "production"
}
```

---

## Run Locally

You can run the application locally by passing environment variables directly:

```bash
APP_NAME=my-app APP_VERSION=v2.1.0 APP_ENV=local API_KEY=test-key go run ./cmd/server
```

Then test it:

```bash
curl localhost:8080/
curl localhost:8080/health
curl localhost:8080/version
```

---

## Build Docker Image

```bash
docker build -t yourusername/go-k8s-app:v1 .
```

---

## Run Docker Container

```bash
docker run -p 8080:8080 \
  -e APP_NAME=go-docker-app \
  -e APP_VERSION=v1.0.1 \
  -e APP_ENV=docker \
  -e API_KEY=secret123 \
  yourusername/go-k8s-app:v1
```

Then test it:

```bash
curl localhost:8080/
curl localhost:8080/health
curl localhost:8080/version
```

---

## Deploy to Kubernetes

Apply all Kubernetes manifests:

```bash
kubectl apply -f k8s/
```

Check resources:

```bash
kubectl get pods
kubectl get svc
kubectl get configmap
kubectl get secret
```

---

## Kubernetes Resources Used

This project uses the following Kubernetes objects:

- Deployment
- Service
- ConfigMap
- Secret

---

## Environment Variables

The application supports the following environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_NAME` | Application name | `go-k8s-app` |
| `APP_VERSION` | Application version | `v1.0.0` |
| `APP_ENV` | Running environment | `production` |
| `API_KEY` | Example secret key | `super-secret-api-key` |

---

## Notes

This repository uses **demo-only configuration values and secrets** for educational and portfolio purposes.

Do **not** store real production secrets in Git repositories.

---

## Learning Goals

This project was built to practice:

- Structuring a Go backend service
- Writing a production-style Dockerfile
- Understanding Docker image build flow
- Deploying an application to Kubernetes
- Managing application configuration with ConfigMap and Secret

---

## Future Improvements

Possible next improvements for this project:

- Add CI/CD with GitHub Actions
- Add Ingress
- Add Horizontal Pod Autoscaler (HPA)
- Add graceful shutdown
- Add structured logging
- Add unit tests
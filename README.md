# Go Kubernetes App

A small production-style Go web service: containerized with a multi-stage
Docker build and deployed to Kubernetes.

## Tech Stack

- Go (standard library only, no external dependencies)
- Docker
- Kubernetes

## Features

- JSON HTTP API with a documented 404 for unknown paths
- Health check endpoint
- Version endpoint
- Environment-based configuration
- Graceful shutdown on SIGINT/SIGTERM
- Multi-stage Docker build producing a static, non-root container
- Kubernetes Deployment, Service, ConfigMap and Secret
- Unit tests and a GitHub Actions pipeline

## Project Structure

```bash
go-k8s-app/
├── cmd/
│   └── server/
│       └── main.go          # entrypoint: routing, timeouts, graceful shutdown
│
├── internal/
│   ├── config/
│   │   └── config.go        # environment-based configuration
│   ├── handler/
│   │   └── handler.go       # HTTP handlers
│   └── response/
│       └── response.go      # JSON response helper
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
├── Makefile
├── go.mod
└── README.md
```

## API Endpoints

Every response is JSON. Any path that is not one of the three below returns
`404` with `{"error":"not found","path":"..."}`.

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

Used by the liveness and readiness probes.

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

## Configuration

| Variable      | Description                          | Default         |
| ------------- | ------------------------------------ | --------------- |
| `APP_NAME`    | Application name                     | `go-k8s-app`    |
| `APP_VERSION` | Application version                  | `v1.0.0`        |
| `APP_ENV`     | Running environment                  | `development`   |
| `API_KEY`     | Example secret key                   | `default-secret`|
| `ADDR`        | Listen address                       | `:8080`         |

An empty variable is treated as unset, which matters because Kubernetes
delivers a missing ConfigMap key as an empty string.

---

## Run Locally

```bash
make run
```

Or with explicit values:

```bash
APP_NAME=my-app APP_VERSION=v2.1.0 APP_ENV=local ADDR=:8080 go run ./cmd/server
```

Then test it:

```bash
curl localhost:8080/
curl localhost:8080/health
curl localhost:8080/version
```

The server logs `shutdown signal received, draining connections` and exits
with status 0 on `Ctrl-C`.

---

## Make Targets

Run `make help` for the full list.

| Target         | Action                                            |
| -------------- | ------------------------------------------------- |
| `make check`   | Everything CI runs: formatting, vet, race tests    |
| `make test`    | Run the test suite                                |
| `make cover`   | Race tests with a coverage summary                |
| `make build`   | Build the static server binary                    |
| `make smoke`   | Build and probe the running server's endpoints    |
| `make docker-build` | Build the container image                    |
| `make k8s-dry-run`  | Validate the manifests against the cluster    |

---

## Build Docker Image

```bash
make docker-build
# or
docker build -t go-k8s-app:v1 .
```

The builder stage is pinned to the Go release in `go.mod`. The official
`golang` images set `GOTOOLCHAIN=local`, so that image's Go is final: pinning
an older release fails the build rather than downloading a newer toolchain.

---

## Run Docker Container

```bash
docker run -p 8080:8080 \
  -e APP_NAME=go-docker-app \
  -e APP_VERSION=v1.0.1 \
  -e APP_ENV=docker \
  -e API_KEY=secret123 \
  go-k8s-app:v1
```

The container runs as a non-root user and includes a `HEALTHCHECK` that
probes `/health`.

---

## Deploy to Kubernetes

Push the image to a registry you control and update `image:` in
`k8s/deployment.yaml` first — it ships with a placeholder value on purpose.
An image reference whose first path segment contains uppercase letters is
parsed as a registry hostname, so `Inc-cryp/go-k8s-app:v1` is not pullable.

Then apply the manifests:

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

- Deployment (3 replicas, liveness and readiness probes, non-root,
  `terminationGracePeriodSeconds: 30` to cover the app's 10s drain)
- Service (NodePort)
- ConfigMap
- Secret

---

## Notes

This repository uses **demo-only configuration values and secrets** for
educational and portfolio purposes.

Do **not** store real production secrets in Git repositories.

---

## Testing

```bash
make check
```

Covers the configuration loader, all four handlers, the JSON response helper
and the route table — including a regression test asserting that unknown
paths return 404 instead of falling through to the home handler.

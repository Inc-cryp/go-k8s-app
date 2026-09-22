GO     ?= go
BINARY ?= server
CMD    ?= ./cmd/server
IMAGE  ?= ghcr.io/inc-cryp/go-k8s-app:dev

.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: fmt
fmt: ## Format all Go source
	$(GO) fmt ./...

.PHONY: fmt-check
fmt-check: ## Fail if any file is not gofmt-clean
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then \
		echo "gofmt reported unformatted files:"; \
		echo "$$out"; \
		exit 1; \
	fi

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: fmt-check vet ## Run all static checks

.PHONY: test
test: ## Run the test suite
	$(GO) test ./...

.PHONY: test-race
test-race: ## Run the test suite with the race detector
	$(GO) test -race -count=1 ./...

.PHONY: cover
cover: ## Run the test suite and print per-package coverage
	$(GO) test -race -count=1 -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -n 1

.PHONY: check
check: lint test-race ## Run everything CI runs

.PHONY: build
build: ## Build the server binary
	CGO_ENABLED=0 $(GO) build -trimpath -o $(BINARY) $(CMD)

.PHONY: run
run: ## Run the server locally
	$(GO) run $(CMD)

.PHONY: smoke
smoke: build ## Build and probe the running server's endpoints
	@ADDR=127.0.0.1:18080 ./$(BINARY) & \
	pid=$$!; \
	trap 'kill $$pid 2>/dev/null' EXIT; \
	for i in $$(seq 1 50); do \
		if curl -sf http://127.0.0.1:18080/health >/dev/null 2>&1; then break; fi; \
		sleep 0.1; \
	done; \
	echo "GET /health  -> $$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:18080/health)"; \
	echo "GET /version -> $$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:18080/version)"; \
	echo "GET /        -> $$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:18080/)"; \
	echo "GET /missing -> $$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:18080/missing) (want 404)"; \
	kill -TERM $$pid; \
	wait $$pid 2>/dev/null; \
	echo "shutdown exit code: $$?"

.PHONY: docker-build
docker-build: ## Build the container image
	docker build -t $(IMAGE) .

.PHONY: k8s-dry-run
k8s-dry-run: ## Validate the Kubernetes manifests against the cluster schema
	kubectl apply --dry-run=client -f k8s/

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out

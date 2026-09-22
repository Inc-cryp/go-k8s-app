# ---------- STAGE 1: BUILD ----------
# Keep this tag in sync with the `go` directive in go.mod. The official
# golang images set GOTOOLCHAIN=local, so the image's Go is final: pinning a
# release older than go.mod's minimum fails the build instead of downloading
# a newer toolchain.
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Dependencies are resolved from go.mod alone, so this layer is only
# invalidated when the dependency set actually changes, not on every edit.
COPY go.mod ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 produces a static binary that runs on a bare Alpine base.
# -trimpath keeps local paths out of the binary, and the -s -w linker flags
# drop the symbol table and DWARF data.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/server \
    ./cmd/server

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.22

# wget is used by the container healthcheck.
RUN apk add --no-cache ca-certificates wget \
    && adduser -D -H -u 10001 appuser

WORKDIR /app

COPY --from=builder /out/server /app/server

USER appuser

EXPOSE 8080

# The service has no upstream dependencies, so the process being able to
# answer its own health endpoint is a meaningful liveness signal.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/health || exit 1

CMD ["/app/server"]

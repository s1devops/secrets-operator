# Build the binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY secrets/ secrets/

# Build
RUN CGO_ENABLED=0 go build -a -o secrets-operator main.go

FROM debian:buster

RUN set -eux; \
    apt-get update \
    && apt-get install -y \
        ca-certificates \
        gnupg \
    && apt-get clean \
    && apt-get autoremove --purge -y \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=builder /workspace/secrets-operator /usr/bin/
RUN useradd -ms /bin/bash secret-operator
USER secret-operator:65532

ENTRYPOINT ["/usr/bin/secrets-operator"]

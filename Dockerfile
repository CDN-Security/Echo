FROM golang:1.25 AS builder

RUN go install github.com/goreleaser/goreleaser/v2@latest
WORKDIR /app

# Copy source code
COPY . .

# Check if the current commit is tagged, and build accordingly
RUN if git describe --tags --exact-match >/dev/null 2>&1; then \
      echo "Commit is tagged. Creating a release build."; \
      goreleaser build --clean; \
    else \
      echo "Commit is not tagged. Creating a snapshot build."; \
      goreleaser build --clean --snapshot; \
    fi

FROM alpine:3.14

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/dist/http-echo-server_linux_amd64_v1/http-echo-server /usr/local/bin/http-echo-server
COPY --from=builder /app/config.example.yaml /etc/http-echo-server/config.yaml

ENTRYPOINT [ "/usr/local/bin/http-echo-server", "--config", "/etc/http-echo-server/config.yaml" ]

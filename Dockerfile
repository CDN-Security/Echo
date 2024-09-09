FROM golang:1.22 AS builder

RUN go install github.com/goreleaser/goreleaser@latest
WORKDIR /app
COPY ./.git/ ./.git/
RUN git reset --hard HEAD
RUN goreleaser build --clean --id=http-echo-server --snapshot

FROM alpine:3.14
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/dist/http-echo-server_linux_amd64_v1/http-echo-server /usr/local/bin/http-echo-server
COPY --from=builder /app/config.example.yaml /etc/http-echo-server/config.yaml
ENTRYPOINT [ "/usr/local/bin/http-echo-server", "--config", "/etc/http-echo-server/config.yaml" ]

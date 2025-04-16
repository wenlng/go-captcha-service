# Build phase
FROM --platform=$BUILDPLATFORM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o go-captcha-service ./cmd/go-captcha-service

# Run phase (default binary)
FROM scratch AS binary

WORKDIR /app

COPY --from=builder /app/go-captcha-service .
COPY config.json .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080 50051

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/go-captcha-service", "--health-check"] || exit 1

CMD ["/app/go-captcha-service"]

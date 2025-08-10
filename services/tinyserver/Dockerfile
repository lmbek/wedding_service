# syntax=docker/dockerfile:1

FROM golang:1.24.4 AS builder
WORKDIR /app/services/tinyserver
# Use wildcard to copy any go.* files for caching
COPY services/tinyserver/go.* ./
RUN go mod download || true
# Copy repository for module context
COPY . /app
# Build tinyserver
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/tinyserver .

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/tinyserver /app/tinyserver
EXPOSE 80
ENTRYPOINT ["/app/tinyserver"]

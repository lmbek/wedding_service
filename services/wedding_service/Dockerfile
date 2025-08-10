# syntax=docker/dockerfile:1

FROM golang:1.24.6 AS builder
# Build from the wedding_service module directory
WORKDIR /app/services/wedding_service
# Copy module files first for caching
COPY services/wedding_service/go.* ./
RUN go mod download || true
# Copy the whole repository so the replace path (../..) exists in build stage
COPY . /app
# Build the wedding_service binary from this module
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags dockerdev -o /out/wedding_service .

# Use Alpine for the runtime image
FROM alpine:3.22
WORKDIR /app
# Install CA certificates for HTTPS outbound requests, if any
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /out/wedding_service /app/wedding_service
# Note: Frontend assets are mounted at runtime via docker-compose bind mount
ENV GIN_MODE=release

ENTRYPOINT ["/app/wedding_service"]

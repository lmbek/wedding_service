# syntax=docker/dockerfile:1

FROM golang:1.24.4 AS builder
# Build from the wedding_service module directory
WORKDIR /app/services/wedding_service
# Copy module files first for caching
COPY services/wedding_service/go.* ./
RUN go mod download || true
# Copy the whole repository so the replace path (../..) exists in build stage
COPY . /app
# Build the wedding_service binary from this module
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/wedding_service .

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/wedding_service /app/wedding_service
# Static assets
COPY --from=builder /app/services/wedding_service/webserver/website/frontend/out /app/webserver/website/frontend/out
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["/app/wedding_service"]

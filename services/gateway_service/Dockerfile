# syntax=docker/dockerfile:1

FROM golang:1.24.4 AS builder
WORKDIR /app/services/gateway_service
COPY services/gateway_service/go.* ./
RUN go mod download || true
COPY services/gateway_service/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/gateway_service .

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/gateway_service /app/gateway_service
# Include self-signed certificates so default CERT_PATH/KEY_PATH work out of the box
COPY --from=builder /app/services/gateway_service/certificate/self_sign_cert /app/services/gateway_service/certificate/self_sign_cert
# Certificates and other runtime data under /data (mounted volume)
VOLUME ["/data"]
EXPOSE 80 443
ENTRYPOINT ["/app/gateway_service"]

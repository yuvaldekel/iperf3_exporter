FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o iperf3_exporter ./cmd/iperf3_exporter

FROM alpine:3.23
RUN apk add --no-cache iperf3
USER iperf3
WORKDIR /exporter
COPY --from=builder /app/iperf3_exporter .
ENTRYPOINT ["/exporter/iperf3_exporter"]

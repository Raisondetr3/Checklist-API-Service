FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o checklist-api-service ./cmd/api-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata wget
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN mkdir -p logs && chown -R appuser:appgroup /app

COPY --from=builder /app/checklist-api-service .

RUN chmod +x checklist-api-service
USER appuser
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./checklist-api-service"]
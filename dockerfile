# Stage 1: Build
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ internal/
COPY cmd/ cmd/

RUN CGO_ENABLED=1 GOOS=linux go build -o model-registry ./cmd/model-registry


# Stage 2: Runtime
FROM alpine:3.21.3

WORKDIR /app

RUN addgroup -S appgroup && adduser -S -G appgroup appuser

COPY --from=builder /app/model-registry .

COPY --chmod=100 entrypoint.sh /entrypoint.sh

RUN mkdir -p /app/storage /app/config /app/models

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
CMD ["./model-registry"]
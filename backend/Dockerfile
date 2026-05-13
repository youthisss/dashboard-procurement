# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server /app/server

RUN addgroup -S app && adduser -S -G app -h /app app \
    && mkdir -p /app/uploads \
    && chown -R app:app /app

USER app

EXPOSE 8080

CMD ["./server"]

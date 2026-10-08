FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -s /sbin/nologin appuser

WORKDIR /app
COPY --from=builder /app/server ./server

USER appuser
EXPOSE 8080

CMD ["./server"]
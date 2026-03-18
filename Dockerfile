FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /ssh-portfolio ./cmd/server

# ---

FROM alpine:3.21

RUN apk add --no-cache openssh-keygen ca-certificates

COPY --from=builder /ssh-portfolio /usr/local/bin/ssh-portfolio
COPY deploy/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN adduser -D -h /app portfolio
USER portfolio
WORKDIR /app

# SSH keys, DB stored here — mount a volume for persistence
RUN mkdir -p /app/.ssh /app/data

ENV SSH_HOST=0.0.0.0
ENV SSH_PORT=22
ENV HOST_KEY_DIR=/app/.ssh
ENV DB_PATH=/app/data/analytics.db

EXPOSE 22

ENTRYPOINT ["/entrypoint.sh"]
CMD ["ssh-portfolio"]

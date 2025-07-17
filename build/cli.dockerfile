FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
            go build -o cli-app ./cmd/main.go


FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY battleship-lesta-start.ru.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

COPY --from=builder /app/cli-app ./cli

ENTRYPOINT ["./cli"]

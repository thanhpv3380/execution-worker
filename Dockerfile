FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o main ./cmd

FROM golang:1.24-alpine

RUN addgroup -S sandboxgroup && adduser -S sandboxuser -G sandboxgroup

WORKDIR /home/sandboxuser
COPY --from=builder /app/main ./main
RUN chmod +x ./main

RUN mkdir /tmp && chown sandboxuser:sandboxgroup /tmp

USER sandboxuser

WORKDIR /home/sandboxuser

CMD ["./main"]

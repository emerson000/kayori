FROM golang:alpine AS builder

WORKDIR /app

COPY src /app/src
COPY go.mod go.sum /app/

RUN cd /app/src && \
    go mod tidy && \
    go build -o /app/server server.go

# Create a smaller image for the final executable
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 3001

CMD ["./server"]

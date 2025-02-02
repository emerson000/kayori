FROM golang:alpine AS builder

WORKDIR /app

COPY src /app/src
COPY go.mod go.sum /app/

RUN cd /app/src && \
    go mod tidy && \
    go build -o /app/server server.go

# Create a smaller image for the final executable
FROM golang:alpine

# Install Redis
RUN apk add --no-cache redis

WORKDIR /root/

COPY --from=builder /app/server .

COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

COPY redis.conf /etc/redis.conf

EXPOSE 3001
EXPOSE 6379

ENTRYPOINT [ "./entrypoint.sh" ]
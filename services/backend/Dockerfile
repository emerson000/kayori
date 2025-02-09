FROM golang:alpine AS builder

WORKDIR /app

COPY . /app

RUN cd /app && \
    go mod tidy && \
    go build -o /app/server cmd/server/main.go

# Create a smaller image for the final executable
FROM golang:alpine

# Install Redis
RUN apk add --no-cache redis

WORKDIR /root/

COPY --from=builder /app/server .

COPY entrypoint.sh /root/entrypoint.sh

RUN chmod +x /root/entrypoint.sh

COPY redis.conf /etc/redis.conf

EXPOSE 3001
EXPOSE 6379

ENTRYPOINT [ "/root/entrypoint.sh" ]
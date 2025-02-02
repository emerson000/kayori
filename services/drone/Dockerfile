FROM golang:alpine AS builder

WORKDIR /app

COPY . /app

RUN cd /app && \
    go mod tidy && \
    go build -o /app/drone ./cmd/drone/main.go

# Create a smaller image for the final executable
FROM golang:alpine

WORKDIR /root/

COPY --from=builder /app/drone .

COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]
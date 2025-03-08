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

COPY ./scripts/ /root/scripts

COPY ./entrypoint.sh .

RUN chmod +x entrypoint.sh

WORKDIR /root/scripts/

# Install Python
RUN apk add --no-cache python3 \
    py3-pip \
    py3-nltk \
    py3-scikit-learn \
    py3-numpy

RUN python3 -m venv .venv --system-site-packages

RUN .venv/bin/pip install --upgrade pip setuptools wheel

RUN .venv/bin/pip install -r ./requirements.txt

WORKDIR /root/

ENTRYPOINT [ "/root/entrypoint.sh" ]
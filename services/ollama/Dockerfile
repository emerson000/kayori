FROM ollama/ollama

EXPOSE 11434

VOLUME /root/.ollama

RUN apt-get update && apt-get install -y curl

RUN ollama serve & \
    until curl -s http://localhost:11434/ | grep "Ollama is running"; do \
    echo "Waiting for Ollama app to be ready..."; \
    sleep 5; \
    done && \
    ollama pull nomic-embed-text

RUN ollama serve & \
    until curl -s http://localhost:11434/ | grep "Ollama is running"; do \
    echo "Waiting for Ollama app to be ready..."; \
    sleep 5; \
    done && \
    ollama pull hermes3:3b

CMD ["serve"]
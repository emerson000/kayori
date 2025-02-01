package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello, Server1"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("New client connected!")
	conn.WriteMessage(websocket.TextMessage, []byte("Welcome to WebSocket server!"))

	// Kafka reader configuration
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "Test",
		GroupID: "wsHandlerGroup",
	})
	defer kafkaReader.Close()

	go func() {
		for {
			m, err := kafkaReader.ReadMessage(context.Background())
			if err != nil {
				if err == io.EOF {
					fmt.Println("Kafka read error: EOF")
					break
				}
				fmt.Println("Kafka read error:", err)
				break
			}
			fmt.Printf("Message from Kafka: %s\n", string(m.Value))
			conn.WriteMessage(websocket.TextMessage, m.Value)
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("WebSocket closed:", err)
				break
			}
			fmt.Println("Read error:", err)
			break
		}
		fmt.Println("Received:", string(message))
		conn.WriteMessage(websocket.TextMessage, append([]byte("Server echo: "), message...))
	}
}

func main() {
	http.HandleFunc("/api/users", usersHandler)
	http.HandleFunc("/ws", wsHandler)

	server := &http.Server{Addr: ":3001"}
	fmt.Println("Server listening on port 3001")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

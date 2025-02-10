package controllers

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var (
	connections = make(map[*websocket.Conn]chan []byte)
	connMutex   sync.Mutex
)

func WebSocketHandler(pubsub *redis.PubSub) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func WebSocketConnection(pubsub *redis.PubSub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		writeQueue := make(chan []byte, 10)
		connMutex.Lock()
		connections[c] = writeQueue
		connMutex.Unlock()
		defer func() {
			connMutex.Lock()
			delete(connections, c)
			connMutex.Unlock()
			c.Close()
		}()
		var writeMutex sync.Mutex
		// Goroutine to handle all writes to the websocket
		go func() {
			for msg := range writeQueue {
				writeMutex.Lock()
				if c == nil {
					log.Println("write: nil *Conn")
					writeMutex.Unlock()
					break
				}
				if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println("write:", err)
					writeMutex.Unlock()
					break
				}
				writeMutex.Unlock()
			}
		}()
		// Listen for messages from Redis and broadcast them to all connections
		go func() {
			ch := pubsub.Channel()
			for msg := range ch {
				connMutex.Lock()
				for conn, queue := range connections {
					select {
					case queue <- []byte(msg.Payload):
					default:
						log.Printf("dropping message for connection: %v", conn)
					}
				}
				connMutex.Unlock()
			}
		}()
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			writeMutex.Lock()
			if c == nil {
				log.Println("write: nil *Conn")
				writeMutex.Unlock()
				break
			}
			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				writeMutex.Unlock()
				break
			}
			writeMutex.Unlock()
		}
	})
}

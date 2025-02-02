package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"kayori.io/drone/internal/kafka"
	"kayori.io/drone/internal/providers/rss"
)

// Define the possible job types (this can come from your Kafka messages)
const (
	RssType = "rss"
)

// Mock message structure for demonstration
type Message struct {
	Type    string
	Payload interface{}
}

func main() {
	fmt.Println("Drone awaiting dispatch")

	consumer := kafka.NewConsumer()

	err := consumer.Start(func(value []byte) {
		processMessage(value)
	})
	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}
}

type DroneJob struct {
	Id      string      `json:"id"`
	Service ServiceType `json:"service"`
	Task    interface{} `json:"task"`
}

type ServiceType string

const (
	RssService ServiceType = "rss"
)

func processMessage(value []byte) {
	var msg DroneJob
	err := json.Unmarshal(value, &msg)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Starting job: %v", msg.Service)

	err = unmarshalTask(&msg, value)
	if err != nil {
		log.Printf("Error unmarshaling task: %v", err)
		return
	}
	switch msg.Service {
	case RssService:
		rss.ProcessTask(msg.Task.(*rss.Task), postJSON)
	}

	log.Printf("Processed job: %+v", msg.Service)
}

func postJSON(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	return nil
}

func unmarshalTask(msg *DroneJob, value []byte) error {
	switch msg.Service {
	case RssService:
		var raw struct {
			Task rss.Task `json:"task"`
		}
		if err := json.Unmarshal(value, &raw); err != nil {
			return fmt.Errorf("error unmarshaling task for service %s: %v", msg.Service, err)
		}
		msg.Task = &raw.Task
	}
	return nil
}

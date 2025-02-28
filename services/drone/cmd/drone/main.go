package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"kayori.io/drone/internal/data"
	"kayori.io/drone/internal/kafka"
	"kayori.io/drone/internal/providers/rss"
	"kayori.io/drone/internal/providers/utilities"
)

// Define the possible job types (this can come from your Kafka messages)
const (
	RssType     = "rss"
	Deduplicate = "deduplicate"
)

const maxRetries = 5

// Mock message structure for demonstration
type Message struct {
	Type    string
	Payload interface{}
}

func main() {
	fmt.Println("Drone awaiting dispatch")

	consumer := kafka.NewConsumer()
	retries := 0

	for {
		err := consumer.Start(func(value []byte) {
			if err := processMessage(value); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		})
		if err != nil {
			if retries >= maxRetries {
				log.Fatalf("Error starting consumer after %d retries: %v", maxRetries, err)
			}
			log.Printf("Error starting consumer: %v. Retrying in 30 seconds... (Attempt %d/%d)", err, retries+1, maxRetries)
			time.Sleep(30 * time.Second)
			retries++
			continue
		} else {
			log.Println("Consumer started successfully")
		}
		break
	}
}

type DroneJob struct {
	Id       string      `json:"id"`
	Service  ServiceType `json:"service"`
	Category string      `json:"category"`
	Task     interface{} `json:"task"`
}

type ServiceType string

const (
	RssService         ServiceType = "rss"
	DeduplicateService ServiceType = "deduplicate"
)

func processMessage(value []byte) error {
	var msg DroneJob
	err := json.Unmarshal(value, &msg)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return err
	}

	log.Printf("Starting job: %v", msg.Service)

	err = unmarshalTask(&msg, value)
	if err != nil {
		log.Printf("Error unmarshaling task: %v", err)
		return err
	}

	err = processServiceTask(msg)
	if err != nil {
		objectId, err1 := bson.ObjectIDFromHex(msg.Id)
		if err1 != nil {
			log.Printf("Error converting job ID to ObjectID: %v", err1)
			return err1
		}
		data.SetJobStatus(objectId, "failed")
	}

	log.Printf("Processed job: %+v", msg.Service)
	return nil
}

func processServiceTask(msg DroneJob) error {
	switch msg.Service {
	case RssService:
		return rss.ProcessTask(msg.Id, msg.Task.(*rss.Task), postJSON)
	case DeduplicateService:
		return utilities.ProcessTask(msg.Id, msg.Task.(*utilities.Task), postJSON)
	default:
		return fmt.Errorf("unknown service type: %v", msg.Service)
	}
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
	case DeduplicateService:
		var raw struct {
			Task utilities.Task `json:"task"`
		}
		if err := json.Unmarshal(value, &raw); err != nil {
			return fmt.Errorf("error unmarshaling task for service %s: %v", msg.Service, err)
		}
		msg.Task = &raw.Task
	}
	return nil
}

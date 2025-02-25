package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	kafkaLib "kayori.io/backend/data/kafka"
	"kayori.io/backend/data/mongorm"
	"kayori.io/backend/models"
	"kayori.io/backend/routes"
)

func main() {
	app := fiber.New()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	client, err := mongorm.Connect("mongodb://root:kayori@mongo:27017")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("kayori")

	// Subscribe to the "drone-status" topic
	pubsub := rdb.Subscribe(context.Background(), "drone-status")
	defer pubsub.Close()

	// Kafka writer
	kafkaWriter := kafkaLib.GetKafkaWriter()
	defer kafkaLib.CloseKafkaWriter()

	// Register routes
	routes.RegisterRoutes(app, db, rdb, kafkaWriter, pubsub)

	go startScheduler(db, kafkaWriter)

	app.Listen(":3001")
}

// startScheduler starts a scheduler that fetches jobs from the MongoDB collection and runs them at the appropriate time
func startScheduler(db *mongo.Database, kafkaWriter *kafka.Writer) {
	for {
		log.Println("Checking for jobs")
		jobs, err := fetchJobs(db)
		if err != nil {
			log.Printf("Error fetching jobs: %s", err)
			continue
		}

		for _, job := range jobs {
			log.Printf("Checking job %s", job.Title)
			duration := job.Schedule.Duration
			intervalType := job.Schedule.Interval
			var nextRun time.Time
			switch intervalType {
			case "hours":
				nextRun = job.Schedule.LastRun.Add(time.Duration(duration) * time.Hour)
			case "minutes":
				nextRun = job.Schedule.LastRun.Add(time.Duration(duration) * time.Minute)
			case "days":
				nextRun = job.Schedule.LastRun.AddDate(0, 0, duration)
			default:
				nextRun = job.Schedule.LastRun
			}

			if time.Now().After(nextRun) {
				log.Printf("Running job %s", job.Title)
				go runJob(job, kafkaWriter, db)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func fetchJobs(db *mongo.Database) ([]models.Job, error) {
	jobs := make([]models.Job, 0)
	if err := (&models.Job{}).ReadAll(context.Background(), db, "jobs", bson.M{"schedule.schedule": true}, &jobs, nil); err != nil {
		return nil, err
	}
	return jobs, nil
}

func runJob(job models.Job, kafkaWriter *kafka.Writer, db *mongo.Database) {
	job.Schedule.LastRun = time.Now()
	updateStatement := bson.M{
		"$set": bson.M{
			"schedule.last_run": job.Schedule.LastRun,
		},
	}
	if err := job.Update(context.Background(), db, "jobs", updateStatement); err != nil {
		log.Printf("Error updating job %s: %s", job.Title, err)
		return
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		log.Printf("Error marshalling job %s: %s", job.Title, err)
		return
	}

	if err := kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Value: jobData,
		},
	); err != nil {
		log.Printf("Error writing job %s to Kafka: %s", job.Title, err)
	}
}

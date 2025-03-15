package utilities

import (
	"encoding/json"
	"log"
	"sync"

	"kayori.io/drone/internal/data"
)

type Task struct {
	Jobs  []string `json:"jobs"`
	Field string   `json:"field"`
}

type artifact map[string]interface{}

func ProcessTask(jobId string, task *Task, postJSON func(url string, data interface{}) error) error {
	log.Printf("Starting deduplication job: %v", jobId)
	deduplicate := func(job string) {
		log.Printf("Deduplicating artifacts from job: %v using field %v", job, task.Field)
		path := "/api/jobs/" + job + "/artifacts"
		resp, err := data.Get(path)
		if err != nil {
			log.Printf("Error fetching artifacts: %v", err)
			return
		}
		var parsedJson []artifact
		err = json.Unmarshal([]byte(resp), &parsedJson)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			return
		}

		seenValues := make(map[string]bool)
		for _, art := range parsedJson {
			value := getFieldValue(art, task.Field)
			artifactId := art["id"].(string)
			if seenValues[value] {
				if value == "" {
					continue
				}
				log.Printf("Duplicate artifact found: %v, %v, %v", artifactId, seenValues[value], value)
				deletePath := "/api/artifacts/" + artifactId
				err := data.Delete(deletePath)
				if err != nil {
					log.Printf("Error deleting artifact %v: %v", artifactId, err)
				} else {
					log.Printf("Deleted duplicate artifact %v with value %v", artifactId, value)
				}
			} else {
				seenValues[value] = true
			}
		}
	}
	if len(task.Jobs) > 0 {
		var wg sync.WaitGroup
		for _, currentJob := range task.Jobs {
			wg.Add(1)
			go func(job string) {
				defer wg.Done()
				deduplicate(job)
			}(currentJob)
		}
		wg.Wait()
	}
	return nil
}

func getFieldValue(art artifact, field string) string {
	if value, ok := art[field]; ok {
		return value.(string)
	}
	return ""
}

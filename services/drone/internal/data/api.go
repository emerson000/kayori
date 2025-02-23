package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var hostname = GetHostname()

func GetHostname() string {
	if value, exists := os.LookupEnv("BACKEND_HOSTNAME"); exists {
		return value
	}
	return "http://backend:3001"
}

type JobUpdate struct {
	Status string `json:"status"`
}

func SetJobStatus(jobId bson.ObjectID, status string) error {
	if status != "failed" && status != "done" {
		return fmt.Errorf("invalid status: %s", status)
	}
	job := JobUpdate{Status: status}
	return put(fmt.Sprintf("/api/jobs/%s", jobId.Hex()), job)
}

func put(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, hostname+url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	return nil
}

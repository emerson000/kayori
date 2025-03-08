package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func Get(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, hostname+url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating GET request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

func Delete(url string) error {
	req, err := http.NewRequest(http.MethodDelete, hostname+url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	return nil
}

func Post(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, hostname+url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("JSON Request Body: %s", string(jsonData))
		return fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	return nil
}

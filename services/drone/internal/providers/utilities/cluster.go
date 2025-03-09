package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"kayori.io/drone/internal/data"
)

type ClusterTask struct {
	EntityType string `json:"entity_type"`
	After      int64  `json:"after"`
}

func ProcessClusterTask(jobId string, task *ClusterTask, postJSON func(url string, data interface{}) error) error {
	log.Printf("Starting cluster job: %v", jobId)
	entity_type := task.EntityType
	log.Printf("Processing cluster task for entity type: %v", entity_type)
	currentPage := 1
	limit := 100
	hasMore := true
	after := task.After
	if after == 0 {
		after = time.Now().Add(-168 * time.Hour).Unix()
	}
	// Cache to store unmarshaled JSON for each page
	pageCache := make(map[int]interface{})

	// Synchronously calculate the total number of pages and cache the responses
	for hasMore {
		log.Printf("Fetching data for page %d", currentPage)
		path := "/api/entities/" + entity_type + "?columns=[\"title\"]&limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(currentPage) + "&after=" + strconv.FormatInt(after, 10)
		body, resp, err := data.GetBodyAndResponse(path)
		if err != nil {
			log.Printf("Error fetching cluster data: %v", err)
			return err
		}
		if body == "" {
			log.Printf("No data found for entity type: %v", entity_type)
			break
		}
		var jsonData interface{}
		err = json.Unmarshal([]byte(body), &jsonData)
		if err != nil {
			log.Printf("Error unmarshalling JSON body: %v", err)
			return err
		}
		pageCache[currentPage] = jsonData
		hasMoreHeader := resp.Header.Get("x-has-more")
		hasMore, err = strconv.ParseBool(hasMoreHeader)
		if err != nil {
			log.Printf("Error parsing x-has-more header: %v", err)
			return err
		}
		currentPage++
	}

	// Join all pages in pageCache
	var allData []interface{}
	for i := 1; i < currentPage; i++ {
		pageData, ok := pageCache[i].([]interface{})
		if !ok {
			log.Printf("Error asserting page data to []interface{}")
			return fmt.Errorf("error asserting page data to []interface{}")
		}
		allData = append(allData, pageData...)
	}

	allDataBytes, err := json.Marshal(allData)
	if err != nil {
		log.Printf("Error marshalling all data: %v", err)
		return err
	}

	// Process the response
	cmd := exec.Command("scripts/.venv/bin/python3", "scripts/cluster_headlines.py")
	cmd.Stdin = bytes.NewReader(allDataBytes)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running Python script: %v", err)
		log.Printf("Stderr: %s", stderr.String())
		return err
	}

	var clusters map[string][]map[string]interface{}
	err = json.Unmarshal(out.Bytes(), &clusters)
	if err != nil {
		log.Printf("Error unmarshalling JSON output: %v", err)
		return err
	}

	delete(clusters, "-1")
	if len(clusters) == 0 {
		log.Printf("No clusters found in the output")
		return nil
	}
	for _, cluster := range clusters {
		ids := make([]string, 0)
		for _, item := range cluster {
			if id, ok := item["_id"]; ok {
				ids = append(ids, id.(string))
			}
		}
		body := map[string]interface{}{
			"artifacts": ids,
		}
		err := data.Post("/api/clusters", body)
		if err != nil {
			log.Printf("Error posting cluster data: %v", err)
			return err
		}
	}

	log.Printf("Total pages processed: %d", currentPage-1)
	return nil
}

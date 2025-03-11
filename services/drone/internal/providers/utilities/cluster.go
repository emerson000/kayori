package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	pageCache := make(map[int]interface{})

	for hasMore {
		body, resp, err := fetchData(entity_type, limit, currentPage, after)
		if err != nil {
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
		hasMore, err = parseHasMoreHeader(resp)
		if err != nil {
			return err
		}
		currentPage++
	}

	allData, err := joinPages(pageCache, currentPage)
	if err != nil {
		return err
	}

	clusters, err := processResponse(allData)
	if err != nil {
		return err
	}

	return postClusters(clusters)
}

func fetchData(entityType string, limit, currentPage int, after int64) (string, *http.Response, error) {
	log.Printf("Fetching data for page %d", currentPage)
	path := fmt.Sprintf("/api/entities/%s?columns=[\"title\"]&limit=%d&page=%d&after=%d", entityType, limit, currentPage, after)
	body, resp, err := data.GetBodyAndResponse(path)
	if err != nil {
		log.Printf("Error fetching cluster data: %v", err)
		return "", nil, err
	}
	return body, resp, nil
}

func parseHasMoreHeader(resp *http.Response) (bool, error) {
	hasMoreHeader := resp.Header.Get("x-has-more")
	hasMore, err := strconv.ParseBool(hasMoreHeader)
	if err != nil {
		log.Printf("Error parsing x-has-more header: %v", err)
		return false, err
	}
	return hasMore, nil
}

func joinPages(pageCache map[int]interface{}, currentPage int) ([]interface{}, error) {
	var allData []interface{}
	for i := 1; i < currentPage; i++ {
		pageData, ok := pageCache[i].([]interface{})
		if !ok {
			log.Printf("Error asserting page data to []interface{}")
			return nil, fmt.Errorf("error asserting page data to []interface{}")
		}
		allData = append(allData, pageData...)
	}
	return allData, nil
}

func processResponse(allData []interface{}) (map[string]map[string]interface{}, error) {
	allDataBytes, err := json.Marshal(allData)
	if err != nil {
		log.Printf("Error marshalling all data: %v", err)
		return nil, err
	}

	cmd := exec.Command("scripts/.venv/bin/python3", "scripts/cluster_headlines.py")
	cmd.Stdin = bytes.NewReader(allDataBytes)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Create a pipe for status updates
	statusReader, statusWriter, err := os.Pipe()
	if err != nil {
		log.Printf("Error creating pipe: %v", err)
		return nil, err
	}
	defer statusReader.Close()
	defer statusWriter.Close()

	// Pass the write-end of the pipe to the child process
	cmd.ExtraFiles = []*os.File{statusWriter}

	// Goroutine to read status updates
	go func() {
		statusBuf := make([]byte, 1024)
		for {
			n, err := statusReader.Read(statusBuf)
			if err != nil {
				break
			}
			log.Printf("Status update: %s", string(statusBuf[:n]))
		}
	}()

	err = cmd.Run()
	if err != nil {
		log.Printf("Error running Python script: %v", err)
		log.Printf("Stderr: %s", stderr.String())
		return nil, err
	}
	log.Print("Python script completed successfully")
	log.Print("Starting to parse JSON output")
	var clusters map[string]map[string]interface{}
	err = json.Unmarshal(out.Bytes(), &clusters)
	if err != nil {
		log.Printf("Error unmarshalling JSON output: %v", err)
		return nil, err
	}
	log.Print("Finished parsing JSON output")
	delete(clusters, "-1")
	if len(clusters) == 0 {
		log.Printf("No clusters found in the output")
		return nil, nil
	}
	return clusters, nil
}

func postClusters(clusters map[string]map[string]interface{}) error {
	log.Print("Posting clusters to the server")
	for _, cluster := range clusters {
		ids := make([]string, 0)
		for _, item := range cluster["articles"].([]interface{}) {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				log.Printf("Error asserting item to map[string]interface{}")
				return fmt.Errorf("error asserting item to map[string]interface{}")
			}
			article, ok := itemMap["article"]
			if !ok {
				log.Printf("Error finding article in item")
				return fmt.Errorf("error finding article in item")
			}
			articleMap, ok := article.(map[string]interface{})
			if !ok {
				log.Printf("Error asserting article to map[string]interface{}")
				return fmt.Errorf("error asserting article to map[string]interface{}")
			}
			if id, ok := articleMap["_id"]; ok {
				ids = append(ids, id.(string))
			}
		}
		body := map[string]interface{}{
			"artifacts": ids,
			"centroid":  cluster["centroid"],
		}
		err := data.Post("/api/clusters", body)
		if err != nil {
			log.Printf("Error posting cluster data: %v", err)
			return err
		}
	}
	return nil
}

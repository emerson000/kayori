package utilities

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	"kayori.io/drone/internal/data"
)

type ClusterTask struct {
	EntityType string `json:"entity_type"`
}

func ProcessClusterTask(jobId string, task *ClusterTask, postJSON func(url string, data interface{}) error) error {
	log.Printf("Starting cluster job: %v", jobId)
	entity_type := task.EntityType
	log.Printf("Processing cluster task for entity type: %v", entity_type)
	path := "/api/entities/" + entity_type + "?columns=[\"title\"]&limit=100"
	resp, err := data.Get(path)
	if err != nil {
		log.Printf("Error fetching cluster data: %v", err)
	}
	if resp == "" {
		log.Printf("No data found for entity type: %v", entity_type)
	}
	cmd := exec.Command("scripts/.venv/bin/python3", "scripts/cluster_headlines.py")
	cmd.Stdin = bytes.NewReader([]byte(resp))
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
	return nil
}

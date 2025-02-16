package models

import (
	"github.com/gocql/gocql"
	"kayori.io/backend/data/mongorm"
)

type NewsArticle struct {
	mongorm.Model
	ArtifactID  gocql.UUID `json:"artifact_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	Published   string     `json:"published"`
	Timestamp   int64      `json:"timestamp"`
	Author      string     `json:"author"`
	Categories  []string   `json:"categories"`
	ServiceID   string     `json:"service_id"`
	Checksum    string     `json:"checksum"`
	JobId       string     `json:"job_id"`
}

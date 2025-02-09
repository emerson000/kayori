package models

type NewsArticle struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Published   string   `json:"published"`
	Timestamp   int64    `json:"timestamp"`
	Author      string   `json:"author"`
	Categories  []string `json:"categories"`
	SourceID    string   `json:"source_id"`
	Checksum    string   `json:"checksum"`
	JobId       string   `json:"job_id"`
}

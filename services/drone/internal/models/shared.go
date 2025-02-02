package models

type NewsArticle struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Url           string   `json:"url"`
	Published     string   `json:"published"`
	PublishedUnix string   `json:"published_unix"`
	Author        string   `json:"author"`
	Categories    []string `json:"categories"`
	Source        string   `json:"source"`
	Checksum      string   `json:"checksum"`
}

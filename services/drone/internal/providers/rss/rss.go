package rss

import (
	"crypto/sha256"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/v2/bson"
	"kayori.io/drone/internal/data"
	"kayori.io/drone/internal/models"
)

func calculateChecksum(article models.NewsArticle) string {
	data := article.Title + article.URL + article.Description + article.Author + article.Published + article.ServiceID
	if article.Categories != nil {
		for _, category := range article.Categories {
			data += category
		}
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}

type Task struct {
	Url  string   `json:"url"`
	Urls []string `json:"urls"`
}

func ProcessTask(jobId string, task *Task, postJSON func(url string, data interface{}) error) {
	processURL := func(url string) {
		log.Printf("RSS URL: %+v", url)
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)

		if err != nil {
			log.Fatalf("Error parsing feed: %v", err)
			return
		}

		var wg sync.WaitGroup
		for _, item := range feed.Items {
			wg.Add(1)
			go func(item *gofeed.Item) {
				defer wg.Done()
				var newsArticle models.NewsArticle
				newsArticle.Title = item.Title
				if item.Link != "" {
					newsArticle.URL = item.Link
				}
				newsArticle.Description = item.Description
				if item.Author != nil {
					newsArticle.Author = item.Author.Name
				}
				if item.Published != "" {
					if item.PublishedParsed != nil {
						newsArticle.Timestamp = *item.PublishedParsed
						newsArticle.Published = item.PublishedParsed.Format(time.RFC3339)
					}
				}
				if item.Categories != nil {
					newsArticle.Categories = item.Categories
				}
				newsArticle.ServiceID = url
				newsArticle.Service = "rss"
				newsArticle.JobId, err = bson.ObjectIDFromHex(jobId)
				if err != nil {
					log.Fatalf("Error converting job ID to ObjectID: %v", err)
					return
				}
				newsArticle.Checksum = calculateChecksum(newsArticle)
				err = postJSON(data.GetHostname()+"/api/entities/news_articles", newsArticle)
				if err != nil {
					if err.Error() == "received non-OK response: 409 Conflict" {
						log.Printf("Article already exists: %v", newsArticle.Title)
					} else {
						log.Fatalf("Error publishing to backend API: %v", err)
					}
					return
				}
				log.Printf("Found news article: %v", item.Title)
			}(item)
		}
		wg.Wait()
		objId, err := bson.ObjectIDFromHex(jobId)
		if err != nil {
			log.Fatalf("Error converting job ID to ObjectID: %v", err)
			return
		}
		err = data.SetJobStatus(objId, "done")
		if err != nil {
			log.Fatalf("Error setting job status: %v", err)
			return
		}
		log.Printf("Finished parsing RSS feed: %v", feed.Title)
	}

	if len(task.Urls) > 0 {
		var wg sync.WaitGroup
		for _, url := range task.Urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				processURL(url)
			}(url)
		}
		wg.Wait()
	} else {
		processURL(task.Url)
	}
}

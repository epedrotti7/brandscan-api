package service

import (
	"brandscan-api/internal/client"
	"brandscan-api/internal/db"
	"context"
	"log"

	"github.com/beeker1121/goque"
)

type SearchTask struct {
	Query string `json:"query"`
	Email string `json:"email"`
}

type SearchService struct {
	Client       *client.CustomSearchClient
	QueueService *QueueService
	EmailAPIKey  string
	EmailSender  string
}

func NewSearchService(client *client.CustomSearchClient, queue *goque.Queue, emailAPIKey string) *SearchService {
	queueService := NewQueueService(queue, nil)
	searchService := &SearchService{Client: client, QueueService: queueService, EmailAPIKey: emailAPIKey}
	queueService.SearchService = searchService
	return searchService
}

func (s *SearchService) ProcessTask(task SearchTask) error {
	results, err := s.Client.Search(task.Query)
	if err != nil {
		log.Printf("Failed to search: %v", err)
		return err
	}

	collection := db.Client.Database("brandscan").Collection("requests")
	_, err = collection.InsertOne(context.Background(), map[string]interface{}{
		"query":   task.Query,
		"email":   task.Email,
		"domains": results,
	})
	if err != nil {
		log.Printf("Failed to save to MongoDB: %v", err)
		return err
	}

	emailBody := "Here are the search results:<br><br>" + formatDomains(results)
	if err := SendEmail(s.EmailAPIKey, s.EmailSender, task.Email, "Search Results", emailBody); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func formatDomains(domains []string) string {
	result := ""
	for _, domain := range domains {
		result += domain + "<br><br>"
	}
	return result
}

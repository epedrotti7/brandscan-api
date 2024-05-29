package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/beeker1121/goque"
)

type QueueService struct {
	Queue         *goque.Queue
	SearchService *SearchService
}

func NewQueueService(queue *goque.Queue, searchService *SearchService) *QueueService {
	return &QueueService{Queue: queue, SearchService: searchService}
}

func (qs *QueueService) EnqueueTask(task SearchTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	_, err = qs.Queue.Enqueue(data)
	return err
}

func (qs *QueueService) ProcessQueue() {
	for {
		item, err := qs.Queue.Dequeue()
		if err != nil {
			if err == goque.ErrEmpty {
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("Error dequeuing item: %v", err)
			continue
		}

		var task SearchTask
		if err := json.Unmarshal(item.Value, &task); err != nil {
			log.Printf("Error unmarshalling task: %v", err)
			continue
		}

		if err := qs.SearchService.ProcessTask(task); err != nil {
			log.Printf("Error processing task: %v", err)
		}
	}
}

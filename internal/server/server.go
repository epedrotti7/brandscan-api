package server

import (
	"brandscan-api/internal/client"
	"brandscan-api/internal/db"
	"brandscan-api/internal/routes"
	"brandscan-api/internal/service"
	"log"
	"os"

	"github.com/beeker1121/goque"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.ConnectMongo()

	apiKey := os.Getenv("API_KEY")
	cx := os.Getenv("CX")
	mailerSendAPIKey := os.Getenv("MAILER_SEND_API_KEY")
	emailSender := os.Getenv("EMAIL_SENDER")

	customSearchClient := client.NewCustomSearchClient(apiKey, cx)

	queue, err := goque.OpenQueue("task_queue")
	if err != nil {
		log.Fatalf("Failed to open queue: %v", err)
	}
	defer queue.Close()

	searchService := service.NewSearchService(customSearchClient, queue, mailerSendAPIKey)
	searchService.EmailSender = emailSender

	go searchService.QueueService.ProcessQueue()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	routes.SetupRoutes(e, searchService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

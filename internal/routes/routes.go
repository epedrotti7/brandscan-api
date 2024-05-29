package routes

import (
	"brandscan-api/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, searchService *service.SearchService) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/search-ads", func(c echo.Context) error {
		type request struct {
			Query string `json:"query"`
			Email string `json:"email"`
		}

		var req request
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		task := service.SearchTask{
			Query: req.Query,
			Email: req.Email,
		}

		if err := searchService.QueueService.EnqueueTask(task); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusAccepted, map[string]string{"message": "Your request is being processed"})
	})
}

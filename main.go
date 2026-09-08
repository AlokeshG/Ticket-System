package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ticket-system/handlers"
	"ticket-system/middleware"
	"ticket-system/storage"
)

func main() {
	db := storage.InitDB()
	defer db.Close()

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Public authentication routes
	r.POST("/auth/register", handlers.Register(db))
	r.POST("/auth/login", handlers.Login(db))

	// Protected ticket routes
	tickets := r.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.POST("", handlers.CreateTicket(db))
		tickets.GET("", handlers.GetTickets(db))
		tickets.GET("/:id", handlers.GetTicket(db))
		tickets.PATCH("/:id/status", handlers.UpdateTicketStatus(db))
	}

	// Start server
	r.Run(":8080")
}

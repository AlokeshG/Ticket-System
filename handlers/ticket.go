package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// CreateTicket creates a new ticket for the logged-in user.
func CreateTicket(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// We will get user_id from JWT middleware.
		userID, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userIDInt, ok := userID.(int)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user id",
			})
			return
		}

		var request CreateTicketRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		if request.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "title is required",
			})
			return
		}

		// New tickets always start with "open".
		result, err := db.Exec(`
			INSERT INTO tickets (
				user_id,
				title,
				description,
				status
			)
			VALUES (?, ?, ?, ?)
		`,
			userIDInt,
			request.Title,
			request.Description,
			"open",
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create ticket",
			})
			return
		}

		ticketID, err := result.LastInsertId()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get ticket id",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":          ticketID,
			"user_id":     userIDInt,
			"title":       request.Title,
			"description": request.Description,
			"status":      "open",
			"message":     "ticket created successfully",
		})
	}
}

// GetTickets returns only tickets belonging to the logged-in user.
func GetTickets(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userIDInt, ok := userID.(int)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user id",
			})
			return
		}

		rows, err := db.Query(`
			SELECT
				id,
				user_id,
				title,
				description,
				status,
				created_at,
				updated_at
			FROM tickets
			WHERE user_id = ?
			ORDER BY id DESC
		`, userIDInt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch tickets",
			})
			return
		}

		defer rows.Close()

		tickets := make([]gin.H, 0)

		for rows.Next() {

			var (
				id          int
				ticketUser  int
				title       string
				description string
				status      string
				createdAt   time.Time
				updatedAt   time.Time
			)

			err := rows.Scan(
				&id,
				&ticketUser,
				&title,
				&description,
				&status,
				&createdAt,
				&updatedAt,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to read ticket",
				})
				return
			}

			tickets = append(tickets, gin.H{
				"id":          id,
				"user_id":     ticketUser,
				"title":       title,
				"description": description,
				"status":      status,
				"created_at":  createdAt,
				"updated_at":  updatedAt,
			})
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to read tickets",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tickets": tickets,
		})
	}
}

// GetTicket returns one ticket only if it belongs to the logged-in user.
func GetTicket(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userIDInt, ok := userID.(int)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user id",
			})
			return
		}

		ticketID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid ticket id",
			})
			return
		}

		var (
			id          int
			ticketUser  int
			title       string
			description string
			status      string
			createdAt   time.Time
			updatedAt   time.Time
		)

		err = db.QueryRow(`
			SELECT
				id,
				user_id,
				title,
				description,
				status,
				created_at,
				updated_at
			FROM tickets
			WHERE id = ? AND user_id = ?
		`, ticketID, userIDInt).Scan(
			&id,
			&ticketUser,
			&title,
			&description,
			&status,
			&createdAt,
			&updatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ticket not found",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch ticket",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":          id,
			"user_id":     ticketUser,
			"title":       title,
			"description": description,
			"status":      status,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}
}

// UpdateTicketStatus updates the status of a user's own ticket.
func UpdateTicketStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userIDInt, ok := userID.(int)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user id",
			})
			return
		}

		ticketID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid ticket id",
			})
			return
		}

		var request UpdateStatusRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		// Only these three statuses are allowed.
		if request.Status != "open" &&
			request.Status != "in_progress" &&
			request.Status != "closed" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid status",
			})
			return
		}

		// Find the ticket and make sure it belongs to this user.
		var currentStatus string

		err = db.QueryRow(`
			SELECT status
			FROM tickets
			WHERE id = ? AND user_id = ?
		`, ticketID, userIDInt).Scan(&currentStatus)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ticket not found",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch ticket",
			})
			return
		}

		// Validate status transition.
		if currentStatus == "open" &&
			request.Status != "in_progress" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "open ticket can only move to in_progress",
			})
			return
		}

		if currentStatus == "in_progress" &&
			request.Status != "closed" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "in_progress ticket can only move to closed",
			})
			return
		}

		if currentStatus == "closed" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "closed ticket cannot be reopened",
			})
			return
		}

		// Update status.
		_, err = db.Exec(`
			UPDATE tickets
			SET status = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ? AND user_id = ?
		`,
			request.Status,
			ticketID,
			userIDInt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update ticket status",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      ticketID,
			"status":  request.Status,
			"message": "ticket status updated successfully",
		})
	}
}

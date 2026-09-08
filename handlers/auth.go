package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"ticket-system/utils"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var request RegisterRequest

		// Read JSON request body
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		// Remove unnecessary spaces from email
		request.Email = strings.TrimSpace(request.Email)

		// Validate email
		if request.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email is required",
			})
			return
		}

		// Validate password
		if request.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "password is required",
			})
			return
		}

		// Check if user already exists
		var existingUserID int

		err := db.QueryRow(
			"SELECT id FROM users WHERE email = ?",
			request.Email,
		).Scan(&existingUserID)

		if err == nil {
			// User already exists
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already exists",
			})
			return
		}

		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(request.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to hash password",
			})
			return
		}

		// Insert user into database
		result, err := db.Exec(
			"INSERT INTO users (email, password_hash) VALUES (?, ?)",
			request.Email,
			string(hashedPassword),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create user",
			})
			return
		}

		userID, err := result.LastInsertId()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get user id",
			})
			return
		}

		// Return successful response
		c.JSON(http.StatusCreated, gin.H{
			"id":      userID,
			"email":   request.Email,
			"message": "user registered successfully",
		})
	}
}

// Login authenticates a user and returns a JWT.
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var request LoginRequest

		// Read JSON request body.
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		// Remove unnecessary spaces from email.
		request.Email = strings.TrimSpace(request.Email)

		// Validate email.
		if request.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email is required",
			})
			return
		}

		// Validate password.
		if request.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "password is required",
			})
			return
		}

		// Find the user in the database.
		var (
			userID       int
			passwordHash string
		)

		err := db.QueryRow(
			"SELECT id, password_hash FROM users WHERE email = ?",
			request.Email,
		).Scan(&userID, &passwordHash)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		// Compare the password entered by the user
		// with the hashed password stored in the database.
		err = bcrypt.CompareHashAndPassword(
			[]byte(passwordHash),
			[]byte(request.Password),
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		}

		// Generate JWT.
		token, err := utils.GenerateToken(userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate token",
			})
			return
		}

		// Return JWT.
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

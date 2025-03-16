package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/jimsyyap/error_app/backend/pkg/models"
)

// SessionHandler handles session-related operations
type SessionHandler struct {
	DB *gorm.DB
}

// SessionRequest represents a session creation request
type SessionRequest struct {
	OpponentName *string `json:"opponent_name"`
	Location     *string `json:"location"`
	Notes        *string `json:"notes"`
}

// SessionSummary represents a session's error summary
type SessionSummary struct {
	TotalErrors  int            `json:"total_errors"`
	ErrorsByType map[string]int `json:"errors_by_type"`
}

// StartSession starts a new match session
func (h *SessionHandler) StartSession(c *gin.Context) {
	var req SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Continue even if binding fails, as all fields are optional
	}

	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user already has an active session
	var activeSession models.MatchSession
	result := h.DB.Where("user_id = ? AND end_time IS NULL", userID).First(&activeSession)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already has an active session",
			"session_id": activeSession.SessionID,
		})
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Create new session
	session := models.MatchSession{
		UserID:       userID,
		StartTime:    time.Now(),
		OpponentName: req.OpponentName,
		Location:     req.Location,
		Notes:        req.Notes,
	}

	if err := h.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"session_id": session.SessionID})
}

// EndSession ends an active match session
func (h *SessionHandler) EndSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Find the session
	var session models.MatchSession
	if err := h.DB.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if already ended
	if !session.IsActive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session already ended"})
		return
	}

	// End the session
	if err := session.End(h.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session ended successfully"})
}

// GetSessions gets all user's sessions
func (h *SessionHandler) GetSessions(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var sessions []models.MatchSession
	if err := h.DB.Where("user_id = ?", userID).Order("start_time DESC").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// GetActiveSession gets the user's active session if any
func (h *SessionHandler) GetActiveSession(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var session models.MatchSession
	result := h.DB.Where("user_id = ? AND end_time IS NULL", userID).First(&session)
	if result.Error != nil {
		if errors.Is(result.

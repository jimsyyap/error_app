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

// ErrorHandler handles error logging operations
type ErrorHandler struct {
	DB *gorm.DB
}

// ErrorLogRequest represents an error log creation request
type ErrorLogRequest struct {
	SessionID   uuid.UUID `json:"session_id" binding:"required"`
	ErrorTypeID int       `json:"error_type_id" binding:"required"`
}

// Common errors
var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrSessionNotFound = errors.New("session not found")
	ErrNotActiveSession = errors.New("session is not active")
	ErrErrorTypeNotFound = errors.New("error type not found")
)

// LogError logs a new error
func (h *ErrorHandler) LogError(c *gin.Context) {
	var req ErrorLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify session exists and belongs to the user
	var session models.MatchSession
	if err := h.DB.Where("session_id = ? AND user_id = ?", req.SessionID, userID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if session is active
	if !session.IsActive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot log errors to a completed session"})
		return
	}

	// Verify error type exists
	var errorType models.ErrorType
	if err := h.DB.First(&errorType, req.ErrorTypeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Error type not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Create error log
	errorLog := models.ErrorLog{
		SessionID:   req.SessionID,
		ErrorTypeID: req.ErrorTypeID,
		Timestamp:   time.Now(),
	}

	if err := h.DB.Create(&errorLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error_id": errorLog.ErrorID,
		"message": "Error logged successfully",
	})
}

// UndoLastError removes the last error in a session
func (h *ErrorHandler) UndoLastError(c *gin.Context) {
	type UndoRequest struct {
		SessionID uuid.UUID `json:"session_id" binding:"required"`
	}

	var req UndoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify session exists and belongs to the user
	var session models.MatchSession
	if err := h.DB.Where("session_id = ? AND user_id = ?", req.SessionID, userID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Get the last error
	var lastError models.ErrorLog
	result := h.DB.Where("session_id = ?", req.SessionID).Order("timestamp DESC").First(&lastError)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No errors to undo"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Delete the last error
	if err := h.DB.Delete(&lastError).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to undo error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Last error undone successfully"})
}

// GetErrorTypes gets all error types
func (h *ErrorHandler) GetErrorTypes(c *gin.Context) {
	var errorTypes []models.ErrorType
	if err := h.DB.Find(&errorTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve error types"})
		return
	}

	c.JSON(http.StatusOK, errorTypes)
}

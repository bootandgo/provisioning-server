package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"provisioning-server/models"
	storepkg "provisioning-server/store"
)

func RegisterServerHandler(store storepkg.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			SerialNumber string `json:"serial_number" binding:"required"`
			IPAddress    string `json:"ip_address" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check for existing server
		existing, err := store.FindServerBySerialNumber(req.SerialNumber)
		if err != nil && !errors.Is(err, storepkg.ErrServerNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check server existence"})
			return
		}

		if existing != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "serial number already exists"})
			return
		}

		server := &models.Server{
			ID:           uuid.New().String(),
			SerialNumber: req.SerialNumber,
			IPAddress:    req.IPAddress,
			Status:       "pending",
			CreatedAt:    time.Now(),
		}

		if err := store.CreateServer(server); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register server"})
			return
		}

		c.JSON(http.StatusCreated, server)
	}
}

func ListServersHandler(store storepkg.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		servers, err := store.ListServers()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list servers"})
			return
		}

		c.JSON(http.StatusOK, servers)
	}
}

func ApproveServerHandler(store storepkg.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ServerID string `json:"server_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := store.ApproveServer(req.ServerID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "server approved"})
	}
}

func ServerStatusHandler(store storepkg.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		serverID := c.Param("id")

		server, err := store.GetServerByID(serverID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
			return
		}

		c.JSON(http.StatusOK, server)
	}
}

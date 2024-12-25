package handlers

import (
	clients "feed-service/client"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MachineFeedHandler struct {
	Client clients.MachineServiceClient
}

func NewMachineFeedHandler(client clients.MachineServiceClient) *MachineFeedHandler {
	return &MachineFeedHandler{Client: client}
}

func (h *MachineFeedHandler) GetFeeds(c *gin.Context) {
	machineID := c.Param("machineId")
	if machineID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "machineId path parameter is required"})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	// Fetch machine name
	machineName, err := h.Client.FetchMachineName(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch machine name", "details": err.Error()})
		return
	}

	// Fetch repairs and sessions
	repairs, err := h.Client.FetchRepairs(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repairs", "details": err.Error()})
		return
	}

	sessions, err := h.Client.FetchSessions(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions", "details": err.Error()})
		return
	}

	// Combine and sort data
	feed := append(repairs, sessions...)
	sort.Slice(feed, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, feed[i].Timestamp)
		t2, _ := time.Parse(time.RFC3339, feed[j].Timestamp)
		return t1.After(t2)
	})

	// Implement pagination
	totalItems := len(feed)
	startIndex := (page - 1) * size
	endIndex := startIndex + size

	if startIndex >= totalItems {
		c.JSON(http.StatusOK, gin.H{
			"machineId":   machineID,
			"machineName": machineName,
			"feeds":       []interface{}{},
			"page":        page,
			"size":        size,
			"totalPages":  (totalItems + size - 1) / size,
			"totalItems":  totalItems,
		})
		return
	}

	if endIndex > totalItems {
		endIndex = totalItems
	}

	paginatedFeed := feed[startIndex:endIndex]

	// Respond with the feed
	c.JSON(http.StatusOK, gin.H{
		"machineId":   machineID,
		"machineName": machineName,
		"feeds":       paginatedFeed,
		"page":        page,
		"size":        size,
		"totalPages":  (totalItems + size - 1) / size,
		"totalItems":  totalItems,
	})
}

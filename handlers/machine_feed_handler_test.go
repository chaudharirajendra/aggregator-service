package handlers

import (
	clients "feed-service/client"

	"feed-service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestMachineFeedHandler(t *testing.T) {
	mockClient := &clients.MockHTTPMachineServiceClient{
		MockFetchRepairs: func(machineID string) ([]models.FeedItem, error) {
			return []models.FeedItem{
				{Type: "repair", Details: "Repair 1", Timestamp: "2023-12-25T10:00:00Z"},
			}, nil
		},
		MockFetchSessions: func(machineID string) ([]models.FeedItem, error) {
			return []models.FeedItem{
				{Type: "session", Details: "Session 1", Timestamp: "2023-12-25T11:00:00Z"},
			}, nil
		},
		MockFetchMachineName: func(machineID string) (string, error) {
			return "Machine ABC", nil
		},
	}

	handler := NewMachineFeedHandler(mockClient)

	router := gin.Default()
	router.GET("/machine-feeds/:machineId", handler.GetFeeds)

	req, _ := http.NewRequest("GET", "/machine-feeds/123?page=1&size=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{
		"machineId":"123",
		"machineName":"Machine ABC",
		"feeds":[{"type":"session","details":"Session 1","timestamp":"2023-12-25T11:00:00Z"}],
		"page":1,
		"size":1,
		"totalPages":2,
		"totalItems":2
	}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

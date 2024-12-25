package clients

import (
	"encoding/json"
	"errors"
	"feed-service/models"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type MachineServiceClient interface {
	FetchRepairs(machineID string) ([]models.FeedItem, error)
	FetchSessions(machineID string) ([]models.FeedItem, error)
	FetchMachineName(machineID string) (string, error)
}

type HTTPMachineServiceClient struct {
	RepairsServiceURL     string
	SessionsServiceURL    string
	MachineNameServiceURL string
	HTTPClient            *http.Client
}

type MockHTTPMachineServiceClient struct {
	MockFetchRepairs     func(machineID string) ([]models.FeedItem, error)
	MockFetchSessions    func(machineID string) ([]models.FeedItem, error)
	MockFetchMachineName func(machineID string) (string, error)
}

func NewHTTPMachineServiceClient(repairsURL, sessionsURL, machineNameURL string) *HTTPMachineServiceClient {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.HTTPClient.Timeout = 10 * time.Second

	return &HTTPMachineServiceClient{
		RepairsServiceURL:     repairsURL,
		SessionsServiceURL:    sessionsURL,
		MachineNameServiceURL: machineNameURL,
		HTTPClient:            client.StandardClient(),
	}
}

func (c *HTTPMachineServiceClient) FetchMachineName(machineID string) (string, error) {
	url := c.MachineNameServiceURL + "?machine_id=" + machineID
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch machine name: " + resp.Status)
	}

	var response struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	return response.Name, nil
}

// FetchRepairs fetches repair data
func (c *HTTPMachineServiceClient) FetchRepairs(machineID string) ([]models.FeedItem, error) {
	url := c.RepairsServiceURL + "?machine_id=" + machineID
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch repairs: " + resp.Status)
	}

	var repairs []models.FeedItem
	if err := json.NewDecoder(resp.Body).Decode(&repairs); err != nil {
		return nil, err
	}
	return repairs, nil
}

// FetchSessions fetches session data
func (c *HTTPMachineServiceClient) FetchSessions(machineID string) ([]models.FeedItem, error) {
	url := c.SessionsServiceURL + "?machine_id=" + machineID
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch sessions: " + resp.Status)
	}

	var sessions []models.FeedItem
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (m *MockHTTPMachineServiceClient) FetchRepairs(machineID string) ([]models.FeedItem, error) {
	if m.MockFetchRepairs != nil {
		return m.MockFetchRepairs(machineID)
	}
	return nil, nil
}

func (m *MockHTTPMachineServiceClient) FetchSessions(machineID string) ([]models.FeedItem, error) {
	if m.MockFetchSessions != nil {
		return m.MockFetchSessions(machineID)
	}
	return nil, nil
}

func (m *MockHTTPMachineServiceClient) FetchMachineName(machineID string) (string, error) {
	if m.MockFetchMachineName != nil {
		return m.MockFetchMachineName(machineID)
	}
	return "", nil
}

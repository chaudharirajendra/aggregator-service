package models

type FeedItem struct {
	Type      string `json:"type"`
	Details   string `json:"details"`
	Timestamp string `json:"timestamp"`
}

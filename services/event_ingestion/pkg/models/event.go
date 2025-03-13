package models

import "time"

type Event struct {
	EventType string    `json:"event_type"`
	Id        string    `json:"id"`
	PageUrl   string    `json:"page_url"`
	Timestamp time.Time `json:"timestamp"`
	UserId    string    `json:"user_id"`
}

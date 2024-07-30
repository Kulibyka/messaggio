package models

import "time"

type Message struct {
	ID          int64     `json:"id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	Processed   bool      `json:"processed"`
	ProcessedAt time.Time `json:"processed_at"`
}

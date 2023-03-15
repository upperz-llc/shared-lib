package sharedlib

import "time"

type AlarmType int

const (
	Connection = iota
	MonitoringState
)

type Alarm struct {
	ID              string    `json:"id"`
	Type            int       `json:"type"`
	ClientID        string    `json:"client_id"`
	CreatedAt       time.Time `json:"created_at"`
	ClosedAt        time.Time `json:"closed_at"`
	AckedAt         time.Time `json:"acked_at"`
	AckedBy         string    `json:"acked_by"`
	AckedCheckCount int       `json:"acked_check_count"`
	Acked           bool      `json:"acked"`
	Active          bool      `json:"active"`
}

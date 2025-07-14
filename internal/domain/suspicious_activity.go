package domain

import "time"

type Severity string

const (
	LowActivity      Severity = "low"
	MediumActivity   Severity = "medium"
	HighActivity     Severity = "high"
	CriticalActivity Severity = "critical"
)

type SuspiciousActivity struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	ActivityType string     `json:"activity_type"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	Description  string     `json:"description"`
	Metadata     []byte     `json:"metadata"`
	Severity     string     `json:"severity"`
	Resolved     *bool      `json:"resolved"`
	CreatedAt    *time.Time `json:"created_at"`
}

type CreateSuspiciousActivityAction struct {
	UserID       int64
	ActivityType string
	IPAddress    string
	UserAgent    string
	Description  string
	Metadata     []byte
	Severity     *Severity
}

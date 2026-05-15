package types

import "time"

type Timer struct {
	ID          int       `json:"id"`
	ProjectName string    `json:"project_name"`
	StartedTime time.Time `json:"started_time"`
}

type TimerListResponse struct {
	Timers []Timer `json:"timers"`
}

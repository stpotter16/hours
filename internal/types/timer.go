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

type AddTimerRequest struct {
	StartedAt time.Time `json:"started_at"`
	StoppedAt time.Time `json:"stopped_at"`
}

type StopTimerRequest struct {
	StoppedAt time.Time `json:"stopped_at"`
}

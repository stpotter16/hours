package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/stpotter16/hours/internal/store"
	"github.com/stpotter16/hours/internal/store/sqlite"
	"github.com/stpotter16/hours/internal/types"
)

func postTimers(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectName := r.PathValue("name")

		var req types.AddTimerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && !req.StartedAt.IsZero() && !req.StoppedAt.IsZero() {
			if req.StoppedAt.Before(req.StartedAt) {
				http.Error(w, "stopped_at must be after started_at", http.StatusBadRequest)
				return
			}
			if err := s.AddTimer(r.Context(), projectName, req.StartedAt, req.StoppedAt); err != nil {
				log.Printf("timerPost add: %v", err)
				http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}

		if err := s.StartTimer(r.Context(), projectName); err != nil {
			if errors.Is(err, sqlite.ErrTimerAlreadyRunning) {
				http.Error(w, "Timer already running for project", http.StatusConflict)
				return
			}
			log.Printf("timerPost: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func deleteTimers(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectName := r.PathValue("name")

		stoppedAt := time.Now()
		var req types.StopTimerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && !req.StoppedAt.IsZero() {
			stoppedAt = req.StoppedAt
		}

		if err := s.StopTimer(r.Context(), projectName, stoppedAt); err != nil {
			if errors.Is(err, sqlite.ErrNoActiveTimer) {
				http.Error(w, "No active timer for project", http.StatusNotFound)
				return
			}
			log.Printf("timerDelete: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getTimers(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timers, err := s.GetTimers(r.Context())
		if err != nil {
			log.Printf("timerGet: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(types.TimerListResponse{Timers: timers}); err != nil {
			log.Printf("timerGet encode: %v", err)
		}
	}
}

package handlers

import (
	"log"
	"net/http"

	"github.com/stpotter16/hours/internal/store"
)

func postProjects(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := s.CreateProject(r.Context()); err != nil {
			log.Printf("projectPost: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}

func getProjects(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := s.GetProjects(r.Context()); err != nil {
			log.Printf("projectGet: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

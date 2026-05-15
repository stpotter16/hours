package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/stpotter16/hours/internal/parse"
	"github.com/stpotter16/hours/internal/store"
	"github.com/stpotter16/hours/internal/types"
)

func postProjects(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createProjectReq, err := parse.ParseProjectCreateRequest(r)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if _, err := s.CreateProject(r.Context(), createProjectReq.Name); err != nil {
			log.Printf("projectPost: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}
}

func getProjects(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := s.GetProjects(r.Context())
		if err != nil {
			log.Printf("projectGet: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(types.ProjectListResponse{Projects: projects}); err != nil {
			log.Printf("projectGet encode: %v", err)
		}
	}
}

func deleteProjects(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectName := r.PathValue("name")
		if err := s.DeleteProject(r.Context(), projectName); err != nil {
			log.Printf("projectDelete: %v", err)
			http.Error(w, "Server issue - try again later", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

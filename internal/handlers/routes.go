package handlers

import (
	"net/http"

	"github.com/stpotter16/hours/internal/handlers/authentication"
	"github.com/stpotter16/hours/internal/handlers/middleware"
	"github.com/stpotter16/hours/internal/handlers/sessions"
	"github.com/stpotter16/hours/internal/store"
)

func addRoutes(
	mux *http.ServeMux,
	store store.Store,
	sessionManager sessions.SessionManger,
	authenticator authentication.Authenticator,
) {
	// Auth
	mux.HandleFunc("POST /login", loginPost(authenticator, sessionManager))

	// Session authenticated API endpoints
	apiAuthRequired := middleware.NewApiAuthenticationRequiredMiddleware(sessionManager)
	mux.Handle("POST /projects", apiAuthRequired(postProjects(store)))
	mux.Handle("GET /projects", apiAuthRequired(getProjects(store)))
	mux.Handle("DELETE /projects/{name}", apiAuthRequired(deleteProjects(store)))
}

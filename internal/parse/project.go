package parse

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/stpotter16/hours/internal/types"
)

func ParseProjectCreateRequest(r *http.Request) (types.ProjectCreateRequest, error) {
	var req types.ProjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return types.ProjectCreateRequest{}, err
	}

	if req.Name == "" {
		return types.ProjectCreateRequest{}, errors.New("Project name required")
	}

	return req, nil
}

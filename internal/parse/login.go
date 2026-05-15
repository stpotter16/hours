package parse

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/stpotter16/hours/internal/types"
)

func ParseLoginPost(r *http.Request) (types.LoginRequest, error) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return types.LoginRequest{}, err
	}
	if req.Passphrase == "" {
		return types.LoginRequest{}, errors.New("passphrase is required")
	}
	return req, nil
}

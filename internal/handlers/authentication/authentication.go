package authentication

import (
	"context"
	"errors"

	"github.com/stpotter16/hours/internal/types"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Authenticator struct {
	passphrase string
}

func New(getenv func(string) string) (Authenticator, error) {
	passphrase := getenv("HOURS_PASSPHRASE")
	if passphrase == "" {
		return Authenticator{}, errors.New("could not locate passphrase environment variable")
	}

	a := Authenticator{
		passphrase: passphrase,
	}

	return a, nil
}

func (a Authenticator) AuthenticateUser(ctx context.Context, req types.LoginRequest) (int, error) {

	if req.Passphrase != a.passphrase {
		return 0, ErrInvalidCredentials
	}

	return 1, nil
}

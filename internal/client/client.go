package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/stpotter16/hours/internal/types"
)

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client
}

func New(getenv func(string) string) (Client, error) {
	raw := getenv("HOURS_HOST")
	if raw == "" {
		return Client{}, errors.New("Missing HOURS_HOST environment variable")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return Client{}, err
	}
	if u.Host == "" {
		return Client{}, fmt.Errorf("Invalid server address: %s", raw)
	}

	return Client{baseUrl: u, httpClient: &http.Client{Timeout: 10 * time.Second}}, nil
}

func (c Client) Login(passphrase string) error {
	body, err := json.Marshal(types.LoginRequest{Passphrase: passphrase})
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Post(c.baseUrl.JoinPath("/login").String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Login failed - %s", resp.Status)
	}

	var session string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "X-HOURS-SESSION" {
			session = cookie.Value
			break
		}
	}
	if session == "" {
		return errors.New("No session cookie in login response")
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("Could not find user config directory: %v", err)
	}
	sessionDir := filepath.Join(dir, "hours")
	if err = os.MkdirAll(sessionDir, 0700); err != nil {
		return fmt.Errorf("Could not create session config directory: %v", err)
	}
	return os.WriteFile(filepath.Join(sessionDir, "session"), []byte(session), 0600)
}

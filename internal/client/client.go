package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/stpotter16/hours/internal/types"
)

const sessionCookieName = "X-HOURS-SESSION"

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client
}

func New(host string) (Client, error) {
	if host == "" {
		return Client{}, errors.New("No host configured — run `hours configure` or set HOURS_HOST")
	}
	u, err := url.Parse(host)
	if err != nil {
		return Client{}, err
	}
	if u.Host == "" {
		return Client{}, fmt.Errorf("Invalid server address: %s", host)
	}

	return Client{baseUrl: u, httpClient: &http.Client{Timeout: 10 * time.Second}}, nil
}

func (c Client) Login(ctx context.Context, passphrase string) error {
	body, err := json.Marshal(types.LoginRequest{Passphrase: passphrase})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl.JoinPath("/login").String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Could not build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Login failed - %s", resp.Status)
	}

	var session string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == sessionCookieName {
			session = cookie.Value
			break
		}
	}
	if session == "" {
		return errors.New("No session cookie in login response")
	}

	return saveSession(session)
}

func (c Client) ListProjects(ctx context.Context) (types.ProjectListResponse, error) {
	req, err := c.newAuthRequest(ctx, "GET", c.baseUrl.JoinPath("/projects").String(), nil)
	if err != nil {
		return types.ProjectListResponse{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return types.ProjectListResponse{}, fmt.Errorf("Could not send GET /projects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.ProjectListResponse{}, fmt.Errorf("list projects failed: %s", resp.Status)
	}

	var result types.ProjectListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.ProjectListResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c Client) CreateProject(ctx context.Context, name string) error {
	body, err := json.Marshal(types.ProjectCreateRequest{Name: name})
	if err != nil {
		return err
	}
	req, err := c.newAuthRequest(ctx, "POST", c.baseUrl.JoinPath("/projects").String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not send POST /projects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create projects failed: %s", resp.Status)
	}

	return nil
}

func (c Client) DeleteProject(ctx context.Context, name string) error {
	req, err := c.newAuthRequest(ctx, "DELETE", c.baseUrl.JoinPath("/projects/"+name).String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not send DELETE /projects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete projects failed: %s", resp.Status)
	}

	return nil
}

func (c Client) newAuthRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	session, err := loadSession()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("Could not build request: %v", err)
	}
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: session})
	return req, nil
}

func sessionPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("Could not find user config directory: %v", err)
	}
	return filepath.Join(dir, "hours", "session"), nil
}

func loadSession() (string, error) {
	path, err := sessionPath()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Could not read session — have you logged in? %v", err)
	}
	return string(b), nil
}

func saveSession(session string) error {
	path, err := sessionPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("Could not create session config directory: %v", err)
	}
	return os.WriteFile(path, []byte(session), 0600)
}

package types

import "time"

type Project struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	CreatedTime      time.Time `json:"created_time"`
	LastModifiedTime time.Time `json:"last_modified_time"`
}

type ProjectListResponse struct {
	Projects []Project `json:"projects"`
}

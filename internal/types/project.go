package types

import "time"

type Project struct {
	ID               int
	Name             string
	CreatedTime      time.Time
	LastModifiedTime time.Time
}

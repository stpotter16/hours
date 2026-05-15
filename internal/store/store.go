package store

import (
	"context"
	"time"

	"github.com/stpotter16/hours/internal/types"
)

type Store interface {
	// Projects
	GetProjects(ctx context.Context) ([]types.Project, error)
	CreateProject(ctx context.Context, name string) (int, error)
	DeleteProject(ctx context.Context, name string) error

	// Timers
	StartTimer(ctx context.Context, projectName string) error
	StopTimer(ctx context.Context, projectName string, stoppedAt time.Time) error
	AddTimer(ctx context.Context, projectName string, startedAt, stoppedAt time.Time) error
	GetTimers(ctx context.Context) ([]types.Timer, error)
}

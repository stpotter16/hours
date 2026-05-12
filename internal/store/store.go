package store

import (
	"context"

	"github.com/stpotter16/hours/internal/types"
)

type Store interface {
	// Projects
	GetProjects(ctx context.Context) ([]types.Project, error)
	CreateProject(ctx context.Context) (int, error)
}

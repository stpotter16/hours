package store

import (
	"context"

	"github.com/stpotter16/hours/internal/types"
)

type Store interface {
	// Clicks
	GetClicks(ctx context.Context) ([]types.Click, error)
	CreateClick(ctx context.Context) (int, error)
}

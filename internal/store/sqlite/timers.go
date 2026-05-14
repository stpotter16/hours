package sqlite

import (
	"context"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
)

var ErrNoActiveTimer = errors.New("no active timer for project")
var ErrTimerAlreadyRunning = errors.New("timer already running for project")

func isUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}

func (s Store) StartTimer(ctx context.Context, projectName string) error {
	now := formatTime(time.Now().UTC())
	_, err := s.db.Exec(ctx,
		`INSERT INTO timer_entries (project_id, started_at)
		 SELECT id, ? FROM projects WHERE name = ?`,
		now, projectName,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrTimerAlreadyRunning
		}
		return err
	}
	return nil
}

func (s Store) StopTimer(ctx context.Context, projectName string) error {
	now := formatTime(time.Now().UTC())
	result, err := s.db.Exec(ctx,
		`UPDATE timer_entries
		 SET stopped_at = ?
		 WHERE stopped_at IS NULL
		   AND project_id = (SELECT id FROM projects WHERE name = ?)`,
		now, projectName,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoActiveTimer
	}
	return nil
}

package sqlite

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

var ErrNoActiveTimer = errors.New("no active timer for project")
var ErrTimerAlreadyRunning = errors.New("timer already running for project")
var ErrProjectHasTimers = errors.New("project has timer entries")

func isUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}

func isForeignKeyConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey
}

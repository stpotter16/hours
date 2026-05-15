package sqlite

import (
	"context"
	"time"

	"github.com/stpotter16/hours/internal/types"
)

func (s Store) GetProjects(ctx context.Context) ([]types.Project, error) {
	rows, err := s.db.Query(ctx,
		`SELECT
			p.id,
			p.name,
			p.created_time,
			p.last_modified_time,
			COALESCE(SUM(
				unixepoch(COALESCE(te.stopped_at, datetime('now'))) - unixepoch(te.started_at)
			), 0) AS total_seconds
		FROM projects p
		LEFT JOIN timer_entries te ON te.project_id = p.id
		GROUP BY p.id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []types.Project

	for rows.Next() {
		var project types.Project
		var createdTime string
		var lastModifiedTime string
		if err := rows.Scan(&project.ID, &project.Name, &createdTime, &lastModifiedTime, &project.TotalSeconds); err != nil {
			return nil, err
		}
		project.CreatedTime, err = parseTime(createdTime)
		if err != nil {
			return nil, err
		}
		project.LastModifiedTime, err = parseTime(lastModifiedTime)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s Store) CreateProject(ctx context.Context, name string) (int, error) {
	now := formatTime(time.Now().UTC())
	result, err := s.db.Exec(ctx,
		`INSERT INTO projects
			(name, created_time, last_modified_time)
		VALUES ( ?, ?, ? )`,
		name,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s Store) DeleteProject(ctx context.Context, name string) error {
	_, err := s.db.Exec(ctx,
		`DELETE FROM projects
		 WHERE name = ?
		`,
		name,
	)
	if err != nil {
		if isForeignKeyConstraintError(err) {
			return ErrProjectHasTimers
		}
		return err
	}
	return nil
}

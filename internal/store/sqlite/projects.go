package sqlite

import (
	"context"
	"time"

	"github.com/stpotter16/hours/internal/types"
)

func (s Store) GetProjects(ctx context.Context) ([]types.Project, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, name, created_time, last_modified_time from projects`,
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
		if err := rows.Scan(&project.ID, &project.Name, &createdTime, &lastModifiedTime); err != nil {
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

	return err
}

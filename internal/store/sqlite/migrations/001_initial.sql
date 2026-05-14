CREATE TABLE IF NOT EXISTS session (
    session_key TEXT PRIMARY KEY,
    value BLOB,
    expires_at TEXT NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_time TEXT NOT NULL,
    last_modified_time TEXT NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS timer_entries (
    id          INTEGER PRIMARY KEY,
    project_id  INTEGER NOT NULL REFERENCES projects(id),
    started_at  TEXT NOT NULL,
    stopped_at  TEXT
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS one_active_timer_per_project
    ON timer_entries(project_id)
    WHERE stopped_at IS NULL;

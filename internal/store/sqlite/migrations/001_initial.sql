CREATE TABLE IF NOT EXISTS session (
    session_key TEXT PRIMARY KEY,
    value BLOB,
    expires_at TEXT NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    created_time TEXT NOT NULL,
    last_modified_time TEXT NOT NULL
) STRICT;

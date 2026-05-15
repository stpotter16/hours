# hours

A CLI tool for tracking time spent on projects.

## Installation

Clone the repo and build the client binary for your machine:

```sh
make client/release
```

This writes the binary to `dist/hours`. Move it somewhere on your `$PATH`:

```sh
mv dist/hours /usr/local/bin/hours
```

## Configuration

Before using the CLI, point it at your server:

```sh
hours configure
```

You'll be prompted for:
- **Host URL** — the address of your hours server (e.g. `https://hours.example.com`)
- **Passphrase** — optional; if set, `hours login` will use it automatically without prompting

Configuration is saved to `~/.config/hours/config.json`.

You can also override the host for a single command with the `HOURS_HOST` environment variable:

```sh
HOURS_HOST=http://localhost:8080 hours projects list
```

## Authentication

```sh
hours login
```

Authenticates with the server and saves a session to `~/.config/hours/session`. Most commands require a valid session.

## Projects

```sh
# List all projects with total time worked
hours projects list

# Create a project
hours projects create <name>

# Delete a project (only works if no time has been recorded)
hours projects delete <name>
```

Example output from `hours projects list`:

```
ID   NAME         TIME WORKED    CREATED      LAST MODIFIED
1    my-project   2h 30m 00s     2026-05-12   2026-05-14
2    other        45s            2026-05-13   2026-05-13
```

## Timers

```sh
# Start a timer on a project
hours timers start <project>

# Stop a running timer
hours timers stop <project>

# Stop a timer as of a past time
hours timers stop <project> "1h30m ago"

# Add a historical time entry (when you forgot to track)
hours timers add <project> <start> <stop>
hours timers add my-project "3h ago" "1h ago"

# List all currently running timers
hours timers list
```

Durations use Go's duration format: `1h`, `30m`, `1h30m`, etc. The `ago` suffix is optional.

package timeutil

import (
	"fmt"
	"strings"
	"time"
)

func ParseRelativeTime(s string) (time.Time, error) {
	s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "ago"))
	d, err := time.ParseDuration(strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid duration %q — use a format like \"2h30m ago\" or \"45m ago\"", s)
	}
	return time.Now().Add(-d), nil
}

func FormatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

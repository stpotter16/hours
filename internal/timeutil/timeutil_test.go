package timeutil_test

import (
	"testing"
	"time"

	"github.com/stpotter16/hours/internal/timeutil"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{0, "0s"},
		{45, "45s"},
		{60, "1m 00s"},
		{90, "1m 30s"},
		{3600, "1h 00m 00s"},
		{3661, "1h 01m 01s"},
		{7384, "2h 03m 04s"},
	}

	for _, tt := range tests {
		got := timeutil.FormatDuration(tt.seconds)
		if got != tt.want {
			t.Errorf("FormatDuration(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestParseRelativeTime(t *testing.T) {
	tests := []struct {
		input     string
		wantDelta time.Duration
		wantErr   bool
	}{
		{"1h ago", time.Hour, false},
		{"30m ago", 30 * time.Minute, false},
		{"2h30m ago", 2*time.Hour + 30*time.Minute, false},
		{"1h", time.Hour, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		before := time.Now()
		got, err := timeutil.ParseRelativeTime(tt.input)
		after := time.Now()

		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseRelativeTime(%q) expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRelativeTime(%q) unexpected error: %v", tt.input, err)
			continue
		}

		earliest := before.Add(-tt.wantDelta)
		latest := after.Add(-tt.wantDelta)
		if got.Before(earliest) || got.After(latest) {
			t.Errorf("ParseRelativeTime(%q) = %v, want ~%v ago", tt.input, got, tt.wantDelta)
		}
	}
}

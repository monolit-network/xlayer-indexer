package models

import (
	"testing"
	"time"
)

func TestNextFreePlanRefillAt(t *testing.T) {
	tests := []struct {
		name string
		now  time.Time
		want time.Time
	}{
		{
			name: "before UTC midnight",
			now:  time.Date(2026, time.July, 22, 23, 59, 59, 0, time.UTC),
			want: time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "exactly at UTC midnight",
			now:  time.Date(2026, time.July, 22, 0, 0, 0, 0, time.UTC),
			want: time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "non UTC input",
			now:  time.Date(2026, time.July, 22, 2, 30, 0, 0, time.FixedZone("UTC+3", 3*60*60)),
			want: time.Date(2026, time.July, 22, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextFreePlanRefillAt(tt.now); !got.Equal(tt.want) {
				t.Fatalf("NextFreePlanRefillAt(%s) = %s, want %s", tt.now, got, tt.want)
			}
		})
	}
}

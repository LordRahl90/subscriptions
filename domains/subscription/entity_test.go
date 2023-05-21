package subscription

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestEndDateWithTrial(t *testing.T) {
	t.Parallel()
	table := []struct {
		name string
		arg  Subscription
		exp  time.Duration
	}{
		{
			name: "with trial period",
			arg: Subscription{
				Duration:    3,
				TrialPeriod: 1,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			},
			exp: time.Duration(4 * month),
		},
		{
			name: "without trial period",
			arg: Subscription{
				Duration: 3,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			},
			exp: time.Duration(3 * month),
		},
		{
			name: "one month",
			arg: Subscription{
				Duration: 1,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			},
			exp: time.Duration(1 * month),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.EndDate()
			diff := got.Sub(tt.arg.CreatedAt)
			require.Equal(t, tt.exp, diff)
		})
	}
}

func TestStatusToString(t *testing.T) {
	t.Parallel()
	table := []struct {
		name, exp string
		arg       Status
	}{
		{
			name: "cancelled",
			arg:  StatusCancelled,
			exp:  "cancelled",
		},
		{
			name: "active",
			arg:  StatusActive,
			exp:  "active",
		},
		{
			name: "paused",
			arg:  StatusPaused,
			exp:  "paused",
		},
		{
			name: "ended",
			arg:  StatusEnded,
			exp:  "ended",
		},
		{
			name: "unknown",
			arg:  -1,
			exp:  "unknown",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.String()
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestStatusFromString(t *testing.T) {
	t.Parallel()
	table := []struct {
		name, arg string
		exp       Status
	}{
		{
			name: "cancelled",
			arg:  "cancelled",
			exp:  StatusCancelled,
		},
		{
			name: "active",
			arg:  "active",
			exp:  StatusActive,
		},
		{
			name: "paused",
			arg:  "paused",
			exp:  StatusPaused,
		},
		{
			name: "ended",
			arg:  "ended",
			exp:  StatusEnded,
		},
		{
			name: "unknown",
			arg:  "unknown",
			exp:  -1,
		},
		{
			name: "cancelled unknowne",
			arg:  "cancelles",
			exp:  -1,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := FromString(tt.arg)
			require.Equal(t, tt.exp, got)
		})
	}
}

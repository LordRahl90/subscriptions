package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSwipeFromString(t *testing.T) {
	table := []struct {
		name, args string
		exp        Swipe
	}{
		{
			name: "yes",
			args: "yes",
			exp:  SwipeYes,
		},
		{
			name: "great",
			args: "great",
			exp:  SwipeInvalid,
		},
		{
			name: "no",
			args: "no",
			exp:  SwipeNo,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := SwipeFromString(tt.args)
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestSwipeToString(t *testing.T) {
	table := []struct {
		name, exp string
		args      Swipe
	}{
		{
			name: "yes",
			exp:  "yes",
			args: SwipeYes,
		},
		{
			name: "invalid",
			exp:  "invalid",
			args: SwipeInvalid,
		},
		{
			name: "no",
			exp:  "no",
			args: SwipeNo,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.String()
			require.Equal(t, tt.exp, got)
		})
	}

}

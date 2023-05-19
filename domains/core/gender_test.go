package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenderFromString(t *testing.T) {
	table := []struct {
		name, args string
		exp        Gender
	}{
		{
			name: "male",
			args: "male",
			exp:  GenderMale,
		},
		{
			name: "female",
			args: "female",
			exp:  GenderFemale,
		},
		{
			name: "others",
			args: "others",
			exp:  GenderOthers,
		},
		{
			name: "unknown",
			args: "unknown",
			exp:  GenderUnknown,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := GenderFromString(tt.args)
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestGenderToString(t *testing.T) {
	table := []struct {
		name, exp string
		args      Gender
	}{
		{
			name: "male",
			exp:  "male",
			args: GenderMale,
		},
		{
			name: "female",
			exp:  "female",
			args: GenderFemale,
		},
		{
			name: "others",
			exp:  "others",
			args: GenderOthers,
		},
		{
			name: "unknown",
			exp:  "unknown",
			args: GenderUnknown,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.String()
			require.Equal(t, tt.exp, got)
		})
	}
}

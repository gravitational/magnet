package magnet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportsFromReader(t *testing.T) {
	testCases := []struct {
		comment string
		input   string
		output  map[string]string
	}{
		{
			comment: "parses input into a set of environment variables and strips the prefix",
			input: `
MAGNET_VERSION=v1.0
MAGNET_PKG_VERSION=v2.0
			`,
			output: map[string]string{
				"VERSION":     "v1.0",
				"PKG_VERSION": "v2.0",
			},
		},
		{
			comment: "ignores input without the expected predix",
			input: `
MAGNET_VERSION=v1.0
PKG_VERSION=v2.0
			`,
			output: map[string]string{
				"VERSION": "v1.0",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.comment, func(t *testing.T) {
			env, err := ImportEnvFromReader(strings.NewReader(tc.input))
			require.NoError(t, err)
			require.Equal(t, env, tc.output)
		})
	}
}

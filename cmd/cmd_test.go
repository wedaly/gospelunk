package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedStdout string
		expectedStderr string
	}{
		{
			name:           "inspect",
			args:           []string{"inspect", "-f", "../pkg/inspect/testdata/testmodule001/localvar.go", "-l", "6", "-c", "2"},
			expectedStdout: "localVar ../pkg/inspect/testdata/testmodule001/localvar.go:6:2\n",
			expectedStderr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := ExecuteInTest(tc.args)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStdout, stdout)
			assert.Equal(t, tc.expectedStderr, stderr)
		})
	}
}

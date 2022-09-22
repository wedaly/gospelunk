package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmd(t *testing.T) {
	testCases := []struct {
		name           string
		dir            string
		args           []string
		expectedStdout string
		expectedStderr string
	}{
		{
			name:           "inspect",
			dir:            "../pkg/inspect/testdata/testmodule001",
			args:           []string{"inspect", "-f", "localvar.go", "-l", "6", "-c", "2"},
			expectedStdout: "localVar localvar.go:6:2\n",
			expectedStderr: "",
		},
		{
			name: "list",
			dir:  "../pkg/list/testdata/testmodule001",
			args: []string{"list", "./..."},
			expectedStdout: `MyVar defs.go:3:5
MyConst defs.go:5:7
MyStruct defs.go:7:6
MyInterface defs.go:17:6
MyFunc defs.go:21:1
`,
			expectedStderr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := ExecuteInTest(tc.args, tc.dir)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStdout, stdout)
			assert.Equal(t, tc.expectedStderr, stderr)
		})
	}
}

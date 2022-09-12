package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	testCases := []struct {
		name         string
		tmplString   string
		buildTmplVal func(*testing.T) any
		expected     string
	}{
		{
			name:       "rel path function",
			tmplString: "{{ . | RelPath }}",
			buildTmplVal: func(t *testing.T) any {
				wd, err := os.Getwd()
				require.NoError(t, err)
				return filepath.Join(wd, "foobar.txt")
			},
			expected: "foobar.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpl, err := Template(tc.tmplString)
			require.NoError(t, err)
			tmplVal := tc.buildTmplVal(t)

			var sb strings.Builder
			err = tmpl.Execute(&sb, tmplVal)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, sb.String())
		})
	}
}

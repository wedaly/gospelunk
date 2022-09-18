package list

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wedaly/gospelunk/pkg/file"
)

func TestList(t *testing.T) {
	pkg := Package{
		Name: "testmodule001",
		ID:   "github.com/wedaly/gospelunk/pkg/list/testdata/testmodule001",
	}

	testPkg := Package{
		Name: "testmodule001",
		ID:   "github.com/wedaly/gospelunk/pkg/list/testdata/testmodule001 [github.com/wedaly/gospelunk/pkg/list/testdata/testmodule001.test]",
	}

	pkgPath := func(p string) string {
		path := filepath.Join("testdata", "testmodule001", p)
		absPath, err := filepath.Abs(path)
		require.NoError(t, err)
		return absPath
	}

	testCases := []struct {
		name     string
		dir      string
		patterns []string
		opts     Options
		expected Result
	}{
		{
			name:     "public only",
			dir:      "testdata/testmodule001",
			patterns: []string{"./..."},
			opts:     Options{},
			expected: Result{
				Defs: []Definition{
					{Name: "MyVar", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 3, Column: 5}},
					{Name: "MyConst", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 5, Column: 7}},
					{Name: "MyStruct", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 7, Column: 6}},
					{Name: "MyStruct.MyField", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 8, Column: 2}},
					{Name: "MyInterface", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 17, Column: 6}},
					{Name: "MyInterface.String", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 18, Column: 2}},
					{Name: "MyFunc", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 21, Column: 1}},
				},
			},
		},
		{
			name:     "include private",
			dir:      "testdata/testmodule001",
			patterns: []string{"./..."},
			opts:     Options{IncludePrivate: true},
			expected: Result{
				Defs: []Definition{
					{Name: "MyVar", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 3, Column: 5}},
					{Name: "MyConst", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 5, Column: 7}},
					{Name: "MyStruct", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 7, Column: 6}},
					{Name: "MyStruct.MyField", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 8, Column: 2}},
					{Name: "MyStruct.privateField", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 9, Column: 2}},
					{Name: "privateStruct", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 12, Column: 6}},
					{Name: "privateStruct.PublicField", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 13, Column: 2}},
					{Name: "privateStruct.privateField", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 14, Column: 2}},
					{Name: "MyInterface", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 17, Column: 6}},
					{Name: "MyInterface.String", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 18, Column: 2}},
					{Name: "MyFunc", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 21, Column: 1}},
					{Name: "privateFunc", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 25, Column: 1}},
				},
			},
		},
		{
			name:     "include tests",
			dir:      "testdata/testmodule001",
			patterns: []string{"./..."},
			opts:     Options{IncludeTests: true},
			expected: Result{
				Defs: []Definition{
					{Name: "MyVar", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 3, Column: 5}},
					{Name: "MyConst", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 5, Column: 7}},
					{Name: "MyStruct", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 7, Column: 6}},
					{Name: "MyStruct.MyField", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 8, Column: 2}},
					{Name: "MyInterface", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 17, Column: 6}},
					{Name: "MyInterface.String", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 18, Column: 2}},
					{Name: "MyFunc", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 21, Column: 1}},
					{Name: "TestStruct", Pkg: testPkg, Loc: file.Loc{Path: pkgPath("defs_test.go"), Line: 3, Column: 6}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oldWd, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(oldWd)

			err = os.Chdir(tc.dir)
			require.NoError(t, err)

			result, err := List(tc.patterns, tc.opts)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

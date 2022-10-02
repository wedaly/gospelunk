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
					{Name: "MyInterface", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 17, Column: 6}},
					{Name: "MyFunc", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 21, Column: 1}},
				},
			},
		},
		{
			name:     "include struct fields",
			dir:      "testdata/testmodule001",
			patterns: []string{"./..."},
			opts:     Options{IncludeStructFields: true},
			expected: Result{
				Defs: []Definition{
					{Name: "MyVar", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 3, Column: 5}},
					{Name: "MyConst", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 5, Column: 7}},
					{Name: "MyStruct", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 7, Column: 6}},
					{Name: "MyStruct.MyField", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 8, Column: 2}},
					{Name: "MyInterface", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 17, Column: 6}},
					{Name: "MyFunc", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 21, Column: 1}},
				},
			},
		},
		{
			name:     "include interface methods",
			dir:      "testdata/testmodule001",
			patterns: []string{"./..."},
			opts:     Options{IncludeInterfaceMethods: true},
			expected: Result{
				Defs: []Definition{
					{Name: "MyVar", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 3, Column: 5}},
					{Name: "MyConst", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 5, Column: 7}},
					{Name: "MyStruct", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 7, Column: 6}},
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
			opts:     Options{IncludePrivate: true, IncludeStructFields: true},
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
					{Name: "MyInterface", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 17, Column: 6}},
					{Name: "MyFunc", Pkg: pkg, Loc: file.Loc{Path: pkgPath("defs.go"), Line: 21, Column: 1}},
					{Name: "TestStruct", Pkg: testPkg, Loc: file.Loc{Path: pkgPath("defs_test.go"), Line: 3, Column: 6}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withWorkingDir(t, tc.dir, func(t *testing.T) {
				result, err := List(tc.patterns, tc.opts)
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			})
		})
	}
}

func TestListWithCGo(t *testing.T) {
	cgoRelPath := filepath.Join("testdata", "testmodule002", "cgo.go")
	cgoPath, err := filepath.Abs(cgoRelPath)
	require.NoError(t, err)

	withWorkingDir(t, "testdata/testmodule002", func(t *testing.T) {
		result, err := List([]string{"."}, Options{})
		require.NoError(t, err)

		expected := Result{
			Defs: []Definition{
				{
					Loc:  file.Loc{Path: cgoPath, Line: 6, Column: 6},
					Name: "MyStruct",
					Pkg: Package{
						Name: "testmodule002",
						ID:   "github.com/wedaly/gospelunk/pkg/list/testdata/testmodule002",
					},
				},
				{
					Loc:  file.Loc{Path: cgoPath, Line: 8, Column: 1},
					Name: "Random",
					Pkg: Package{
						Name: "testmodule002",
						ID:   "github.com/wedaly/gospelunk/pkg/list/testdata/testmodule002",
					},
				},
			},
		}
		assert.Equal(t, expected, result)
	})
}

func TestListWithImports(t *testing.T) {
	withWorkingDir(t, "testdata/testmodule003", func(t *testing.T) {
		result, err := List([]string{"."}, Options{IncludeImports: true})
		require.NoError(t, err)

		// Exact definitions may vary based on Go stdlib, so check that imported packages are included.
		defPkgs := make(map[string][]string, 0)
		for _, def := range result.Defs {
			defPkgs[def.Pkg.ID] = append(defPkgs[def.Pkg.ID], def.Name)
		}
		assert.Equal(t, 3, len(defPkgs))
		assert.Equal(t, []string{"MyFunc"}, defPkgs["github.com/wedaly/gospelunk/pkg/list/testdata/testmodule003"])
		assert.Equal(t, []string{"SubFunc"}, defPkgs["github.com/wedaly/gospelunk/pkg/list/testdata/testmodule003/subpkg"])
		assert.Greater(t, len(defPkgs["fmt"]), 1)
	})
}

func withWorkingDir(t *testing.T, dir string, f func(t *testing.T)) {
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)
	err = os.Chdir(dir)
	require.NoError(t, err)
	f(t)
}

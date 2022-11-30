package inspect

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wedaly/gospelunk/pkg/file"
)

func TestInspectLocalVariableDefinedInSameFunction(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule001/localvar.go",
		Line:   7,
		Column: 32,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "localVar",
		Type: "int",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "localVar",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule001/localvar.go"),
					Line:   6,
					Column: 2,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructTypeLiteral(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule002/struct.go",
		Line:   11,
		Column: 7,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "TestStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule002.TestStruct",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "TestStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule002/struct.go"),
					Line:   5,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructLiteralField(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule002/struct.go",
		Line:   12,
		Column: 3,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "StringField",
		Type: "string",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "StringField",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule002/struct.go"),
					Line:   6,
					Column: 2,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructSelectionField(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule002/struct.go",
		Line:   16,
		Column: 23,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "IntField",
		Type: "int",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "IntField",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule002/struct.go"),
					Line:   7,
					Column: 2,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructMethodCall(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule004/methods.go",
		Line:   28,
		Column: 23,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "String",
		Type: "func() string",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "String",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule004/methods.go"),
					Line:   12,
					Column: 22,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructMethodCallOnReturnedValue(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule004/methods.go",
		Line:   29,
		Column: 32,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "String",
		Type: "func() string",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "String",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule004/methods.go"),
					Line:   22,
					Column: 26,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectInterfaceMethodCall(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule004/methods.go",
		Line:   33,
		Column: 27,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "MyString",
		Type: "func() string",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "MyString",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule004/methods.go"),
					Line:   7,
					Column: 2,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectInterfaceEmbeddedMethodCall(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule004/methods.go",
		Line:   34,
		Column: 27,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "String", result.Name)
	assert.Equal(t, "func() string", result.Type)
	require.Equal(t, 1, len(result.Relations))
	assert.Equal(t, RelationKindDef, result.Relations[0].Kind)
	assert.Contains(t, result.Relations[0].Path, "src/fmt")
}

func TestInspectFuncCallSameFile(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule003/func.go",
		Line:   14,
		Column: 9,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "FuncInSameFile",
		Type: "func() int",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "testmodule003",
				Name: "FuncInSameFile",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule003/func.go"),
					Line:   9,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectFuncCallOtherFile(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule003/func.go",
		Line:   15,
		Column: 9,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "FuncInOtherFile",
		Type: "func() int",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "testmodule003",
				Name: "FuncInOtherFile",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule003/other.go"),
					Line:   3,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectFuncCallOtherPkg(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule003/func.go",
		Line:   16,
		Column: 14,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "FuncInOtherPkg",
		Type: "func(x int) int",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "subpkg",
				Name: "FuncInOtherPkg",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule003/subpkg/subpkgfunc.go"),
					Line:   3,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectFuncCallStdlib(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule003/func.go",
		Line:   17,
		Column: 6,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Printf", result.Name)
	assert.Equal(t, "func(format string, a ...any) (n int, err error)", result.Type)
	require.Equal(t, 1, len(result.Relations))
	assert.Equal(t, RelationKindDef, result.Relations[0].Kind)
	assert.Equal(t, "fmt", result.Relations[0].Pkg)
	assert.Equal(t, "Printf", result.Relations[0].Name)
	assert.Contains(t, result.Relations[0].Path, "src/fmt")
}

func TestInspectPkgNameInSelection(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule003/func.go",
		Line:   17,
		Column: 2,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "fmt",
		Type: "",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "testmodule003",
				Name: "fmt",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule003/func.go"),
					Line:   4,
					Column: 2,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectFileWithComments(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule005/comments.go",
		Line:   12,
		Column: 6,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Println", result.Name)
	assert.Equal(t, "func(a ...any) (n int, err error)", result.Type)
	require.Equal(t, 1, len(result.Relations))
	assert.Equal(t, RelationKindDef, result.Relations[0].Kind)
	assert.Equal(t, "fmt", result.Relations[0].Pkg)
	assert.Equal(t, "Println", result.Relations[0].Name)
	assert.Contains(t, result.Relations[0].Path, "src/fmt")
}

func TestInspectIntegerLiteral(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule006/const.go",
		Line:   5,
		Column: 18,
	}, "")
	assert.EqualError(t, err, "Could not find AST node of type *ast.Ident at location testdata/testmodule006/const.go:5:18")
	assert.Nil(t, result)
}

func TestInspectIntegerConst(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule006/const.go",
		Line:   13,
		Column: 23,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "intConst",
		Type: "untyped int",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "main",
				Name: "intConst",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule006/const.go"),
					Line:   5,
					Column: 7,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectNil(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule006/const.go",
		Line:   9,
		Column: 10,
	}, "")
	require.NoError(t, err)
	expected := &Result{
		Name: "nil",
		Type: "untyped nil",
	}
	assert.Equal(t, expected, result)
}

func TestInspectEmbeddedStruct(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule007/struct.go",
		Line:   6,
		Column: 10,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "EmbeddedStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule007/subpkg.EmbeddedStruct",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "subpkg",
				Name: "EmbeddedStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule007/subpkg/struct.go"),
					Line:   3,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectFileWithCGo(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule008/cgo.go",
		Line:   6,
		Column: 6,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	expected := &Result{
		Name: "MyStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule008.MyStruct",
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  "testmodule008",
				Name: "MyStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule008/cgo.go"),
					Line:   6,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectInterfaceWithImpl(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule009/iface.go",
		Line:   3,
		Column: 7,
	}, "testdata/testmodule009") // TODO: back to testdata...

	require.NoError(t, err)
	require.NotNil(t, result)

	expected := &Result{
		Name: "MyInterface",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule009.MyInterface",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule009",
				Name: "MyInterface",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/iface.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "interfaceImplementation",
				Pkg:  "subpkg",
				Name: "MyInterfaceSubpkgImpl",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/subpkg/impl.go"),
					Line:   7,
					Column: 6,
				},
			},
			{
				Kind: "interfaceImplementation",
				Pkg:  "testmodule009",
				Name: "MyInterfaceImpl",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "interfaceImplementation",
				Pkg:  "testmodule009",
				Name: "MyInterfacePointerImpl",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   13,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectInterfaceWithImplMethod(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule009/iface.go",
		Line:   5,
		Column: 3,
	}, "testdata/testmodule009") // TODO: back to testdata...

	require.NoError(t, err)
	require.NotNil(t, result)

	expected := &Result{
		Name: "IfaceMethodTwo",
		Type: "func(param github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule009.ParamType) int",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule009",
				Name: "IfaceMethodTwo",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/iface.go"),
					Line:   5,
					Column: 2,
				},
			},
			{
				Kind: "interfaceImplementation",
				Pkg:  "subpkg",
				Name: "IfaceMethodTwo",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/subpkg/impl.go"),
					Line:   13,
					Column: 32,
				},
			},
			{
				Kind: "interfaceImplementation",
				Pkg:  "testmodule009",
				Name: "IfaceMethodTwo",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   9,
					Column: 26,
				},
			},
			{
				Kind: "interfaceImplementation",
				Pkg:  "testmodule009",
				Name: "IfaceMethodTwo",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   19,
					Column: 34,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func BenchmarkInspect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Inspect(file.Loc{
			Path:   "testdata/testmodule003/func.go",
			Line:   15,
			Column: 9,
		}, "")
		require.NoError(b, err)
	}
}

func absPath(t *testing.T, path string) string {
	absPath, err := filepath.Abs(path)
	require.NoError(t, err)
	return absPath
}

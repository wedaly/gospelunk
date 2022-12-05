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
	}, "testdata/testmodule001", AllRelationKinds)
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
	}, "testdata/testmodule002", AllRelationKinds)
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
	}, "testdata/testmodule002", AllRelationKinds)
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
	}, "testdata/testmodule002", AllRelationKinds)
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
	}, "testdata/testmodule004", AllRelationKinds)
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
	}, "testdata/testmodule004", AllRelationKinds)
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
	}, "testdata/testmodule004", AllRelationKinds)
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
	}, "testdata/testmodule004", AllRelationKinds)
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
	}, "testdata/testmodule004", AllRelationKinds)
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
	}, "testdata/testmodule003", AllRelationKinds)
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
	}, "testdata/testmodule003", AllRelationKinds)
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
	}, "testdata/testmodule003", AllRelationKinds)
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
	}, "testdata/testmodule003", AllRelationKinds)
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
	}, "testdata/testmodule005", AllRelationKinds)
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
	}, "testdata/testmodule006", AllRelationKinds)
	assert.EqualError(t, err, "Could not find AST node of type *ast.Ident at location testdata/testmodule006/const.go:5:18")
	assert.Nil(t, result)
}

func TestInspectIntegerConst(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule006/const.go",
		Line:   13,
		Column: 23,
	}, "testdata/testmodule006", AllRelationKinds)
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
	}, "testdata/testmodule006", AllRelationKinds)
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
	}, "testdata/testmodule007", AllRelationKinds)
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
	}, "testdata/testmodule008", AllRelationKinds)
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
	}, "testdata/testmodule009", AllRelationKinds)

	require.NoError(t, err)
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
				Kind: "reference",
				Pkg:  "testmodule009",
				Name: "MyInterface in declaration of _",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   23,
					Column: 7,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule009",
				Name: "MyInterface in declaration of _",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   24,
					Column: 7,
				},
			},
			{
				Kind: "reference",
				Pkg:  "subpkg",
				Name: "MyInterface in declaration of _",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/subpkg/impl.go"),
					Line:   17,
					Column: 9,
				},
			},
			{
				Kind: "implementation",
				Pkg:  "testmodule009",
				Name: "MyInterfaceImpl",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "implementation",
				Pkg:  "testmodule009",
				Name: "MyInterfacePointerImpl",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   13,
					Column: 6,
				},
			},
			{
				Kind: "implementation",
				Pkg:  "subpkg",
				Name: "MyInterfaceSubpkgImpl",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/subpkg/impl.go"),
					Line:   7,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectImplOfInterface(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule009/impl.go",
		Line:   3,
		Column: 7,
	}, "testdata/testmodule009", []RelationKind{RelationKindIface})

	require.NoError(t, err)
	expected := &Result{
		Name: "MyInterfaceImpl",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule009.MyInterfaceImpl",
		Relations: []Relation{
			{
				Kind: "interface",
				Pkg:  "testmodule009",
				Name: "MyInterface",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/iface.go"),
					Line:   3,
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
	}, "testdata/testmodule009", AllRelationKinds)

	require.NoError(t, err)
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
				Kind: "implementation",
				Pkg:  "testmodule009",
				Name: "MyInterfaceImpl.IfaceMethodTwo()",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   9,
					Column: 26,
				},
			},
			{
				Kind: "implementation",
				Pkg:  "testmodule009",
				Name: "MyInterfacePointerImpl.IfaceMethodTwo()",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/impl.go"),
					Line:   19,
					Column: 34,
				},
			},
			{
				Kind: "implementation",
				Pkg:  "subpkg",
				Name: "MyInterfaceSubpkgImpl.IfaceMethodTwo()",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/subpkg/impl.go"),
					Line:   13,
					Column: 32,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectImplOfInterfaceMethod(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule009/impl.go",
		Line:   5,
		Column: 26,
	}, "testdata/testmodule009", []RelationKind{RelationKindIface})

	require.NoError(t, err)
	expected := &Result{
		Name: "IfaceMethodOne",
		Type: "func() string",
		Relations: []Relation{
			{
				Kind: "interface",
				Pkg:  "testmodule009",
				Name: "MyInterface.IfaceMethodOne",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule009/iface.go"),
					Line:   4,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectInterfaceWithImplAndIfaceInDifferentPkgs(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule010/subpkgWithIface/iface.go",
		Line:   5,
		Column: 7,
	}, "testdata/testmodule010", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyInterface",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule010/subpkgWithIface.MyInterface",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "subpkgWithIface",
				Name: "MyInterface",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule010/subpkgWithIface/iface.go"),
					Line:   5,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "subpkgWithIface",
				Name: "MyInterface in Print() params",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule010/subpkgWithIface/iface.go"),
					Line:   9,
					Column: 14,
				},
			},
			{
				Kind: "implementation",
				Pkg:  "subpkgWithImpl",
				Name: "MyInterfaceImpl",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule010/subpkgWithImpl/impl.go"),
					Line:   3,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructWithReference(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule011/subpkg/def.go",
		Line:   3,
		Column: 7,
	}, "testdata/testmodule011", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule011/subpkg.MyStruct",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "subpkg",
				Name: "MyStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "main",
				Name: "MyStruct in main() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/main.go"),
					Line:   8,
					Column: 14,
				},
			},
			{
				Kind: "reference",
				Pkg:  "subpkg",
				Name: "receiver in MyStruct.MyMethod()",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   15,
					Column: 9,
				},
			},
			{
				Kind: "reference",
				Pkg:  "subpkg",
				Name: "MyStruct in PrintMyStruct() params",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/ref.go"),
					Line:   7,
					Column: 22,
				},
			},
			{
				Kind: "interface",
				Pkg: "subpkg",
				Name: "MyInterface",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   11,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructFieldWithReference(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule011/subpkg/def.go",
		Line:   4,
		Column: 2,
	}, "testdata/testmodule011", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "Foo",
		Type: "string",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "subpkg",
				Name: "Foo",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   4,
					Column: 2,
				},
			},
			{
				Kind: "reference",
				Pkg:  "main",
				Name: "Foo in main() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/main.go"),
					Line:   8,
					Column: 23,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructMethodWithReference(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule011/subpkg/def.go",
		Line:   15,
		Column: 20,
	}, "testdata/testmodule011", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyMethod",
		Type: "func() string",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "subpkg",
				Name: "MyMethod",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   15,
					Column: 19,
				},
			},
			{
				Kind: "reference",
				Pkg:  "subpkg",
				Name: "MyMethod in PrintMyStruct() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/ref.go"),
					Line:   9,
					Column: 21,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectFunctionWithReference(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule011/subpkg/def.go",
		Line:   7,
		Column: 7,
	}, "testdata/testmodule011", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyFunc",
		Type: "func() string",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "subpkg",
				Name: "MyFunc",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   7,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "subpkg",
				Name: "MyFunc in CallMyFunc() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/ref.go"),
					Line:   13,
					Column: 7,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectIfaceMethodWithReference(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule011/subpkg/def.go",
		Line:   12,
		Column: 3,
	}, "testdata/testmodule011", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyMethod",
		Type: "func() string",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "subpkg",
				Name: "MyMethod",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   12,
					Column: 2,
				},
			},
			{
				Kind: "reference",
				Pkg:  "subpkg",
				Name: "MyMethod in PrintMyInterface() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/ref.go"),
					Line:   18,
					Column: 21,
				},
			},
			{
				Kind: "implementation",
				Pkg:  "subpkg",
				Name: "MyStruct.MyMethod()",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule011/subpkg/def.go"),
					Line:   15,
					Column: 19,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructReturnedFromStructMethod(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule012/main.go",
		Line:   14,
		Column: 7,
	}, "testdata/testmodule012", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "SecondStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule012.SecondStruct",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "main",
				Name: "SecondStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule012/main.go"),
					Line:   14,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "main",
				Name: "SecondStruct returned by FirstStruct.ReturnSecondStruct()",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule012/main.go"),
					Line:   18,
					Column: 43,
				},
			},
			{
				Kind: "reference",
				Pkg:  "main",
				Name: "SecondStruct in FirstStruct.ReturnSecondStruct() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule012/main.go"),
					Line:   19,
					Column: 9,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructFieldFromReturnValue(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule012/main.go",
		Line:   15,
		Column: 2,
	}, "testdata/testmodule012", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "B",
		Type: "string",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "main",
				Name: "B",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule012/main.go"),
					Line:   15,
					Column: 2,
				},
			},
			{
				Kind: "reference",
				Pkg:  "main",
				Name: "B in main() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule012/main.go"),
					Line:   6,
					Column: 50,
				},
			},
			{
				Kind: "reference",
				Pkg:  "main",
				Name: "B in FirstStruct.ReturnSecondStruct() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule012/main.go"),
					Line:   19,
					Column: 22,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructWithVarDeclaration(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule013/values.go",
		Line:   3,
		Column: 7,
	}, "testdata/testmodule013", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule013.MyStruct",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule013",
				Name: "MyStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/values.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyStruct in declaration of myStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/values.go"),
					Line:   8,
					Column: 11,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectTypeWithVarAndConstDeclaration(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule013/values.go",
		Line:   5,
		Column: 7,
	}, "testdata/testmodule013", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyInt",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule013.MyInt",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule013",
				Name: "MyInt",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/values.go"),
					Line:   5,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyInt in declaration of x, y",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/values.go"),
					Line:   9,
					Column: 11,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyInt in declaration of z",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/values.go"),
					Line:   12,
					Column: 9,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectTypeWithReferenceInFuncTypeArg(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule013/func.go",
		Line:   5,
		Column: 7,
	}, "testdata/testmodule013", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyArg",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule013.MyArg",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule013",
				Name: "MyArg",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/func.go"),
					Line:   5,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyArg in params for func type MyFunc",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/func.go"),
					Line:   9,
					Column: 20,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyArg in ConstructStructWithFuncField() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/func.go"),
					Line:   17,
					Column: 21,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectTypeWithReferenceInFuncTypeReturn(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule013/func.go",
		Line:   7,
		Column: 7,
	}, "testdata/testmodule013", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyReturn",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule013.MyReturn",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule013",
				Name: "MyReturn",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/func.go"),
					Line:   7,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyReturn returned by func type MyFunc",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/func.go"),
					Line:   9,
					Column: 27,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyReturn in ConstructStructWithFuncField() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/func.go"),
					Line:   17,
					Column: 28,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyReturn in ConstructStructWithFuncField() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/func.go"),
					Line:   18,
					Column: 11,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectTypeWithEmbeddedStructRef(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule013/embed.go",
		Line:   3,
		Column: 7,
	}, "testdata/testmodule013", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyEmbeddedStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule013.MyEmbeddedStruct",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule013",
				Name: "MyEmbeddedStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/embed.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyEmbeddedStruct embedded in struct MyParentStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/embed.go"),
					Line:   6,
					Column: 2,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectTypeReturnedByMethodWithPointerReceiver(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule013/method.go",
		Line:   3,
		Column: 7,
	}, "testdata/testmodule013", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyMethodReturn",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule013.MyMethodReturn",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule013",
				Name: "MyMethodReturn",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/method.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyMethodReturn returned by MyStructWithPointerRecv.MyMethod()",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/method.go"),
					Line:   7,
					Column: 46,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule013",
				Name: "MyMethodReturn in MyStructWithPointerRecv.MyMethod() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule013/method.go"),
					Line:   8,
					Column: 9,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectLocalVar(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule014/localref.go",
		Line:   6,
		Column: 2,
	}, "testdata/testmodule014", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "x",
		Type: "string",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule014",
				Name: "x",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule014/localref.go"),
					Line:   6,
					Column: 2,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule014",
				Name: "x in funcWithLocalVar() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule014/localref.go"),
					Line:   7,
					Column: 19,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectPrivateVar(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule014/privatevar.go",
		Line:   5,
		Column: 5,
	}, "testdata/testmodule014", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "privateVar",
		Type: "int",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule014",
				Name: "privateVar",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule014/privatevar.go"),
					Line:   5,
					Column: 5,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule014",
				Name: "privateVar in funcUsingPrivateVar() body",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule014/privatevar.go"),
					Line:   8,
					Column: 19,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructWithReferenceInTest(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule015/def.go",
		Line:   3,
		Column: 6,
	}, "testdata/testmodule015", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule015.MyStruct",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule015",
				Name: "MyStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule015/def.go"),
					Line:   3,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule015",
				Name: "MyStruct in MyFunc() params",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule015/def.go"),
					Line:   7,
					Column: 15,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectReferenceToStructInTest(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule015/def_test.go",
		Line:   10,
		Column: 7,
	}, "testdata/testmodule015", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule015.MyStruct",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule015",
				Name: "MyStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule015/def.go"),
					Line:   3,
					Column: 6,
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}

func TestInspectStructDefinedAndReferencedInTest(t *testing.T) {
	result, err := Inspect(file.Loc{
		Path:   "testdata/testmodule015/def_test.go",
		Line:   5,
		Column: 6,
	}, "testdata/testmodule015", AllRelationKinds)

	require.NoError(t, err)
	expected := &Result{
		Name: "MyTestStruct",
		Type: "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule015.MyTestStruct",
		Relations: []Relation{
			{
				Kind: "definition",
				Pkg:  "testmodule015",
				Name: "MyTestStruct",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule015/def_test.go"),
					Line:   5,
					Column: 6,
				},
			},
			{
				Kind: "reference",
				Pkg:  "testmodule015",
				Name: "MyTestStruct in declaration of testVar",
				Loc: file.Loc{
					Path:   absPath(t, "testdata/testmodule015/def_test.go"),
					Line:   7,
					Column: 13,
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
		}, "", AllRelationKinds)
		require.NoError(b, err)
	}
}

func absPath(t *testing.T, path string) string {
	absPath, err := filepath.Abs(path)
	require.NoError(t, err)
	return absPath
}

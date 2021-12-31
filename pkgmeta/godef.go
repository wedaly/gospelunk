package pkgmeta

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"unicode"

	"github.com/pkg/errors"
)

// GoDefKind represents the kind of a Go definition (function, value, interface, etc.)
// This is less specific than the full Go type.
type GoDefKind int

const (
	GoDefKindValue = GoDefKind(iota)
	GoDefKindFunc
	GoDefKindTypeStruct
	GoDefKindStructField
	GoDefKindStructMethod
	GoDefKindTypeInterface
	GoDefKindInterfaceMethod
	GoDefKindTypeOther
)

func (k GoDefKind) String() string {
	switch k {
	case GoDefKindValue:
		return "value"
	case GoDefKindFunc:
		return "func"
	case GoDefKindTypeStruct:
		return "struct"
	case GoDefKindStructField:
		return "field"
	case GoDefKindStructMethod:
		return "method"
	case GoDefKindTypeInterface:
		return "interface"
	case GoDefKindInterfaceMethod:
		return "method"
	case GoDefKindTypeOther:
		return "type"
	default:
		return "unknown"
	}
}

// GoDef represents a top-level definition in a Go package.
type GoDef struct {
	Name     string
	Kind     GoDefKind
	LineNum  int
	Exported bool
}

// LoadGoDefs loads top-level Go definitions from a Go source file.
func LoadGoDefs(reader io.Reader) ([]GoDef, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", reader, parser.SkipObjectResolution)
	if err != nil {
		return nil, errors.Wrapf(err, "parser.ParseFile")
	}

	var defs []GoDef
	ast.Inspect(f, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.ValueSpec:
			loadDefsFromValueSpec(x, fset, &defs)
			return false

		case *ast.TypeSpec:
			loadDefsFromTypeSpec(x, fset, &defs)
			return false

		case *ast.FuncDecl:
			loadDefsFromFuncDecl(x, fset, &defs)
			return false

		default:
			return true
		}
	})

	return defs, nil
}

func loadDefsFromValueSpec(valueSpec *ast.ValueSpec, fset *token.FileSet, defs *[]GoDef) {
	lineNum := fset.Position(valueSpec.Pos()).Line
	for _, nameIdent := range valueSpec.Names {
		if nameIdent != nil {
			valueName := nameIdent.Name
			*defs = append(*defs, GoDef{
				Name:     valueName,
				Kind:     GoDefKindValue,
				LineNum:  lineNum,
				Exported: isExported(valueName),
			})
		}
	}
}

func loadDefsFromTypeSpec(typeSpec *ast.TypeSpec, fset *token.FileSet, defs *[]GoDef) {
	if typeSpec.Name == nil {
		return
	}

	lineNum := fset.Position(typeSpec.Pos()).Line
	typeName := typeSpec.Name.Name
	kind := GoDefKindTypeOther

	switch x := typeSpec.Type.(type) {
	case *ast.StructType:
		kind = GoDefKindTypeStruct
		loadDefsFromStructType(x, fset, typeName, defs)
	case *ast.InterfaceType:
		kind = GoDefKindTypeInterface
		loadDefsFromInterfaceType(x, fset, typeName, defs)
	}

	*defs = append(*defs, GoDef{
		Name:     typeName,
		Kind:     kind,
		LineNum:  lineNum,
		Exported: isExported(typeName),
	})
}

func loadDefsFromStructType(structType *ast.StructType, fset *token.FileSet, typeName string, defs *[]GoDef) {
	for _, field := range structType.Fields.List {
		lineNum := fset.Position(field.Pos()).Line
		for _, nameIdent := range field.Names {
			if nameIdent != nil {
				fieldName := nameIdent.Name
				*defs = append(*defs, GoDef{
					Name:     fmt.Sprintf("%s.%s", typeName, fieldName),
					Kind:     GoDefKindStructField,
					LineNum:  lineNum,
					Exported: isExported(fieldName),
				})
			}
		}
	}
}

func loadDefsFromInterfaceType(interfaceType *ast.InterfaceType, fset *token.FileSet, typeName string, defs *[]GoDef) {
	for _, method := range interfaceType.Methods.List {
		lineNum := fset.Position(method.Pos()).Line
		for _, nameIdent := range method.Names {
			if nameIdent != nil {
				methodName := nameIdent.Name
				*defs = append(*defs, GoDef{
					Name:     fmt.Sprintf("%s.%s", typeName, methodName),
					Kind:     GoDefKindInterfaceMethod,
					LineNum:  lineNum,
					Exported: isExported(methodName),
				})
			}
		}
	}
}

func loadDefsFromFuncDecl(funcDecl *ast.FuncDecl, fset *token.FileSet, defs *[]GoDef) {
	if funcDecl.Name == nil {
		return
	}
	lineNum := fset.Position(funcDecl.Pos()).Line
	funcName := funcDecl.Name.Name
	name, kind := funcName, GoDefKindFunc
	if funcDecl.Recv != nil {
		name = fmt.Sprintf("%s.%s", findFuncRecvName(funcDecl), funcName)
		kind = GoDefKindStructMethod
	}
	*defs = append(*defs, GoDef{
		Name:     name,
		Kind:     kind,
		LineNum:  lineNum,
		Exported: isExported(funcName),
	})
}

func findFuncRecvName(funcDecl *ast.FuncDecl) string {
	var typeName string
	for _, field := range funcDecl.Recv.List {
		ast.Inspect(field.Type, func(node ast.Node) bool {
			if ident, ok := node.(*ast.Ident); ok {
				typeName = ident.Name
				return false
			}
			return true
		})
	}
	return typeName
}

func isExported(name string) bool {
	// Check if the first letter in the name is uppercase.
	for _, r := range name {
		if unicode.IsUpper(r) {
			return true
		} else {
			break
		}
	}
	return false
}

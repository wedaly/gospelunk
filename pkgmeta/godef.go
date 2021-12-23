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

// GoDef represents a top-level definition in a Go package.
type GoDef struct {
	Name     string
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
	*defs = append(*defs, GoDef{
		Name:     typeName,
		LineNum:  lineNum,
		Exported: isExported(typeName),
	})

	switch x := typeSpec.Type.(type) {
	case *ast.StructType:
		loadDefsFromStructType(x, fset, typeName, defs)
	case *ast.InterfaceType:
		loadDefsFromInterfaceType(x, fset, typeName, defs)
	}
}

func loadDefsFromStructType(structType *ast.StructType, fset *token.FileSet, typeName string, defs *[]GoDef) {
	for _, field := range structType.Fields.List {
		lineNum := fset.Position(field.Pos()).Line
		for _, nameIdent := range field.Names {
			if nameIdent != nil {
				fieldName := nameIdent.Name
				*defs = append(*defs, GoDef{
					Name:     fmt.Sprintf("%s.%s", typeName, fieldName),
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
	*defs = append(*defs, GoDef{
		Name:     funcName,
		LineNum:  lineNum,
		Exported: isExported(funcName),
	})
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
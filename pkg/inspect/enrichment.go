package inspect

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/wedaly/gospelunk/pkg/file"
)

type enrichResultFunc func(result *Result, pkg *packages.Package, loc file.Loc) error

func enrichResultNameAndType(result *Result, pkg *packages.Package, loc file.Loc) error {
	ident, err := astNodeAtLoc[*ast.Ident](pkg, loc)
	if err != nil {
		return err
	}

	obj, err := typeObjUseOrDefForAstIdent(ident, pkg)
	if err != nil {
		return err
	}

	var typeName string
	if obj.Type() != types.Typ[types.Invalid] {
		typeName = obj.Type().String()
	}

	result.Name = ident.Name
	result.Type = typeName
	return nil
}

func enrichResultDefRelation(result *Result, pkg *packages.Package, loc file.Loc) error {
	ident, err := astNodeAtLoc[*ast.Ident](pkg, loc)
	if err != nil {
		return err
	}

	obj, err := typeObjUseOrDefForAstIdent(ident, pkg)
	if err != nil {
		// If we can't lookup the Use or Def for this object,
		// assume it has no definition relation.
		return nil
	}

	if !obj.Pos().IsValid() {
		return nil
	}

	result.Relations = append(result.Relations, Relation{
		Kind: RelationKindDef,
		Pkg:  pkgNameForTypeObj(obj),
		Name: obj.Name(),
		Loc:  fileLocForTypeObj(pkg, obj),
	})

	return nil
}

func typeObjUseOrDefForAstIdent(ident *ast.Ident, pkg *packages.Package) (types.Object, error) {
	obj, ok := pkg.TypesInfo.Uses[ident]
	if !ok {
		obj = pkg.TypesInfo.Defs[ident]
	}

	if obj == nil {
		return nil, fmt.Errorf("AST identifier %q has no type object use or definition", ident)
	}

	return obj, nil
}

func pkgNameForTypeObj(obj types.Object) string {
	if obj.Pkg() == nil {
		return ""
	}
	return obj.Pkg().Name()
}

func fileLocForTypeObj(pkg *packages.Package, obj types.Object) file.Loc {
	if !obj.Pos().IsValid() {
		return file.Loc{}
	}
	position := pkg.Fset.Position(obj.Pos())
	return file.Loc{
		Path:   position.Filename,
		Line:   position.Line,
		Column: position.Column,
	}
}

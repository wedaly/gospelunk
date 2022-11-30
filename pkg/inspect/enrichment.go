package inspect

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"

	"golang.org/x/tools/go/packages"

	"github.com/wedaly/gospelunk/pkg/file"
)

type enrichResultFunc func(result *Result, pkg *packages.Package, loc file.Loc, searchDir string) error

func enrichResultNameAndType(result *Result, pkg *packages.Package, loc file.Loc, searchDir string) error {
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

func enrichResultDefRelation(result *Result, pkg *packages.Package, loc file.Loc, searchDir string) error {
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

func enrichResultIfaceImplRelation(result *Result, pkg *packages.Package, loc file.Loc, searchDir string) error {
	ifaceName, ifaceType := interfaceNameAndTypeAtFileLoc(pkg, loc)
	if ifaceType == nil || ifaceType.Empty() {
		return nil
	}

	methodName := methodNameForInterfaceAtLoc(pkg, loc, ifaceType) // Empty string if not on method identifier.

	searchPkgs, err := loadGoPackagesEqualToOrImportingPkg(pkg.ID, searchDir)
	if err != nil {
		return err
	}

	relationSet := make(map[Relation]struct{}, 0)
	for _, searchPkg := range searchPkgs {
		// Lookup the interface type either in the package or its imports.
		// We need this to check if other types in the package implement the interface.
		// (We can't use ifaceType because it comes from a different package, so isn't comparable to types in this pkg.)
		pkgIfaceType := interfaceTypeInPkgScopeWithName(searchPkg, ifaceName)
		if searchPkg.ID == pkg.ID {
			pkgIfaceType = interfaceTypeInPkgScopeWithName(searchPkg, ifaceName)
		} else if importedPkg, ok := searchPkg.Imports[pkg.ID]; ok {
			pkgIfaceType = interfaceTypeInPkgScopeWithName(importedPkg, ifaceName)
		}

		if pkgIfaceType == nil {
			continue
		}

		// Search every reference in this package for implementations of the interface.
		seen := make(map[types.Object]struct{}, 0)
		for _, obj := range searchPkg.TypesInfo.Uses {
			if obj == nil || obj.Type() == types.Typ[types.Invalid] {
				continue
			}

			if _, ok := seen[obj]; ok {
				// Skip objects we've already processed.
				continue
			}
			seen[obj] = struct{}{}

			if _, ok := obj.Type().(*types.Named); !ok {
				// Filter for only named types.
				continue
			}

			if _, ok := obj.(*types.Var); ok {
				// Exclude variables (including method receivers).
				continue
			}

			if types.Identical(obj.Type().Underlying(), pkgIfaceType) {
				// Interfaces always implement themselves, so skip the one we're looking for.
				continue
			}

			// Check if this type OR a pointer to this type implements the interface.
			if types.Implements(obj.Type(), pkgIfaceType) || types.Implements(types.NewPointer(obj.Type()), pkgIfaceType) {
				if methodName == "" {
					// If we're not looking for a specific method, the relation points to the implementation of the interface type.
					r := Relation{
						Kind: RelationKindIfaceImpl,
						Pkg:  pkgNameForTypeObj(obj),
						Name: obj.Name(),
						Loc:  fileLocForTypeObj(searchPkg, obj),
					}
					relationSet[r] = struct{}{}
				} else {
					// If we're looking for a specific method, the relation points to the implementation of the method.
					methodObj, _, _ := types.LookupFieldOrMethod(obj.Type(), true, searchPkg.Types, methodName)
					if methodObj != nil {
						r := Relation{
							Kind: RelationKindIfaceImpl,
							Pkg:  pkgNameForTypeObj(methodObj),
							Name: methodObj.Name(),
							Loc:  fileLocForTypeObj(searchPkg, methodObj),
						}
						relationSet[r] = struct{}{}
					}
				}
			}
		}
	}

	// Ensure relations returned in a consistent order.
	relations := make(RelationSlice, 0, len(relationSet))
	for r := range relationSet {
		relations = append(relations, r)
	}
	sort.Stable(relations)
	result.Relations = append(result.Relations, relations...)
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

func interfaceNameAndTypeAtFileLoc(pkg *packages.Package, loc file.Loc) (string, *types.Interface) {
	typeSpec, err := astNodeAtLoc[*ast.TypeSpec](pkg, loc)
	if err != nil {
		// If this isn't a typespec, it can't be an interface declaration.
		return "", nil
	}

	ifaceDefObj, ok := pkg.TypesInfo.Defs[typeSpec.Name]
	if !ok {
		// Couldn't find the definition for the typespec name.
		return "", nil
	}

	ifaceType, ok := ifaceDefObj.Type().Underlying().(*types.Interface)
	if !ok {
		// This is a typeSpec, but not for an interface.
		return "", nil
	}

	return ifaceDefObj.Name(), ifaceType
}

func interfaceTypeInPkgScopeWithName(pkg *packages.Package, name string) *types.Interface {
	ifaceDefObj := pkg.Types.Scope().Lookup(name)
	if ifaceDefObj == nil {
		return nil
	}

	ifaceType, ok := ifaceDefObj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil
	}

	return ifaceType
}

func methodNameForInterfaceAtLoc(pkg *packages.Package, loc file.Loc, ifaceType *types.Interface) string {
	methodIdent, err := astNodeAtLoc[*ast.Ident](pkg, loc)
	if err != nil || methodIdent == nil {
		return ""
	}

	methodObj, _, _ := types.LookupFieldOrMethod(ifaceType, true, pkg.Types, methodIdent.Name)
	if methodObj == nil {
		return ""
	}

	return methodObj.Name()
}

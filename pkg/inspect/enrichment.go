package inspect

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strings"

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

func enrichResultRefRelation(result *Result, pkg *packages.Package, loc file.Loc, searchDir string) error {
	ident, err := astNodeAtLoc[*ast.Ident](pkg, loc)
	if err != nil {
		return err
	}

	if _, ok := pkg.TypesInfo.Defs[ident]; !ok {
		return nil
	}

	targetPosition := pkg.Fset.Position(ident.Pos())

	loadMode := (packages.NeedName |
		packages.NeedSyntax |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedTypesInfo)

	includeTests := isGoTestFile(loc.Path)
	searchPkgs, err := loadGoPackagesMatchingPredicate(searchDir, loadMode, includeTests, func(candidate skeletonPkg) bool {
		return candidate.ImportPath == pkg.PkgPath || (ident.IsExported() && candidate.ImportsPkg(pkg.PkgPath))
	})
	if err != nil {
		return err
	}

	relationSet := make(map[Relation]struct{})
	for _, searchPkg := range searchPkgs {
		for refIdent, refObj := range searchPkg.TypesInfo.Uses {
			refPosition := searchPkg.Fset.Position(refObj.Pos())
			if refPosition != targetPosition {
				continue
			}

			r := Relation{
				Kind: RelationKindRef,
				Pkg:  searchPkg.Name,
				Name: nameForRefRelation(searchPkg, refIdent.Pos(), refIdent.Name),
				Loc:  fileLocForIdent(searchPkg, refIdent),
			}
			relationSet[r] = struct{}{}
		}
	}

	result.Relations = append(result.Relations, relationSetToSortedSlice(relationSet)...)
	return nil
}

func nameForRefRelation(pkg *packages.Package, pos token.Pos, identName string) string {
	path := pkg.Fset.Position(pos).Filename
	root, err := astFileForPath(pkg, path)
	if err != nil {
		return identName
	}

	refName := identName
	ast.Inspect(root, func(node ast.Node) bool {
		if node == nil || pos < node.Pos() || pos > node.End() {
			return false
		}

		switch node := node.(type) {
		case *ast.ValueSpec:
			valueNameStrings := make([]string, 0, len(node.Names))
			for _, valueNameIdent := range node.Names {
				valueNameStrings = append(valueNameStrings, valueNameIdent.Name)
			}
			refName = fmt.Sprintf("%s in declaration of %s", identName, strings.Join(valueNameStrings, ", "))
			return false

		case *ast.FuncDecl:
			funcName := node.Name.Name
			if node.Recv != nil && len(node.Recv.List) > 0 {
				for _, field := range node.Recv.List {
					if recvTypeIdent, ok := field.Type.(*ast.Ident); ok {
						funcName = fmt.Sprintf("%s.%s", recvTypeIdent.Name, funcName)
						break
					} else if recvTypeStar, ok := field.Type.(*ast.StarExpr); ok {
						if recvTypeIdent, ok := recvTypeStar.X.(*ast.Ident); ok {
							funcName = fmt.Sprintf("%s.%s", recvTypeIdent.Name, funcName)
							break
						}
					}
				}
			}

			if node.Body != nil && !(pos < node.Body.Lbrace || pos > node.Body.Rbrace) {
				refName = fmt.Sprintf("%s in %s() body", identName, funcName)
			} else if !(pos < node.Type.Params.Opening || pos > node.Type.Params.Closing) {
				refName = fmt.Sprintf("%s in %s() params", identName, funcName)
			} else if node.Recv != nil && !(pos < node.Recv.Opening || pos > node.Recv.Closing) {
				refName = fmt.Sprintf("receiver in %s()", funcName)
			} else {
				for _, resultField := range node.Type.Results.List {
					if pos >= resultField.Pos() && pos <= resultField.End() {
						refName = fmt.Sprintf("%s returned by %s()", identName, funcName)
						break
					}
				}
			}

			return false

		case *ast.TypeSpec:
			if structType, ok := node.Type.(*ast.StructType); ok {
				for _, field := range structType.Fields.List {
					if pos >= field.Pos() || pos <= field.End() {
						if len(field.Names) == 0 {
							refName = fmt.Sprintf("%s embedded in struct %s", identName, node.Name.Name)
						}
						break
					}
				}
			} else if funcType, ok := node.Type.(*ast.FuncType); ok {
				if !(pos < funcType.Params.Opening || pos > funcType.Params.Closing) {
					refName = fmt.Sprintf("%s in params for func type %s", identName, node.Name.Name)
				} else {
					for _, resultField := range funcType.Results.List {
						if pos >= resultField.Pos() && pos <= resultField.End() {
							refName = fmt.Sprintf("%s returned by func type %s", identName, node.Name.Name)
							break
						}
					}
				}
			}
			return true

		default:
			return true
		}
	})

	return refName
}

func enrichResultImplRelation(result *Result, pkg *packages.Package, loc file.Loc, searchDir string) error {
	ifaceName, ifaceType := interfaceNameAndTypeAtFileLoc(pkg, loc)
	if ifaceType == nil || ifaceType.Empty() {
		return nil
	}

	methodName := methodNameForTypeAtLoc(pkg, loc, ifaceType) // Empty string if not on method identifier.

	loadMode := (packages.NeedName |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedTypesInfo |
		packages.NeedImports)

	includeTests := isGoTestFile(loc.Path)
	searchPkgs, err := loadGoPackagesMatchingPredicate(searchDir, loadMode, includeTests, func(candidate skeletonPkg) bool {
		return candidate.ImportPath == pkg.PkgPath || candidate.ImportsPkg(pkg.PkgPath)
	})
	if err != nil {
		return err
	}

	relationSet := make(map[Relation]struct{})
	for _, searchPkg := range searchPkgs {
		// Lookup the interface type either in the package or its imports.
		// We need this to check if other types in the package implement the interface.
		// (We can't use ifaceType directly because it comes from a different package, so it isn't comparable to types in this pkg.)
		var pkgIfaceType *types.Interface
		if searchPkg.PkgPath == pkg.PkgPath {
			pkgIfaceType = interfaceTypeInPkgScopeWithName(searchPkg, ifaceName)
		} else if importedPkg, ok := searchPkg.Imports[pkg.PkgPath]; ok {
			pkgIfaceType = interfaceTypeInPkgScopeWithName(importedPkg, ifaceName)
		}

		if pkgIfaceType == nil {
			continue
		}

		// Search every reference in this package for implementations of the interface.
		seen := make(map[types.Object]struct{})
		for _, obj := range searchPkg.TypesInfo.Uses {
			if obj == nil || obj.Type() == types.Typ[types.Invalid] {
				continue
			}

			if _, ok := obj.(*types.TypeName); !ok {
				// Filter for only type names.
				continue
			}

			if _, ok := seen[obj]; ok {
				// Skip objects we've already processed.
				continue
			}
			seen[obj] = struct{}{}

			if types.Identical(obj.Type().Underlying(), pkgIfaceType) {
				// Interfaces always implement themselves, so skip the one we're looking for.
				continue
			}

			// Check if this type OR a pointer to this type implements the interface.
			if types.Implements(obj.Type(), pkgIfaceType) || types.Implements(types.NewPointer(obj.Type()), pkgIfaceType) {
				if methodName == "" {
					// If we're not looking for a specific method, the relation points to the implementation of the interface type.
					r := Relation{
						Kind: RelationKindImpl,
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
							Kind: RelationKindImpl,
							Pkg:  pkgNameForTypeObj(methodObj),
							Name: fmt.Sprintf("%s.%s()", obj.Name(), methodObj.Name()),
							Loc:  fileLocForTypeObj(searchPkg, methodObj),
						}
						relationSet[r] = struct{}{}
					}
				}
			}
		}
	}

	result.Relations = append(result.Relations, relationSetToSortedSlice(relationSet)...)
	return nil
}

func enrichResultIfaceRelation(result *Result, pkg *packages.Package, loc file.Loc, searchDir string) error {
	if typeSpec, err := astNodeAtLoc[*ast.TypeSpec](pkg, loc); err == nil {
		return enrichResultIfaceRelationFromTypeSpec(result, pkg, loc, searchDir, typeSpec)
	} else if funcDecl, err := astNodeAtLoc[*ast.FuncDecl](pkg, loc); err == nil {
		return enrichResultIfaceRelationFromFuncDecl(result, pkg, loc, searchDir, funcDecl)
	}

	return nil
}

func enrichResultIfaceRelationFromTypeSpec(result *Result, pkg *packages.Package, loc file.Loc, searchDir string, typeSpec *ast.TypeSpec) error {
	ident, err := astNodeAtLoc[*ast.Ident](pkg, loc)
	if err != nil || ident != typeSpec.Name {
		// Not on the name of the typespec, so skip it.
		return nil
	}

	implObj, ok := pkg.TypesInfo.Defs[typeSpec.Name]
	if !ok {
		// Can't find the definition, so skip it.
		return nil
	}

	implType := implObj.Type().Underlying()

	if _, ok := implType.(*types.Interface); ok {
		// Skip if this is an interface.
		return nil
	}

	loadMode := (packages.NeedName |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedTypesInfo |
		packages.NeedImports)

	includeTests := isGoTestFile(loc.Path)
	searchPkgs, err := loadGoPackagesMatchingPredicate(searchDir, loadMode, includeTests, func(candidate skeletonPkg) bool {
		return candidate.ImportPath == pkg.PkgPath || candidate.ImportsPkg(pkg.PkgPath)
	})
	if err != nil {
		return err
	}

	relationSet := make(map[Relation]struct{}, 0)
	for _, searchPkg := range searchPkgs {
		// Lookup the impl type either in the package or its imports.
		// We need this to find interfaces in this package that implement the target implementation.
		// (We can't use implType directly because it comes from a different package, so it isn't comparable to types in this pkg.)
		var pkgImplType types.Type
		if searchPkg.PkgPath == pkg.PkgPath {
			pkgImplType = implTypeInPkgScopeWithName(searchPkg, implObj.Name())
		} else if importedPkg, ok := searchPkg.Imports[pkg.PkgPath]; ok {
			pkgImplType = implTypeInPkgScopeWithName(importedPkg, implObj.Name())
		}

		if pkgImplType == nil {
			continue
		}

		// Search every reference in this package for an interface implementing the target type.
		seen := make(map[types.Object]struct{})
		for _, obj := range searchPkg.TypesInfo.Uses {
			if obj == nil || obj.Type() == types.Typ[types.Invalid] {
				continue
			}

			ifaceType, ok := obj.Type().Underlying().(*types.Interface)
			if !ok {
				// Not an interface.
				continue
			}

			if _, ok := seen[obj]; ok {
				// Skip intefaces we've already seen.
				continue
			}

			if _, ok := obj.Type().(*types.Named); !ok {
				// Filter for only named types.
				continue
			}

			if _, ok := obj.(*types.Var); ok {
				// Exclude variables (including method receivers).
				continue
			}

			// Check if the interface implements this type OR a pointer to this type.
			if types.Implements(pkgImplType, ifaceType) || types.Implements(types.NewPointer(pkgImplType), ifaceType) {
				r := Relation{
					Kind: RelationKindIface,
					Pkg:  pkgNameForTypeObj(obj),
					Name: obj.Name(),
					Loc:  fileLocForTypeObj(searchPkg, obj),
				}
				relationSet[r] = struct{}{}
			}
		}
	}

	result.Relations = append(result.Relations, relationSetToSortedSlice(relationSet)...)
	return nil
}

func enrichResultIfaceRelationFromFuncDecl(result *Result, pkg *packages.Package, loc file.Loc, searchDir string, funcDecl *ast.FuncDecl) error {
	// TODO
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

func fileLocForIdent(pkg *packages.Package, ident *ast.Ident) file.Loc {
	position := pkg.Fset.Position(ident.Pos())
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

func implTypeInPkgScopeWithName(pkg *packages.Package, name string) types.Type {
	implDefObj := pkg.Types.Scope().Lookup(name)
	if implDefObj == nil {
		return nil
	}
	return implDefObj.Type()
}

func methodNameForTypeAtLoc(pkg *packages.Package, loc file.Loc, targetType types.Type) string {
	methodIdent, err := astNodeAtLoc[*ast.Ident](pkg, loc)
	if err != nil || methodIdent == nil {
		return ""
	}

	methodObj, _, _ := types.LookupFieldOrMethod(targetType, true, pkg.Types, methodIdent.Name)
	if methodObj == nil {
		return ""
	}

	return methodObj.Name()
}

func relationSetToSortedSlice(relationSet map[Relation]struct{}) []Relation {
	relations := make(RelationSlice, 0, len(relationSet))
	for r := range relationSet {
		relations = append(relations, r)
	}
	sort.Stable(relations)
	return relations
}

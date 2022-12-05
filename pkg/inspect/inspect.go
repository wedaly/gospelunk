package inspect

import (
	"github.com/wedaly/gospelunk/pkg/file"
)

type Result struct {
	Name      string
	Type      string
	Relations []Relation
}

func Inspect(loc file.Loc, searchDir string, includeRelKinds []RelationKind) (*Result, error) {
	pkg, err := loadGoPackageForFileLoc(loc)
	if err != nil {
		return nil, err
	}

	enrichments := []enrichResultFunc{enrichResultNameAndType}
	for _, relKind := range includeRelKinds {
		if e := enrichmentForRelKind(relKind); e != nil {
			enrichments = append(enrichments, e)
		}
	}

	var result Result
	for _, enrichFunc := range enrichments {
		if err := enrichFunc(&result, pkg, loc, searchDir); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func enrichmentForRelKind(relKind RelationKind) enrichResultFunc {
	switch relKind {
	case RelationKindDef:
		return enrichResultDefRelation
	case RelationKindRef:
		return enrichResultRefRelation
	case RelationKindImpl:
		return enrichResultImplRelation
	case RelationKindIface:
		return enrichResultIfaceRelation
	default:
		return nil
	}
}

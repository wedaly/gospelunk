package inspect

import (
	"github.com/wedaly/gospelunk/pkg/file"
)

type Result struct {
	Name      string
	Type      string
	Relations []Relation
}

func Inspect(loc file.Loc, searchDir string) (*Result, error) {
	pkg, err := loadGoPackageForFileLoc(loc)
	if err != nil {
		return nil, err
	}

	enrichments := []enrichResultFunc{
		enrichResultNameAndType,
		enrichResultDefRelation,
		enrichResultIfaceImplRelation,
	}

	var result Result
	for _, enrichFunc := range enrichments {
		if err := enrichFunc(&result, pkg, loc, searchDir); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

package inspect

import (
	"fmt"

	"github.com/wedaly/gospelunk/pkg/file"
)

type RelationKind string

const (
	// The relation between a usage and its definition.
	RelationKindDef = RelationKind("definition")

	// The relation between an interface and an implementation.
	RelationKindIfaceImpl = RelationKind("interfaceImplementation")
)

var AllRelationKinds []RelationKind
var AllRelationKindStrings []string

func init() {
	AllRelationKinds = []RelationKind{RelationKindDef, RelationKindIfaceImpl}
	for _, r := range AllRelationKinds {
		AllRelationKindStrings = append(AllRelationKindStrings, string(r))
	}
}

func RelationKindFromString(s string) (RelationKind, error) {
	for _, r := range AllRelationKindStrings {
		if s == r {
			return RelationKind(s), nil
		}
	}
	return RelationKind(""), fmt.Errorf("Invalid relation kind %q", s)
}

func RelationKindsFromStrings(relKindStrings []string) ([]RelationKind, error) {
	relKinds := make([]RelationKind, 0, len(relKindStrings))
	for _, s := range relKindStrings {
		rk, err := RelationKindFromString(s)
		if err != nil {
			return nil, err
		}
		relKinds = append(relKinds, rk)
	}
	return relKinds, nil
}

// Relation represents a relationship between an identifier to some other part of the codebase.
type Relation struct {
	file.Loc
	Kind RelationKind
	Pkg  string
	Name string
}

type RelationSlice []Relation

// Len implements sort.Interface#Len
func (rs RelationSlice) Len() int {
	return len(rs)
}

// Less implements sort.Interface#Less
func (rs RelationSlice) Less(i, j int) bool {
	if rs[i].Pkg != rs[j].Pkg {
		return rs[i].Pkg < rs[j].Pkg
	} else if rs[i].Name != rs[j].Name {
		return rs[i].Name < rs[j].Name
	} else if rs[i].Loc.Path != rs[j].Loc.Path {
		return rs[i].Loc.Path < rs[j].Loc.Path
	} else if rs[i].Loc.Line != rs[j].Loc.Line {
		return rs[i].Loc.Line < rs[j].Loc.Line
	} else {
		return rs[i].Loc.Column < rs[j].Loc.Column
	}
}

// Swap implements sort.Interface#Swap
func (rs RelationSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

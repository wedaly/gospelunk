package inspect

import "github.com/wedaly/gospelunk/pkg/file"

type RelationKind string

const (
	// The relation between a usage and its definition.
	RelationKindDef = RelationKind("definition")

	// The relation between an interface and an implementation.
	RelationKindIfaceImpl = RelationKind("interfaceImplementation")
)

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
	}
	return rs[i].Name < rs[j].Name
}

// Swap implements sort.Interface#Swap
func (rs RelationSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

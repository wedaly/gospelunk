package inspect

import "github.com/wedaly/gospelunk/pkg/file"

type RelationKind string

const (
	// The relation between a usage and its definition.
	RelationKindDef = RelationKind("definition")
)

// Relation represents a relationship between an identifier to some other part of the codebase.
type Relation struct {
	file.Loc
	Kind RelationKind
	Pkg  string
	Name string
}

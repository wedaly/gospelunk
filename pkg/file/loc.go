package file

import "fmt"

// Loc specifies a location in a file.
type Loc struct {
	Path   string
	Line   int
	Column int
}

func (loc Loc) String() string {
	return fmt.Sprintf("%s:%d:%d", loc.Path, loc.Line, loc.Column)
}

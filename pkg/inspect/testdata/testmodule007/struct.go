package testmodule007

import "github.com/wedaly/gospelunk/pkg/inspect/testdata/testmodule007/subpkg"

type StructWithEmbedding struct {
	*subpkg.EmbeddedStruct
}

package pkgmeta

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// PkgHash is a hash of file info for the package's go.mod and all source files.
type PkgHash [sha256.Size]byte

func (h PkgHash) String() string {
	return base64.StdEncoding.EncodeToString(h[:])
}

func (h PkgHash) Empty() bool {
	for _, b := range h {
		if b > 0 {
			return false
		}
	}
	return true
}

// HashFileInfo returns a hash of file info for the package's go.mod and all source files.
func HashFileInfo(pkg Package) (PkgHash, error) {
	h := sha256.New()
	for _, p := range pathsToHash(pkg) {
		fileInfo, err := os.Stat(p)
		if err != nil {
			return PkgHash{}, errors.Wrapf(err, "os.Stat")
		}

		h.Write([]byte(fileInfo.Name()))
		binary.Write(h, binary.LittleEndian, fileInfo.Size())
		binary.Write(h, binary.LittleEndian, fileInfo.ModTime().UnixMilli())
	}

	var hash PkgHash
	copy(hash[:], h.Sum(nil)) // byte slice to fixed-size byte array
	return hash, nil
}

func pathsToHash(pkg Package) []string {
	var paths []string
	if pkg.Module != nil && pkg.Module.GoMod != "" {
		paths = append(paths, pkg.Module.GoMod)
	}

	for _, filename := range pkg.GoFiles {
		paths = append(paths, filepath.Join(pkg.Dir, filename))
	}

	return paths
}

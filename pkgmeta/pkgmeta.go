package pkgmeta

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Package represents metadata about a package.
type Package struct {
	Name       string   // package name
	Dir        string   // directory containing package sources
	ImportPath string   // import path of the package
	GoFiles    []string // filenames of .go source files
}

// Lookup retrieves package metadata from the `go list` command.
// pkgPatterns can include an exact package name ("fmt", "github.com/foo/bar")
// or a relative name (".", "./...").  See `go help packages` for details.
// Packages are resolved relative to the Go module from the current working directory;
// this will fail if the current working directory is not part of a Go module.
func Lookup(pkgPatterns []string) ([]Package, error) {
	args := []string{"-e", "-json"}
	args = append(args, pkgPatterns...)
	data, err := execGoListCmd(args...)
	if err != nil {
		return nil, err
	}
	return parseGoListJsonOutput(data)
}

// ListDirs retrieves a list of package directories matching pkgPatterns.
// See `go help packages` for details about the format of pkgPatterns.
func ListDirs(pkgPatterns []string) ([]string, error) {
	args := []string{"-e", "-find", "-f", "{{ .Dir }}"}
	args = append(args, pkgPatterns...)
	data, err := execGoListCmd(args...)
	if err != nil {
		return nil, err
	}
	dirs := parseGoListDirsOutput(data)
	return dirs, nil
}

func execGoListCmd(args ...string) ([]byte, error) {
	var out bytes.Buffer
	args = append([]string{"list"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "Command.Run")
	}
	return out.Bytes(), nil
}

func parseGoListJsonOutput(data []byte) ([]Package, error) {
	var result []Package
	decoder := json.NewDecoder(bytes.NewReader(data))
	for {
		var pkgMeta Package
		if err := decoder.Decode(&pkgMeta); err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "json.Decoder")
		}
		result = append(result, pkgMeta)
	}
	return result, nil
}

func parseGoListDirsOutput(data []byte) []string {
	var lines []string
	var sb strings.Builder
	for _, b := range data {
		if b == '\n' {
			if sb.Len() > 0 {
				lines = append(lines, sb.String())
				sb.Reset()
			}
		} else {
			sb.WriteByte(b)
		}
	}

	if sb.Len() > 0 {
		lines = append(lines, sb.String())
	}
	return lines
}

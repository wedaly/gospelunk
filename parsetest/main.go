package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"github.com/wedaly/gospelunk/pkgmeta"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("USAGE: %s DIR\n", os.Args[0])
		os.Exit(1)
	}

	root := os.Args[1]

	var fileCount, defCount int
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".go" {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			godefs, err := pkgmeta.LoadGoDefs(f)
			if err != nil {
				return err
			}

			fileCount++
			defCount += len(godefs)
		}
		return nil
	})

	fmt.Printf("Parsed %d files with %d defs", fileCount, defCount)
}

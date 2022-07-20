package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"github.com/wedaly/gospelunk/pkgmeta"
)

const NUM_WORKERS = 16

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("USAGE: %s DIR\n", os.Args[0])
		os.Exit(1)
	}

	root := os.Args[1]

	var mu sync.Mutex
	var fileCount, defCount int
	var wg sync.WaitGroup
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".go" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()

				f, err := os.Open(path)
				if err != nil {
					return err
				}
				defer f.Close()

				godefs, err := pkgmeta.LoadGoDefs(f)
				if err != nil {
					return err
				}

				mu.Lock()
				fileCount++
				defCount += len(godefs)
				mu.Unlock()

			}(path)
		}
		return nil
	})
	wg.Wait()

	fmt.Printf("Parsed %d files with %d defs", fileCount, defCount)
}

package output

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func Template(tmplString string) (*template.Template, error) {
	tmpl, err := template.New("template").
		Funcs(template.FuncMap{"RelPath": RelPathTplFunc()}).
		Parse(tmplString)

	if err != nil {
		return nil, fmt.Errorf("template.Parse: %w", err)
	}

	return tmpl, nil
}

func RelPathTplFunc() func(string) string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	return func(path string) string {
		relpath, err := filepath.Rel(cwd, path)
		if err != nil {
			relpath = path
		}
		return relpath
	}
}

package output

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

func Template(tmplString string) (*template.Template, error) {
	tmpl, err := template.New("template").
		Funcs(template.FuncMap{"RelPath": RelPathTplFunc()}).
		Parse(tmplString)

	if err != nil {
		return nil, errors.Wrapf(err, "template.Parse")
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

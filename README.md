go-symbol-search
================

CLI tool to find symbols in Go projects quickly.

Key Features
------------

-	Single command to search your project, the Go standard library, and Go module dependencies.
-	Fast lookup and reindexing.
-	Easy integration with terminal-based workflows (CLI, fzf, editors, ...).

Installation
------------

```
go get github.com/wedaly/go-symbol-search@latest
```

Usage
-----

Search for definitions:

```
gss def Foobar
```

The first search might take longer because it needs to index the codebase; subsequent searches will be faster. If you want, you can create or update the index ahead-of-time like this:

```
gss index        # reindex anything that has changed
gss index -clean # rebuild the index from scratch
```

By default, the tool searches the package in the current working directory and its direct dependencies (imports). You can override this by passing a packages list:
```
gss def Foobar ./foopkg
gss def Foobar ./...
```
For details about how to specify the packages list, see `go help packages`.

The tool searches both exported and unexported symbols for packages in the package list. For dependencies (imported packages), it searches only exported symbols.

The search term is a glob pattern that matches the fully qualified symbol name, so you can also do this:

```
gss def Foo*
gss def foopkg.Foo*
gss def foopkg.Foo*Bar
```

The default output format looks like this:

```
foo/bar.go:4:25 Foobar func() bool
```

You can override the output format by specifying a Go template:

```
gss def Foo -o '{{ .Path }}:{{ .Line }}:{{ .Column }} {{ .Type }} {{ .Name }}'
```

These template variables are defined:

| Variable | Meaning                                                             |
|----------|---------------------------------------------------------------------|
| Path     | Absolute path to the file containing the symbol.                    |
| Line     | Line number of the symbol (1-indexed).                              |
| Column   | Column of the first character of the symbol (1-indexed, byte count) |
| Type     | The Go type of the symbol.                                          |
| Name     | The name of the symbol.                                             |

In addition to Go's [predefined global template functions](https://pkg.go.dev/text/template#hdr-Functions), these functions are available:

| Function | Args          | Meaning                                                                           |
|----------|---------------|-----------------------------------------------------------------------------------|
| RelPath  | path (string) | Transform an absolute path to a relative path from the current working directory. |
| Filename | path (string) | Returns the filename of the last component in the path.                           |

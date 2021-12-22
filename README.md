gospelunk
================

CLI tool to quickly find definitions in Go projects and dependencies.

Key Features
------------

-	Single command to search your project, the Go standard library, and Go module dependencies.
-	Fast lookup and reindexing.
-	Easy integration with terminal-based workflows (CLI, fzf, editors, ...).

Installation
------------

```
go install github.com/wedaly/gospelunk@latest
```

Usage
-----

Create or update the index:

```
gospelunk index              # package in current directory
gospelunk index -pkgs=./...  # packages in current directory and all subdirectories (recursive)
gospelunk index -imports     # index packages imported by the package in the current directory
```

Search for definitions:

```
gospelunk find Foobar             # find "Foobar" in the current package
gospelunk find -pkg=./... Foobar  # include packages in subdirectories
gospelunk find -imports Foobar    # search packages imported by the current package
```

The results include all symbols that contain the search query.

By default, the `gss find` commands will print a warning if the index is out-of-date.
You can ask `gss find` to automatically reindex like this:
```
gospelunk find -reindex Foobar
```

The default output format looks like this:

```
foo/bar.go:4:25 github.com/example/foo.Foobar
```

You can override the output format by specifying a Go template:

```
gospelunk find -o '{{ .Path }}:{{ .Line }}:{{ .Column }} {{ .Type }} {{ .Name }}' Foo
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

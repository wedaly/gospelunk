gospelunk
=========

Go spelunking! CLI tool to quickly find definitions in Go projects and dependencies.

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
gospelunk index .          # index the package in current directory
gospelunk index ./...      # index packages in current directory and all subdirectories (recursive)
gospelunk index -imports . # index packages imported by the package in the current directory
```

Search for definitions:

```
gospelunk find Foobar           # find "Foobar" in the current package
gospelunk find Foobar ./...     # include packages in subdirectories
gospelunk find -imports Foobar  # include definitions from imports
```

The results include all symbols that contain the search query. The default output format looks like this:

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

Building from Source
--------------------

1.	[Install protoc](https://developers.google.com/protocol-buffers/docs/downloads) by downloading the package and following instructions in the README.
2.	Install the Go protocol buffers plugin:

	```
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	```

3.	Run `make`

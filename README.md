gospelunk
=========

Go spelunking! CLI tool to quickly find definitions in Go projects and dependencies.

Key Features
------------

-	Single command to search your project, the Go standard library, and Go module dependencies.
-	Fast lookup and reindexing.
-	Easy integration with terminal-based workflows (CLI, fzf, editors, ...).

Project Status: Alpha
---------------------

-	Things might change
-	Things might be broken

Installation
------------

```
go install github.com/wedaly/gospelunk@latest
```

Usage
-----

Create or update the index:

```
gospelunk index .      # index the package in current directory
gospelunk index ./...  # index packages in current directory and all subdirectories (recursive)
gospelunk index -i .   # index packages imported by the package in the current directory
```

The index command will skip packages if they haven't changed since they were last indexed.

Search for definitions:

```
gospelunk find Foobar       # find "Foobar" in the current package
gospelunk find Foobar ./... # include packages in subdirectories
gospelunk find -i Foobar    # include definitions from imports
```

The query is a [regex](https://github.com/google/re2/wiki/Syntax) that matches the definition name.

The results include all symbols that contain the search query. The default output format looks like this:

```
foo/bar.go:123 func Foobar
```

You can override the output format by specifying a Go template:

```
gospelunk find -f '{{ .Path }}:{{ .LineNum }} {{ .Kind }} {{ .Name }}' Foo
```

These template variables are defined:

| Variable | Meaning                                                      |
|----------|--------------------------------------------------------------|
| Path     | Absolute path to the file containing the symbol.             |
| LineNum  | Line number of the symbol (1-indexed).                       |
| Name     | The name of the definition.                                  |
| Kind     | The kind of the definition ("value", "func", "struct", etc.) |

In addition to Go's [predefined global template functions](https://pkg.go.dev/text/template#hdr-Functions), these functions are available:

| Function | Args          | Meaning                                                                           |
|----------|---------------|-----------------------------------------------------------------------------------|
| RelPath  | path (string) | Transform an absolute path to a relative path from the current working directory. |
| BasePath | path (string) | Returns the last component in the path.                                           |

For example:

```
gospelunk find -f '{{ .Path | BasePath }} {{ .Name }}' .
```

Index Database
--------------

gospelunk stores the search index in a [BoltDB database](https://github.com/etcd-io/bbolt) at `$XDG_DATA_HOME/gospelunk/index.db`. The default locations are:

-	Linux: `~/.local/share/gospelunk`
-	macOS: `~/Library/Application Support/gospelunk`

Integrations
------------

### Bash function, fzf, less

If you have [fzf](https://github.com/junegunn/fzf) installed, you can add this to your .bashrc:

```
gofind () {
	less $(gospelunk find -i -f "+{{.LineNum}} {{.Path|RelPath}} {{.Kind}} {{.Name}}" $@ | fzf | cut -d " " -f 1-2)
}
```

After `source .bashrc` you can search packages and imports in the current working directory like this:

```
gofind Foobar
```

### Aretext menu commands

Run `aretext -editconfig`, then add this rule:

```
- name: gospelunk-commands
  pattern: "**/*.go"
  config:
    menuCommands:
    - name: gospelunk index
      shellCmd: gospelunk index -i ./... | less
    - name: gospelunk find
      shellCmd: gospelunk find -i -f "{{.Path|RelPath}}:{{.LineNum}}:{{.Kind}} {{.Name}}" "^(.+\.)?${WORD}$" $FILEPATH
      mode: fileLocations
```

For more details, see ["Custom Menu Commands" in the aretext docs](https://aretext.org/docs/custom-menu-commands/).

### Git post-checkout hook

Add a [post-checkout git hook](https://git-scm.com/docs/githooks#_post_checkout) to reindex when changing branches (this will replace any existing post-checkout hooks!):

```
cat << EOF > .git/hooks/post-checkout
#!/usr/bin/env sh
echo "Indexing Go packages (gospelunk)"
gospelunk index -q -i ./...
EOF
chmod +x .git/hooks/post-checkout
```

Building from Source
--------------------

1.	[Install protoc](https://developers.google.com/protocol-buffers/docs/downloads) by downloading the package and following instructions in the README.
2.	Install the Go protocol buffers plugin:

	```
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	```

3.	Run `make`

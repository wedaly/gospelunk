gospelunk
=========

Go spelunking! CLI tool to quickly find things in Go projects.

Project Status: Alpha
---------------------

-	Things might change
-	Things might break

Installation
------------

Once you have [installed Go](https://go.dev/doc/install), run this:

```
go install github.com/wedaly/gospelunk@latest
```

Commands
--------

### List

To list definitions in Go packages:

```
gospelunk list ./...
```

-	You can specify packages using the same format as other go commands. See `go help packages` for details.
-	You can use the `--template` parameter to customize the Go template used to render the output.
-	Use `--include-private` to include non-exported definitions.
-	Use `--include-tests` to include definitions from "_test.go" files.

### Inspect

To lookup type information and the definition for an identifier in a Go file:

```
gospelunk inspect -f <FILE> -l <LINE> -c <COLUMN>
```

-	Line and column numbers are 1-indexed, and the column unit is bytes.
-	You can use the `--template` parameter to customize the Go template used to render the output.

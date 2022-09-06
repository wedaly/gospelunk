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

To lookup type information and the definition for an identifier in a Go file:

```
gospelunk index -f <FILE> -l <LINE> -c <COLUMN>
```

-	Line and column numbers are 1-indexed, and the column unit is bytes.
-	You can use the `--template` parameter to customize the Go template used to render the output.

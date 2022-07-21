# gospelunk v2 spec

## Overview

Go spelunking! Browser-based tool for exploring large Go projects.

* Quickly find definitions and references in your projects, the Go standard library, and module dependencies.
* Efficient on-the-fly search -- never wait for an index to rebuild.
* Minimalist interface in plain HTML.


## CLI

Start the server:
```
gospelunk --address <host>:<port> <directory>
```

Go modules in `directory` (recursive) will be searched. Possible to specify multiple directories. If the same module appears in multiple locations, gospelunk will print a warning and choose the one with the shortest path. If `directory` is ommitted, default to `$GOPATH/src`.


## URL structure

### GET /

Main entry point to find something in one of your projects.

* search box
* checkbox "include tests"
* help text describing filter syntax (see below)
* submit button with action `POST /search/defs`


### POST /search/defs

POST params:
* `query` (string) specifies the search query.
* `includeTests` (bool) specifies whether to include tests in the results.

The query is parsed into a sequence of "filters", separated by whitespace:

* `in:<pkg>` search only package paths matching `pkg` regex (including substring matches).
* `type:<type>` search only definitions of the specified type (`func`, `struct`, `interface`, `method`, `string`, etc.)
* all other terms are interpreted as regex patterns matching the definition name.

Filters with the same prefix are OR'd together; filters with different prefixes are AND'd. For example, `Foo in:foo in:bar type:string type:func` would search in packages matching "foo" OR "bar" and would include only definitions with type "string" OR "func".

The response is an HTML page containing search results:

* package, linked to `/go/{module}/{pkg}`
* definition type ("func", "struct", etc.)
* definition name, linked to `/go/{module}/{pkg}/defs/{name}`


### GET /go/{module}

Show information about a Go module:

* name
* list of packages, linked to `/go/{module}/{pkg}`


### GET /go/{module}/{pkg}

Show information about a Go package:

* name
* module, linked to `/go/{module}`
* imported packages, linked to `/go/{module}/{pkg}/imports/{module}/{pkg}`
* definitions (grouped by type, public/private, and test/non-test), linked to `/go/{module}/{pkg}/defs/{name}`


### GET /go/{module}/{pkg}/defs/{name}

Show information about a Go definition.

* name
* type (struct, func, etc.)
* package, linked to `/go/{module}/{pkg}`
* file location, linked to `/go/{module}/{pkg}/files/{name}.go?line={line}`
* references, linked to `/go/{module}/{pkg}/defs/{name}/refs`
* docstring
* for structs, list of fields and methods, linked to `/go/{module}/{pkg}/defs/{name}`
* for interfaces, list of method specs


### GET /go/{module}/{pkg}/defs/{name}/refs

Search for references to the definition.

The response is an HTML page containing the search results:

* package, linked to `/go/{module}/{pkg}`
* file name and line number, linked to `/go/{module}/{pkg}/files/{name}.go?line={line}`


### GET /go/{module}/{pkg}/files/{name}.go

* full path to file, linked to `file://{path}`
* copy path to clipboard button
* contents of file, with line numbers

The optional anchor `#l[0-9]+` can be used to specify a line number in the file.

Within the file contents:

* The top-level `package` statement links to `/go/{module}/{pkg}`.

* Each import statement links to `/go/{module}/{pkg}/imports/{module}/{pkg}`.

* Links for top-level definitions and types. For definitions within the package, link to `/go/{module}/{pkg}/defs/{name}`, and for definitions in other packages, link to `/go/{module}/{pkg}/imports/{module}/{pkg}/defs/{name}`. This needs to handle user-defined package names correctly ("import f golang.org/x/foobar")

* Local variables within a function link to the line number anchor in the file where the variable is defined. This needs to handle variable shadowing correctly.


### GET /file/{path}

Lookup the Go module and package containing `path` (which should be absolute).
If the package exists, and the file has a ".go" extension, redirect to `/go/{pkg}/files/{name}.go`
Otherwise, return a 404 not found.


### GET /go/{module}/{pkg}/imports/{module}

Same as `/go/{module}`, but construct links relative to `/go/{module}/{pkg}/imports`.


### GET /go/{module}/{pkg}/imports/{module}/{pkg}

Same as `/go/{module}/{pkg}`, but construct links relative to `/go/{module}/{pkg}/imports`.


### GET /go/{module}/{pkg}/imports/{module}/{pkg}/defs/{name}

Same as `/go/{module}/{pkg}/defs/{name}`, but construct links relative to `/go/{module}/{pkg}/imports`.


### GET /go/{module}/{pkg}/imports/{module}/{pkg}/files/{name}.go

Same as `/go/{module}/{pkg}/files/{name}.go`, but construct links relative to `/go/{module}/{pkg}/imports`.

This is necessary to resolve imports relative to the original module.

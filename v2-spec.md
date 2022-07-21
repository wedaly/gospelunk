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

Go modules in `directory` (recursive) will be searched. Possible to specify multiple directories. If the same module appears in multiple locations, gospelunk will print a warning and choose the one with the shortest path. If `directory` is ommitted, use the current working directory.


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

TODO: lookup a module


### GET /go/{module}/{pkg}

TODO: lookup a package


### GET /go/{module}/{pkg}/defs/{name}

TODO: lookup a def


### GET /go/{module}/{pkg}/defs/{name}/refs

Search for references to the definition.

The response is an HTML page containing the search results:

* package, linked to `/go/{module}/{pkg}`
* file name and line number, linked to `/go/{module}/{pkg}/files/{name}.go?line={line}`


### GET /go/{module}/{pkg}/files/{name}.go

TODO: optional ?line


### GET /file/{path}

Lookup the Go module and package containing `path` (which should be absolute).
If the package exists, and the file has a ".go" extension, redirect to `/go/{pkg}/files/{name}.go`
Otherwise, return a 404 not found.


### GET /go/{module}/{pkg}/imports/{module}

TODO

### GET /go/{module}/{pkg}/imports/{module}/{pkg}

TODO

### GET /go/{module}/{pkg}/imports/{module}/{pkg}/defs/{name}

TODO

### GET /go/{module}/{pkg}/imports/{module}/{pkg}/files/{name}.go

TODO

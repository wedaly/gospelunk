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
* help text describing filter syntax

The search box parses the query into a sequence of "filters", separated by whitespace:

* `in:<pkg>` search only package paths matching `pkg` regex (including substring matches).
* `type:<type>` search only definitions of the specified type (`func`, `struct`, `interface`, `method`, `string`, etc.)
* all other terms are interpreted as regex patterns matching the definition name.

Filters with the same prefix are OR'd together; filters with different prefixes are AND'd. For example, `Foo in:foo in:bar type:string type:func` would search in packages matching "foo" OR "bar" and would include only definitions with type "string" OR "func".



```
GET /
  search interface

POST /search/defs
POST /search/refs
  run a search of all modules in root dirs
  show results page
  slightly different UI for defs/refs (for refs, already know what we're 
  looking for, so no need to show the symbol name in each result)
  ref search isn't module specific (looking for everything referring to that mod+pkg+symbol,
  even if it's using a different module version)

GET /file/{path}
  path should be absolute
  lookup module/pkg for path
  redirect to /go/{module}/{pkg}/files/{name}.go

GET /go/{module}
GET /go/{module}/{pkg}
GET /go/{module}/{pkg}/defs/{name}
GET /go/{module}/{pkg}/files/{name}.go
GET /go/{module}/{pkg}/imports/{module}
GET /go/{module}/{pkg}/imports/{module}/{pkg}
GET /go/{module}/{pkg}/imports/{module}/{pkg}/defs/{name}
GET /go/{module}/{pkg}/imports/{module}/{pkg}/files/{name}.go
```

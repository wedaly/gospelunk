Work Breakdown
==============

Milestone 1: Index and search project files
-------------------------------------------

[ ] initialize a SQLite database in XDG_DATA
[ ] parse a single Go file and write symbols to SQLite
[ ] search exact substring match in the index
[ ] walk working directory recursively and index all Go files
[ ] search glob patterns in the index
[ ] skip reindexing files that haven't changed, based on mtime and md5 checksum
[ ] determine fully qualified name for each symbol based on package and go.mod
[ ] handle import aliases when determining fully qualified names

Milestone 2: Index and search dependencies
------------------------------------------

[ ] index Go stdlib
[ ] index Go module dependencies
[ ] handle Go module rewrites correctly

Milestone 3: Quality of life improvements
-----------------------------------------

[ ] display progress bar in tty
[ ] optional template parameter for output format
[ ] optionally output the type of the symbol (function, int, etc.)
[ ] optionally filter by type of symbol
[ ] command to index without searching
[ ] command to clear the index

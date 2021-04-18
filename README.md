# peanut

[![BSD3](https://img.shields.io/badge/license-BSD3-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/travis/jimsmart/peanut/master.svg)](https://travis-ci.org/jimsmart/peanut)
[![codecov](https://codecov.io/gh/jimsmart/peanut/branch/master/graph/badge.svg)](https://codecov.io/gh/jimsmart/peanut)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimsmart/peanut)](https://goreportcard.com/report/github.com/jimsmart/peanut)
[![Used By](https://img.shields.io/sourcegraph/rrc/github.com/jimsmart/peanut.svg)](https://sourcegraph.com/github.com/jimsmart/peanut)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jimsmart/peanut)

peanut is a [Go](https://golang.org/) package to write tagged data structs in a variety of formats.

Its primary purpose is to provide a single consistent interface
for easy, ceremony-free persistence of record-based struct data.

Each distinct struct type is written to an individual file (or table), each named according to the name of the struct. Field/column names in each file/table are derived from struct tags. All writers use the same tags.

Currently supported formats are CSV, TSV, Excel (.xlsx), JSON Lines (JSONL), and SQLite.
Additional writers are also provided to assist with testing and debugging.
Mutiple writers can be combined using MultiWriter.

## About

TODO

## Quickstart

### Installation

Get the package:

```bash
go get github.com/jimsmart/peanut
```

Use the package within your code:

```go
import "github.com/jimsmart/peanut"
```

### Example Code

See GoDocs.

## Documentation

GoDocs [https://godoc.org/github.com/jimsmart/peanut](https://godoc.org/github.com/jimsmart/peanut)

## Testing

To run the tests execute `go test` inside the project folder.

For a full coverage report, try:

```bash
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
```

## License

Package peanut is copyright 2020-2021 by Jim Smart and released under the [BSD 3-Clause License](LICENSE.md).

## History

- v1.0.0: Repository made public (April 2021).

# peanut

[![BSD3](https://img.shields.io/badge/license-BSD3-blue.svg)](LICENSE.md)
[![Build Status](https://github.com/jimsmart/peanut/actions/workflows/main.yml/badge.svg)](https://github.com/jimsmart/peanut/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/jimsmart/peanut/branch/master/graph/badge.svg)](https://codecov.io/gh/jimsmart/peanut)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimsmart/peanut?cache-buster)](https://goreportcard.com/report/github.com/jimsmart/peanut)
[![Used By](https://img.shields.io/sourcegraph/rrc/github.com/jimsmart/peanut.svg)](https://sourcegraph.com/github.com/jimsmart/peanut)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jimsmart/peanut)

peanut is a [Go](https://golang.org/) package to write tagged data structs to disk in a variety of formats.

Its primary purpose is to provide a single consistent interface
for easy, ceremony-free persistence of record-based struct data.

Each distinct struct type is written to an individual file (or table), each named according to the name of the struct. Field/column names in each file/table are derived from struct tags. All writers use the same tags.

Currently supported formats are CSV, TSV, Excel (.xlsx), JSON Lines (JSONL), and SQLite.
Additional writers are also provided to assist with testing and debugging.
Mutiple writers can be combined using MultiWriter.

All writers perform atomic file operations, writing data to a temporary location and moving
it to the final output location when Close is called.

## About

When building an app or tool that needs to output data consisting of
multiple different record types to disk, perhaps with requirements that
change over time (whether during development or after initial deployment),
perhaps requiring multiple output formats (during development/testing,
or as final output) — is where peanut might be 'the right tool for the job'.

Ideal for use as an output solution for, e.g. data conversion tools,
part of an ETL pipeline, data-acquistion or extraction tools/apps, web-scrapers,
structured logging, persistence of captured data/metadata/events,
job reporting, etc.
Whether building an ad-hoc tool as a quick hack, or as part of a bigger,
more serious project.

peanut initially evolved as part of a larger closed-source project,
is tried and tested, and production-ready.

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

### API

All peanut writers implement this interface:

```go
type Writer interface {
    Write(r interface{}) error
    Close() error
    Cancel() error
}
```

### Usage

1. Tag some structs.
2. Initialise a `peanut.Writer` to use.
3. Collect and assign data into tagged structs.
4. Use `Write()` to write records, repeating until done.
5. Call `Close()` or `Cancel()` to finish.

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

- v1.0.2 (2021-04-19) Fixup handling of uints.
- v1.0.1 (2021-04-19) Repository made public.

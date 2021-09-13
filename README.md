# scoop_search

A _fast_ way to search [scoop](https://scoop.sh) manifests.

## Features

- Far faster searching than scoop's built in search functionality
- That's about it, but it's _considerably_ faster!

## Installation

With go installed:

```shell
go get -uv github.com/tehbilly/scoop_search
```

## Usage

```shell
Usage:
  scoop_search [flags]

Flags:
  -h, --help              help for scoop-search
  -m, --in-mem            Use in-memory search index (default true)
  -j, --json              Output JSON instead of table
  -n, --num-results int   Number of results to show (default 10)
  -v, --verbose           Output operational (debug) messages
```

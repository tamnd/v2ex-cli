# v2ex

Browse V2EX topics, nodes and members

`v2ex` is a single pure-Go binary. It reads the public V2EX JSON API over
plain HTTPS, shapes the responses into clean records, and pipes into the rest
of your tools. No API key, nothing to run alongside it.

## Install

```bash
go install github.com/tamnd/v2ex-cli/cmd/v2ex@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/v2ex-cli/releases), or run
the container image:

```bash
docker run --rm ghcr.io/tamnd/v2ex:latest --help
```

## Usage

```bash
# hot and latest topic lists
v2ex hot
v2ex latest
v2ex hot -n 10 -o json

# single topic by ID
v2ex topic 1000

# node by name
v2ex node go

# member profile
v2ex member Livid

# replies for a topic
v2ex replies 1000
```

Output defaults to a table on a TTY and JSONL when piped:

```bash
v2ex hot -o json
v2ex hot -o csv
v2ex hot --fields rank,title,url
```

## Development

```
cmd/v2ex/   thin main, wires cli.Root into fang
cli/        cobra command tree
v2ex/       HTTP client and data models
pkg/render/ output rendering (table/json/jsonl/csv/tsv)
docs/       tago documentation site
```

```bash
make build      # ./bin/v2ex
make test       # go test ./...
make vet        # go vet ./...
```

## Releasing

Push a version tag and GitHub Actions runs GoReleaser, which builds the
archives, Linux packages, the multi-arch GHCR image, checksums, SBOMs, and a
cosign signature:

```bash
git tag v0.1.0
git push --tags
```

The Homebrew and Scoop steps self-disable until their tokens exist, so the first
release works with no extra secrets.

## License

Apache-2.0. See [LICENSE](LICENSE).

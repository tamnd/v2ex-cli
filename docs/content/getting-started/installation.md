---
title: "Installation"
description: "Install v2ex from a release, with go install, or from source."
weight: 20
---

## Prebuilt binaries

Every [release](https://github.com/tamnd/v2ex-cli/releases) carries archives for Linux, macOS,
and Windows on amd64 and arm64, plus deb, rpm, and apk packages for Linux.
Download, unpack, put `v2ex` on your `PATH`, done. The `checksums.txt`
on each release is signed with keyless [cosign](https://docs.sigstore.dev/) if
you want to verify before running.

## With Go

```bash
go install github.com/tamnd/v2ex-cli/cmd/v2ex@latest
```

That puts `v2ex` in `$(go env GOPATH)/bin`, which is `~/go/bin` unless
you moved it. Make sure that directory is on your `PATH`.

## From source

```bash
git clone https://github.com/tamnd/v2ex-cli
cd v2ex-cli
make build        # produces ./bin/v2ex
./bin/v2ex version
```

## Container image

```bash
docker run --rm ghcr.io/tamnd/v2ex:latest --help
```

## Checking the install

```bash
v2ex version
```

prints the version and exits.

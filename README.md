# publisher-oci

[![Latest Release](https://img.shields.io/github/v/release/SemRels/publisher-oci?label=version&color=blue)](https://github.com/SemRels/publisher-oci/releases/latest)
[![CI](https://github.com/SemRels/publisher-oci/actions/workflows/ci.yml/badge.svg)](https://github.com/SemRels/publisher-oci/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/SemRels/publisher-oci)](LICENSE)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/SemRels/publisher-oci/badge)](https://scorecard.dev/viewer/?uri=github.com/SemRels/publisher-oci)
[![REUSE status](https://api.reuse.software/badge/github.com/SemRels/publisher-oci)](https://api.reuse.software/info/github.com/SemRels/publisher-oci)

OCI publisher plugin for semrel. It publishes release artifacts to OCI registries using `oras push`.

## Repository Layout

```text
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
```

## Installation

Published binaries are distributed through releases and synchronized to `registry.semrel.io`.

## Development

```bash
go build ./cmd/plugin
go test ./...
```

## Configuration

Configure in `.semrel.yaml`:

```yaml
plugins:
  - uses: "@semrel/publisher-oci"
    phase: release
    args:
      ref: ghcr.io/semrels/myapp:{version}
      artifacts: dist/myapp-linux-amd64,dist/myapp-linux-arm64
```

## Options

Every option can be set as an entry under `args` in `.semrel.yaml`. semrel passes each entry to the plugin as the environment variable shown in the table, so you can also set that variable directly (for example in your CI environment) instead of using `args`.

| Argument | Environment variable | Required | Description | Default |
| --- | --- | --- | --- | --- |
| `ref` | `SEMREL_PLUGIN_REF` | Required | Destination reference. Supports the `{version}` token. | None |
| `artifacts` | `SEMREL_PLUGIN_ARTIFACTS` | Optional | Comma-separated list of artifacts to push. | None |
| `artifacts_json` | `SEMREL_PLUGIN_ARTIFACTS_JSON` | Optional | JSON array of artifacts, used as a fallback for `artifacts`. | None |
| `artifact` | `SEMREL_PLUGIN_ARTIFACT` | Optional | Single artifact, used as a fallback. | None |

Release context:

| Variable | Description |
| --- | --- |
| `SEMREL_VERSION` / `SEMREL_NEXT_VERSION` | Required release version. |
| `SEMREL_DRY_RUN` | Whether semrel is running in dry-run mode. |

Dependencies:

- `oras` must be installed and available on `PATH`.

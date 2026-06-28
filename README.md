# publisher-oci

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
	- uses: publisher-oci
		args:
			ref: ghcr.io/semrels/myapp:{version}
			artifacts: dist/myapp-linux-amd64,dist/myapp-linux-arm64
```

Runtime inputs:

- `SEMREL_VERSION` / `SEMREL_NEXT_VERSION` (required)
- `SEMREL_DRY_RUN`
- `SEMREL_PLUGIN_REF` (required; supports `{version}` token)
- `SEMREL_PLUGIN_ARTIFACTS` (CSV)
- `SEMREL_PLUGIN_ARTIFACTS_JSON` (JSON array fallback)
- `SEMREL_PLUGIN_ARTIFACT` (single artifact fallback)

Dependencies:

- `oras` must be installed and available on `PATH`.

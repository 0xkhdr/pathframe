# Contributing to Pathframe

## Before editing

Read [PHILOSOPHY.md](PHILOSOPHY.md) and [ARCHITECTURE.md](ARCHITECTURE.md) before changing product behavior or architecture. Coding agents must also follow [AGENTS.md](AGENTS.md).

Each task must name its dependency, owned files or component, governing invariant, focused verification, and observable acceptance criteria. Keep changes small enough to implement without loading the entire repository.

## Repository map

| Path | Responsibility |
| --- | --- |
| `cmd/pathframe` | Executable entry point |
| `internal/adapters/{cli,mcp}` | CLI and typed MCP projections |
| `internal/app` | Canonical application operations |
| `internal/{workflow,artifacts,context,delegation,verification,recovery}` | Deterministic domain behavior |
| `internal/integrations/{codex,claude}` | Generated host assets and capability declarations |
| `internal/store` | Journal and projection persistence |
| `schemas` | Persisted and cross-process JSON contracts |
| `tests/journey` | End-to-end product journeys |

Dependencies point inward: adapters call `internal/app`; application operations coordinate domains; domains do not import adapters or host SDKs. Put a workflow rule in its owning domain once and expose it through both adapters.

## Development loop

Use Go 1.26. Run the focused package test while editing, then the full affected suite:

```sh
go test ./internal/<package>
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
GOOS=linux GOARCH=amd64 go build ./cmd/pathframe
```

Run the relevant test in `tests/journey` when behavior crosses an adapter, persistence, delegation, verification, or recovery boundary. Use `go test -race ./...` for changes to process execution or shared state. The complete release matrix lives in [docs/release.md](docs/release.md).

## Engineering rules

- Use Go 1.26. Linux amd64 is the only initially supported production platform.
- Prefer the standard library. The standard `flag` package is the fixed CLI choice.
- Add no dependency without a concrete need and operational-cost explanation. The official MCP Go SDK stays pinned inside the stdio MCP adapter.
- Keep application operations canonical. CLI, MCP, host guidance, and hooks translate; they do not decide workflow state.
- Preserve explicit approval, `execution_policy`, sequential delegation, timeout-bound argv verification, and journal-based reconstruction.
- Never make Git history or conversation history authoritative.
- Never use deletion of `.pathframe/` as normal recovery.
- Do not add Specd compatibility or copy a legacy subsystem without an approved, evidence-backed decision record.

## Documentation and decisions

Update the smallest owning page when a public contract or supported behavior changes: `README.md` for orientation, `docs/CLI.md` for commands, `docs/ARTIFACTS.md` for authored formats, and `docs/INTEGRATIONS.md` for host setup. Accepted implementation decisions are fixed constraints, not recurring questions. Follow [the decision process](docs/decisions/README.md) only when new repository evidence forces a deviation or a new architecture/scope decision blocks work.

## Completion gate

A change is complete only when focused and affected tests pass, its relevant journey is demonstrated, documentation matches behavior, and unresolved risks are recorded.

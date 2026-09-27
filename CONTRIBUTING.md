# Contributing to Pathframe

## Start with the current stage

Read [the roadmap](docs/planning/IMPLEMENTATION-ROADMAP.md), its detailed stage file, [PHILOSOPHY.md](PHILOSOPHY.md), and [ARCHITECTURE.md](ARCHITECTURE.md). Do not begin a later stage while the current gate has unresolved failures.

Each task must name its dependency, owned files or component, governing invariant, focused verification, and observable acceptance criteria. Keep changes small enough to implement without loading the entire repository.

## Engineering rules

- Use Go 1.26. Linux amd64 is the only initially supported production platform.
- Prefer the standard library. The standard `flag` package is the fixed CLI choice.
- Add no dependency without a concrete need and operational-cost explanation. The official MCP Go SDK is pinned only when Stage 3 starts and stays inside the stdio MCP adapter.
- Keep application operations canonical. CLI, MCP, host guidance, and hooks translate; they do not decide workflow state.
- Preserve explicit approval, `execution_policy`, sequential delegation, timeout-bound argv verification, and journal-based reconstruction.
- Never make Git history or conversation history authoritative.
- Never use deletion of `.pathframe/` as normal recovery.
- Do not add Specd compatibility or copy a legacy subsystem without an approved, evidence-backed decision record.

## Verification

Stage 0 will establish these standard checks:

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
GOOS=linux GOARCH=amd64 go build ./cmd/pathframe
```

Run the focused package tests first, then the full affected suite and relevant journey. Record unresolved risks before requesting the human stop/go decision.

## Documentation and decisions

Update documentation in the same task when a public contract or supported behavior changes. Accepted implementation decisions are fixed constraints, not recurring questions. Follow [the decision process](docs/decisions/README.md) only when new repository evidence forces a deviation or a new architecture/scope decision blocks work.

## Stage gate

A stage exits only when implementation is complete, focused and affected tests pass, its journey is demonstrated, documentation matches behavior, risks are recorded, and a human gives the stop/go decision.

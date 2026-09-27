# Contributing to Pathframe

## Start with the product boundaries

Read [PHILOSOPHY.md](PHILOSOPHY.md) and [ARCHITECTURE.md](ARCHITECTURE.md) before changing product behavior or architecture.

Each task must name its dependency, owned files or component, governing invariant, focused verification, and observable acceptance criteria. Keep changes small enough to implement without loading the entire repository.

## Engineering rules

- Use Go 1.26. Linux amd64 is the only initially supported production platform.
- Prefer the standard library. The standard `flag` package is the fixed CLI choice.
- Add no dependency without a concrete need and operational-cost explanation. The official MCP Go SDK stays pinned inside the stdio MCP adapter.
- Keep application operations canonical. CLI, MCP, host guidance, and hooks translate; they do not decide workflow state.
- Preserve explicit approval, `execution_policy`, sequential delegation, timeout-bound argv verification, and journal-based reconstruction.
- Never make Git history or conversation history authoritative.
- Never use deletion of `.pathframe/` as normal recovery.
- Do not add Specd compatibility or copy a legacy subsystem without an approved, evidence-backed decision record.

## Verification

Run these standard checks:

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
GOOS=linux GOARCH=amd64 go build ./cmd/pathframe
```

Run the focused package tests first, then the full affected suite and relevant journey. Record unresolved risks before requesting the human stop/go decision.

## Documentation and decisions

Update documentation in the same task when a public contract or supported behavior changes. Accepted implementation decisions are fixed constraints, not recurring questions. Follow [the decision process](docs/decisions/README.md) only when new repository evidence forces a deviation or a new architecture/scope decision blocks work.

## Completion gate

A change is complete only when focused and affected tests pass, its relevant journey is demonstrated, documentation matches behavior, and unresolved risks are recorded.

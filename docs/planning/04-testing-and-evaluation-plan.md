# Testing and Evaluation Plan

## Test layers

- Unit: pure transitions, validators, reference resolution, replay, containment, budgets.
- Golden: text/JSON orientation, templates, packets, diagnostics, generated host files.
- Integration: filesystem persistence, stdio MCP, process execution, Git observation.
- Journey: compiled binary in temporary repositories; no package internals.
- Host integration: supported Codex/Claude configuration and native sequential worker contract.
- Evaluation: activation classification and cold-agent template completion.

Fixtures live under `testdata/`; black-box journeys under `tests/journey/`. Tests never depend on user home configuration or network. SDK pinning is the only planned non-stdlib production dependency.

## Required matrices

- Every change phase and task state: legal transitions, illegal transition response, one executable recovery.
- Modes: Quick, Standard, High-risk; missing fields/references/questions; approval invalidation.
- Integrations: must-use, offer, must-not-use; fresh-session resume; broken versions/configuration.
- Delegation: success, missing role/capability, oversized required context, conflicting/abandoned lease, worker statuses, Brain bypass.
- Verification: pass, nonzero, timeout, interruption, output truncation, path escape, post-verification content change.
- Doctor: malformed artifact, corrupt projection, incomplete journal tail, stale lease, ambiguous changes, unavailable host.

## Standard commands

After Stage 0:

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go build ./cmd/pathframe
```

Stage-specific commands add focused package tests and named journey tests. Stage 9 adds `go test -race ./...`, cross-build/CI matrices, security fixtures, performance benchmarks, and clean-machine installation tests.

## Journey evidence

Each stage records its command, fixture, expected observable result, and failure recovery in its stage document. A gate needs focused tests, full affected suite, one black-box journey, aligned docs, recorded risks, and human stop/go.

## Planning consistency check

All product capabilities map to Stages 0–9. Each stage has tests and recovery. Dependencies flow forward: operations before adapters, Codex before Claude, both before delegation, packets before workers, sequential completion before any parallel consideration. Quick mode requires only intent/tasks. No recovery requires `.pathframe/` deletion or Git rewriting.


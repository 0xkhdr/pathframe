# Pathframe Implementation Roadmap

## Current state

Stages 0 through 8 are approved. Stage 9 is implemented and awaiting its human gate: Linux amd64 install/support evidence, security and performance suites, portability cross-builds, onboarding, release guidance, and explicit limitations are present.

## Architecture

One canonical application layer coordinates deterministic workflow, artifacts, context, delegation, verification, and recovery domains. CLI and MCP are thin adapters over those operations. Filesystem, process, optional Git observation, MCP SDK, and host integrations remain edge adapters. Authored artifacts express intent; append-only transitions and bounded runs preserve facts; `state.json` is rebuildable. Brain owns reasoning and acceptance; Pinky executes only explicitly delegated tasks.

## Dependency order and status

| Stage | Capability | Depends | Status |
| --- | --- | --- | --- |
| [0](stages/00-foundation.md) | Foundation and constitution | none | approved |
| [1](stages/01-navigator.md) | Navigator/orientation/recovery kernel | 0 | approved |
| [2](stages/02-templates-okf.md) | OKF artifacts and approval | 1 | approved |
| [3](stages/03-codex.md) | Codex-native planning | 2 | approved |
| [4](stages/04-claude-code.md) | Claude Code parity | 3 | approved |
| [5](stages/05-context-and-roles.md) | Roles, context, task packets | 4 | approved |
| [6](stages/06-pinky-brain.md) | Sequential delegation | 5 | approved |
| [7](stages/07-verification-completion.md) | Verification and completion | 6 | approved |
| [8](stages/08-doctor-recovery.md) | Complete Doctor/recovery | 7 | approved |
| [9](stages/09-hardening.md) | Adoption hardening and platform proof | 8 | implemented; awaiting human gate |

Stages do not overlap gates. Waves show readiness only, not permission for unsafe concurrent mutation. Parallel Pinkies remain deferred beyond Stage 9.

## Plan set

- [Repository analysis](00-repository-analysis.md)
- [Architecture](01-architecture-plan.md)
- [Contracts and artifacts](02-contracts-and-artifacts.md)
- [Integrations](03-integration-plan.md)
- [Testing and evaluations](04-testing-and-evaluation-plan.md)
- [Risks and decisions](05-risk-and-decision-register.md)
- Detailed Stage 0–9 plans linked in table above

## Major fixed decisions

Go 1.26/Linux amd64 first; standard-library `flag`; official MCP SDK via stdio; explicit human ready approval; explicit execution policy with no delegated fallback; owned minimal OKF profile and read-only optional Aido; honest sequential host baseline; Pathframe-run bounded argv verification; journal/run reconstruction; no initial Specd compatibility.

## Cross-stage gates

Each stage requires implementation complete, focused tests, full affected suite, demonstrated journey, aligned docs, recorded risks, and human stop/go. No later stage starts with unresolved acceptance failures. Every non-terminal state needs an executable recovery. No recovery depends on deletion or Git rewriting.

## Current gate

Do not begin post-Stage-9 work until a human approves the stable-release gate. Verify Stage 9 with:

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test -race ./...
go test ./... -run Journey
go test ./... -run Security
go test ./... -run Install
go test ./... -bench 'Discovery|Validation|Packet' -run '^$'
GOOS=linux GOARCH=amd64 go build ./cmd/pathframe
GOOS=linux GOARCH=arm64 go build ./cmd/pathframe
GOOS=darwin GOARCH=amd64 go build ./cmd/pathframe
GOOS=darwin GOARCH=arm64 go build ./cmd/pathframe
GOOS=windows GOARCH=amd64 go build ./cmd/pathframe
git diff --check
```

Acceptance: Linux amd64 passes clean install and the full sequential journey; budgets and security checks pass; documentation matches behavior; other systems remain targets until native install/journey evidence exists; no Specd compatibility exists; limitations are explicit. Stable release requires a human go decision.

# Pathframe Implementation Roadmap

## Current state

Stages 0 through 4 are approved. Stage 5 is implemented and awaiting its human gate: role contracts, four-layer bounded context, explicit byte budgets and omissions, `pathframe.task/v1` preview, and deterministic frontier/waves are present. Worker execution remains unimplemented.

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
| [5](stages/05-context-and-roles.md) | Roles, context, task packets | 4 | implemented; awaiting human gate |
| [6](stages/06-pinky-brain.md) | Sequential delegation | 5 | planned |
| [7](stages/07-verification-completion.md) | Verification and completion | 6 | planned |
| [8](stages/08-doctor-recovery.md) | Complete Doctor/recovery | 7 | planned |
| [9](stages/09-hardening.md) | Adoption hardening and platform proof | 8 | planned |

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

## Next implementation task after the gate

Do not begin Stage 6 until a human approves the Stage 5 gate. After approval, the next task is S6-T1 in [Stage 6](stages/06-pinky-brain.md).

Stage 5 remains at its stop/go gate until then. Verify it with:

```sh
go test ./internal/context ./internal/delegation ./internal/app
go test ./tests/journey -run 'Packet|Frontier|Context'
go test ./...
go vet ./...
```

Acceptance: human and machine packet views agree; required context never truncates; unrelated context is absent; optional omissions are visible; frontier and waves are deterministic; cycle and budget failures provide non-destructive recovery; no worker launches and no Stage 6 capability is present.

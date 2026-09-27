# Pathframe Implementation Roadmap

## Current state

Repository contains only an empty README and two authoritative planning documents. No production source, tests, module, CI, integration, or managed state exists. Go 1.26.4 Linux amd64 is available. Git history confirms the accepted decisions; no contradiction or new blocking decision was found.

## Architecture

One canonical application layer coordinates deterministic workflow, artifacts, context, delegation, verification, and recovery domains. CLI and MCP are thin adapters over those operations. Filesystem, process, optional Git observation, MCP SDK, and host integrations remain edge adapters. Authored artifacts express intent; append-only transitions and bounded runs preserve facts; `state.json` is rebuildable. Brain owns reasoning and acceptance; Pinky executes only explicitly delegated tasks.

## Dependency order and status

| Stage | Capability | Depends | Status |
| --- | --- | --- | --- |
| [0](stages/00-foundation.md) | Foundation and constitution | none | planned |
| [1](stages/01-navigator.md) | Navigator/orientation/recovery kernel | 0 | planned |
| [2](stages/02-templates-okf.md) | OKF artifacts and approval | 1 | planned |
| [3](stages/03-codex.md) | Codex-native planning | 2 | planned |
| [4](stages/04-claude-code.md) | Claude Code parity | 3 | planned |
| [5](stages/05-context-and-roles.md) | Roles, context, task packets | 4 | planned |
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

## Exact first implementation task

`S0-T1 — Create Pathframe product constitution documents.`

Starting from the two root authority files, replace the empty `README.md` and create `PHILOSOPHY.md`, `ARCHITECTURE.md`, `CONTRIBUTING.md`, `AGENTS.md`, and `docs/decisions/README.md`. Preserve the exact product statement and foundational philosophy. Record all ten accepted decisions, initial non-goals, layer/dependency rules, Stage 0 Linux amd64 support claim, deviation-record process, and prohibition on copying legacy subsystems without a traced reason. Do not add Go code or production behavior. Verify with:

```sh
rg -n "Pathframe is a local, deterministic|The agent reasons" README.md PHILOSOPHY.md ARCHITECTURE.md
rg -n "Go 1.26|Linux amd64|flag|stdio|approval|execution_policy|okf-markdown/v1|sequential|timeout|history.jsonl|Specd" README.md PHILOSOPHY.md ARCHITECTURE.md CONTRIBUTING.md AGENTS.md docs/decisions/README.md
git diff --check
```

Acceptance: a cold contributor can state purpose, boundaries, supported platform, non-goals, and deviation process from these documents alone; no production file exists.


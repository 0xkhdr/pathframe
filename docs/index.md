# Documentation guide

Purpose: provide a developer-friendly reading path and identify the owner of each Pathframe topic.

All files in `docs/` use lowercase kebab-case names. A filename states the page's subject. Each page starts with a purpose sentence.

## First read

Read these pages in order when joining the project:

1. [Project README](../README.md) — product promise, status, boundaries, and document map.
2. [Core concepts](concepts.md) — actors, workflow, artifacts, delegation, verification, persistence, and recovery vocabulary.
3. [Philosophy](../PHILOSOPHY.md) — product principles and non-goals.
4. [Architecture](../ARCHITECTURE.md) — dependency direction, fixed decisions, and implemented stages.
5. [Contributing](../CONTRIBUTING.md) — repository layout, development loop, and completion gate.

## Work with Pathframe

1. [Installation](installation.md) — build, install, update, and uninstall.
2. [Getting started](getting-started.md) — complete a small direct workflow.
3. [Planning artifacts](planning-artifacts.md) — author changes, tasks, roles, and context.
4. [CLI reference](cli-reference.md) — find commands and operation contracts.
5. [Host integrations](host-integrations.md) — use Codex or Claude Code through typed MCP operations.
6. [Recovery](recovery.md) — diagnose interruption or damaged machine-owned state.

## Operate and maintain Pathframe

- [Security model](security-model.md) — trust boundaries and enforced controls.
- [Compatibility and limitations](compatibility-and-limitations.md) — support claims and known limits.
- [Release checklist](release-checklist.md) — release evidence and approval gate.
- [Decision records](decisions/index.md) — intentional deviations and new architecture decisions.

## Documentation ownership

Update one owning page instead of repeating the same contract:

- Product purpose and boundaries: `README.md`, `PHILOSOPHY.md`.
- Architecture and fixed decisions: `ARCHITECTURE.md`.
- Developer workflow: `CONTRIBUTING.md`.
- Vocabulary and lifecycle model: `docs/concepts.md`.
- User-facing commands: `docs/cli-reference.md`.
- Persisted planning formats: `docs/planning-artifacts.md`.
- Host-generated assets and MCP use: `docs/host-integrations.md`.
- Recovery behavior: `docs/recovery.md`.
- Security, support, and release evidence: their named pages.

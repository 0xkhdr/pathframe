# Architecture Plan

## Product constitution

> Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.

> The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.

## Layers and dependency rule

```text
cmd/pathframe
  -> internal/adapters/cli ---------+
  -> internal/adapters/mcp         |
  -> internal/integrations/*       v
                              internal/app
                                    |
        +---------------------------+---------------------------+
        v             v             v             v             v
 internal/workflow artifacts     context      delegation     recovery
        +---------------------------+-------------+-------------+
                                    v
                    internal/store | verification | gitobserve
```

Adapters translate. `app` owns use-case sequencing. Domain packages own deterministic rules. Edge packages own I/O. No adapter calls another adapter, and no domain package imports a host SDK.

## Proposed repository layout

```text
cmd/pathframe/                 executable
internal/app/                  canonical operations and results
internal/workflow/             phases, task states, transitions, frontier
internal/artifacts/            OKF profile, schemas, parsing, validation
internal/context/              reference resolution and packet assembly
internal/delegation/           preflight, leases, result reconciliation
internal/verification/         bounded argv process execution
internal/recovery/             diagnosis, replay, repair plans
internal/store/                project discovery and filesystem persistence
internal/gitobserve/           optional read-only Git facts
internal/adapters/cli/         flag dispatcher and text/JSON projection
internal/adapters/mcp/         official SDK and typed stdio tools
internal/integrations/codex/   generated Codex assets/capabilities
internal/integrations/claude/  generated Claude Code assets/capabilities
schemas/                       cross-process/persisted JSON schemas
templates/                     embedded Markdown and host guidance
testdata/                      fixtures and golden results
tests/journey/                 black-box journeys
docs/                          user and architecture documentation
```

## Canonical operations

`app.Service` exposes typed use cases: Orient, AssessRequest, CreateChange, GetTemplate, ValidatePlan, ApprovePlan, GetNext, Recover, PrepareDelegation, SubmitResult, RunVerification, AcceptTask, RequestChanges, and Doctor. CLI and MCP invoke these same operations and project the same versioned results.

## State and ownership

- Authored: `.pathframe/project.yaml`, knowledge, roles, change metadata, Markdown artifacts, task files.
- Durable machine facts: append-only `history.jsonl` plus bounded immutable run records.
- Replaceable projection: `state.json`, rebuilt by deterministic replay.
- Git and conversation history: observations only, never reconstruction sources.
- Writes use project-contained paths, validated components, temporary sibling files, sync/close, and rename where replacement is needed.

## Accepted decisions encoded

Go 1.26; Linux amd64 first; standard `flag`; no-argument orientation; official MCP Go SDK via stdio at Stage 3; explicit approval for every `planning -> ready`; explicit task execution policy; Pathframe-owned `okf-markdown/v1`; optional read-only Aido; sequential shared-workspace delegation; Pathframe-run argv verification with timeouts; journal/run replay; no initial Specd compatibility.

## Versioning and change discipline

Version only persisted or process boundaries. Internal interfaces change with their callers. Material edits to approved plan artifacts invalidate approval; routine state transitions do not. Schema additions remain backward-readable within v1 where optional; incompatible meaning requires v2 and migration/replay policy.

## Security boundary

Treat artifact content, paths, host results, and command arguments as untrusted. Reject path escape and symlink escape. Never invoke an implicit shell. Bound output and duration. Never claim containment from advisory scope. Doctor repairs machine-owned projections only.

## Deferred architecture

No parallel workers, custom sandbox, Streamable HTTP, generic MCP command runner, plugin registry, deployment orchestration, cryptographic ledger, strict Git binding, or Specd reader.


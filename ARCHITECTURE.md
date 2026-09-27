# Pathframe Architecture

## Product boundary

> Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.

> The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.

Pathframe determines legal transitions, validates artifacts, computes readiness, assembles context, validates delegation, runs approved verification, and diagnoses recovery. Host agents interpret intent, make semantic decisions, implement work, launch native subagents, and accept outcomes.

## Dependency direction

```text
Human CLI / Codex / Claude Code
              |
       CLI and MCP adapters
              |
     canonical application operations
              |
workflow | artifacts | context | delegation | verification | recovery
              |
 filesystem persistence | optional Git observation | host/process edges
```

Planned package direction:

- `cmd/pathframe` starts the executable and depends on the CLI adapter.
- `internal/adapters/{cli,mcp}` translates external requests into `internal/app` operations.
- `internal/integrations/{codex,claude}` owns generated host files and capability declarations.
- `internal/app` coordinates use cases without owning domain rules.
- `internal/{workflow,artifacts,context,delegation,verification,recovery}` owns deterministic domain behavior.
- `internal/{store,gitobserve}` and process/host adapters own I/O.

Adapters must not call other adapters. Domain packages must not import host SDKs. CLI, MCP, skills, commands, and hooks must not duplicate lifecycle decisions.

## Fixed implementation decisions

1. Use Go 1.26. Support Linux amd64 first. Other named platforms become supported only after Stage 9 CI and journey validation.
2. Implement the CLI with the standard `flag` package and a small subcommand dispatcher. Running `pathframe` with no arguments provides orientation.
3. Implement local stdio MCP with the official MCP Go SDK, pinned when Stage 3 begins and confined to the MCP adapter. Do not add Streamable HTTP or custom JSON-RPC initially.
4. Require explicit human approval for every `planning -> ready` transition. Material changes require validation and approval again.
5. Every task declares `execution_policy`. Quick tasks default to `brain`; delegated failure never transfers authority to Brain. Changing approved policy requires reapproval.
6. Pathframe owns minimal `okf-markdown/v1` and versioned artifact schemas. Aido is optional and read-only; neither product writes the other's directory or shares lifecycle state.
7. Baseline delegation is one sequential host-native subagent in the shared workspace with a task packet and returned final result. Scope is advisory unless enforcement is declared.
8. Pathframe executes approved verification as structured argument arrays inside the project with a required timeout. No implicit shell. Worker verification is supplemental; completion also requires semantic acceptance.
9. Authored artifacts, append-only `history.jsonl`, and bounded run records are reconstruction sources. `state.json` is rebuildable. Git and conversation history are not authorities.
10. Initial release has no Specd import, compatibility reader, or automatic `.specd/` detection. Stage 9 may reconsider a one-time importer only with real demand and representative fixtures.

## Cross-boundary contracts

Version only persisted or process boundaries: `pathframe.workflow/v1`, `pathframe.artifact-instructions/v1`, `pathframe.task/v1`, `pathframe.task-result/v1`, and `pathframe.integration/v1`. Internal abstractions change with their callers and do not need registries.

## State ownership

Humans and agents author change artifacts. Pathframe owns transition and run records plus replaceable projections. Filesystem paths, host results, and verification arguments are untrusted. Pathframe rejects project-root and symlink escapes, bounds output and duration, and never executes verification through an implicit shell.

## Stage 1 state and recovery

The Navigator domain owns phase/task-state transition tables and deterministic legal actions without filesystem imports. The store discovers `.pathframe`, rejects escaping change paths, appends `pathframe.transition/v1` records, replays only newline-complete records, and atomically replaces `state.json`. The application layer selects a change, coordinates replay and transitions, and returns `pathframe.workflow/v1`; CLI text and JSON are projections of that result.

Stage 1 does not implement authored artifacts, templates, approval workflows, MCP, delegation, or verification execution. Those remain gated by later stages.

## Stage 2 artifacts and approval

The `artifacts` domain owns the strict, flat-front-matter `okf-markdown/v1` profile, embedded templates, parsing, validation, task-per-file contracts, and normalized material identity. Quick requires intent and tasks; Standard adds requirements and design; High-risk adds risks, rollout, and recovery. The application layer owns create, template, check, and explicit human approval operations. Approval is a journal transition bound to the validated identity; canonical orientation detects a material mismatch and journals recovery to `replanning`. Authored files are preserved. CLI remains a projection of these operations. No host, MCP, delegation, context, Aido, or verification execution boundary was added.

## Stage 3 Codex planning integration

The Codex generator owns repository skill, project MCP configuration, optional `SessionStart` orientation, and `pathframe.integration/v1`. The MCP adapter alone imports the pinned official SDK and exposes only typed planning operations backed by `internal/app`. Request classification remains Brain reasoning; the application operation deterministically maps classified facts to `must_use`, `offer`, or `must_not_use`. Explicit approval is handed back through the typed validation operation only after the human approves. Generated assets are updated only when their manifest hash proves Pathframe ownership.

## Stage 4 Claude Code planning integration

The Claude Code generator owns a project skill, slash command, `.mcp.json`, optional supported `SessionStart` command hook, and a host-specific `pathframe.integration/v1` manifest. These assets call the same typed MCP and canonical application operations as Codex. Hash ownership makes install/update non-destructive, and host-specific Doctor checks report drift without editing workflow artifacts.

## Stage 5 context and roles

The `context` domain owns strict role parsing, contained foundation/change references, required/optional byte budgets, omission accounting, and `pathframe.task/v1` assembly. Task and runtime layers are selected canonical facts rather than accumulated files. At the Stage 5 gate, the `delegation` domain owned only deterministic dependency frontier/wave projection, and `app.PrepareDelegation` was a read-only preview.

## Stage 6 sequential delegation

The `delegation` domain validates declared host capabilities, enforces delegated policy, owns one exclusive active lease, bounds append-only run records, validates `pathframe.task-result/v1`, and produces reconciliation without acceptance. Host integrations declare only native sequential/shared-workspace/result-return facts and provide Pinky instructions; Brain performs the native launch. `app.PrepareDelegation`, `SubmitResult`, `ReleaseLease`, and `CheckBrainEdit` are canonical operations shared by CLI/MCP. Scope remains advisory unless the host manifest explicitly declares enforcement. Failed launch, worker failure, and released leases never change `execution_policy` or authorize Brain fallback.

## Stage 7 verification and completion

The `verification` domain validates project-contained working directories, executes approved argv directly with required timeouts, independently bounds stdout/stderr, computes relevant repository content identity, and compares reported changed files with advisory scope. Append-only run events retain submissions, every verification outcome, and Brain review decisions. `app.RunVerification`, `AcceptTask`, and `RequestChanges` are shared by CLI and MCP. Passing mechanical verification never implies semantic acceptance; acceptance rejects stale content and advances the completed-task frontier or completes the change. No shell, destructive revert, parallel worker, or expanded Doctor behavior was added.

## Stage 8 Doctor and recovery

The `recovery` domain owns typed diagnosis, deterministic projection comparison/rebuild, abandoned-lease classification, and normalized task-contract identity. `app.Doctor` composes owning validators and exposes read-only preview plus narrowly safe repair through CLI and MCP. Authored artifacts and complete journal/run records are never machine-repaired. Cancellation and lease recovery retain append-only run evidence, while replanning preserves only completions whose task contracts remain unchanged.

## Stage 9 hardening and support boundary

Linux amd64 is the only supported production platform because it has native CI, installation, and complete sequential journey evidence. Other named targets are cross-build-only portability targets. The repository installer atomically replaces only the executable and never touches project state. Security tests cover managed-directory and context symlinks, portable result paths, and literal structured argv execution. Performance benchmarks cover discovery, plan validation, and packet assembly. No new lifecycle contract, Specd compatibility, parallel worker, or release/deployment subsystem was added.

## Legacy reuse rule

Pathframe is a fresh project, not Specd renamed. Do not copy legacy code until a scoped task traces its actual behavior and tests, confirms license and architectural fit, and records why reuse is smaller and safer than a fresh implementation.

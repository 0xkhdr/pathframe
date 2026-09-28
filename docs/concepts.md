# Core concepts

Purpose: define Pathframe vocabulary and show how the complete development path fits together.

## Product model

Pathframe is a local deterministic protocol for structured software changes. It records the approved path and validates legal operations. It does not reason about product intent or implement code by itself.

The **Brain** is the primary coding agent, such as Codex or Claude Code. Brain interprets intent, designs work, coordinates execution, reviews results, and makes semantic decisions.

A **human** owns explicit plan approval and human-only lifecycle decisions such as pause and cancel. Every transition from `planning` to `ready` requires explicit human approval.

**Pinky** is one host-native subagent explicitly delegated one approved task. Pinky receives a bounded packet, works within the declared role and write scope, and returns a structured result. Pinky cannot approve plans, accept its own work, or change `execution_policy`.

## Canonical operations and adapters

A **canonical application operation** is the single implementation of a Pathframe use case in `internal/app`. CLI and MCP are adapters over these operations. They may render inputs and outputs differently, but they must not duplicate workflow decisions.

The **CLI adapter** is the human-facing `pathframe` command. The **MCP adapter** exposes typed tools to coding agents through local stdio. Pathframe has no generic raw-command MCP tool.

A **typed operation** has a named input and output contract. This prevents agents from constructing shell commands, guessing argument order, or bypassing application rules.

## Change lifecycle

A **change** is one structured unit of development work under `.pathframe/changes/<id>/`. Its phase describes the whole change:

```text
exploring -> planning -> ready -> executing -> reviewing -> ready or done
                    ^                    |
                    |                    +-> blocked -> executing
                    +----- replanning <--+

Any non-terminal phase -> paused -> prior phase
Any non-terminal phase -> cancelled
```

- `exploring`: change exists but planning has not started.
- `planning`: authored artifacts are being prepared and validated.
- `ready`: current plan identity has explicit human approval and work may start.
- `executing`: an approved task is active.
- `reviewing`: implementation was submitted and awaits verification or semantic review.
- `blocked`: execution cannot continue through the current path; recovery remains available.
- `replanning`: approved material changed or the path must be revised.
- `paused`: human suspended a non-terminal phase; the prior phase is retained for resume.
- `done`: all approved tasks completed. Terminal.
- `cancelled`: human ended the change. Terminal.

A **transition** is an actor, action, source phase, and target phase recorded in `history.jsonl`. Illegal transitions return current state, a stable reason code, and supported recovery actions.

## Plans and artifacts

A **planning mode** sets minimum rigor:

- `quick`: `intent.md` and task files.
- `standard`: quick artifacts plus `requirements.md` and `design.md`.
- `high-risk`: standard artifacts plus `risks.md`, `rollout.md`, and `recovery.md`.

An **artifact** is Pathframe-owned `okf-markdown/v1`: human-readable Markdown with flat front matter and required sections. `change.yaml` identifies the change, schema, profile, and mode. Each task has its own file.

A **task contract** defines objective, `execution_policy`, role, dependencies, references, required reads, write scope, constraints, structured verification commands, observable acceptance, questions, and assumptions.

**Plan validation** checks structure and deterministic completeness. It does not decide whether the plan is a good product decision.

A **plan identity** is the normalized identity of material plan content. Whitespace-only normalization does not change it. Content, references, mode, or execution-policy changes do. A changed identity invalidates approval and moves work toward replanning.

**Approval** binds explicit human consent to one valid plan identity. Validation success alone never implies approval.

## Tasks, roles, and execution authority

Task states are `pending`, `ready`, `active`, `submitted`, `completed`, `blocked`, and `changes_requested`. Dependencies determine when a pending task becomes ready.

`execution_policy` assigns implementation authority:

- `brain`: Brain may implement the task directly.
- `delegated`: Pinky must implement the task. Delegation failure never grants Brain fallback authority.

A **role** is a `pathframe.role/v1` file in `.pathframe/roles/`. It defines mission, required and optional foundation reads, allowed actions, and required return fields.

**Write scope** lists project-relative files a task may change. It is advisory unless a host explicitly declares enforcement. Pathframe detects and reports scope violations; it does not claim filesystem containment.

## Dependency planning and context

The **frontier** is the sorted set of approved tasks whose dependencies are complete. The **next task** is one deterministic frontier task selected for sequential execution. **Waves** group tasks by dependency readiness; they describe structure but do not authorize parallel mutation.

A **task packet** is `pathframe.task/v1`. It combines the approved task contract with selected context and current workflow facts. Its layers are:

- `foundation`: role reads from the project.
- `change`: task-required reads from the change directory.
- `task`: the selected task contract.
- `runtime`: current canonical workflow facts.

The **context budget** limits packet bytes. Required context must fit completely or packet creation fails. Optional context is included in declared order while space remains. Every omission records its layer, path, size, and reason.

## Delegation

**Preflight** validates host capabilities, delegated policy, explicit role, non-empty write scope, and verification commands before delegation.

A **lease** is one exclusive active delegation record. It binds a task, host, and random lease ID. Baseline Pathframe permits only one sequential Pinky. Brain launches the host-native worker; Pathframe does not launch it.

A **task result** is bounded `pathframe.task-result/v1` data. It includes lease identity, status, summary, changed files, worker-reported checks, discoveries, questions, and risks. Submission validates and records the result but never accepts the task.

## Verification and review

**Pathframe verification** runs approved argument arrays directly. Each command has a required timeout, a project-contained working directory, and bounded stdout and stderr. No implicit shell is used.

**Worker verification** is supplemental evidence reported by Pinky. It never replaces Pathframe verification.

A **content identity** binds verification to relevant repository content. If content changes after verification, acceptance is stale and must be verified again.

**Semantic acceptance** is Brain's separate judgment that implementation meets the approved objective and acceptance criteria. A task completes only after fresh passing Pathframe verification and semantic acceptance. Brain may instead request changes without changing execution authority.

## Persistence and recovery

**Authored state** includes `change.yaml`, artifact Markdown, task files, and role files. Humans and agents own this content.

`history.jsonl` is the append-only workflow transition journal. **Run records** are bounded append-only delegation, verification, and review events. These records plus authored artifacts reconstruct lifecycle state.

`state.json` is a replaceable projection for fast orientation. It is not an authority. Doctor may rebuild it from complete source records.

**Doctor** returns typed diagnoses. With `--repair`, it may rebuild `state.json`, discard an incomplete trailing journal fragment, or release an abandoned lease through recorded events. Doctor never rewrites authored artifacts, complete journal or run records, or Git history.

**Replanning** preserves authored files and append-only evidence. After reapproval, only completed tasks whose normalized contracts remain unchanged contribute to progress.

## Integrations and external boundaries

A **host integration** is a manifest-owned set of generated skill, MCP, and optional session-orientation files for Codex or Claude Code. Ownership hashes prevent Pathframe from overwriting user-modified files.

**Aido** is optional read-only knowledge. Pathframe does not require it, write its directory, or share lifecycle state with it.

**Specd** has no initial Pathframe importer, reader, or automatic detection. This is an accepted product boundary, not missing recovery behavior.

Git is optional observation, not workflow authority. Conversation history is also not reconstruction authority. Every supported non-terminal state must remain recoverable without deleting `.pathframe/` or rewriting Git history.

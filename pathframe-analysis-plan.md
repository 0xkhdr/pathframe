# Pathframe Analysis Plan

## Purpose

This document defines what Pathframe must become, the boundaries its implementation must preserve, and a staged path for building and validating the project without recreating the complexity of Specd.

It is an analysis and product architecture plan. It deliberately does not prescribe every package, function, or database representation. The coding agent must first inspect the new repository and convert this analysis into a repository-specific implementation plan.

## Product definition

> **Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.**

Pathframe gives humans and agents a shared answer to five questions:

1. Where are we?
2. What is ready now?
3. What comes next?
4. Why is progress blocked?
5. Where can a human intervene?

Its foundational statement is:

> **The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.**

## Product identity

| Item | Value |
| --- | --- |
| Product | Pathframe |
| Repository | `0xkhdr/pathframe` |
| CLI | `pathframe` |
| Managed directory | `.pathframe/` |
| Suggested Go module | `github.com/0xkhdr/pathframe` |
| Host coordinator | Brain |
| Delegated worker | Pinky |
| Workflow kernel | Navigator |
| Diagnosis and recovery | Doctor |

## Core product promises

### Orientation

A new user or cold coding agent can run one operation and understand:

- whether Pathframe is configured;
- whether a change is active;
- the current phase;
- progress and blockers;
- the recommended next action;
- other legal options;
- whether human action is required.

### Navigation

The same state always produces the same set of legal next actions. Pathframe recommends one action but exposes alternatives when they are valid.

### Coordination

Pathframe turns an approved change into bounded task contracts that a host agent can execute directly or delegate to Pinky workers.

### Recovery

Every reachable non-terminal state has a supported route to resume, retry, replan, repair, pause, or cancel. Deleting `.pathframe/` is never the normal recovery path.

## Philosophy and engineering rules

### Path before enforcement

Every refusal must first explain the current state, then explain why the requested transition is illegal, and finally provide a supported recovery action.

### Progressive rigor

Pathframe supports different levels of planning rather than applying maximum ceremony to every change:

| Mode | Appropriate for | Required artifacts |
| --- | --- | --- |
| Quick | Small, well-understood change | Intent and task contracts |
| Standard | Multi-step feature or bug | Intent, requirements, design, tasks |
| High-risk | Security, migration, contract, cross-system work | Standard plus risks, rollout, and recovery |

The mode may be raised when new risk is discovered. It must never be silently lowered.

### Agent-native and human-readable

- Markdown is the human-readable artifact surface.
- Typed operations are the agent surface.
- CLI and MCP project the same application operations.
- Agents must not construct raw CLI strings when a typed tool exists.
- Machine identifiers are round-tripped by integrations rather than remembered by the model.

### Context is selected, not accumulated

Workers receive only the project knowledge, change facts, task contract, and runtime state required for their assignment.

### Explicit delegation authority

- Brain coordinates.
- Pinky changes implementation files for delegated tasks.
- If a Pinky cannot be launched, Brain reports a blocker.
- Brain must not silently implement a delegated task itself.
- Direct Brain execution is allowed only when the task policy explicitly selects it.

### Host capability honesty

Pathframe may declare scope and validate results, but it must not claim filesystem containment unless the host provides it. Capabilities such as subagents, sandboxes, hooks, and write restrictions must be detected or declared explicitly.

### Complexity must be earned

Do not restore a Specd mechanism merely because it already exists. Reuse it only when it directly supports the simplified Pathframe workflow.

## Accepted implementation decisions

These decisions are approved constraints for repository planning and implementation:

| Concern | Decision |
| --- | --- |
| Language and platform | Use Go 1.26. Initially support Linux amd64. Treat Linux arm64, macOS amd64/arm64, and Windows amd64 as portability targets until Stage 9 proves them with CI and end-to-end journeys. |
| CLI | Use Go's standard `flag` package with a small subcommand dispatcher. Add no CLI framework. The no-argument command renders orientation rather than bare usage. |
| MCP | Use local stdio and the official MCP Go SDK, pinned when Stage 3 begins. Confine the SDK to the MCP adapter. Defer Streamable HTTP and do not hand-roll JSON-RPC. |
| Approval | Require explicit human approval for `planning -> ready` in every mode. Quick mode may validate and obtain approval in one interaction. Material plan changes require revalidation and reapproval; routine transitions do not. |
| Task execution | Quick tasks default to explicit `execution_policy: brain`. Other tasks also declare their policy. A delegated task never falls back to Brain, and changing approved execution policy requires reapproval. |
| OKF and Aido | Pathframe owns a minimal `okf-markdown/v1` profile plus versioned Pathframe schemas. Aido is optional and read-only. There is no shared lifecycle, required service, shared package, automatic synchronization, or cross-writing. |
| Host capability baseline | Require only one sequential host-native subagent, a supplied task packet, shared-workspace access, and a returned final result. Scope is advisory unless enforcement is declared. Do not assume sandboxing, worktrees, cancellation, resumption, hooks, or parallelism. |
| Verification | Pathframe executes approved verification with structured argument arrays, a contained working directory, and a timeout. Initial execution does not invoke an implicit shell. Worker reports are supplemental; mechanical verification and semantic acceptance are both required. |
| Reconstruction | Authored artifacts define intended work; an append-only transition journal and bounded run records preserve lifecycle facts; `state.json` is a replaceable projection. Git and conversations are not reconstruction authorities. |
| Specd | Provide no import, compatibility reader, or automatic `.specd/` detection initially. Reconsider only a one-time importer at Stage 9 if real demand and fixtures justify it. |

A later implementation may refine internal names or representations without reopening these product choices. Reversing one requires a newer explicit human decision and a recorded rationale.

## Non-goals for the initial product

- Becoming an LLM or agent runtime
- Supporting every coding agent
- Mandatory Git branches or commits
- Cryptographic evidence or audit guarantees
- Release and deployment orchestration
- Full project-management functionality
- Organization-wide identity and permission management
- Mandatory human approval at every transition
- Parallel code-writing workers in the first complete version
- A security sandbox implemented inside Pathframe

## Applicability model

Pathframe must define when a coding agent should activate it.

### Must use Pathframe

- The user explicitly asks to use Pathframe.
- A relevant Pathframe change is already active.
- The request continues an existing Pathframe workflow.

### Should offer Pathframe

- The implementation has multiple dependent steps.
- The work should survive a new session or agent.
- Requirements and acceptance criteria need agreement.
- Multiple roles or delegated workers would be useful.
- The change spans components or carries meaningful risk.

The agent asks once and waits for the user's choice.

### Must not activate automatically

- Explanation or question answering
- Read-only inspection or review
- Brainstorming before formalization
- Trivial isolated edits
- Non-development operational work
- Requests where the user explicitly chooses direct execution

Activation and non-activation examples must be part of the Codex and Claude Code integration evaluations.

## Workflow model

### Change phases

```text
exploring -> planning -> ready -> executing -> reviewing -> done
                ^          |          |           |
                |          v          v           v
                +------ replanning <- blocked <- corrections
```

| Phase | Meaning | Normal owner |
| --- | --- | --- |
| `exploring` | The problem and alternatives are still being understood. | Human + Brain |
| `planning` | Change artifacts are being created and checked. | Brain + human |
| `ready` | The plan is accepted and at least one task is runnable. | Brain |
| `executing` | A Brain or Pinky worker owns a task. | Worker |
| `reviewing` | Results are checked against requirements and acceptance. | Brain/reviewer |
| `blocked` | Progress needs a named decision, configuration repair, or replan. | Named owner |
| `replanning` | The plan is changing without destroying the change. | Brain + human |
| `done` | The change outcome is accepted. | Human or configured policy |
| `paused` | Automatic progression is intentionally suspended. | Human |
| `cancelled` | Work ended without claiming the outcome was achieved. | Human |

### Task states

Keep task activity small and separate from change phase:

```text
pending -> ready -> active -> submitted -> completed
                       |          |
                       v          v
                    blocked <- changes_requested
```

### Transition requirements

Each transition definition must contain:

- source phase/state;
- target phase/state;
- actor allowed to request it;
- deterministic preconditions;
- state written;
- canonical result;
- failure reason;
- at least one supported recovery action.

### No-dead-end invariant

Tests must execute a recovery journey from every non-terminal state. A text message that merely suggests recovery is insufficient.

## User experience model

### One obvious entry point

Running `pathframe` with no arguments must provide a useful home view. It must never respond with a bare usage dump when it can orient the user.

Example:

```text
Pathframe

Active change: improve-payment-retry
Phase: executing
Progress: 3/7 tasks complete
Active task: T4 - Add retry idempotency guard
Worker: pinky/backend-engineer

Recommended next action:
  Inspect the worker result

Other options:
  Pause the change
  Replan the change
  Diagnose the workflow
```

### Canonical navigation result

Every surface projects a versioned result such as:

```json
{
  "schema": "pathframe.workflow/v1",
  "change": "improve-payment-retry",
  "phase": "ready",
  "progress": {"completed": 3, "total": 7},
  "recommended": {
    "action": "delegate_task",
    "task": "T4",
    "role": "backend-engineer",
    "agent_allowed": true
  },
  "alternatives": ["inspect_plan", "pause", "replan", "cancel"],
  "human_required": false
}
```

Human text is a projection of this result, not a separate decision engine.

## Artifact and template model

### Proposed layout

```text
.pathframe/
  project.yaml
  knowledge/
  roles/
  changes/
    improve-payment-retry/
      change.yaml
      intent.md
      requirements.md
      design.md
      tasks/
        T1.md
        T2.md
      history.jsonl
      state.json
      runs/
```

The implementation plan may refine this layout, but it must preserve clear ownership:

- humans and agents author change artifacts;
- Pathframe owns runtime state and run records;
- broad project knowledge may be referenced from Aido or other OKF sources;
- task artifacts reference knowledge rather than duplicating it.

### OKF-compatible templates

Generated Markdown follows the Pathframe-owned `okf-markdown/v1` profile. Each artifact also declares its versioned Pathframe schema. The profile boundary keeps the core lifecycle independent from document rendering and from Aido.

Aido may be consumed as an optional, read-only source of approved knowledge or business intent. An explicit import or reference never starts implementation. Pathframe never writes `.aido/`; Aido never writes `.pathframe/`. A Pathframe intent seeded from Aido becomes an independent Pathframe artifact with its own validation, approval, and lifecycle. Pathframe remains fully usable when Aido is absent or removed.

Templates contain:

1. typed front matter;
2. clear human-facing headings;
3. agent instructions explaining how to fill each section;
4. examples where ambiguity is likely;
5. stable identifiers and references;
6. validation rules associated with the template version.

### Task-per-file model

Do not return to a wide Markdown table. Each task file contains:

- stable ID and title;
- execution policy: Brain or delegated;
- role/persona;
- objective;
- requirement references;
- dependencies;
- required read context;
- declared write scope;
- constraints;
- verification commands;
- observable completion criteria;
- questions and unresolved assumptions.

## Context engineering model

### Context layers

| Layer | Content | Lifetime |
| --- | --- | --- |
| Foundation | Architecture, conventions, glossary | Project |
| Change | Intent, requirements, decisions | Change |
| Task | Objective, refs, dependencies, scope, verification | Task |
| Runtime | Diff summary, test results, discoveries, blockers | Attempt |

### Assembly rules

- Required context is never silently truncated.
- Optional context may be omitted with a visible reason.
- References are preferred over repeated prose.
- A Pinky receives no unrelated tasks.
- Brain receives worker summaries and artifact references, not entire transcripts.
- If required context exceeds the budget, Pathframe recommends splitting the task or raising the budget explicitly.

### Role contract

Roles remain concise and stable:

```yaml
id: backend-engineer
mission: Implement backend behavior and focused automated tests.
reads:
  - architecture.backend
  - conventions.php
allowed_actions:
  - edit
  - test
returns:
  - summary
  - changed_files
  - verification
  - discoveries
  - questions
```

The role defines stable behavior. The task defines change-specific work. Knowledge files define project facts.

## Pinky and the Brain protocol

### Brain responsibilities

- Assess whether Pathframe applies.
- Orient the user.
- Create or select changes.
- Help author artifacts.
- Ask for human decisions.
- Read the ready frontier.
- Perform delegation preflight.
- Launch a host-native Pinky.
- Submit and reconcile results.
- Explain blockers and recovery.

### Pinky responsibilities

- Accept one task packet.
- Work within the assigned role and context.
- Modify only authorized files when host enforcement is available.
- Run declared verification.
- Return a structured result.
- Report discoveries instead of silently changing the plan.

### Delegation preflight

Before delegation, Pathframe validates:

- selected role exists;
- host declares a supported subagent capability;
- required task fields are complete;
- required context fits its budget;
- write scope is present and safe;
- verification is defined;
- no conflicting task lease exists;
- the integration version is compatible.

Failure creates a configuration blocker. Brain does not inherit permission to code.

The required host baseline is deliberately small: launch one host-native subagent with a task packet in the shared repository workspace, wait for its final response, and return that response for submission. The integration declares stronger capabilities explicitly. Unknown capabilities are unsupported. Declared write scope is advisory unless the host reports enforcement, and Pathframe never describes an advisory scope as containment.

### Task packet

The canonical `pathframe.task/v1` packet contains:

- change and task identities;
- role contract;
- objective and acceptance;
- requirement/design references;
- selected context;
- dependencies and completed prerequisites;
- write scope;
- verification;
- result schema;
- host assurance/capability facts.

### Worker result

The canonical `pathframe.task-result/v1` result supports:

- `completed`;
- `failed`;
- `blocked`;
- `needs_replan`.

It reports summary, changed files, verification, discoveries, questions, and risks. Pinky never marks itself accepted; Brain or a reviewer reconciles the result.

Pinky's verification report is diagnostic. After submission, Pathframe runs the approved structured verification commands and records their bounded output, exit status, duration, timeout or interruption, and repository content identity. Commands use argument arrays, a project-contained working directory, and a required timeout; implicit shell execution is not supported initially. A task completes only after this mechanical verification and semantic acceptance.

## Codex and Claude Code integration

### Three required surfaces

#### Generated guidance

`pathframe init` generates focused skills or commands for Codex and Claude Code. Their descriptions contain precise activation and non-activation rules. Detailed workflow material is loaded only when activated.

#### Typed tools

Expose intent-level MCP tools rather than a generic shell command:

- `pathframe_orient`
- `pathframe_assess_request`
- `pathframe_create_change`
- `pathframe_get_template`
- `pathframe_validate_plan`
- `pathframe_get_next`
- `pathframe_prepare_delegation`
- `pathframe_submit_result`
- `pathframe_recover`

#### Lifecycle hooks

Where the host supports them:

- session start injects current orientation;
- pre-delegation runs preflight;
- post-subagent requires result reconciliation;
- pre-edit warns or blocks Brain during delegated execution;
- session end surfaces unresolved active work and resume instructions.

Hooks call canonical Pathframe operations. They contain no independent workflow logic.

### Integration order

Codex is implemented and validated first. Claude Code follows using the same canonical operation contracts. Do not attempt broad tool support before both vertical slices are reliable.

## Git policy

The default workflow observes Git but does not make Git the lifecycle engine.

- No mandatory commit before ordinary transitions.
- No automatic commit, branch, reset, stash, or push.
- Moving `HEAD` does not automatically destroy the workflow.
- Changed-file observations may help validate task scope.
- Strict commit/evidence binding may later exist as an opt-in high-assurance policy.

## Recovery model

### Reconstructable state

Authored artifacts are durable intent. A minimal append-only transition journal records approvals and lifecycle changes, while bounded run records capture leases, submissions, verification, and reconciliation. `state.json` is a replaceable machine projection reconstructed from these sources. Run metadata is retained initially, with fixed limits on potentially large output; compaction waits for measured need. Git and conversation history are not reconstruction sources.

### Doctor responsibilities

`pathframe doctor` diagnoses:

- missing or incompatible integration files;
- invalid configuration;
- malformed artifacts;
- broken references;
- unknown roles;
- abandoned task leases;
- inconsistent projected state;
- unavailable host capabilities.

`pathframe doctor --repair` may rebuild machine-owned projection files. It must not silently rewrite authored intent.

### Required recovery routes

| Failure | Supported action |
| --- | --- |
| Worker crashed | Resume or release and retry the task |
| Wrong plan | Replan while preserving unaffected work |
| Missing role/configuration | Repair configuration and resume |
| Verification failed | Keep task active and retry |
| Scope needs expansion | Amend task and reconcile approval impact |
| State projection corrupt | Rebuild machine-owned state |
| User wants to stop | Pause or cancel without deletion |
| Multiple changes ambiguous | Require explicit selection |

## Architecture boundaries

### Suggested layers

```text
Codex skill / Claude command / Human CLI
                   |
              MCP + CLI adapters
                   |
          Application operations
                   |
workflow | artifacts | context | delegation | recovery
                   |
       filesystem + optional Git observations
```

### Deterministic boundary

Pathframe deterministically owns:

- artifact parsing and validation;
- readiness and dependency projection;
- legal state transitions;
- task frontier and waves;
- context manifest assembly;
- delegation preflight;
- result-shape and scope validation;
- recovery diagnosis;
- canonical CLI/MCP results.

The host agent owns:

- interpreting user intent;
- drafting artifact content;
- semantic design decisions;
- implementation;
- semantic review;
- determining whether discoveries change product intent.

### Cross-boundary contracts

Version only contracts that cross a process or persisted artifact boundary:

- `pathframe.workflow/v1`
- `pathframe.artifact-instructions/v1`
- `pathframe.task/v1`
- `pathframe.task-result/v1`
- `pathframe.integration/v1`

Avoid registries and ownership machinery for purely internal abstractions.

## Staged build plan

Each stage must produce a usable vertical slice, its own validation evidence, and a clear stop/go decision. The coding agent must not begin a later stage while the current stage has unresolved acceptance failures.

### Stage 0 - Repository foundation and product constitution

#### Goal

Create a clean Pathframe repository whose documentation prevents implementation drift.

#### Deliverables

- Go 1.26 module and minimal CLI entry point for Linux amd64
- standard-library `flag` parsing with a small subcommand dispatcher
- `README.md` with product statement and first-use vision
- `PHILOSOPHY.md` with principles and non-goals
- `ARCHITECTURE.md` with initial boundaries
- contributor and agent guidance
- basic CI for format, vet, tests, and build
- decision log for intentional deviations from this analysis

#### Validation

- Clean build on the primary development platform
- One smoke test for `pathframe --version` and `pathframe --help`
- Documentation review confirms no inherited Specd terminology or guarantees
- No copied legacy subsystem without a recorded reason

#### Gate

Proceed only when a new contributor can explain Pathframe’s purpose and non-goals from repository documents.

### Stage 1 - Navigator kernel and orientation

#### Goal

Make state and next actions understandable before implementing templates or orchestration.

#### Deliverables

- Project root discovery
- Minimal change/phase model
- Transition table
- Canonical `pathframe.workflow/v1` result
- Home command, `status`, and `next`
- Pause, resume, replan, and cancel transitions
- Text and JSON projections
- No-dead-end test matrix

#### Validation

- Golden tests for human and JSON views
- Transition-table tests
- Journey tests for every initial phase
- Recovery action executed from every non-terminal phase
- Cold user can identify state and options in under one minute

#### Gate

Do not proceed if orientation requires reading documentation or if any reachable phase has no recovery route.

### Stage 2 - Template engine and OKF profile

#### Goal

Generate self-describing, progressively rigorous planning artifacts.

#### Deliverables

- Template/profile interface
- Pathframe-owned `okf-markdown/v1` profile and versioned Pathframe artifact schemas
- Quick, Standard, and High-risk modes
- `new`, template instruction, and `check` operations
- Task-per-file parser
- Stable reference validation
- Explicit human plan-approval transition for every mode, with material-change invalidation
- Clear ownership of authored versus machine files

#### Validation

- Fixture changes for all three modes
- Round-trip parser tests
- Tests for incomplete placeholders, broken references, missing acceptance, missing verification, and unresolved questions
- Cold-agent exercise: fill templates using generated instructions only
- Quick mode proves no unused artifact is mandatory

#### Gate

Proceed only when template instructions are sufficient without hidden prompt knowledge.

### Stage 3 - Codex-native planning path

#### Goal

Make Pathframe native to Codex rather than a CLI the model must memorize.

#### Deliverables

- Generated Codex skill
- Activation and non-activation rules
- Typed MCP operations over local stdio using the official MCP Go SDK for orientation, change creation, templates, validation, approval handoff, and next action
- Integration version manifest
- Optional supported session-start orientation
- Negative activation evaluation set

#### Validation

- End-to-end Codex journey from user request to ready plan
- Agent never guesses CLI argument order
- Must-use, offer, and must-not-use prompt evaluations
- Fresh-session resume journey
- Schema parity tests between CLI and MCP results

#### Gate

Do not add Claude Code or delegation until Codex can complete the planning journey without raw command construction.

### Stage 4 - Claude Code parity

#### Goal

Deliver equivalent planning behavior through Claude Code’s native mechanisms.

#### Deliverables

- Generated Claude Code skill/commands
- MCP configuration support
- Relevant lifecycle hooks
- Shared operation-contract tests
- Tool-specific installation/update/doctor checks

#### Validation

- Same journey suite as Codex
- Equivalent state transitions and canonical results
- Fresh-session resume
- Broken integration configuration produces an actionable Doctor result

#### Gate

Proceed when both hosts satisfy the same product journeys without duplicating domain logic.

### Stage 5 - Context packages and roles

#### Goal

Produce bounded task packets before launching workers.

#### Deliverables

- Role/persona schema
- Foundation/change/task/runtime context resolver
- Required and optional context classification
- Context budgets and omission reporting
- `pathframe.task/v1` packet
- Packet preview for humans
- Task frontier and sequential wave projection

#### Validation

- Required context never truncates
- Irrelevant task context is absent
- Reference resolution and budget tests
- Role/task/knowledge separation tests
- Human preview agrees with the machine packet

#### Gate

Do not launch Pinkies until packet quality can be inspected and evaluated independently.

### Stage 6 - Pinky and the Brain vertical slice

#### Goal

Safely delegate and reconcile one sequential task.

#### Deliverables

- Host capability declaration using the minimal sequential shared-workspace baseline
- Delegation preflight
- Task lease
- Codex subagent adapter
- Claude Code subagent adapter
- Pinky worker instructions
- `pathframe.task-result/v1`
- Result submission and Brain reconciliation
- Brain edit guard for delegated tasks and explicit `execution_policy` enforcement
- Explicit no-fallback blocker

#### Validation

- Successful delegation journey on both hosts
- Missing role, missing host capability, oversized context, and conflicting lease tests
- Brain does not edit code after delegation failure
- Pinky cannot change plan state or accept itself
- `needs_replan` result returns control to Brain correctly

#### Gate

Proceed only when delegation failure is safe and understandable before any parallel execution work.

### Stage 7 - Verification and task completion

#### Goal

Make completion meaningful without recreating the heavy evidence system.

#### Deliverables

- Pathframe-run structured verification with argument arrays, contained working directories, timeouts, and bounded results
- Changed-file comparison against advisory scope
- Brain/reviewer acceptance transition
- Changes-requested loop
- Task completion and next-task selection
- Change-level completion policy

#### Validation

- Passing, failing, timed-out, and interrupted verification journeys
- Changes after worker verification are detected by a content/result identity appropriate to the chosen design
- Scope violations produce keep/revert/replan options rather than a dead end
- Completed work advances the deterministic frontier

#### Gate

Do not call a task complete unless verification and semantic acceptance are both represented explicitly.

### Stage 8 - Doctor and recovery completion

#### Goal

Eliminate delete-and-restart recovery.

#### Deliverables

- `doctor` diagnosis catalog
- Safe repair operations
- Machine-state reconstruction from authored artifacts, the transition journal, and bounded run records
- Abandoned lease recovery
- Replan impact analysis
- Preservation of unaffected completed tasks
- Cancellation and retained history

#### Validation

- Reproduce every historical Specd dead end and demonstrate a supported recovery
- Corrupt state fixtures
- Worker crash and session interruption journeys
- Repair never silently edits authored intent

#### Gate

Pathframe is not feature-complete until destructive restart is unnecessary for supported failure scenarios.

### Stage 9 - Full product hardening

#### Goal

Prepare the sequential product for real adoption.

#### Deliverables

- Cross-platform CI for claimed platforms
- Installation and update path
- Explicit confirmation that the initial release has no Specd compatibility; evaluate a one-time importer only if real demand and representative fixtures exist
- Performance budgets for repository discovery, plan validation, and packet assembly
- Security review of path handling and command execution
- Stable docs and onboarding journey
- Release process based on demonstrated support, not speculative claims

#### Validation

- Clean-machine installation tests
- End-to-end journeys on supported platforms
- Malformed input and path-containment tests
- Usability test: new user creates, executes, interrupts, and recovers a change
- Codex and Claude compatibility matrix

#### Gate

Release the first stable version only with explicit known limitations and proven user journeys.

### Deferred stage - Parallel Pinkies

Parallel mutation is not part of initial full development. Add it only after sequential delegation is reliable and users demonstrate a need.

Prerequisites:

- disjoint-scope analysis;
- host worktree or sandbox isolation;
- deterministic lease and merge policy;
- conflict reconciliation;
- cancellation propagation;
- end-to-end concurrent journey testing.

## Cross-stage quality requirements

Every stage plan must include:

- affected packages and files;
- new or changed public contracts;
- migrations, if any;
- focused tests;
- end-to-end journey tests;
- failure and recovery cases;
- documentation changes;
- explicit non-goals;
- rollback or correction strategy;
- acceptance checklist;
- commands used for verification.

Every stage ends with:

1. implementation complete;
2. focused tests passing;
3. full affected suite passing;
4. journey demonstrated;
5. documentation aligned;
6. unresolved risks recorded;
7. human stop/go decision.

## Decision status

The implementation decisions formerly required during planning are resolved in **Accepted implementation decisions**. The coding agent must carry them into repository-specific contracts and stages rather than reopening them. Any new decision that materially changes architecture or scope must state its options, recommended default, deferral consequence, and first blocked stage.

## Success metrics

- A cold user or agent orients in under one minute.
- Native integrations do not guess CLI arguments.
- Activation evaluation measures false activation and missed activation.
- No supported state requires deletion to recover.
- Brain never silently takes over a delegated task.
- Delegation configuration failures occur before code mutation.
- Quick mode creates no unused mandatory artifact.
- Required Pinky context is complete and unrelated context is excluded.
- A fresh session resumes the correct change and next action.
- Sequential end-to-end change completion succeeds on every claimed platform.

## Completion definition

Pathframe reaches initial full development when Stages 0 through 9 pass their gates and the following journey works through both Codex and Claude Code:

1. User submits a qualifying development request.
2. Agent correctly offers or activates Pathframe.
3. User understands the current phase and options.
4. Agent creates and fills the appropriate templates.
5. Plan is validated and accepted.
6. Navigator selects the next task.
7. Brain prepares a bounded packet and delegates to Pinky.
8. Pinky changes code, verifies it, and submits structured results.
9. Brain reconciles the result and advances the path.
10. A failure or interruption is recovered without deleting the change.
11. Remaining tasks complete sequentially.
12. The final outcome is reviewed and the change becomes done.

Parallel Pinkies, strict Git evidence, and organization-scale controls remain separate future decisions.

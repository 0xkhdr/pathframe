# Specd Reset

## Complete analysis and initial product plan

**Prepared for:** Khedr  
**Date:** 27 September 2026  
**Status:** Initial direction for discussion, not an implementation specification

---

## Executive summary

Specd began with a valuable and focused promise: give a human and any coding agent the same deterministic view of a development journey. At every moment, both should be able to answer:

1. Where are we?
2. What is ready now?
3. What comes next?
4. Why is progress blocked?
5. Where can the human intervene?

The current `specd-cli` proves that lifecycle, approval, scope, evidence, and durable transitions can be enforced locally. It also demonstrates the cost of treating every desirable assurance property as part of the first product. The workflow now contains enough ceremony and recovery constraints that the mechanism can obscure the path it was created to explain.

The reset should not discard the central idea. It should narrow it:

> **Specd is a deterministic development-path protocol for humans and coding agents. It creates structured change artifacts, exposes the current and next legal actions, and coordinates bounded task execution.**

Specd should not be an agent runtime, a security sandbox, a Git workflow, a release system, a universal policy engine, or a replacement for the host agent. The host supplies intelligence and execution. Specd supplies shared state, task contracts, context packages, transition rules, and recovery.

The recommended product is built around five decisions:

- **Progressive activation:** ordinary questions and trivial edits do not enter Specd automatically. Planned, multi-step, risky, or explicitly requested work does.
- **One obvious entry point:** `specd` or `specd status` always explains whether Specd applies, where the user is, available options, and the recommended next action.
- **Tool-native integrations:** Codex and Claude Code receive generated skills/commands and MCP tools with exact schemas. Agents select operations rather than construct CLI arguments from memory.
- **Pinky and the Brain as a protocol:** Brain remains the host agent. Pinkies are host-native subagents. Specd produces delegation packets and validates results; it does not become a proprietary multi-agent runtime.
- **Recoverable state:** every non-terminal state has a supported repair, cancel, or replan route. Deleting `.specd` must never be the normal recovery procedure.

The first release should deliberately exclude hard Git coupling, per-transition human approvals, append-only evidence machinery, production profiles, release qualification, and parallel code mutation. These can return only after observed usage proves they are needed.

---

## 1. The original product problem

Chat-based coding agents have no stable shared development position. A conversation may contain an idea, partial analysis, an implementation attempt, a changed requirement, and a failed test without a canonical answer to “what stage are we in?” Another agent joining later must infer state from prose and may continue from the wrong assumption.

Specd exists to externalize that state into inspectable files and deterministic transitions. Its job is not to make the agent intelligent. Its job is to make the development path legible.

### 1.1 Core outcome

A user should be able to open any supported coding agent and ask:

> “What is happening with this change?”

The agent should obtain a canonical answer from Specd, including:

- the selected change;
- the current phase;
- progress through tasks;
- any active delegation;
- ready, blocked, and completed work;
- the recommended next action;
- alternative legal actions;
- whether human input is required.

The same answer should be available directly from the CLI without an agent.

### 1.2 The three product promises

Specd should make only three foundational promises:

| Promise | Meaning |
| --- | --- |
| **Orientation** | A human or cold agent can immediately understand the current state. |
| **Navigation** | Specd returns the next legal action and its required inputs. |
| **Coordination** | Work is converted into bounded task contracts that can be delegated and reconciled. |

Enforcement supports these promises. It is not the product by itself.

---

## 2. What went wrong in the previous direction

The current implementation accumulated solutions to real problems, but too many were promoted into universal workflow requirements.

### 2.1 Assurance displaced navigation

The present workflow emphasizes exact approval identities, revisions, Git baselines, evidence applicability, ledgers, atomic recovery, maturity claims, and release proof. Those are valuable for a high-assurance execution harness. They are not all necessary to tell a developer and agent what to do next.

This caused a priority inversion:

> Specd became better at refusing an invalid transition than helping the user understand and recover from it.

### 2.2 Commands became an agent memory test

Agents guessed command names, argument ordering, identifiers, revisions, and attempt IDs even when MCP existed. This reveals an integration failure, not merely an agent failure.

If an agent must remember:

```text
specd complete <change> <task> --revision <observed-revision>
```

then the protocol is leaking into natural-language reasoning. A native tool should expose typed fields, valid choices, and descriptions. Opaque identifiers should be returned and round-tripped by the integration, not copied by the model.

### 2.3 Activation was undefined

Agents sometimes invoked Specd for unrelated tasks and sometimes ignored it for substantial development work. The tool lacked a crisp applicability decision.

This must become explicit:

- Specd should not govern ordinary explanation, repository reading, brainstorming, formatting, or a trivial isolated edit unless requested.
- Specd should govern explicitly requested Specd work, multi-step implementation, coordinated delegation, work requiring a durable plan, or changes where interruption/restart matters.
- Ambiguous cases should produce a recommendation, not silently activate.

### 2.4 The user lacked a home screen

The user did not know how to start, which stage was active, or which options were available. A collection of precise commands is not a product entry point.

Running `specd` with no arguments should be useful. It should act as a project dashboard:

```text
Specd is installed for this repository.

Active change: improve-payment-retry
Phase: execution
Progress: 3/7 tasks complete
Active task: T4 - add idempotency guard
Worker: pinky/backend-engineer

Recommended: inspect the worker result
  specd next

Other options:
  specd pause
  specd replan
  specd cancel
```

### 2.5 Recovery was exceptional instead of designed

Corrupted or stale flows frequently ended with deleting state and starting over. This violates a central deterministic-workflow principle: a state machine is useful only if every reachable state has a defined exit.

Recovery must be part of the main model, not a troubleshooting appendix. Every phase needs:

- resume;
- retry;
- replan;
- pause;
- cancel;
- repair/doctor.

### 2.6 Git became workflow state

Git commits, clean-tree requirements, exact baselines, and approved file lists introduced friction and dead ends. Git is excellent for source history, but it should not be the primary lifecycle engine for the simplified product.

The reset should use Git as optional context:

- observe the current branch and diff;
- optionally suggest a checkpoint;
- never create commits without an explicit host/user request;
- do not require a commit to transition between ordinary task states;
- do not invalidate the whole workflow merely because `HEAD` moved.

Higher-assurance Git binding can later exist as an opt-in policy.

### 2.7 Brain violated delegation when configuration failed

When Pinky configuration was incomplete, Brain sometimes implemented the task itself. That collapses the separation of responsibility precisely when the system is least certain.

The correct invariant is:

> If a task requires delegated execution and no valid worker can be launched, Brain must return a configuration blocker. It may not silently become the worker.

Fallback should be an explicit policy chosen by the user, never implicit behavior.

---

## 3. Lessons from the reference systems

### 3.1 OpenSpec

OpenSpec demonstrates several useful product choices:

- changes are self-contained folders;
- accepted behavior and proposed changes are separated;
- artifact dependencies determine what is ready;
- a quick path coexists with expanded control;
- agent-specific skills and slash commands are generated for many tools;
- plans are fluid and can be revised instead of forcing a rigid waterfall;
- progressive rigor keeps most changes lightweight;
- `instructions ... --json` provides agent-consumable guidance.

OpenSpec explicitly favors “fluid not rigid,” “iterative not waterfall,” and “easy not complex.” It also separates behavior requirements from implementation design. These are strong corrections for Specd’s current rigidity. See [OpenSpec Concepts](https://github.com/Fission-AI/OpenSpec/blob/main/docs/concepts.md) and [OpenSpec README](https://github.com/Fission-AI/OpenSpec/blob/main/README.md).

What Specd should not copy:

- a large catalog of user-facing workflow commands;
- dependence on the model correctly following long generated prompts;
- artifact existence as a proxy for artifact quality;
- broad support for many tools before Codex and Claude Code are excellent.

### 3.2 Kiro

Kiro provides a coherent end-user experience around three layers:

- **Specs** structure requirements, design, and tasks.
- **Steering** supplies persistent project knowledge with inclusion modes.
- **Hooks** react to events, inject context, run checks, or block dangerous operations.

Kiro also supports requirements-first and design-first workflows rather than assuming every change begins the same way. Its spec UI gives users visible progress and task tracking. See [Kiro Specs](https://kiro.dev/docs/specs/), [Steering](https://kiro.dev/docs/steering/), and [Hooks](https://kiro.dev/docs/hooks/).

The main lesson for Specd is that context has different lifetimes:

| Context | Lifetime | Specd equivalent |
| --- | --- | --- |
| Project conventions | Persistent | Knowledge/steering references |
| Change intent | Change lifetime | Change artifacts |
| Task contract | One task | Delegation packet |
| Execution events | One run | Attempt/result records |

Specd should not inject everything into every agent turn. It should resolve only the context needed for the current action.

### 3.3 Buzz

Buzz’s most relevant philosophy is: **the infrastructure is the pipe, not the brain**. Humans and agents provide intelligence; the platform provides shared state and communication. Buzz also uses protocol seams, structured events, personas, explicit access boundaries, and lifecycle hooks. Its agent boundary favors stdin/stdout and JSON, while MCP hooks can intercept lifecycle points without exposing every mechanism to the model. See [Buzz Vision](https://github.com/block/buzz/blob/main/VISION.md), [Buzz Agent](https://github.com/block/buzz/blob/main/crates/buzz-agent/README.md), and [MCP-Driven Hooks](https://github.com/block/buzz/blob/main/docs/MCP_DRIVEN_HOOKS.md).

For Specd, this implies:

- Specd owns workflow truth, not agent reasoning.
- The host agent remains Brain.
- Host-native subagent APIs launch Pinkies.
- Lifecycle hooks can refresh state or prevent Brain from ending with unresolved work.
- CLI and MCP should project the same operation model.
- A missing host capability must remain visible; Specd must not pretend to enforce what the host cannot enforce.

### 3.4 Combined lesson

The useful synthesis is:

| Source | Adopt | Avoid |
| --- | --- | --- |
| OpenSpec | Change folders, artifact graph, quick path, progressive rigor, generated skills | Prompt-only enforcement and command sprawl |
| Kiro | Persistent steering, visible phases, hooks, alternate planning paths | Lock-in to one IDE/runtime |
| Buzz | Protocol-first seams, lifecycle interception, personas, “pipe not brain” | Becoming a collaboration platform or agent runtime |

---

## 4. Proposed product philosophy

### 4.1 Product statement

> **Specd is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.**

### 4.2 Principles

1. **Path before enforcement.** First explain the state and next action; then enforce the small set of rules that makes the path trustworthy.
2. **Progressive rigor.** The workflow grows with risk and complexity. A small bug should not carry the ceremony of a cross-system migration.
3. **Agent-native, human-readable.** MCP/skills provide typed agent operations; Markdown and CLI views remain understandable without an agent.
4. **One canonical model.** CLI, MCP, generated skills, and HTML/text views project the same operation and state schemas.
5. **Recovery is a feature.** Every active state has a supported recovery path.
6. **The host executes.** Specd prepares and validates work; Codex or Claude Code performs reasoning and launches subagents.
7. **Context is selected, not accumulated.** Each role receives the smallest sufficient context package.
8. **Explicit delegation authority.** Brain coordinates; Pinky executes. Failure to delegate is a blocker, not permission for Brain to code.
9. **No invisible activation.** The user can always see whether Specd governs the request and why.
10. **Observed needs earn complexity.** Advanced assurance features return only after repeated real-world friction justifies them.

### 4.3 Non-goals for the reset

- Full process containment or sandboxing
- Generic multi-agent runtime
- Mandatory Git branching or commits
- Release management
- CI/CD orchestration
- Organization-wide project management
- Universal compatibility with every coding agent
- Cryptographic audit guarantees
- Concurrent writes by multiple Pinkies in the initial release

---

## 5. The new user workflow

### 5.1 Applicability: should this request use Specd?

The integration should classify the request before entering the workflow.

#### Use Specd when

- the user explicitly requests Specd;
- implementation has multiple dependent steps;
- the work should survive a new agent/session;
- more than one role or subagent is useful;
- requirements or acceptance criteria need agreement;
- the change affects several components or carries meaningful risk;
- the repository already has an active Specd change relevant to the request.

#### Do not automatically use Specd when

- answering or explaining;
- reading or reviewing without requested changes;
- brainstorming before the user chooses to formalize;
- trivial isolated edits;
- operational tasks unrelated to repository development;
- the user explicitly asks for direct execution.

#### Ambiguous request

The agent should ask once:

```text
This looks like a multi-step repository change. Use Specd to create a durable
plan and tracked execution path, or handle it directly?
```

This decision should be encoded in the generated Codex/Claude skill description, because activation descriptions are part of the tool contract.

### 5.2 The lifecycle

Use a small lifecycle with recoverable transitions:

```text
exploring -> planning -> ready -> executing -> reviewing -> done
                ^          |          |           |
                |          v          v           v
                +------ replanning <- blocked <- corrections
```

Definitions:

| Phase | Meaning | Owner |
| --- | --- | --- |
| `exploring` | The problem is still being understood; no execution commitment. | Human + Brain |
| `planning` | Templates are being filled and validated. | Brain, with human input |
| `ready` | The plan is accepted and at least one task is runnable. | Brain |
| `executing` | A Pinky owns an active task. | Pinky |
| `reviewing` | Results are being reconciled against acceptance criteria. | Brain, optional reviewer |
| `blocked` | Progress needs a named decision or repair. | Named owner |
| `replanning` | The plan is being changed without destroying the change. | Human + Brain |
| `done` | The change-level outcome is satisfied. | Human or configured policy |

Approval should be one explicit transition from `planning` to `ready`. It can be interactive but should not require terminal identity tricks in the default profile. The human’s action is meaningful; its authentication strength is a separate policy.

### 5.3 The one obvious command

Running `specd` with no arguments should render the home view. `specd next` should return the same navigation model in focused form.

Every result includes:

- current phase;
- recommended action;
- all legal alternatives;
- blocking reason and owner, if any;
- whether an agent may proceed automatically;
- exact structured invocation, when relevant.

Example:

```json
{
  "schema": "specd.workflow/v1",
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

---

## 6. Artifacts and OKF-compatible Markdown

Specd should generate Markdown templates that explain themselves to both humans and agents. The templates are not blank forms; they are executable guidance.

### 6.1 Relationship to OKF

Use OKF as the document organization profile for generated knowledge: stable identifiers, explicit purpose, typed relationships, concise sections, and traceable references. Specd remains the execution owner; it may consume project knowledge from Aido or other OKF sources, but should not become the knowledge system itself.

Because the exact OKF schema may evolve independently, Specd should implement an `okf-markdown/v1` profile behind a template interface rather than hard-code “OKF” assumptions throughout the parser.

### 6.2 Minimal change layout

```text
.specd/
  project.yaml
  changes/
    improve-payment-retry/
      intent.md
      requirements.md
      design.md          # created only when needed
      tasks/
        T1.md
        T2.md
      state.json         # machine-owned, replaceable projection
      runs/              # bounded worker/result records
```

Why task files instead of a single table:

- commands containing pipes no longer fight Markdown syntax;
- lists use native bullets rather than custom separators;
- task changes produce readable diffs;
- one Pinky receives one self-contained task artifact;
- task files can be validated and hashed independently;
- future tasks can carry different role schemas without a wide table.

### 6.3 Template anatomy

Each template should contain three layers:

1. **Machine metadata** in YAML front matter.
2. **Agent instructions** in comments or clearly delimited guidance blocks.
3. **Human-readable content** using stable headings.

Example task:

````markdown
---
okf: okf-markdown/v1
kind: execution-task
id: T4
status: pending
role: backend-engineer
requires:
  - REQ-payment-idempotency
depends_on:
  - T3
write_scope:
  - app/Payments/RetryService.php
  - tests/Feature/Payments/RetryServiceTest.php
---

# Add the retry idempotency guard

<!--
Agent instructions:
- State the observable outcome, not a vague activity.
- List only files this task is expected to modify.
- Reference requirement IDs rather than copying requirements.
- Provide one focused verification command.
- If required information is missing, mark Questions instead of guessing.
-->

## Objective

Prevent a retried payment request from producing a second charge.

## Context

- [Payment retry requirement](../requirements.md#req-payment-idempotency)
- [Retry design](../design.md#decision-idempotency-key)

## Constraints

- Preserve the existing payment gateway adapter.
- Do not alter successful first-attempt behavior.

## Verification

```command
php artisan test --filter=RetryServiceTest
```

## Done when

- Replaying the same idempotency key returns the original result.
- No second gateway charge is created.

## Questions

None.
````

### 6.4 Progressive artifact sets

Not every change needs every artifact.

| Mode | Generated artifacts | Use case |
| --- | --- | --- |
| **Quick** | `intent.md`, task files | Small, understood change |
| **Standard** | intent, requirements, design, task files | Multi-step feature or bug |
| **High-risk** | standard + risks, rollout, migration/recovery sections | Contracts, data migration, security, cross-system changes |

The mode is chosen during planning and can be raised later. It should never be silently lowered.

---

## 7. Pinky and the Brain, redesigned

The old idea remains valuable, but it should be a host-independent orchestration protocol rather than embedded orchestration logic that assumes a specific agent runtime.

### 7.1 Responsibilities

#### Brain: host agent

Brain may:

- understand user intent;
- ask clarifying questions;
- select or create a Specd change;
- request templates and instructions;
- help draft planning artifacts;
- call `specd next`;
- select from the ready frontier;
- launch a Pinky through the host’s native subagent mechanism;
- submit the Pinky result to Specd;
- explain state and blockers to the user.

Brain may not modify implementation files for a task designated `delegated`, unless the user explicitly changes the execution policy.

#### Pinky: task worker

Pinky may:

- read its delegation packet;
- inspect task-relevant repository context;
- modify only the declared write scope when the host supports restriction;
- run allowed verification commands;
- return a structured result.

Pinky may not:

- change the plan;
- mark its own task accepted;
- start another task;
- widen its scope;
- approve the change;
- silently delegate again in the initial version.

#### Specd: coordinator and contract owner

Specd may:

- derive task readiness and waves;
- generate bounded delegation packets;
- validate configuration before delegation;
- record leases and worker results;
- reconcile results into state;
- identify the next action;
- block illegal transitions;
- explain recovery.

Specd does not launch models directly in the initial version.

### 7.2 Delegation packet

`specd task packet T4 --json` should return everything a Pinky requires and nothing unrelated:

```json
{
  "schema": "specd.task/v1",
  "change": "improve-payment-retry",
  "task": "T4",
  "role": {
    "id": "backend-engineer",
    "mission": "Implement backend behavior with focused tests",
    "constraints": ["Do not redesign unrelated modules"]
  },
  "objective": "Prevent duplicate charge creation during retries",
  "requirement_refs": ["REQ-payment-idempotency"],
  "read_context": ["...selected excerpts or paths..."],
  "write_scope": ["app/Payments/RetryService.php", "tests/Feature/Payments/RetryServiceTest.php"],
  "verify": ["php artisan test --filter=RetryServiceTest"],
  "result_schema": "specd.task-result/v1"
}
```

### 7.3 Result contract

Pinky returns:

```json
{
  "task": "T4",
  "outcome": "completed",
  "summary": "Added idempotency lookup before gateway invocation.",
  "changed_files": ["app/Payments/RetryService.php", "tests/Feature/Payments/RetryServiceTest.php"],
  "verification": [{"command": "php artisan test --filter=RetryServiceTest", "exit_code": 0}],
  "discoveries": [],
  "questions": [],
  "risks": []
}
```

Specd validates the shape and task scope. Brain judges semantic adequacy, optionally with a reviewer Pinky. In the simple profile, successful declared verification plus Brain acceptance completes the task.

### 7.4 Waves and parallelism

Waves remain a useful projection of dependencies, but the first reset release should execute one code-writing Pinky at a time. This avoids merge coordination before the protocol is stable.

Specd can expose the full ready frontier and recommend parallel-safe work, but parallel mutation should remain disabled unless:

- write scopes are disjoint;
- the host supports isolated worktrees or sandboxes;
- result reconciliation is defined;
- the user enables it.

### 7.5 Delegation preflight

Before Brain launches a Pinky, `specd delegate preflight` validates:

- role exists;
- required template fields are complete;
- host declares a compatible subagent capability;
- task packet fits the configured context budget;
- write scope is present and safe;
- verification is defined;
- no conflicting active lease exists.

Failure returns a single blocker with repair steps. Brain must not implement the task as an implicit fallback.

---

## 8. Native Codex and Claude Code integration

MCP alone is insufficient. A tool can exist while the agent fails to activate it, selects the wrong operation, or guesses arguments. Native integration needs three coordinated surfaces.

### 8.1 Generated skill/command layer

`specd init` should detect or accept explicit targets:

```text
specd init --agent codex
specd init --agent claude-code
specd init --agent all
```

It generates:

- a focused Specd workflow skill for Codex;
- Claude Code slash commands or skills;
- a small managed block in `AGENTS.md` and/or `CLAUDE.md`;
- MCP configuration guidance or local server configuration;
- a compatibility manifest recording generated integration versions.

The skill description must precisely state when to activate and when not to activate. Detailed commands live in references loaded on demand.

### 8.2 Typed MCP layer

Do not expose a generic `run_specd_command(command: string)` tool. Expose intent-level operations:

- `specd_orient`
- `specd_assess_request`
- `specd_create_change`
- `specd_get_template`
- `specd_validate_plan`
- `specd_get_next`
- `specd_prepare_delegation`
- `specd_submit_result`
- `specd_recover`

Each tool should have a narrow JSON schema, descriptions, enum choices, and examples. Change/task selectors should be optional when Specd can resolve them unambiguously; ambiguity should return candidates.

Opaque revisions and attempt IDs should not be user/model-supplied unless unavoidable. Prefer a transition token returned by one tool and consumed automatically by the next integration call.

### 8.3 CLI parity

CLI and MCP are adapters over the same application operations:

```text
CLI command ----\
                > operation service -> canonical result
MCP tool -------/
```

The operation service owns validation and result schemas. No MCP-specific business logic and no independently maintained prompt copy.

### 8.4 Lifecycle hooks

Use supported host hooks where they create deterministic value:

- **Session start:** inject active change and current phase.
- **Before task delegation:** run preflight.
- **After subagent result:** require result submission/reconciliation.
- **Before Brain edits code:** warn or block when an active delegated task exists.
- **Before session end:** surface unresolved active work and the resume command.

Hooks must remain thin. They call Specd and relay its result; they do not duplicate state logic.

### 8.5 Cold-start test

Every integration release should pass this test with a new agent and no chat history:

1. Detect whether the current request needs Specd.
2. Find the active change or explain that none exists.
3. Report the current phase and options.
4. Identify the next ready task.
5. Prepare the correct delegation packet.
6. Stop rather than code when delegation configuration is broken.
7. Resume correctly after a new session.

Codex and Claude Code need separate end-to-end fixtures because their skill, hook, and subagent surfaces differ.

---

## 9. Context engineering model

### 9.1 Context layers

Use four explicit layers:

| Layer | Examples | Inclusion |
| --- | --- | --- |
| **Foundation** | architecture, conventions, glossary | Selected by relevance; never copied wholesale |
| **Change** | intent, requirements, decisions | Included for planning and review |
| **Task** | objective, referenced requirements, dependencies, scope | Always in Pinky packet |
| **Runtime** | diff, test output, discoveries, blocker | Created during execution and summarized on return |

### 9.2 Context assembly rules

- References are preferred over duplicated prose.
- Required context is never silently truncated.
- Optional context may be omitted with visible reasons.
- Context budgets are per role and task, not global.
- Pinky receives no unrelated tasks.
- Brain receives summaries and links for completed Pinky work, not entire worker transcripts.
- A task that cannot fit its required context should be split or explicitly granted a larger budget.

### 9.3 Persona design

A role/persona should be a small contract, not theatrical prose:

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

The task carries domain-specific instructions. The role carries stable behavior. Project knowledge remains separately referenced. This prevents one enormous prompt from mixing persona, project rules, and task requirements.

### 9.4 Discoveries and plan drift

Pinkies will discover incorrect assumptions. They need a legitimate response other than guessing or editing the plan.

Result outcomes should include:

- `completed`
- `failed`
- `blocked`
- `needs_replan`

`needs_replan` returns discoveries and suggested impact. Brain transitions the change to replanning, preserving completed task history while recalculating affected tasks. Human approval is requested only if the change materially alters intent, acceptance criteria, risk, or scope.

---

## 10. Recovery model

### 10.1 State should be reconstructable

Authored artifacts are durable truth. `state.json` is a convenient projection and should be reconstructable from task metadata and bounded run records. A malformed projection should be repairable with:

```text
specd doctor --repair
```

Repair should never rewrite authored intent without confirmation.

### 10.2 Supported recovery operations

| Situation | Recovery |
| --- | --- |
| Worker crashed | Expire/release task lease; retry or resume |
| Wrong task plan | `replan` while preserving prior results |
| Missing role/config | `doctor` reports exact configuration repair |
| Test failure | Task remains active; fix and retry |
| Scope must widen | Propose task amendment; Brain/human accepts; regenerate packet |
| Agent edited outside scope | Report discrepancy; choose keep/revert/replan explicitly |
| State projection malformed | Rebuild from artifacts and run records |
| User no longer wants change | Cancel with reason; keep artifacts |
| Wrong change selected | Switch explicitly; never infer across ambiguous active changes |

### 10.3 No dead-end invariant

For every non-terminal state, automated tests must prove at least one documented legal transition. The test is not merely “a refusal has a next string”; it executes the recovery journey.

### 10.4 Destructive recovery

Deleting a change remains available only as an explicit destructive command with confirmation. It is not troubleshooting advice. Normal repair preserves the change identity and authored documents.

---

## 11. Recommended command and tool model

Keep the human CLI small:

| Command | Purpose |
| --- | --- |
| `specd` | Home/orientation view |
| `specd init` | Initialize repository and selected agent integrations |
| `specd new` | Start guided change creation |
| `specd status` | Inspect current or selected change |
| `specd next` | Return recommended and alternative actions |
| `specd check` | Validate authored artifacts and configuration |
| `specd approve` | Accept a plan for execution |
| `specd pause` | Stop automatic progression safely |
| `specd replan` | Re-enter planning without destroying progress |
| `specd resume` | Resume a paused or repaired change |
| `specd cancel` | End a change while preserving its record |
| `specd doctor` | Diagnose and repair setup/state |

Task start, packets, leases, result submission, and reconciliation should primarily be MCP/application operations used by the host integration. Advanced CLI subcommands may exist for debugging, but should not dominate normal user documentation.

---

## 12. Architecture direction

### 12.1 Layers

```text
Generated Codex skill / Claude command / Human CLI
                       |
                  MCP + CLI adapters
                       |
               Application operations
                       |
     workflow | artifacts | tasks | delegation | recovery
                       |
          filesystem + optional Git observations
```

### 12.2 Canonical models

Keep a small set of versioned contracts:

- `specd.workflow/v1`
- `specd.artifact-instructions/v1`
- `specd.task/v1`
- `specd.task-result/v1`
- `specd.integration/v1`

Avoid creating a registry for every internal type. Version only cross-boundary contracts.

### 12.3 Determinism boundary

Deterministic:

- parsing artifacts;
- validating schemas and references;
- deriving artifact readiness;
- deriving task waves/frontier;
- selecting legal transitions;
- creating context manifests;
- validating worker results;
- recovery diagnosis;
- rendering CLI/MCP results.

Non-deterministic and host-owned:

- interpreting user intent;
- drafting artifact content;
- choosing between semantically valid designs;
- coding;
- semantic review;
- deciding whether newly discovered information changes intent.

This is the clean meaning of “the agent reasons; the harness coordinates and enforces.”

---

## 13. Initial delivery plan

### Phase 0 - Preserve and learn

**Goal:** capture the old system before redesign.

Actions:

- Tag or archive the current `specd-cli` implementation.
- Recover the Pinky/Brain package if a commit, bundle, or local archive can be located.
- Convert the user-reported failures into journey fixtures.
- Write a one-page product charter and explicit non-goals.

Exit criteria:

- Every reported problem has a reproducible narrative.
- Current behavior is preserved for reference.
- No new implementation begins before the charter is accepted.

**Repository note:** the supplied `0xkhdr/specd` repository currently clones as an empty repository and public code search returned no Pinky/Brain implementation. The plan therefore treats the orchestration behavior described by the user as authoritative, but implementation recovery requires a commit SHA, Git bundle, local clone, release archive, or another surviving ref.

### Phase 1 - Workflow kernel and orientation

**Goal:** prove the simple lifecycle and make the product understandable without an agent.

Build:

- project/change discovery;
- minimal lifecycle and transition table;
- `specd`, `status`, `next`, `pause`, `resume`, `replan`, `cancel`;
- canonical JSON result;
- `doctor` with read-only diagnosis;
- no-dead-end transition tests.

Exit criteria:

- A new user can identify the current stage and options in under one minute.
- Every reachable state has an executed recovery journey.
- No Git commit is required for normal transitions.

### Phase 2 - Template engine and OKF profile

**Goal:** generate useful self-describing artifacts.

Build:

- template interface;
- `okf-markdown/v1` profile;
- Quick, Standard, and High-risk artifact modes;
- template instruction projection for humans and agents;
- task-per-file format;
- artifact dependency graph and validation;
- `new`, `check`, and `approve`.

Exit criteria:

- A cold agent can fill each template without external command knowledge.
- Validation catches missing questions, unverifiable acceptance, broken references, and missing verification.
- A small change can use Quick mode without empty ceremony.

### Phase 3 - Codex-native vertical slice

**Goal:** one excellent native integration.

Build:

- generated Codex skill;
- exact activation/non-activation rules;
- typed MCP tools;
- session-start orientation if supported;
- one complete flow from request assessment through approved plan;
- integration compatibility/version check.

Exit criteria:

- Codex never constructs raw command strings during the tested journey.
- It uses Specd for qualifying work and avoids it for negative examples.
- Cold-start and session-resume tests pass.

### Phase 4 - Claude Code parity

**Goal:** equivalent workflow behavior through Claude Code’s native surfaces.

Build:

- generated Claude skill/commands;
- MCP configuration;
- relevant lifecycle hooks;
- parity suite against the canonical operations.

Exit criteria:

- Codex and Claude Code produce equivalent state transitions and task packets.
- Tool-specific generated files differ only where invocation mechanisms require it.

### Phase 5 - Pinky and the Brain vertical slice

**Goal:** safely delegate one task.

Build:

- roles/personas;
- task packet generation;
- delegation preflight;
- task lease;
- structured result submission;
- Brain code-edit guard while a delegated task is active;
- explicit blocker when no Pinky can launch;
- sequential wave execution.

Exit criteria:

- Brain delegates a task and never implements it itself under delegation policy.
- Missing configuration stops before code mutation and explains repair.
- Pinky receives only task-relevant context.
- Result reconciliation determines the next task.

### Phase 6 - Recovery and replanning

**Goal:** eliminate delete-and-restart recovery.

Build:

- `doctor --repair`;
- state reconstruction;
- expired worker recovery;
- plan amendment impact analysis;
- replan while preserving unaffected completed tasks;
- cancel/archive without deletion.

Exit criteria:

- Every historical dead-end scenario has a supported recovery journey.
- Corrupt projection recovery does not alter authored content.

### Phase 7 - Carefully earned concurrency

**Goal:** parallel Pinkies only after sequential orchestration is trustworthy.

Build only if justified:

- disjoint-scope analysis;
- host worktree/sandbox capability negotiation;
- parallel leases;
- deterministic result merge/review;
- conflict recovery.

Exit criteria:

- Parallel mode cannot start without isolation capability.
- Conflicts are surfaced as reconcile/replan states, never silently merged.

---

## 14. Testing strategy

### 14.1 Product journeys

Tests should tell user stories, not only protect internal abstractions:

1. New user asks “how do I start?”
2. Trivial request correctly avoids Specd.
3. Multi-step change recommends Specd and waits for consent.
4. Agent creates and fills Standard templates.
5. User sees phase, progress, and options.
6. Codex selects typed operations without guessing arguments.
7. Claude Code resumes after a fresh session.
8. Brain delegates to a correctly configured Pinky.
9. Missing Pinky configuration blocks Brain from coding.
10. Pinky reports a discovery requiring replanning.
11. Failed verification remains recoverable.
12. Corrupted state is rebuilt.
13. User pauses, resumes, and cancels.
14. Plan changes preserve unaffected completion.

### 14.2 Activation evaluation

Maintain an evaluation set of user prompts labeled:

- must use Specd;
- should offer Specd;
- must not use Specd.

Measure false activation and missed activation separately. Both matter.

### 14.3 Context evaluation

For every task packet, test:

- required facts are present;
- irrelevant task content is absent;
- references resolve;
- write scope is complete;
- role instructions are concise;
- packet stays within its budget;
- omitted optional context is reported.

### 14.4 Recovery evaluation

Inject failures after every state mutation and confirm that `status` and `doctor` produce a supported path. Focus first on user-observed failures, not theoretical exhaustive durability.

---

## 15. Decisions to make before implementation

The following decisions materially shape the design and should be discussed before code:

1. **OKF contract:** Which exact OKF version/grammar should Specd target, and is OKF authoritative or one selectable template profile?
2. **Approval semantics:** Is default approval a simple recorded user decision, with stronger identity only in an opt-in policy?
3. **Completion authority:** Does Brain accept a successful Pinky result automatically, or must a separate review step exist for Standard mode?
4. **Host capability:** Which Codex and Claude Code subagent APIs are considered the initial supported baseline?
5. **Direct execution:** Can the user choose `execution: brain` for small tasks, while delegated tasks forbid Brain edits?
6. **Knowledge boundary:** Which project knowledge is owned by Specd versus referenced from Aido/OKF documents?
7. **Compatibility:** Is the reset a new major version with no state migration, or must existing `.specd` roots be importable?

Recommended defaults:

- OKF is a versioned template profile.
- Approval is simple by default, stronger identity optional.
- Brain accepts task results after verification in Standard mode; high-risk mode adds an independent reviewer.
- Support Codex and Claude Code only until both pass the same journeys.
- Allow explicit `brain` execution for Quick mode; delegated tasks never silently fall back.
- Aido owns broad knowledge; Specd owns change and execution artifacts.
- Ship as a clean major-version reset with a read-only importer, not an in-place mutation.

---

## 16. Success measures

Specd succeeds when it becomes boring and dependable.

Measure:

- time for a new user to start the first change;
- percentage of sessions where the agent selects the correct next operation first try;
- false and missed activation rates;
- number of delete-and-restart recoveries;
- percentage of delegation failures caught before code mutation;
- average required context per Pinky;
- number of times Brain edits code despite delegated policy;
- number of user questions answered by the home/status view without documentation;
- task result acceptance and replanning rates;
- workflow completion across fresh sessions.

Initial product targets:

- zero guessed CLI argument journeys in native integrations;
- zero states whose documented recovery is deletion;
- zero implicit Brain fallback after delegation failure;
- under one minute to orient a new user or cold agent;
- under five primary commands needed for an ordinary user journey;
- Quick mode creates no unused mandatory artifact.

---

## 17. Immediate next steps

1. Accept or revise the product statement and non-goals.
2. Resolve the exact OKF profile and Aido/Specd ownership boundary.
3. Recover the original Pinky/Brain source through a commit SHA or archive, if available.
4. Write the lifecycle transition table and no-dead-end recovery matrix.
5. Design the human home view and canonical `specd.workflow/v1` JSON together.
6. Design one OKF task template and one task packet/result pair.
7. Build a mocked Codex journey before implementing orchestration.
8. Implement the smallest vertical slice: orient -> plan -> approve -> next.
9. Add one delegated task only after the single-agent path is understandable.

Do not begin by porting the current state/evidence machinery. Begin with the user journey, typed contracts, and recovery invariants. Reuse current code only where it directly serves the simplified model.

---

## Conclusion

Specd’s original idea remains strong. The mistake was not pursuing determinism; it was allowing determinism to become synonymous with maximum enforcement.

The reset should make the development path explicit before making it strict. A human should always understand the stage and options. A coding agent should never guess commands. Brain should coordinate without stealing Pinky’s work. Pinky should receive a narrow, sufficient context package. Broken configurations and changed plans should lead to visible recovery states, not deletion.

The resulting philosophy is concise:

> **The agent reasons. Specd makes the path visible, packages the work, validates transitions, and preserves recovery.**

That is small enough to build, distinct enough to matter, and extensible enough to regain advanced assurance later when real experience earns it.

---

## References

- [Current Specd CLI](https://github.com/0xkhdr/specd-cli)
- [Archived Specd repository](https://github.com/0xkhdr/specd)
- [OpenSpec](https://github.com/Fission-AI/OpenSpec)
- [OpenSpec concepts](https://github.com/Fission-AI/OpenSpec/blob/main/docs/concepts.md)
- [OpenSpec agent contract](https://github.com/Fission-AI/OpenSpec/blob/main/docs/agent-contract.md)
- [Kiro specs](https://kiro.dev/docs/specs/)
- [Kiro spec best practices](https://kiro.dev/docs/specs/best-practices/)
- [Kiro steering](https://kiro.dev/docs/steering/)
- [Kiro hooks](https://kiro.dev/docs/hooks/)
- [Kiro agent skills](https://kiro.dev/docs/skills/)
- [Buzz](https://github.com/block/buzz)
- [Buzz vision](https://github.com/block/buzz/blob/main/VISION.md)
- [Buzz agent](https://github.com/block/buzz/blob/main/crates/buzz-agent/README.md)
- [Buzz MCP-driven lifecycle hooks](https://github.com/block/buzz/blob/main/docs/MCP_DRIVEN_HOOKS.md)

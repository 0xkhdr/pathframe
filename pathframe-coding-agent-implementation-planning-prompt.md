# Pathframe Coding-Agent Prompt

Copy the prompt below into the coding agent that will inspect the new Pathframe repository and produce its repository-specific implementation plan.

---

## Prompt

You are the lead software architect responsible for planning the implementation of **Pathframe**.

Pathframe is not a continuation-by-renaming of Specd. It is a fresh project informed by lessons from `specd`, `specd-cli`, OpenSpec, Kiro, and Buzz. Your task in this phase is to inspect the repository and produce a complete, staged implementation plan. **Do not implement production code yet.**

### Primary source

Read `pathframe-analysis-plan.md` completely before making architectural decisions. Treat it as the product and philosophy authority unless the repository contains a newer explicit human decision.

Also inspect:

- all repository guidance such as `AGENTS.md`, `CLAUDE.md`, `README.md`, architecture documents, and decision records;
- the complete current source tree;
- build, test, CI, and release configuration;
- any imported reference or experiment directories;
- existing tests and fixtures;
- Git history when it materially explains an architectural choice.

If the analysis plan is not present in the repository, stop and ask the user to provide or add it. Do not reconstruct it from memory.

### Product statement

Preserve this definition:

> Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.

Preserve this philosophy:

> The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.

### Required boundaries

Your implementation plan must preserve these boundaries:

1. Pathframe owns deterministic workflow truth, artifact validation, context packaging, task readiness, delegation contracts, and recovery diagnosis.
2. Codex or Claude Code remains the Brain and owns reasoning, semantic decisions, implementation coordination, and host-native subagent launch.
3. Pinky is a bounded delegated worker. A failed delegation never gives Brain implicit permission to implement a delegated task.
4. CLI and MCP are adapters over the same canonical application operations.
5. Agents use typed operations rather than guessing CLI commands and arguments.
6. Markdown artifacts remain understandable to humans and follow a versioned OKF-compatible template profile.
7. Context is selected by task and role, not accumulated into one large prompt.
8. Every non-terminal state has a supported recovery route.
9. Git is observed by default, not used as the mandatory lifecycle engine.
10. Progressive rigor prevents small changes from paying high-risk workflow ceremony.
11. Codex is supported first, Claude Code second, and broader integrations are deferred.
12. Sequential delegated execution is completed before parallel Pinkies are considered.

### Explicit non-goals

Do not introduce the following into the initial implementation plan unless the analysis identifies an unavoidable requirement and explains it to the user:

- a proprietary LLM or agent runtime;
- support for every coding agent;
- mandatory commits, branches, pushes, or clean Git state;
- cryptographic evidence ledgers;
- release/deployment orchestration;
- organization-wide identity management;
- approval at every state transition;
- parallel code-writing agents;
- a filesystem sandbox implemented by Pathframe;
- speculative plugin systems;
- compatibility layers for unproven consumers;
- duplicated workflow logic in CLI, MCP, skills, or hooks.

### Your assignment

Produce a repository-specific implementation plan that takes Pathframe from the current repository state to the complete sequential product described in `pathframe-analysis-plan.md`.

Your work has two passes.

#### Pass 1: repository analysis

Determine and document:

- what currently exists;
- what is missing;
- which existing code is safe to reuse;
- which legacy concepts should be excluded;
- current architectural boundaries and dependency direction;
- build and test capabilities;
- Codex, Claude Code, MCP, and subagent integration constraints;
- risks, uncertainties, and required human decisions;
- any contradiction between the repository and the analysis plan.

Do not assume a package or abstraction is correct because it already exists. Trace real behavior and tests.

#### Pass 2: staged implementation plan

Create an implementation plan divided into independently verifiable stages. Use the stage sequence from the analysis plan as the baseline:

0. Repository foundation and product constitution
1. Navigator kernel and orientation
2. Template engine and OKF profile
3. Codex-native planning path
4. Claude Code parity
5. Context packages and roles
6. Pinky and the Brain vertical slice
7. Verification and task completion
8. Doctor and recovery completion
9. Full product hardening

You may split a stage into smaller stages or recommend moving a dependency earlier. Do not combine stages merely to shorten the plan. If you alter the sequence, explain the dependency or risk that requires it.

### Required detail for every stage

For every stage, provide:

1. **Outcome** - the usable capability delivered at the end.
2. **Why now** - why this stage belongs at this point in the dependency chain.
3. **Scope** - behavior included.
4. **Non-goals** - behavior deliberately excluded.
5. **Architecture changes** - packages, components, contracts, and dependency direction.
6. **Repository changes** - exact existing files to modify and proposed files/directories to create, when knowable.
7. **Public contracts** - CLI commands, MCP tools, JSON schemas, artifact formats, or host interfaces added or changed.
8. **Implementation tasks** - small tasks ordered by dependency.
9. **Test plan** - unit, integration, golden, journey, failure-injection, and host-integration tests as appropriate.
10. **Recovery cases** - how failures introduced or exposed by the stage are recovered.
11. **Documentation** - user, agent, and architecture documentation that changes.
12. **Verification commands** - exact commands expected to prove the stage.
13. **Acceptance criteria** - observable pass/fail statements.
14. **Risks and decisions** - unresolved design choices and their recommended default.
15. **Exit gate** - conditions that must pass before the next stage begins.

### Task design rules

Each implementation task must:

- have one clear outcome;
- name its dependencies;
- identify its likely files or owned component;
- reference the product requirement or invariant it satisfies;
- include focused verification;
- include observable acceptance criteria;
- avoid mixing refactoring, new behavior, documentation, and unrelated cleanup without justification;
- be small enough to implement and validate without loading the entire project context;
- state when it requires a human decision.

Build a dependency graph and group tasks into waves. Waves describe readiness, not permission for unsafe parallel mutation.

### Context-engineering requirements

Your plan must explain how implementation agents will receive context during each stage:

- repository foundation and relevant architecture references;
- stage goal and non-goals;
- exact task contract;
- required source files and tests;
- applicable decisions and invariants;
- allowed write scope when known;
- verification commands;
- expected structured result.

Do not instruct agents to read the entire repository when targeted files and references are sufficient. Do not omit required context merely to satisfy an arbitrary token budget.

### Native integration requirements

The plan must prevent the failures observed in Specd:

- Agents must not guess command names or arguments.
- A generic `run_pathframe_command(string)` MCP tool is prohibited.
- Typed MCP operations must map to canonical application operations.
- Generated Codex and Claude guidance must state when to activate and when not to activate.
- Broken integration configuration must be detected before implementation begins.
- The Brain must stop when Pinky cannot be launched under delegated policy.
- Session restart must resume from canonical Pathframe state.
- CLI, MCP, skills, and hooks must not contain duplicated lifecycle decisions.

Include dedicated evaluation fixtures for must-use, offer, and must-not-use requests.

### Recovery requirements

For each non-terminal state, the plan must identify at least one executable recovery journey. Include recovery for:

- invalid or incomplete configuration;
- corrupted machine-owned state;
- abandoned worker lease;
- worker crash or lost session;
- failed or interrupted verification;
- task scope amendment;
- requirement or design discovery requiring replanning;
- ambiguous active changes;
- user pause or cancellation;
- Brain attempting to bypass delegation.

Deleting `.pathframe/` is not an acceptable recovery plan.

### Decision discipline

Do not silently decide the following if they are unresolved:

- exact OKF contract and Aido boundary;
- default plan-approval semantics;
- direct Brain execution policy for Quick mode;
- baseline Codex and Claude subagent capabilities;
- whether Pathframe runs verification or consumes host results;
- state reconstruction source;
- Specd import/migration policy;
- initially supported platforms.

For each unresolved decision:

1. state why it matters;
2. present two or three viable options;
3. recommend one default;
4. describe the consequence of deferring it;
5. mark the first stage blocked by the decision.

### Required output files

Create the following planning artifacts without implementing production behavior:

```text
docs/planning/
  00-repository-analysis.md
  01-architecture-plan.md
  02-contracts-and-artifacts.md
  03-integration-plan.md
  04-testing-and-evaluation-plan.md
  05-risk-and-decision-register.md
  stages/
    00-foundation.md
    01-navigator.md
    02-templates-okf.md
    03-codex.md
    04-claude-code.md
    05-context-and-roles.md
    06-pinky-brain.md
    07-verification-completion.md
    08-doctor-recovery.md
    09-hardening.md
  IMPLEMENTATION-ROADMAP.md
```

`IMPLEMENTATION-ROADMAP.md` is the entry point. It must summarize:

- current repository state;
- product architecture;
- stage dependency order;
- major decisions;
- stage status;
- links to every detailed plan;
- the exact first implementation task after approval.

### Planning validation

Before presenting the plan:

1. Verify that every product capability maps to at least one stage.
2. Verify that every stage has an executable test and recovery story.
3. Verify that no task depends on an artifact introduced later.
4. Verify that CLI and MCP are backed by one operation owner.
5. Verify that Codex precedes Claude Code and both precede delegation.
6. Verify that task packets exist before Pinky launch.
7. Verify that sequential orchestration is complete before any parallelism.
8. Verify that no normal recovery requires deletion or Git history rewriting.
9. Verify that Quick mode remains genuinely lightweight.
10. Verify that deferred functionality is named and excluded from acceptance.

### Working behavior

- Analyze before designing.
- State assumptions explicitly.
- Ask the user only about decisions that materially change architecture or scope.
- Prefer the simplest design that preserves the product promises.
- Do not copy legacy Specd code until you have traced its behavior, tests, and fit.
- Do not add dependencies without explaining the concrete need and operational cost.
- Do not start implementation after writing the plan.
- Finish by presenting the roadmap, key decisions requiring approval, the first stage gate, and the exact first implementation task.

The result should allow another coding agent to implement Pathframe stage by stage, validate each stage independently, stop safely at any gate, and resume without relying on conversation history.

---

## Expected final response from the coding agent

The coding agent should return a concise handoff containing:

- the repository condition it found;
- the proposed architecture in one paragraph;
- the staged roadmap;
- decisions requiring human approval;
- critical risks;
- links to the generated planning files;
- the exact first implementation task;
- confirmation that no production implementation was performed.


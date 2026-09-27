# Pathframe Philosophy

> Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.

> The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.

## Principles

### Show the path before enforcing it

Every refusal must state the current state, explain why the requested transition is illegal, and provide at least one supported recovery action.

### Keep reasoning in the agent

Pathframe owns deterministic workflow mechanics. Brain owns interpretation, design, semantic decisions, implementation coordination, and review. Pathframe must not grow into an LLM or proprietary agent runtime.

### Use progressive rigor

Quick changes require intent and task contracts. Standard changes add requirements and design. High-risk changes add risks, rollout, and recovery. Risk may raise rigor but never silently lower it.

### Keep artifacts human-readable and operations typed

Markdown artifacts follow the Pathframe-owned `okf-markdown/v1` profile. CLI and MCP expose the same canonical operations. Agents use typed operations instead of guessing CLI commands or argument order.

### Select context instead of accumulating it

Task packets contain only required foundation, change, task, and runtime context. Required context is never silently truncated. Optional omissions are visible.

### Make delegation authority explicit

Every task declares `execution_policy`. Quick tasks default to `brain`. A delegated task belongs to Pinky for implementation. If sequential host-native delegation fails, Brain reports a blocker and does not inherit permission to code.

### Be honest about host capabilities

Shared-workspace write scope is advisory unless the host explicitly declares enforcement. Pathframe must not describe advisory scope as containment or assume worktrees, sandboxes, cancellation, resumption, hooks, or parallel workers.

### Preserve recovery

Authored artifacts define intended work. Append-only `history.jsonl` transitions and bounded run records preserve lifecycle facts. `state.json` is a replaceable projection. Git and conversation history are not reconstruction authorities. Deleting `.pathframe/` is never normal recovery.

### Let complexity earn its place

Do not copy a legacy subsystem because it exists. Trace its behavior, tests, license, and fit first. Prefer standard library and narrow internal contracts. Version only persisted or cross-process boundaries.

## Initial non-goals

- A proprietary LLM or agent runtime
- Support for every coding agent
- Mandatory commits, branches, pushes, or clean Git state
- Cryptographic evidence or audit guarantees
- Release and deployment orchestration
- Full project management or organization-wide identity controls
- Approval at every routine transition
- Parallel code-writing Pinkies
- A Pathframe-built filesystem sandbox
- Streamable HTTP or custom JSON-RPC
- Speculative plugin or compatibility frameworks
- Specd import, compatibility reading, or automatic `.specd/` detection

Parallel workers, strict Git evidence, broader hosts, and a possible one-time read-only-source Specd importer remain future decisions. They are not part of initial product acceptance.

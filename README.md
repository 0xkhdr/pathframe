# Pathframe

> Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.

> The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.

Pathframe will give humans and coding agents one local view of current work: where a change is, what is ready, what is blocked, what comes next, and how to recover. It will package approved tasks for direct or delegated execution without becoming an agent runtime.

## Project status

Pathframe Stage 5 is implemented and awaiting its human stop/go gate. It adds role contracts, bounded four-layer context resolution, visible budget omissions, inspectable `pathframe.task/v1` packets, and deterministic task frontier/wave projection over the approved planning path shared by Codex and Claude Code.

The first supported production platform will be Linux amd64 with Go 1.26. Linux arm64, macOS amd64/arm64, and Windows amd64 remain portability targets until Stage 9 proves them through CI and end-to-end journeys.

## Intended first use

Running `pathframe` without arguments provides project and active-change orientation instead of bare usage. See the [CLI reference](docs/CLI.md), [artifact profile](docs/ARTIFACTS.md), [Codex integration](docs/CODEX.md), and [Claude Code integration](docs/CLAUDE-CODE.md). Both hosts share typed planning and packet-preview operations. Later stages add sequential worker launch, verification execution, and broader Doctor recovery.

## Boundaries

- Pathframe owns deterministic workflow truth, artifact validation, context packaging, task readiness, delegation contracts, verification records, and recovery diagnosis.
- Codex or Claude Code remains the Brain and owns reasoning, semantic decisions, coordination, and host-native subagent launch.
- Pinky is a bounded delegated worker. Failed delegation never authorizes Brain to implement a delegated task.
- CLI and MCP are adapters over one canonical application layer.
- Markdown remains human-readable; agents use typed operations.
- Git is observed when useful but is not the lifecycle engine.
- Every supported non-terminal state must have an executable recovery route.

## Initial non-goals

Pathframe will not initially provide a proprietary LLM runtime, universal agent support, mandatory Git branches or commits, cryptographic evidence, release orchestration, organization identity management, a custom filesystem sandbox, parallel code-writing workers, speculative plugins, or Specd compatibility.

## Planning and contribution

Start with [the implementation roadmap](docs/planning/IMPLEMENTATION-ROADMAP.md). Architecture and product constraints are summarized in [ARCHITECTURE.md](ARCHITECTURE.md) and [PHILOSOPHY.md](PHILOSOPHY.md). See [CONTRIBUTING.md](CONTRIBUTING.md) and [AGENTS.md](AGENTS.md) before changing the repository. Record any intentional deviation from an accepted decision using [the decision process](docs/decisions/README.md).

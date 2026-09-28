# Pathframe

> Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.

> The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.

Pathframe gives humans and coding agents one local view of current work: where a change is, what is ready, what is blocked, what comes next, and how to recover. It packages approved tasks for direct or delegated execution without becoming an agent runtime.

## Project status

Pathframe Stages 0 through 9 are implemented and approved. The sequential product includes bounded delegation and verification, Doctor recovery, Linux amd64 installation and hardening evidence, and explicit compatibility and security limits.

The supported production platform is Linux amd64 with Go 1.26. Linux arm64, macOS amd64/arm64, and Windows amd64 remain portability targets with cross-build evidence only.

## Documentation

Running `pathframe` without arguments provides project and active-change orientation instead of bare usage.

| Need | Read |
| --- | --- |
| Understand Pathframe | [Documentation guide](docs/index.md), then [Core concepts](docs/concepts.md) |
| Install and try Pathframe | [Installation](docs/installation.md), then [Getting started](docs/getting-started.md) |
| Use the product | [CLI reference](docs/cli-reference.md), [Planning artifacts](docs/planning-artifacts.md), [Host integrations](docs/host-integrations.md) |
| Recover interrupted work | [Recovery](docs/recovery.md) |
| Evaluate production use | [Compatibility and limitations](docs/compatibility-and-limitations.md), [Security model](docs/security-model.md) |
| Contribute or release | [Contributing](CONTRIBUTING.md), [Release checklist](docs/release-checklist.md), [Decision records](docs/decisions/index.md) |

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

## Contributing

Start with [CONTRIBUTING.md](CONTRIBUTING.md). Architecture and product constraints are fixed in [ARCHITECTURE.md](ARCHITECTURE.md) and [PHILOSOPHY.md](PHILOSOPHY.md); coding agents must also follow [AGENTS.md](AGENTS.md).

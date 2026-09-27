# Repository Agent Guidance

## Authority

Read `pathframe-analysis-plan.md` and `pathframe-coding-agent-implementation-planning-prompt.md` before architecture work. Treat their accepted implementation decisions as fixed unless newer explicit human evidence contradicts them. Follow `docs/planning/IMPLEMENTATION-ROADMAP.md` and the active stage plan.

## Product rules

- Preserve: “Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.”
- Preserve: “The agent reasons. Pathframe makes the path visible, packages the work, validates transitions, and preserves recovery.”
- Keep Brain responsible for reasoning and Pinky limited to explicitly delegated work.
- Never let delegated failure change `execution_policy` or authorize Brain fallback.
- Keep CLI and MCP over the same canonical application operations.
- Use typed operations; never add a generic raw-command MCP tool.
- Keep `okf-markdown/v1` human-readable and Pathframe-owned.
- Keep Aido optional and read-only.
- Complete sequential delegation before considering parallel Pinkies.
- Keep every non-terminal state recoverable without deleting `.pathframe/` or rewriting Git history.

## Implementation constraints

- Go 1.26; Linux amd64 first.
- Standard-library `flag`; no CLI framework.
- No-argument `pathframe` must orient.
- Official MCP Go SDK via local stdio, pinned in Stage 3 and isolated in the MCP adapter.
- Explicit human approval before every `planning -> ready`; material edits require reapproval.
- Verification uses structured argv, project-contained workdir, required timeout, bounded output, and no implicit shell.
- Authored artifacts plus append-only `history.jsonl` and bounded run records reconstruct replaceable `state.json`.
- No initial Specd importer, reader, or `.specd/` detection.

## Working method

Implement only the approved task and stage. Read targeted source, tests, contracts, and decisions. Reuse repository code before adding an abstraction or dependency. Keep adapter and domain dependencies pointing inward as described in `ARCHITECTURE.md`.

Do not copy legacy code until behavior, tests, license, and fit are traced and the reason is recorded. Do not implement deferred features for future flexibility.

Run focused verification, the full affected suite, and the stage journey. Report summary, changed files, commands, risks, and questions. Stop at every stage gate for human approval.

## Deviations

Record intentional deviations in `docs/decisions/` using its process. A deviation record must identify evidence, affected accepted decision or boundary, options, recommendation, consequences, and first blocked stage. Do not silently reverse accepted decisions.

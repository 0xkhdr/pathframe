# Repository Analysis

## Baseline

- `main` at `e949218`; remote `0xkhdr/pathframe`.
- Tracked files: empty `README.md`, `pathframe-analysis-plan.md`, and the implementation-planning prompt.
- No source, module, tests, fixtures, CI, release configuration, architecture record, or repository-local agent guidance exists.
- Local toolchain is Go 1.26.4 on Linux amd64. This matches the accepted primary platform.
- Git history shows a deliberate reset from Specd language to Pathframe, followed by explicit resolution of ten implementation decisions.

## Reuse and exclusions

Only product language, boundaries, accepted decisions, phases, contracts, and stage order in the two root planning documents are reusable. No implementation can be reused because none exists here. Neighboring repositories were not treated as repository content or dependencies.

Do not restore Specd subsystems, names, imports, compatibility detection, evidence machinery, mandatory Git lifecycle, parallel workers, plugin frameworks, or a proprietary agent runtime. Reference projects may later supply test ideas only after their behavior and license are traced in a scoped task.

## Missing capabilities

Everything from Stage 0 through Stage 9 remains unimplemented. Stage 0 must establish the Go module, minimal CLI, product constitution, guidance, tests, and CI before workflow behavior begins.

## Initial dependency direction

`cmd/pathframe` and `internal/adapters/{cli,mcp}` depend on `internal/app`. `internal/app` coordinates domain packages. Domain packages (`workflow`, `artifacts`, `context`, `delegation`, `recovery`, `verification`) depend only on narrow storage/host interfaces and standard-library value types. Filesystem, Git observation, process execution, MCP SDK, and host-specific integration code stay at adapter edges. CLI, MCP, generated guidance, and hooks contain no lifecycle policy.

## Build and test baseline

No commands currently build or test production code. Planning files can be checked with `find docs/planning -type f` and link/path inspection. Stage 0 introduces `go test ./...`, `go vet ./...`, `gofmt`, and Linux amd64 build checks.

## Integration constraints

- Codex first, Claude Code second.
- Typed MCP operations over local stdio; official SDK pinned only at Stage 3 and isolated in `internal/adapters/mcp`.
- Minimum host delegation: one sequential native subagent, shared workspace, supplied packet, returned result.
- Scope is advisory unless a host declares enforcement. Unknown capability means unsupported.
- Delegated execution never falls back to Brain.

## Risks and uncertainties

- Host configuration formats and hook support may change before Stages 3–6. Keep generated files behind host adapters and fixture-test exact supported versions.
- Go 1.26 is available locally but CI runner support must be proven in Stage 0.
- Markdown parsing can become complex. Own a deliberately small profile and reject unsupported syntax.
- Durable journal semantics need crash-safe append and deterministic replay tests before lifecycle claims.
- No representative Specd recovery fixtures exist in this repository. Stage 8 must encode historical dead-end scenarios from the authority document, not import legacy code.

## Contradictions and decisions

No repository evidence contradicts accepted decisions. No new architecture-blocking human decision is required now. Limits such as output byte caps, lease staleness, and context budgets are implementation constants/configuration choices recorded when first introduced; they do not reopen product policy.


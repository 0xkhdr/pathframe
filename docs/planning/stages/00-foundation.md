# Stage 0: Repository Foundation and Product Constitution

## Outcome and why now

Buildable Linux amd64 Go 1.26 CLI skeleton, governing docs, guidance, and CI. Every later package and contract depends on this baseline.

## Scope and non-goals

Add module, version/help CLI, no-argument placeholder orientation, smoke tests, README, philosophy, architecture, contributing/agent guidance, decision log, and Linux CI. No workflow persistence, `.pathframe/` mutation, MCP, templates, or production lifecycle behavior.

## Architecture and repository changes

Create `go.mod`; `cmd/pathframe/main.go`; `internal/adapters/cli/{cli.go,cli_test.go}`; `internal/app/version.go`; `.github/workflows/ci.yml`; `README.md`; `PHILOSOPHY.md`; `ARCHITECTURE.md`; `CONTRIBUTING.md`; `AGENTS.md`; `docs/decisions/README.md`. CLI uses `flag`; `main` depends on CLI adapter, which may read version only. No dependency besides standard library.

## Public contracts and migrations

Commands: `pathframe`, `pathframe --help`, `pathframe --version`. No arguments prints an honest “not initialized” orientation, not bare usage. Version injection may use a package variable with `dev` default. No migration.

## Ordered tasks and dependencies

| Task | Depends | Files/component | Invariant and focused check |
| --- | --- | --- | --- |
| S0-T1 Create constitution docs | none | README, PHILOSOPHY, ARCHITECTURE, CONTRIBUTING, AGENTS, decision log | Preserve product statement, philosophy, boundaries, accepted decisions; `rg` review |
| S0-T2 Bootstrap module/CLI | T1 | `go.mod`, `cmd/pathframe`, `internal/app`, CLI adapter | Go 1.26, `flag`, no-arg orientation; focused CLI tests |
| S0-T3 Add smoke tests | T2 | CLI tests | help/version/no-arg observable output |
| S0-T4 Add CI | T3 | workflow | format, vet, test, Linux amd64 build |
| S0-T5 Run contributor gate | T4 | docs/checklist | cold reader explains purpose/non-goals; human stop/go |

Wave 0: T1. Wave 1: T2. Wave 2: T3. Wave 3: T4. Wave 4: T5. Sequential by dependency.

## Context packet for implementers

Read both root authority files, this stage, architecture plan, and repository analysis. Write only listed Stage 0 files. Return summary, changed files, commands, risks, and questions. Do not load legacy repositories.

## Tests, journey, and recovery

Unit/smoke: help, version, no-arg output. Journey: clean checkout builds and a new user runs `pathframe`. CI/toolchain failure leaves docs/code intact; correct runner or code and rerun. Bad skeleton change is reverted by ordinary review, not state migration.

## Documentation

All constitution and contributor documents above. Record intentional deviations in `docs/decisions/`; none expected.

## Verification commands

```sh
test "$(go env GOVERSION)" = go1.26.4 || go version
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
GOOS=linux GOARCH=amd64 go build ./cmd/pathframe
rg -n "Pathframe is a local, deterministic|The agent reasons" README.md PHILOSOPHY.md ARCHITECTURE.md
```

## Acceptance, risks, and exit gate

- Commands build and smoke tests pass on Linux amd64.
- No-argument output orients; help/version work.
- Docs state product, boundaries, non-goals, accepted decisions, and no Specd guarantees.
- No legacy subsystem or third-party dependency appears.
- CI runs all four checks.
- Human confirms a cold contributor can explain purpose and non-goals.

Risk R-09 blocks exit if CI lacks Go 1.26. Correction: fix CI image/setup; never lower accepted version silently. Exit requires standard seven-point gate and human stop/go.


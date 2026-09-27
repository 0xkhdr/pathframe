# Stage 3: Codex-Native Planning Path

## Outcome and why now

Codex can assess, create, validate, approve-handoff, orient, and resume through typed tools without guessing CLI syntax. Canonical operations and artifacts already exist.

## Scope and non-goals

Pin official MCP Go SDK; local stdio only; generated Codex skill/manifest; typed planning tools; activation evaluations; optional supported session orientation. No Claude files, delegation, Streamable HTTP, custom JSON-RPC, or generic command tool.

## Architecture and files

Create `internal/adapters/mcp/{server,tools}.go`; `internal/integrations/codex/{generate,doctor}.go`; `templates/integrations/codex/`; `schemas/integration-v1.schema.json`; `testdata/integrations/codex/`; `testdata/evals/`; MCP/Codex journeys. SDK imports exist only under MCP adapter.

## Contracts and migrations

Pin SDK in `go.mod/go.sum`. Add `pathframe.integration/v1` and planning subset of nine named MCP tools. Add CLI integration install/update/doctor adapter commands if needed for humans; agents use typed MCP. Existing projects require no lifecycle migration.

## Tasks and waves

- S3-T1 (S2): pin SDK/version; map typed request/response schemas to app operations.
- S3-T2 (T1): stdio server and planning tools.
- S3-T3 (S2): Codex templates, activation rules, integration manifest.
- S3-T4 (T2,T3): install/update/diagnosis with non-destructive generated-file policy.
- S3-T5 (T4): parity, fresh-session, and must-use/offer/must-not-use evaluations.

Waves: T1; T2/T3; T4; T5. Context includes app contracts, exact adapter, generated templates, and fixture; never whole repo.

## Tests and recovery

Schema/parity tests compare app/CLI/MCP JSON. Stdio integration tests prove no raw command construction. Broken/missing/outdated config produces Doctor steps; regeneration changes machine-owned integration assets only. Failed server startup leaves workflow state untouched and CLI orientation usable.

## Documentation and correction

Add Codex setup, activation, update, troubleshooting, and supported capability docs. SDK failure is corrected within adapter or pin; domain code remains unchanged.

## Verification

```sh
go test ./internal/adapters/mcp ./internal/integrations/codex ./internal/app
go test ./tests/journey -run CodexPlanning
go test ./... && go vet ./...
```

## Acceptance and gate

Codex completes request-to-ready journey through typed tools; no guessed argument order; all activation classes pass; fresh session resumes canonical state; CLI/MCP results match; SDK remains adapter-local. Do not start Claude or delegation until human stop/go.


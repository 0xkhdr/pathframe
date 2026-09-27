# Stage 4: Claude Code Parity

## Outcome and why now

Claude Code completes the same planning and resume journey through native guidance/MCP configuration. Codex proves shared contracts first.

## Scope and non-goals

Generate Claude skill/commands, MCP settings, only supported lifecycle hooks, manifest, update/Doctor checks, and shared evaluations. No duplicated domain rules, delegation, or invented hook capability.

## Architecture and files

Create `internal/integrations/claude/{generate,doctor}.go`; `templates/integrations/claude/`; `testdata/integrations/claude/`; Claude journeys. Reuse MCP adapter, app operations, schemas, and semantic evaluation cases.

## Contracts and migrations

Same `pathframe.integration/v1` with host-specific capability values. Hooks call typed/canonical operations and are optional when unsupported. Generated-file upgrades are integration migrations only, previewed and diagnosable.

## Tasks and waves

- S4-T1 (S3 gate): map supported Claude surfaces/capabilities.
- S4-T2 (T1): generate skill/commands/MCP config and manifest.
- S4-T3 (T2): supported hooks as thin operation calls.
- S4-T4 (T2,T3): install/update/Doctor fixtures.
- S4-T5 (T4): run shared planning, activation, broken-config, and resume journeys.

Implementers read Stage 3 contracts and Claude-owned files only. Sequential waves prevent generator fixture conflicts.

## Tests and recovery

Golden generated files, parity with Codex canonical results, fresh-session resume, and broken configuration. Unsupported hook is declared absent and workflow remains operable through tools. Doctor repairs generated assets without editing authored plans.

## Documentation and correction

Add Claude setup/capability matrix/troubleshooting. Host drift changes only Claude generator/fixtures and declared version.

## Verification

```sh
go test ./internal/integrations/claude ./internal/adapters/mcp ./internal/app
go test ./tests/journey -run ClaudePlanning
go test ./...
```

## Acceptance and gate

Both hosts pass same planning semantics and resume journey; differences are explicit capabilities; no lifecycle logic is duplicated; broken config yields actionable diagnosis. Human stop/go before packet work.

Implementation evidence: Claude Code uses project-native `.claude/skills/pathframe/SKILL.md`, `.claude/commands/pathframe.md`, and `.mcp.json`. Optional orientation uses only the supported `SessionStart` command hook in `.claude/settings.json`, calling canonical `status --json`. Manifest hashes protect user content during install/update; Doctor reports actionable drift. Both hosts share the same activation fixtures and MCP/application planning operations.

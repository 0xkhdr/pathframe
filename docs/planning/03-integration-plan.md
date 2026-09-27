# Integration Plan

## Shared rule

CLI, MCP, generated guidance, and hooks call `internal/app`; none owns transitions. Machine identifiers pass through typed results. Generic `run_pathframe_command(string)` is prohibited.

## CLI

Stage 0 adds `cmd/pathframe` and `internal/adapters/cli` using `flag` plus a small dispatcher. No arguments calls Orient. Later commands are thin mappings: `status`, `next`, `new`, `check`, `approve`, `pause`, `resume`, `replan`, `cancel`, `doctor`, integration install/update, packet preview, submit, verify, accept, and request-changes. `--json` returns canonical contract JSON.

## MCP

Stage 3 pins official MCP Go SDK inside `internal/adapters/mcp`. Local stdio only. Tools: `pathframe_orient`, `pathframe_assess_request`, `pathframe_create_change`, `pathframe_get_template`, `pathframe_validate_plan`, `pathframe_get_next`, `pathframe_prepare_delegation`, `pathframe_submit_result`, and `pathframe_recover`. Later stages activate tools only after underlying app operation exists. Schema parity tests compare MCP and CLI/app results.

## Codex

Stage 3 generates a focused skill and integration manifest from embedded templates. Guidance contains must-use, offer-once, and must-not-use rules. Optional session-start orientation is enabled only when supported. Stage 6 adds a host adapter contract that accepts one packet, launches one native subagent in shared workspace, waits, and returns its final result.

## Claude Code

Stage 4 generates skill/commands, MCP settings, supported hooks, and manifest. Hooks only invoke canonical operations. Stage 6 adds the same sequential launch/result contract. Missing hook capabilities reduce declared assurance; they do not get simulated or claimed.

## Capability model

`pathframe.integration/v1` declares host/version, sequential subagent launch, shared workspace, result return, and optional enforcement/hook capabilities. Unknown or absent capabilities fail preflight. Advisory write scope is labeled advisory.

## Installation and diagnosis

Generated files carry template/integration version and a content marker. Install/update previews changes and never overwrites unknown user content silently. Doctor checks missing files, incompatible versions, malformed configuration, unavailable capabilities, and stale generated assets; repair regenerates only machine-owned integration files with backup/preview behavior defined by the host adapter.

## Evaluations

`testdata/evals/{must-use,offer,must-not-use}.yaml` covers explicit activation, active-change continuation, multi-step offers, trivial edits, read-only review, explanation, and explicit direct execution. Both hosts run the same semantic cases plus format-specific installation/resume/broken-config journeys.

## Order

Codex planning must pass before Claude Code work. Both planning paths must pass before packet work. Packets must pass inspection before delegation. Sequential delegation must pass before completion/recovery hardening. Parallelism remains deferred.


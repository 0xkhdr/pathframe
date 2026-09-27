# Risk and Decision Register

## Accepted decisions

| ID | Fixed decision | Enforced first |
| --- | --- | --- |
| AD-01 | Go 1.26; Linux amd64 initially | Stage 0 |
| AD-02 | Standard `flag`; no-argument orientation | Stage 0/1 |
| AD-03 | Official MCP Go SDK, stdio, adapter-local | Stage 3 |
| AD-04 | Human approval for every `planning -> ready`; material edits reapprove | Stage 2 |
| AD-05 | Explicit execution policy; Quick defaults Brain; no delegated fallback | Stage 2/6 |
| AD-06 | Pathframe-owned `okf-markdown/v1`; Aido optional read-only | Stage 2/5 |
| AD-07 | Minimum sequential shared-workspace native subagent; honest assurance | Stage 6 |
| AD-08 | Pathframe executes bounded argv verification without implicit shell | Stage 7 |
| AD-09 | Authored intent + journal + runs reconstruct replaceable state | Stage 1/8 |
| AD-10 | No initial Specd compatibility; conditional importer review at Stage 9 | Stage 9 |

These are not approval requests.

## Risks

| ID | Risk | Mitigation / recovery | First gate |
| --- | --- | --- | --- |
| R-01 | Host formats/capabilities drift | Stages 3–4 pin `pathframe.integration/v1`, generated-file hashes, host fixtures, and actionable diagnosis; host changes stay in host generators or the MCP adapter | 3/4 |
| R-02 | Parser accepts ambiguous Markdown | Stage 2 uses flat strict front matter, exact headings, structured verification, and negative fixtures; profile changes remain versioned | 2 |
| R-03 | Journal append interrupted | Ignore/diagnose incomplete final record; replay prior complete events | 1/8 |
| R-04 | Projection diverges | Compare replay identity; Doctor rebuilds machine projection | 1/8 |
| R-05 | Path/symlink escape | Resolve against project root and reject escapes before I/O/exec | 0/7 |
| R-06 | Delegation falsely implies containment | Stage 5 packets label write scope `advisory` and host assurance `not_evaluated`; Stage 6 may strengthen assurance only from declared host facts | 5/6 |
| R-07 | Required context exceeds budget | Stage 5 blocks packet without truncation and reports split-task or explicit `budget_bytes` increase recovery | 5 |
| R-08 | Verification becomes stale after edits | Bind results to deterministic repository content identity | 7 |
| R-09 | Go 1.26 unavailable in CI | Stage 0 gate fails; fix runner, do not lower version silently | 0 |
| R-10 | Scope expands toward Specd | Constitution review and deviation record at every gate | all |
| R-11 | Concurrent CLI writers can race between replay and append | Stage 1 is single-writer; diagnose an invalid chain without rewriting complete records; add cross-process serialization when concurrent mutation is in scope | 9 |
| R-12 | A repository may already own `.codex/config.toml` or hooks | Stage 3 refuses to overwrite unmanifested or modified files and reports a move/restore/reinstall recovery; native config merging waits for a proven safe ownership contract | 3 |
| R-13 | A repository may already own `.mcp.json` or `.claude/settings.json` | Stage 4 refuses to overwrite or merge unmanifested/modified files; users move custom content aside or omit the optional hook. Native merge support waits for a proven ownership contract | 4 |
| R-14 | Host crash leaves an active lease | Stage 8 diagnoses leases older than the named 24-hour recovery window and records release before explicit retry or replan; no fallback occurs | 6/8 |
| R-15 | Worker or verification output grows without bound | Stage 6 rejects task results above 256 KiB; Stage 7 independently bounds stdout and stderr with explicit truncation | 6/7 |
| R-16 | Repository content changes after a passing verification | Stage 7 binds verification to a deterministic relevant-content identity and refuses stale semantic acceptance | 7 |
| R-17 | Replanning invalidates already accepted work | Stage 8 binds completion to normalized task-contract identity; unchanged tasks remain complete and changed tasks return to the frontier | 8 |
| R-18 | Cross-build success is mistaken for platform support | Compatibility docs claim only Linux amd64; every other target requires native clean-install and complete-journey CI before promotion | 9 |
| R-19 | Installer interruption replaces a working binary with a partial file | Stage beside the destination and atomically rename only after copy and permission checks; failed pre-rename updates leave the prior binary intact | 9 |
| R-20 | Host releases drift beyond repository fixtures | Versioned manifests and Doctor checks remain the support boundary; rerun host journeys after updates and reduce claims on failure | 9 |

## New decisions requiring approval

None. Exact numeric defaults for output limits, timeouts, and context budgets should be conservative named constants when implemented, tested at boundaries, and changed without contract-version churn unless persisted meaning changes.

## Deferred decisions

Parallel Pinkies, strict Git evidence, importer creation, broader hosts, Streamable HTTP, and organization controls remain outside Stage 0–9 acceptance. Stage 9 may only recommend a read-only-source importer when real demand and representative fixtures exist.

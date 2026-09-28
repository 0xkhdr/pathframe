# Compatibility and limitations

“Supported” means native build, clean installation, full sequential journey, and the Stage 9 verification suite pass. A cross-build alone establishes only a portability target.

Evidence for the initial release:

| Surface | Status | Evidence |
| --- | --- | --- |
| Linux amd64, Go 1.26 | Supported | native CI build, race/full/security/install suites, and complete sequential journeys |
| Linux arm64 | Portability target | cross-build only; no install or journey claim |
| macOS amd64/arm64 | Portability target | cross-build only; no install or journey claim |
| Windows amd64 | Portability target | cross-build only; no install or journey claim |
| Codex | Supported contract | generated `pathframe.integration/v1`, local stdio MCP, sequential shared-workspace subagent contract, repository journey tests |
| Claude Code | Supported contract | generated `pathframe.integration/v1`, local stdio MCP, sequential shared-workspace subagent contract, repository journey tests |
| Aido | Optional, read-only | referenced knowledge only; no shared lifecycle or cross-writing |
| Specd | Incompatible by design | no detection, reader, importer, or compatibility layer |

Host product releases do not have a stable version handshake in this repository. Run `pathframe codex-doctor` or `pathframe claude-doctor` after host updates; a failed check removes the support claim for that local configuration until repaired and re-tested.

## Known limitations

- Pathframe coordinates one sequential Pinky. Parallel workers, worktree orchestration, and merge policy are not implemented.
- Write scope is advisory unless the host explicitly declares enforcement. Pathframe does not sandbox workers or verification programs.
- Repository mutation assumes one Pathframe writer at a time; concurrent processes may be diagnosed as an invalid journal chain.
- Host configuration installation is ownership-based and refuses to merge unknown existing files.
- Git is observational, not a lifecycle or recovery authority.
- Run records bound output but are not cryptographic audit evidence.
- Release publication and deployment orchestration are outside Pathframe; maintainers produce and verify release artifacts using the checklist.

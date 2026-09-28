# Codex and Claude Code integrations

Build or install `pathframe`, then install the integration from the project root:

| Host | Install | Check |
| --- | --- | --- |
| Codex | `pathframe codex-install --session-orientation` | `pathframe codex-doctor` |
| Claude Code | `pathframe claude-install --session-orientation` | `pathframe claude-doctor` |

Restart the host after installation. Claude Code prompts for project MCP approval; Codex requires the project to be trusted and project hooks to be reviewed. Session orientation is optional; omit `--session-orientation` if the repository already owns the host settings file.

| Asset | Codex | Claude Code |
| --- | --- | --- |
| Repository guidance | `.agents/skills/pathframe/SKILL.md` | `.claude/skills/pathframe/SKILL.md` and `.claude/commands/pathframe.md` |
| MCP configuration | `.codex/config.toml` | `.mcp.json` |
| Optional session orientation | `.codex/hooks.json` | `.claude/settings.json` |

Generated assets are manifest-owned. Pathframe refuses to overwrite an unknown or modified path; preserve or move custom content, restore the generated version if needed, reinstall, and rerun the host Doctor command.

## Activation and delegation

Both hosts use the same rules. Activate Pathframe when explicitly requested or resuming an active workflow. Offer it once for multi-step, durable, agreement-sensitive, or risky development work. Do not activate it for explanations, read-only review, brainstorming, trivial isolated edits, non-development work, or explicit direct execution.

The local stdio server exposes typed planning, delegation, verification, and review operations. Use those operations instead of constructing CLI commands. `pathframe_prepare_delegation` without a host is a read-only packet preview; Pathframe itself never launches a worker.

For a ready delegated task, call `pathframe_prepare_delegation` with `host: codex` or `host: claude-code`. After preflight creates the exclusive lease, Brain launches one native subagent with the returned packet and Pinky rules, then submits Pinky's unchanged `pathframe.task-result/v1` with `pathframe_submit_result`.

Brain runs bounded Pathframe verification and separately accepts or requests changes. Pinky cannot approve itself or alter plan state; worker checks are supplemental. Stale content blocks acceptance, and advisory scope violations require an explicit keep, revert, or replan decision. Failed launch or execution never changes `execution_policy`; `pathframe_check_brain_edit` continues to deny Brain edits. Recover a lost lease by exact ID with `pathframe_release_delegation`, then retry or replan.

Run the host Doctor command after host updates. Failed MCP startup does not mutate workflow state; CLI orientation remains available.

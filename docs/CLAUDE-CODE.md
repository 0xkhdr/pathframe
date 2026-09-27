# Claude Code integration

Install project-native planning support from the repository root:

```sh
pathframe claude-install --session-orientation
pathframe claude-doctor
```

Restart Claude Code and approve the project MCP server when prompted. The generated skill is `.claude/skills/pathframe/SKILL.md`; `/pathframe` is provided by `.claude/commands/pathframe.md`; `.mcp.json` starts the existing local-stdio Pathframe MCP server. The optional `.claude/settings.json` hook runs canonical `status --json` at supported `SessionStart` events.

| Capability | Claude Code | Codex |
| --- | --- | --- |
| Repository guidance | project skill + slash command | repository skill |
| Typed planning operations | stdio MCP | stdio MCP |
| Session orientation | optional `SessionStart` command hook | optional `SessionStart` MCP hook |

Pathframe refuses to overwrite an existing or modified generated path. Move custom content aside or restore the manifest-owned version, rerun `claude-install`, then rerun `claude-doctor`. Install without `--session-orientation` when the repository already owns `.claude/settings.json`; planning remains available through the skill and typed tools.

Activation is identical on both hosts: use Pathframe when explicitly requested or resuming an active workflow; offer it once for multi-step, durable, agreement-sensitive, or risky development work; do not activate it for explanations, read-only review, brainstorming, trivial isolated edits, non-development work, or explicit direct execution.

# Codex integration

Build or install `pathframe`, then run from the project root:

```sh
pathframe codex-install --session-orientation
pathframe codex-doctor
```

Restart Codex after installation. The generated repository skill lives at `.agents/skills/pathframe/SKILL.md`; trusted-project MCP configuration lives at `.codex/config.toml`. Session-start orientation is optional and lives at `.codex/hooks.json`; Codex requires review and trust for project hooks.

The skill activates Pathframe when explicitly requested or when a workflow is active or continued. It offers Pathframe once for dependent, durable, agreement-sensitive, or meaningfully risky development work. It does not activate for explanations, read-only review, brainstorming, trivial isolated edits, non-development operations, or explicit direct execution.

The stdio server exposes typed planning tools for orientation, request assessment, change creation, templates, validation/approval handoff, next action, and recovery. It exposes no raw-command tool. `prepare_delegation` and `submit_result` remain unavailable until their application operations exist in later approved stages.

If diagnosis reports a modified generated file, preserve or move the custom file before reinstalling. Pathframe never overwrites unknown content. Failed MCP startup does not mutate workflow state; CLI orientation remains available.

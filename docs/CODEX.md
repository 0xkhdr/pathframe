# Codex integration

Build or install `pathframe`, then run from the project root:

```sh
pathframe codex-install --session-orientation
pathframe codex-doctor
```

Restart Codex after installation. The generated repository skill lives at `.agents/skills/pathframe/SKILL.md`; trusted-project MCP configuration lives at `.codex/config.toml`. Session-start orientation is optional and lives at `.codex/hooks.json`; Codex requires review and trust for project hooks.

The skill activates Pathframe when explicitly requested or when a workflow is active or continued. It offers Pathframe once for dependent, durable, agreement-sensitive, or meaningfully risky development work. It does not activate for explanations, read-only review, brainstorming, trivial isolated edits, non-development operations, or explicit direct execution.

The stdio server exposes typed planning and delegation tools. `pathframe_prepare_delegation` without `host` remains a read-only preview. With `host: codex`, it validates declared sequential/shared-workspace/result-return capabilities and creates one exclusive lease before Brain launches one native subagent. Pinky returns `pathframe.task-result/v1`; Brain passes it unchanged to `pathframe_submit_result`. Pathframe reconciles but does not accept or complete the task. `pathframe_release_delegation` recovers a lost lease by exact ID. `pathframe_check_brain_edit` denies delegated work even after failure. Scope is advisory because Codex does not declare write-scope enforcement.

If diagnosis reports a modified generated file, preserve or move the custom file before reinstalling. Pathframe never overwrites unknown content. Failed MCP startup does not mutate workflow state; CLI orientation remains available.

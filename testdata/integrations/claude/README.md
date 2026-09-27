# Claude Code integration fixture

Tests generate the project skill, slash command, `.mcp.json`, optional `SessionStart` hook, and `pathframe.integration/v1` manifest in a temporary repository. Absolute executable paths are intentionally generated at test time; semantic assertions cover activation guidance, configuration shape, declared capabilities, ownership hashes, safe updates, diagnosis, and resume behavior.

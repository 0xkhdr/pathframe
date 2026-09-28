# Getting started

Purpose: complete the smallest direct Pathframe workflow and learn its approval boundary.

Install Pathframe, enter a repository, and run `pathframe`. With no arguments it reports whether the project is configured, the active phase, blockers, and the next legal action.

For a small direct change:

```sh
pathframe new --change demo --mode quick
pathframe template --mode quick --artifact intent
pathframe template --mode quick --artifact task
pathframe check --change demo
pathframe approve --change demo
pathframe next --change demo
```

Fill the generated Markdown before `check`. `approve` represents explicit human approval; never call it merely because validation passed. For Codex or Claude Code, follow [Host integrations](host-integrations.md) and use the typed MCP operations instead of constructing CLI commands.

During execution, `pathframe doctor --change demo` diagnoses recovery. Use the reported pause, release/retry, replan, repair, or cancel operation. Do not delete `.pathframe/` or rewrite Git history.

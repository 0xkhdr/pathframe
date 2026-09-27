# Command-line reference

Running `pathframe` with no arguments is equivalent to `pathframe status`. Both answer whether the project is configured, which change and phase are active, task progress, blockers, the recommended next action, legal alternatives, and whether a human is required.

```text
pathframe [--json] [--change ID] [status|next|new|template|check|approve|pause|resume|replan|cancel|codex-install|codex-doctor]
```

- `status` and `next` render the canonical `pathframe.workflow/v1` orientation.
- `pause`, `resume`, `replan`, and `cancel` append a human transition to the selected change journal and rebuild `state.json`.
- `--json` renders the canonical object; text is a projection of the same object.
- `--change ID` is required when more than one change exists.
- `new --change ID --mode quick|standard|high-risk` creates a planning change and only its required missing artifacts.
- `template --mode MODE --artifact KIND` prints the owned Markdown template; `--json` returns its typed filling instructions.
- `check --change ID` validates artifacts, tasks, references, questions, acceptance, and structured verification.
- `approve --change ID` is the explicit human approval operation and binds the validated material plan identity.
- `codex-install [--session-orientation]` installs manifest-owned Codex skill and stdio MCP assets; the optional flag adds supported session-start orientation.
- `codex-doctor` validates the integration version and generated-file hashes without editing workflow state.

Flags may follow the command. Illegal transitions exit nonzero after returning current state, a stable reason code, and executable recovery actions. An incomplete final journal record is diagnosed and safely discarded without inventing state; a missing, corrupt, or divergent projection is rebuilt from complete journal records.

See [Planning artifacts](ARTIFACTS.md) for the profile and recovery rules and [Codex integration](CODEX.md) for setup and activation. `pathframe mcp` is the generated local stdio entry point, not a generic command executor. Stage 3 intentionally has no Claude support, delegation, context packet, or verification execution commands.

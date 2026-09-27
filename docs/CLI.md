# Command-line reference

Running `pathframe` with no arguments is equivalent to `pathframe status`. Both answer whether the project is configured, which change and phase are active, task progress, blockers, the recommended next action, legal alternatives, and whether a human is required.

```text
pathframe [--json] [--change ID] [status|next|new|template|check|approve|packet|submit-result|verify|accept-task|request-changes|lease-release|edit-check|doctor|pause|resume|replan|cancel|codex-install|codex-doctor|claude-install|claude-doctor]
```

- `status` and `next` render the canonical `pathframe.workflow/v1` orientation.
- `pause`, `resume`, `replan`, and `cancel` append a human transition to the selected change journal and rebuild `state.json`.
- `--json` renders the canonical object; text is a projection of the same object.
- `--change ID` is required when more than one change exists.
- `doctor --change ID [--repair]` previews typed recovery diagnoses and optionally applies only safe machine-state repairs. See [Doctor and recovery](DOCTOR.md).
- `new --change ID --mode quick|standard|high-risk` creates a planning change and only its required missing artifacts.
- `template --mode MODE --artifact KIND` prints the owned Markdown template; `--json` returns its typed filling instructions.
- `check --change ID` validates artifacts, tasks, references, questions, acceptance, and structured verification.
- `approve --change ID` is the explicit human approval operation and binds the validated material plan identity.
- `packet --change ID [--task ID] [--budget-bytes N]` previews the bounded packet. Add `--host codex|claude-code` to run preflight and acquire one exclusive lease before Brain launches the native Pinky. Pathframe itself never launches a host worker.
- `submit-result --result-file PATH` validates the leased `pathframe.task-result/v1`, records it, and reconciles to reviewing or blocked. It never accepts or completes a task.
- `verify --change ID --task ID --timeout-ms N [--workdir PATH] [--max-bytes N] [--changed-files LIST]` runs every approved argv directly, never through a shell. Workdir must resolve inside the project; timeout is required; stdout/stderr are independently bounded. A directly authorized Brain task supplies its changed files here; a delegated task reuses Pinky's submitted list. Results record exit, duration, timeout/interruption, truncation, content identity, and advisory-scope violations.
- `accept-task --change ID --task ID --reason TEXT [--keep-scope-violations]` records Brain's semantic acceptance only after fresh passing Pathframe verification. The explicit scope flag records a keep decision; Pathframe never reverts files automatically.
- `request-changes --change ID --task ID --reason TEXT` records Brain's rejection and leaves retry/replan recovery available under the unchanged execution policy.
- `lease-release --change ID --lease-id ID` explicitly releases a lost worker lease and preserves a recoverable blocker.
- `edit-check --change ID --task ID` enforces Brain edit authority. Delegated policy is denied regardless of delegation failure.
- `codex-install [--session-orientation]` installs manifest-owned Codex skill and stdio MCP assets; the optional flag adds supported session-start orientation.
- `codex-doctor` validates the integration version and generated-file hashes without editing workflow state.
- `claude-install [--session-orientation]` installs manifest-owned Claude Code guidance, stdio MCP configuration, and optional supported session-start hook.
- `claude-doctor` validates the Claude Code integration version and generated-file hashes without editing workflow state.

Flags may follow the command. Illegal transitions exit nonzero after returning current state, a stable reason code, and executable recovery actions. An incomplete final journal record is diagnosed and safely discarded without inventing state; a missing, corrupt, or divergent projection is rebuilt from complete journal records.

See [Planning artifacts](ARTIFACTS.md), [Codex integration](CODEX.md), and [Claude Code integration](CLAUDE-CODE.md). `pathframe mcp` is the generated local stdio entry point, not a generic command executor. Worker-reported checks remain supplemental; task completion requires both Pathframe verification and Brain semantic acceptance.

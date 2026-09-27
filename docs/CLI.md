# Command-line reference

Running `pathframe` with no arguments is equivalent to `pathframe status`. Both answer whether the project is configured, which change and phase are active, task progress, blockers, the recommended next action, legal alternatives, and whether a human is required.

```text
pathframe [--json] [--change ID] [status|next|pause|resume|replan|cancel]
```

- `status` and `next` render the canonical `pathframe.workflow/v1` orientation.
- `pause`, `resume`, `replan`, and `cancel` append a human transition to the selected change journal and rebuild `state.json`.
- `--json` renders the canonical object; text is a projection of the same object.
- `--change ID` is required when more than one change exists.

Flags may follow the command. Illegal transitions exit nonzero after returning current state, a stable reason code, and executable recovery actions. An incomplete final journal record is diagnosed and safely discarded without inventing state; a missing, corrupt, or divergent projection is rebuilt from complete journal records.

Stage 1 intentionally has no commands for creating changes, authoring templates, approval, delegation, MCP, or verification execution.

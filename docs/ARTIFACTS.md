# Planning artifacts

Pathframe owns the human-readable `okf-markdown/v1` profile. Each Markdown file starts with flat `key: value` front matter and uses the exact `##` sections returned by `pathframe template`. Unknown prose is allowed inside sections; missing sections, placeholders, unresolved questions, broken references, and shell-string verification are rejected.

Modes require:

| Mode | Files |
| --- | --- |
| Quick | `intent.md`, `tasks/*.md` |
| Standard | Quick plus `requirements.md`, `design.md` |
| High-risk | Standard plus `risks.md`, `rollout.md`, `recovery.md` |

`change.yaml` declares `pathframe.change/v1`, the profile, stable change ID, and mode. Each task is one `pathframe.task/v1` file with stable ID/title, explicit `execution_policy`, role, references, dependencies, objective, required reads, write scope, constraints, JSON argv verification, observable acceptance, questions, and assumptions. Quick templates explicitly default policy to `brain`.

Use `pathframe new --change ID --mode MODE` to create only missing files. Raising mode adds missing artifacts without overwriting authored content; lowering mode is rejected. `pathframe template --mode MODE --artifact KIND` prints a template, or add `--json` for `pathframe.artifact-instructions/v1`.

`pathframe check --change ID` validates the complete plan. `pathframe approve --change ID` validates again and records explicit human approval of the normalized material identity. Whitespace-only normalization does not change that identity. Content, references, mode, or execution-policy edits do; the next canonical operation moves a ready plan to `replanning`, preserves files and history, and requires validation and approval again.

Pathframe owns `history.jsonl` and `state.json`. Humans and agents own `change.yaml`, artifact Markdown, and task files. Aido is neither required nor written or detected by these operations.

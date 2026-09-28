# Planning artifacts

Purpose: define Pathframe's authored planning files, validation rules, roles, and context packets.

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

## Roles and context packets

Roles live at `.pathframe/roles/<id>.yaml` with `schema: pathframe.role/v1`, an ID and mission, JSON-array `reads`, `optional_reads`, `allowed_actions`, and `returns`. Role reads are project-relative foundation context. Task `Required Reads` are change-relative context; the selected task contract and current canonical workflow facts supply task and runtime context. Paths and symlink targets must remain inside their owning root.

`pathframe packet` assembles `pathframe.task/v1` only for an approved frontier task. The sorted frontier and dependency waves describe readiness; `next` selects one task for sequential execution and does not authorize parallel mutation. The default budget is 65536 bytes and `--budget-bytes` changes it explicitly. Required context is complete or packet creation blocks with split/raise recovery. Optional context is included in declared order while it fits; every omission records its layer, path, size, and reason. Unreferenced tasks and knowledge are excluded. Declared write scope is advisory, never described as filesystem containment.

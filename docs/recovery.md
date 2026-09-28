# Doctor and recovery

Purpose: diagnose damaged or interrupted workflows and describe safe recovery paths.

`pathframe doctor --change ID` returns `pathframe.diagnosis/v1` without changing files. `--repair` applies only the repair named by a diagnosis: rebuilding replaceable `state.json`, discarding an incomplete trailing journal fragment, or releasing an abandoned lease through recorded events.

Doctor never edits authored Markdown or `change.yaml`, deletes `.pathframe/`, rewrites Git history, or repairs malformed complete journal/run records. Those diagnoses require a human to repair the named source and rerun Doctor.

## Interruption playbook

Run Doctor, review its evidence, and use `--repair` only for its safe machine-state repairs. For authored artifacts, roles, references, verification failures, or capabilities, fix the named input and use the reported workflow recovery: resume, retry with the same execution policy, replan, or cancel.

Cancellation appends the workflow transition and, when work is leased, a `cancelled` run event before removing the replaceable active lease. History remains available.

## Replanning

Semantic completion records the normalized task-contract identity. After reapproval, completed tasks whose contracts are unchanged remain completed; changed or removed task contracts do not contribute to the new frontier. Execution-policy changes remain material and require explicit human reapproval.

# Contracts and Artifacts

## Contract catalog

| Contract | Owner | First stage | Purpose |
| --- | --- | --- | --- |
| `pathframe.workflow/v1` | `workflow` + `app` | 1 | Canonical orientation, legal actions, blockers |
| `okf-markdown/v1` | `artifacts` | 2 | Minimal Markdown/profile rules |
| `pathframe.artifact-instructions/v1` | `artifacts` | 2 | Typed template-filling instructions |
| `pathframe.intent/v1`, `requirements/v1`, `design/v1`, `task/v1` artifact front matter | `artifacts` | 2 | Human-authored intent and work |
| `pathframe.integration/v1` | integrations | 3 | Host identity, version, capabilities |
| `pathframe.task/v1` packet | `context` | 5 | Bounded executable task context |
| `pathframe.task-result/v1` | `delegation` | 6 | Worker submission, never acceptance |
| transition event/run records | `store` | 1/6/7 | Durable lifecycle facts and bounded attempts |

Schemas live in `schemas/*.schema.json`; matching Go structs live with owning package. Golden fixtures prove encoding. Avoid a central schema registry.

## Workflow result minimum

Fields: schema, project root, selected change, phase, task progress, active task/worker, recommended action, legal alternatives, blockers with recovery operations, human-required flag, and stable opaque identifiers. Text is a projection of this object.

## Authored layout

```text
.pathframe/project.yaml
.pathframe/knowledge/
.pathframe/roles/<role>.yaml
.pathframe/changes/<change>/change.yaml
.pathframe/changes/<change>/{intent,requirements,design,risks,rollout,recovery}.md
.pathframe/changes/<change>/tasks/<task>.md
```

Quick requires intent plus task files only. Standard adds requirements and design. High-risk adds risks, rollout, and recovery. Mode may rise, never silently fall.

Every task contains stable ID/title, `execution_policy`, optional role, objective, references, dependencies, required reads, write scope, constraints, verification argv/workdir/timeout, observable acceptance, questions, and assumptions. Quick defaults to explicit `brain`; changing policy after approval invalidates approval.

## Machine-owned layout

```text
.pathframe/changes/<change>/history.jsonl
.pathframe/changes/<change>/state.json
.pathframe/changes/<change>/runs/<run>.json
```

Events include unique ID, schema, timestamp, actor class, source/target state, subject IDs, reason, and data needed for replay. Runs store leases, submissions, verification, reconciliation, content identity, exit/duration/timeout, and bounded stdout/stderr. Truncation is explicit. `state.json` contains replay position and projection only.

## Approval and materiality

Approval event binds the validated authored-artifact content identity. Any change to plan-bearing artifact content or execution policy is material and returns the change to planning/replanning. Machine records, formatting-only normalization explicitly proven by the parser, and routine transitions are not plan changes. Every mode requires explicit human approval.

## Task packet and result

Packet includes change/task IDs, role, objective/acceptance, refs and resolved required context, dependencies, write scope plus enforcement assurance, verification, constraints, result schema, integration capabilities, omissions with reasons, and budget accounting. Required context never truncates.

Result status is `completed`, `failed`, `blocked`, or `needs_replan`; fields include summary, changed files, supplemental verification, discoveries, questions, and risks. Submission cannot complete or accept a task.

## Recovery contracts

Every rejected operation returns current state, stable reason code, and at least one executable typed recovery operation. No recovery instructs deletion or Git history rewriting. Ambiguous active changes require explicit selection.


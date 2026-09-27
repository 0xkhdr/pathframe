# Stage 5: Context Packages and Roles

## Outcome and why now

Brain previews a complete, bounded `pathframe.task/v1` packet and deterministic sequential frontier before any worker launch.

## Scope and non-goals

Role schema, four context layers, reference resolution, required/optional classification, budgets, omission reasons, packet preview, task DAG/frontier/waves. No subagent launch, lease, result submission, or parallel mutation.

## Architecture and files

Create `internal/context/{model,resolve,budget,packet}.go`; `internal/delegation/frontier.go`; `internal/app/prepare.go`; `schemas/{role,task-packet}-v1.schema.json`; role/context/packet fixtures; CLI/MCP packet preview mapping; journeys.

## Contracts and migrations

Add role v1 and `pathframe.task/v1`. Extend `pathframe_prepare_delegation` to validate/preview without launch. Add packet preview human projection. Existing tasks missing role under `brain` remain valid; delegated tasks missing roles block with repair guidance.

## Tasks and waves

- S5-T1 (S4): role and context-reference schemas.
- S5-T2 (T1): safe reference resolver across foundation/change/task/runtime layers.
- S5-T3 (T2): required/optional budgets and visible omissions.
- S5-T4 (S2,T3): packet assembler and human preview.
- S5-T5 (S2): dependency DAG, ready frontier, sequential waves/cycle diagnostics.
- S5-T6 (T4,T5): CLI/MCP parity and packet-quality journeys.

Agents get exact task, referenced fixture, owning package, schemas, and focused test. Write scope excludes host adapters.

## Tests and recovery

Tests prove required context never truncates, optional omission is reported, unrelated tasks are absent, references cannot escape root, roles stay separate from task facts, cycles block with replan action, and preview equals serialized packet. Oversized required context recovers by splitting task or explicit budget increase.

## Documentation and correction

Document role authoring, context layers, budget semantics, assurance labels, and packet preview. Packet schema corrections require fixtures and compatibility review; authored sources remain unchanged.

## Verification

```sh
go test ./internal/context ./internal/delegation ./internal/app
go test ./tests/journey -run 'Packet|Frontier|Context'
go test ./...
```

## Acceptance and gate

Human preview matches machine packet; required context is complete; irrelevant context absent; all omissions visible; frontier deterministic; cycles and budget failures recover without deletion. No worker launches. Human approves packet quality before Stage 6.


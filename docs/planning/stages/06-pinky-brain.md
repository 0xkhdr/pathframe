# Stage 6: Pinky and Brain Vertical Slice

**Status:** implemented; awaiting human stop/go gate.

## Outcome and why now

Codex and Claude Brain safely launch one sequential native Pinky, receive a structured result, and reconcile it without implicit fallback. Packet quality is already independently proven.

## Scope and non-goals

Capability declarations, preflight, durable lease/run, two host adapters, worker guidance, result validation/submission, Brain edit guard, no-fallback blocker. No Pathframe verification acceptance yet, parallel workers, sandbox claims, cancellation, worktrees, or resumable workers.

## Architecture and files

Create `internal/delegation/{capability,preflight,lease,result,reconcile}.go`; `internal/integrations/{codex,claude}/worker.go`; worker templates; `schemas/task-result-v1.schema.json`; run fixtures and delegation journeys. App owns PrepareDelegation/SubmitResult; host adapters only launch/return.

## Contracts and migrations

Activate full `pathframe_prepare_delegation` and `pathframe_submit_result`. Add `pathframe.task-result/v1`, lease/run v1, host capability facts. Existing approved delegated tasks need a compatible manifest and role before launch; otherwise remain blocked, not reassigned.

## Tasks and waves

- S6-T1 (S5): capability, preflight, reason/recovery codes.
- S6-T2 (T1): exclusive task lease and bounded run persistence.
- S6-T3 (T1): result parser/validator and reconciliation states.
- S6-T4 (T1,T2): Codex sequential launcher + worker instructions.
- S6-T5 (T1,T2): Claude sequential launcher + worker instructions.
- S6-T6 (T2-T5): edit guard, no-fallback, end-to-end host journeys.

Tasks T4/T5 may form readiness wave but mutate separate host packages. Context includes packet/schema, host capability fixture, and exact adapter.

## Tests and recovery

Cover success, missing role/capability, oversized context, conflicting lease, malformed result, four worker statuses, lost session, and Brain bypass. Failure keeps delegated policy and reports repair/retry/replan. Crash leaves diagnosable lease; release/retry is explicit. `needs_replan` returns control to Brain. Pinky cannot approve itself or transition plan state.

## Documentation and correction

Document Brain/Pinky authority, advisory scope, declared enforcement, lease lifecycle, and no fallback. Adapter correction never changes task policy.

## Verification

```sh
go test ./internal/delegation ./internal/integrations/codex ./internal/integrations/claude ./internal/app
go test ./tests/journey -run Delegation
go test ./...
```

## Acceptance and gate

Both hosts complete one sequential delegation; every preflight failure occurs before mutation; conflicting leases block; Brain cannot silently edit or inherit authority; Pinky only submits; recovery is executable. Human stop/go before completion semantics.

# Stage 1: Navigator Kernel and Orientation

## Outcome and why now

Deterministic project discovery, lifecycle navigation, persisted journal/projection, and recovery from every non-terminal state. Templates and integrations need stable operations first.

## Scope and non-goals

Implement phases/task states, transition table, root discovery, append-only transition journal, replayed `state.json`, Orient/GetNext and pause/resume/replan/cancel, text/JSON views. No artifact templates, approval, MCP, delegation, or verification execution.

## Architecture and files

Create `internal/workflow/{model,transition,frontier}.go`; `internal/store/{discover,journal,projection}.go`; `internal/app/{orient,navigate,recover}.go`; extend CLI; add `schemas/workflow-v1.schema.json`, `testdata/workflow/`, `tests/journey/navigator_test.go`. `app` coordinates store/workflow; workflow has no filesystem imports.

## Contracts and migrations

Add `pathframe.workflow/v1`; transition event v1; CLI `status`, `next`, `pause`, `resume`, `replan`, `cancel`, and `--json`. Initial `.pathframe` layout is new, so no migration. Illegal transitions return current state, reason, and executable recovery.

## Tasks, dependencies, waves

- S1-T1 (S0): define models, transition table, reason/action codes; unit matrix.
- S1-T2 (T1): project-root discovery and safe contained paths; temp-repo tests.
- S1-T3 (T1,T2): journal append/replay and atomic projection; crash-tail tests.
- S1-T4 (T1,T3): app operations and deterministic frontier.
- S1-T5 (T4): CLI text/JSON projections and goldens.
- S1-T6 (T5): black-box journey for each phase and recovery action.

Waves: T1/T2 after Stage 0 where write scopes do not overlap; then T3; T4; T5; T6. Read architecture, contracts, this stage, touched packages/tests. Return structured task result.

## Test and recovery plan

Table tests cover all legal/illegal transitions and actors. Golden tests prove text derives from canonical JSON. Journeys execute pause/resume/replan/cancel from reachable states. Corrupt projection triggers a diagnostic and rebuild from complete journal records; incomplete last journal line does not invent state. Ambiguous changes require selection.

## Documentation and rollback

Update architecture and CLI reference; document state ownership and recovery. Correction strategy: fix transition/replay logic and rebuild projection. Never rewrite authored files or Git history.

## Verification

```sh
go test ./internal/workflow ./internal/store ./internal/app ./internal/adapters/cli
go test ./tests/journey -run Navigator
go test ./...
go vet ./...
```

## Acceptance and gate

Same input state yields same legal actions; no-arg/status/next answer five orientation questions; schema validates; every non-terminal state completes one recovery journey; projection rebuild is deterministic; no recovery deletes `.pathframe`. Exit only after goldens, full suite, journey, docs, risks, and human stop/go.


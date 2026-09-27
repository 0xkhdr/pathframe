# Stage 8: Doctor and Recovery Completion

## Outcome and why now

Every supported failure can be diagnosed and recovered without delete-and-restart, while preserving authored intent and unaffected completed work.

## Scope and non-goals

Diagnosis catalog, safe projection repair, full replay, abandoned lease release/retry, replan impact, retained cancellation/history, integration/config diagnosis, ambiguous-change selection. No authored-content auto-repair, Git rollback, run compaction, or legacy import.

## Architecture and files

Create `internal/recovery/{diagnose,repair,replay,lease,replan}.go`; `internal/app/doctor.go`; extend CLI/MCP `pathframe_recover`; diagnosis schema/fixtures; corrupt-state and interruption journeys. Recovery reads owning domain validators rather than duplicating rules.

## Contracts and migrations

Add typed diagnosis codes, severity, affected subject, evidence, safe repair operation, human-required flag. `doctor --repair` only rebuilds/removes replaceable machine projections or releases a lease through a recorded event. Journal/run schemas remain sources; migration is deterministic replay.

## Tasks and waves

- S8-T1 (S7): diagnosis catalog mapped to executable recoveries.
- S8-T2 (T1): replay/rebuild and incomplete-tail handling.
- S8-T3 (T1): abandoned lease and interrupted-run recovery.
- S8-T4 (T1): replan impact preserving unaffected completed tasks.
- S8-T5 (T2-T4): Doctor CLI/MCP, repair preview, corrupt fixtures.
- S8-T6 (T5): historical dead-end and session-interruption journeys.

Agents receive one diagnosis family, source invariants, fixtures, and write scope. No task requires reading all recovery cases.

## Tests and recovery

Exercise missing config, malformed artifact, corrupt projection, broken ref, unknown role, stale lease, unavailable capability, failed/interrupted verification, scope amendment, discovery/replan, ambiguous changes, pause/cancel, and Brain bypass. If repair itself fails, preserve sources, report partial machine-file effects, and allow rerun.

## Documentation and correction

Add Doctor catalog, repair ownership, interruption playbook, replan preservation rules. Defective projection can always be discarded and rebuilt by Doctor; authored artifacts and complete journal/run records are never silently edited.

## Verification

```sh
go test ./internal/recovery ./internal/store ./internal/app
go test ./tests/journey -run 'Doctor|Recovery|Interrupted|Replan'
go test ./...
```

## Acceptance and gate

Every named failure has an executable journey; projection rebuild is deterministic; abandoned work resumes/releases; unaffected completions survive replan; cancellation retains history; repair never edits authored intent; no guidance deletes `.pathframe` or rewrites Git. Human declares recovery complete.


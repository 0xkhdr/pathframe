# Stage 2: Templates and OKF Profile

Status: implemented; awaiting human stop/go gate.

## Outcome and why now

Self-describing Quick, Standard, and High-risk artifacts validate into an explicitly approved ready plan. Navigator operations already provide state and recovery.

## Scope and non-goals

Own minimal `okf-markdown/v1`, versioned artifact schemas, task-per-file parsing, references, material-change identity, new/template/check/approve operations. Aido is optional read-only input only. No MCP, host generation, delegation, or broad Markdown engine.

## Architecture and files

Create `internal/artifacts/{profile,frontmatter,parse,validate,templates,identity}.go`; `internal/app/{changes,templates,approval}.go`; embedded `templates/artifacts/`; `schemas/{artifact-instructions,intent,requirements,design,task}-v1.schema.json`; mode fixtures under `testdata/artifacts/`; extend CLI and journeys.

## Contracts and migrations

Add `okf-markdown/v1`, `pathframe.artifact-instructions/v1`, artifact front matter, CLI `new`, `template`, `check`, `approve`. Task policy is explicit; Quick defaults to `brain`. Approval binds validated plan identity. Material edits or policy change require revalidation/reapproval. Existing Stage 1 state with no artifacts remains discoverable and gets actionable initialization guidance.

## Tasks and waves

- S2-T1 (S1): freeze minimal profile and schema fixtures.
- S2-T2 (T1): strict front matter/section parser and round trips.
- S2-T3 (T2): mode/template generator and instruction contract.
- S2-T4 (T2): references, placeholders, acceptance, verification, question validation.
- S2-T5 (T3,T4): Create/Validate/Approve operations and materiality identity.
- S2-T6 (T5): CLI/goldens and cold-agent exercises.

Waves: T1; T2; T3/T4; T5; T6. Agent reads profile/schema fixture plus owned package, not whole repo.

## Tests and recovery

Round-trip/golden fixtures for all modes. Negative cases: placeholder, missing acceptance/verification, broken ref, unresolved question, silent mode lowering, policy edit, stale approval. Failed validation stays planning with precise fixes. Material edit moves to replanning/planning while preserving artifacts/history. Mode increase regenerates only missing artifacts; authored content is never silently overwritten.

## Documentation and correction

Add artifact/profile and approval docs. Profile defects require version-compatible parser fix or new schema version with explicit migration; rollback generator changes without rewriting user-authored content.

## Verification

```sh
go test ./internal/artifacts ./internal/app ./internal/adapters/cli
go test ./tests/journey -run 'Quick|Standard|HighRisk|Approval'
go test ./...
```

## Acceptance and gate

All modes generate only required files; Quick has no unused mandatory artifact; instructions suffice for a cold agent; every plan needs explicit human approval; material edits invalidate approval; Aido absence changes nothing; validation failures have executable recovery. Human approves profile usability before Stage 3.

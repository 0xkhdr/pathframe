# Stage 7: Verification and Task Completion

## Outcome and why now

Submitted work completes only after Pathframe-run mechanical verification bound to repository content and separate semantic acceptance.

## Scope and non-goals

Structured argv runner, contained workdir, required timeout, bounded output, interruption, content identity, advisory scope comparison, accept/request-changes loop, task/change completion/frontier advancement. No shell strings, deployment checks, sandbox, or cryptographic ledger.

## Architecture and files

Create `internal/verification/{command,runner,output,identity,scope}.go`; `internal/app/{verify,review,complete}.go`; run-record extensions and verification fixtures; CLI/MCP adapters; journeys. Runner uses `os/exec.CommandContext` directly and validates workdir under project root.

## Contracts and migrations

Verification command fields: argv array, relative workdir, positive timeout, optional bounded environment policy. Run result records exit status, duration, timeout/interruption, explicit truncation, output, and content identity. Add RunVerification, AcceptTask, RequestChanges app operations. Older run records remain replayable with absent verification meaning “not verified.”

## Tasks and waves

- S7-T1 (S6): command/result models, containment and limit validation.
- S7-T2 (T1): bounded context-aware process runner.
- S7-T3 (T1): deterministic relevant-content identity and changed-file scope comparison.
- S7-T4 (T2,T3): verification orchestration/run persistence.
- S7-T5 (T4): semantic accept/request-changes and frontier completion.
- S7-T6 (T5): CLI/MCP projections and full journeys.

Implementer context contains security rules, runner package, run schema, fixtures; write scope never includes shell wrappers.

## Tests and recovery

Pass, nonzero, timeout, interrupt, oversized output, missing binary, path/symlink escape, stale identity, and scope violation. Failure keeps task active/submitted and offers retry, amend/replan, keep/revert decision as appropriate. Request changes returns task to an executable state. No automatic destructive revert.

## Documentation and correction

Document argv syntax, timeout/output limits, identity meaning, advisory scope, review authority. Runner bug correction invalidates affected verification and reruns; it does not rewrite worker output.

## Verification

```sh
go test ./internal/verification ./internal/app
go test ./tests/journey -run 'Verification|Completion|ChangesRequested'
go test ./... && go vet ./...
```

## Acceptance and gate

No implicit shell path exists; workdir cannot escape; all outcomes persist bounded facts; changed content invalidates stale success; scope violations expose recovery; mechanical success plus semantic acceptance are both necessary; completed task advances frontier. Human stop/go.


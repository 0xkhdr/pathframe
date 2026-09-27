# Stage 9: Full Product Hardening

## Outcome and why now

Sequential Pathframe is installable, secure, measured, documented, and supported only on platforms proven by CI and complete journeys.

## Scope and non-goals

Cross-platform proof, installation/update, security review, performance budgets, clean-machine/onboarding journeys, compatibility matrix, release process, known limitations. Evaluate but do not automatically build a Specd importer. No parallel Pinkies, strict Git evidence, organization controls, deployment orchestration, or broad host support.

## Architecture and files

Extend `.github/workflows/`; add `docs/{install,getting-started,security,limitations,release}.md`; `docs/compatibility.md`; `scripts/` only where Go tooling cannot express packaging; `tests/install/`, security fixtures, benchmarks, and release metadata. Production package changes require a discovered hardening defect, not speculative abstraction.

## Contracts and migrations

No new lifecycle contract expected. Claim Linux arm64, macOS amd64/arm64, or Windows amd64 only after matching build, install, and journey evidence. Release versioning documents compatibility. Initial release explicitly has no Specd reader/detection.

## Tasks and waves

- S9-T1 (S8): threat model and malformed/path/command security suite.
- S9-T2 (S8): benchmarks and budgets for discovery, validation, packet assembly.
- S9-T3 (S8): reproducible install/update/uninstall and clean-machine tests.
- S9-T4 (T3): platform CI matrix and complete sequential journeys.
- S9-T5 (T1-T4): onboarding usability, compatibility matrix, known limits.
- S9-T6 (T5): evidence-based importer decision record and release checklist.

Platform jobs are a readiness wave, not concurrent repository mutation. Each task gets relevant platform/config/test context only.

## Tests and recovery

Race/full tests, malformed inputs, symlink/path traversal, command injection attempts, update rollback, interrupted install, clean machine, and full create/execute/interrupt/recover journey on each claimed platform. Failed platform evidence leaves it a target, not supported. Failed update restores prior binary/config while preserving `.pathframe`.

## Documentation and correction

Finalize installation, onboarding, security, compatibility, limitations, release, and support policy. Correct claims downward immediately if evidence fails; never claim speculative portability.

## Verification

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test -race ./...
go test ./... -run Journey
go test ./... -run Security
go test ./... -run Install
go test ./... -bench 'Discovery|Validation|Packet' -run '^$'
```

CI runs equivalent platform-specific commands and end-to-end journeys.

## Acceptance and gate

Every claimed platform passes clean install and full sequential journey; budgets pass; security findings are resolved/accepted; docs match behavior; Codex/Claude matrix is current; no Specd compatibility exists unless a separately approved evidence-backed task was created; limitations are explicit. Stable release requires human go decision.


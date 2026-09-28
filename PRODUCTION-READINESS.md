# Pathframe Production Readiness Assessment

Date: 2026-09-28  
Audited revision: `0b65aaf54a30595d53cf0ee04b6f720285e332de`  
Audited platform: Linux amd64 with Go 1.26.4

## Executive summary

Pathframe is technically ready for controlled pilots in real Linux amd64 repositories. It is not yet ready to be presented as a generally available production product.

The core implementation is substantially more mature than a prototype. It has a coherent product boundary, deterministic lifecycle rules, human-readable planning artifacts, explicit approval, typed CLI and MCP operations over one application layer, bounded sequential delegation, structured verification, append-only recovery evidence, and focused security controls. Its full race-enabled test suite, journey tests, installer tests, release checks, benchmarks, and portability builds pass locally. Upstream CI also passes for the audited revision.

The remaining blockers are primarily productization and external proof rather than missing lifecycle architecture:

1. No Git tag or GitHub release exists, so the documented download installer cannot currently install `latest` or `v1.0.0`.
2. No license file or declared repository license exists, preventing safe external adoption and redistribution.
3. The repository's release checklist has not been completed with a human go decision, reproducible release artifacts, checksums, and a published tag.
4. Codex and Claude Code integration contracts are tested, but complete journeys have not been demonstrated against actual released host products and their native subagent mechanisms.
5. There is no evidence yet from independent, representative repositories or users.
6. Operational limits—one writer, one sequential Pinky, advisory write scope, no sandbox, and Linux amd64-only production support—must remain explicit during initial adoption.

The recommended path is therefore: fix distribution and licensing first, execute the existing release process without expanding product scope, validate the supported host journeys on real installations, then run bounded pilots across representative repositories. Promote Pathframe from pilot-ready to production-ready only when those gates have durable evidence.

## Product being evaluated

Pathframe defines itself as:

> Pathframe is a local, deterministic development-path protocol that helps humans and coding agents create, understand, execute, interrupt, and recover structured software changes.

Its intended division of responsibility is sound:

- Brain reasons about intent, design, implementation coordination, and semantic acceptance.
- Pathframe records and validates the approved path, packages task context, controls deterministic transitions, executes approved verification, and preserves recovery evidence.
- Pinky performs only work explicitly delegated through an approved task contract.
- A human explicitly approves every transition from planning to ready.

This assessment treats the accepted decisions in `PHILOSOPHY.md` and `ARCHITECTURE.md` as fixed. Production work should harden and prove that design, not replace it with an agent runtime, generic command executor, project-management system, sandbox, parallel worker framework, or Specd compatibility layer.

## Readiness conclusion

### Current classification

| Area | Assessment | Meaning |
| --- | --- | --- |
| Core lifecycle | Ready for pilot | Deterministic phases, task states, approval, interruption, and recovery are implemented and tested. |
| Planning artifacts | Ready for pilot | Quick, standard, and high-risk contracts are validated and material edits invalidate approval. |
| CLI | Ready for pilot | Human-facing operations are complete, no-argument orientation works, and journey coverage exists. |
| Typed MCP | Ready for pilot | Official pinned SDK, local stdio transport, and typed operations are implemented and tested. |
| Delegation | Ready for bounded pilot | One sequential Pinky and authority preservation are implemented; real host execution still needs field proof. |
| Verification and review | Ready for pilot | Structured argv, containment, timeouts, output bounds, freshness, and separate semantic acceptance are present. |
| Recovery | Ready for pilot | Journal replay, projection rebuild, lease recovery, replanning, pause, and cancellation are tested. |
| Security model | Appropriate for pilot | Important boundaries are enforced and non-goals are stated, but Pathframe deliberately is not a sandbox. |
| Performance | Ready for initial use | Current microbenchmarks show negligible protocol overhead. Large-repository evidence is still required. |
| Distribution | Blocked | No release or tag exists, so the advertised network installer has nothing to download. |
| Legal adoption | Blocked | No license exists. |
| Host compatibility | Partially proven | Contract and simulated MCP journeys pass; native released-host journeys are not recorded. |
| External validation | Not proven | No representative pilot reports, compatibility fixtures, or independent usage evidence exist. |

### Overall verdict

Pathframe can add value today when built from a pinned commit and used by informed pilot teams on Linux amd64. It should not yet be advertised as a generally available production release. The correct next work is to close release and evidence gaps, not to add speculative features.

## Evidence collected

### Repository and architecture review

The assessment reviewed:

- `PHILOSOPHY.md` and `ARCHITECTURE.md`;
- repository guidance and accepted decision records;
- CLI and MCP adapters;
- canonical application operations;
- workflow, artifacts, context, delegation, verification, recovery, store, and integration packages;
- public documentation, security boundaries, compatibility limits, installer, CI, and release checklist;
- unit, security, installer, benchmark, and end-to-end journey tests.

The repository was clean before and after the read-only audit. No product code was modified.

### Verification executed

The following checks passed locally:

```sh
test -z "$(gofmt -l .)"
go vet ./...
go test -race ./...
go test ./... -run Journey
go test ./... -run Security
go test ./... -run Install
go test ./... -bench 'Discovery|Validation|Packet' -run '^$'
GOOS=linux GOARCH=amd64 go build ./cmd/pathframe
GOOS=linux GOARCH=arm64 go build ./cmd/pathframe
GOOS=darwin GOARCH=amd64 go build ./cmd/pathframe
GOOS=darwin GOARCH=arm64 go build ./cmd/pathframe
GOOS=windows GOARCH=amd64 go build ./cmd/pathframe
git diff --check
```

The race-enabled suite passed across all packages, including `tests/install` and `tests/journey`.

Observed benchmark results on the audit machine were approximately:

| Benchmark | Result |
| --- | ---: |
| Discovery | 26.7 microseconds/op |
| Plan validation | 168.8 microseconds/op |
| Packet assembly | 1.34 microseconds/op |

These results show that Pathframe's protocol mechanics are not currently a meaningful performance cost. They do not establish behavior for very large plans, repositories, contexts, histories, or run-record sets.

### Upstream evidence

The public repository's latest CI run for the audited revision completed successfully. The repository was created very recently, has no tags or releases, reports no license, and has no meaningful external adoption signal yet. These facts do not weaken the internal implementation evidence, but they do limit claims about production maturity.

## Value to developers

Pathframe addresses several real failure modes in agent-assisted development.

### Durable agreement before implementation

An agent conversation can lose context or silently reinterpret a request. Pathframe binds approval to normalized material plan content. A material plan edit invalidates approval and requires explicit reapproval. This makes the implementation boundary visible and durable without making Git or chat history authoritative.

### Explicit implementation authority

Every task declares `execution_policy`. A delegated task belongs to Pinky; failure does not silently authorize Brain to implement it. This avoids an important class of authority drift where a coordinator starts editing after a worker fails.

### Selected context instead of accumulated context

Task packets contain declared foundation, change, task, and runtime context with required and optional byte-budget behavior. Required context cannot be silently truncated, and optional omissions are reported. This is valuable in long agent sessions and larger repositories where arbitrary context accumulation becomes expensive and unreliable.

### Verification that cannot silently become shell execution

Verification commands are argument arrays, executed directly in a contained project workdir with a required timeout and bounded stdout and stderr. Passing mechanical verification remains separate from Brain's semantic acceptance. Content identity prevents accepting results after relevant files have changed.

### Interruption and recovery as normal behavior

Append-only transition history and bounded run records reconstruct replaceable state. Doctor diagnoses projection damage and abandoned work while preserving authored artifacts and complete evidence. Developers do not need to delete `.pathframe/` or rewrite Git history to recover.

### Shared human and agent view

No-argument orientation and typed MCP operations expose the same canonical application state. Humans can inspect the path using the CLI while agents use typed operations without guessing commands or argument order.

## Where Pathframe should and should not be used

### Strong initial use cases

- Multi-step changes with task dependencies.
- Changes that need explicit human agreement before implementation.
- Work likely to span sessions or be interrupted.
- Agent workflows that delegate a bounded task to one subagent.
- Risky work needing structured rollout and recovery artifacts.
- Work where fresh verification and semantic acceptance must both be recorded.
- Teams evaluating agent behavior and needing a human-readable local evidence trail.

### Poor initial use cases

- Explanations, read-only reviews, brainstorming, or non-development requests.
- Trivial isolated edits whose coordination overhead exceeds their risk.
- Work requiring parallel code-writing workers or automated merge orchestration.
- Environments requiring filesystem containment or hostile-code sandboxing.
- Concurrent multi-process mutation of the same Pathframe state.
- Platforms other than Linux amd64 when production support is required.
- Release or deployment orchestration.

The existing request-assessment rules already encode these exclusions. Production work should preserve them rather than maximizing automatic activation.

## Strengths of the current implementation

### Coherent dependency direction

Adapters call canonical application operations, which coordinate domain packages. MCP dependencies are confined to the MCP adapter. Domain rules are not duplicated in host assets. This is a maintainable boundary and should remain part of release review.

### Deterministic workflow behavior

Transition tables, task frontiers, approval identity, reconciliation, and recovery actions are deterministic. Illegal transitions report state and supported recovery instead of merely failing.

### Human-readable persisted contracts

The Pathframe-owned `okf-markdown/v1` profile is strict enough for deterministic validation while remaining editable and reviewable. Versioning is limited to persisted or cross-process boundaries rather than applied speculatively to internal abstractions.

### Honest host boundary

Pathframe does not pretend to launch or sandbox host-native workers. Host capability declarations, advisory scope, and sequential execution limits are explicit. This is safer than claiming containment that the host does not provide.

### Recovery-oriented persistence

Authored artifacts plus append-only history and run evidence remain the reconstruction sources. `state.json` is disposable. This avoids making a mutable cache or Git state the only authority.

### Focused security controls

The implementation rejects managed-directory and context symlink escapes, path traversal, absolute or non-portable worker paths, uncontained verification workdirs, implicit shell verification, unbounded output, and unsafe overwriting of generated host configuration. The installer refuses symlink destinations and stages replacements atomically.

### Scope discipline

The repository has avoided speculative parallelism, a proprietary agent runtime, generic raw-command MCP access, broad platform claims, and legacy compatibility. That restraint is a production strength.

## Blocking gaps and recommended work

## 1. Legal and governance readiness

### Evidence

There is no root `LICENSE` file and the public repository reports no detected license.

### Risk

External developers and organizations cannot safely assume permission to use, modify, distribute, embed, or contribute to the software. This alone blocks normal open-source production adoption.

### Required decision

A human owner must choose the license. This cannot be inferred by a coding agent. The decision should consider intended commercial use, redistribution, patent terms, contribution expectations, and compatibility with the official MCP SDK and all transitive dependencies.

### Recommended implementation scope

- Add the chosen license text as `LICENSE`.
- Add the license identifier to repository metadata and release artifacts where appropriate.
- Document the contribution licensing expectation in `CONTRIBUTING.md`.
- Generate and review a dependency/license inventory for the release.
- Decide whether a security contact file or private reporting mechanism is needed beyond the current prose.

### Acceptance evidence

- Repository hosting detects the intended license.
- Legal owner explicitly approves the license choice.
- Dependency license review records no incompatible term.
- Packaged source and release artifacts include the license.

## 2. Release engineering and installability

### Evidence

The installer downloads `pathframe_linux_amd64.tar.gz` and its SHA-256 file from GitHub Releases. No tag or GitHub release exists, and the `releases/latest` API returns 404. Therefore the documented one-line installation command does not currently work.

### Risk

Users cannot follow the primary installation instructions. Building from an untagged branch is not a stable supply-chain or rollback practice.

### Recommended implementation scope

Use the existing release checklist rather than creating a new release framework:

1. Choose an initial semantic version.
2. Run the checklist from a clean checkout on native Linux amd64.
3. Build with `-trimpath` and the version linker flag.
4. Produce exactly the documented Linux amd64 archive and checksum.
5. Build twice in the documented environment and compare checksums.
6. Test clean installation, update, failed-update preservation, uninstall, and rollback using the produced artifact.
7. Obtain the required explicit human go decision.
8. Create the tag and GitHub release and attach both files.
9. Test the public installer using both `latest` and the exact version.
10. Preserve prior release assets for rollback.

Avoid adding package managers, signing infrastructure, automatic release publication, or multi-platform archives until the first supported release works and demand justifies them.

### Acceptance evidence

- `curl .../scripts/install.sh | sh` succeeds on a clean supported machine.
- `PATHFRAME_VERSION=<tag>` installs the exact version.
- `pathframe --version` reports the release version.
- The downloaded checksum is verified before extraction.
- A failed or interrupted update preserves the prior executable.
- Release files match the documented names.
- A human-approved release record links the commit, tag, CI result, artifacts, and checksum.

## 3. Native host compatibility

### Evidence

Generated Codex and Claude Code assets, manifests, Doctor checks, MCP startup, and typed tool journeys are tested. Repository journeys connect an MCP client to Pathframe but do not launch actual released Codex or Claude Code products and their native sequential subagents.

### Risk

Host configuration formats, trust prompts, hook behavior, MCP startup conventions, or subagent result behavior may differ from the repository's simulated contract. A green internal journey is necessary but not sufficient for a support claim tied to external products.

### Recommended implementation scope

- Define a small tested compatibility record containing host name, tested host version, Pathframe version, operating system, installation options, and journey result.
- On actual Codex and Claude Code installations, execute:
  - integration installation;
  - host restart and project trust/approval;
  - session orientation when enabled;
  - typed orientation and request assessment;
  - plan creation and explicit human approval;
  - one direct Brain task;
  - one delegated Pinky task;
  - worker failure without Brain authority transfer;
  - lost lease release and retry;
  - verification and semantic acceptance;
  - material edit and reapproval;
  - host Doctor after installation and after deliberate configuration drift.
- Record only observable results and sanitized diagnostics.
- State tested versions in compatibility documentation without claiming an untested stable handshake.

Do not add a proprietary host launcher. Brain must continue to launch native workers, and Pathframe must continue to own only preflight, leases, packets, result validation, verification, and reconciliation.

### Acceptance evidence

- Complete native journeys pass for both claimed hosts on Linux amd64.
- Generated files are accepted by the hosts without manual undocumented modification.
- Doctor detects deliberate drift and provides a workable recovery.
- Delegated failure leaves `execution_policy` unchanged and Brain editing denied.
- The compatibility page names exact tested host and Pathframe versions and test dates.

## 4. Representative repository pilots

### Evidence

Current journeys use controlled temporary repositories and fixtures. No independent real-repository evidence exists.

### Risk

Fixture-based tests may miss scale, repository conventions, nested workdirs, large context, unusual verification commands, multiple changes, long histories, developer comprehension, and recovery friction.

### Recommended pilot matrix

Use at least three repositories with different characteristics:

| Pilot | Characteristics | Primary evidence sought |
| --- | --- | --- |
| Small service or library | Fast tests, simple tree, one language | Onboarding clarity and quick/standard workflow cost. |
| Medium application | Multiple packages and slower verification | Context selection, output bounds, workdirs, interruption, and review freshness. |
| Monorepo or multi-component project | Nested components, shared foundations, multiple commands | Discovery, context budgets, scope reporting, performance, and multi-change orientation. |

Each pilot should include at least one real change that is completed, one interrupted and resumed change, one material replan, one failed verification, one requested-changes cycle, and one Doctor recovery exercise.

Collect:

- repository size and shape without proprietary source content;
- plan mode and task count;
- time spent authoring and approving artifacts;
- packet size and optional omissions;
- number and duration of verification commands;
- recovery events and whether instructions were sufficient;
- false-positive or missed scope findings;
- user confusion and manual workarounds;
- whether the workflow prevented a real mistake;
- whether the overhead was judged worthwhile.

### Acceptance evidence

- At least three representative pilots complete without deleting `.pathframe/` or rewriting Git history.
- Every non-terminal state encountered has a documented executable recovery.
- No pilot requires bypassing canonical operations or manually editing machine-owned state.
- High-impact usability findings are fixed or explicitly documented before general availability.
- The production support statement distinguishes measured evidence from design claims.

## 5. Security and trust-boundary validation

### Evidence

The current security design properly treats paths, artifacts, manifests, journal records, worker results, and verification arguments as untrusted. Security-focused tests pass. The product explicitly does not sandbox Pinky or verification commands.

### Residual risks

- Verification executes repository-provided programs with the invoking user's authority.
- Host-native Pinky can write outside advisory scope unless the host enforces containment.
- Checksums protect download integrity relative to the release metadata but are not publisher signatures.
- Run records are operational evidence, not cryptographic audit logs.
- One-writer assumptions can be violated by users or automation.
- Parser and persisted-boundary robustness has internal tests but no fuzzing or independent review evidence.

### Recommended implementation scope

- Add Go fuzz targets only for high-value trust boundaries: artifact/front-matter parsing, JSONL replay with crash tails, portable worker paths, result decoding, and manifest ownership input. Keep the set small and persist only valuable regressions.
- Run the race suite and fuzz smoke budget in CI if runtime remains bounded.
- Perform a focused manual threat review of installer replacement, managed-directory discovery, symlink races, verification process termination, journal append/replay, and generated host configuration ownership.
- Confirm file permissions for authored and machine-owned state are appropriate for local repositories.
- Add a clear security reporting file or repository configuration if public use begins.
- Consider signed release provenance only after the first release path is stable or if pilot users require it.

Do not claim sandboxing, containment, cryptographic auditability, or protection from malicious repository code.

### Acceptance evidence

- Trust-boundary fuzz targets run for an agreed CI budget without crashes.
- Any discovered crash or traversal case has a regression test.
- Security documentation matches actual enforcement and explicitly retains the no-sandbox warning.
- Installer and release threat review has a recorded human disposition.

## 6. Reliability and concurrency boundary

### Evidence

Pathframe assumes one writer. Concurrent processes may produce a chain Doctor diagnoses as invalid. Sequential delegation is deliberate and accepted.

### Risk

Editors, multiple agent sessions, hooks, or developers may invoke mutating commands concurrently in real repositories. Diagnosis after corruption is weaker than preventing predictable concurrent mutation.

### Recommended approach

First collect pilot evidence. If concurrent mutation is observed, add the smallest cross-process serialization mechanism around canonical mutations, with bounded acquisition and an actionable refusal. It must preserve append-only history and never delete state. Avoid distributed locking or a daemon.

Until evidence justifies implementation:

- state the one-writer rule prominently in getting-started and troubleshooting guidance;
- ensure Doctor identifies the condition accurately;
- include concurrent invocation in pilot exercises;
- avoid hooks that perform mutating operations automatically.

### Acceptance evidence

- Either pilots show the documented one-writer practice is sufficient, or a narrowly scoped local locking decision is recorded and implemented.
- Interrupted writers remain recoverable.
- Lock or concurrency failure messages state the current condition and supported recovery.

## 7. Performance and bounded-storage evidence

### Evidence

Current discovery, validation, and packet benchmarks are fast. Output and run records are bounded, but long-lived histories and large representative repository shapes have not been measured.

### Recommended implementation scope

- Add benchmark fixtures representing many changes, many tasks, long complete journals, bounded run records, and foundation/context files near budget limits.
- Measure wall time and allocations for orientation, validation, replay/Doctor, and packet assembly.
- Define evidence-based budgets rather than arbitrary optimization targets.
- Confirm repeated Doctor and orientation operations remain acceptable after long pilot histories.

Do not add caches, indexes, databases, or background services unless measurements show a user-visible problem. `state.json` should remain a replaceable projection.

### Acceptance evidence

- Representative-scale operations remain within agreed interactive latency and memory budgets on the supported platform.
- Bounds are documented and tested where exhaustion has a recovery implication.
- No optimization changes reconstruction authority or dependency direction.

## 8. Documentation and onboarding

### Evidence

The documentation is unusually complete for the project's age and accurately describes core concepts, CLI operations, integration installation, security, limitations, and recovery. The principal defect is that installation instructions describe releases that do not exist yet.

### Recommended implementation scope

- Until a release exists, mark source build as the only available installation method or publish the release immediately after completing its gate.
- Add one complete, copyable direct Brain journey and one complete delegated journey based on released behavior.
- State where each command must run and which files the user is expected to edit.
- Add a concise troubleshooting map for MCP startup, modified generated configuration, abandoned lease, stale verification, material plan change, and projection damage.
- Link compatibility evidence to exact tested host versions.
- Keep known limitations near adoption instructions rather than buried in architecture material.

Avoid duplicating lifecycle rules across many pages. Each contract should remain owned by its existing documentation page.

### Acceptance evidence

- A new user can install the released binary and complete a quick direct workflow using documentation alone.
- A host user can install the integration, pass Doctor, and complete the native journey.
- Every documented command and asset name is checked against the release binary.
- Documentation never calls advisory scope enforcement or a verification process a sandbox.

## 9. Observability and support evidence

### Current boundary

Pathframe is intentionally local and does not require telemetry. Production readiness does not require adding remote analytics or a hosted service.

### Recommended implementation scope

- Preserve deterministic JSON outputs for automation and support reproduction.
- Ensure errors retain stable reason codes, current state, and supported recovery actions.
- Define a sanitized support bundle procedure that gathers version, platform, workflow orientation, Doctor output, manifest diagnostics, and selected bounded records without silently including authored source context.
- Document exactly what users should review before sharing.
- Establish an issue template for bugs and a private route for vulnerabilities.

### Acceptance evidence

- A pilot failure can be diagnosed using local structured outputs without requiring the user's repository source.
- Support instructions do not ask users to publish secrets or entire `.pathframe/` directories.
- Recovery guidance maps stable codes to supported actions.

## 10. Adoption and value validation

### Risk

A technically correct workflow may still impose more planning overhead than developers will accept. Production value must be demonstrated, not inferred from passing tests.

### Recommended measures

For each pilot, compare appropriate Pathframe changes with the team's prior workflow using:

- setup and planning time;
- implementation restarts caused by lost context;
- number of material misunderstandings found before coding;
- recovery time after interruption;
- verification or acceptance performed against stale work;
- delegated work that exceeded scope;
- developer confidence in current state and next action;
- percentage of requests correctly excluded as trivial or read-only;
- artifact maintenance burden.

Do not optimize for activation count. A successful activation policy should decline work where Pathframe adds no value.

### Acceptance evidence

- Pilot developers can name concrete mistakes prevented or recovery time saved.
- The overhead is acceptable for the target change categories.
- Activation exclusions prevent routine small work from being burdened.
- Negative evidence results in narrower positioning rather than speculative feature expansion.

## Recommended implementation sequence

The sequence below respects existing architecture and keeps the first production release small.

### Gate 0: Human decisions

Owner: human project owner.

- Choose the license.
- Choose the initial semantic version.
- Confirm whether the target is a public open-source release or a restricted internal pilot.
- Define the minimum host versions available for native validation.

Exit condition: legal and release choices are explicit; agents do not guess them.

### Gate 1: Release truthfulness

Owners: documentation and release engineering.

- Add the approved license.
- Make installation documentation match current availability.
- Complete the existing release checklist.
- Publish and test the first Linux amd64 release.

Exit condition: a clean supported machine can install a pinned, licensed release whose checksum and version are verifiable.

### Gate 2: Native host proof

Owners: Codex integration, Claude Code integration, MCP adapter, and journey evidence.

- Run complete journeys on actual supported host releases.
- Fix only observed contract mismatches.
- Record exact compatibility evidence.

Exit condition: both advertised host contracts work end to end without undocumented manual repair.

### Gate 3: Representative pilots

Owners: product evaluation, workflow, context, verification, and recovery domains as findings require.

- Execute the defined repository pilot matrix.
- Exercise interruption, failure, stale verification, reapproval, and Doctor recovery.
- Record user value and friction.

Exit condition: representative real changes complete recoverably and the product demonstrably saves coordination or recovery effort.

### Gate 4: Targeted hardening

Owners: only the domains implicated by evidence.

- Add focused fuzzing and threat-review regressions.
- Address measured scale or concurrency failures.
- Improve onboarding and support diagnostics from pilot findings.

Exit condition: every release-blocking pilot finding is fixed, explicitly accepted, or reflected in a narrower support claim.

### Gate 5: Production declaration

Owner: human project owner.

- Repeat the full release checklist from a clean revision.
- Review security, compatibility, pilot results, and remaining limitations.
- Make an explicit production go/no-go decision.
- Publish the supported matrix and known limitations with the release.

Exit condition: production status is supported by artifacts and dated evidence, not only by repository claims.

## Domain ownership for implementation planning

An implementation-planning agent should route work to existing owners instead of creating a cross-cutting production package.

| Required work | Existing owner or location | Boundary to preserve |
| --- | --- | --- |
| License choice | Human governance | Agent may apply the decision, never choose it silently. |
| Release build and publication | `docs/release-checklist.md`, `scripts/install.sh`, CI/repository release process | Installer touches only the executable; releases do not mutate project state. |
| Installation behavior | `scripts/install.sh`, `tests/install` | Atomic replacement, checksum verification, symlink refusal, rollback preservation. |
| CLI behavior | `internal/adapters/cli` | Projection only; lifecycle rules stay in canonical operations/domains. |
| MCP behavior | `internal/adapters/mcp` | Typed tools only; official SDK remains isolated; no raw command tool or HTTP transport. |
| Host compatibility | `internal/integrations/codex`, `internal/integrations/claude`, `tests/journey` | Brain launches native workers; generated files remain ownership-protected. |
| Workflow correctness | `internal/workflow` | Explicit approval and recoverable transitions remain fixed. |
| Artifact behavior | `internal/artifacts` | Human-readable Pathframe-owned profile; no Specd compatibility. |
| Context scale/findings | `internal/context` | Required context never silently truncates; optional omissions stay visible. |
| Delegation findings | `internal/delegation` | One sequential Pinky; failure never transfers Brain authority. |
| Verification findings | `internal/verification` | Structured argv, contained workdir, timeout, bounded output, no shell. |
| Recovery findings | `internal/recovery`, `internal/store` | Authored artifacts and complete append-only evidence are never machine-rewritten. |
| Security regressions | Owning domain plus focused tests | No unsupported sandbox or audit claims. |
| Documentation | Existing smallest owning page | Avoid duplicated lifecycle truth. |

## Production acceptance checklist

Pathframe should be declared ready for production use only when all applicable items are true.

### Legal and distribution

- [ ] Human-approved license is present and detected.
- [ ] Dependency licenses have been reviewed.
- [ ] A semantic version tag and GitHub release exist.
- [ ] Linux amd64 archive and checksum are published with documented names.
- [ ] Latest and pinned-version installers succeed on clean machines.
- [ ] Interrupted and failed updates preserve the previous executable.
- [ ] Release version is embedded and reported correctly.
- [ ] Release checklist has an explicit human go record.

### Correctness and security

- [ ] Formatting, vet, full tests, race tests, journeys, security tests, installer tests, and builds pass from the release commit.
- [ ] Trust-boundary fuzz or equivalent robustness evidence is recorded.
- [ ] Security documentation matches actual enforcement.
- [ ] No operation claims worker or verification sandboxing.
- [ ] All supported non-terminal states retain an executable recovery path.

### Host compatibility

- [ ] Native Codex journey passes on a named host version.
- [ ] Native Claude Code journey passes on a named host version.
- [ ] Host Doctor detects generated-asset drift.
- [ ] Delegated failure does not alter authority.
- [ ] Compatibility documentation includes dates and exact tested versions.

### Real-world evidence

- [ ] Representative small, medium, and monorepo-style pilots have run.
- [ ] Direct, delegated, interrupted, replanned, failed-verification, and recovered changes have been exercised.
- [ ] No normal recovery deleted `.pathframe/` or rewrote Git history.
- [ ] User value and overhead have been measured.
- [ ] Release-blocking findings are resolved or support claims are narrowed.

### Documentation and support

- [ ] Installation instructions work verbatim.
- [ ] Direct and delegated tutorials work against the release binary.
- [ ] Known limitations are visible before adoption.
- [ ] Security and bug reporting routes exist.
- [ ] Sanitized diagnostic collection is documented.

## Risks that should remain explicit after release

Even after the production gate passes, the following are intentional limits rather than defects:

- Linux amd64 is the only initially supported production platform.
- Only one sequential Pinky is supported.
- Write scope is advisory unless the host explicitly enforces it.
- Pathframe is not a sandbox.
- Verification programs execute with the invoking user's permissions.
- Git and conversation history are not workflow authorities.
- Run evidence is not a cryptographic audit trail.
- Release and deployment orchestration remain outside Pathframe.
- Specd remains incompatible by design.
- Parallel code-writing workers remain deferred.

These limits should be expanded only in response to real demand and representative evidence, with an accepted decision or deviation record where required.

## Final recommendation

Do not redesign Pathframe before release. Its core architecture is coherent, its automated evidence is strong, and its narrow product boundary is a competitive advantage.

The shortest responsible path to production is:

1. obtain the human license and version decisions;
2. publish and verify the first licensed Linux amd64 release;
3. prove the Codex and Claude Code journeys on actual host installations;
4. run representative repository pilots;
5. apply only evidence-driven hardening;
6. repeat the release checklist and make an explicit human production decision.

Until those steps are complete, describe Pathframe as **pilot-ready on Linux amd64**, not production-ready. For the multi-step, interruption-prone, agent-delegated changes it targets, it already has a credible mechanism for reducing ambiguity, authority drift, stale verification, and unrecoverable workflow state. The remaining work is to make that value installable, legally usable, externally demonstrated, and operationally supported.

# Pathframe Workflow Report: `demo-endpoint`

## Outcome

The Laravel change was implemented and its focused feature test passes. Pathframe planning, explicit approval, Brain edit authorization, replanning, and Doctor completed through typed MCP operations. The workflow could not record Pathframe verification or semantic task acceptance because the connected MCP toolset did not expose the required typed operations.

Final observed Pathframe phase: `ready` (`0/1` tasks completed).

## Requested Change

- Create a quick Pathframe change named `demo-endpoint`.
- Use `execution_policy: brain`.
- Add `GET /pathframe-demo`.
- Return HTTP 200 with exact JSON `{"status":"ok"}`.
- Add a focused Laravel feature test.
- Require explicit human approval before implementation.
- Run Pathframe verification with a timeout and exact changed files.
- Semantically accept the task and run Doctor.

## Workflow Timeline

1. Read the repository's Pathframe skill instructions.
2. Called `pathframe_orient` for `demo-endpoint`. It reported that `.pathframe/changes` did not yet exist.
3. Called `pathframe_assess_request`. Result: `must_use`, because Pathframe was explicitly requested.
4. Called `pathframe_create_change` in `quick` mode.
5. Retrieved the typed `intent` and `task` templates with `pathframe_get_template`.
6. Inspected the Laravel route and Pest conventions.
7. Authored:
   - `.pathframe/changes/demo-endpoint/intent.md`
   - `.pathframe/changes/demo-endpoint/tasks/T1.md`
8. Ran `pathframe_validate_plan` with `human_approved: false` until the plan was valid.
9. Presented the plan and stopped for explicit approval.
10. After approval, called `pathframe_validate_plan` with `human_approved: true`.
11. Called `pathframe_get_next`; recommended action was `execute`.
12. Called `pathframe_check_brain_edit`; result: Brain editing explicitly authorized.
13. Implemented the route and focused feature test.
14. Ran the focused test locally with a 120-second process timeout. It passed: 1 test, 2 assertions.
15. Used the typed recovery diagnosis operation as the exposed Pathframe Doctor operation. It found an undeclared task role.
16. Called typed `replan`, changed the task role from `laravel-implementer` to `none`, revalidated, and stopped for renewed explicit approval.
17. After renewed approval, recorded the new approved plan identity and rechecked Brain edit authority.
18. Ran Doctor again. Result: healthy, no diagnoses.
19. Could not call Pathframe verification or acceptance because those operations were absent from the connected MCP capability set.

## Plan Identities

- Initial approved identity: `4fd1bcde6403f1054b37848b52badd604853340fd8ffc6d553a21cd025c475fa`
- Corrected approved identity: `91928e06b562fde95e840b90df4cf93a89eb48f14c79b36f347aef763fc00754`

The second identity superseded the first after changing the undeclared role to `none`.

## Changed Application Files

- `routes/web.php`
  - Added `GET /pathframe-demo`.
  - Returns `response()->json(['status' => 'ok'])`.
- `tests/Feature/PathframeDemoTest.php`
  - Requests `/pathframe-demo` as JSON.
  - Asserts HTTP 200.
  - Asserts the exact JSON object `{"status":"ok"}`.

Pathframe also created and maintained workflow artifacts under `.pathframe/changes/demo-endpoint/`.

## Verification Evidence

Command executed by the coding agent:

```text
timeout 120s php artisan test tests/Feature/PathframeDemoTest.php
```

Result:

```text
PASS Tests\Feature\PathframeDemoTest
1 passed (2 assertions)
Duration: 0.46s
```

This is valid host-side test evidence. It is not recorded as Pathframe verification because `pathframe_run_verification` was unavailable.

## What Went Well

- Activation was deterministic: the request was correctly classified as `must_use`.
- Templates made the expected artifact structure explicit.
- Plan validation caught malformed front matter and invalid references before approval.
- Approval identity prevented silent plan mutation.
- Brain edit authority was explicit and machine-checked before code changes.
- Replanning correctly required renewed human approval after an approved artifact changed.
- Doctor found the undeclared role and was healthy after correction.
- The implementation stayed minimal: one inline route and one focused test; no controller or dependency.
- The test passed with the required behavior and exact response contract.
- No Pathframe CLI workaround was used; all Pathframe state operations used typed MCP calls.

## Blockers and Friction

### Missing typed MCP operations

The installed Pathframe source defines these operations, but the current coding-agent MCP capability set did not expose them:

- `pathframe_run_verification`
- `pathframe_accept_task`
- `pathframe_request_changes`

Runtime inspection showed the verification and acceptance functions as `undefined`. Therefore the agent could not:

- submit the required positive timeout;
- submit the exact changed-file list;
- create a fresh Pathframe verification record;
- perform semantic acceptance;
- advance the task to completed.

This was the final blocking condition. The code works, but Pathframe remains at `ready` with `0/1` tasks completed.

### Role validation occurred too late

The initial task used `role: laravel-implementer`. Plan validation accepted it, but Doctor later reported `unknown_role` because `.pathframe/roles/laravel-implementer.yaml` did not exist. This forced replanning and a second human approval. For a Brain task, `role: none` was sufficient.

### Orientation on a fresh repository returned a raw filesystem error

The first orientation call returned an `lstat .../.pathframe/changes: no such file or directory` error. A fresh, configured project should ideally produce a structured empty-state response rather than a filesystem-level error.

### Temporary projection diagnostic

The typed replan response included `corrupt_projection` while moving into `replanning`. Subsequent plan validation, approval, and Doctor succeeded, so it did not remain an active fault. The diagnostic lacked enough context to distinguish an expected rebuild requirement from actual corruption.

## What Native Pathframe Usage Looks Like to This Coding Agent

Native usage means invoking schema-defined Pathframe MCP tools directly from the agent runtime. The agent passes structured objects and receives structured results. It does not compose shell commands, parse human CLI output, or mutate workflow state files manually.

The successful native sequence was:

```text
pathframe_orient
pathframe_assess_request
pathframe_create_change
pathframe_get_template
pathframe_validate_plan(human_approved=false)
human approval
pathframe_validate_plan(human_approved=true)
pathframe_get_next
pathframe_check_brain_edit
Brain edits
pathframe_recover(action=replan) when required
pathframe_recover(repair=false) for diagnosis/Doctor
```

Expected completion sequence:

```text
pathframe_run_verification(
  change="demo-endpoint",
  task="T1",
  workdir=".",
  timeout_ms=<positive value>,
  changed_files=[
    "routes/web.php",
    "tests/Feature/PathframeDemoTest.php"
  ]
)
pathframe_accept_task(
  change="demo-endpoint",
  task="T1",
  reason=<semantic acceptance reason>
)
Pathframe Doctor
```

The exact verification input schema could not be obtained from the connected toolset, so the agent did not guess or use the CLI as a compatibility bypass.

## Recommendations

1. **Expose the complete workflow toolset atomically.** Planning tools should not be available without verification, review, acceptance, change-request, and Doctor tools from the same compatible version.
2. **Add an MCP capability/version handshake.** Return server version, integration version, tool schema version, and missing required operations during orientation.
3. **Make orientation fail early on incomplete integrations.** Before creating a change, report that the session cannot complete the requested lifecycle.
4. **Validate roles during plan validation.** Reject undeclared roles before approval. Permit or recommend `role: none` automatically for Brain tasks.
5. **Return a structured fresh-project state.** Missing `.pathframe/changes` should mean no active changes, not a raw `lstat` error.
6. **Clarify projection diagnostics.** Include evidence, severity, recovery, and whether the diagnostic is expected during replan.
7. **Provide a dedicated typed Doctor tool.** If `pathframe_recover(repair=false)` is Doctor, name or document it consistently so agents do not need to infer equivalence.
8. **Include exact changed files in verification output.** Echo the submitted list and report additions, omissions, and scope violations deterministically.
9. **Allow capability preflight before approval.** Confirm that the session can execute every required post-approval operation before asking the human to approve.
10. **Keep semantic acceptance separate from test execution.** The current intended design is correct: passing commands are evidence; Brain still checks whether behavior and scope satisfy acceptance.

## Current Risks and Recovery

- Application risk is low: the endpoint is fixed, isolated, and covered by a focused passing test.
- Workflow risk remains: Pathframe has no recorded verification or semantic acceptance for `T1`.
- The workflow should resume only when the missing typed MCP operations are exposed. Then run fresh Pathframe verification with a positive timeout and the two exact application files, semantically review the result, accept `T1`, and run Doctor once more.

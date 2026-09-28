# Pathframe Decision Records

Purpose: define when architecture decisions need records and how humans approve them.

This directory records intentional deviations from the accepted Pathframe plan and genuinely new architecture or scope decisions. It is not a backlog and must not reopen settled choices without direct repository evidence.

## Accepted baseline

The baseline is fixed: Go 1.26 and Linux amd64 first; standard `flag`; local stdio through the official MCP Go SDK; explicit ready approval; explicit `execution_policy` with no delegated fallback; Pathframe-owned `okf-markdown/v1` and optional read-only Aido; one sequential shared-workspace subagent; Pathframe-run structured verification with a required timeout and no implicit shell; authored intent plus `history.jsonl` and bounded runs as reconstruction sources; and no initial Specd compatibility.

## When to write a record

Write a record only when:

- repository evidence directly contradicts an accepted constraint;
- a new decision materially changes architecture or scope;
- two viable choices cannot be deferred without blocking a stage; or
- an intentional implementation deviation is necessary.

Do not write a record for routine implementation details, internal names, or a choice that the active plan already fixes.

## File and status

Use `NNNN-short-title.md`. Start with `proposed`. A human may mark it `accepted`, `rejected`, or `superseded`. Never treat a proposed record as permission to reverse an accepted decision.

## Required content

```markdown
# NNNN: Decision title

Status: proposed
Date: YYYY-MM-DD
Stage: first blocked stage

## Evidence and problem
## Constraints preserved
## Options
## Recommended default
## Consequence of deferral
## Impact and recovery
## Human decision
```

State two or three viable options when possible. Explain operational cost, migration or rollback, affected contracts/files, and how work recovers if the choice fails.

## Legacy code

A record proposing reuse from Specd or another repository must identify exact source files, behavior, tests, license, dependency cost, rejected concepts, and why reuse fits Pathframe boundaries. Copying a subsystem merely because it exists is prohibited.

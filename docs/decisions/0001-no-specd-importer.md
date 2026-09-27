# 0001: Do not build a Specd importer

Status: accepted
Date: 2026-09-28
Stage: 9

## Evidence and problem

Stage 9 requires reconsidering a one-time importer only when real demand and representative fixtures exist. This repository contains neither user demand nor approved Specd fixtures, and the complete Pathframe journey does not require compatibility.

## Constraints preserved

The accepted initial-release decision forbids Specd detection, readers, and compatibility. Pathframe keeps its owned human-readable formats and deterministic recovery sources.

## Options

1. Keep no compatibility surface.
2. Build a one-time read-only-source importer after demand and representative fixtures are supplied.
3. Add continuous compatibility or automatic detection, which is outside approved scope.

## Recommended default

Choose option 1 for the initial release. Option 2 requires a separately approved, evidence-backed task.

## Consequence of deferral

Existing Specd users must author Pathframe artifacts explicitly. No Stage 9 capability is blocked.

## Impact and recovery

No code, schema, or migration is added. If future evidence justifies import, add a new decision and isolated task; failure must leave the source untouched and must not create a partially active change.

## Human decision

The initial no-compatibility decision was already accepted in the product plan. Stage 9 records that no new evidence justifies changing it.

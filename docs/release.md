# Release checklist

Stable release requires an explicit human go decision at the Stage 9 gate.

1. Choose a semantic version and build with `-trimpath` and `-X github.com/0xkhdr/pathframe/internal/app.Version=<version>`.
2. Run formatting, vet, race/full tests, named journey/security/install suites, benchmarks, schema tests, and `git diff --check` exactly as listed in the Stage 9 plan.
3. Run Linux amd64 clean install, update, failed-update preservation, uninstall, and the complete sequential create/execute/interrupt/recover journey.
4. Confirm every supported-platform claim has a native build, install, and journey result. Keep cross-build-only systems as portability targets.
5. Review [security](security.md), [compatibility](compatibility.md), and [limitations](limitations.md); accept or resolve every open finding.
6. Verify Codex and Claude Code generated integration manifests and Doctor checks against the release binary.
7. Produce the Linux amd64 archive and SHA-256 checksum from a clean checkout. Rebuild once and compare checksums before publication.
8. Tag only after the human go decision. Preserve prior artifacts so update rollback remains possible.

If any required evidence fails, fix it and repeat the checklist or reduce the corresponding support claim. Do not publish a stable release from a cross-build result alone.

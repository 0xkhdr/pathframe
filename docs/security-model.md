# Security model

Purpose: define trust boundaries, enforced controls, and protections Pathframe does not provide.

Pathframe treats project paths, authored artifacts, worker results, host manifests, journal records, and verification arguments as untrusted input.

Controls in the initial release:

- managed change identifiers and context paths reject traversal and symlink escape;
- `.pathframe` project creation rejects a symlinked managed directory;
- worker changed-file paths use portable slash syntax and reject absolute paths, traversal, and backslashes;
- verification executes structured argv directly, with a project-contained resolved workdir, required timeout, and bounded stdout/stderr;
- generated host files use ownership hashes and are not overwritten after user modification;
- state repair is limited to machine-owned projections; authored artifacts and complete append-only evidence are preserved;
- the installer stages a complete executable and refuses symlink destinations.

Pathframe is not a sandbox. A verification executable or host-native Pinky has the operating-system permissions of the invoking user, and declared write scope is advisory unless the host explicitly reports enforcement. Review task commands and run only trusted repository content.

Report vulnerabilities privately to the repository maintainers. Include the affected version, platform, reproduction, impact, and whether project data or credentials were exposed. Do not include secrets in a public issue.

# Known limitations

- Linux amd64 is the only supported production platform. Other named targets have compile evidence only.
- Pathframe coordinates one sequential Pinky. Parallel workers, worktree orchestration, and merge policy are not implemented.
- Write scope is advisory unless the host explicitly declares enforcement. Pathframe does not sandbox workers or verification programs.
- Repository mutation assumes one Pathframe writer at a time; concurrent processes may be diagnosed as an invalid journal chain.
- Host configuration installation is ownership-based and refuses to merge unknown existing files.
- Git is observational, not a lifecycle or recovery authority.
- Run records bound output but are not cryptographic audit evidence.
- There is no Specd detection, reader, or importer.
- Release publication and deployment orchestration are outside Pathframe; maintainers produce and verify release artifacts using the checklist.

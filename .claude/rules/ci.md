---
paths:
  - ".github/**"
  - ".devcontainer/**"
  - ".claude/settings.json"
  - "justfile"
  - ".sonarcloud.properties"
---

# CI, the gate, the release and the container

- **The checks** are this repository's own recipes, in the `justfile`, and run in the toolchain stage of `.devcontainer/Dockerfile`, which is spec's (R13). CI calls NovaTrail's shared workflows of `.github`, pinned by SHA with the version beside it, all at one version, and Dependabot moves them: "Secrets", "DCO", "Documents" and "Go" — with `same-bits`, "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules" (R13, R18, R34). Its own jobs are "Proto" — `just ci-proto` and `just ci-release` — and "Gate". The findings of SonarCloud, CodeQL and Scorecard are closed by rule 12 of RUN.md.
- **The gate.** "Gate" runs `just gate` of `sim` in the toolchain stage of `sim`'s devcontainer, on the code of the pull request in place of the commit of learner-model that `sim`'s `gate.json` names, by the tag and the commit of `sim` that `SIM_TAG` and `SIM_COMMIT` of `.github/workflows/ci.yml` name; red when a number leaves its band, for the better or the worse (R20). Neither Dependabot nor a shared workflow raises them: a task does, by the order of releases of R103 spec.
- **The release** comes out by itself through the shared `release.yml` of `.github`, on the line `0.1`, for a commit of `main` whose CI has passed and which changes a path of its `paths` — the code but its tests, the schema, the license, the `justfile`, `.devcontainer/Dockerfile` and `release.yml` itself; documents and tests change no release (R16). Its one file is the schema's archive. Nobody tags by hand.
- **The container is spec's.** `.devcontainer/` and `.claude/settings.json` are what `just -f spec/justfile new-repo-draft learner-model go 3.8 <tag>` builds for the kind `go` at the pinned tag: a move to a new tag compares them byte for byte. gopls comes from its lock file, `.devcontainer/gopls/`, by the kind `go` (R19), and `just ci-lint` holds it to `GOPLS_VERSION`; the plugin `gopls-lsp` of `.claude/settings.json` runs it (R143 spec; R35). No decision records a difference of the container from spec's so far.

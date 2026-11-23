---
paths:
  - ".github/**"
  - ".devcontainer/**"
  - ".claude/settings.json"
  - "justfile"
  - ".sonarcloud.properties"
---

# CI, the gate, the release and the container

- **The checks** are this repository's own recipes, in the `justfile`, and run in the toolchain stage of `.devcontainer/Dockerfile`, which is spec's. CI calls NovaTrail's shared workflows of `.github`, pinned by SHA with the version beside it, all at one version, and Dependabot moves them: "Secrets", "DCO", "Documents" and "Go" — with `same-bits`, "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules". Its own jobs are "Proto" — `just ci-proto` and `just ci-release`, in the image of the stage by the image steps of `.github`'s composite action `toolchain`, without its save and its Go caches, since it runs buf, tar and gzip alone and `go.yml`'s job `toolchain` saves the image on a push to `main` — and the gate. The findings of SonarCloud, CodeQL and Scorecard are closed before the next task, as the rules of RUN.md say.
- **The gate.** "Gate" is the check after four jobs "Gate / part k" on `ubuntu-24.04-arm`: each runs `sim`'s `just gate-part k 4` and "Gate" its `just gate-check bands`, in the toolchain stage of `sim`'s devcontainer, with `sim`'s cache of Go, on the code of the pull request in place of the commit of learner-model that `sim`'s `gate.json` names, by the tag and the commit of `sim` that `SIM_TAG` and `SIM_COMMIT` of `.github/workflows/ci.yml` name; red when a number leaves its band, for the better or the worse. Neither Dependabot nor a shared workflow raises them: a task does, by spec's order of releases of the model.
- **The release** comes out by itself through the shared `release.yml` of `.github`, on the line `0.1`, for a commit of `main` whose CI has passed and which changes a path of its `paths` — the code but its tests, the schema, the license, the `justfile`, `.devcontainer/Dockerfile` and `release.yml` itself; documents and tests change no release. It builds the stage `toolchain-core`, which `toolchain` adds nothing to. Its one file is the schema's archive. Nobody tags by hand.
- **The container is spec's.** `.devcontainer/` and `.claude/settings.json` are what `just -f spec/justfile new-repo-draft learner-model go 3.8 <tag>` builds for the kind `go` at the pinned tag: a move to a new tag compares them byte for byte. gopls comes from its lock file, `.devcontainer/gopls/`, by the kind `go`, and `just ci-lint` holds it to `GOPLS_VERSION`; the plugin `gopls-lsp`, which `CLAUDE_CODE_PLUGIN_DIRS` of `.devcontainer/devcontainer.json` loads from `spec/claude-plugins/`, runs it. No decision records a difference of the container from spec's.

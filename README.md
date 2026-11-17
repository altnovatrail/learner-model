# learner-model

[![CI](https://github.com/altnovatrail/learner-model/actions/workflows/ci.yml/badge.svg)](https://github.com/altnovatrail/learner-model/actions/workflows/ci.yml)
[![CodeQL](https://github.com/altnovatrail/learner-model/actions/workflows/codeql.yml/badge.svg)](https://github.com/altnovatrail/learner-model/actions/workflows/codeql.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/altnovatrail/learner-model)](go.mod)
[![Coverage](https://codecov.io/gh/altnovatrail/learner-model/branch/main/graph/badge.svg)](https://codecov.io/gh/altnovatrail/learner-model)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/altnovatrail/learner-model/badge)](https://scorecard.dev/viewer/?uri=github.com/altnovatrail/learner-model)
[![License](https://img.shields.io/github/license/altnovatrail/learner-model)](LICENSE)

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)

The NovaTrail learner model: how well a learner has mastered each skill, estimated from the evidence of attempts. It is a Go library of pure computation with no input or output of its own: `api` hands it the skill graph and the evidence, and the engine reads its estimates to pick the next lesson or task. Its model is a Rasch model with a guessing floor, updated Elo-style, with the level kept for each module, forgetting, evidence weighed by its outcome and provenance, and a cautious "mastered". The evidence log is the source of truth and the state a cache of it, recomputed bit for bit on any build, WebAssembly included.

A component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It comes together with engine and sim: no change to the model goes in without passing sim's gate.

## Releases

Each change of the code or the schema merged to `main` comes out by itself as the next patch, with its tag. Go code takes the library as a Go module, and the release's one file, the schema of the evidence and the state, is verified with spec's recipe ([docs/README.md](docs/README.md), "Releases"):

```sh
go get github.com/altnovatrail/learner-model@<tag>
gh release download <tag> -R altnovatrail/learner-model
just -f spec/justfile release-verify learner-model <tag> learner-model-proto-<tag>.tar.gz
```

## Documents

- [docs/README.md](docs/README.md) — the index of the documents, the detail of the releases and how to work on the repository.
- [docs/architecture.md](docs/architecture.md) — the component and how it is built; [docs/decisions.md](docs/decisions.md) — the decision log.
- [RUN.md](RUN.md) — the work plan; [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md) — the working rules.

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

The NovaTrail learner model: how well a learner has mastered each skill, estimated from the evidence of attempts, and each skill's state. It is a Go library of pure computation with no input or output of its own. Its model is a Rasch model with a guessing floor, updated Elo-style, with the level kept for each module, forgetting, evidence weighed by its outcome and provenance, and a cautious "mastered". The evidence log is the source of truth and the state a cache of it, recomputed bit for bit on any build, WebAssembly included.

It is a component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring:

- `api`, the platform's service, calls it when an attempt ends: it hands it the evidence and the skill graph and keeps the evidence log and the state;
- `engine` reads its estimates to pick the next lesson or task;
- `sim`, the learner simulator, judges every change of the model before it is released.

Its card is in the architecture of [spec](https://github.com/altnovatrail/spec); the detail, the decisions and how to work on it are in [docs/](docs/README.md).

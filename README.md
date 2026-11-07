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

The NovaTrail learner model: how well a learner has mastered each skill, estimated from the evidence of attempts — the outcome of each attempt, weighed by that outcome and by the provenance of the task's key. It is a Go library of pure computation with no input or output of its own: `api` hands it the skill graph and the evidence, the time comes from outside and it takes no seed, needing no randomness, and the engine reads its estimates to pick the next lesson or task. Its starting point is the learner model of MathTrail's standalone — a Rasch model with a guessing floor, updated Elo-style, with a trial series and a cautious "mastered" — with the level kept for each module rather than for the whole learner model ([R03](docs/decisions.md)), with forgetting added as each skill's retention, whose stability grows with spaced practice ([R04](docs/decisions.md)), and with evidence weighed by its outcome, the provenance of its key and the series it belongs to: each item's first verdict counts, a hint only withholds "mastered", and a page of many items counts for less than their sum ([R06](docs/decisions.md)). Whether a learner has an anchor skill of the shared taxonomy, which links one module to another, it reads from the mastered skills that teach it ([R07](docs/decisions.md)). The log of evidence is the source of truth and the state a cache of it: each evidence record carries every observation that recompute reads, appending a record gives the recomputed state bit for bit on any build, and a retired skill's evidence counts in full for each skill that replaced it ([R08](docs/decisions.md)). Its code imports nothing but Go's standard library and gives the same bits on every build, WebAssembly included, so that it could also run on the device ([R09](docs/decisions.md)), and it encodes evidence and state in protobuf's binary form by a schema of its own ([R10](docs/decisions.md)).

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.9.5), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, engine and sim; it comes together with engine and sim, since the learner simulator needs both the model and the rule of choice, and no change to the model after its first version goes in without passing it.

## Status

Detailing is done: the first epic of [RUN.md](RUN.md) decided what the component's card left open, and section 3 of [docs/architecture.md](docs/architecture.md) is what to build. One question came later with О-25 — how much an answer typed in place of handwriting weighs — and spec answered it: as much as handwriting (section 2.7). The code comes in epic E2 of RUN.md, and so far there is its skeleton and the first part of the model: the Go module; its types — the graph, the task, the evidence, the state and what reading gives — with the schema of the evidence and the state in `proto/` and the codec of their bytes; its package `detmath` — exp, log, the logistic function, a power and the retention of a skill, the same bits on every build, which the engine may take too; the level of each module and the mastery of each skill — `BandByGrade`, `Append` and `Recompute`, the step of one piece of evidence and the trial series, held to standalone's test vectors ([R03](docs/decisions.md)); and its checks — `just ci-lint`, `ci-test`, `ci-proto`, `ci-vuln` and `ci-licenses`, which CI runs too ([R13](docs/decisions.md)). `ci-test` holds the code to the same bits on every build from the start: a test of the rules of code, the tests on amd64 with and without fusion, on 386, in js/wasm and on arm64, and a comparison of the assembler with fusion and without it ([R09](docs/decisions.md)). Its first version, checked by its own tests, is to be the baseline the simulator is built on, released signed before the simulator's first tag, and every change after it passes the simulator against a criterion written down before the first candidate ([R12](docs/decisions.md)).

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, what was deliberately deferred, and what to build.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.9.5), checked out here as a git submodule at `v0.9.5`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).

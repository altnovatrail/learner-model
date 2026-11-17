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

A component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, engine and sim; it comes together with engine and sim, since the learner simulator needs both the model and the rule of choice, and no change to the model after its first version goes in without passing it.

## Releases

Releases come out by themselves: each commit of `main` whose CI has passed and which changes the code, the schema or what builds the release's file gets the next patch of the line `0.1` — the first, `v0.1.0`, came out on 2026-11-07 — with its tag, by .github's shared release workflow ([R16](docs/decisions.md)). Nobody tags by hand, and a release is immutable.

- **Go code** takes the library as a Go module, by the tag: `go get github.com/altnovatrail/learner-model@v0.1.0`. `go.sum` and the Go checksum database hold the version from then on.
- **The release's one file** is the schema of the evidence and the state, for what reads them without Go: `learner-model-proto-<tag>.tar.gz`, with `LICENSE`, `buf.yaml` and the `.proto` files of `proto/`, signed without a key. Its `.proto` files are the ones api puts into a learner's export, under `api/proto/novatrail/learnermodel/v1/`, beside the JSON of the evidence log and the state they describe, so their comments say what each field is without the code (R135 spec). Verify it where it enters NovaTrail with spec's recipe, `just -f spec/justfile release-verify learner-model <tag> learner-model-proto-<tag>.tar.gz`, with its bundle `<file>.sigstore.json` beside it. `just release-files <tag>`, in the toolchain stage of the tag's `.devcontainer/Dockerfile`, builds the same bytes from the tag's tree.

## Documents

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, what was deliberately deferred, and how it is built.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [docs/sources.md](docs/sources.md) — the external facts the documents rest on, each with its source, its date and its quote.
- [docs/measurements.md](docs/measurements.md) — the results of measurements.
- [RUN.md](RUN.md) — the work plan, in epics.
- `run/` — the plan's other tables: [the epics for later](run/EPICS.md), [the open questions](run/QUESTIONS.md), [the technical debt](run/DEBT.md) and [the tasks outside an epic](run/TASKS.md).

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec), checked out here as a git submodule: clone this repository with `git clone --recurse-submodules`.

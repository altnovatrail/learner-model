# learner-model

[![CI](https://github.com/altnovatrail/learner-model/actions/workflows/ci.yml/badge.svg)](https://github.com/altnovatrail/learner-model/actions/workflows/ci.yml)
[![CodeQL](https://github.com/altnovatrail/learner-model/actions/workflows/codeql.yml/badge.svg)](https://github.com/altnovatrail/learner-model/actions/workflows/codeql.yml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/altnovatrail/learner-model/badge)](https://scorecard.dev/viewer/?uri=github.com/altnovatrail/learner-model)
[![License](https://img.shields.io/github/license/altnovatrail/learner-model)](LICENSE)

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_learner-model&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model)

The NovaTrail learner model: how well a learner has mastered each skill, estimated from the evidence of attempts — the outcome of each attempt, weighed by that outcome and by the provenance of the task's key. It is a Go library of pure computation with no input or output of its own: `api` hands it the skill graph and the evidence, the time and the seed come from outside, and the engine reads its estimates to pick the next lesson or task. Its starting point is the learner model of MathTrail's standalone — a Rasch model with a guessing floor, updated Elo-style, with a trial series and a cautious "mastered" — with the level kept for each module rather than for the whole learner model ([R03](docs/decisions.md)), and with forgetting added.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.6.3), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, engine and sim; it comes together with engine and sim, since the learner simulator needs both the model and the rule of choice, and no change to the model goes in without passing it.

## Status

Detailing. The documents come first: the first epic of [RUN.md](RUN.md) decides what the component's card left open, and the code follows in the order RUN.md sets. There is no code here yet: it comes with engine and sim, and the simulator's criterion is written down before the first candidate.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, and what is still to decide.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.6.3), checked out here as a git submodule at `v0.6.3`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).

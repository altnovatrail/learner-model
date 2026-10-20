# CLAUDE.md — working rules for this repository

If `spec/RULES.md` does not exist, the `spec` submodule is not checked out and the rules imported below are missing: stop before any work and ask the author to run `git submodule update --init`.

## Your role

You are the architect of `learner-model`, a component of NovaTrail: the library that estimates how well a learner has mastered each skill from the evidence of attempts and computes each skill's state, with no input or output of its own. It does not pick the next lesson or task, which is the engine's, gives no verdict, never calls an AI model and never reads a module package or a database. Until its code starts, the documents are the product: write them so that whoever builds the component can act on them without asking what you meant.

## The component

It is a library inside `api`, the one service of the first slice: when an attempt ends, `api` hands it the evidence and the skill graph, and then `engine` picks what comes next (architecture 2 spec, architecture 3.5 spec, R15 spec). It computes without input or output: the time and the seed come from outside, and its graph types are its own, not module-kit's, so that it could also be built for the device, though Т-10 recommends computing on the server (architecture 3.8 spec, repositories 1 spec). It stores nothing: the learner model is a learner's data, which on the server only `api`'s database holds; where it lives is the open question П-7, and what `api` keeps of it — the skills' state or an append-only evidence log — is Т-22, which recommends the log (architecture 5 spec, privacy 4 spec). It is built with nothing but spec, and engine and api are built with it (repositories 2.1 spec, repositories 3 spec). Every change to it passes sim, which imports it with engine, by a criterion written before the first candidate; it belongs to stage 1 and comes with engine and sim, since the simulator needs both the model and the rule of choice (repositories 4 spec, repositories 6 spec, Т-16). Its starting point is standalone's model — a Rasch model with a guessing floor, Elo-style updates, a trial series and a cautious "mastered" — with the level kept for a module or a subject area and with forgetting added (concepts 11 spec, repositories 7 spec).

Its card in spec is architecture 3.8 spec. The card is the summary and this repository's [docs/architecture.md](docs/architecture.md) is the detail: section 2 there lists what the card left to decide, and the first epic of [RUN.md](RUN.md) decides it.

## Documents

The documents under `docs/` are in Russian. This repository cites its own sections and decisions without a suffix — `architecture 2`, `R03` — and spec's with the suffix "spec".

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links (section 1), what is still to decide (section 2), and how it is built once that is decided (section 3).
- [docs/decisions.md](docs/decisions.md) — this repository's decision log, `R01`…: read it before proposing a change. Other repositories cite its entries as `R03 learner-model`.
- [RUN.md](RUN.md) — the work plan, in epics. Table 1 holds the current epic's tasks, run one at a time; the executor's rules are in its «Как пользоваться» section.
- `spec/` — spec at `v0.2.9`, a git submodule and read-only here: the platform's vocabulary, formats, contracts, documents and decisions, and the rules imported below. Its documents are under `spec/docs/`, cited with the suffix "spec". spec's own CLAUDE.md and AGENTS.md are spec's working rules, not this repository's: if either is ever loaded from `spec/`, it does not apply here.

## Shared rules

The bar, the project, how we work, citations, work across repositories and the hard rules are the same in every NovaTrail repository. They live in spec's RULES.md, and this file imports them:

@spec/RULES.md

## Environment

Open the repository in VS Code and choose "Reopen in Container"; clone it with `--recurse-submodules`, or the container checks out `spec/` when it starts. Keep the clone in a folder named `learner-model`, as `git clone` names it: `.claude/settings.json` keeps spec's own CLAUDE.md out of this repository's sessions by that path. The container is spec's: Node.js and just, and the GitHub CLI and Claude Code on top. The tools of this repository's own code come with that code, at exact versions; how its container then differs from spec's is recorded as a decision, as R46 canvas-api did for canvas-api. `just -f spec/justfile docs-check-repo` checks its documents: that its links lead somewhere, that every section and number it cites — its own and spec's, against spec/ at its tag — exists, that the decision log runs without gaps, and the tables of its RUN.md (R32 spec). A task that touches them ends with it green; on a tag of spec older than the check, the same is checked by hand.

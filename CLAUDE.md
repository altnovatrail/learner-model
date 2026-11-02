# CLAUDE.md — working rules for this repository

If `spec/RULES.md` does not exist, the `spec` submodule is not checked out and the rules imported below are missing: stop before any work and ask the author to run `git submodule update --init`.

## Your role

You are the architect of `learner-model`, a component of NovaTrail: the library that estimates how well a learner has mastered each skill from the evidence of attempts and computes each skill's state, with no input or output of its own. It does not pick the next lesson or task, which is the engine's, gives no verdict, never calls an AI model and never reads a module package or a database. Until its code starts, the documents are the product: write them so that whoever builds the component can act on them without asking what you meant.

## The component

It is a library inside `api`, the one service of the first slice: when an attempt ends, `api` hands it the evidence and the skill graph, and then `engine` picks what comes next (architecture 2 spec, architecture 3.5 spec, R15 spec). It computes without input or output: the time and the seed come from outside, and its graph types are its own, not module-kit's, so that it could also be built for the device, though Т-10 recommends computing on the server (architecture 3.8 spec, repositories 1 spec). It stores nothing: the learner model is a learner's data, and it lives in `api`'s database as an append-only evidence log, the source of truth, with the skills' state as a cache marked with the version of learner-model and the release of modules. What an evidence record holds and how much it weighs is this library's to decide, and `api` keeps it as it is; the one condition is that the state an appended record gives equals the state recomputed from the log (R80 spec, which closed П-7 and Т-22; privacy 4 spec). An evidence record carries every observation that recompute reads, in a form that only grows; the state is a fold of the log, so appending gives the recomputed state bit for bit on any build; and a retired skill's evidence counts in full for each skill that replaced it (R08). It is built with nothing but spec, and engine and api are built with it (repositories 2.1 spec, repositories 3 spec). Every change to it passes sim, which imports it with engine, by a criterion written before the first candidate; it belongs to stage 1 and comes with engine and sim, since the simulator needs both the model and the rule of choice (repositories 4 spec, repositories 6 spec, Т-16). In spec's cross-repository plan, COORDINATION.md, its detailing needed only the move to spec's line 0.6, and is under way; its code waits for sim's gate, the stand's first tag, which comes after the first modules, and api's code waits for it (R76 spec). Its starting point is standalone's model — a Rasch model with a guessing floor, Elo-style updates, a trial series and a cautious "mastered" — with the level kept for each module (R03), with forgetting added as each skill's retention, whose stability grows with spaced practice (R04), and with evidence weighed by outcome, key provenance and series (R06; concepts 11 spec, repositories 7 spec). Whether a learner has an anchor skill of the shared taxonomy it reads, at read time, from the mastered skills that teach it (R07).

Its card in spec is architecture 3.8 spec. The card is the summary and this repository's [docs/architecture.md](docs/architecture.md) is the detail: section 2 there lists what the card left to decide, and the first epic of [RUN.md](RUN.md) decides it.

## Documents

The documents under `docs/` are in Russian. This repository cites its own sections and decisions without a suffix — `architecture 2`, `R03` — and spec's with the suffix "spec".

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links (section 1), what is still to decide (section 2), and how it is built once that is decided (section 3).
- [docs/decisions.md](docs/decisions.md) — this repository's decision log, `R01`…: read it before proposing a change. Other repositories cite its entries as `R03 learner-model`.
- [RUN.md](RUN.md) — the work plan, in epics. Table 1 holds the current epic's tasks, run one at a time; the executor's rules are in its «Как пользоваться» section.
- `spec/` — spec at `v0.6.6`, a git submodule and read-only here: the platform's vocabulary, formats, contracts, documents and decisions, and the rules imported below. Its documents are under `spec/docs/`, cited with the suffix "spec". spec's own CLAUDE.md and AGENTS.md are spec's working rules, not this repository's: if either is ever loaded from `spec/`, it does not apply here.

## Shared rules

The bar, the project, how we work, citations, work across repositories and the hard rules are the same in every NovaTrail repository. They live in spec's RULES.md, and this file imports them:

@spec/RULES.md

## Environment

Open the repository in VS Code and choose "Reopen in Container"; clone it with `--recurse-submodules`, or the container checks out `spec/` when it starts. Keep the clone in a folder named `learner-model`, as `git clone` names it: `.claude/settings.json` keeps spec's own CLAUDE.md out of this repository's sessions by that path. The container is spec's: Node.js, Go, buf and just, and the GitHub CLI and Claude Code on top. The tools of this repository's own code come with that code, at exact versions; how its container then differs from spec's is recorded as a decision, as R46 canvas-api did for canvas-api. `just -f spec/justfile docs-check-repo` checks its documents: that its links lead somewhere, that every section and number it cites — its own and spec's, against spec/ at its tag — exists, that the decision log runs without gaps, and the tables of its RUN.md (R32 spec). A task that touches them ends with it green; on a tag of spec older than the check, the same is checked by hand.

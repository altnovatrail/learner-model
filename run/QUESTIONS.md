# QUESTIONS.md — questions to the author of `learner-model`

Table 3 of [RUN.md](../RUN.md): the author's open questions that tasks leave.

- Never read whole: take the rows the task and its sections name — `grep -n '^| В1 ' run/QUESTIONS.md`.
- A row is a question a task leaves open (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main`. It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В60.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В59 | «В реестр spec»: may the register R of a repository's decision log have gaps, as the tables В and Д already do (`gaps: true` in `spec/tools/docs-check/config.mjs`)? Now its numbers go without gaps and every R cited anywhere, the archives `docs/plan/` too, must exist, so no entry can leave a log: here R02, R05, R11, R15, R17, R19, R23, R26, R27, R33, R38, R44, R46 and R48 of [docs/decisions.md](../docs/decisions.md) stay as one sentence each that holds no decision (Д24). Proposed: a number that left is never given again; the archives may cite it, as they cite a removed record of sources and measurements (`spec/tools/docs-check/sections.mjs`); `docs/`, RUN.md and `run/` still cite only entries that exist | В58 | (a) `gaps` for R, and the archives' exception for R; (b) as now: an entry that holds nothing stays as a stub; (c) a section "Retired" in the log, where an entry leaves its text and keeps its number | (a): "no gaps" and "an archive's citation exists" give almost nothing — a removal shows in the diff, and a log only grows by new numbers —, while no reuse and existing citations in the live documents stay; a stub costs every session that reads the log |

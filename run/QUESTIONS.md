# QUESTIONS.md — questions to the author of `learner-model`

Table 3 of [RUN.md](../RUN.md): the author's open questions from tasks' «Открыто» (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| В1 ' run/QUESTIONS.md`.
- A row comes from a task's «Открыто» (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main` (R81 spec; RULES.md, Across repositories). It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В56.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В54 | «В реестр spec»: answered by the author on 2026-11-16, through the orchestrator: (a) — T68 ([docs/plan/E8.md](../docs/plan/E8.md)); the row leaves once `spec` records it. the order of row 4.20 and row 4.21 of `spec`'s `COORDINATION.md` against the order of releases of R103 spec. T68 changes what the model computes (R45), so "Gate", on the bands `sim` `v0.1.25` recorded on `learner-model` v0.1.8, is red on its pull request — 12 numbers out of their bands, measurements 9.2 — and holds the merge and the release; R103 spec has `sim` record the change at the commit of its pull request and cut a tag first, the pull request then raise `SIM_TAG` and `SIM_COMMIT` and be merged, while row 4.21, `sim`'s short run of the three rules, waits for row 4.20's release | T68 | (a) row 4.21 runs on the commit of T68's pull request: its short run by the criterion written before it, then the bands of the gate recorded on that commit and a tag of `sim`; T68's pull request raises the tag and is merged and released; 4.21 depends on 4.20's pull request, not its release. (b) A row of `sim` of its own before 4.21 records the bands on the commit and cuts the tag, and 4.21 runs after the release as planned — one more task of `sim`. (c) Release T68 without the gate — against R12 and R103 spec | (a): R103 spec's order as it stands, no extra task, and `sim` judges the three rules before they are released, as R20 asks of a change of what the model computes |
| В55 | «В реестр spec»: row 4.20 of `spec`'s `COORDINATION.md`, the three rules of R157 spec in the model, is done here, but `spec`'s plan at v0.18.14 still has it `[ ]`. T68 built them (R45): the start a band's width below the grade's band, an adult's correction as evidence of kind 3, the guard after three failures in a row, and the level's band in `Read` for T178 spec; `sim` judged them on the commit of the pull request (row 4.21, В54), and pull request #84 was merged and released as `learner-model` v0.1.13, the gate on `sim` `v0.1.28` | T70; R45, R46 | (a) `spec` marks row 4.20 `[x]` with T68 and v0.1.13, and row 6.28 of `api` takes `Read`'s band from v0.1.13; (b) `spec` keeps it open until `api` takes v0.1.13 | (a): the row is the library's work, which is merged and released; what `api` does with it is row 6.28 of its own |

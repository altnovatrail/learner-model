# QUESTIONS.md — questions to the author of `learner-model`

Table 3 of [RUN.md](../RUN.md): the author's open questions from tasks' «Открыто» (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| В1 ' run/QUESTIONS.md`.
- A row comes from a task's «Открыто» (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main` (R81 spec; RULES.md, Across repositories). It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В52.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В51 | The budget of the runs of the exploration T61 of E8, the round of fairness by the grade, and the machine's time. A decisive run of `sim` is 44 populations, 68 cells, 2,000 learners a cell: the base of `fairness-1` took 43 min 28 s by the clock on a machine free at its start (R58 sim), and a busy machine takes longer. The cells `D-grade1`, `D-grade3` and `D-grade5` alone, which T59 runs for the analysis, are about 2% of the CPU time of a decisive run (R55 sim). In T43 and T52 the author gave about ten runs, the base first (R32; T52, the author's decision of 2026-11-14) | T58 | (a) about ten decisive runs, the base first, with a screening of the sets on the cells `D-` alone before them; (b) about six decisive runs, the base first; (c) up to fifteen, if T59 finds more than one cause | (a): the same budget gave the judgement `model-2` its candidate number for number (R32), and a screening on the cells `D-` costs little and drops the sets that do not move `M11.same` |

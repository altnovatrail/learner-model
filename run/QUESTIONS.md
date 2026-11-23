# QUESTIONS.md — questions to the author of `learner-model`

Table 3 of [RUN.md](../RUN.md): the author's open questions that tasks leave.

- Never read whole: take the rows the task and its sections name — `grep -n '^| В1 ' run/QUESTIONS.md`.
- A row is a question a task leaves open (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main`. It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В59.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В58 | «Решено исполнителем по рекомендации, до проверки автора»: R01, R02, R05, R11, R15, R17, R19, R23, R26, R27 and R33 of [docs/decisions.md](../docs/decisions.md) record the moves to `spec` v0.2.9, v0.6.3, v0.6.6, v0.9.3, v0.9.5, v0.9.9, v0.9.16, v0.12.0, v0.14.0, v0.15.0 and v0.16.3, and their headings say the repository stands on those tags, while the log's opening, as the template has it, says a move to a tag is its commit with no entry, and the repository stands on v0.19.12. T82, T84, T85, T86 and T87 rewrote them by the rule of citations alone, their meaning kept | T82, T84, T85, T86, T87 | (a) keep them as rewritten; removing or retiring them changes the log and is a task of its own on the author's word; (b) remove the eleven now: what they hold is the history of the moves, which their commits keep; (c) keep them, with headings in the past tense, "stood on" | (a): the task edits citations and keeps the meaning, and what the eleven still hold — why the container is `spec`'s, why a move takes the latest release — stays until the author decides |

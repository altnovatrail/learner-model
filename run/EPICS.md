# EPICS.md — epics for later of `learner-model`

Table 2 of [RUN.md](../RUN.md).

- A row has a result and a basis, no tasks yet. Row order is the order epics are taken, until the author decides otherwise.
- Read by the choice and planning of the next epic, not by each task. Once an epic is split into tasks, they go to table 1 and the body of RUN.md, and its row leaves (RUN.md, Epics).

## Table 2. Epics for later

| ID | Epic | Result | Basis |
|---|---|---|---|
| E6 | The "mastered" threshold paired with `engine` | The threshold "mastered" and the corridor's target (Д16), and with them the failure boundary ½, which `engine` reads too (Д7), as a pair candidate: judged by `sim`'s pair criterion, with the pair's release order by R103 spec — chosen and released, or the exit recorded. Not yet (T71, 2026-11-17; sources 6.15): `engine` has no current epic and no candidate of the corridor's target — its E5 ended without a judgement, and its table 2 is empty until a reason comes —; `sim` lays out no round by plan (R59 sim) and runs no long fitting of constants while there are fewer than 100 children (R69 sim; Д261 spec); item 4 of `spec`'s `COORDINATION.md` has no row for it. Opens when all three hold: a reason of its own for the threshold, named by this repository or by `engine` — a fault the gate or a measurement finds, a need of the product such as Д16's risk on a choice of two or three or a free answer, or a row of `spec`; `engine` proposes its target for the pair; `sim` can judge it — whether R69 sim's stop covers a pair's judgement is the author's, asked then, unless there are 100 children. Then this repository plans E6 and writes a row «В реестр spec» for the parts of `engine` and `sim` | Д7, Д16; R103 spec; R17; R30; R42; R59 sim, R69 sim; sources 6.15 |

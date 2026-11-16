# EPICS.md — epics for later of `learner-model`

Table 2 of [RUN.md](../RUN.md) (R153 spec, R154 spec).

- A row has a result and a basis, no tasks yet. Row order is the order epics are taken, until the author decides otherwise.
- Read by the choice and planning of the next epic, not by each task. Once an epic is split into tasks, they go to table 1 and the body of RUN.md, and its row leaves (RUN.md, Epics).

## Table 2. Epics for later

| ID | Epic | Result | Basis |
|---|---|---|---|
| E6 | The "mastered" threshold paired with `engine` | The threshold "mastered" and the corridor's target (Д16), and with them the failure boundary ½, which `engine` reads too (Д7), as a pair candidate: judged by `sim`'s pair criterion, with the pair's release order by R103 spec — chosen and released, or the exit recorded. Take it when `engine` proposes its half: its E5 in table 2 is the second round of the corridor at the same threshold, and it is running its E7 "Lesson consolidation" (RUN.md `engine` on `main`, read 2026-11-13), so there is no pair yet | Д7, Д16; R103 spec; R17; E5 RUN.md `engine` |

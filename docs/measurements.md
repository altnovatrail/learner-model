# learner-model — measurements

> The results of the library's measurements: what was measured, how, on what, and what numbers came out. The document that describes how it is built — [architecture.md](architecture.md) — cites this one and holds only what was decided from the measurements (RULES.md, How we work; R143 spec). A measurement is cited by its section, as any document is: `measurements 1`, and from another repository `measurements 1 learner-model`.
>
> An external fact on which a measurement stands is a record of [sources.md](sources.md). Sections 1–3 were moved here word for word on 2026-11-13 from [architecture.md](architecture.md), subsections 3.6 and 3.7, by task T49 (R35), and the conclusions stayed there; "rule N" in them is a rule of subsection 3.6 of architecture. The tokens spent by the sessions of an epic, measured at its archive by the rules of epics of [RUN.md](../RUN.md), are section 4.

## 1. The same bits: the port `detmath` and `math`

### 1.1. T05's one-off program, 2026-11-02

A one-off program in the folder of T05's session, 2026-11-02 ([plan/E1.md](plan/E1.md), T05):
- the port with `float64()` gave the same bits of exp, log, σ, the power, the retention and the likelihood of the trial series on a grid — on 114,561 inputs of exp and 200,009 inputs of log — in six builds: amd64 at GOAMD64 v1 and v3, v1 with `GODEBUG=cpu.all=off`, 386 and js/wasm. On 386 and in wasm the port equals `math.Exp` and `math.Log` bit for bit;
- `math.Exp` on amd64 diverges from the port by up to 2 ulp on 18,069 inputs of 114,561, and with `cpu.all=off` on 15,562: the bits of `math.Exp` change with the processor even for one binary;
- the program itself still had `1 + r.Float64()*1000` without `float64()`, and at v3 the retention came out different — rule 1 at work;
- the differential check: the port's assembly matched on all six targets with fusion, and the control diverged on each. The build cache does not swallow the output of `-S` — a repeated build prints it again.

### 1.2. The tests of `detmath`, T14, 2026-11-07

On go1.27.3 in the container:
- the port equals `math.Exp` and `math.Log` bit for bit on 386 and in js/wasm, on 200,021 inputs of each — over the whole domain and as raw bits, NaN with any payload too. On amd64 `math.Exp` departs from the port by up to 2 ulp, and `math.Log` on normal numbers by not a single one;
- the own assembly of `math` on amd64 diverges from the port also where the port is right: `math.Log` takes a subnormal number for a number about 2⁻¹⁰²³ — `math.Log(5e-324)` gives −709.09, while ln 2⁻¹⁰⁷⁴ = −744.44, as the port gives; `math.Exp(709.782712893384)` gives +∞, and the port the last finite number; a NaN with a payload `math.Log` returns as its own NaN, and the port as the one that came in. `detmath` clamps the arguments within these edges;
- a removed `float64()` at `s*(hfsq+R)` in the port of `log` was caught by all three checks of that time: the test of rules — from T26 "Go / rules" —, `fma-diff` on six targets and the golden files in the run at GOAMD64 v3, while in the run at v1 they are green — fusion changes bits only where the compiler does it.

## 2. The time of `Readout.Task`

After T40 (R29), 2026-11-11. The benchmark is `BenchmarkReadoutTask`, a synthetic module of the form of `olympiad-math` 0.1.1 in four states: 697–755 ns per call before and 286–346 after, and `engine.Choose` 0.96–1.04 ms before and 0.57–0.69 after.

## 3. The schema and `buf`

T05's experiment: `buf convert` with `use_proto_names=true` turned evidence encoded by hand into `{"kind":1,"practice_time":"1762041600","x":-0,"items":[{"place":1,"first_verdict":"VERDICT_CORRECT"}]}` and back into the same bytes, and `buf breaking` on a removed field answered "Previously present field "3" with name "x" on message "Evidence" was deleted" ([plan/E1.md](plan/E1.md), T05).

## 4. The tokens of an epic's sessions

Measured at each epic's archive, by RUN.md, Epics: `just -f spec/justfile tokens-report --since … --until … --row` over this repository's transcripts on the machine where the epic ran, one row of the table of measurements 1.9 spec. The reference of the threshold of R145 spec is this repository's first row of at least 300 main calls after T50 (measurements 1.9 spec): the row of E7, 128.2 thousand for a main call.

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-learner-model` — E7, T51–T57: the reference | 2026-11-13T22:41 … 2026-11-16T01:48 | 7 | 355 / 86 | 128.2 / 73.2 | 0 | $26.24 |

- **E7**, measured 2026-11-16 by T58 (R40). The window runs from the merge of T50, c0db4dc at 22:40 UTC, after which E7 was unfolded, to the merge of T57, 1e2dc90 at 01:47:40 UTC; it holds the unfolding, T51, the planning of T52 and T56–T57, all on `claude-opus-5-5`. T58, which archives the epic, is a session after the window.
- **The threshold** is not reached: this row is the reference itself, and `docs-budgets.json` has not been raised — T49 set it, and T50 removed the limits held until the translation, neither of which counts as a raise (the 2026-11-13 amendment to R145 spec).

## 5. The start by the grade and the cells `D-`, T59, 2026-11-16

What the model does under grades 1, 3 and 5 of olympiad maths, where it starts a band below the truth, at it and a band above it (R55 sim), measured for the analysis of T59 (R41).

### 5.1. The harness and the trace

- **The runs.** `sim` at the tag `v0.1.21`, commit bc9aaaf (sources 6.12), from a clean clone in the session's folder; the base is `learner-model` v0.1.8 and `engine` v0.1.6, by its `go.mod`. The cells `D-grade1`, `D-grade3` and `D-grade5`, 2,000 working learners each, the rule `engine`: `sim run -populations D-grade1,D-grade3,D-grade5 -learners 2000 -set working -rule engine -part k/20` for k from 0 to 19, then `sim merge` — 6,001 lines, SHA-256 `80323302ceed586c…04b`, 2 min 19 s by the clock on 20 cores.
- **The trace** is a program added to that clone and published nowhere: it runs each learner of a cell as `runOne` of `internal/bench/bench.go` does and writes a row, tab-separated, for each of the first 200 first answers of items with a verdict, the answers `measure.Traced` counts: the population, n, k from 0, S, the model's P with retention and without (`PModel`, `PFresh`), P*, the band's place and the difficulty of the task, θ after the attempt, θ*, the mean θ + δ and the mean θ* + δ* of the task's primary skills before it, how many of them have a state, the attempt's items, whether the answer had help, and the kind of the choice; numbers with five decimals. Pairs are read as `measure.Fair` reads them — the first of ten answers in a row with |θ − θ_ref| < 0.5, a learner for whom it never holds counting as 200. On the base the trace gives the numbers of R58 sim to every digit they have: `M11.answers` 32.52 and 83.33, `M11.same` 0.812 and 0.5525. Its file, 1,200,000 rows, is `542fb926bf765a1e…32f`, 3 min 40 s by the clock.
- **The variants** are commits on top of `main` e3aaee6 in a clone in the session's folder, the trace built on each through a `go.work` with the clone of `sim` and the variant, and the three cells run by it; none goes into the repository. In `model.go`: σ₀ (`startSpread`) 2.0, 3.0 and 4.0; the weight of the trial series (`trialWeight`) 10; the floor of K_θ (`levelStepFloor`) 0.1; K₀ of θ (`levelStep0`) 0.4; K₀ of δ (`correctionStep0`) 0, a diagnostic; σ₀ 3.0 with K₀ 0.4, and the same with the weight of the series 10. In `step.go`, a diagnostic: the step without forgetting — in `answers`, P₁ in place of P and w = 1, as `nostep-base` of R31.

### 5.2. The base

| | Grade 1 | Grade 3 | Grade 5 |
|---|---|---|---|
| `M11.answers` | 32.52 | — | 83.33 |
| `M11.same` | 0.812 | — | 0.5525 |
| never within 200 answers | 0.021 | — | 0.314 |
| their part of `M11.answers` | 4.2 | — | 62.8 |
| k of those who come: mean, median of all | 28.9, 16 | — | 29.9, 30 |
| θ − θ* after the 5th, 50th, 200th answer | −0.58, −0.41, −0.30 | +0.01, −0.04, −0.07 | +0.48, +0.25, +0.10 |
| \|θ − θ*\| after them | 1.04, 0.71, 0.47 | 1.00, 0.64, 0.41 | 0.96, 0.59, 0.37 |
| \|θ − θ_ref\| after them | 0.59, 0.37, 0.24 | — | 1.00, 0.62, 0.40 |

- **The trial series' tasks.** Answers 1–3 are one skill's tasks of difficulty 1, 3 and 5, answers 4–5 the next skill's of difficulty 1 and 3. Under grades 1 and 3 the first skill is of band 1 for all 2,000 learners, and the second of band 1 for 996 and band 2 for 1,004: all five tasks under grade 1 are the reference's for every learner. Under grade 5 the first skill is of band 2 for all, and the second of bands 1, 2 and 3 for 978, 860 and 162; no learner has the reference's first task. No item of the series has guessing: c worked back from P is 0 on all of them, and 20 % of the answers after the series are a choice of five.
- **The trial series' outcome.** The learner n answers by the same fractions in the same order under every grade (R55 sim), and the tasks under grade 5 are harder for the same truth: no learner has fewer wrong answers in the series under grade 5 than under grade 3. 502 learners have under grade 5 the reference's five tasks a band up; 288 of them answered all five alike, and their θ after the 5th answer is the reference's plus 2.5 to every digit of the trace. By the wrong answers in the series, under grade 5 and under grade 3:

  | Wrong, grade and reference | Share | Never | `M11.same` | θ − θ_ref after the 5th, 50th, 200th | θ − θ* after the 5th |
  |---|---|---|---|---|---|
  | grade 5: 0, 0 | 0.015 | 0.933 | 0 | +2.50, +1.45, +0.77 | +2.62 |
  | 1, 1 | 0.110 | 0.995 | 0 | +2.50, +1.53, +0.92 | +1.40 |
  | 2, 2 | 0.058 | 0.914 | 0.034 | +2.02, +1.27, +0.82 | +0.75 |
  | 3, 3 | 0.015 | 0.931 | 0 | +1.70, +1.02, +0.69 | +0.18 |
  | 1, 0 | 0.204 | 0.005 | 0.922 | +0.14, +0.09, +0.03 | +1.15 |
  | 2, 1 | 0.296 | 0.061 | 0.708 | +0.62, +0.36, +0.22 | +0.45 |
  | 3, 2 | 0.077 | 0.058 | 0.774 | +0.61, +0.39, +0.26 | −0.39 |
  | 2, 0 | 0.096 | 0.691 | 0.220 | −1.44, −0.87, −0.51 | +0.08 |
  | 3, 1 | 0.074 | 0.223 | 0.486 | −0.85, −0.53, −0.34 | −0.61 |
  | 3, 0 | 0.010 | 1.000 | 0 | −3.10, −1.87, −1.13 | −0.92 |
  | grade 1: 0, 0 | 0.325 | 0.057 | 0.642 | −0.80, −0.45, −0.28 | +0.36 |
  | 1, 1 | 0.483 | 0.005 | 0.856 | −0.51, −0.35, −0.22 | −0.83 |
  | 2, 2 | 0.157 | 0 | 0.987 | −0.42, −0.28, −0.19 | −1.47 |

  Rows of fewer than 20 learners are left out. Under grade 5 the four rows of equal wrong answers are 0.201 of the learners, and 0.958 of them never come; the rows 2, 0 and 3, 0 are the other way round — the reference, right on all five easier tasks, stays high, and grade 5 is near the truth.
- **After the series.** The step of θ is K_θ · (S − P), with K_θ = 0.2 / (1 + 0.05 · n) down to 0.05 from n = 60, and δ of each primary skill by K_δ = 0.4 / (1 + 0.05 · m) times its share of S − P (architecture 3.1). Under grade 5, the mean S − P is −0.067 over answers 6–20 and −0.025 over 51–100, and the mean δ of the tasks' skills −0.13 against the reference's −0.06 over 51–100 and −0.21 against −0.11 over 101–200; under grade 1, +0.057 and +0.023, δ +0.03 and +0.02. A gap of 2.5 at the end of the series is 1.5 after the 50th answer and 0.9 after the 200th (the row 1, 1 above).

### 5.3. The variants

| Variant | `M11.answers`, grade 1 | `M11.same`, grade 1 | `M11.answers`, grade 5 | `M11.same`, grade 5 | Never, grade 5 | Never of equal wrong in the series, grade 5 | \|θ − θ*\| of grade 3 after the 50th, 200th | SHA-256 of the trace |
|---|---|---|---|---|---|---|---|---|
| base | 32.52 | 0.812 | 83.33 | 0.5525 | 0.314 | 0.958 | 0.64, 0.41 | `542fb926bf765a1e…32f` |
| σ₀ 2.0 | 116.82 | 0.2855 | 87.98 | 0.5095 | 0.289 | 0.970 | 0.55, 0.35 | `90fa4a2eee94468c…408` |
| σ₀ 3.0 | 17.45 | 0.9135 | 81.63 | 0.5725 | 0.325 | 0.948 | 0.68, 0.45 | `9061a8b165f2461d…209` |
| σ₀ 4.0 | 50.29 | 0.687 | 115.71 | 0.3425 | 0.413 | 0.425 | 0.67, 0.44 | `80ec0c8f764796c0…0d6` |
| series of weight 10 | 27.49 | 0.8465 | 68.20 | 0.5945 | 0.188 | 0.517 | 0.54, 0.36 | `ed5f1d35d3f048bc…bc6` |
| floor of K_θ 0.1 | 27.30 | 0.842 | 68.35 | 0.5655 | 0.144 | 0.537 | 0.60, 0.33 | `ce561f17101bf741…db9` |
| K₀ of θ 0.4 | 24.66 | 0.870 | 54.86 | 0.669 | 0.125 | 0.435 | 0.45, 0.29 | `64ea3500769003ae…84c` |
| K₀ of δ 0 | 29.70 | 0.833 | 73.39 | 0.570 | 0.180 | 0.622 | 0.60, 0.32 | `0d23945cd0ce6626…062` |
| step without forgetting | 32.69 | 0.7945 | 80.47 | 0.5595 | 0.285 | 0.928 | 0.64, 0.43 | `0e29c612554ee94e…a0e` |
| σ₀ 3.0, K₀ of θ 0.4 | 14.35 | 0.939 | 54.81 | 0.662 | 0.129 | 0.410 | 0.48, 0.30 | `af244ccdfb8d1608…e19` |
| σ₀ 3.0, K₀ of θ 0.4, series 10 | 12.87 | 0.936 | 42.54 | 0.735 | 0.073 | 0.209 | 0.43, 0.28 | `d8cb98db3044ad4c…970` |

- Every variant runs the reference under itself too, so each row's pairs are of one model. The rows have no intervals: a share near 0.55 of 2,000 learners has a standard error of about 0.011, and the judgement's yardstick is `sim`'s (R57 sim).
- The share of learners with the reference's wrong answers in the series under grade 1 stays 1.000 but for σ₀ 2.0, 0.847, and σ₀ 4.0, 0.559: a prior that much narrower or wider moves θ in the series enough for the engine to give grade 1 other tasks than grade 3. Under grade 5 the share is 0.201 but for σ₀ 2.0, 0.136, and σ₀ 4.0, 0.400.
- The runs measure the three cells `D-` alone, in `G0`: no learning and no break. What the variants do to the other 65 cells and to the criterion of the model is not measured here.

## 6. The time of "Gate" on `sim` `v0.1.23`, T60, 2026-11-16

How long the gate of `sim` `v0.1.23`, commit 61ed988, takes on the code of `main` cb8c5e0, for the 2026-11-16 amendment to R20. The earlier tags' times stand in the amendments to R20 themselves: 118 seconds on 20 cores for `v0.1.14` in T51, and 7 min 28 s of the job on a runner.

- **In the container:** `just gate exact` from a clone of `sim` at `61ed988` in the session's folder, `GATE_LEARNER_MODEL` the repository's folder, 20 cores: 2 min 17.9 s by the clock, 1,997 s user and 48 s system (`time`); 69 numbers in their bands and to the bit, no learner breaks an invariant, code 0.
- **On a runner:** the job "Gate" of CI on pull request #79, commit 01fd0ab, `ubuntu-24.04` without a Go cache: 7 min 16 s, of which the build of the stage `toolchain` 22 s and the gate 6 min 40 s, by the times of its steps (`gh api repos/altnovatrail/learner-model/actions/runs/<run>/jobs`); all 69 numbers in their bands. Below the budget `timeout-minutes: 15`, as 7 min 28 s on `v0.1.14` in T51.

# learner-model — measurements

> The results of the library's measurements: what was measured, how, on what, and what numbers came out. The document that describes how it is built — [architecture.md](architecture.md) — cites this one and holds only what was decided from the measurements. A claim cites a measurement by its number, `measurements 1`, and from another repository `measurements 1 learner-model`.
>
> An external fact on which a measurement stands is a record of [sources.md](sources.md). A section's heading names what was measured and its date, not a task. "Rule N" in measurements 1–3 is the N-th rule of code for the same bits in [architecture.md](architecture.md), which holds the conclusions. The tokens spent by the sessions of an epic, measured at its archive by the plan's rules for epics in [RUN.md](../RUN.md), are measurements 4, and the analysis of E8's, which `spec`'s threshold for a new analysis called for, is measurements 10.

## 1. The same bits: the port `detmath` and `math`

### 1.1. A one-off program, 2026-11-02

A one-off program in the folder of a session, 2026-11-02 ([plan/E1.md](plan/E1.md)):
- the port with `float64()` gave the same bits of exp, log, σ, the power, the retention and the likelihood of the trial series on a grid — on 114,561 inputs of exp and 200,009 inputs of log — in six builds: amd64 at GOAMD64 v1 and v3, v1 with `GODEBUG=cpu.all=off`, 386 and js/wasm. On 386 and in wasm the port equals `math.Exp` and `math.Log` bit for bit;
- `math.Exp` on amd64 diverges from the port by up to 2 ulp on 18,069 inputs of 114,561, and with `cpu.all=off` on 15,562: the bits of `math.Exp` change with the processor even for one binary;
- the program itself still had `1 + r.Float64()*1000` without `float64()`, and at v3 the retention came out different — rule 1 at work;
- the differential check: the port's assembly matched on all six targets with fusion, and the control diverged on each. The build cache does not swallow the output of `-S` — a repeated build prints it again.

### 1.2. The tests of `detmath`, 2026-11-07

On go1.27.3 in the container:
- the port equals `math.Exp` and `math.Log` bit for bit on 386 and in js/wasm, on 200,021 inputs of each — over the whole domain and as raw bits, NaN with any payload too. On amd64 `math.Exp` departs from the port by up to 2 ulp, and `math.Log` on normal numbers by not a single one;
- the own assembly of `math` on amd64 diverges from the port also where the port is right: `math.Log` takes a subnormal number for a number about 2⁻¹⁰²³ — `math.Log(5e-324)` gives −709.09, while ln 2⁻¹⁰⁷⁴ = −744.44, as the port gives; `math.Exp(709.782712893384)` gives +∞, and the port the last finite number; a NaN with a payload `math.Log` returns as its own NaN, and the port as the one that came in. `detmath` clamps the arguments within these edges;
- a removed `float64()` at `s*(hfsq+R)` in the port of `log` was caught by all three checks of that time: the test of rules — later the job "Go / rules" —, `fma-diff` on six targets and the golden files in the run at GOAMD64 v3, while in the run at v1 they are green — fusion changes bits only where the compiler does it.

## 2. The time of `Readout.Task`

After `Readout.Task` was made faster, 2026-11-11. The benchmark is `BenchmarkReadoutTask`, a synthetic module of the form of `olympiad-math` 0.1.1 in four states: 697–755 ns per call before and 286–346 after, and `engine.Choose` 0.96–1.04 ms before and 0.57–0.69 after.

## 3. The schema and `buf`

An experiment of 2026-11-02: `buf convert` with `use_proto_names=true` turned evidence encoded by hand into `{"kind":1,"practice_time":"1762041600","x":-0,"items":[{"place":1,"first_verdict":"VERDICT_CORRECT"}]}` and back into the same bytes, and `buf breaking` on a removed field answered "Previously present field "3" with name "x" on message "Evidence" was deleted" ([plan/E1.md](plan/E1.md)).

## 4. The tokens of an epic's sessions

Measured at each epic's archive, by the plan's rules for epics: `just -f spec/justfile tokens-report --since … --until … --row` over this repository's transcripts on the machine where the epic ran, one row of the table of measurements 1.9 spec. The reference of `spec`'s threshold for a new analysis — the context of a main call more than 25 % above it, over at least 300 main calls — is this repository's first row of at least 300 main calls after the translation of its documents (measurements 1.9 spec): the row of E7, 128.2 thousand for a main call.

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-learner-model` — E7: the reference | 2026-11-13T22:41 … 2026-11-16T01:48 | 7 | 355 / 86 | 128.2 / 73.2 | 0 | $26.24 |
| `-workspaces-learner-model` — E8 | 2026-11-16T01:48 … 2026-11-16T16:17 | 4 | 545 / 87 | 173.3 / 65.4 | 0 | $36.10 |

- **E7**, measured 2026-11-16 at its archive. The window runs from the merge of the translation of the documents, c0db4dc at 22:40 UTC, after which E7 was unfolded, to the merge of its last task, 1e2dc90 at 01:47:40 UTC; it holds the unfolding, three of its tasks and the planning of another, all on `claude-opus-5-5`. The session that archives the epic is after the window.
- **The threshold** is not reached: this row is the reference itself, and `docs-budgets.json` has not been raised — it was set when the plan took `spec`'s structure of documents, and the limits held until the translation were removed with the translation, neither of which counts as a raise.
- **E8**, measured 2026-11-16 at its archive. The window runs from the merge of E7's last task, 1e2dc90 at 01:47:40 UTC, after which E8 was unfolded, to the merge of its exploration, 7ff2b89 at 16:16:32 UTC; it holds four tasks, a session each, all on `claude-opus-5-5`. Two tasks were deferred and had no session; the session that archives the epic is after the window.
- **The threshold** is met for E8: 173.3 thousand for a main call is 35 % above the reference, 128.2, over 545 main calls, against the 25 % of `spec`'s threshold. By session, over all its calls, the subagents' with the main ones — so below the main calls' own figure —, the exploration is 187.7 thousand with 234 main calls, the analysis of the start 165.6 with 138, E7's archive with E8's plan 144.4 with 111 and `sim` v0.1.23 for "Gate" 83.2 with 62. What entered between calls carried 100.2 M of context in all: the output of Bash 47.4 M, `docs/decisions.md` read by Bash 87 times 11.2 M and RUN.md 73 times 4.3 M. `docs-budgets.json` has not been raised since it was set. The analysis found the length of the sessions the cause, not the documents: over each session's first 80 main calls E8 carries 122.1 thousand a call, below E7's 127.9 — measurements 10.

## 5. The start by the grade and the cells `D-`, 2026-11-16

What the model does under grades 1, 3 and 5 of olympiad maths, where it starts a band below the truth, at it and a band above it, measured for the analysis of why the estimate under grade 5 comes to grade 3's later than under grade 1.

### 5.1. The harness and the trace

- **The runs.** `sim` at the tag `v0.1.21`, commit bc9aaaf (sources 6.12), from a clean clone in the session's folder; the base is `learner-model` v0.1.8 and `engine` v0.1.6, by its `go.mod`. The cells `D-grade1`, `D-grade3` and `D-grade5`, 2,000 working learners each, the rule `engine`: `sim run -populations D-grade1,D-grade3,D-grade5 -learners 2000 -set working -rule engine -part k/20` for k from 0 to 19, then `sim merge` — 6,001 lines, SHA-256 `80323302ceed586c…04b`, 2 min 19 s by the clock on 20 cores.
- **The trace** is a program added to that clone and published nowhere: it runs each learner of a cell as `runOne` of `internal/bench/bench.go` does and writes a row, tab-separated, for each of the first 200 first answers of items with a verdict, the answers `measure.Traced` counts: the population, n, k from 0, S, the model's P with retention and without (`PModel`, `PFresh`), P*, the band's place and the difficulty of the task, θ after the attempt, θ*, the mean θ + δ and the mean θ* + δ* of the task's primary skills before it, how many of them have a state, the attempt's items, whether the answer had help, and the kind of the choice; numbers with five decimals. Pairs are read as `measure.Fair` reads them — the first of ten answers in a row with |θ − θ_ref| < 0.5, a learner for whom it never holds counting as 200. On the base the trace gives the numbers that `sim`'s decision log records to every digit they have: `M11.answers` 32.52 and 83.33, `M11.same` 0.812 and 0.5525. Its file, 1,200,000 rows, is `542fb926bf765a1e…32f`, 3 min 40 s by the clock.
- **The variants** are commits on top of `main` e3aaee6 in a clone in the session's folder, the trace built on each through a `go.work` with the clone of `sim` and the variant, and the three cells run by it; none goes into the repository. In `model.go`: σ₀ (`startSpread`) 2.0, 3.0 and 4.0; the weight of the trial series (`trialWeight`) 10; the floor of K_θ (`levelStepFloor`) 0.1; K₀ of θ (`levelStep0`) 0.4; K₀ of δ (`correctionStep0`) 0, a diagnostic; σ₀ 3.0 with K₀ 0.4, and the same with the weight of the series 10. In `step.go`, a diagnostic: the step without forgetting — in `answers`, P₁ in place of P and w = 1, as the rejected variant `nostep-base` of forgetting.

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
- **The trial series' outcome.** The learner n answers by the same fractions in the same order under every grade, and the tasks under grade 5 are harder for the same truth: no learner has fewer wrong answers in the series under grade 5 than under grade 3. 502 learners have under grade 5 the reference's five tasks a band up; 288 of them answered all five alike, and their θ after the 5th answer is the reference's plus 2.5 to every digit of the trace. By the wrong answers in the series, under grade 5 and under grade 3:

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
- **After the series.** The step of θ is K_θ · (S − P), with K_θ = 0.2 / (1 + 0.05 · n) down to 0.05 from n = 60, and δ of each primary skill by K_δ = 0.4 / (1 + 0.05 · m) times its share of S − P. Under grade 5, the mean S − P is −0.067 over answers 6–20 and −0.025 over 51–100, and the mean δ of the tasks' skills −0.13 against the reference's −0.06 over 51–100 and −0.21 against −0.11 over 101–200; under grade 1, +0.057 and +0.023, δ +0.03 and +0.02. A gap of 2.5 at the end of the series is 1.5 after the 50th answer and 0.9 after the 200th (the row 1, 1 above).

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

- Every variant runs the reference under itself too, so each row's pairs are of one model. The rows have no intervals: a share near 0.55 of 2,000 learners has a standard error of about 0.011, and the judgement's yardstick is `sim`'s.
- The share of learners with the reference's wrong answers in the series under grade 1 stays 1.000 but for σ₀ 2.0, 0.847, and σ₀ 4.0, 0.559: a prior that much narrower or wider moves θ in the series enough for the engine to give grade 1 other tasks than grade 3. Under grade 5 the share is 0.201 but for σ₀ 2.0, 0.136, and σ₀ 4.0, 0.400.
- The runs measure the three cells `D-` alone, in `G0`: no learning and no break. What the variants do to the other 65 cells and to the criterion of the model is not measured here.

## 6. The time of "Gate" on `sim` `v0.1.23`, 2026-11-16

How long the gate of `sim` `v0.1.23`, commit 61ed988, takes on the code of `main` cb8c5e0, for the check "Gate" of CI. On `v0.1.14` the gate took 118 seconds on 20 cores, and its job 7 min 28 s on a runner.

- **In the container:** `just gate exact` from a clone of `sim` at `61ed988` in the session's folder, `GATE_LEARNER_MODEL` the repository's folder, 20 cores: 2 min 17.9 s by the clock, 1,997 s user and 48 s system (`time`); 69 numbers in their bands and to the bit, no learner breaks an invariant, code 0.
- **On a runner:** the job "Gate" of CI on pull request #79, commit 01fd0ab, `ubuntu-24.04` without a Go cache: 7 min 16 s, of which the build of the stage `toolchain` 22 s and the gate 6 min 40 s, by the times of its steps (`gh api repos/altnovatrail/learner-model/actions/runs/<run>/jobs`); all 69 numbers in their bands. Below the budget `timeout-minutes: 15`, as 7 min 28 s on `v0.1.14`.

## 7. The exploration of the round of fairness, 2026-11-16

What the model's numbers of the start and the step do to `M11` of grades 1 and 5 and to the rest of the model's criterion, measured by the exploration of the round of fairness, which stopped before every set was read.

### 7.1. The harness and the runs

- **The harness** is `sim` `v0.1.23`, commit 61ed988, the tag of "Gate" and of the criterion `model-fair-1` (sources 6.13), from a clean clone in the session's folder; the base is `learner-model` v0.1.8 and `engine` v0.1.6, by its `go.mod`. The machine is the container's, 20 cores, shared with other sessions: the load average stood from 35 to 130 through the runs, so the times below are of a busy machine.
- **The same bits as before.** The three cells `D-` of the base, `sim run -populations D-grade1,D-grade3,D-grade5 -learners 2000 -set working -rule engine -part k/20` for k from 0 to 19 and `sim merge`, give the file `80323302ceed586c…04b`, the SHA-256 of measurements 5.1, in 6 min 34 s. The decisive run of the base, `just decisive engine <file>`, gives `9e89cb7cfd0f1f0c…763`, exactly `base` of `internal/criterion/paths.json` of `sim`, in 2 h 25 min by the clock, 49,217 s of user time.
- **The variants** are commits on top of `main` 26be23c, whose Go code is v0.1.8's, in a clone in the session's folder; none goes into the repository. Their names: `kX` — K₀ of θ (`levelStep0`) X; `fX` — the floor of K_θ (`levelStepFloor`) X; `wX` — the weight of the trial series (`trialWeight`) X; `s3` — σ₀ (`startSpread`) 3.0; `dX` — K₀ of δ (`correctionStep0`) X; `dtX` — a decay of K_θ of its own, `levelStepDecay` X, in place of the shared 0.05 for θ alone; `kbXnN` — a rule: K₀ of θ is X while the weight of answers n is below N, and the base's 0.2 from N on. So `k04d02` is K₀ of θ 0.4 and K₀ of δ 0.2, commit c8b4ce9.
- **To repeat them.** A name's numbers drop the decimal point after a leading 0: `k04` is 0.4, `k025` 0.25, `f0075` 0.075, `d015` 0.15, `dt10` 0.10, `dt30` 0.30; `w10` is 10, `w15` 15 and `w75` 7.5; N of `kbXnN` is whole. Each variant is one commit on `main` 26be23c that edits `model.go` alone: the constants of its name in the block of constants, the others as they are; `dtX` adds `levelStepDecay = X` beside `stepDecay` and puts it in place of `stepDecay` in `levelStep`, so that K_δ keeps 0.05; `kbXnN` adds `levelBoost0 = X` and `levelBoostWeight = N` and starts `levelStep` with `k0 := float64(levelStep0)`, `if n < levelBoostWeight { k0 = levelBoost0 }`, then `max(k0/(1+stepDecay·n), levelStepFloor)`. The commits of the session's clone are named above; a commit built anew has another SHA and the same numbers, and a file of results names its commit in its header, so its SHA-256 differs while every number is the same. The reader of `M11` is the three calls named in the next item, over the files of the screening.
- **The screening** runs each variant's three cells `D-` as the base's above, the bench built through a `go.work` with the clone of `sim` and the variant, `-commit` naming the variant's commit, 2–16 minutes a set. A program of the session's clone of `sim`, published nowhere, reads `M11` off two such files as `fair` and the target of `model-fair-1` read it: `bench.Fair`, `measure.ReadPair` of `measure.FairNumbers`, and the target (Δ1 + Δ5) / 2 from the same 2,000 resamples. Its numbers are `sim judge`'s to every digit on `k04w10`, `kb04n50w10` and `k04d02`.
- **The cells of risk** are the twelve populations on which `k04w10` and `kb04n50w10` failed "no worse" — `W-repeat`, `S-adult`, `F-none`, `G3`, `G3-drop`, `G6`, `G0-exact`, `W-badday`, `W-error`, `G0-miss0.25`, `G0-skills0.3`, `E-leave` —, run as the cells `D-` are, 10–30 minutes a set. Their lines and those of the cells `D-` replace the base's in a copy of the base's decisive file, by population, graph and learner, and `sim judge` reads the copy. On `kb04n50w10` the copy gave the same ten failures as its decisive run, number for number; a copy reads no failure off the other cells, which stay the base's.
- **The decisive runs** are `CANDIDATE_LEARNER_MODEL=<folder of the variant> just decisive engine <file>`, read by `sim judge` with a judgement file in the session's folder: `"judgement": "model-fair-1"`, `"criterion": "model"`, `"target": "M11.same"`, the base's file and the candidate's. The report reads them as decisive runs; it is the exploration's, not the judgement's.

### 7.2. The screening on the cells `D-`

`M11.same` of grades 1 and 5, the share of each that never comes within 200 answers, and the target with its interval; the base is 0.812 and 0.5525, never 0.021 and 0.314.

| Set | `M11.same` 1; 5 | Never 1; 5 | Target |
|---|---|---|---|
| `k04` | 0.870; 0.669 | 0.012; 0.1245 | +0.087 [0.075; 0.100] |
| `k04f01` | 0.870; 0.669 | 0.0035; 0.0665 | +0.087 [0.075; 0.100] |
| `k04w10` | 0.8955; 0.7165 | 0.011; 0.0655 | +0.124 [0.109; 0.139] |
| `f01w10` | 0.859; 0.631 | 0.0055; 0.0645 | +0.063 [0.049; 0.078] |
| `k03f0075w75` | 0.880; 0.630 | 0.0065; 0.110 | +0.073 [0.059; 0.086] |
| `k04w15` | 0.8885; 0.7735 | 0.013; 0.042 | +0.149 [0.133; 0.165] |
| `k05w10` | 0.909; 0.778 | 0.0055; 0.0345 | +0.161 [0.146; 0.177] |
| `k06w10` | 0.9145; 0.8125 | 0.004; 0.0195 | +0.181 [0.166; 0.196] |
| `k08w10` | 0.922; 0.853 | 0.0025; 0.0045 | +0.205 [0.190; 0.220] |
| `s3` | 0.9135; 0.5725 | 0.004; 0.3245 | +0.061 [0.049; 0.072] |
| `w10` | 0.8465; 0.5945 | 0.0455; 0.1875 | +0.038 [0.024; 0.054] |
| `w75` | 0.8565; 0.5765 | 0.024; 0.2415 | +0.034 [0.021; 0.048] |
| `s3w10` | 0.9025; 0.6235 | 0.025; 0.1945 | +0.081 [0.068; 0.095] |
| `s3k04` | 0.939; 0.662 | 0.0035; 0.1285 | +0.118 [0.107; 0.131] |
| `s3k04w10` | 0.936; 0.735 | 0.0075; 0.0725 | +0.153 [0.139; 0.168] |
| `s3k04f01w10d03` | 0.940; 0.7425 | 0; 0.024 | +0.159 [0.144; 0.174] |
| `s3k06w10` | 0.953; 0.8285 | 0.0015; 0.0175 | +0.209 [0.195; 0.223] |
| `k025` | 0.8445; 0.582 | 0.0145; 0.2745 | +0.031 [0.023; 0.039] |
| `k04dt10` | 0.847; 0.594 | 0.0135; 0.2605 | +0.038 [0.028; 0.048] |
| `k04dt15` | 0.800; 0.561 | 0.0215; 0.304 | −0.002 [−0.010; 0.007] |
| `k06dt30` | 0.7785; 0.5505 | 0.021; 0.3155 | −0.018 [−0.026; −0.010] |
| `kb04n30` | 0.849; 0.626 | 0.0195; 0.210 | +0.055 [0.043; 0.068] |
| `kb04n50` | 0.871; 0.673 | 0.018; 0.1715 | +0.090 [0.077; 0.102] |
| `kb04n50w10` | 0.8975; 0.7205 | 0.020; 0.095 | +0.127 [0.112; 0.142] |
| `kb06n50w10` | 0.922; 0.8225 | 0.0235; 0.057 | +0.190 [0.175; 0.205] |
| `s3kb06n50w10` | 0.956; 0.831 | 0.011; 0.055 | +0.211 [0.197; 0.226] |
| `d02` | 0.813; 0.5665 | 0.018; 0.2645 | +0.008 [−0.000; 0.016] |
| `d03` | 0.8055; 0.5635 | 0.015; 0.2885 | +0.002 [−0.005; 0.009] |
| `k03d02` | 0.8625; 0.6235 | 0.0135; 0.152 | +0.061 [0.050; 0.071] |
| `k035d02` | 0.867; 0.6555 | 0.012; 0.1045 | +0.079 [0.067; 0.090] |
| `k04d02` | 0.882; 0.688 | 0.007; 0.082 | +0.103 [0.091; 0.114] |
| `k04d03` | 0.882; 0.677 | 0.0105; 0.0955 | +0.097 [0.084; 0.110] |
| `k045d015` | 0.8845; 0.732 | 0.008; 0.037 | +0.126 [0.113; 0.139] |
| `k05d01` | 0.908; 0.7505 | 0.0045; 0.0235 | +0.147 [0.133; 0.161] |

Fourteen more sets — the floor and σ₀ 3.0 beside the others, more `dtX` and `kbXnN`, K₀ of δ 0.25, and K₀ of θ 0.45 with K₀ of δ 0.2 — fall between these rows and change none of the readings below. Of all 48, grade 5 is at most 0.853, `k08w10`; grade 1 passes 0.90 with K₀ of θ 0.5 and up at the base's decay of K_θ, or with σ₀ 3.0. `M1` of grade 3 after the 200th answer is 0.519 on the base, 0.498 with `k04`, 0.492 with `k04w10` and `k04d02`, and 0.514 with the floor 0.1 added to `k04w10`: none of these sets brings the grades together by moving grade 3 off the truth.

### 7.3. The decisive runs

| Run | Commit | SHA-256 | Time | Target | Score, a reference | "No worse" failed |
|---|---|---|---|---|---|---|
| base | — | `9e89cb7cfd0f1f0c…763` | 2 h 25 min | — | — | — |
| `k04w10` | d3471b0 | `37223ff093c28e03…95f` | 1 h 39 min | +0.124 [0.109; 0.139] | +0.044 [0.031; 0.057] | 13 |
| `kb04n50w10` | f08930e | `97d781a1d792c7b1…dbb` | 1 h 56 min | +0.127 [0.112; 0.142] | +0.028 [0.016; 0.040] | 10 |
| `kb06n50w10` | 17d544d | `0f47e0da19a7f26e…4de` | 1 h 57 min | +0.190 [0.175; 0.205] | +0.039 [0.027; 0.051] | 23 |
| `k04d02` | c8b4ce9 | `82261d96e984b070…3ca` | 1 h 53 min | +0.103 [0.091; 0.114] | +0.051 [0.039; 0.062] | 18 |

- All four are better than the base by the target, and none holds every constraint: the report's outcome is "the base stays as it is" each time. `k04d02` reaches 0.882 [0.868; 0.897] under grade 1 and 0.688 [0.669; 0.708] under grade 5.
- `k04w10`: `M1.rms.200` on `G6`, `G0-exact`, `G0-miss0.25`, `G0-skills0.3`, `F-none`, `W-error`, `W-repeat` of school maths and `E-leave` of the essays, `M3.settle` on `G3-drop` and Brier and `M7.hard` on `S-adult` of olympiad maths, `M5.caught` on `F-none` and `M8.shift` on `W-repeat` of school maths. On `W-repeat` `M1.rms.200` is worse by up to 0.333 at the tolerance 0.01.
- `kb04n50w10`, the same set with the larger K₀ only below n = 50: the same kinds, ten of them, `W-repeat` worse by up to 0.307 — the step after the 50th answer is not where the price is.
- `kb06n50w10`: 23, `M1.rms.200` on sixteen cells of school maths among them, `W-repeat` worse by up to 0.518.
- `k04d02`: `W-repeat` holds, and the failures move elsewhere — `M1.rms.200` on `G5` of both maths, `G0-skills1` of both, `G0-exact` of school maths and `S-scale3.5`; `M6.never` on `G2`, `G3`, `G0-miss1`, `G0-skills1` and `F-C1` of olympiad maths; `E3.retention` on `F-model`, `F-exp`, `F-C1` and `F-growth5` and `E3.late` on `F-growth5`; `M3.settle` on `G3-drop`; Brier on `S-adult`. Its score's group `forgets` is +0.235 [0.148; 0.318].
- Two decisive runs, `s3k04w10` and `k06w10`, were stopped unfinished when `k04w10` failed: a stronger step could only be dearer, and the cells of risk were cheaper for the rest.

### 7.4. The cells of risk

"No worse" failed on the twelve populations of 7.1 alone, by `sim judge` of the copies.

| Set | Failed | What |
|---|---|---|
| `k04` | 10 | `W-repeat` `M1.rms.200` by up to 0.344; `M1.rms.200` on five more cells of school maths and on `E-leave`; `G3-drop` `M3.settle`; `S-adult` Brier; `W-repeat` `M8.shift` |
| `w10` | 6 | `S-adult` `M7.hard` and Brier; `F-none` `M5.caught`; `M1.rms.200` on `G6`, `W-error` and `W-repeat` of school maths |
| `s3` | 17 | `M1.rms.200` and Brier on `G3-drop`, `G6`, `G0-exact`, `G0-skills0.3`, `F-none` and others, of both maths; `S-adult` `M7.hard` by up to 0.039 |
| `s3w10` | 12 | `M1.rms.200` on nine cells of school maths; `G3-drop` `M3.settle`; `S-adult` Brier; `F-none` `M5.caught` |
| `k03d02` | 10 | `G3` and `G3-drop`: `M1.rms.200`, `M3.settle`, `M3.unsettled`, `M6.never`; `W-repeat` Brier |
| `k035d02` | 10 | the same kinds as `k03d02` |
| `k04d03` | 4 | `W-repeat` `M1.rms.200` by up to 0.188; `G3-drop` `M3.settle`; `S-adult` Brier; `F-none` `M5.caught` |
| `k04d02` | 4 | `G3` `M6.never`; `G3-drop` `M3.settle`; `G0-exact` `M1.rms.200` of school maths; `S-adult` Brier |
| `k05d01` | 8 | `M6.never` on `G3`, `G0-exact` and `F-none` of olympiad maths; `G3-drop` `M3.settle`; `M1.rms.200` on `G6` and `G0-exact` of school maths; `S-adult` Brier and `M7.hard` |

`k035d025` ran and was not read: the exploration stopped first.

### 7.5. `W-repeat`

On `W-repeat` of school maths — a repeat of a task whose solution has gone to the device, answered with P* = 0.95 — the base's estimate of the skills is already high: `M1.mean` is +1.20 after the 50th answer and +2.72 after the 200th, `M1.rms.200` 2.92, while on olympiad maths it is +0.02. With `k04` it is +1.54 and +3.06, `M1.rms.200` 3.26; through the 20th answer the two are equal to the bit. A repeat is help and keeps "mastered" from being set, but its answer weighs in full, and a wider step of θ moves the estimate up by more on every repeat. K₀ of δ 0.2 with K₀ of θ 0.4 holds `W-repeat` (`k04d02`); 0.3 does not (`k04d03`).

## 8. The time of CI with the tests in two runs and the gate in four jobs, 2026-11-16

What the tests in two runs and the gate in four jobs gave: the time of `ci-test`, of the gate and of "Proto" before and after, beside measurements 3.1 spec.

- **`ci-test` in the container**, 20 cores, Go 1.27.4, the code of `main` ffd27d6, each command with an empty `GOCACHE` of its own (`time`), twice: one run with `-race -covermode=atomic -coverpkg=./... -coverprofile`, as before — 70.2 s and 70.8 s by the clock, 105.3 s and 103.7 s user; the two runs that replace it — `-race` 13.2 s and 13.8 s, then the coverage with `-covermode=set` on the same cache 5.0 s and 5.1 s, 18.2 s and 18.9 s together. With a warm cache `just ci-test` takes 11.1 s. Both profiles have 2,260 blocks, and the same blocks are marked as run, 99.1 % of the statements.
- **The gate in the container**, the same code: `sim`'s `just gate-part k 4` for k = 0…3 one after the other, from a clone of `sim` at 0a31313, `GATE_LEARNER_MODEL` this checkout — 140 s together, each part on all 20 cores; `just gate-check bands` on the four — 3 s, "The gate is green"; `just gate-check exact` — "every number is the recorded ones to the last bit". On runners the parts run at once, each on its own four processors.
- **Before, on runners:** the runs of CI on the pushes to `main` 7ff2b89 and ffd27d6, 43905689460 and 43910268369 (`gh api repos/altnovatrail/learner-model/actions/runs/<run>/jobs`, 2026-11-16): the run 7 min 24 s and 11 min 43 s from its start to its last job; "Gate" 437 s and 646 s, of them the step `gate` 394 s and 619 s and the build of `sim`'s stage 24 s and 20 s; "Go / test" 113 s and 128 s, of them `ci-test` 75 s and 84 s and the image 33 s and 37 s; "Go / arm64" 205 s and 191 s, of them `ci-test` 130 s and 127 s; "Proto" 36 s and 29 s, of them the build of the image 22 s and 24 s.
- **After, on runners:** the runs of CI on pull request #82, 43915596079 on d31dfde and 43916796535 on 47ae913, read the same way the same day. The first missed every entry of the image — the Dockerfile's hash is new, and `main` has saved none for it —, the second too, since a pull request saves none; `sim`'s image is built in each job of the gate, its Dockerfile unchanged but its entry of arm64 new. The run of 47ae913 took 3 min 30 s from its start to its last job, the run of d31dfde 6 min 44 s, of them up to 3 min 20 s in the queue for a runner: `Go / test` waited from 17:12:36 to 17:16:03. "Gate / part k" 120–126 s each, of them `gate-part` 89–95 s and the image 21–24 s; "Gate" 41 s and 48 s, of them `gate-check` 11 s and 13 s; the path of the gate 3 min 15 s on 47ae913, from the start of the first part to the end of "Gate", against 437 s and 646 s of the one job before. "Go / test" 83 s and 99 s, of them `ci-test` 28 s and 39 s; "Go / arm64" 84 s and 87 s, of them `ci-test` 30 s and 29 s; "Proto" 53 s and 29 s, of them the image 34 s and 22 s. `codecov/patch` — "Coverage not affected", and all 69 numbers of the gate in their bands.

## 9. The three rules on `sim`'s gate, 2026-11-16

What chose the guard after three failures in a row and what the rules do to the gate's numbers.

### 9.1. Three failures in a row, a computation of runs

The chance that a run of independent answers, each a failure with chance q, holds three failures in a row somewhere, by the recurrence over the length of the last run of failures (a one-off computation in the session):

| q | 20 answers | 50 | 100 | 200 |
|---|---|---|---|---|
| 0.225 — the corridor's target 0.775 | 0.153 | 0.356 | 0.592 | 0.836 |
| 0.1 | 0.016 | 0.043 | 0.085 | 0.164 |

### 9.2. The variants of the guard on the gate

`sim`'s `just gate bands`, from a clone of `sim` at 0a31313 — the gate of `v0.1.25`, the bands recorded on `learner-model` v0.1.8 —, `GATE_LEARNER_MODEL` this checkout, in the container, about 2 min a run. How often the guard fired a learner's module down was counted by a temporary counter of the session, not committed.
- **The start alone**, a band below the grade's band, no guard: out of the bands are `M1.rms.200` of `G0` in olympiad maths, 0.561 against 0.499 ± 2 × 0.016; `M7.error.10` of `G1`, 1.758 against 1.103 in olympiad maths and 1.261 against 1.022 in grade-3 maths; `M3.lag` of `G2` in olympiad maths, −1.222 against −1.024, and of `G2-half`, −0.756 against −0.520 in olympiad maths and −1.245 against −1.056 in grade-3 maths; `E1.corridor` of `G2` in olympiad maths, 0.236 against 0.255; `E2.stuck` of `G2-half` in olympiad maths and of `A-linked`, 0.016 against 0; and `M4.false.all` and `M6.late` of `G0` and `E3.late` of `F-model` in grade-3 maths are "not read" in the gate's report.
- **Three tasks in a row not solved at the end, lowering at each third:** 12,886 firings; `M1.rms.200` of `G0` in olympiad maths 1.501, `M6.never` 0.164 against 0.103, and every cell above but `G0`'s further out.
- **The same, lowering once a run:** 9,581 firings; `M3.lag` of `G2` in olympiad maths −1.460, `G3`'s `M3.unsettled` 0.964 against 0.856, and `F-model` in olympiad maths out on `M2.brier.after7`, `M5.caught`, `M5.never` and `E3.late`, 0.196 against 0.507.
- **The P of the attempts stuck on**, by the state before each: most had the model's P below 0.2 — tasks the learner was expected to fail.
- **The sum of P at least half the answers, once a run:** 2 firings, and each of the gate's 62 numbers is the start alone's to its printed digits, the same out of their bands.
- **The guard's failure as the model counts it, the same sum above the guessing floor, R̄ · σ:** the review found that by P a choice of two is never too hard. The gate's 62 numbers are the start alone's again; the firings were not counted.

### 9.3. The constraint `fair`


A rough look of `sim` at 0a31313 on 300 learners a cell, the working set: `just rough engine` of the base, `learner-model` v0.1.8, and of the code of the three rules with `CANDIDATE_LEARNER_MODEL`, the commit of a copy of the checkout whose Go code is this pull request's, the sum of P of the guard of 9.2, then `sim judge` of a judgement file of the criterion `model` with the one candidate, not committed. The report names no failure of the constraint `fair` on `D-grade1` and `D-grade5`. The candidate holds not every constraint, and the base would stay: 158 "no worse" checks fail, among them `M1.rms.200` of `G0` in olympiad maths, worse by at most 0.080 against the tolerance 0.048, of `G1` by 0.287 against 0.122, `M7.error.10` of `G1` in olympiad maths by 0.799 against 0.267, `M2.brier.all` in grade-3 maths on every cell `G`, and `M3.settle` of `G3` in grade-3 maths by 12.9 answers against 5.8; and the floors of `M4.false.all` and `M6.late` of `G0` in grade-3 maths, which the candidate's run does not read. On the gate the start alone moves these numbers, and the guard keeps them (measurements 9.2): the risk that the three rules in place of a calibration accept until there are children to calibrate on. `sim`'s own judgement of the three rules holds them to its criterion.

## 10. The context of E8's sessions, 2026-11-17

The analysis that `spec`'s threshold called for, which measurements 4 found met for E8: why a main call of E8 carries 173.3 thousand of context against the reference's 128.2, and what of it to remove.

- **How.** `just -f spec/justfile tokens-report --top 25` over this repository's transcripts on the machine where both epics ran, for the window of each row of measurements 4 — E7 from 2026-11-13T22:40:00Z to 2026-11-16T01:47:40Z, E8 from 2026-11-16T01:47:40Z to 16:16:32Z —, its `--json` for the main calls of each session — the columns of their context and its largest below —, and a script of Node over the same transcripts for the context of each main call by its number in the session: the input, the cache read and the cache written of its `usage`, once a message — the last column below; it counts the exploration's session as 235 main calls, one more than `tokens-report`. The script is in no repository; it prints numbers alone, never what a transcript says. E7's window counts 358 main calls here at 128.6 thousand against the 355 at 128.2 of measurements 4, since it starts a minute earlier.
- **By session, main calls alone**, thousand a call:

| Epic | Session | Main calls | Context of a call | Largest | Over the first 80 calls |
|---|---|---|---|---|---|
| E7 | the plan of E7 | 73 | 128.8 | 190.7 | 128.8 |
| E7 | `sim` v0.1.14 for "Gate" | 66 | 104.8 | 144.4 | 104.8 |
| E7 | the planning of a candidate of forgetting | 70 | 163.8 | 246.9 | 163.8 |
| E7 | `spec` v0.18.3 | 79 | 115.2 | 175.2 | 115.2 |
| E7 | the plan in `run/` | 62 | 127.1 | 182.9 | 127.1 |
| E7 | all five | 350 | 127.9 | 246.9 | 127.9 |
| E8 | E7's archive and E8's plan | 111 | 161.8 | 237.7 | 136.9 |
| E8 | the analysis of the start | 138 | 179.8 | 267.5 | 136.4 |
| E8 | `sim` v0.1.23 for "Gate" | 62 | 93.1 | 130.1 | 93.1 |
| E8 | the exploration | 234 | 196.3 | 327.7 | 115.5 |
| E8 | all four | 545 | 173.3 | 327.7 | 122.1 |

  E7's other two sessions in the window, the end of the translation and a check of `spec`'s, are 8 calls. Every session starts at 26–40 thousand and grows over its first 60 calls by 1.8–3.0 thousand a call in E7 and 1.7–2.7 in E8.
- **The cause is the length of the sessions.** Over each session's first 80 main calls E8 carries 122.1 thousand a call, 4.5 % below the 127.9 of E7's five sessions, and over its first 60, 107.2 against 118.8. E7's longest session has 79 main calls; E8's three of 111, 138 and 234 add 243 calls beyond the 80th, at 200–330 thousand each, and they alone lift the mean from 122.1 to 173.3. The context grows with what a session has gathered (measurements 1.5 spec), so a session twice as long carries more per call at the same pace.
- **What entered between calls** is E7's in size: by Bash 1,276,147 characters in 591 results against 1,261,456 in 398; `docs/decisions.md` by Bash 294,693 characters in 87 reads against 293,746 in 78, and E7 read 103,212 more by Read; RUN.md 173,370 against 209,746, by Bash and Read together. What they carried grew — Bash 47.36 M against 20.99 M, `docs/decisions.md` 11.17 M against 4.81 M, RUN.md 6.19 M against 5.08 M, the model's own output 21.94 M against 7.99 M — because the same text rides on more calls after it in a longer session. `docs/decisions.md` grew from 400,235 to 431,799 bytes over E8 and RUN.md was 28–31 KB; within the budgets of `docs-budgets.json`, unraised since they were set.
- **Why the sessions were long.** E7's archive wrote the epic, measured E7 and planned the analysis of the start; the analysis and the exploration ran `sim`'s cells — 39 and 60 calls of runs —, the exploration its 48 sets screened and decided. The exploration also waited: 19 commands in the background, then 76 main calls that only waited for them — 70 loops `until grep -q …; do sleep …; done` in the foreground, each holding its call until the run ended or its timeout of 10 to 60 minutes did, and then called again, the rest Monitor and the stopping of tasks —, at 203 thousand a call on average, 15.4 M of the session's 45.9 M, $3.08 of its $14.62 at the cache's price, $0.20 a million (`spec/tools/tokens/prices.mjs`). By the script, the exploration's other 159 calls carry 191.7 thousand against 195.4 over all 235: the waits cost, but lift the mean little. Over 10 hours 41 minutes, from 05:34 to 16:15 UTC, the exploration had no compaction.
- **Smaller things.** 7 Reads of a tool's output that Claude Code kept in a file — RUN.md cut at 40,000 characters, a part of `sim`'s decision log and the outputs of background commands — brought 224,516 characters, 4.37 M carried; 86 notifications of background commands, 4.51 M; the attachment of the cost ceiling on 625 calls, 2.25 M. None is the documents'.

## 11. The time of recomputation by a rough step, 2026-11-02

What showed that recomputation needs no snapshot: whether a fold of the whole log is fast enough to run at each change of the cache's mark.

- **How.** A one-off Node script in the session's folder, not kept, with a rough model of the step — not the library's code, which did not exist yet —, folded synthetic logs from the empty state ([docs/plan/E1.md](plan/E1.md)).
- **Numbers.** 10,000 pieces of evidence in 13–18 ms; 50,000 pages of 10 items in about 55 ms. The time grows linearly with the log.

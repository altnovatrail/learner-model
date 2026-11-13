# learner-model — measurements

> The results of the library's measurements: what was measured, how, on what, and what numbers came out. The document that describes how it is built — [architecture.md](architecture.md) — cites this one and holds only what was decided from the measurements (RULES.md, How we work; R143 spec). A measurement is cited by its section, as any document is: `measurements 1`, and from another repository `measurements 1 learner-model`.
>
> An external fact on which a measurement stands is a record of [sources.md](sources.md). Sections 1–3 were moved here word for word on 2026-11-13 from [architecture.md](architecture.md), subsections 3.6 and 3.7, by task T49 (R35), and the conclusions stayed there; "rule N" in them is a rule of subsection 3.6 of architecture. The measurement of the tokens spent by the sessions of an epic, by the rules of epics of [RUN.md](../RUN.md), will take its place here as a section of its own when it is made for the first time.

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

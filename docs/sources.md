# learner-model — sources

> The facts of the outside world the repository's documents under `docs/` rest on and no run of ours checks: research, a law, a standard's normative text, a provider's terms and policies, a price, a license, a page that changes in place with no version, a text outside the repositories. Here are their records — the address, the date it was read and a quote in the source's language — and [architecture.md](architecture.md) and the decision log hold the claim and a citation of its one record by its number, `sources 1.1`, not a chain. A record's heading names its topic, not a task.
>
> No record is needed for a library's documentation or code at a version or commit — a link beside the claim —, for a behaviour a test or a measurement checks — they are its evidence —, or for spec's tags and the files of NovaTrail's and MathTrail's repositories — git, at a commit.
>
> A record is a subsection "N.M" of a topic: one source, or one reading of several sources for one decision. A record's number never changes and is never given again: a new one goes at the end of its topic, a new topic at the end of the file, and a record removed leaves a gap. A source read again with the same text changes the record's date; with another text, the claims and decisions on it are rewritten in the same change and the record holds the new reading alone — where a decision waits on the author, an open question, and both readings stand until it is closed. A citation inside a quote is elided as "[…]", and "section N" inside a record is a section of the text the record reads. The model is `spec/docs/sources.md`. The results of measurements are not here but in [measurements.md](measurements.md).

## 1. MathTrail: standalone

No records: a file of `mathtrail/mathtrail-standalone` is git, a link to the file at its commit beside the claim.

## 2. Learner models and the scoring of answers

### 2.1. The Rasch model

- G. Rasch, [On General Laws and the Meaning of Measurement in Psychology](https://digitalassets.lib.berkeley.edu/math/ucb/text/math_s4_v4_article-22.pdf), Proceedings of the Fourth Berkeley Symposium on Mathematical Statistics and Probability, vol. 4, 1961, pp. 321–333, read 2026-11-02 from the text layer of the scan:
  - section 3 — "For the responses to the items of an intelligence test the following model is suggested… The probability that person number v gives a correct answer to item number i is" — then formula (3.1). The text layer renders the formula with errors, so it is not repeated here: that the Rasch model and Elo have one form is said by Pelánek 2016, below;
  - Rasch's 1960 book has no open copy, and it was not read. Standalone cites the 1961 article for the same reason ([`research/literature/learner-models.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/research/literature/learner-models.md), commit c6cbeb2).

### 2.2. The guessing floor

- A. Birnbaum, [Some Latent Trait Models and Their Use in Inferring an Examinee's Ability](https://faculty.ucmerced.edu/jvevea/classes/290_21/readings/week%209/Birnbaum.pdf), in F. M. Lord, M. R. Novick, Statistical Theories of Mental Test Scores, 1968, ch. 17, read 2026-11-02 from the text layer of the scan:
  - section 17.3 — "it further assumes that if he does not know it he will guess, and, with probability c_g, will guess correctly".

### 2.3. Elo in learning

- R. Pelánek, [Applications of the Elo Rating System in Adaptive Educational Systems](https://www.fi.muni.cz/~xpelanek/publications/CAE-elo.pdf), Computers & Education 98 (2016), 169–179; the authors' version, read 2026-11-02:
  - section 2.2:
    - the update — "θs := θs + K(correct_si − P(correct_si = 1))";
    - "In this formulation the Elo rating system is closely related to the Rasch model (one parameter model) used in IRT. These two approaches share the same functional form of the model for predicting the probability of a correct answer, they differ only in the parameter estimation procedure";
    - on guessing — "for question with k options the probability of a correct answer becomes" 1/k + (1 − 1/k) / (1 + e^(−(θs − di)));
  - section 2.3 — "An intuitive improvement used in most Elo extensions is to use an "uncertainty function" instead of a constant K"; on Math Garden — "Klinkenberg et al. (2011) use uncertainty initialized as U = 1 and updated after each answer as follows:" U := U − 1/40 + D/30, "where D is a time from a previous attempt (in days)". The article Klinkenberg et al. 2011 was not read;
  - section 4.3 — "It is also useful to incorporate forgetting, particularly when we are modeling declarative knowledge (in contrast to procedural skills, which are the typical focus of student modeling)";
  - Elo's 1978 book has no open copy, and it was not read. Standalone does not cite it either.

### 2.4. The global level and the correction

- J. Nižnan, R. Pelánek, J. Řihák, [Student Models for Prior Knowledge Estimation](https://files.eric.ed.gov/fulltext/ED560514.pdf), Proceedings of the 8th International Conference on Educational Data Mining, 2015, pp. 109–116, read 2026-11-02:
  - section 2.1 — "We use an uncertainty function U(n) = α/(1 + βn), where n is the number of previous updates to the estimated parameters";
  - section 2.3 — "In addition to the global skill θs the model now uses also the concept skill θsc … we just correct the global skill by the concept skill".

### 2.5. The first attempt

- [KDD Cup 2010: Educational Data Mining Challenge, Data Format](https://pslcdatashop.web.cmu.edu/KDDCup/rules_data_format.jsp), read 2026-11-02: "Correct First Attempt: the tutor's evaluation of the student's first attempt on the step—1 if correct, 0 if an error"; an error of the first attempt there also counts a request for a hint — "the first attempt for the step was an error (incorrect attempt or hint request)".

### 2.6. Testlets

- E. T. Bradlow, H. Wainer, X. Wang, A Bayesian Random Effects Model for Testlets, Psychometrika 64 (1999), 153–168, [doi:10.1007/BF02294533](https://doi.org/10.1007/BF02294533). The article is closed access; the abstract was read, through OpenAlex, 2026-11-02: "all items given to an examinee are unlikely to be conditionally independent (given examinee proficiency). Models that assume conditional independence will overestimate the precision with which examinee proficiency is measured. Overstatement of precision may lead to inaccurate inferences such as prematurely ending an examination in which the stopping rule is based on the estimated standard error of examinee proficiency".

### 2.7. The design effect and the effective sample size

- G. Kalton, J. M. Brick, T. Lê, [Estimating components of design effects for use in sample design](https://unstats.un.org/unsd/hhsurveys/pdf/chapter_6.pdf), chapter VI in Household Sample Surveys in Developing and Transition Countries, UN, ST/ESA/STAT/SER.F/96, March 2005, read 2026-11-02 from the text layer:
  - "The effective sample size is the size of an unrestricted sample that would yield the same level of precision for the survey estimate as that attained by the complex design";
  - the cluster effect — formula (18), 1 + (b − 1)ρ, — "depends on two factors: the subsample size within selected PSUs (b) and the intra-class correlation (ρ)". The text layer renders the formulas and the letter ρ with errors, so the formula is written here by its meaning;
  - the chapter traces the term "design effect" to Kish 1965 — the book was not read, and nothing is claimed from it.

### 2.8. Elo with several skills

- Pelánek 2016, above, section 2.4, read 2026-11-02: the multidimensional extension takes the case where "each item belongs to a single knowledge component", and after an answer updates every skill j with the correlation c_ij: θ_sj := θ_sj + c_ij · K · (correct − P). There is no division of the step among the primary skills of one item there, and the 1/J of this repository's step is derived from the derivative, not taken from there.

## 3. Forgetting and repetition

### 3.2. HLR and the Leitner system

- B. Settles, B. Meeder, [A Trainable Spaced Repetition Model for Language Learning](https://aclanthology.org/P16-1174.pdf), Proceedings of the 54th Annual Meeting of the Association for Computational Linguistics, 2016, pp. 1848–1858, read 2026-11-02:
  - section 3.2 — Leitner boxes with intervals of 1, 2, 4 days and on: "All cards start out in the 1-day box", and on an error a card moves to a box with a shorter interval;
  - section 3.3 — "This posits that memory decays exponentially over time": p = 2^(−Δ/h), where h is "the half-life or measure of strength in the learner's long-term memory"; and "Under the assumption that half-life should increase exponentially with each repeated exposure (a common practice in spacing and lag effect research)" the half-life is a function of features whose weights are learned from Duolingo's data.

### 3.3. SM-2

- P. Wozniak, [SuperMemo 2: Algorithm](https://super-memory.com/english/ol/sm2.htm), from the 1990 master's thesis, on the site since 1998-05-10, read 2026-11-02: "With all items associate an E-Factor equal to 2.5"; "I(1):=1", "I(2):=6", "for n>2: I(n):=I(n-1)*EF", the interval — "(in days)"; "If the quality response was lower than 3 then start repetitions for the item from the beginning without changing the E-Factor".

### 3.4. The forgetting curve

- H. Ebbinghaus, [Memory: A Contribution to Experimental Psychology, chapter VII](https://psychclassics.yorku.ca/Ebbinghaus/memory7.htm), 1885, translation 1913, read 2026-11-02, § 29, Discussion of Results: "forgetting would be very rapid at the beginning of the process and very slow at the end" and "the quotients of the amounts retained and the amounts forgotten were inversely as the logarithms of the times".

### 3.5. The form of forgetting

- J. T. Wixted, E. B. Ebbesen, On the Form of Forgetting, Psychological Science 2 (1991), 409–415, [doi:10.1111/j.1467-9280.1991.tb00175.x](https://doi.org/10.1111/j.1467-9280.1991.tb00175.x). The article itself was not read: the open copy is a scan without a text layer. The abstract was read, through OpenAlex, 2026-11-02: "In all cases, the course of forgetting was best described by a simple power function of time relative to five other reasonable alternatives (linear, exponential, exponential-power, hyperbolic, and logarithmic)".

### 3.6. The forgetting of skills

- W. Arthur Jr., W. Bennett Jr., P. L. Stanush, T. L. McNelly, Factors That Influence Skill Decay and Retention: A Quantitative Review and Analysis, Human Performance 11 (1998), 57–101, [doi:10.1207/s15327043hup1101_3](https://doi.org/10.1207/s15327043hup1101_3). The article is closed access; the abstract was read, through OpenAlex, 2026-11-02: "Results indicate that there is substantial skill loss with nonpractice or nonuse, with the amount of skill loss ranging from an effect size (d) of -0.01 immediately after training to a d of -1.4 after more than 365 days of nonuse" and "Physical, natural, and speed-based tasks were less susceptible to skill loss than cognitive, artificial, and accuracy-based tasks".

### 3.7. Intervals in mathematics

- D. Rohrer, K. Taylor, [The Effects of Overlearning and Distributed Practice on the Retention of Mathematics Knowledge](https://files.eric.ed.gov/fulltext/ED505642.pdf), Applied Cognitive Psychology 20 (2006), 1209–1224; the authors' manuscript in ERIC ED505642, read 2026-11-02 from the text layer. The subjects are university students, "116 undergraduates" in the first experiment. General Discussion — "In Experiment 1, distributing 10 practice problems across two sessions instead of massing all 10 problems in the same session had no effect on one-week test scores but virtually doubled four-week test scores. In Experiment 2, increasing the number of problems solved in a single session from three to nine had virtually no effect on test scores on either the one- or four-week test".

### 3.8. The ACT-R activation model

- P. I. Pavlik Jr., J. R. Anderson, Practice and forgetting effects on vocabulary memory: an activation-based model of the spacing effect, Cognitive Science 29 (2005), 559–586, [PubMed 21702785](https://pubmed.ncbi.nlm.nih.gov/21702785/). The article itself was not read: the authors' copy did not open. The abstract was read, 2026-11-02: "each time an item is practiced it receives an increment of strength but that these increments decay as a power function of time" and "This mechanism limits long-term benefits from further practice at higher levels of activation and produces the spacing effect".

### 3.9. Pelánek's correction for time

- R. Pelánek, [Modeling Students' Memory for Application in Adaptive Educational Systems](https://www.fi.muni.cz/~xpelanek/publications/memory-geography.pdf), Proceedings of the 8th International Conference on Educational Data Mining, 2015, pp. 480–483; the authors' version, read 2026-11-02 from the text layer:
  - section 2.2 — "instead of P(m) to use P(m + f(t)), where t is the time (in seconds) from the last attempt", and "the value f(t) corresponds to temporal increase in memory activation due to (short) time from previous exposure of an item". The data gave no simple function f, and it is a step function, fitted to the data, and past the last step — "0 otherwise";
  - the data — places on a map, slepemapy.cz, section 3.

### 3.10. Glicko

- M. E. Glickman, [The Glicko system](http://www.glicko.net/glicko/glicko.pdf), undated, read 2026-11-02: "In the Glicko system, a player's rating changes only from game outcomes, but his/her RD changes both from game outcomes and also from the passage of time when not playing"; after t periods without games — RD = min(√(RD²_old + c² · t); 350).

## 4. Go

### 4.1. Floating-point operations in Go

- [The Go Programming Language Specification](https://go.dev/ref/spec#Floating_point_operators), the section Floating-point operators, Language version go1.27 (May 26, 2026), read 2026-10-19: "An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements, and produce a result that differs from the value obtained by executing and rounding the instructions individually".

### 4.2. Explicit conversion against fusion in Go

- [The Go Programming Language Specification](https://go.dev/ref/spec#Floating_point_operators), the section Floating-point operators, Language version go1.27 (May 26, 2026), read 2026-11-02: "An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding" and the example "r = float64(x*y) + z" among those where FMA is forbidden.

### 4.3. The sources of Go 1.27.3

- the directory `src` of Go 1.27.3 in the container (`go env GOROOT`), read 2026-11-02:
  - `math/exp_asm.go` — `//go:build amd64 || arm64 || loong64 || riscv64 || s390x` and `const haveArchExp = true`: `math.Exp` has its own assembly implementation on these architectures, and on the rest, WASM included, one in Go (`math/exp_noasm.go`);
  - `math/exp_amd64.go`, line 11, — `var useFMA = cpu.X86.HasAVX && cpu.X86.HasFMA`: the path in `math/exp_amd64.s` is chosen by the processor at run time. It can be turned off through `internal/cpu/cpu_x86.go` — `option{Name: "fma", Feature: &X86.HasFMA}`, that is, `GODEBUG=cpu.fma=off`;
  - `math/log_asm.go` — `//go:build amd64 || s390x`;
  - `math/pow.go`, line 128, — `a1 = Exp(yf * Log(x))`;
  - the compiler's rules: in `cmd/compile/internal/ssa/_gen/ARM64.rules` — `(FADDD a (FMULD  x y)) && a.Block.Func.useFMA(v) => (FMADDD  a x y)`, and in `AMD64.rules` — `(ADDS(S|D) (MULS(S|D) x y) z) && buildcfg.GOAMD64 >= 3 && z.Block.Func.useFMA(v) => (VFMADD231S(S|D) z x y)`.

### 4.4. Fusion, conversions and order in Go

- [The Go Programming Language Specification](https://go.dev/ref/spec), Language version go1.27 (May 26, 2026), the copy `doc/go_spec.html` of Go 1.27.3 in the container, read 2026-11-02:
  - Floating-point operators — "An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements" and "An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding"; among the examples where fusion is allowed — `t = x*y; r = t + z`;
  - Conversions — "In all non-constant conversions involving floating-point or complex values, if the result type cannot represent the value the conversion succeeds but the result value is implementation-dependent";
  - For statements with range clause — "The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next".

### 4.5. The compiler and standard library of Go 1.27.3

- the directory `src` in the container (`/usr/local/go`), read 2026-11-02:
  - the fusion rules, each with the condition `useFMA`: `cmd/compile/internal/ssa/_gen/ARM64.rules`, lines 1798–1809 — `FMADDD`, `FMSUBD`, `FNMSUBD`, `FNMADDD`; `PPC64.rules`, 1013–1014; `S390X.rules`, 1266–1267; `RISCV64.rules`, 808–810; `LOONG64.rules`, 770–781; `AMD64.rules`, 1501 — only `ADDS(S|D)` and only `buildcfg.GOAMD64 >= 3`. `Wasm.rules`, `386.rules`, `MIPS.rules` and `MIPS64.rules` have no such rules;
  - `generic.rules`, line 1331, — `(Div64F x (Const64F <t> [c])) && reciprocalExact64(c) => (Mul64F x (Const64F <t> [1/c]))`: a division by a constant with an exact reciprocal becomes a multiplication;
  - `cmd/compile/internal/ssagen/ssa.go`, line 2709, — a conversion of `float64` to `float64` gives `OpRound64F`, and the fusion rules do not pass through it;
  - `cmd/compile/internal/ssa/func.go`, `useFMA`, — "useFMA allows targeted debugging w/ GOFMAHASH", and `base/hashdebug.go`, lines 64–70, — the value "n" selects nothing: `-d=fmahash=n` turns off all the fusion rules;
  - `cmd/compile/internal/base/flag.go`, lines 268–273, — "quietly disable the convert hash changes", and `AMD64.rules`, lines 168–180, — a different conversion of numbers to integers on amd64 is behind this switch;
  - `math/sqrt.go` — "Return correctly rounded sqrt"; `Wasm.rules`, line 349, — `(Sqrt ...) => (F64Sqrt ...)`;
  - `math/exp.go` and `math/log.go` — pure Go `exp`, `expmulti` and `log`, with the Sun Microsystems notice from FreeBSD msun: "Permission to use, copy, modify, and distribute this software is freely granted, provided that this notice is preserved"; Go's license — `LICENSE`, BSD-3-Clause;
  - `internal/cpu/cpu_x86.go`, lines 95–102, — `fma` is among the options `GODEBUG` can turn off only at GOAMD64 below v3, and `runtime/extern.go` — "cpu.all=off disables the use of all optional instruction set extensions";
  - `runtime/asm_amd64.s`, line 69, — a v3 build does not start on a processor without v3: "This program can only be run on AMD64 processors with v3 microarchitecture support";
  - `encoding/json` is built from the v2 files — in `internal/buildcfg/exp.go` the base experiments have `JSONv2: true`;
  - `slices/sort.go` — "This sort is not guaranteed to be stable"; `unicode/tables.go` — `const Version = "17.0.0"`;
  - `lib/wasm/go_js_wasm_exec` — the js/wasm tests run in Node, `exec node --stack-size=8192 "$DIR/wasm_exec_node.js"`; `go tool dist list -json` — `js/wasm` and `wasip1/wasm` have `"FirstClass": false`.

### 4.6. Build constraints by file name in Go 1.27.3

- the directory `src` in the container, read 2026-11-07: `go/build/build.go`, `goodOSArchFile`, — the name of a file up to the first dot, without the first word and without `_test`, ends in `_GOOS`, `_GOARCH` or `_GOOS_GOARCH`; the lists of systems and architectures are `KnownOS` and `KnownArch` in `internal/syslist/syslist.go`. A build with `-gcflags='-S -d=fmahash=n'` prints, at every place where a fusion rule would fire, the line "fmahash triggered [DISABLED]" with a hash, but without the file's name and line — where exactly is shown by the difference of the assembler.

### 4.7. Go and WebAssembly

- [Go Wiki: WebAssembly](https://go.dev/wiki/WebAssembly), read 2026-11-02: "Set `GOOS=js` and `GOARCH=wasm` environment variables to compile for WebAssembly"; "Go 1.21 introduced WASI as a supported platform. To build for WASI, use the `wasip1` port"; "At present, Go generates large Wasm files, with the smallest possible size being around ~2MB" and "Using `gz` compression reduces the ~2MB (minimum file size) example WASM file down to around 500kB"; "TinyGo supports a subset of the Go language".

### 4.8. protobuf for Go

- [`google.golang.org/protobuf/proto`, `MarshalOptions`](https://pkg.go.dev/google.golang.org/protobuf/proto#MarshalOptions), v1.36.12, read 2026-11-02: `Deterministic` — "the same message will always be serialized to the same bytes within the same binary", and "the deterministic serialization is NOT canonical across languages. It is not guaranteed to remain stable over time".

## 5. CI, GitHub and tools

### 5.1. GitHub's runners

- [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), read 2026-11-02: the labels `ubuntu-24.04-arm`, `ubuntu-22.04-arm` and `ubuntu-26.04-arm` are in the table of standard runners for public repositories, and "Use of the standard GitHub-hosted runners is free and unlimited on public repositories". The page does not name the processors of the x64 runners. `altnovatrail/learner-model` is public (`gh repo view`, 2026-11-02).

### 5.2. The shared Go wrapper

- [`.github/workflows/go.yml`](https://github.com/altnovatrail/.github/blob/af40f949676fb3a01abb0a1e6ffbc99ce824660b/.github/workflows/go.yml), `altnovatrail/.github`, `main` af40f94, read 2026-11-02 through `gh api`: "by four recipes of that repository's own — ci-lint, ci-test, ci-vuln and ci-licenses — one job each, in the toolchain stage of the caller's .devcontainer/Dockerfile, where the caller pins its Go"; the jobs run on `ubuntu-24.04`, without arm64 and without WASM.

### 5.3. The shared workflows of `.github` v0.1.15

- the release [v0.1.15](https://github.com/altnovatrail/.github/releases/tag/v0.1.15), commit 89d3cee, read 2026-11-06; `gh release verify v0.1.15 -R altnovatrail/.github` — "Release v0.1.15 verified!":
  - `go.yml` — "ci-test writes the coverage of the tests as one Go coverage profile, coverage/go.out, and the job coverage sends it to Codecov by OIDC, with no secret"; `ci-licenses` gets `ALLOWED_LICENSES` — "the ASF's Category A as SPDX identifiers";
  - `codeql.yml` — "Besides actions, javascript-typescript and go are taken", and Go is read from the Go build of the `toolchain` stage; at v0.1.0 the input `languages` took only `javascript-typescript`;
  - `secrets.yml`, `dco.yml`, `documents.yml` and `scorecard.yml` did not change between v0.1.0 and v0.1.15 (`git diff --stat v0.1.0 v0.1.15`);
  - the texts of the calling files, `codecov.yml` and the order of the move to coverage — [`docs/architecture.md`](https://github.com/altnovatrail/.github/blob/main/docs/architecture.md) of `.github`, section 3.3, `main`, read 2026-11-06.

### 5.4. The shared workflows of `.github` v0.1.21

- the release [v0.1.21](https://github.com/altnovatrail/.github/releases/tag/v0.1.21), commit 706f406, read 2026-11-08 in a clone in a temporary folder; `gh release verify v0.1.21 -R altnovatrail/.github` — "Release v0.1.21 verified!". With it, the earlier "on `ubuntu-24.04`, without arm64 and without WASM" above holds only without the input `same-bits`:
  - `go.yml` — "With same-bits, four more jobs hold the caller's bits the same on every build": ""Go / targets" runs the tests on the builds where the bits are most at risk, and builds and vets wasip1; "Go / fma-diff" compares the assembler of each package with fusion and without it, with a control of its own here; "Go / arm64" runs ci-test and the targets of its host on an arm64 runner; "Go / rules" holds the code to the rules of code by the tool tools/coderules of .github, of this workflow's own version, and the caller's lists, .github/coderules.json"; "Without the input the workflow is what it was";
  - `tools/coderules/lists.go` — the fields of the lists `code-imports`, `test-imports`, `test-math`, `requires` and `nested-modules`, "A field it does not know is a mistake";
  - between v0.1.15 and v0.1.21 `codeql.yml` and `scorecard.yml` only raised the versions of actions, `documents.yml` calls `spec/tools/toolchain/toolchain.sh`, and `release.yml` changed only the npm branch (`git diff v0.1.15 v0.1.21`).

### 5.5. The shared workflows of `.github` v0.1.24

- the release [v0.1.24](https://github.com/altnovatrail/.github/releases/tag/v0.1.24), commit fab102d, read 2026-11-13 in a clone in a temporary folder; `gh release verify v0.1.24 -R altnovatrail/.github` — "Release v0.1.24 verified!"; the latest release (`gh release list -R altnovatrail/.github`, 2026-11-13). Between v0.1.21 and v0.1.24 (`git diff v0.1.21 v0.1.24 -- .github/workflows/`):
  - `go.yml` — ""Go / targets" runs the tests on the builds where the bits are most at risk, wasip1 among them": `GOOS=wasip1 GOARCH=wasm go test -count=1` with `-exec` — `node --stack-size=8192` and a runner on `node:wasi`, `new WASI({ version: "preview1", args, env: process.env, preopens: { "/": "/" }, returnOnExit: true })`; the same in "Go / arm64". Also, each job in the `toolchain` stage takes its image from the GitHub Actions cache by the key of the OS, the architecture and the hash of the Dockerfile, and builds it only on a miss, "so that an outage of what the Dockerfile downloads from does not fail the checks while the Dockerfile stays the same". The inputs and the permissions of the call are the same;
  - `documents.yml` — the same cache of the stage for `spec/tools/toolchain/toolchain.sh` through `TOOLCHAIN_IMAGE_FILE`, which the script reads since `spec` v0.13.0;
  - `release.yml` and `scorecard.yml` — only the links in comments; `codeql.yml`, `secrets.yml` and `dco.yml` did not change.

### 5.6. The tools of the checks and gopls

- the Go module proxy, `https://proxy.golang.org/<module>/@latest`, read 2026-11-06: `golang.org/x/vuln` — v1.9.0 of 2026-10-11, `github.com/google/go-licenses/v2` — v2.0.1 of 2025-09-08, `golang.org/x/tools/gopls` — v0.24.0 of 2026-10-14. The licenses — by GitHub's API, `repos/<repository>/license`, 2026-11-06: govulncheck and gopls — BSD-3-Clause (`golang/vuln` and `golang/tools`), go-licenses — Apache-2.0 (`google/go-licenses`). Since `spec` v0.9.16 the lock file of gopls is the one of `spec`'s template: `.devcontainer/` from the output of the recipe `new-repo` for `learner-model` of the kind `go` differs from the former one only in comments, and `go.sum` is the same byte for byte (`diff -r`, 2026-11-09).

### 5.7. Codecov

- read 2026-11-06: the Codecov app is installed for the organization on all repositories — `gh api orgs/altnovatrail/installations`, `"app_slug":"codecov","repository_selection":"all"`; `learner-model` is already in Codecov but not activated — `https://api.codecov.io/api/v2/github/altnovatrail/repos/learner-model/`, `"active":false,"activated":false`.

### 5.8. The shared workflows of `.github` v0.1.31 and `go test -covermode`

- the release [v0.1.31](https://github.com/altnovatrail/.github/releases/tag/v0.1.31), commit 93c848f, read 2026-11-16 in a clone in a temporary folder: immutable, of 2026-11-16 16:11 UTC, the latest (`gh release list -R altnovatrail/.github`, `gh release view v0.1.31 --json isImmutable`); `gh release verify v0.1.31 -R altnovatrail/.github` — "✓ Release v0.1.31 verified!". Between v0.1.24 and v0.1.31 (`git diff v0.1.24 v0.1.31 -- .github/workflows/`) the `on.workflow_call` of `go.yml`, `documents.yml`, `release.yml`, `codeql.yml`, `secrets.yml`, `dco.yml` and `scorecard.yml` did not change, and:
  - `go.yml` — "The caches of Go — the modules, the builds and golangci-lint's — come from the same cache, one entry for each check, which a push to main alone saves"; "On a miss the job toolchain alone builds it and hands it on to the others by an artifact of the run, and a push to main alone saves it"; "lint, vuln and licenses, which start no browser and run no test, run in the stage toolchain-core, Node, just, Go and buf, where the Dockerfile has one that holds less than toolchain, by an image of its own"; "Go / targets" — "which run at once, a job "Go / target (<target>)" each"; the job `arm64` takes the image by the steps "Toolchain image from the cache", "Toolchain image" with `STAGES: toolchain` and "Toolchain image into the cache", the file `${{ runner.temp }}/novatrail-toolchain.tar` and the key `toolchain-${{ runner.os }}-${{ runner.arch }}-${{ hashFiles('.devcontainer/Dockerfile') }}`, and runs `just ci-test` as the job `check` does; the job `check` uploads `coverage/go.out` as the artifact `go-coverage`, "A recipe that wrote none fails "Go / test"";
  - `release.yml` — "The script builds the stage toolchain-core, Node, just, Go and buf, which every recipe of a release needs and no more, by TOOLCHAIN_STAGE"; it signs and checks the files several at once;
  - `codeql.yml` — Go "from go.dev's archive by ARG GO_VERSION and GO_SHA256_AMD64", "each once as ARG <name>=<value>" in the stage `toolchain` and the stages it is built from, with no image.
- `go help testflag` of Go 1.27.4 in the container, read 2026-11-16: "-covermode set,count,atomic … The default is "set" unless -race is enabled, in which case it is "atomic". … set: bool: does this statement run? … atomic: int: count, but correct in multithreaded tests; significantly more expensive."

### 5.9. The shared workflows and the action `toolchain` of `.github` v0.1.33

- the release [v0.1.33](https://github.com/altnovatrail/.github/releases/tag/v0.1.33), commit f1c7e9e, read 2026-11-18 in a clone in a temporary folder: immutable, of 2026-11-18 01:35 UTC, the latest (`gh release list -R altnovatrail/.github`, `gh release view v0.1.33 --json isImmutable`); `gh release verify v0.1.33 -R altnovatrail/.github` — "✓ Release v0.1.33 verified!". Between v0.1.31 and v0.1.33 (`git diff --stat v0.1.31 v0.1.33 -- .github/workflows .github/actions`) of what a caller calls only `go.yml` changed, with its `on.workflow_call` the same, and `.github/actions/toolchain/action.yml` is new:
  - `docs/decisions.md` of `.github`, the entry on the composite action — "`.github` has a composite action, `.github/actions/toolchain`, "Toolchain image and Go caches": the image of the caller's stage `toolchain` and the caches of Go for a job of the caller's own that runs Go in that stage, by the steps that the job `arm64` of the shared `go.yml` had, and `arm64` takes the action itself"; its step "Toolchain image" sets `SAVE: ${{ github.event_name == 'push' && github.ref == 'refs/heads/main' }}`, and "Toolchain image into the cache" saves the image on a miss on a push to `main`;
  - the same log, the entry on the caches of Go — "on a miss of the exact key every restore of the Go caches — in `check`, `target`, `fma-diff` and `rules`, and in the composite action … — takes the newest entry of the same check and Dockerfile by the restore key `go-<OS>-<architecture>-<check>-<hash of .devcontainer/Dockerfile>-`, the key without the hash of `go.sum`".

## 6. Other NovaTrail repositories

No records: a file of another NovaTrail repository is git, a link to the repository at its commit beside the claim.

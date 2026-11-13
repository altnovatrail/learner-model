# learner-model — sources

> The external facts that the documents of the repository's `docs/` rest on: standalone's model and code, research, the documentation of Go, protobuf, GitHub and Codecov, the documents, code and releases of other NovaTrail repositories. Here lie their records — the address, the version or the date of the page, the date read and the quote in the language of the source, — and [architecture.md](architecture.md) holds the claim and a reference to the record by its number, as to a section: `sources 1.1` (RULES.md, The bar; R143 spec).
>
> A record is a subsection "N.M" of a topic: one source, or one reading of several sources for one decision. The number of a record does not change: a new record goes at the end of its topic, a new topic at the end of the file. A source read again is an item with its date in the same record, and the earlier text is not edited. The model is sources spec. The results of measurements are not here but in [measurements.md](measurements.md).
>
> Records 1.1–6.10 were moved here word for word on 2026-11-13 from [architecture.md](architecture.md), section 4 "Sources", by task T49 (R35): the heading of a record is the name of the former item, and its first line names the subsections of architecture and the decisions for which the item was read. "Section N" inside a record is a section of the text it refers to. In architecture a pointer remains in the place of the former section: the decision log and the archives refer to it, and they are not rewritten.

## 1. MathTrail: standalone

### 1.1. Standalone's learner model

For architecture 2.1–2.5.

- [mathtrail/mathtrail-standalone, `SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md), the branch `main`, commit c6cbeb2, read 2026-10-19:
  - section 2.1 — "a Rasch model with a guessing floor, updated Elo-style after every answer" and "P = 0.2 + 0.8 · σ(θ + δ_topic − β)";
  - section 2.2 — "Only correctness enters the formula" and "The hint, the pace and the run of failures are recorded (04-profile) and change what is chosen next (section 3) — never the rating"; on the exception — "One kind of task is read for more than its correctness: a probe";
  - section 2.5 — "mastered" is set if the learner at the cautious level would answer the level's middle task "with **P ≥ 0.775**"; it holds at the level — "The topic counts as mastered only while its recommended point (2.3) stays at that level or below"; and when the rule changes, "The masteries of the earlier rule are cleared";
  - section 2.8 — the test vectors of the rating package: "The implementation computes in `float64` and must match the four decimals printed here to within 1e-4".

### 1.2. Standalone's profile

For architecture 2.4 and 3.5.

- [`docs/architecture/04-profile.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/architecture/04-profile.md), the same commit, read 2026-10-19:
  - on the summary by topics — "This is the part that is never pruned — it is the summary that lets the history window stay short"; on the window — "`recent` — the last **20** entries, answers and skipped tasks alike";
  - the section Versions and migrations — "Reading an older version runs the migration chain in memory — pure functions, one per step, each with its golden file in `testdata/`".

### 1.3. Standalone's harness

For architecture 3.2.

- [`docs/learners.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/learners.md), the same commit, read 2026-10-19 and 2026-11-02: the course of a run — "A task is answered two minutes after it is issued, the next is issued at once, and every tenth answer starts a new day"; the table of generators, G3-drop — "their level falls by 1.0 at once, at an answer from the 50th to the 150th, as after a long break".

### 1.4. Standalone's model code

For architecture 2.1, 2.5.

- [`internal/domain/rating/`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/rating), the same commit, read 2026-10-19: the rating package also holds the corridor, `corridor.go`; the files without `_test` import only `cmp`, `math` and `slices`.

### 1.5. Standalone's learner model in detail

For architecture 3.1.

- [`SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md), commit c6cbeb2, read 2026-11-02. Checked against `main` 89986bc of 2026-11-01: the package `internal/domain/rating` has not changed since c6cbeb2, and in SPEC 2 only the paragraphs on probes (R311 standalone) and on tasks from a picture (R322 standalone) changed, with no formulas:
  - section 2.1:
    - "β = the shift of its level + (difficulty − 3)";
    - "Each level stands 2.5 above the one before it, so the two hardest difficulties of a level overlap the two easiest of the next and no two points coincide" and "2.5 is a guess with no data behind it";
    - on the start — "It is written into `ratings.start` when the profile is created, and θ starts there. The grade decides nothing else: changed later, it moves neither θ, nor θ₀, nor the trial series of 2.2.1" (О-56 standalone);
  - section 2.2:
    - "θ += K_θ · (S − P)";
    - "K_δ = K₀ / (1 + a · n)" and "K_θ = max(K₀ / (1 + a · n), 0.05)" with K₀ 0.2 for θ and 0.4 for δ and a = 0.05;
    - "β is not updated in v1";
  - section 2.2.1:
    - "the first five answers are a trial series";
    - "a normal prior centred on the start times the chance of each answer given; its spread is σ₀ = 2.5";
    - "the corrections of the topics do not move";
    - "A grid rather than Newton's method from θ₀, because L can have two peaks";
  - section 2.5:
    - "the topic becomes mastered at the highest level it is taught at, at or below the level of the task, whose middle task — difficulty 3, β = the level's shift — the child would answer with P ≥ 0.775 at the cautious level λ = θ + δ − √v";
    - "v = 1 / (1/σ₀² + I · n) + 1 / (I · (20 + m))" and "0.9383 | ln(0.71875 / 0.28125)";
    - "mastery never moves down";
    - "On the learners' bench 3.7 % of these masteries are false, where the earlier rule's were 63.5 %, and they come in 4.6 answers against 9.9" (R187 standalone).

### 1.6. Standalone's model code in detail

For architecture 3.1.

- [`internal/domain/rating/`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/rating), commit c6cbeb2, read 2026-11-02:
  - `ladder.go` — `levelShift = 2.5`, and above `Start` — "A grade outside 1 to 6 starts from the nearest level's shift rather than from no number at all";
  - `trial.go` — `TrialAnswers = 5`, `startSpread = levelShift`, `gridStep = 0.05`;
  - `mastery.go` — above `MasteredAt`: "the highest of the levels it is taught at, at or below the level of the task answered, whose middle task the child at the cautious level would answer at the corridor's middle or better".

### 1.7. Standalone's decisions on the model

For architecture 3.1.

- [`docs/decisions.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/decisions.md), commit c6cbeb2, read 2026-11-02:
  - R79 standalone — "The grade picks the child's starting point θ₀, the shift of its level, which is written into `ratings.start` when the profile is created and never moved again";
  - R187 standalone — "Mastery is declared only above the level held, and lost after two wrong answers in a row, as before".

### 1.8. Answers, hints, skips and traps in standalone

For architecture 3.3.

- [`SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/SPEC.md), [`docs/decisions.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/decisions.md) and [`docs/architecture/04-profile.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/architecture/04-profile.md), the branch `main`, commit 89986bc, read 2026-11-02:
  - SPEC 2.2 — "Only correctness enters the formula" and "The hint, the pace and the run of failures are recorded (04-profile) and change what is chosen next (section 3) — never the rating, which stays an honest estimate of what the child knows"; on the probe — "A right answer to a probe lifts its topic" (R299 standalone) and "Whether a task is a probe is read off where the child stands when it is answered";
  - SPEC 2.5 — "The hint would make the answer measure the hint, and five answers keep a lucky start from mastering a topic";
  - R93 standalone — ""I don't know" is an answer, and a wrong one" and "a child who has seen the solution has not solved the task; counting it as wrong keeps the rating an honest estimate";
  - R98 standalone — on a skipped task: "Nothing that reads answers reads it — not the rating, the trial series, the misconception map, the streaks or the rule" and "punishing that would punish a stray tap";
  - R101 standalone — "The answer is recorded once";
  - `04-profile.md`, the field `traps` of a topic's summary — "How often this child fell for each trap in this topic, ever. The rule picks the two most frequent".

### 1.9. The log and versions in standalone

For architecture 3.5.

- [mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc), commit 89986bc on `main`, read 2026-11-02 through `gh api`:
  - `docs/decisions.md`, R165 standalone, among what was rejected — "the levels worked back from the answers, which breaks the day the step of the rating changes";
  - `docs/architecture/04-profile.md`, Versions and migrations, — "Reading a newer version is refused", because "two revisions are live during every Cloud Run rollout, and an old instance rewriting a new file would quietly drop whatever it did not understand", and "Unknown fields at the current version are ignored on read and are not written back".

## 2. Learner models and the scoring of answers

### 2.1. The Rasch model

For architecture 3.1.

- G. Rasch, [On General Laws and the Meaning of Measurement in Psychology](https://digitalassets.lib.berkeley.edu/math/ucb/text/math_s4_v4_article-22.pdf), Proceedings of the Fourth Berkeley Symposium on Mathematical Statistics and Probability, vol. 4, 1961, pp. 321–333, read 2026-11-02 from the text layer of the scan:
  - section 3 — "For the responses to the items of an intelligence test the following model is suggested… The probability that person number v gives a correct answer to item number i is" — then formula (3.1). The text layer renders the formula with errors, so it is not repeated here: that the Rasch model and Elo have one form is said by Pelánek 2016, below;
  - Rasch's 1960 book has no open copy, and it was not read. Standalone cites the 1961 article for the same reason ([`research/literature/learner-models.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/research/literature/learner-models.md), commit c6cbeb2).

### 2.2. The guessing floor

For architecture 3.1.

- A. Birnbaum, [Some Latent Trait Models and Their Use in Inferring an Examinee's Ability](https://faculty.ucmerced.edu/jvevea/classes/290_21/readings/week%209/Birnbaum.pdf), in F. M. Lord, M. R. Novick, Statistical Theories of Mental Test Scores, 1968, ch. 17, read 2026-11-02 from the text layer of the scan:
  - section 17.3 — "it further assumes that if he does not know it he will guess, and, with probability c_g, will guess correctly".

### 2.3. Elo in learning

For architecture 3.1 and 3.2.

- R. Pelánek, [Applications of the Elo Rating System in Adaptive Educational Systems](https://www.fi.muni.cz/~xpelanek/publications/CAE-elo.pdf), Computers & Education 98 (2016), 169–179; the authors' version, read 2026-11-02:
  - section 2.2:
    - the update — "θs := θs + K(correct_si − P(correct_si = 1))";
    - "In this formulation the Elo rating system is closely related to the Rasch model (one parameter model) used in IRT. These two approaches share the same functional form of the model for predicting the probability of a correct answer, they differ only in the parameter estimation procedure";
    - on guessing — "for question with k options the probability of a correct answer becomes" 1/k + (1 − 1/k) / (1 + e^(−(θs − di)));
  - section 2.3 — "An intuitive improvement used in most Elo extensions is to use an "uncertainty function" instead of a constant K"; on Math Garden — "Klinkenberg et al. (2011) use uncertainty initialized as U = 1 and updated after each answer as follows:" U := U − 1/40 + D/30, "where D is a time from a previous attempt (in days)". The article Klinkenberg et al. 2011 was not read;
  - section 4.3 — "It is also useful to incorporate forgetting, particularly when we are modeling declarative knowledge (in contrast to procedural skills, which are the typical focus of student modeling)";
  - Elo's 1978 book has no open copy, and it was not read. Standalone does not cite it either.

### 2.4. The global level and the correction

For architecture 3.1.

- J. Nižnan, R. Pelánek, J. Řihák, [Student Models for Prior Knowledge Estimation](https://files.eric.ed.gov/fulltext/ED560514.pdf), Proceedings of the 8th International Conference on Educational Data Mining, 2015, pp. 109–116, read 2026-11-02:
  - section 2.1 — "We use an uncertainty function U(n) = α/(1 + βn), where n is the number of previous updates to the estimated parameters";
  - section 2.3 — "In addition to the global skill θs the model now uses also the concept skill θsc … we just correct the global skill by the concept skill".

### 2.5. The first attempt

For architecture 3.3.

- [KDD Cup 2010: Educational Data Mining Challenge, Data Format](https://pslcdatashop.web.cmu.edu/KDDCup/rules_data_format.jsp), read 2026-11-02: "Correct First Attempt: the tutor's evaluation of the student's first attempt on the step—1 if correct, 0 if an error"; an error of the first attempt there also counts a request for a hint — "the first attempt for the step was an error (incorrect attempt or hint request)".

### 2.6. Testlets

For architecture 3.3.

- E. T. Bradlow, H. Wainer, X. Wang, A Bayesian Random Effects Model for Testlets, Psychometrika 64 (1999), 153–168, [doi:10.1007/BF02294533](https://doi.org/10.1007/BF02294533). The article is closed access; the abstract was read, through OpenAlex, 2026-11-02: "all items given to an examinee are unlikely to be conditionally independent (given examinee proficiency). Models that assume conditional independence will overestimate the precision with which examinee proficiency is measured. Overstatement of precision may lead to inaccurate inferences such as prematurely ending an examination in which the stopping rule is based on the estimated standard error of examinee proficiency".

### 2.7. The design effect and the effective sample size

For architecture 3.3.

- G. Kalton, J. M. Brick, T. Lê, [Estimating components of design effects for use in sample design](https://unstats.un.org/unsd/hhsurveys/pdf/chapter_6.pdf), chapter VI in Household Sample Surveys in Developing and Transition Countries, UN, ST/ESA/STAT/SER.F/96, March 2005, read 2026-11-02 from the text layer:
  - "The effective sample size is the size of an unrestricted sample that would yield the same level of precision for the survey estimate as that attained by the complex design";
  - the cluster effect — formula (18), 1 + (b − 1)ρ, — "depends on two factors: the subsample size within selected PSUs (b) and the intra-class correlation (ρ)". The text layer renders the formulas and the letter ρ with errors, so the formula is written here by its meaning;
  - the chapter traces the term "design effect" to Kish 1965 — the book was not read, and nothing is claimed from it.

### 2.8. Elo with several skills

For architecture 3.3.

- Pelánek 2016, above, section 2.4, read 2026-11-02: the multidimensional extension takes the case where "each item belongs to a single knowledge component", and after an answer updates every skill j with the correlation c_ij: θ_sj := θ_sj + c_ij · K · (correct − P). There is no division of the step among the primary skills of one item there, and the 1/J in section 3.3 is derived from the derivative, not taken from there.

## 3. Forgetting and repetition

### 3.1. Spaced repetition: FSRS

For architecture 3.2.

- [The Algorithm](https://github.com/open-spaced-repetition/awesome-fsrs/wiki/The-Algorithm), the wiki of `open-spaced-repetition/awesome-fsrs`, commit 7bcf7dc, read 2026-11-02:
  - the notation — "$S$: Stability (interval when $R=90%$)";
  - FSRS-4.5 — R(t, S) = (1 + FACTOR · t / S)^DECAY, "where $R(t,S)=0.9$ when $t=S$", and "In FSRS-4.5, $DECAY=-0.5$ and $FACTOR=\cfrac{19}{81}$"; "This version uses 17 parameters", and for FSRS-6 — "This version uses 21 parameters";
  - FSRS v1 and v3 — an exponential law, R(t, S) = 0.9^(t/S);
  - FSRS v4, the growth of stability after a successful repetition — "The larger the value of $S$, the smaller the $SInc$ value. This means that the higher the stability of the memory, the harder it becomes to make the memory even more stable" and "The smaller the value of $R$, the larger the $SInc$ value. This means that the spacing effect accumulates over time".

### 3.2. HLR and the Leitner system

For architecture 3.2.

- B. Settles, B. Meeder, [A Trainable Spaced Repetition Model for Language Learning](https://aclanthology.org/P16-1174.pdf), Proceedings of the 54th Annual Meeting of the Association for Computational Linguistics, 2016, pp. 1848–1858, read 2026-11-02:
  - section 3.2 — Leitner boxes with intervals of 1, 2, 4 days and on: "All cards start out in the 1-day box", and on an error a card moves to a box with a shorter interval;
  - section 3.3 — "This posits that memory decays exponentially over time": p = 2^(−Δ/h), where h is "the half-life or measure of strength in the learner's long-term memory"; and "Under the assumption that half-life should increase exponentially with each repeated exposure (a common practice in spacing and lag effect research)" the half-life is a function of features whose weights are learned from Duolingo's data.

### 3.3. SM-2

For architecture 3.2.

- P. Wozniak, [SuperMemo 2: Algorithm](https://super-memory.com/english/ol/sm2.htm), from the 1990 master's thesis, on the site since 1998-05-10, read 2026-11-02: "With all items associate an E-Factor equal to 2.5"; "I(1):=1", "I(2):=6", "for n>2: I(n):=I(n-1)*EF", the interval — "(in days)"; "If the quality response was lower than 3 then start repetitions for the item from the beginning without changing the E-Factor".

### 3.4. The forgetting curve

For architecture 3.2.

- H. Ebbinghaus, [Memory: A Contribution to Experimental Psychology, chapter VII](https://psychclassics.yorku.ca/Ebbinghaus/memory7.htm), 1885, translation 1913, read 2026-11-02, § 29, Discussion of Results: "forgetting would be very rapid at the beginning of the process and very slow at the end" and "the quotients of the amounts retained and the amounts forgotten were inversely as the logarithms of the times".

### 3.5. The form of forgetting

For architecture 3.2.

- J. T. Wixted, E. B. Ebbesen, On the Form of Forgetting, Psychological Science 2 (1991), 409–415, [doi:10.1111/j.1467-9280.1991.tb00175.x](https://doi.org/10.1111/j.1467-9280.1991.tb00175.x). The article itself was not read: the open copy is a scan without a text layer. The abstract was read, through OpenAlex, 2026-11-02: "In all cases, the course of forgetting was best described by a simple power function of time relative to five other reasonable alternatives (linear, exponential, exponential-power, hyperbolic, and logarithmic)".

### 3.6. The forgetting of skills

For architecture 3.2.

- W. Arthur Jr., W. Bennett Jr., P. L. Stanush, T. L. McNelly, Factors That Influence Skill Decay and Retention: A Quantitative Review and Analysis, Human Performance 11 (1998), 57–101, [doi:10.1207/s15327043hup1101_3](https://doi.org/10.1207/s15327043hup1101_3). The article is closed access; the abstract was read, through OpenAlex, 2026-11-02: "Results indicate that there is substantial skill loss with nonpractice or nonuse, with the amount of skill loss ranging from an effect size (d) of -0.01 immediately after training to a d of -1.4 after more than 365 days of nonuse" and "Physical, natural, and speed-based tasks were less susceptible to skill loss than cognitive, artificial, and accuracy-based tasks".

### 3.7. Intervals in mathematics

For architecture 3.2.

- D. Rohrer, K. Taylor, [The Effects of Overlearning and Distributed Practice on the Retention of Mathematics Knowledge](https://files.eric.ed.gov/fulltext/ED505642.pdf), Applied Cognitive Psychology 20 (2006), 1209–1224; the authors' manuscript in ERIC ED505642, read 2026-11-02 from the text layer. The subjects are university students, "116 undergraduates" in the first experiment. General Discussion — "In Experiment 1, distributing 10 practice problems across two sessions instead of massing all 10 problems in the same session had no effect on one-week test scores but virtually doubled four-week test scores. In Experiment 2, increasing the number of problems solved in a single session from three to nine had virtually no effect on test scores on either the one- or four-week test".

### 3.8. The ACT-R activation model

For architecture 3.2.

- P. I. Pavlik Jr., J. R. Anderson, Practice and forgetting effects on vocabulary memory: an activation-based model of the spacing effect, Cognitive Science 29 (2005), 559–586, [PubMed 21702785](https://pubmed.ncbi.nlm.nih.gov/21702785/). The article itself was not read: the authors' copy did not open. The abstract was read, 2026-11-02: "each time an item is practiced it receives an increment of strength but that these increments decay as a power function of time" and "This mechanism limits long-term benefits from further practice at higher levels of activation and produces the spacing effect".

### 3.9. Pelánek's correction for time

For architecture 3.2.

- R. Pelánek, [Modeling Students' Memory for Application in Adaptive Educational Systems](https://www.fi.muni.cz/~xpelanek/publications/memory-geography.pdf), Proceedings of the 8th International Conference on Educational Data Mining, 2015, pp. 480–483; the authors' version, read 2026-11-02 from the text layer:
  - section 2.2 — "instead of P(m) to use P(m + f(t)), where t is the time (in seconds) from the last attempt", and "the value f(t) corresponds to temporal increase in memory activation due to (short) time from previous exposure of an item". The data gave no simple function f, and it is a step function, fitted to the data, and past the last step — "0 otherwise";
  - the data — places on a map, slepemapy.cz, section 3.

### 3.10. Glicko

For architecture 3.2.

- M. E. Glickman, [The Glicko system](http://www.glicko.net/glicko/glicko.pdf), undated, read 2026-11-02: "In the Glicko system, a player's rating changes only from game outcomes, but his/her RD changes both from game outcomes and also from the passage of time when not playing"; after t periods without games — RD = min(√(RD²_old + c² · t); 350).

## 4. Go

### 4.1. Floating-point operations in Go

For architecture 2.5.

- [The Go Programming Language Specification](https://go.dev/ref/spec#Floating_point_operators), the section Floating-point operators, Language version go1.27 (May 26, 2026), read 2026-10-19: "An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements, and produce a result that differs from the value obtained by executing and rounding the instructions individually".

### 4.2. Explicit conversion against fusion in Go

For architecture 2.5 and 3.5.

- [The Go Programming Language Specification](https://go.dev/ref/spec#Floating_point_operators), the section Floating-point operators, Language version go1.27 (May 26, 2026), read 2026-11-02: "An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding" and the example "r = float64(x*y) + z" among those where FMA is forbidden.

### 4.3. The sources of Go 1.27.3

For architecture 2.5 and 3.5.

- the directory `src` of Go 1.27.3 in the container (`go env GOROOT`), read 2026-11-02:
  - `math/exp_asm.go` — `//go:build amd64 || arm64 || loong64 || riscv64 || s390x` and `const haveArchExp = true`: `math.Exp` has its own assembly implementation on these architectures, and on the rest, WASM included, one in Go (`math/exp_noasm.go`);
  - `math/exp_amd64.go`, line 11, — `var useFMA = cpu.X86.HasAVX && cpu.X86.HasFMA`: the path in `math/exp_amd64.s` is chosen by the processor at run time. It can be turned off through `internal/cpu/cpu_x86.go` — `option{Name: "fma", Feature: &X86.HasFMA}`, that is, `GODEBUG=cpu.fma=off`;
  - `math/log_asm.go` — `//go:build amd64 || s390x`;
  - `math/pow.go`, line 128, — `a1 = Exp(yf * Log(x))`;
  - the compiler's rules: in `cmd/compile/internal/ssa/_gen/ARM64.rules` — `(FADDD a (FMULD  x y)) && a.Block.Func.useFMA(v) => (FMADDD  a x y)`, and in `AMD64.rules` — `(ADDS(S|D) (MULS(S|D) x y) z) && buildcfg.GOAMD64 >= 3 && z.Block.Func.useFMA(v) => (VFMADD231S(S|D) z x y)`.

### 4.4. Fusion, conversions and order in Go

For architecture 3.6.

- [The Go Programming Language Specification](https://go.dev/ref/spec), Language version go1.27 (May 26, 2026), the copy `doc/go_spec.html` of Go 1.27.3 in the container, read 2026-11-02:
  - Floating-point operators — "An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements" and "An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding"; among the examples where fusion is allowed — `t = x*y; r = t + z`;
  - Conversions — "In all non-constant conversions involving floating-point or complex values, if the result type cannot represent the value the conversion succeeds but the result value is implementation-dependent";
  - For statements with range clause — "The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next".

### 4.5. The compiler and standard library of Go 1.27.3

For architecture 3.6 and 3.7.

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

For architecture 3.6; the 2026-11-07 amendment to R09.

- the directory `src` in the container, read 2026-11-07: `go/build/build.go`, `goodOSArchFile`, — the name of a file up to the first dot, without the first word and without `_test`, ends in `_GOOS`, `_GOARCH` or `_GOOS_GOARCH`; the lists of systems and architectures are `KnownOS` and `KnownArch` in `internal/syslist/syslist.go`. A build with `-gcflags='-S -d=fmahash=n'` prints, at every place where a fusion rule would fire, the line "fmahash triggered [DISABLED]" with a hash, but without the file's name and line — where exactly is shown by the difference of the assembler.

### 4.7. Go and WebAssembly

For architecture 3.6.

- [Go Wiki: WebAssembly](https://go.dev/wiki/WebAssembly), read 2026-11-02: "Set `GOOS=js` and `GOARCH=wasm` environment variables to compile for WebAssembly"; "Go 1.21 introduced WASI as a supported platform. To build for WASI, use the `wasip1` port"; "At present, Go generates large Wasm files, with the smallest possible size being around ~2MB" and "Using `gz` compression reduces the ~2MB (minimum file size) example WASM file down to around 500kB"; "TinyGo supports a subset of the Go language".

### 4.8. protobuf for Go

For architecture 3.7.

- [`google.golang.org/protobuf/proto`, `MarshalOptions`](https://pkg.go.dev/google.golang.org/protobuf/proto#MarshalOptions), v1.36.12, read 2026-11-02: `Deterministic` — "the same message will always be serialized to the same bytes within the same binary", and "the deterministic serialization is NOT canonical across languages. It is not guaranteed to remain stable over time".

## 5. CI, GitHub and tools

### 5.1. GitHub's runners

For architecture 3.6.

- [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), read 2026-11-02: the labels `ubuntu-24.04-arm`, `ubuntu-22.04-arm` and `ubuntu-26.04-arm` are in the table of standard runners for public repositories, and "Use of the standard GitHub-hosted runners is free and unlimited on public repositories". The page does not name the processors of the x64 runners. `altnovatrail/learner-model` is public (`gh repo view`, 2026-11-02).

### 5.2. The shared Go wrapper

For architecture 3.6.

- [`.github/workflows/go.yml`](https://github.com/altnovatrail/.github/blob/af40f949676fb3a01abb0a1e6ffbc99ce824660b/.github/workflows/go.yml), `altnovatrail/.github`, `main` af40f94, read 2026-11-02 through `gh api`: "by four recipes of that repository's own — ci-lint, ci-test, ci-vuln and ci-licenses — one job each, in the toolchain stage of the caller's .devcontainer/Dockerfile, where the caller pins its Go"; the jobs run on `ubuntu-24.04`, without arm64 and without WASM.

### 5.3. The shared workflows of `.github` v0.1.15

For architecture 3.6; R13.

- the release [v0.1.15](https://github.com/altnovatrail/.github/releases/tag/v0.1.15), commit 89d3cee, read 2026-11-06; `gh release verify v0.1.15 -R altnovatrail/.github` — "Release v0.1.15 verified!":
  - `go.yml` — "ci-test writes the coverage of the tests as one Go coverage profile, coverage/go.out, and the job coverage sends it to Codecov by OIDC, with no secret"; `ci-licenses` gets `ALLOWED_LICENSES` — "the ASF's Category A as SPDX identifiers";
  - `codeql.yml` — "Besides actions, javascript-typescript and go are taken", and Go is read from the Go build of the `toolchain` stage; at v0.1.0 the input `languages` took only `javascript-typescript`;
  - `secrets.yml`, `dco.yml`, `documents.yml` and `scorecard.yml` did not change between v0.1.0 and v0.1.15 (`git diff --stat v0.1.0 v0.1.15`);
  - the texts of the calling files, `codecov.yml` and the order of the move to coverage — [`docs/architecture.md`](https://github.com/altnovatrail/.github/blob/main/docs/architecture.md) of `.github`, section 3.3, `main`, read 2026-11-06.

### 5.4. The shared workflows of `.github` v0.1.21

For architecture 3.6; R18.

- the release [v0.1.21](https://github.com/altnovatrail/.github/releases/tag/v0.1.21), commit 706f406, read 2026-11-08 in a clone in a temporary folder; `gh release verify v0.1.21 -R altnovatrail/.github` — "Release v0.1.21 verified!". With it, the earlier "on `ubuntu-24.04`, without arm64 and without WASM" above holds only without the input `same-bits`:
  - `go.yml` — "With same-bits, four more jobs hold the caller's bits the same on every build": ""Go / targets" runs the tests on the builds where the bits are most at risk, and builds and vets wasip1; "Go / fma-diff" compares the assembler of each package with fusion and without it, with a control of its own here; "Go / arm64" runs ci-test and the targets of its host on an arm64 runner; "Go / rules" holds the code to the rules of code by the tool tools/coderules of .github, of this workflow's own version, and the caller's lists, .github/coderules.json"; "Without the input the workflow is what it was";
  - `tools/coderules/lists.go` — the fields of the lists `code-imports`, `test-imports`, `test-math`, `requires` and `nested-modules`, "A field it does not know is a mistake";
  - between v0.1.15 and v0.1.21 `codeql.yml` and `scorecard.yml` only raised the versions of actions, `documents.yml` calls `spec/tools/toolchain/toolchain.sh`, and `release.yml` changed only the npm branch (`git diff v0.1.15 v0.1.21`).

### 5.5. The shared workflows of `.github` v0.1.24

For architecture 3.6; R34.

- the release [v0.1.24](https://github.com/altnovatrail/.github/releases/tag/v0.1.24), commit fab102d, read 2026-11-13 in a clone in a temporary folder; `gh release verify v0.1.24 -R altnovatrail/.github` — "Release v0.1.24 verified!"; the latest release (`gh release list -R altnovatrail/.github`, 2026-11-13). Between v0.1.21 and v0.1.24 (`git diff v0.1.21 v0.1.24 -- .github/workflows/`):
  - `go.yml` — ""Go / targets" runs the tests on the builds where the bits are most at risk, wasip1 among them": `GOOS=wasip1 GOARCH=wasm go test -count=1` with `-exec` — `node --stack-size=8192` and a runner on `node:wasi`, `new WASI({ version: "preview1", args, env: process.env, preopens: { "/": "/" }, returnOnExit: true })`; the same in "Go / arm64". Also, each job in the `toolchain` stage takes its image from the GitHub Actions cache by the key of the OS, the architecture and the hash of the Dockerfile, and builds it only on a miss, "so that an outage of what the Dockerfile downloads from does not fail the checks while the Dockerfile stays the same". The inputs and the permissions of the call are the same;
  - `documents.yml` — the same cache of the stage for `spec/tools/toolchain/toolchain.sh` through `TOOLCHAIN_IMAGE_FILE`, which the script reads since `spec` v0.13.0;
  - `release.yml` and `scorecard.yml` — only the links in comments; `codeql.yml`, `secrets.yml` and `dco.yml` did not change.

### 5.6. The tools of the checks and gopls

For architecture 3.6; R13, R14, R19.

- the Go module proxy, `https://proxy.golang.org/<module>/@latest`, read 2026-11-06: `golang.org/x/vuln` — v1.9.0 of 2026-10-11, `github.com/google/go-licenses/v2` — v2.0.1 of 2025-09-08, `golang.org/x/tools/gopls` — v0.24.0 of 2026-10-14. The licenses — by GitHub's API, `repos/<repository>/license`, 2026-11-06: govulncheck and gopls — BSD-3-Clause (`golang/vuln` and `golang/tools`), go-licenses — Apache-2.0 (`google/go-licenses`). Since `spec` v0.9.16 the lock file of gopls is the one of `spec`'s template: `.devcontainer/` from the output of the recipe `new-repo` for `learner-model` of the kind `go` differs from the former one only in comments, and `go.sum` is the same byte for byte (`diff -r`, 2026-11-09).

### 5.7. Codecov

For architecture 3.6; R13.

- read 2026-11-06: the Codecov app is installed for the organization on all repositories — `gh api orgs/altnovatrail/installations`, `"app_slug":"codecov","repository_selection":"all"`; `learner-model` is already in Codecov but not activated — `https://api.codecov.io/api/v2/github/altnovatrail/repos/learner-model/`, `"active":false,"activated":false`.

## 6. Other NovaTrail repositories

### 6.1. The log and the cache in `api`

For architecture 3.5.

- [altnovatrail/api](https://github.com/altnovatrail/api/tree/d896943), commit d896943 on `main`, read 2026-11-02 through `gh api`:
  - R06 api, «Что решено», — `api` stores a piece of evidence as «его тип как есть и столбцы для себя: ID профиля, ID попытки, дайджест пакета, модуль, время и порядковый номер. Содержания в свидетельстве нет», and the cache as «строка на профиль: состояние, версия `learner-model`, закреплённый релиз `modules` и номер последнего учтённого свидетельства»; the export — «журнал, текущее состояние и выборы с кодами причин»;
  - the 2026-11-02 amendments to R06 api — «Пока идёт выкладка, запросы может обслуживать и прежний образ: пометку новее своей — по релизу или окну — он не переписывает, а состояние для ответа считает, не записывая его» and «выбор движка хранит и своё время»;
  - R08 api — «Свидетельство о навыке без пары не идёт ни одному навыку, а о навыке, которого граф текущей версии не знает, — ждёт в журнале, пока следующая версия его не назовёт»;
  - architecture 3.7 api, «Свидетельство», — «Нужен ли весу свидетельства факт «в уроке» — попытка знает своё прохождение, — решает `learner-model`».

### 6.2. The repeat of a choice in `engine`

For architecture 3.5.

- [altnovatrail/engine](https://github.com/altnovatrail/engine/tree/08bc48a), commit 08bc48a on `main`, read 2026-11-02 through `gh api`: in RUN.md, T03 engine decides «что входит в запись выбора и что из журнала свидетельств нужно, чтобы повторить выбор по сохранённым версиям», and T04 engine sends the fields that the types of `learner-model` lack «строкой таблицы 3 «в реестр spec», к T05 learner-model». The detailing of `engine` has not begun.

### 6.3. The `sim` harness after the judgement `model-2`

For architecture 1 and 3.6; the 2026-11-13 amendment to R20; the 2026-11-13 amendment to R32.

- [altnovatrail/sim at `v0.1.13`](https://github.com/altnovatrail/sim/tree/v0.1.13), read 2026-11-13: the release is immutable, of 2026-11-13, the latest, the tag lightweight, on the commit `7c3201622165070519a489f7dfd2066dbf06d189`, the merge of pull request #58 of `sim` (`gh release view v0.1.13 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — "Release v0.1.13 verified!"). `gate.json` — `learner-model` at a063835 without a version, `engine` v0.1.6 at 57e0ab8, `modules` v0.2.5, `spec` v0.16.3, 69 numbers, the bands written anew — `"T47: the judgement model-2 chose theta, learner-model a063835, and the held-out learners confirmed it (judgements/model-2/report.md; R46)"`; `go.mod` — `learner-model` `v0.1.7-0.20261113040121-a063835087c7`. The report of the judgement is `judgements/model-2/report.md`: "Chosen and confirmed: theta." Between `v0.1.10` and `v0.1.13` the recipe `just gate` and `internal/gate/` did not change, `Learners` is 250, and in the README, "Gate" and "Taking a release" only "bench" became "simulator" (`git diff v0.1.10 v0.1.13`).

### 6.4. The `sim` harness with the sizes after `learner-model` v0.1.6

For architecture 1 and 3.6; the 2026-11-12 amendment to R20.

- [altnovatrail/sim at `v0.1.10`](https://github.com/altnovatrail/sim/tree/v0.1.10), read 2026-11-12: the release is immutable, of 2026-11-12, the latest, the tag lightweight, on the commit `ff921be033917605b7be9d8d194a3c41a90c280c` (`gh release view v0.1.10 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — "Release v0.1.10 verified!"). `gate.json` — `learner-model` v0.1.6 at 7d796e2, `engine` v0.1.6 at 57e0ab8, `modules` v0.2.5 — `olympiad-math` and `grade-3-math` 0.3.0, `english-essays` 0.2.0, — `spec` v0.15.0, 69 numbers with bands written anew on 250 learners, — "T46: the gate takes 250 learners a cell, not 150, after learner-model v0.1.6 made its readout faster". `internal/gate/gate.go` — `const Learners = 250`: "250 keep the slowest run of the job on 150, scaled, at about ten minutes at every caller, whose gate keeps no cache of Go, within the twelve that 300 would only just meet"; the recipe `just decisive` — `-learners 2000`. Between `v0.1.8` and `v0.1.10` there are two tags: `v0.1.9` — `time` and the leading question of format 0.4 (T42, R41 `sim`), the packages of `modules` v0.2.5 with them (T43, R42 `sim`) and the base on `learner-model` v0.1.6 (T45, R43 `sim`), and `v0.1.10` — the sizes (T46; the 2026-11-12 amendments to R31 and R33 `sim`). The README section "Gate", the input `GATE_LEARNER_MODEL` and the section "Taking a release" are the same (`git diff d957fc7 ff921be -- README.md justfile`): in the README only "Status" and the size of the deciding run in "Judging a candidate" changed.

### 6.5. The `sim` harness on `spec` v0.15.0

For architecture 1 and 3.6; the second 2026-11-11 amendment to R20.

- [altnovatrail/sim at `v0.1.8`](https://github.com/altnovatrail/sim/tree/v0.1.8), read 2026-11-11: the release is immutable, of 2026-11-11, the latest, the tag lightweight, on the commit `d957fc7e9ab411fb1f2e9442b99e2a888d8d6775` (`gh release view v0.1.8 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — "Release v0.1.8 verified!"). Between `v0.1.6` and `v0.1.8` — the move of `sim` to `spec` v0.15.0 (T38, R39 `sim`) and the judgement that holds a candidate of one library to the second library of the base (T40, R40 `sim`); in `gate.json` only the dependency `spec` changed, and the libraries, the packages, the bands and the recipe `just gate` are the same (`git diff v0.1.6 v0.1.8`).

### 6.6. The `sim` harness on `engine` v0.1.6

For architecture 1 and 3.6; the 2026-11-11 amendment to R20.

- [altnovatrail/sim at `v0.1.6`](https://github.com/altnovatrail/sim/tree/v0.1.6), read 2026-11-11: the release is immutable, of 2026-11-11, the latest, the tag on the commit `a97be529249042595e544da6192555fa783ff9ec` (`gh release view v0.1.6 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — "Release v0.1.6 verified!"). `gate.json` — `engine` v0.1.6 at 57e0ab8, `learner-model` v0.1.4 at acfc4f2, `modules` v0.2.0, `spec` v0.12.0, 150 learners a cell (T41 `sim`; the 2026-11-11 amendments to R26 and R31 `sim`). The recipe `just gate` writes the workspace by the hidden recipe `workspace`, and `GATE_LEARNER_MODEL` takes a folder, as before, or a commit by its full SHA; README, "Gate" — "A caller puts its own code in place of a library with `GATE_LEARNER_MODEL` or `GATE_ENGINE`: a folder, or a commit by its full SHA" (`git diff v0.1.2 v0.1.6`). `sim`'s measurement on the slowest runner — "Gate" of `learner-model` on 150 learners about 11.2 minutes against a budget of 15 (the 2026-11-11 amendment to R31 `sim`).

### 6.7. The `sim` harness on a new pin

For architecture 1 and 3.6; the 2026-11-09 amendment to R20.

- [altnovatrail/sim at `v0.1.2`](https://github.com/altnovatrail/sim/tree/v0.1.2), read 2026-11-09: the release is immutable, of 2026-11-09, the tag on the commit `e61f90aa1bd894d1ce8f35547e0757d38f28c0b2` (`gh release view v0.1.2 -R altnovatrail/sim`, `git ls-remote`, `gh release verify` — "Release v0.1.2 verified!"). The recipe `just gate` and the input `GATE_LEARNER_MODEL` are the same as at `v0.1.0`, and "Taking a release" of the README did not change (`git diff v0.1.0 v0.1.2`). `gate.json` — `modules` `v0.1.4`, `olympiad-math` 0.1.1, the cell `A-linked` on two maths modules, 100 learners a cell; the libraries in it are still `learner-model` `v0.1.0` at d2620a1 and `engine` `v0.1.0` at c6c5bb9, as in `go.mod`. `sim`'s own CI on pull request #40 of `sim` — "Gate / amd64" 11 minutes, "Gate / arm64" 6 minutes (RUN.md of `sim`, T30, «Сделано»). "Gate" here on pull request #41 — 7 min 10 s, of which the build of the `toolchain` stage is 23 s (RUN.md, T29).

### 6.8. The `sim` harness and its gate

For architecture 1 and 3.6; R20.

- [altnovatrail/sim at `v0.1.0`](https://github.com/altnovatrail/sim/tree/v0.1.0), read 2026-11-09: the release is immutable, of 2026-11-09, the tag on the commit `09e5a84f8a26fedcf5c8996bf81bf29c616917a2` (`gh release view v0.1.0 -R altnovatrail/sim`, `git ls-remote`). README, "Gate" — "A caller puts its own code in place of a library with `GATE_LEARNER_MODEL=<folder>` or `GATE_ENGINE=<folder>`", and "Taking a release" — "write both the tag and the SHA of its commit in its workflow… check out sim at that SHA, then run `just gate`… in the toolchain stage of its devcontainer"; the exit code is 0 when the gate is green, 1 when it is red or the input is rejected, 2 on a call the recipe does not know. `gate.json` names `learner-model` `v0.1.0` at d2620a1 and `engine` `v0.1.0` at c6c5bb9; there is no reusable workflow, on purpose (R05 `sim`, «Отвергнуто»). `sim`'s measurement on a runner — about four minutes on `ubuntu-24.04`, «кеш Go почти ничего не даёт» (RUN.md of `sim`, T24).

### 6.9. The Go module `spec`

For architecture 3.6.

- [`spec.go`](../spec/spec.go) in the submodule at v0.6.6: `//go:embed all:format`.

### 6.10. The neighbours on the build

For architecture 3.6.

- read 2026-11-02 through `gh api`:
  - `engine`, `main` 08bc48a: T05 engine decides «что нужно движку, чтобы однажды выбирать и на устройстве», and asks «проверять ли сборку для устройства в CI с первого кода»; T04 engine sends the fields that the types of `learner-model` lack «строкой таблицы 3 «в реестр spec», к T05 learner-model»; in section 2.4 — «у навыка — ID, название, описание, ярус, пререквизиты и оцениваемость»;
  - `sim`, `main` 067f624, docs/architecture.md, section 2.3, — «полосы стража в standalone записаны на amd64 при `GOAMD64` v1 и на процессоре с FMA. На arm64, на amd64 при `GOAMD64` v3 и выше и на процессоре без FMA последние цифры чисел могут разойтись, хотя и в пределах полос»;
  - `app`, `main` 375a09b, docs/architecture.md, section 3.1, — «Приложение — веб, для планшета и компьютера, без офлайна и нативных приложений», on React 19, Vite 8 and TypeScript 7 (R03 app).

# learner-model — architecture

> The detail of the card, architecture 3.8 spec, at the tag of `spec` that [CLAUDE.md](../CLAUDE.md) names. The card stays the summary: if the detail changes what the card says, the card is changed in `spec` by the same route as the questions to its registers (RULES.md, Across repositories). Section 1 is what the component does and does not do, its data and links; section 2, the card's items, decided, with what is deliberately deferred, and the questions of `spec` that reach the model; section 3, how it is built.
>
> Why it was decided so and not otherwise is in [decisions.md](decisions.md); the work plan is in [RUN.md](../RUN.md). The concepts are concepts spec. The documents and decisions of `spec` lie in `spec/docs/` and are cited with the suffix "spec".

## 1. Purpose and boundaries

**Does:**
- estimates how the learner masters each skill from the evidence. A piece of evidence is the outcome of an attempt, not of a single submission (concepts 11 spec; architecture 4.1 spec);
- keeps the skill's state: the estimate of mastery, its uncertainty, the number of pieces of evidence, the time of the last practice, the frequent traps, "mastered" or not (concepts 11 spec), and beyond that the stability and the failures in a row (R04, R06) — sections 3.1–3.3;
- takes forgetting into account: without practice a skill's retention fades and the probability of a correct answer falls, while spaced practice grows the skill's stability and forgetting slows down. Time does not remove "mastered" but gives a review date by which `engine` returns the skill (concepts 11, 12 spec; R04) — section 3.2;
- starts from standalone's model: a Rasch model with a guessing floor, Elo updates, a trial series, a cautious "mastered" (architecture 3.8 spec; concepts 11 spec; repositories 7 spec). The main differences from it are the level per module, not for the whole model (R03), and forgetting (concepts 11 spec; R04); the others are guessing, the weights of evidence and the rules in place of calibration, below. The model with a level per module is section 3.1, forgetting is section 3.2, the weights of evidence are section 3.3;
- takes a task's hardness by its band and its difficulty from 1 to 5 within the band, and starts the learner a band's width below the band to which the module manifest assigns the learner's grade, below the lowest band too. How far apart the bands are and where a learner outside the bands starts is `learner-model`'s (architecture 3.8 spec; R51 spec; R03) — section 3.1;
- holds the card's three rules in place of calibration: that start a band below, an adult's correction of the module's level — a fact of the evidence log that `api` appends, evidence of kind 3 — and a guard after three failures in a row (R157 spec; R45) — section 3.1, "The rules of R157 spec";
- takes guessing into account by the kind of answer, counted by the item's answer slots: 1/k for a choice among k options — by touch or by hand —, 0 for a free answer. In standalone it was hard-wired as 0.2 for five options (architecture 3.8 spec; concepts 11 spec; R03) — section 3.1;
- weighs a piece of evidence by the outcome of the attempt and by key provenance (architecture 3.8 spec; concepts 9.4, 11 spec; R06) — section 3.3:
  - the step takes the first verdict of each item, while a hint and other help only keep "mastered" from being set;
  - a series of items of one attempt weighs less than their sum (concepts 11 spec): a workbook page of ten items counts as 2.5 answers;
  - an error is evidence against the task's primary skills but not against its supporting ones (concepts 6 spec), and the answer gives the supporting ones nothing;
- knows a module's skill by its full ID — the module, the kind and the local ID, `novatrail/olympiad-math/skill/counting-gaps` (concepts 3 spec; R46 spec), and an anchor skill of the shared taxonomy by its ID, `core.<name>`. Whether the learner knows an anchor skill is `learner-model`'s to derive from the skills that teach it in the learner's modules (architecture 3.8 spec; concepts 4, 11 spec; R50 spec): the learner masters it when at least one skill that teaches it is mastered (R07) — section 3.4;
- carries progress over when a module gets a major version, by the table `retired_skills` of the skill graph: the evidence of a retired skill goes to each skill that replaced it; `[]` means the progress goes nowhere (architecture 3.8 spec; concepts 3 spec; R49 spec). The table is closed: the server carries progress over (R49 spec). `api` cuts off the evidence of `wrong_key` and hands over the log and the graphs of the current versions of the log's modules with `retired_skills`, in the types of `learner-model`, and `learner-model` does the transfer, by the card's conditions: a module's skill that the graph of the current version knows neither in `skills` nor in `retired_skills` goes nowhere; an anchor skill goes as it is if the dictionary of some format version in `api`'s window knows it, and otherwise nowhere; a trap's skills go by the same rules; the recomputation opens no packages of earlier versions (architecture 3.8 spec; the 2026-11-02 amendment to R80 spec; R08 api). Each replacement gets the retired skill's share in full, and derives "mastered" from the transferred evidence (R08) — section 3.5;
- does not count the evidence of attempts on a version of a module that the module index withdrew for a wrong key, `wrong_key`: the model is recomputed from the log without them (architecture 3.8 spec; R63 spec). `api` cuts them off, by the package digest, before handing the log to the library (R08 api; the 2026-11-02 amendment to R80 spec): the library gets the log already without them.

**Does not:**
- does not read module packages or the database: the skill graph and the evidence are handed to it by `api` (architecture 3.8 spec);
- does not choose the next lesson or task. `engine` does that, and the learner model is one of its inputs (concepts 12 spec; R03 spec). Lessons are a rule of `engine` (R05 engine; architecture 3.8 spec; R107 spec). A worked example in a lesson is not a task and does not serve as evidence, and a task of a lesson is the same kind of attempt as one outside a lesson (R58 spec; architecture 3.5 spec);
- does not call AI and gives no verdict. `module-kit` gives the verdict, and what comes here is the attempt's outcome (architecture 1 spec; R03 spec);
- does not show the model. The skill map, ratings, ranks, "strengths" and "what to develop" are derived from it, and they do not affect the model back (concepts 11 spec). `skill-map` draws the map, from ready data from `api`: it derives the skill's state from the graph and the model's ready marks (R03 skill-map; architecture 3.3 spec). A module's rank — the number of the band in which the module's level stands — would need a field of the contract and a rule of `learner-model` by which the level falls into a band; there are no ranks in the first slice, and no such rule either, until О-13 decides that ranks are visible (the 2026-11-07 addition to О-13).

**Data and links:**
- it has no data of its own: these are computations without input or output, and the time comes from outside (architecture 3.8 spec; repositories 1 spec). The model takes no seed: it needs no randomness and imports no generator, and no seed comes to `engine` either (architecture 3.8 spec, "Common"; R09; R59, R107 spec) — section 3.6. If randomness is ever needed, the generator and the derivation of seeds come from `spec`: they are one for the whole platform — PCG, as in Go's `math/rand/v2`, and the places of randomness, as in standalone's harness ([format text](../spec/format/README.md), section 22; R27, R59 spec). There are no learners' data and no secrets in the repository, as in no repository: the examples and tests are synthetic (RULES.md, Hard rules);
- the learner model is the learner's data. It lives in `api`'s database, the only store of learners' data on the server among the platform's parts (R80 spec; architecture 1, 5 spec; R15 spec); it does not go to AI (privacy 4 spec). It can be exported in an open format and deleted (R09 spec); the export format is [the part `export/`](../spec/export/README.md) of `spec`: the log and the state are JSON by the `.proto` schema of this library, and the `.proto` files themselves lie beside them (R135 spec; section 3.5);
- what is kept is the append-only evidence log, and it is the source of truth, while the skills' state is a cache with a mark of the version of `learner-model`, the release of `modules` and `api`'s window of format versions (R80 spec and its fourth amendment of 2026-11-02). There is one piece of evidence per finished attempt: what is in it and how much it weighs is `learner-model`'s, and `api` keeps it as it is. The condition for the library is one: the state appended with a new piece of evidence matches the state recomputed from the log (R80 spec). What a piece of evidence carries and how the appended state matches the recomputed one bit for bit is R08 — section 3.5; how `api` keeps it, R06 and R08 api;
- `api` calls the library, in its own process: on an attempt's outcome the evidence goes into `learner-model`, then `engine` chooses the next lesson or task (architecture 3.5, 4.1 spec). The library itself calls no one;
- the engine choice keeps the versions of `engine` and `learner-model`, the package digests, the seed and the reason code, the release of `modules`, the number of the log's last piece of evidence it was made by, the time of the choice, the grade by which the model computed the starting point of modules without evidence, whether the pool was limited to tasks without AI, the skill from the skill map, the age category at the time of the choice and a mark of which of the digests are the engine's modules; the seed in it is empty while the engine takes no randomness. The choice is repeated from the evidence log, without AI (architecture 1 spec; R80, R107 spec);
- at build time it depends only on `spec` — on its documents and rules in the submodule —, and `engine`, `sim` and `api` are built with it (repositories 2.1, 3 spec). The code imports only Go's standard library, without the Go module of `spec`, gives the same bits on any build, WebAssembly included, and builds for the device as it is; CI checks it (R09) — section 3.6. The graph types are its own and simple, not from `module-kit`, so the library can be built for the device too (architecture 3.8 spec; Т-10). `learner-model` defines them, and `engine` takes them from it (R28 spec); what they are and how the evidence and the state are encoded is R10 — section 3.7;
- `sim` imports it together with `engine` and runs it on thousands of virtual learners of different populations (architecture 3.9 spec; concepts 15 spec). The library's CI runs `sim`'s gate, the check "Gate", at the tag and the commit of `sim` that `SIM_TAG` and `SIM_COMMIT` of `.github/workflows/ci.yml` name, with `go.work` written on the fly and not committed: `sim`, the pull request's code and `engine` at the commit from `sim`'s gate file. The gate runs on every pull request and every commit of `main` and holds the merge: a change that worsens a criterion is not merged (repositories 4 spec; R103 spec; R20). It is two-sided: a number outside its band turns red for the better too, so a change that moves the numbers first records the bands in `sim` (section 3.6). The model is changed only if the new variant has passed the harness by a criterion recorded before the first run (concepts 11 spec; R163 standalone); the first version is the baseline, not a candidate (R12). A new format version of a module `learner-model` takes with a new `sim` by a move of its own, and keeps its gate on the earlier tag (repositories 4 spec; R44 spec);
- stage 1, "Format and core" (product 11 spec). It comes together with `engine` and `sim`: the harness needs both the model and the choice rule, and the harness and the criterion come before the first candidate (repositories 6 spec). Where the work stands is [RUN.md](../RUN.md); what may come next, [run/EPICS.md](../run/EPICS.md).

## 2. To decide when detailed

The card's items "To decide when detailed" that are `learner-model`'s, a subsection for each, 2.1–2.6: each is decided by an entry of the log, built as section 3 says, and names what is deliberately deferred; the card carries the decisions (architecture 3.8 spec). Subsections 2.7 and 2.8 are not items of the card but questions of `spec` that reach the model — О-25 and format 0.4 —, each decided as its subsection says.

`learner-model` and `engine` have one card (architecture 3.8 spec). It names the level and forgetting — subsections 2.1 and 2.2 (R03, R04) — as decided rows of `learner-model`, and has seven items:
- the weights of evidence, subsection 2.3 (R06);
- the repetition rule, an item of `engine` (R06 engine): repetition is a task on a mastered skill whose date has come, on the module's turn of repetition — once in three of its choices —, and when the others have nothing to give, outside the turn too; the rule does not read the daily goal (R107 spec). The model gives it the review date — subsection 2.2 (R04);
- evidence: what is in it and how much it weighs, on the condition that the state appended with a new piece of evidence matches the state recomputed from the log (R80 spec), — subsection 2.4 (R08);
- a task — not a template — whose solution has already gone to the learner's device: the engine does not give it again at once (R82 spec). An item of `engine` (R03 engine): for 24 hours it is not a candidate, and when the choice gave nothing, the second pass gives it too (R107 spec). Here there is only how much an attempt at such a task weighs, subsection 2.3;
- the build for the device (Т-10) — subsection 2.5 (R09, R10): the build and the same bits, and with them the types and the encoding (R08);
- how knowing an anchor skill is derived from the skills that teach it — subsection 2.6 (R50 spec; R07). It is last, so that the numbers 2.4 and 2.5, which the documents cite, keep. The second half of the item — what the engine does with "does not know" and with a cycle of prerequisites between modules through anchor skills — is `engine`'s (R04 engine; R107 spec);
- lessons (R58 spec), an item of `engine`: when to give a lesson and how to choose between lessons (R05 engine; R107 spec). It has no subsection here — a lesson gives the learner model only what section 1 says.

Items 2.4 and 2.5 concern both libraries, and 2.4 also `api`. Т-10 is a question of `spec`'s register and open, and its recommendation carries what is decided here. Item 2.4 holds the condition of R80 spec.

### 2.1. Level per module or per subject area

- **Gist:** what the learner's level is counted for — per module or per subject area —, what the learner's state consists of, where the learner starts, how far apart the bands are, what a task's guessing is and where the threshold "mastered" lives. An item of the card (architecture 3.8 spec; R23, R51 spec).
- **Decided:** R03. The level is per module, the correction is the skill's, and the scale of bands, the starting point, the trial series, guessing by the item's answer slots, the probability of a correct answer and the threshold "mastered" are section 3.1. The start a band's width below the grade's, an adult's correction and the guard after three failures in a row are the rules of R157 spec in place of calibration (R45) — section 3.1, "The rules of R157 spec". The card, concepts 11 and 15 spec and the 2026-11-02 amendment to R23 spec carry it.
- **Deferred:**
  - the initial estimate by age, by self-assessment and by a short check (concepts 11 spec). Age is kept only as a category (privacy 4 spec), and the bands do not name it; self-assessment — for example, a CEFR level — is needed by language modules, and there are none in the first slice (R41 spec); a short check is a choice of tasks, that is, a rule of `engine`. Return to it with the first language module or when the harness shows that the grade and the trial series are not enough;
  - a module's own step of bands — a field of the next format version (R44 spec), if the harness shows that one step for all modules is not enough;
  - the numbers that the harness `sim` will check (R163 standalone): the step of 2.5 between bands and the spread of the trial series σ₀ = 2.5 — standalone too calls them a guess and a number of simulation; I without guessing taken into account, because of which the uncertainty is understated on a choice of two; guessing 0 for `predicate`, where there are many correct answers;
  - a calibration of the start and the step by the grade. The start at the grade's band, which the start a band below replaced, was not fair by `sim`'s bound: under grades 1 and 5 of olympiad maths too few learners come to the estimate of grade 3, because of a tail that the trial series starts and the step of θ closes slowly (R58 sim; R41; measurements 5), and no set of the model's numbers of the start and the step holds the criterion (R42; measurements 7). Long runs of fitting wait until there are 100 children (R42); the rules of R157 spec stand in their place, `engine` does not change, and the author accepted the residual risk in dpia 7 spec (R157 spec).

### 2.2. The model of forgetting

- **Gist:** how the estimate of mastery fades without practice, on what it depends how fast, and by what sign `engine` learns that a skill is due for review. An item of the card (architecture 3.8 spec). The repetition rule itself is `engine`'s (concepts 12 spec).
- **Decided:** R04. Forgetting is the retention of a skill whose stability grows with spaced practice; the law and the constants are one for all modules; time does not remove "mastered" but gives a review date; `engine` gets the estimate for the time of the choice and the date; the current time does not enter the cache; after a pause of the module of 14 days or more, n before the step falls to 20, and the step of θ is wider — section 3.2. The cards architecture 3.8 and 3.9 spec and concepts 11, 12 and 15 spec carry it.
- **Deferred:**
  - a module's own constants of forgetting — a field of the next format version (R44 spec). Return to it with the first language module or when the harness shows that one set of constants for all modules is not enough: for words and grammar forgetting is a must (concepts 11 spec), and the sources of the constants are flashcards and words, not maths skills (sources 3.1–3.3);
  - the numbers that the harness will check: C = 0.5, τ₀ = 1 day, growth g = 15 and halving τ. In FSRS the growth of stability also slows down with the stability itself (sources 3.1), and here it does not;
  - the numbers of forgetting together with the forgetting of the module's level: with it they do not delay catching up after a break, but no set of them holds all the constraints (R32). A round of them — a base of dividing τ between 0.4 and 0.5, or a growth of τ that slows down with τ itself (R37) — waits for a candidate that brings a cause of its own: a fault that the gate or a measurement found, a need of the product, a row of `spec` (R59 sim; R40).

  All this is Д6 in table 4 of [run/DEBT.md](../run/DEBT.md).

### 2.3. The weights of evidence

- **Gist:** how much a piece of evidence weighs depending on the outcome of the attempt and on key provenance, how a series of items of one type is reduced to evidence and how a piece of evidence is divided between the task's skills. An item of the card (architecture 3.8 spec).
- **Decided:** R06. The step takes the first verdict of each item, while help — a rung of the ladder before it, a technique on request, a leading question (R21), a solution already seen — only keeps "mastered" from being set. The answers of one attempt weigh 1/(1 + (N − 1)/3) each, one piece of evidence is one step, and "mastered" is decided by the evidence. The item's primary skills get 1/J of the step each, the supporting ones nothing. The outcomes, key provenance, the facts of evidence, traps and the probe are section 3.3. The 2026-11-02 amendment to R82 spec, concepts 9.4 and 11 spec and the cards architecture 3.8 and 3.9 spec carry it.
- **Deferred:**
  - partial credit — the item's answers after the first verdict, by the number of hints. The log carries the number of incorrect answers before the first verdict, and a new version will read more by recomputation, if the harness shows a benefit;
  - a reduced weight of an answer with help and of a repeat of a task: now their weight is full, and a correct answer with a technique or from memory raises the estimate a little;
  - the numbers that the harness will check: ρ = 1/3, the failure boundary ½ and the rule τ for a page;
  - the weight of a key entered by an adult — with tasks from paper, stage 5 (О-11): they need a band and a difficulty on the module's scale and skills not from AI;
  - the weight of feedback by a rubric — with the answer to О-17.

  The numbers are Д7 in table 4 of [run/DEBT.md](../run/DEBT.md).

### 2.4. Evidence and recomputation from the log

- **Gist:** what a piece of evidence carries and how much it weighs, what the library takes and what it gives, and how the state is recomputed from the evidence log — on a change of the version of `learner-model`, of the release of `modules` or of `api`'s window of format versions (the fourth 2026-11-02 amendment to R80 spec), on the withdrawal of a module's version because of `wrong_key` and on a module's major version —, so that the state appended with a new piece of evidence matches the state recomputed from the log (R80 spec). An item of the card (architecture 3.8 spec). How `api` keeps it is its own (R06, R08 api), and the engine choice is repeated from the same log (architecture 1 spec), so the question is shared with `api` and `engine`.
- **Decided:** R08. The state is a fold of the log: appending is the same step as recomputation and gives the same bits on any build of one version of `learner-model`, and there is no snapshot. A piece of evidence carries observations — everything R82 spec names and everything the model reads —, and its kind only grows, and an earlier version reads a new kind, skipping what it does not know. The transfer by `retired_skills` gives each replacement the retired skill's share in full, skills merged in one item are one skill, and "mastered" is derived anew from the transferred evidence. The band's place is as recorded, and the library writes the start band when appending. All this is section 3.5. The 2026-11-02 amendment to R82 spec, architecture 1 spec and the band rule of format 0.3 (the 2026-11-02 amendment to R89 spec) carry it.
- **Deferred:**
  - the fields of stage 5 — the provenance of the band and the difficulty of a task from paper, which an adult or a catalogue gives (the 2026-11-02 addition to О-11), — as a new kind of evidence, when it comes to them;
  - partial credit: the log carries the item's verdict at the end and the number of incorrect answers before the first "correct", and a new version will read it by recomputation, if the harness shows a benefit (section 2.3; Д7 in table 4 of [run/DEBT.md](../run/DEBT.md));
  - the tolerance 1e-9 at the thresholds of sums of weights is right while an attempt has no more than 170 items with an answer, — Д9 in table 4 of [run/DEBT.md](../run/DEBT.md).

### 2.5. The build for the device

- **Gist:** what the library keeps to already now, so that computing on the device — for offline — stays possible, though computing on the server is recommended. An item of the card (architecture 3.8 spec). The requirement is the same for `engine`, so the question is shared with it. The promise of R08 — the same bits on any build of one version — is here too, and with it the library's types and the encoding of the evidence and the state.
- **Decided:** R09 and R10. The library gives the same bits on any build of one version of its own, WebAssembly too, and builds for the device as it is: the rules of code, only Go's standard library, its own exp and log and the checks of CI — section 3.6. The graph types, the description of a task, the operations and the encoding — protobuf's binary form by its own schema — are section 3.7. The 2026-11-02 addition to Т-10 and the card carry it.
- **Deferred:**
  - the glue with the app — `syscall/js` or `//go:wasmexport` — and the size of the WASM build, about 2 MB and about 500 KB after gzip (sources 4.7), — stage 5, with offline;
  - how offline works — computing on the device or a queue that the server chose in advance, — stage 5, by Т-10 (the 2026-11-02 addition to Т-10);
  - TinyGo and gccgo: they have no promise of the same bits (section 3.6);
  - the targets without a run — ppc64le, ppc64, s390x, riscv64 and loong64 with one differential check of the assembler, arm, mips and builds with software floating point —, Д10 in table 4 of [run/DEBT.md](../run/DEBT.md); the run at GOAMD64 v3 rests on the runner's processor — Д11 in the same table.

### 2.6. Knowing an anchor skill

- **Gist:** how the model derives whether the learner knows an anchor skill of the shared taxonomy from the skills that teach it, in the learner's modules, — and from skills with the assessability `rubric` or `none`, which have no verdict (О-17), — and whether such skills have a state. An item of the card (architecture 3.8 spec; R50 spec). The second half of the item — what the engine does with a cycle of prerequisites between modules through anchor skills — is `engine`'s: the format checks each module separately and does not see such a cycle (R50 spec).
- **Decided:** R07. The learner knows an anchor skill when "mastered" is set on at least one skill that teaches it, in any module whose graph came, and in any band; the threshold is the same as for "mastered", and forgetting does not remove knowing. "Does not know" is given only by the modules that the engine gets, and otherwise the anchor skill is "unknown". Each state has a reason and the skills that teach it. Skills with `rubric` and `none` have no state. The model computes all this at reading — section 3.4. Concepts 4 and 11 spec and the cards architecture 3.8 and 3.9 spec carry it.
- **Deferred:**
  - knowing with the band taken into account: a module where a skill is taught in several bands gives knowing already by the lowest of them, and the band needed is a field that the format does not have. In the first slice only school maths of grade 3 teaches anchor skills, and it has one band — Д8 in table 4 of [run/DEBT.md](../run/DEBT.md);
  - a short check instead of a lesson at "no evidence" — a rule of `engine` (concepts 4 spec);
  - knowing by skills with `rubric` — with the answer to О-17, when feedback by a rubric may become evidence.

### 2.7. An answer by a second way and a submission corrected at once (О-25)

- **Gist:** an answer slot with a number may be answered by a second way — a keyboard or on-screen digit buttons, without a path, which an adult turns on in the profile, checked without AI, as a choice of an option is —, which leaves to `learner-model` whether such an answer is evidence and with what weight, and how to count a submission that the learner corrected right away (product 13 spec, О-25). In format v0 an answer slot with a number is filled only by handwriting, and an answer by keyboard comes by a minor change of a later format version ([format text](../spec/format/README.md), section 14).
- **Decided** for the model by the 2026-11-07 addition to О-25, by row В23 (product 13 spec; R15), and the card architecture 3.8 spec says the same; О-25 itself stays open:
  - an answer by a second way is evidence with the weight of a handwritten answer, and a smaller weight is rejected;
  - a submission whose verdict the learner saw the model always counts;
  - only a submission cancelled before the verdict may go uncounted — it is an edit of the answer in `app`, and it does not reach evidence. Whether it needs a confirmation before the check and whether that removes the nonconformity of criterion 3.3.4 that R11 app declared, `app` decides (R90 spec).
- **What it means here:** neither fields of evidence nor code. An answer by a second way comes into the evidence as the same first verdict of the item, with the same answer slot with a number and guessing 0 (R03), and its weight is the weight of any answer (section 3.3). The evidence says nothing about the way of input: the model does not need it.

### 2.8. Format 0.4: the leading question and the check kind `time`

- **Gist:** two parts of format 0.4 concern the model (R108 and R111 spec; [format text](../spec/format/README.md), section 2):
  - a task's optional closed leading question, `leading_question`, — a question of the author that breaks the task down. `api` gives it on "Hint" before the first wrong answer, together with the task's techniques, and the evidence carries that the learner got it (the 2026-11-08 amendment to R82 spec);
  - the check kind `time`: the key is a time of day `H:MM`, hours from 0 to 23 and minutes from 00 to 59, and the answer is compared with it by hours and minutes, each as a number, without a 12-hour dial (R111 spec; [format text](../spec/format/README.md), section 15).

  The word flag of a choice by hand and the substitutions in the path of a template's picture do not concern the model: a choice by hand stays a choice of k options, and the model does not read pictures.
- **Decided:** `time` is a free answer, a value of `CheckKind` with guessing 0, by R03 (R19). The leading question is help, as a technique on request, by a field of its own of the evidence, a fact of the attempt; `api` writes it only from a release of `learner-model` that has the field (R21; the 2026-11-09 amendment to R82 spec). They are `CHECK_KIND_TIME = 5` and `leading_question_on_request = 20` of evidence kind 2 (R28) — sections 3.1, 3.3, 3.5 and 3.7.

## 3. How it is built

Here is what to build. Subsection 3.1 is the decision of section 2.1 (R03), 3.2 of section 2.2 (R04), 3.3 of section 2.3 (R06), 3.4 of section 2.6 (R07), 3.5 of section 2.4 (R08), 3.6 and 3.7 of section 2.5 (R09, R10). Below is the library as a whole, and the rules, numbers, examples and sources are in the subsections.

**What it is.** The Go module `github.com/altnovatrail/learner-model` is a library without input or output in `api`'s process: the root with the types and operations, the package `detmath` — exp, log, σ, power and retention —, the internal codec and the `.proto` schema of the evidence and the state (section 3.7). The code imports only Go's standard library and gives the same bits on any build of one version, WebAssembly too (section 3.6).

**Data:**
- a module's graph is in the library's types; `api` and `sim` build it from the server's part of the package: the bands in order, the skills with their assessability, `teaches`, the bands where they are taught, the prerequisites and the tier, and `retired_skills` (section 3.7);
- a piece of evidence is one per finished attempt, with the number of its kind: the observations of the attempt, not the model's conclusions — the description of the task, the start band, the outcome, the time of practice, a technique at a step and on request, whether the solution went to the device, for each item the submissions, the verdicts and the number of incorrect answers, the traps that the program found (section 3.5);
- the evidence log in `api`'s database is the source of truth, and the state is a cache, a fold of the log in the order of numbers (section 3.5);
- the state: for a module — the level θ, the sum of the weights of answers n, the start band, the answers of the trial series, whether the module's last piece of evidence was still in the series, and the traps that have no skills with a state; for a module's skill with the assessability `verdict` — the correction δ, the sum of weights m, the number of pieces of evidence, the time of the last practice, the stability τ, the traps with their count and the time of the last one, the mastery band and the failures in a row (sections 3.1–3.3). Anchor skills and skills without a verdict have no state of their own (section 3.4).

**Operations** — in the order in which `api` calls them (sections 3.5 and 3.7):
1. `BandByGrade` — the start band by the grade from the profile and the bands of the attempt's version; `api` puts it into the evidence;
2. `Append` — append a piece of evidence to the cache. The state and the evidence with the start band that the appending took, `api` writes in one transaction;
3. `Read` — for the time of the choice, with the grade and the engine's modules: what the engine and the skill map read, and `Readout.Task` — P of the items of a candidate task;
4. `Recompute` — a fold of the log from an empty state: on a change of the cache's mark and on a repeat of a choice.

**The step of one piece of evidence** is a function of the state before it, the evidence and the graph of its module, one and the same at appending and at recomputation (section 3.5):
1. the item's primary skills go by `retired_skills` of the current version's graph: each replacement gets the retired skill's share in full (section 3.5);
2. the item's answer is its first verdict, and help is a rung of the ladder before it, a technique on request and a solution already seen (section 3.3);
3. each of the N answers of the evidence weighs u = 3 / (N + 2), and for a skill by its own items (section 3.3);
4. P of an item is c + (1 − c) · R̄ · σ(θ + δ̄ − β): guessing by the answer slots, the retention of the primary skills, the band and difficulty of the task on the module's scale, and in the trial series without retention, R̄ = 1 (sections 3.1 and 3.2);
5. while n < 5, θ is the maximum of the weighted likelihood of the trial series; after that the Elo step with the retention weight w moves θ and δ, and δ of a primary skill by its share 1/J. After a pause of the module of 14 days or more since the last practice of its skills, n before the step falls to 20: the level is forgotten, and the step of θ is wider (sections 3.1–3.3);
6. n grows by N · u, while m, the number of pieces of evidence, the stability τ by the share of correct answers and the time of practice change for the primary skills of the items with an answer. A trap that the program found goes to the skills from its `skills` in the catalogue, and without them or when they have no state, to the module (sections 3.2 and 3.3);
7. a failure is less than half correct for a skill in a piece of evidence. Two failures in a row remove "mastered", and it is set by a piece of evidence without a failure where the skill has a correct answer without help and m ≥ 5, — in the highest band where the skill is taught, not above the task's band, for which σ of the cautious level against an average task is not below 0.71875. "Mastered" is set anew only in a band above the one where it holds, and it does not move down (sections 3.1 and 3.3).

**At reading, past the cache:** P of an item with retention for the time t₀, and in the trial series without it, at R̄ = 1; for a mastered skill, the cautious probability with retention and the review date or "already" (section 3.2); the starting point of modules without evidence (section 3.1); knowing the anchor skills — "knows", "does not know" or "unknown", with a reason and the skills that teach them (section 3.4).

**The version's constants:** the step of bands 2.5, σ₀ = 2.5, the Elo steps, I and 20 in the uncertainty, the threshold 0.71875, five answers of the series and before "mastered", two failures in a row and the grid step of the series 0.05 (section 3.1); C = 0.5 and F = 19/81, τ₀ = 1 day, g = 15, halving τ after "not yet", and for the forgetting of the level the module's pause d₀ = 14 days and N = 20 (section 3.2); ρ = 1/3 and the failure boundary ½ (section 3.3); the tolerance 1e-9 at the thresholds of sums of weights (section 3.5). The model needs no fields of the format, and the places of the bands are held by the comparison rule of format 0.3 `module-version.bands-kept` (the 2026-11-02 amendment to R89 spec; section 3.5). To change a constant is a new version of `learner-model`, and only through the harness (concepts 11 spec; R163 standalone).

**Links:**
- `api` calls the operations, keeps the log and the cache with a mark of the version of `learner-model`, the release of `modules` and the format window, and with a choice its time (section 3.5); the choice keeps the grade too (architecture 1 spec);
- `engine` reads the output of `Read`, takes the graph types and the description of a task (R28 spec; section 3.7) and may take `detmath` (section 3.6);
- `sim` imports the library together with `engine` and holds the gate of its CI (section 1);
- the skill map gets the model's data from `api` (section 1).

**How the tests and CI check what is built:** the appended equals the recomputed byte for byte on random logs, the transfer vectors, the boundary of the thresholds, reading the kinds and standalone's test vectors with the tolerance 1e-4 (section 3.5); the golden vectors of the state and of reading on all builds, the rules of code and the differential check of the assembler — by the shared jobs of `go.yml` (section 3.6); buf — the schema and the codec (section 3.7).

**What is deliberately deferred** — "Deferred" in the subsections of section 2 and the debts Д5–Д11 in table 4 of [run/DEBT.md](../run/DEBT.md).

### 3.1. The learner's state, the scale and "mastered"

Decided by R03. The starting model is standalone's: a Rasch model with a guessing floor and an Elo update, where the general level is corrected for the topic (SPEC 2 standalone; R79, R187 standalone; sources — sources 1.1, 1.5–1.7 and 2.1–2.4). Here it has a level per module, a correction per skill and guessing by the item's answer slots. The olympiad maths module, where a task is a choice of five options, computes exactly as standalone does while the skills' retention is 1, and standalone's test vectors fit it at Δt = 0 (SPEC 2.8 standalone; forgetting — section 3.2).

**The level is per module.** Each module of the learner has its own level θ, on the scale of that module's bands. The evidence of different modules does not add up into a level: modules are linked only by an anchor skill — by a prerequisite and by the field `teaches` (section 2.6; concepts 4, 11 spec). The model has no subject area.

**An item's answer.** The formulas below are for an item's answer with a verdict, where S = 1 at "correct" and S = 0 at "not yet" ([format text](../spec/format/README.md), section 15). Which answer of an item goes into the step, how much an answer of an attempt of many items weighs and what "unreadable" and "unchecked" give, R06 decided — section 3.3: the step takes the item's first verdict, and one piece of evidence gives one step.

**The learner's state in a module:**
- the level θ;
- n — the number of answers in the module, and after a pause of the module of 14 days or more no more than 20: so the level is forgotten (section 3.2);
- the place of the start band;
- while the trial series has not ended, its answers: the place of the band and the difficulty of the task, the check kind and the number of options of each answer slot, the outcome;
- the count of the guard, and whether an adult's correction ended the trial series (R45);
- whether the last piece of evidence that moved the module was still in the series, — the one that ended it too: it clears the answers of the series, and by n after it this cannot be told (R104 spec; R07).

**A skill's state** — for a module's skill with the assessability `verdict`:
- the correction δ;
- m — the number of answers of items where the skill is primary;
- the number of pieces of evidence — attempts where it, as a primary skill of an item, has an answer (concepts 11 spec; section 3.3);
- the time of the last practice;
- frequent traps — for each trap the number of pieces of evidence where the program found it, and the time of the last one (section 3.3);
- "mastered" — the mastery band or "no";
- failures in a row — pieces of evidence in a row where the skill has less than half correct (section 3.3).

The estimate of mastery of a skill is θ + δ, its uncertainty is v (below). A skill without answers has δ = 0, m = 0 and no "mastered": the model computes this at reading and does not write it into the cache (R08 api; the 2026-11-02 amendment to R80 spec). Anchor skills and skills with `rubric` and `none` are section 3.4, and forgetting — the stability τ and retention — is section 3.2.

**The scale.** The band at place p in the manifest's `bands`, counting from 1, stands at the offset 2.5 · (p − 1). A task of band p with difficulty d is at β = 2.5 · (p − 1) + (d − 3). The two hardest difficulties of a band overlap the two easiest of the next one, and the points do not coincide: difficulty 5 of the first band, β = 2, stands between difficulties 2 and 3 of the second, β = 1.5 and 2.5 (SPEC 2.1 standalone). The band's place is taken in the module version the attempt was on, and a major version does not change the places of earlier bands (section 3.5).

**The starting point** is θ₀ = 2.5 · (p − 1) − 2.5, a band's width below the offset of the start band at place p, the band of the grade (R157 spec; R45). A learner of the lowest band's grades starts at −2.5, below the lowest band, a place on the scale that no band has. The start band is:
- the start band is the one whose grades name the learner's grade ([format text](../spec/format/README.md), section 11);
- a grade that no band names goes to the nearest band, by the distance to its nearest grade, and at an equal distance to the lower one. So standalone's `Start(grade)` leads a grade outside 1–6 to the nearest band;
- without a grade and in a module without grades, the first band: the trial series of five answers moves the learner by a band in either direction.

The start band is fixed by the module's first piece of evidence, each next one carries the same, and the library writes it when appending (section 3.5), and a change of grade does not move it afterwards — as in standalone, where the start "is written into `ratings.start` when the profile is created" and a grade changed later "moves neither θ, nor θ₀" (SPEC 2.1 standalone; О-56 standalone). For a module without evidence the start band is computed at reading, by the grade from the profile (privacy 4 spec; R08 api). Age, self-assessment and a short check do not enter the starting point yet (section 2.1).

**The probability of a correct answer of an item:** P = c + (1 − c) · σ(θ + δ̄ − β), σ(x) = 1 / (1 + e^(−x)).
- δ̄ is the mean δ of the item's primary skills: the item's own skills or, if it has none, the task's skills (concepts 6 spec; [format text](../spec/format/README.md), section 14). Supporting skills do not enter P. With one primary skill this is standalone's formula.
- c is the item's guessing, the product over its answer slots:
  - for `choice`, 1/k, where k is the number of options of the binding `choice` or of the check itself;
  - for `number`, `fraction` and `predicate`, 0;
  - for `time` of format 0.4, also 0: it is a free answer, the value `CheckKindTime` (R108 spec; R19, R28; subsection 2.8).

  An item is correct only when each of its answer slots is correct ([format text](../spec/format/README.md), section 15), so to guess an item is to guess all its answer slots. A workbook box "56 [ ] 65" with the signs <, > and = gives 1/3, two such boxes 1/9, a choice next to a number 0. An item with a rubric has no verdict, and so no P either.
- With forgetting, σ is multiplied by the retention R̄ of the item's primary skills: P = c + (1 − c) · R̄ · σ(θ + δ̄ − β) (section 3.2). The formulas of this section are at R̄ = 1.
- Example: a learner of grade 3 in the olympiad maths module with bands 1–2, 3–4 and 5–6 has the start band "3–4" and starts a band below it, with θ₀ = 0. On a task of band "1–2" of difficulty 3 with five options β = 0, and P = 0.2 + 0.8 · σ(0) = 0.6, as SPEC 2.1 standalone has it for a task of band "3–4" at its start of 2.5. A learner of grade 7 has the band "5–6" and starts with θ₀ = 2.5, and a learner without a grade in an essay module has its first band and starts below it, with θ₀ = −2.5.

**The step after an answer** — as in R187 standalone:
- θ += K_θ · (S − P), K_θ = max(0.2 / (1 + 0.05 · n), 0.05) — the step decreases with the number of answers, but not below 0.05;
- δ += K_δ · (S − P) for each primary skill of the item, K_δ = 0.4 / (1 + 0.05 · m);
- n and m are the answers before this one (SPEC 2.2 standalone), and n after a pause of the module is already after the forgetting of the level (section 3.2).

This is the step of one answer. A piece of evidence of many items gives one step — the sum of such steps, with the weight of the series and the division between skills, and a probe has no raise (section 3.3). β is not updated: the author sets the task's hardness.

**The trial series** is the first five answers in the module. For a task of one item these are five tasks, as in standalone, and in general whole pieces of evidence of the module while n < 5, with the weights of the answers (section 3.3). After each piece of evidence of the series, θ is the point of the maximum of L(θ) = −(θ − θ₀)² / (2σ₀²) + Σᵢ [Sᵢ · ln Pᵢ + (1 − Sᵢ) · ln(1 − Pᵢ)], where σ₀ = 2.5 and Pᵢ is the probability of answer i with its guessing. The maximum is searched for as in standalone: a grid with the step 0.05 on θ₀ ± n · σ₀², then halving the best cell by the sign of the derivative (SPEC 2.2.1 standalone). A grid, not Newton's method, because L can have two peaks. In the series δ do not move, while the counters, the traps and "mastered" are written as usual. After the series θ and δ go by the step above.

**The uncertainty** is v = 1 / (1/σ₀² + I · n) + 1 / (I · (20 + m)), where n and m are the answers with this answer, not before it, as in the step; 20 is the uncertainty of a new correction, as from twenty answers; I = 0.149981 is how much one answer in the middle of the corridor tells about the level at guessing 0.2 (SPEC 2.5 standalone). I does not depend on guessing: on a choice of two an answer tells less, and the uncertainty there is understated. The harness will check this (section 2.1).

**"Mastered":**
- is set after a correct answer without help, in a skill where m ≥ 5 with this answer, by θ and δ after it. What help is and how the rule is read on a piece of evidence of many items is section 3.3;
- is set in the mastery band — the highest of the bands where the skill is taught, not above the task's band, for which σ(λ − the band's offset) ≥ 0.71875 at the cautious level λ = θ + δ − √v. The average task of a band is difficulty 3, its β is the band's offset;
- 0.71875 is P = 0.775 for a task with five options, the middle of standalone's corridor, and in logits a margin of 0.9383 = ln(0.71875 / 0.28125) above the average task (SPEC 2.5 standalone). The threshold does not depend on guessing: "mastered" means "solves", not "guesses";
- the bands where the skill is taught are the bands of tasks where it is a primary skill of the task or of an item (section 2.5);
- "mastered" is set anew only in a band above the one where it holds, and it does not move down. It is removed by two failures in a row — pieces of evidence where the skill has less than half correct (section 3.3). For a task of one item a failure is the answer "not yet", and this is two answers "not yet" in a row, as in R187 standalone;
- in standalone a mastered topic returned into the rotation when its recommended point went above the mastery band (SPEC 2.5 standalone). The recommended point is the corridor, and the corridor is a rule of `engine` (R23 spec): whether to return a skill `engine` decides, and R03 engine decided — it returns it, with the code "in rotation" (concepts 12 spec; R107 spec). The model gives "mastered" as it is, and so knowing an anchor skill (section 3.4) and the skill map read it;
- the threshold is a constant of `learner-model`. It equals the target of `engine`'s corridor at guessing 0.2. If `engine` changes the target, the model's new constant passes `sim` together with it (the 2026-11-02 amendment to R23 spec; R103 spec: the pair goes by one tag of `sim`).

**The rules of R157 spec** (R45) in place of a calibration, beside the start above:
- **the band the level stands in** — of the bands of the module's current version, the one whose middle, at 2.5 · (p − 1), is nearest θ, within 1.25, the upper of two as near; the highest band when θ is above it; none when θ is below −1.25, as at the start below the lowest band, and none without a current version or bands. `Read` gives it for each module, of the starting point for one with no evidence: the skill map shows it, and `api` holds an adult's correction to a band next to it (R157 spec; the 2026-11-16 amendment to R99 spec);
- **an adult's correction** — evidence of kind 3, section 3.5: the band the adult set, "Easier" or "Harder" in `app`. θ becomes the band's offset, its middle task, unless it stands in that band already, and a band above the current version's highest is the highest. Either way it ends the trial series, if it lasts, without its five answers, and the series' answers go, and starts the guard's count anew, so a repeat changes nothing; n stays, no skill moves, and its time is no practice. The library does not hold the band next to the current one: a log that `wrong_key` cut (section 3.5) may put a correction after another level than it was made at;
- **the guard** — past the trial series, a piece of evidence that moves the module is a failure of the module when the sum of R̄ · σ of its answers, their P above the guessing floor, is at least half their number — the model gave the task as fitting, and not by guessing alone — and fewer than half of its items with an answer were solved at the end of the attempt: the final verdict right, or the first one with no final one recorded. The third failure in a row lowers θ by 2.5 after the step, not below −2.5 nor upward from below it; the count goes on, lowering no more, until a piece of evidence that is no failure starts it anew, and a correction starts it anew too. A wrong first answer the learner put right on the ladder is no failure, nor is a task the model already took for too hard: the guard is for the tasks the learner is stuck on while the level stands too high (R45, "Why").

**What `engine` reads:** the probability P of an item for the corridor (R23 spec), the estimate θ + δ and its uncertainty, "mastered" with the mastery band. The review date and retention for the time of the choice are section 3.2. For a module, two more signs of the trial series: whether it has ended and whether the module's last piece of evidence was still in it. In the series the engine does not work through an error and gives no probes, and the second sign is needed so that an error in the series' last answer is not worked through (R104 spec). The model computes the series, and the engine does not repeat its rules.

**Numbers** — the constants of the version of `learner-model`:
- the step of bands 2.5;
- σ₀ = 2.5;
- the initial steps 0.2 for θ and 0.4 for δ, their decrease 0.05 and the floor of the θ step 0.05;
- I and 20 in v;
- the threshold 0.71875;
- five answers of the series, five answers before "mastered" and two failures in a row;
- the grid step of the series 0.05;
- the start a band below, the half band 1.25 of the band the level stands in, the guard's three failures and its floor −2.5 (R45).

The constants of forgetting are section 3.2.

To change them is a new version, and only through the harness (concepts 11 spec; R163 standalone).

**In code** (R03):
- the operations are `BandByGrade` in `band.go`, `Append` and `Recompute` in `fold.go`. An empty state is the zero `State{}` or a state of kind 1 without modules: `Append` accepts both. `Recompute` is `Append` over the log starting from the second of them — it is what an empty log gives —, and an error names the index of the piece of evidence in the log, from 0. `Recompute` checks a module's graph once per module, and `BandByGrade` rejects an empty list of bands, bands not at places 1, 2, … in order and a band ID twice or not in UTF-8. `NewGraph` rejects bands not by places and an ID twice, but accepts a graph without bands: a module without tasks may have no `bands` ([format text](../spec/format/README.md), section 11), and `Read` gives such an engine's module without a state without a starting point, not as an error (section 3.4; R07);
- the step is `step.go`: one step per piece of evidence, the step of section 3.3 — u = 3 / (N + 2), the share 1/J, n and m as sums of weights with the tolerance ε, a failure is less than half correct. The trial series is `trial.go`, the constants and formulas of this section are `model.go`;
- `Append` rejects beyond the codec's checks:
  - a task ID that does not name a module;
  - a difficulty outside 1–5, a place of the task's band or of the start band below 1;
  - a choice without a number of options or with fewer than two, an item with an answer without answer slots;
  - a module without an entry in the input "graphs", a graph of another module or not in the order of `NewGraph`, and a graph where a skill is taught in a band that is not in it;
  - a state not in the canonical order, a module in the trial series whose last piece of evidence was not in it, and a trial answer that the step does not write: without answer slots, with an answer slot of an unknown kind, without the verdict "correct" or "not yet", with a weight not above 0 or after the end of the series.

  A piece of evidence of a kind newer than its own it accepts;
- an item's answer is the first verdict "correct" or "not yet" of an item all of whose answer slots are of a known kind. A piece of evidence without a single answer changes nothing: it creates no module and fixes no band. A module without a current version moves nothing. The start band is from the state, and for a module without one from the evidence;
- retention, the weight w and stability are section 3.2, help, outcomes, provenance and traps section 3.3, and the transfer by `retired_skills` section 3.5;
- the checks:
  - standalone's test vectors, SPEC 2.6, 2.7, 2.9 and 2.12, with the tolerance 1e-4. 2.6 goes from n = 0 by the Elo step, without the trial series, and a test inside the package checks it;
  - the property "the appended equals the recomputed" on 50 random logs;
  - the boundary of the thresholds on pages of 8 and 13 items;
  - the golden vectors of the fold `testdata/fold/` — a piece of evidence and the state after it, in the codec's bytes.

### 3.2. Forgetting and repetition

Decided by R04. Forgetting is the retention of a skill whose stability grows with spaced practice. So spaced repetition is built — SM-2, the Leitner system, HLR, FSRS (sources — sources 3.1–3.3), — and here it is reduced to four constants and lies on top of the model of section 3.1. The module's level θ is forgotten only through n: after a pause of the module the step of θ is wider — below, "Forgetting of the module's level" —, and time does not remove "mastered".

**The time of practice** of a piece of evidence, t, is the time of the attempt's last action, which `api` keeps, not the time of its closing: an abandoned attempt is closed by a job after 7 days (R10 api). The evidence carries t in whole seconds of UTC (R08; section 3.5). Times are instants of UTC, and Δt is their difference in days of 86,400 s. The model knows neither calendar days nor the time zone, and it does not need the family's zone: forgetting adds no new data of the learner (privacy 4 spec).

**A skill's state** (section 3.1) adds the stability τ, in days. The time of the last practice t_last is already in it. τ and t_last are set by the times of practice in the log, and they go into the cache (R80 spec; R08 api).

**Retention** is R(Δt, τ) = (1 + F · Δt / τ)^(−C), C = 0.5, F = 0.9^(−1/C) − 1 = 19/81, where Δt is the days since the skill's last practice:
- R(0) = 1, and R(τ) = 0.9: τ is the pause after which retention is 0.9;
- the law is a power law: retention first falls fast, then slowly;
- a skill without practice is not forgotten: it has no τ, and R = 1.

Stability here is τ: in FSRS it is S, and S in section 3.1 is the outcome of an answer.

**The probability of a correct answer of an item:** P = c + (1 − c) · R̄ · σ(θ + δ̄ − β), where R̄ is the mean R of the item's primary skills, as δ̄ is (section 3.1). The learner solves an item if they remember its skills and the task is within their power; guessing is as before. At R̄ = 1 this is the formula of section 3.1.

**A piece of evidence with the time of practice t.** For each skill that it moves (section 3.3):
1. Δt = max(0, t − t_last) and R = R(Δt, τ) — one per skill for a piece of evidence, however many items carry it, so the order of a page's items changes nothing;
2. the step of section 3.1 with the weight w: θ += K_θ · w · (S − P) and δ += K_δ · w · (S − P), where P is with this R̄, P₁ is the same P at R̄ = 1, and w = R̄ · P₁ · (1 − P₁) / (P · (1 − P)). For a piece of evidence of many items the step is the sum of such steps over its answers, with the weight of the series (section 3.3);
3. stability — no more often than once per skill for a piece of evidence:
   - the skill's first answer — τ = τ₀ = 1 day;
   - a correct answer — τ := τ · (1 + g · (1 − R)), g = 15;
   - "not yet" — τ := max(τ₀, τ / 2).

   When not all of the skill's answers in a piece of evidence are correct, τ changes by the share of correct answers — section 3.3;
4. t_last := max(t_last, t): a time that came backwards does not turn the clock back.

In the trial series P is without retention, R̄ = 1 and w = 1: its answers lie in the module's state without time (section 3.1), and the skills that it finds are almost all new. τ and t_last in the series are written as usual — with R by the time of practice: a piece of evidence always has a time of practice.

w is how many times less an answer tells about the skill than without forgetting: the step of R03 is (S − P), and w · (S − P) is what it turns into if retention is taken into account in the derivative of the log-likelihood of the answer by θ + δ̄. w ≤ 1, and at R̄ = 1 w = 1 — and the step of section 3.1. An incorrect answer after a long pause moves θ and δ little — forgetting explains it —, and halves τ, and the skill soon comes for review again. A correct one moves θ and δ almost as without forgetting, and grows τ strongly. Without w a correct answer after a pause would raise θ and δ many times more strongly than without it, though the skill was not in doubt.

Within a session R ≈ 1 — two minutes after an answer at τ = 1 day R = 0.99984, — and τ hardly grows: practice in a row does not hold for long (Rohrer & Taylor 2006, sources 3.7). τ grows from a correct answer after a pause: at R = 0.9 by 2.5 times, at R = 0.6 by 7 times.

**Forgetting of the module's level** (R04). The module's pause is from the greatest t_last of its skills to the time of practice of the evidence. If it is no less than d₀ = 14 days, before the step n := min(n, N), N = 20:
- K_θ = max(0.2 / (1 + 0.05 · n), 0.05) after a pause is 0.1, if n was greater than 20, not the floor 0.05, at which the step stands at n > 60. Answers after a break move θ more strongly, and the model catches up with the fall of the level faster: without this the step with retention explains it for a long time by the forgetting of skills, since an answer after a pause with a low R̄ tells little about the level (R31);
- with n, v is wider too for each skill of the module, not only for the one that came back, and the cautious level λ and q are lower: "mastered" comes later, and the review date is nearer. `engine` reads this too — the estimate, the uncertainty and "mastered" through `Readout`;
- the rule stands in the step after the check "no answers". A piece of evidence without an answer and a piece of evidence that weighs nothing change nothing after a pause either (section 3.5). At reading there is no rule: `Read` after a pause gives v by the former n until a piece of evidence comes;
- there is no new field: the module's time is derived from the skills' t_last, and the rule lasts through n, which is already in the state. A module without skills knows no pause. A piece of evidence whose answers have no primary skills in the graph, J = 0, moves θ and narrows n, but does not set t_last (section 3.5), so the pause of the next piece of evidence is still counted from the skills: while only such evidence comes into the module, each narrows n again, and it stays at 20;
- in the trial series n < 5 ≤ N, and the rule changes nothing. N is no less than 5, otherwise the module would return into the series.

d₀ = 14 days is the threshold of a skill's pause in `sim`'s measures. N = 20 stands between N = 10, which breaks "no worse" by `M5.caught` on `F-none` of school maths, and N = 40, which hardly moves the catching up (R32).

**"Mastered"** is set and removed as in section 3.1: the check comes after the answer, when t_last = t and R = 1. Time does not remove "mastered".

**The review date** — for a mastered skill, in the mastery band p:
- q = σ(λ − 2.5 · (p − 1)) is the cautious probability of solving the average task of the mastery band, without forgetting; λ = θ + δ − √v (section 3.1);
- Δt* is the pause after which R(Δt*, τ) · q = 0.71875: R* = 0.71875 / q and Δt* = τ · (R*^(−1/C) − 1) / F;
- if q < 0.71875, the threshold does not hold even without forgetting — "already", and Δt* = 0. This happens when θ went down through other skills of the module;
- the date is t_last + max(Δt*, τ₀): no earlier than a day after the last practice. A repetition on the same day hardly grows τ, and a skill mastered just barely would otherwise come for review every few minutes — at q = 0.72 and τ = 1 Δt* = 21 minutes.

A power law has no floor, and there is always a date. The date depends on the module's θ and v, and any piece of evidence of the module moves them, so the model computes it at reading — from the cache and without the current time. The date is the nearer, the smaller the margin of q above the threshold.

**At reading for the time t₀** — it comes from outside, and for the engine choice it is the time that the choice keeps (the 2026-11-02 amendment to R06 api; architecture 1 spec), — the model gives:
- P of an item with R̄ at t₀ — for the corridor. This is the P that the step would take for a piece of evidence at t₀, so in the trial series it is without retention, R̄ = 1, as in the step (R07);
- for a skill with practice — θ + δ, v and R at t₀;
- for a mastered skill also "mastered" with the mastery band — as it is, the cautious probability with retention R · q at t₀ and the review date.

This is not written into the cache: «Время тоже не перенос» (R08 api).

**What `engine` reads** — all this. "Time to review" is t₀ no earlier than the date. When to give a skill with a date, how many repetitions a day and how they share the queue with new skills, consolidation and probes, `engine` decides, and so say the 2026-11-02 amendment to R23 spec and concepts 12 spec; R06 engine decided it (architecture 3.8 spec; R107 spec).

**Example.** A learner in the olympiad maths module mastered a skill in band "3–4", with the offset 2.5. θ + δ = 4.5, n = 100 and m = 20, so v = 0.2327, λ = 4.018 and q = σ(1.518) = 0.820. Then R* = 0.876 and Δt* = 1.288 · τ:
- at τ = 1 the date is in 1.3 days. The learner came back after a day and answered correctly: R(1, 1) = 0.9, τ = 1 · (1 + 15 · 0.1) = 2.5, and the next date is in 3.2 days, if θ, δ and v hardly moved;
- at τ = 10 the date is in 12.9 days. The learner came back after a summer break of 60 days: R(60, 10) = 0.645. A task of band "3–4" of difficulty 3 with five options, β = 2.5, would give without forgetting P = 0.2 + 0.8 · σ(2) = 0.905 — above the corridor, — and with forgetting P = 0.2 + 0.8 · 0.645 · 0.881 = 0.654, and the engine will give an easier task. Here w = 0.246: a correct answer will move θ and δ by 0.085 · K instead of 0.095 · K without forgetting and raise τ to 10 · (1 + 15 · 0.355) = 63.3 days, and "not yet" will move them by −0.161 · K instead of −0.905 · K and lower τ to 5. If the module had no practice at all for 60 days, the module's pause is 60 days, and before the step n := min(100, 20) = 20: K_θ = 0.1 instead of the floor 0.05, so a correct answer will move θ by 0.0085, not by 0.0043, and δ, as without the rule, by 0.017. After it n = m = 21, v = 0.465 instead of 0.228 and q = 0.793 instead of 0.824, and the next date is in 58.7 days after the answer instead of 84.9.

**How the harness checks this.** On the two maths modules of the first slice: the virtual learners have days, sessions and breaks, and their forgetting goes by a family of laws different from the model's law, — with learners who do not forget, and with a fall of the level after a break. The measures are the calibration of P by the length of the pause, in how many answers the model catches up with the real fall, the share of repetitions over time and a false "mastered" after a break. The populations and the criterion are `sim`'s: its card names the populations of virtual learners with a calendar and forgetting and the model's criterion with measures of forgetting, and which of the measures the criterion takes is recorded in `sim` (architecture 3.9 spec; R03–R05 sim; R11, R15). standalone's harness did not know time: "every tenth answer starts a new day", and the learners G3-drop lose their level "as after a long break" by the number of the answer (sources 1.3).

**Numbers** — the constants of the version of `learner-model`: C = 0.5 and F = 19/81, τ₀ = 1 day — both the first stability and the smallest pause before the date, — g = 15 and halving τ after "not yet"; for the forgetting of the level, the module's pause d₀ = 14 days and N = 20. They are one for all modules, and a module's own ones are a field of the next format version, if the harness or the first language module shows the need (R44 spec; Д6 in table 4 of [run/DEBT.md](../run/DEBT.md)). To change them is a new version, and only through the harness (concepts 11 spec; R163 standalone).

**In code** (R04):
- in the step, R of each primary skill is one per piece of evidence, from its last practice to the time of practice, and for a skill without practice R = 1. R̄ of an item is the mean over its J skills, and at J = 0 R̄ = 1. P = c + (1 − c) · R̄ · σ, and θ and δ are moved by w · (S − P). In the trial series R̄ = 1 and w = 1;
- the weight is R̄ · P₁ · (1 − σ) / (P · (1 − R̄σ)), with (1 − c) cancelled, and without guessing (1 − σ) / (1 − R̄σ): so there is no division of a small number by a small number. At R̄ = 1, w = 1 and P = P₁ exactly, by a separate branch, and a piece of evidence without a pause has the bits as without forgetting;
- stability — the skill's first answer gives τ = τ₀, and after that τ changes once per skill for a piece of evidence by the rule of section 3.3 with the share of correct answers S̄, in the trial series too. At S̄ = 1 and 0 the exact branches of `detmath.Pow` give exactly the rules of this section. τ is kept below the greatest `float64`;
- Δt is from whole seconds: the difference is taken without overflow and no greater than 2⁵³ s, and then divided by 86,400;
- `Append` rejects a state with a skill whose τ is less than τ₀: a skill has a state only when it has had practice;
- the forgetting of the level is `forgetLevel` in `step.go`, and the step calls it after the check "no answers". The constants `levelPause` = 14 days and `levelWeight` = 20 stand in `model.go`, among the constants of forgetting. The pause is `pauseDays` from the greatest `LastPractice` of the module's skills; exactly 14 days already turns the rule on, and a time that came backwards does not turn it on;
- at reading, for `Read`, the functions of `model.go`: a skill's retention at t₀, q of a mastered skill and the review date. The date is computed without a power: R* = 0.71875 / q and Δt* = τ · (1/R*² − 1) / F. It turns into seconds by rounding up, with saturation at the greatest `int64`, and retention itself checks this second: while R · q there holds the threshold, the date moves by a second, no more than four times and only while it is less than 2⁵³ seconds; the date itself is bounded only by the saturation at `int64`. At q < 0.71875 — "already";
- the checks:
  - the example of this section and the page example of section 3.3 — as vectors, and the example with a break of the module, where "without the rule" is the same answer of a learner for whom another skill of the module had practice the day before;
  - the edges — the pause, the weight, τ and the date;
  - the edges of the forgetting of the level through `Append`: a second before 14 days and exactly 14, n no greater than 20, the pause from the latest skill, time backwards, a piece of evidence without an answer and a piece of evidence that weighs nothing, J = 0 and a module without skills. The state after the rule is checked against the state of a learner who has n = 20 before the step, to the byte;
  - standalone's vectors — without a pause, as section 3.1 requires;
  - the golden vectors of the fold — with a break of 60 days, after which "mastered" is set again;
  - the golden log `break` of the fold and of reading: more than 20 answers, a break of 30 days, and n after it is 20 and four answers. The other golden logs and standalone's vectors have no pause of a module of 14 days or more at n > 20.

### 3.3. Weights of evidence

Decided by R06. A piece of evidence is one finished attempt (R80 spec), and the model reads the items' answers in it: which answer goes into the step, how much it weighs and what it does to "mastered" and to the stability. A task of one item, solved without help, counts exactly as in sections 3.1 and 3.2, and at R̄ = 1 as in standalone.

**An item's answer is its first verdict.** The step takes the item's first verdict in the attempt — "correct", S = 1, or "not yet", S = 0. Later answers of the same item, after rungs of the ladder, do not go into the step, and the hint number k does not change the step. The guessing floor c describes the first try, and after a wrong answer and a hint the probability is no longer the one in P. So in standalone too an answer is recorded once (R101 standalone), and KDD Cup 2010 measured models by "Correct First Attempt" (sources 2.5).
- "Unreadable" and "unchecked" give no answer: they are not an error (R04 spec), and the item's first verdict will be the next one — "correct" or "not yet".
- When the solution is shown, an item that was empty or not submitted before the showing is "not solved", S = 0: О-7 names the showing so (R82 spec), and in standalone "I don't know" is a wrong answer (R93 standalone). An item that had only the verdict "unreadable" before the showing gives nothing.
- The input method does not change the weight: an answer in an answer slot with a number by keyboard or by digit buttons, when the format introduces them, weighs as an answer in handwriting. A submission whose verdict the learner saw the model always counts: it cannot be cancelled, and if it is the first with a verdict, its verdict is the item's first verdict. One cancelled before the verdict does not reach the evidence (О-25; section 2.7).

**Help.** An item's answer is with help if the attempt had any of these:
- a rung of the ladder before its first verdict — the attempt already had a wrong answer (R82 spec);
- the task's technique on request — anywhere in the attempt, before the item's first verdict or after: it is a fact of the attempt, not of the item;
- the task's leading question on request — from format 0.4, by the field `leading_question_on_request` of evidence of kind 2: "hint" gives it only before the first wrong answer, but, like the technique, it is a fact of the attempt, not of the item (R21, R28);
- the solution of this task — not of the template — has already gone to their device earlier (R82 spec), and the learner could have remembered the key. The solution goes there when the attempt is closed by a correct answer or by the solution shown (R82 spec), so this is almost any repeat of a solved task. A repeat after an attempt that the learner left or that was closed "unchecked" does not count as help: the learner did not see the key;
- the sketch "how to start", when it appears in the format (R57 spec).

An answer with help goes into the step, into the counters and into the stability as it is, but does not set "mastered" — like a hint in standalone: "The hint would make the answer measure the hint" (SPEC 2.5 standalone; R187 standalone). A technique at a step does not count as help: a new profile goes "step by step" by default (concepts 10 spec), and otherwise a learner of this mode would master nothing. A request after a wrong answer gives the words of the current rung, and "explain differently" the same rung in other words (R82 spec): the rung is already there, and they add nothing to the help.

**Series.** The answers of one attempt depend on each other: ten fractions in five minutes are not ten independent pieces of evidence (concepts 11 spec). A model that counts such answers as independent overstates the precision of the estimate and stops checking too early — so it goes with items of one stem, testlets (Bradlow, Wainer & Wang 1999, sources 2.6). So for N answers of an attempt each weighs u = 1 / (1 + (N − 1) · ρ), ρ = 1/3, that is u = 3 / (N + 2), and the whole attempt N · u answers. This is the effective sample size of N observations of one cluster at intraclass correlation ρ (Kalton, Brick & Lê 2005, sources 2.7):

| Items with an answer, N | 1 | 2 | 4 | 10 | 20 | without limit |
|---|---|---|---|---|---|---|
| The attempt weighs, N · u | 1 | 1.5 | 2 | 2.5 | 2.73 | 3 |

- For θ and n, N is all items of the evidence with an answer. For δ and m of skill s, it is N_s, the items with an answer where s is the primary skill, and the weight u_s is by N_s. After a major version of the module, s is the skill after the transfer by `retired_skills`, and its share in the item is the share of the recorded skill it replaced (R08; section 3.5).
- A series is the items of one attempt, whether of one type or not, and the format needs no mark of being of one type. Items with different skills get the discount on θ together, and on δ each skill by its own items.
- For a task of one item, u = 1.

**One step per piece of evidence.** All P, w and K are computed from the state before the evidence, and the shifts add up:
- Δθ = K_θ · u · Σᵢ wᵢ · (Sᵢ − Pᵢ) over all answers of the evidence, and n += N · u;
- Δδ_s = K_δ,s · u_s · Σᵢ wᵢ · (Sᵢ − Pᵢ) / Jᵢ over the answers where s is primary, and m_s += N_s · u_s. Jᵢ is the number of primary skills of item i;
- K_θ and K_δ are by n and m before the evidence, as in section 3.1, and the uncertainty v by them after;
- the order of the items changes nothing, as with R in section 3.2. The sums are computed in a canonical order — items by place, skills by ID — so that what is appended matches what is recomputed bit for bit (sections 2.5 and 3.5).

The 1/J of the step is because P holds the mean δ̄ of the item's primary skills (section 3.1): the derivative of the log-likelihood of an answer by the δ of one of J skills is 1/J of the derivative by δ̄. And m grows by the whole u: m is the number of answers where the skill is primary (section 3.1), and a skill that always stands in a pair would otherwise need twice as many tasks to "mastered".

**The trial series** is whole pieces of evidence of the module while n < 5: a piece of evidence that found n < 5 enters the series whole. The maximum is sought of the weighted likelihood — L(θ) from section 3.1 with the sum Σᵢ uᵢ · [Sᵢ · ln Pᵢ + (1 − Sᵢ) · ln(1 − Pᵢ)], — and for each answer of the series the module's state keeps its u. A page of ten items weighs 2.5 answers, so the series of a workbook module is two pages, and of an olympiad maths module five tasks, as in standalone.

**"Mastered" by evidence.** A piece of evidence fails a skill if fewer than half of the skill's answers in it are correct. `engine` reads this rule too — by the facts of the summary of issues, with its own code: the rule "failure" is one for the model and for the engine (R107 spec). A new boundary, a candidate through the gate, will need an edit of `engine` too; how to release them together R103 spec names only for the pair of the threshold and the target of the corridor, and it is decided with the first such candidate — Д7 in table 4 of [run/DEBT.md](../run/DEBT.md) (R17). The skill's state has failures in a row instead of "not yet" answers in a row (section 3.1). After a piece of evidence, for each skill that has an answer in it:
1. a failure adds to the failures in a row, and evidence without a failure resets them; two failures in a row remove "mastered";
2. evidence without a failure where the skill has a correct answer without help sets "mastered" by the rule of section 3.1, if m ≥ 5 with this evidence.

One piece of evidence cannot both remove and set "mastered": a failure removes it, and only evidence without one sets it. For a task of one item this is exactly the rule of section 3.1: a correct answer is not a failure, "not yet" is a failure.

**Skills.**
- An item's primary skills get the step and the count, as above.
- The supporting ones — and an anchor skill in this role — get nothing: no step, no counters, no stability. They do not enter P (section 3.1), and the answer tells the model nothing about them, while a shift only upward, on a correct answer, would pull their estimate up without measure. The model derives mastery of an anchor skill from the skills that teach it (section 3.4).

**Stability and practice** — to section 3.2. The stability τ, the time of the last practice and the number of pieces of evidence change for a skill that has an answer in the evidence — an item's primary skill. τ changes once per skill per piece of evidence, by the share of correct ones S̄ among the skill's answers:

τ := max(τ₀; τ · (1 + g · (1 − R))^S̄ · (1/2)^(1 − S̄)).

At S̄ = 1 and S̄ = 0 these are the rules of section 3.2. At R = 0.9 eight correct out of ten grow τ by 1.81 times, half by 1.12, and one miss on a page does not halve τ. The skill's first answer is τ = τ₀, as before.

**Outcomes** (concepts 11 spec):

| Outcome | What goes into the step |
|---|---|
| from the first submission | the items' first verdicts — all "correct"; there was no rung, but a technique on request or a repeat of the task make them answers with help |
| after hint number k | the items' first verdicts; an item first answered after a rung is with help |
| solution shown | the items' first verdicts, and empty and not submitted items are "not solved", S = 0 |
| left | only the first verdicts of the items that there were: "left" is not "not solved" (R82 spec) |
| unchecked | the first verdicts of the items given before it; the outcome itself weighs nothing |
| with a sketch hint | the items' first verdicts, all with help: the sketch is help for all answers of the attempt, like a technique on request (R06) |

Practice is for the skills that have an answer in the evidence, whatever the outcome. An attempt that was left or abandoned without wrong answers, and an attempt of a task with a rubric, give no evidence (the 2026-11-02 amendment to R80 spec; О-17), and evidence without a single answer moves nothing. An attempt of a lesson that the run end "not suitable" closed is, for evidence, like a left one: with wrong answers its outcome is "left", without them there is no evidence (the third 2026-11-02 amendment to R80 spec). The weight needs no fact "segment of a lesson", and the evidence does not carry it (section 3.5).

**Provenance** (concepts 9.4 spec):

| What | Provenance | Weight |
|---|---|---|
| key or stem | from the author | 1 |
| key | checked by the solver — a template (R60, R61 spec) | 1 |
| key | entered by an adult — a task from paper, stage 5 | 1 when the task has a band and a difficulty on the module's scale and the skills are not from AI, and otherwise 0 (the 2026-11-02 addition to О-11) |
| key | derived by AI | 0: there is no verdict (О-11) |
| skill tags | from the author | 1 |
| skill tags | derived by AI | 0 |
| trap | by a wrong value, the program found it | 1 in the count of traps |
| trap | AI supposed it | 0: the evidence does not carry it, and the tutor needs it for the words (R82 spec; concepts 9.4 spec). `engine` does not take it: the items' traps in its candidate are by `trap_values` (architecture 3.8 spec; R107 spec) |

A task from paper that AI itself tagged and solved does not move the rating: its key and skills are from AI, and its weight is 0 (concepts 9.4 spec; О-11). In a module everything is from the author, and a template's key is computed by the program (concepts 9.4 spec), so in the first slice the weight by provenance is always 1.

**Facts of evidence** beyond the outcome (R82 spec):

| Fact | What the model does with it |
|---|---|
| traps that the program found | counts them, below |
| items correct to the end, with their verdict | does not read: it reads the items' first verdicts, and the evidence carries both (R08; the 2026-11-02 amendment to R82 spec) |
| technique at a step | nothing: it is not help |
| technique on request | help |
| leading question on request — from format 0.4 (R108 spec) | help, by its own field `leading_question_on_request` of evidence of kind 2 (R21, R28; subsection 2.8); `technique_on_request` is only the technique on request |
| how many times the learner asked "explain differently" | nothing |
| the task's solution has already gone to the device | help |
| the task is a segment of a lesson (R13 api) | nothing: a task of a lesson is the same kind of attempt (R58 spec), and the evidence does not carry it (R08) |

That the learner saw "givens and unknown" is not recorded (R82 spec). The facts of R82 spec that the model does not read now the evidence carries all the same — as observations, so that a new version can read them by recomputation, — but it does not carry the segment of a lesson (R08; section 3.5).

**Frequent traps.** A trap that the program found goes to the skills that the evidence names for it — by its `skills` in the trap catalog ([the format's text](../spec/format/README.md), section 16). A skill has, for each trap, the number of pieces of evidence where it was found, once per piece of evidence, however many items showed it, and the time of the last. A trap without `skills` or with skills without a state counts with the module. Traps do not change the step: an answer with a trap is already "not yet". Which traps are frequent the reader decides — `engine`, — and the contract `app` ↔ `api` does not give the skill map the model's frequent traps (R99 spec). It can choose as in standalone, where the rule takes the two most frequent from the topic's count "ever" (`04-profile.md` standalone; architecture 3.8 spec).

**A probe** — without a lift. The evidence does not know whether the task was a probe — that is a reason of the engine choice, — and the model counts it as any task. In standalone a correct answer to a probe lifted δ beyond the formula (R299 standalone), and a probe was recognized by the state of the choice rule: "Whether a task is a probe is read off where the child stands when it is answered" (SPEC 2.2 standalone). A pause is a rule of `engine`, and whether to bring a skill back from a pause on a correct answer to a probe it decides, while a correct answer at a low P lifts δ strongly anyway (architecture 3.8 spec).

**Feedback by a rubric** — there is no place for its weight: an attempt of a task with a rubric gives no evidence (R80 spec; О-17), and the weight of feedback will come with a new version together with the answer to О-17.

**Who sets the weights.** All of them are constants of the version of `learner-model`. The format has no fields of weights (R44 spec), and nothing is needed from the next format version.

**How this differs from standalone:**
- first verdict: in standalone an answer was also recorded once (R101 standalone), but a task there has one answer, and here items and submissions;
- help does not change the step and does not set "mastered", like a hint there (SPEC 2.2 and R187 standalone), but help here is also the rungs of the ladder and a repeat of the task;
- the solution shown is "not solved", like "I don't know" there (R93 standalone);
- "left" with wrong answers moves the model by these answers, and without them there is no evidence, as with a skip there: "Nothing that reads answers reads it" (R98 standalone);
- the discount of a series and the division among skills are new: in standalone a task has one answer and one topic;
- there is no lift of a probe (R299 standalone);
- a trap also has the time of the last, not only the number.

**Example.** A workbook page: 10 items `number` (c = 0), all with one primary skill; the band is in first place of `bands`, difficulty 3, so β = 0. Before the evidence θ = δ = 0.5, n = 20, m = 8 and τ = 1, and the skill's last practice was a day ago, R = 0.9. The first verdicts are 8 "correct" and 2 "not yet", without help.
- P₁ = σ(1) = 0.7311, P = 0.9 · 0.7311 = 0.6580, w = 0.786;
- u = 1 / (1 + 9/3) = 0.25: the page weighs 2.5 answers;
- Σ w · (S − P) = 10 · 0.786 · (0.8 − 0.658) = 1.117;
- Δθ = 0.1 · 0.25 · 1.117 = 0.028 and Δδ = 0.2857 · 0.25 · 1.117 = 0.080, and n = 22.5 and m = 10.5;
- there is no failure: 8 of 10 correct. There is a correct answer without help, and m ≥ 5, but v = 0.502, λ = 1.108 − 0.708 = 0.400 and σ(0.400) = 0.599 < 0.71875: no "mastered" yet;
- τ = 1 · 2.5^0.8 · 0.5^0.2 = 1.81.

Without the discount θ and δ would shift four times as much, and m would grow by 10: a new skill would gain m ≥ 5 in one page, and "mastered" could come with the very first — the very early "mastered" that concepts 11 spec and Bradlow, Wainer & Wang 1999 warn about. With the discount, not before the second page.

**The numbers** are constants of the version of `learner-model`: ρ = 1/3, so an attempt weighs no more than three answers, and the failure boundary ½. There are no data under them: they, the full weight of an answer with help and the rule of τ for a page the harness will check — Д7 in table 4 of [run/DEBT.md](../run/DEBT.md). Changing them is a new version, and only through the harness (concepts 11 spec; R163 standalone).

**In code** (R06):
- an item's answer is `answerOf` in `step.go`: the first verdict "correct" or "not yet", and when the solution is shown, "not solved", S = 0, for an item without a verdict that was not submitted or was submitted only empty. Only another "unreadable", "unchecked", an unknown verdict or check kind give no answer, and an unknown outcome is read as "unchecked". "Not solved" is an answer: it goes into N, u and the failure, and without answer slots `Append` rejects it;
- help is `helped`:
  - a wrong answer in the attempt before the item's first verdict;
  - a technique or a leading question on request;
  - a solution that has already gone to the device;
  - the outcome "with a sketch" — help for all answers of the attempt (R06).

  Help does not set "mastered", and does not change the step, the counters and τ;
- provenance is `weighs`: weight 1 for a key and skill tags from the author, the solver or an adult: an adult's key moves the rating when the skills are not from AI (concepts 9.4 spec). A key or tags from AI, "not specified" and an unknown value weigh 0 (R06), and such evidence moves nothing, like evidence without answers;
- traps are `countTraps` after the step, once per piece of evidence for each one found: for the skills from its `skills` that have a state after the step, and without them for the module, with the number and the time of the last. A trap does not create a state for a skill;
- the checks:
  - the tables of outcomes, help and provenance — by tests;
  - traps for a skill and for the module;
  - the property "what is appended equals what is recomputed" — on logs with outcomes, help, provenance and traps;
  - the golden vectors of the fold — with a hint, the solution shown, a technique on request, traps and a task with tags from AI.

### 3.4. Knowing an anchor skill

Decided by R07. An anchor skill `core.<name>` belongs to no module, and it has no state of its own. Whether the learner knows it, the model derives at reading from the state of the skills that teach it — by the field `teaches` of the graph ([the format's text](../spec/format/README.md), section 12). This does not go into the cache: knowing depends on the connections and on the graphs that `api` gives at reading, not on the log (R08 api; the 2026-11-02 amendment to R80 spec). There are no new constants.

**Inputs at reading:**
- the skills' state — the cache (sections 3.1–3.3);
- the graphs that `api` gave: of the current versions of all modules of the log, connected or not, and of the connected modules about which there is no evidence yet (the 2026-11-02 amendment to R07 api; R08 api). A skill in the graph has `teaches` and the assessability, and by the module's tasks the bands where it is taught (section 3.1);
- the set of the engine's modules — connected and suiting the learner, whose graphs `engine` gets (the 2026-11-02 amendment to R07 api). A repeat of a choice restores it by the digests of the graphs that the choice keeps (R08 api), as the card says (architecture 3.8 spec).

**A skill that teaches** is a module's skill whose `teaches` names the anchor skill. It can be mastered only if its assessability is `verdict` and there is a band where it is taught — the band of a task where it is a primary skill of the task or of an item (section 3.1). A skill with `rubric` or `none` and a skill with `verdict` without such a band cannot be mastered.

**The state of an anchor skill** — for each anchor skill that a graph of an engine's module names as a prerequisite or in `teaches` — is one of three, in order:
1. **"Knows"** — at least one skill that teaches it has "mastered", in any module whose graph came, and at any band. A disconnected module does not take knowing away: evidence derived "mastered", and it holds until two failures in a row remove it (section 3.3).
2. **"Does not know"** — otherwise, if the engine's modules have a skill that teaches it and can be mastered. The cause:
   - "not mastered" — at least one such skill has answers;
   - "no evidence" — none has answers.
3. **"Unknown"** — in all other cases. The cause:
   - "nobody to teach" — no engine's module teaches the anchor skill, as in concepts 4 spec;
   - "taught only by skills without a verdict" — in the engine's modules it is taught only by skills that cannot be mastered.

With each state the model gives the skills that teach — by full ID (concepts 3 spec), with their estimate, "mastered" and the mastery band, the highest band where they are taught, the retention and the review date (sections 3.1 and 3.2).

**The threshold** is the same as for "mastered" of a module's skill (section 3.1): knowing has no number of its own. Each skill that teaches teaches the anchor skill whole — "a learner who has the skill can do what the description of each says" ([the format's text](../spec/format/README.md), section 12). So there can be any number of skills and modules that teach, and any of them gives knowing.

**The band** is any. The bands of different modules are not comparable (section 3.1), and which band a module that relies on the anchor skill expects, the format does not say. The model gives the mastery band and the highest band where the skill is taught, and knowing that takes the band into account is Д8 in table 4 of [run/DEBT.md](../run/DEBT.md).

**Forgetting** does not remove knowing: knowing is a function of "mastered", and time does not remove "mastered" (section 3.2). The retention and the review date of a skill that teaches the model gives together with it.

**The grade** does not give knowing: it sets only the starting point (concepts 4 spec; section 3.1).

**Skills with `rubric` and `none`** are without a state. They have no evidence: an attempt of a task with a rubric does not give it (R80 spec; О-17), and a skill `none` is only ever supporting, and the answer gives it nothing (section 3.3). An anchor skill that only such skills teach is "unknown" and does not hold a lesson back, and among skills without a verdict the engine goes by the author's order (concepts 4 spec; R58 spec). concepts 4 spec says the same: "unknown" both when nobody teaches and when only skills without a verdict teach.

**An anchor skill as a supporting skill of a task** gets nothing (section 3.3), and it has no direct state. Evidence may name it — then at the transfer it goes as it is (the 2026-11-02 amendment to R80 spec), — but the model does not read it now. The evidence carries it among the supporting skills of the task and of the item — for future versions (R08; section 3.5).

**A module without a current version** gives no graph (R08 api), and its "mastered" is not seen: the anchor skill is then by the other modules.

**The module index** the model does not read. Which modules teach an anchor skill it says (R63 spec), and the engine may offer the learner a module (concepts 4 spec).

**What the engine does** is its rules (architecture 3.8 spec, "To decide when detailed"):
- "unknown" does not hold a lesson back (concepts 4 spec);
- "does not know" holds back a skill that needs the anchor skill, by the prerequisite rule of `engine`. At "no evidence" the engine may hold the lesson back, give a task of a skill that teaches, or a short check;
- knowing goes away when two failures in a row remove "mastered" from the last mastered skill that teaches, and what happens then with the open skills the engine decides;
- a cycle of prerequisites between modules through anchor skills, where both sides are "does not know", holds both skills back, and the format does not see such a cycle (R50 spec).

**Examples** — by the first slice. In the format's examples school maths of grade 3 teaches `core.division-with-remainder` and `core.fractions-of-a-whole`, and olympiad maths relies on them (R50 spec):
- a grade 3 child connected both maths modules and has not solved the school one yet — "does not know", "no evidence";
- the same child after two pages of school maths where the skill of division with remainder is mastered (section 3.3) — "knows" `core.division-with-remainder`;
- a grade 5 learner connected only olympiad maths — "unknown", "nobody to teach", and the lesson is not held back. If they connect school maths too — "does not know", "no evidence", and what happens with the lesson the engine decides;
- an adult disconnected school maths, where the skill is mastered — "knows"; where it is not mastered — "unknown", "nobody to teach": a disconnected module does not hold olympiad skills back;
- an essay module does not name `teaches`, and it has no anchor skills.

In standalone there is one module, and there are no anchor skills there: this rule has no counterpart in it.

**In code** (R07):
- `Read` in `readout.go` — the modules of the state and the engine's modules, by ID. A module with a state has the start band, the skills, the marks of the trial series and the module's traps that have no skills with a state; an engine's module without a state and with a current version has the starting point by grade, and if this version has no bands, as with a module without tasks, without it: a zero `BandRef`, and the engine still gets the module's reading, while an engine's module without a state and without a current version does not come out. Each needs an entry in "graphs", otherwise `ErrInvalid`, as with `Append`;
- anchor skills — each `core.<name>` that a graph of an engine's module with a current version names as a prerequisite or in `teaches`, with the status and the cause by the rule above. The skills that teach are from the modules of the state and of the engine with a current version: "mastered" exists only for a module with a state, and the map "graphs" is not iterated (section 3.6, rule 6). Each teacher has the full ID, the highest band where it is taught and, if it has practice, the estimate, "mastered" with the band, the retention and the date;
- the mastery band and the highest band are a `BandRef` with the ID from the graph of the current version, and without it the ID is empty; the start band of a module with a state is as recorded in the state, with the ID from the evidence;
- the checks are vectors of each state with each cause by the examples above: both maths modules without evidence of the school one, a mastered and a not mastered skill, only olympiad maths, school maths disconnected, skills without a verdict and a module without a current version; the golden vectors of reading `testdata/read/` — with the log `anchors`, where division with remainder is mastered.

### 3.5. Evidence, the log and recomputation

Decided by R08. The append-only evidence log is the source of truth, and the skills' state is `api`'s cache with a mark of the version of `learner-model`, the release of `modules` and the format window (R80 spec; R06 api). Here is what evidence carries, what the library takes and gives, and why what is appended matches what is recomputed. The decision holds the condition of R80 spec, and Т-22 does not come back.

**Operations** are computations without input or output; their names and types are section 3.7 (R10):
- **band by grade** — the grade and the bands of the module's version → the place and ID of the start band, by the rule of section 3.1. `api` calls it when it assembles the evidence, by the grade from the profile and the bands of the version on which the attempt ran;
- **append** — the state, a new piece of evidence and the graphs → the new state and the evidence with the start band that the appending took. `api` writes both in one transaction, under the lock of the cache row (R06 api);
- **recompute** — the log and the graphs → the state: "append" for each piece of evidence of the log in order, from the empty state. The tail of the log that the cache has not taken into account yet `api` appends the same way, both at reading and before appending (the 2026-11-02 amendment to R06 api);
- **read** — the state, the graphs, the time t₀, the grade and the engine's modules → what the engine and the skill map read: an item's P with retention, the estimates, "mastered", the review dates (sections 3.1 and 3.2), the starting point of modules without evidence (section 3.1) and knowing anchor skills (section 3.4). The engine's candidate task is described the same way as a task in evidence. Reading writes nothing into the cache;
- **encode and decode** evidence and the state — canonically, without loss. `api` keeps the bytes as they are.

**The fold.** Recomputation is a fold of the log from the empty state, and appending is one of its steps over the cache, and the step is the same. So what is appended matches what is recomputed by construction — if the step is deterministic and the cache is kept without loss. Everything in sections 3.1–3.4 is a function of the log in order: the answers of the trial series lie in the module's state, the counters, the stability, the failures in a row and "mastered" go by evidence, the time of the last practice and the time of the last trap are the greatest of the times of evidence. What depends on the time of reading, the connections and the profile, the model computes at reading.
- **The order** is the numbers of `api`'s evidence in the profile's log, not the time: a number is taken under the lock of the cache row (R06 api), and the time of practice can go backward too (section 3.2).
- **Locality.** A step reads only the graph of the module of its own evidence. So the graphs that come only at reading — of connected modules without evidence (the 2026-11-02 amendment to R07 api) — and extra graphs change nothing, and a module that enters the log for the first time does not touch others. This is the rule for future versions too: a step that read others' graphs would make the cache depend on modules outside the log, against the 2026-11-02 amendment to R80 spec.
- **The input "graphs"** names, for the module of each piece of evidence being appended, the graph of its current version with `retired_skills` (R08 api) — or "no current version". At recomputation these are all modules of the log, and at appending the modules of the new evidence and of the tail. A module of evidence without an entry is an error: otherwise appending would silently not transfer what recomputation will transfer. Evidence of a module without a current version moves nothing and waits in the log until a version appears, and its release will change the cache's mark (R08 api).

**Evidence** is one per finished attempt (R80 spec). An attempt of a task with a rubric and an attempt that was left or abandoned without wrong answers do not give it (О-17; the 2026-11-02 amendment to R80 spec). It carries observations — what there was in the attempt — not the model's conclusions: the rule of R06 reads them now, and another rule will read them later, by recomputation. It holds everything that R82 spec names and everything that the model reads, and beyond that a little for future versions: the IDs of bands, the supporting skills, the places of items for traps and the number of wrong answers before the first "correct". The first three are the module's content, not the learner's data, and the fourth is what partial credit cannot be introduced without later. Nothing more: data at a minimum (R09 spec).

About the attempt:

| What | What the model does with it |
|---|---|
| the number of the kind of evidence — 1, 2 or 3, the kinds below (R28, R45) | reads the evidence by its kind, below |
| the task — the full ID, and for a template the template's ID (concepts 3, 9.3 spec; R60 spec). The module comes from it | the module's graph for the step |
| the task's band — the place in `bands` of the version on which the attempt ran, and the ID; the difficulty from 1 to 5 | β by place and difficulty (section 3.1); does not read the ID |
| the module's start band — place and ID | θ₀ by place (section 3.1); does not read the ID |
| the task's skills — primary, by full ID, and supporting, by full ID or `core.<name>` | primary — P, the step and the count (sections 3.1 and 3.3); does not read the supporting ones |
| the provenance of the key and of the skill tags — the author, the solver, an adult or AI (concepts 9.4 spec) | the weight (section 3.3) |
| the outcome — from the first submission, after hint number k, with a sketch hint, solution shown, left, unchecked (concepts 11 spec) | section 3.3 |
| the time of practice — the time of the attempt's last action (R10 api), in whole seconds UTC | retention and stability (section 3.2) |
| technique — at a step and on request, whether each was there (R82 spec) | on request — help, at a step — nothing (section 3.3) |
| whether "hint" gave the task's leading question, from kind 2 (R82 spec; R21, R28) | help (section 3.3); in kind 1 not recorded and read as "did not give" |
| how many times the learner asked "explain differently" (R82 spec) | nothing |
| whether the solution of this task went to the device earlier; a template has no such fact (R82 spec) | help (section 3.3) |

About each item, by place, from 1:

| What | What the model does with it |
|---|---|
| the item's own skills — primary and supporting, if it has them ([the format's text](../spec/format/README.md), section 14) | like the task's skills, in their place |
| the answer slots in order — the check kind and, for a choice, the number of options | guessing c (section 3.1) |
| whether the item was submitted, and the reasons of its submissions without a verdict before the first verdict — "unreadable" with the reason `blank` or another, "unchecked" (concepts 9.1 spec) | when the solution is shown, an item not submitted or only empty is "not solved", and one only with another "unreadable" — nothing (section 3.3) |
| the first verdict — "correct" or "not yet" — and the number of the attempt's wrong answers before the submission that gave it | the item's answer, and a rung before it is help (section 3.3) |
| the verdict at the end (R82 spec) and the number of wrong answers before the first "correct" | nothing: partial credit is later (section 2.3) |

About each trap that the program found, once per piece of evidence (R82 spec):

| What | What the model does with it |
|---|---|
| the trap's full ID, its skills from the catalog of the attempt's version — by full ID — and the places of the items where it was found | frequent traps (section 3.3); does not read the places |

**What evidence does not carry:** the IDs of the profile and of the attempt, the package digest, the module's version and the number in the log — these are `api`'s columns and rows (R06, R07 api); the segment of a lesson — it does not change the weights (R58 spec; section 3.3), and it has no reader, and this is the answer to the question of R13 api; the mode "step by step" or "on my own"; that the learner saw "givens and unknown" — this is not recorded (R82 spec); AI's candidate traps (R82 spec; concepts 9.4 spec); the work, what was read, the words of hints and the rungs with their source. The module and the time `api` also keeps in its own columns, but the model reads them from the evidence, and the columns serve `api` itself.

**The kind of evidence** changes only by addition — of an optional field or a new value of a list — like the members of the module index (R63 spec). The number of the kind is in each piece of evidence, and the log is not rewritten (R80 spec), so:
- each version of the library reads evidence of all earlier kinds. A field that old evidence does not have is "not recorded", not "no" and not 0: how to read it is decided by the version that added the field — usually the way the earlier one read its absence;
- and of a newer kind too. `api` rolls back only forward — a pull request with the earlier code (Т-19) — and the earlier version will recompute a log where they already are. An unknown field it skips, an item with an unknown check kind or an unknown verdict it counts as an item without an answer, and an unknown outcome as "unchecked": the items' verdicts before it go, and the outcome itself weighs nothing. Recomputation does not fail and computes the same;
- a field that changes a weight acts only in a version that knows it — under its own cache mark.

**Kind 2** adds field 20 `leading_question_on_request` (R28) — whether "hint" gave the leading question of a task of format 0.4. Kind 2 always writes the field, `false` too, and the strict reader requires it. In kind 1 the field is "not recorded": the writer does not write it, the reader of kind 1 rejects field 20, and `EncodeEvidence` and `Append` reject kind 1 where the field is set — so the bytes of kind 1 from `api`'s log stay the same bytes after reading and writing. "Not recorded" is read as "did not give", that is, the way the earlier version read kind 1: before format 0.4 there were no leading questions. The value `CHECK_KIND_TIME` needs no kind and is read the same in every kind: an unknown value is already an addition that R08 allows.

**Kind 3** adds field 21 `level_set_by_adult` the same way (R45): the record is then no attempt but an adult's correction of the module's level, the second fact of the log, which `api` appends by its method of row 6.28 of `spec`'s `COORDINATION.md`, and `LevelCorrection` gives it. It names no task: `task_id` is `<module>/level/<band ID>`, the module read from it as from an attempt's, `task_band` the band set, `difficulty` 3, both provenances the adult, `practice_time` the time, and it has no skills, items, traps or other facts of an attempt — the codec and `Append` refuse a correction other than `LevelCorrection` gives. An earlier version skips field 21 and reads the record as evidence with no answer, which moves nothing. What the step does with it is section 3.1, "The rules of R157 spec". In the export the record is a line of the log beside the attempts, and the comment of field 21 in `evidence.proto` says what it is (R135 spec).

Standalone did not read a profile newer than itself, because the earlier instance would rewrite the file and lose the unknown (`04-profile.md` standalone; sources 1.2 and 1.9). Nobody rewrites the log, and reading it while skipping the unknown loses nothing.

**The transfer** by `retired_skills` goes inside the step, at recomputation and at appending, by the graph of the module's current version, and the log does not change (R08 api; the 2026-11-02 amendment to R80 spec). Each recorded primary skill x of an item goes to the skills map(x):

| x in the graph of the current version | map(x) |
|---|---|
| in `skills`, with assessability `verdict` | x itself |
| in `retired_skills` | the skills of its `replaced_by` with `verdict`; at `[]` — nobody |
| in `skills`, but now with `rubric` or `none` | nobody: such skills have no state (section 3.4) |
| neither here nor there | nobody, without an error: the evidence waits until the graph names the skill (R08 api) |

- **The share.** J is the number of recorded primary skills of the item whose map(x) is not empty. Skill s gets in the item the share |{x ∈ primary of the item : s ∈ map(x)}| / J of the step of δ. Without a major version this is the 1/J of section 3.3. Each part of a split skill gets the earlier one's share whole, a merged skill the sum of the shares, and a skill that reached no one does not take a share from others.
- **The count and practice** are by the skills after the transfer, once per item: N_s, u_s, m, the stability, the time of the last practice and the number of pieces of evidence are counted as if the item named s. Skills merged in one item are one skill, and their count and practice do not double. N and u for θ, and with them the trial series, do not depend on the skills.
- **δ̄ and R̄** of an item are the means over its J recorded skills, and for each recorded skill δ and R are the means over its skills after the transfer. For a split skill the δ of each part goes exactly as the earlier one's δ, and the item's δ̄ does not change with a major version, whether the item has other primary skills or not. Bit for bit when there are two parts, since (δ + δ) / 2 = δ exactly, and at three or more up to the rounding of the mean. Likewise a merge: two skills in one item give exactly what an item with one replacement would give, and three up to the rounding of the share k/J (R08). Example: an item has primary A and D, J = 2, and A is split into B and C. B, C and D each get half of the step, and δ̄ = (½ · (δ_B + δ_C) + δ_D) / 2 — exactly the earlier (δ_A + δ_D) / 2, since δ_B = δ_C = δ_A.
- **An item without skills after the transfer**, J = 0, moves only θ and n of the module: δ̄ = 0, R̄ = 1 and w = 1. The level is the module's, and the band and the difficulty stay with the task. Such an item does not set a time of practice, so the next piece of evidence counts the module's pause for forgetting the level from the skills as before. Evidence without a single answer does not forget the level: it changes nothing after a pause too (section 3.2).
- **"Mastered"**, the stability and the traps of a replacement are derived from the transferred evidence by the same rules (sections 3.1–3.3), and the flag "mastered" does not carry over separately. The mastery band is by the bands where the replacement is taught in the current version. A part of a split skill that is primary in no task cannot be mastered (section 3.4), and knowing an anchor skill through it goes away.
- **A trap's skills** go by the same mapping, and a trap that has no skills left counts with the module (section 3.3).
- **Anchor skills** among the supporting skills go as they are, if the dictionary of a format version in `api`'s window knows them (the 2026-11-02 amendment to R80 spec). The model does not read them, and it does not need the dictionary yet. When it does, it will be an input of the version that starts reading them, and the cache's mark already carries the window (R06 api).

| What the major version did | What the skill of the current version gets |
|---|---|
| renamed A to A′ | A′ — everything A got |
| split A into B and C | B and C — each everything A got |
| merged A and B into C | C — the items of A and the items of B, and an item with both once, with the share 2/J |
| retired A without a pair, `[]` | nobody; an item where A was the only primary skill moves only θ |
| A is a skill of a withdrawn newest major version that the current version does not know | nobody, until the next version names it in `skills` or `retired_skills` |

**Bands.** Recomputation reads a band's place as recorded: the scale of bands is one for all versions of the module, and the shift 2.5 · (p − 1) for earlier evidence is the same (section 3.1). A major version does not change the places of earlier bands: a band is not removed or reordered, and a new one stands only after the last. With format 0.3 this is the rule of comparing versions `module-version.bands-kept`: the bands of each published version stand first in the new one, in order and with the same IDs; a band's ID does not name another band, and `manifest.bands` is among the places that the format keeps in all versions (the 2026-11-02 amendment to R89 spec; [the format's text](../spec/format/README.md), section 27). The evidence carries the band's ID beside its place, but recomputation does not read it: it is for a table of correspondence of bands, if one is ever introduced, and it cannot be appended to the log later.

**The start band** of a module is in each piece of evidence of the module (section 3.1):
- `api` takes it from the rule "band by grade" when it assembles the evidence;
- the step — both at appending and at recomputation — takes, for a module that already has a state, the band from the state, and for a module without a state, the band of the evidence. Appending gives back the evidence with the band that it took, and `api` keeps it so.

The rule is one for appending and recomputation, so it does not set what is appended apart from what is recomputed. Two first attempts of one module on two devices, begun before and after a change of grade, are appended in turn, under the lock of the cache row, and the second takes the band from the state. When `wrong_key` cuts off the first evidence of a module, the next one gives the band, and it carries the same one: appending took it from the state (R63 spec; R08 api). Evidence of a module without a current version does not create a state, and the band is then pinned by the first evidence that creates it.

**`wrong_key`.** The evidence of attempts on a version withdrawn because of a wrong key `api` cuts off by the package digest before it gives the log (R63 spec; R08 api). Beyond this the library needs nothing: recomputation is a function of the log that came, and the start band, which the first evidence of a module pins, each following one carries too.

**There is no snapshot.** `api`'s cache — a row per profile with the number of the last evidence taken into account (R06 api) — is the state at the end of the log. Recomputation is needed only when the cache's mark changes (R80 spec; R06 api), and a snapshot under the earlier mark is then of no use, and the transfer is not applied to it (the 2026-11-02 amendment to R80 spec). Recomputation is linear in the log: a long log recomputes in milliseconds (measurements 11).

**A repeat of the engine's choice** by the log (architecture 1 spec):
1. the state is a recomputation of the log up to the number of the last evidence that the choice keeps: without the evidence marked `wrong_key` in the index of the kept release, by the graphs with the kept digests and with the version of `learner-model` that the choice names (R80 spec; R06, R08 api);
2. the reading is at the time that the choice keeps (the 2026-11-02 amendment to R06 api), with the engine's modules restored by the digests of the graphs (R08 api), and with the grade that the choice keeps: without it the first choice in a module without evidence would not repeat if the grade has changed since (architecture 1 spec; R11).

A repeat with an earlier version of `learner-model` runs with its code, not the current one: how `api` and the engine keep it, they decide.

**The same bits.** The thresholds — "mastered", n < 5, m ≥ 5, failures in a row — make the state discontinuous: a difference in the last bit sets "mastered" or does not, and from then on the states diverge forever. So "matches" in the condition of R80 spec is bit for bit, and on any build of one version of `learner-model`: on any architecture, GOAMD64 level, processor and Go version. While a rollout is going on, the earlier and the new image append under one mark (the 2026-11-02 amendments to R06 api), and the mark does not carry the build. How the code achieves this R09 decided — section 3.6.
- The thresholds of the sums of weights are compared with a tolerance ε = 1e-9: "n < 5" is n < 5 − ε, and "m ≥ 5" is m ≥ 5 − ε. The bits are the same anyway, so the tolerance is not about matching but about meaning: two pages of 8 and 13 items give exactly 5, 2.4 + 2.6, and a floating-point sum must not decide otherwise. Each piece of evidence with answers adds no less than 1 to n and to m, since N · u = 3N / (N + 2) ≥ 1, so a sum of five terms is already no less than 5, and it is exactly 5 only from five ones, and that floating-point number is exact. At the boundary there can only be a sum of no more than four fractions, and if it is not equal to 5, it stands from 5 no less than by one divided by the product of their denominators N + 2. At N ≤ 170 this is 1/(172 · 171 · 170 · 169) ≈ 1.18 · 10⁻⁹ — more than the tolerance. So the tolerance is right while an attempt has no more than 170 items with an answer — Д9 in table 4 of [run/DEBT.md](../run/DEBT.md). The format's guideline is up to 20 items per page ([the format's text](../spec/format/README.md), section 14).
- The time of practice is whole seconds UTC, and Δt = (t − t_last) / 86,400 is computed from integers, so Δt too has the same bits everywhere.
- The cache is kept without loss: the library sets the encoding, canonically, and `api` keeps the bytes as they are, so neither the database nor a conversion to another form changes the bits.

**Export.** The log, the current state and the choices with reason codes are `api`'s own part of the export (R06 api; R76, R80 spec). The log goes whole, the evidence of withdrawn versions too: it is the learner's data. For the log to be read without the library too (R09 spec), the export carries a description of the kind of evidence and of all its kinds: the `.proto` schema of the needed version with comments in English (section 3.7; R11). The export is JSON by the schema (R135 spec; [part `export/`](../spec/export/README.md)): `api` writes the log and the state as JSON of their messages by the protobuf-to-JSON mapping with field names as in `.proto` — like `buf convert` with `use_proto_names=true` below — and beside it puts `common.proto`, `evidence.proto` and `state.proto` of the version of `learner-model` that computed, from its release file, which `release-verify` accepted, where they lie under `learner-model-proto-<version>/proto/novatrail/learnermodel/v1/`, into `api/proto/novatrail/learnermodel/v1/` of the export. There are no bytes in the export, but they are assembled again from the JSON and `.proto`. For this `learner-model` keeps the comments of `.proto` such that a reader of the export understands each field without code, and its release file, the schema's archive, is the file of the export (section 3.7; R26).

**П-9.** The log is deleted only whole — with the profile or by a term (R80 spec). Recomputation under a new version reads the log from the start, so a term that cut off old evidence would break it, and there is no snapshot from which recomputation could start: it would not survive a change of the mark.

**How the tests check the condition of R80 spec:**
- a property: for random synthetic logs and graphs, appending one at a time, with encoding and decoding of the state between the steps, gives the same bytes as recomputation. The logs are with any evidence filtered out, with extra graphs, with time going backward and with two first pieces of evidence of a module at different grades;
- a module without an entry in the input "graphs" is an error, and "no current version" — the evidence waits;
- transfer vectors: renaming, splitting, merging in one item and in different ones, `[]`, an unknown skill, a skill that became `rubric`, a chain of retirements;
- the boundary of the thresholds: pages of 8 and 13 items give n = 5 and m = 5;
- reading kinds: evidence of each earlier kind and evidence of a kind newer than its own — with an unknown field, check kind, verdict and outcome;
- golden vectors — the state by the log, bit for bit, on the builds that section 3.6 names;
- standalone's test vectors — with a tolerance of 1e-4, as they are printed (SPEC 2.8 standalone), for a task of one item at Δt = 0 and apart from its own vectors bit for bit (section 3.1).

**In code** (R08):
- the transfer is `transfer` in `step.go`, the mapping map(x) by the table above. The graph of the current version names the replacements already as skills of its own version, so a chain of retirements is one step;
- the step:
  - the item's recorded primary skills — by ID;
  - J — those whose map(x) is not empty;
  - δ and R of a recorded skill — the means over its replacements;
  - a skill after the transfer gets k/J of the step of δ, and it counts the item once. Without a retired skill the bits are those of the step of section 3.3;
- traps — by the skills after the transfer, without repeats, and without skills with a state, with the module;
- `Append` rejects a graph whose `Retired` or `replaced_by` entries are not in order, where a retired skill is still a skill of the graph or a replacement is not a skill of the graph;
- the input "graphs", the start band and the tolerance at the thresholds are section 3.1, "In code", and reading the unknown section 3.3, "In code";
- the checks — by the tests from the list above:
  - the property "what is appended equals what is recomputed" — on logs with retired skills, evidence filtered out and an extra graph;
  - transfer vectors — the equality of a recomputation with the old IDs and of a log with the new ones, by the bytes of the state: renaming, splitting next to another skill, merging in one item and in different ones, `[]`, an unknown skill, a skill that became `rubric`, a replacement with `rubric`, a chain and traps;
  - the boundary of the thresholds on a log with a transfer;
  - evidence of a kind newer than its own — by bytes;
  - the golden log `testdata/fold/transfer.txt` — a module through three versions.

**How this differs from standalone.** There is no log there: the profile file is the truth — a summary for each topic, which is not cut, and a window of the last 20 answers (sources 1.2). The profile's format is changed by a chain of migrations — pure functions — and standalone rejected levels derived anew from the answers: "which breaks the day the step of the rating changes" (R165 standalone). Here it is the other way round: a new version recomputes the state by the log, and a change of the step changes the earlier estimates too — so R80 spec requires. When the rule "mastered" changed in standalone, the earlier "mastered" were reset (R187 standalone), and here recomputation derives them anew, by the new rule. The log does not migrate, and instead of a chain of migrations there is the rule of reading kinds.

### 3.6. The build and the same bits

Decided by R09. Here is how the library keeps the promise of R08 — what is appended matches what is recomputed, bit for bit, on any build of one version of `learner-model` (section 3.5) — and what of it is needed so that the library can be built for the device.

**The promise.** The same bits hold for a build by Go's compiler, `gc`: on any architecture it knows, 32-bit too, at any GOAMD64 and GOARM64, on any processor and in any version of Go no lower than the one `go.mod` names. TinyGo and gccgo are outside the promise: they are other compilers, and there are no checks for them. The assumption is the default floating-point environment, rounding to nearest: Go does not change it, but C code called through cgo could.

Addition, subtraction, multiplication, division and square root of IEEE 754 are correctly rounded, so each such operation gives one result on any target. Builds diverge in only three places (sources 4.4 and 4.5):
- **fusion.** The compiler fuses a multiplication with an addition or a subtraction into one operation with one rounding on arm64, ppc64le, s390x, riscv64 and loong64, and on amd64 with an addition, and only at GOAMD64 ≥ 3. On wasm, 386 and mips it does not. A division by a constant with an exact reciprocal it makes a multiplication, so `x/4 + y` fuses too. It also fuses across statements — `t := a*b; r := t + c` — and through inlining, including into the consumer's code, which `api` builds with its own flags. 32-bit arm has its own rules for multiply-add, `MULAD`, without the condition `useFMA` (`ARM.rules`), and the differential check does not see them; whether this instruction rounds the product separately is not checked here — Д10 in table 4 of [run/DEBT.md](../run/DEBT.md);
- **`math`.** `math.Exp` has its own assembly on amd64, arm64, loong64, riscv64 and s390x, and on amd64 it chooses the fused path by the processor right at run time; `math.Log` has it on amd64 and s390x; `math.Pow` goes through them. On wasm and 386 both are pure Go;
- **conversions.** A floating-point number outside the range of an integer the Go specification leaves to the implementation: wasm saturates, amd64 gives something else, and the compiler holds a disabled change of behaviour on amd64, `converthash`, which may switch between versions of Go. The encoders of the standard library do not promise bytes either: `encoding/json` in Go 1.27 is already built from v2.

**The rules of code** — for everything that computes the state and the reading:
1. The result of each `*` and `/` on floating-point numbers goes through `float64(…)` at once, unless it is an operand of another `*` or `/` or of a comparison: `float64(a*b) + c`, but not `a*b + c`, not `t := a*b; t + c` and not `x/4 + y`. No `*=` and `/=` on such numbers, no `complex`, no type parameters in their arithmetic: there an explicit conversion either cannot be written or it is unclear whether it holds. The compiler must round an explicit conversion — "An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding" (the Go 1.27 specification) — and the rules of fusion do not cross it.
2. From `math`, only functions with an exactly defined result: `Sqrt`, `Abs`, `Floor`, `Ceil`, `Trunc`, `Copysign`, `Signbit`, `Frexp`, `Ldexp`, `Float64bits`, `Float64frombits`, `IsNaN`, `IsInf`, `Inf` and the constants. `Sqrt` is correctly rounded: the compiler makes it an instruction on all the architectures it knows, on wasm `F64Sqrt`, and the software version is "Return correctly rounded sqrt". Builds with software floating point — Д10 in table 4 of [run/DEBT.md](../run/DEBT.md). `Exp`, `Log`, `Pow`, `Exp2`, `Log1p`, `Expm1`, trigonometry and `FMA` are not taken.
3. exp and log are its own: a port of the pure Go `exp`, `expmulti` and `log` from `math` of Go 1.27.3 — the code that `math` itself runs on wasm and 386 — with `float64()` at each step. On top of them: σ(x) = 1 / (1 + exp(−x)), the power a^s = exp(s · log a) with exact branches s = 0 and s = 1 — so that the rule of τ of section 3.3 at S̄ = 1 and 0 exactly matches the rules of section 3.2 — and the retention 1 / √(1 + F · Δt / τ) (section 3.2). All this is the public package `detmath`, so that `engine` takes the same port. The port's file carries Go's BSD-3-Clause notice and the Sun Microsystems notice from FreeBSD msun, which asks to be preserved: "Permission to use, copy, modify, and distribute this software is freely granted, provided that this notice is preserved".
4. NaN and ∞ do not arise. The arguments of exp and log are clamped so that the result is finite and not subnormal. The likelihood term of the trial series chooses a branch by the answer S, rather than taking a factor (1 − S): σ at the edge of the grid θ₀ ± n · σ₀² (section 3.1) can be exactly 1 — from an argument of about 36.7 — and (1 − S) · ln(1 − P) would give 0 · (−∞) there. 1 − P at R̄ = 1 is computed as (1 − c) · σ(−x), without a subtraction from one. The encoding rejects NaN and ∞ — as the last barrier.
5. Floating-point numbers into integers only by a function that checks the range; integers into floating-point numbers only with an absolute value no greater than 2⁵³, where the conversion is exact. Everything stored, and the time, is `int64`: `int` on 386 and arm is 32 bits. The time is whole seconds of Unix UTC (section 3.5), and the review date is converted into seconds by rounding up: "due" does not come before the threshold.
6. Order. Sums and the repeated fields of the encoding are in canonical order: items by place, skills, traps and modules by ID (sections 3.3 and 3.7). The iteration order of a Go map is not specified — "The iteration order over maps is not specified" — so neither the computation nor the codec iterates over maps. Sorts have a total order: `slices.SortFunc` "is not guaranteed to be stable". IDs are compared byte by byte — without `strings.ToLower`, `strings.EqualFold` and the tables of the package `unicode`, which have their own version of Unicode (`unicode.Version` is "17.0.0" in Go 1.27.3).
7. No global mutable state — `api` calls the library from many requests at once — and no goroutines of its own. No recursion as deep as the length of the log.
8. Imports only from the list below. No cgo, no `.s` files, no `//go:linkname`, no build tags in the core.

**Dependencies** — only Go's standard library, by the list, in the code and in the tests alike; `go.mod` without `require`.
- Code: `math` — only the functions of rule 2 —, `cmp`, `slices`, `errors`, `strconv`, `strings` — only the byte-wise functions —, `unicode/utf8` and `encoding/binary` — the byte order, not the encoder. A new package enters the list by a change of R09.
- Not taken: the Go module `spec` — it embeds the whole `format/`, `//go:embed all:format`, and that would land in every build, WASM too, while the library's types are its own (R28 spec); `module-kit`; third-party modules — each one more place where bits and bytes can change with its version; `os`, `io/fs`, `net`, `syscall` and `time` — the time comes as a number; `math/rand`, `math/rand/v2` and `hash/maphash` — the model needs no randomness, and should it need some, that is the generator of R59 spec by [the text of the format](../spec/format/README.md), section 22, with its vectors; `sync`, `unsafe`, `reflect` and `runtime`; `encoding/json` and other encoders — the codec is its own (section 3.7).
- Tests have their own list, a wider one: `testing`, `fmt`, `bytes` and `runtime` — to learn that `math` on the build is pure Go, in wasm and on 386; `math.Exp` and `math.Log` — as the reference of precision; `math/rand/v2` — PCG with pinned seeds for the property logs of section 3.5, the same on all targets; `crypto/sha256` — for the hashes of the vectors; `os` and `path/filepath` — to read `testdata/`.

**For the device.** The core builds for the browser, `GOOS=js GOARCH=wasm`, and for WASI, `GOOS=wasip1 GOARCH=wasm`, as it is: without build tags and without code of its own for them. The glue with the application — `syscall/js` or `//go:wasmexport` — is a separate package, and only at stage 5. Neither port is first-class in Go: `go tool dist list -json` gives them `"FirstClass": false`. Go's build for wasm is "around ~2MB", after gzip "around 500kB"; TinyGo gives less, but supports "a subset of the Go language" (sources 4.7), and it has no promise of the same bits. The application `app` is web, without offline and WASM (R03 app), so the device is the browser, and the build for it is js/wasm.

**What CI checks it with.** The shared wrapper `go.yml` of `.github`, of the release that `.github/workflows/ci.yml` pins, calls the repository's recipes `ci-lint`, `ci-test`, `ci-vuln` and `ci-licenses` on `ubuntu-24.04`, in the stage `toolchain` of its Dockerfile (sources 5.2), and with the input `same-bits` four more jobs of its own for the same bits: "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules" (R107 spec; R13, R18). The repository's own jobs — "Proto", where buf and `ci-release` run, and "Gate", the gate of `sim` below — stand in `ci.yml` beside the shared ones. CodeQL reads Go from the build.
- `ci-lint` — `gofmt`, `go vet` and the comparison of the image's gopls version with its lock file `.devcontainer/gopls`: the gopls step and the lock file come from `spec`'s template, and the comparison holds them (R14, R19, R23). `ci-vuln` is govulncheck and `ci-licenses` go-licenses with the wrapper's list of licenses, at the versions the `justfile` names.
- `ci-test` — the tests with `-race`, then again for the coverage into `coverage/go.out`, which the wrapper's job `coverage` sends to Codecov, with `-covermode=set`, on amd64 at GOAMD64 v1: the two runs apart are faster than one with both (R43; measurements 8). "Go / targets" runs the same tests without `-race` and coverage — Go does not support `-race` on all targets, and 386 and wasm are not in its list:
  - at GOAMD64 v3 — where the compiler fuses an addition;
  - at v1 with `GODEBUG=cpu.all=off` — without the processor extensions by which `math` chooses its path at run time;
  - at GOARCH=386;
  - in js/wasm — in the container's Node, through `go_js_wasm_exec` from `$(go env GOROOT)/lib/wasm`. Here a test also compares the port with `math.Exp` and `math.Log` bit for bit: in wasm `math` runs the same pure Go;
  - in wasip1 — in `node:wasi` of the same Node, preview1, with an executor that `go.yml` writes itself: Go's `go_wasip1_wasm_exec` knows only runtimes that the stage lacks (R34).
- The rules of code — "Go / rules", the tool `tools/coderules/` of `.github` of the same version as `go.yml`, on `go/types`: imports and functions of `math` only from the lists, each `*` and `/` on floating-point numbers by rule 1, in the tests too, no `*=` and `/=`, no `complex`, no type parameters in arithmetic, no cgo, `.s`, `//go:linkname` or build tags — a file name that the `go` command reads as a build constraint, `_linux.go` or `_amd64.s`, is a tag too —, `go.mod` without `require` and `replace`, no nested modules but those named; in the code also rules 6 and 7: no range over a map, no function of `strings` and `strconv` that reads tables made from Unicode, no `go` statement. The lists are `.github/coderules.json`: the imports of the code and of the tests, the `math` functions of the tests — `Exp` and `Log` —, an empty `requires` and the nested module `.devcontainer/gopls`, which neither the code nor a test reads: the zip of a Go module carries no nested module (R18). The control that breaks each rule of this repository is `.github`'s `tools/coderules/testdata/learner-model`. The rules forbid build tags, so the code is one on all targets, and the tool runs once, on the host.
- The differential check of assembly. The library's packages are built for arm64, ppc64le, s390x, riscv64, loong64 and amd64 at GOAMD64 v3 twice — as usual and with `-gcflags=-d=fmahash=n`, which turns off all the rules of fusion — and the assembly of `-S` must match, and the second build must not print a single "fmahash triggered". This is "Go / fma-diff", and the control with `a*b + c`, `x/4 + y` and `c − a*b` belongs to `go.yml` itself: its check must catch it on each target. The check sees the library's code; the fusion that inlining and PGO would bring into the consumer's code is excluded by rule 1, and that is checked by "Go / rules".
- "Go / arm64" — on `ubuntu-24.04-arm`, the same stage `toolchain` of the same Dockerfile, the same `ci-test` and the runs of "Go / targets" for its host, on arm64 hardware. arm64 runners are free for public repositories (sources 5.1).
- buf — `buf lint`, `buf breaking` and `buf convert` (section 3.7).

**`detmath` in code** (R09):
- `Exp`, `Log`, `Sigmoid`, `Pow` and `Retention` and the constant F, `RetentionFactor`, over the port of `exp`, `expmulti` and `log` in `detmath/port.go`, with the BSD-3-Clause notices of Go and of Sun Microsystems;
- the clamp by rule 4: the argument of `Exp` and `Sigmoid` is in [−708; 708], the argument of `Log` in the positive normal numbers, from 2⁻¹⁰²² to the largest `float64`; for `Pow` the base is as for `Log`, and the exponent is in the finite numbers, so s · ln a is never 0 · ∞ and 1ˢ = 1 for any s; for `Retention` Δt is in [0; the largest `float64`], τ in the positive normal numbers, and F · Δt / τ is no greater than the largest `float64`. NaN goes to the lower bound: for Δt that is 0 and R = 1. So no function gives NaN, ∞ or a subnormal number, whatever comes in;
- the golden files — `detmath/testdata/golden.txt`, a table "input → bits", and `detmath/testdata/sweep.txt`, the SHA-256 of the bits of each function on 100,000 PCG inputs with pinned seeds, half of them raw bits, every kind of `float64`. `DETMATH_UPDATE=1 go test ./detmath` writes them, and they change only with the version of `learner-model`.

**The release** (R16) — `release.yml`, the shared release of `.github` of the release it pins, in the stage `toolchain-core` (R43, R48), on the line its `RELEASE_LINE` names: a commit of `main` with green CI that changes one of its `paths` — what consumers of the Go module build and what builds the release file — gets the next patch of the line, a tag and a signed file, the archive of the schema (section 3.7). Consumers take the Go module by this tag. The release waits for the gate too: it goes after green CI on `main`.

**The gate of `sim`** (R20) — the jobs "Gate / part k" and the check "Gate" in `ci.yml`, on each pull request and each commit of `main`; "Gate" is a required check of the ruleset of `main`:
- `sim` — as source from git, by the tag and the SHA of its commit, `SIM_TAG` and `SIM_COMMIT` in the workflow's `env`; what each pin brought — the base, the modules, the bands, the criterion — is R20. The step "sim tag", in each part, checks by `git ls-remote` that the tag stands on this commit, as `sim`'s recipe checks the libraries' tags. Neither Dependabot nor a shared workflow raises them: the tag is raised by a session of this repository, in the order of releases of R103 spec (R15);
- four jobs "Gate / part k", each `sim`'s `just gate-part k 4` in the stage `toolchain` of `sim`'s devcontainer of the same tag, and the check "Gate" after them, `just gate-check bands`, which merges the parts and holds them as `just gate` does, all on `ubuntu-24.04-arm`, with `sim`'s cache of Go and the limit of time of `sim`'s `gate.yml` a job (R43). The pull request's code is put in by `GATE_LEARNER_MODEL`: the recipe writes `go.work` in a temporary folder and does not commit it, `engine` is taken by the commit from `sim`'s `gate.json`, and the harness's `go.mod` is left untouched. The table of numbers against the bands goes into the run's summary;
- the gate is two-sided: a number beyond its band is a red gate, for the better or for the worse. A change that moves the numbers first records the bands in `sim`, and only with its new tag does the model's pull request turn green (R103 spec; 3.9 `sim`). The gate also holds the invariants I1–I4 and the build: the model's `go` line is the one `sim`'s `gate.json` names, and the libraries are built from `go.work`.

**Golden vectors** — files of `testdata/`, the same for all targets:
- the state by the log — the bytes of the encoding (section 3.7) after each piece of evidence;
- the output of the reading at the time t₀ — the P of an item, the retention, the review date, "mastered", the starting point, the mastery of anchor skills (sections 3.1–3.4);
- the table "input → bits" for exp, log, σ, the power and the retention: it tells a divergence in them from a divergence in the model.
- the bytes of the codec — `testdata/wire/`, evidence and state in binary form and in JSON (section 3.7).

The vectors change only with the version of `learner-model`: the cache's mark is this version (R80 spec). The test vectors of standalone are with a tolerance of 1e-4, as they are printed (section 3.5).

**What the experiment showed** — a one-off program and the tests of `detmath` (measurements 1): the port with `float64()` gives the same bits on all six builds, and on 386 and in wasm it equals `math.Exp` and `math.Log` bit for bit; the own assembly of `math` on amd64 diverges from the port — `math.Exp` by up to 2 ulp, `math.Log` at the edges of the domain — and the bits of `math.Exp` change with the processor even for one binary. Code without `float64()` changes bits only where the compiler fuses operations — rule 1 at work — and this is caught by "Go / rules", "Go / fma-diff" and the golden files in the run at GOAMD64 v3. `detmath` clamps the arguments within the edges of `math`.

**Т-10.** Compute on the server, as Т-10 recommends: the learner model lives in `api`'s database (R80 spec). A build for the device stays possible, and with the same bits: the device, if it starts computing, will get from the same log the same state as the server. How offline works — whether the device computes by itself, or the server chooses a queue in advance and the device only writes evidence — is decided by stage 5. `engine` needs the same discipline to repeat a choice on any build. So the recommendation of Т-10 is supplemented, while the question itself is open (R11), and from the engine's side too: `engine` gives the same choice on any build of one version, WebAssembly too, and it does not need `detmath` — it computes no exp or log (R10 engine; R107 spec).

**How this differs from standalone.** standalone's ratings package takes `math` (sources 1.4), and its test vectors are compared with a tolerance of 1e-4: "The implementation computes in `float64` and must match the four decimals printed here to within 1e-4" (SPEC 2.8 standalone). The harness `sim`, which inherits it, says so outright: the guard's bands are recorded on amd64 at GOAMD64 v1 and on a processor with FMA, and on arm64 and at v3 «последние цифры чисел могут разойтись» (architecture 2.3 sim). Here the same bits are a requirement: the state is a fold of the log with thresholds (section 3.5).

### 3.7. Types and the encoding

Decided by R10. The types are its own and simple, not from `module-kit` and not from `spec` (R28 spec; section 3.6), and `engine` takes the graph's types from here. The names below are the Go names; the comments of the code and the schema are in English (RULES.md, Hard rules).

**Packages** of the Go module `github.com/altnovatrail/learner-model`:
- the root — the types and the operations;
- `detmath` — exp, log, σ, the power and the retention (section 3.6);
- `internal/wire` — protobuf's binary form, on which the root's codec stands: `codec.go`, `codec_check.go` and `codec_read.go`;
- `proto/novatrail/learnermodel/v1/` — the schema of evidence and state, the package `novatrail.learnermodel.v1`, as with `spec`'s contracts (R83 spec).

**The graph** — per module. `api` and `sim` build it from the part of the module package for the server (R62 spec) by a function of the library, and `engine` gets the same values (R28 spec):

| Field | What is in it | Who reads it |
|---|---|---|
| module | the module's full ID | all |
| bands | in the order of the module manifest's `bands`: the place — a number from 1, the ID and the grades | the scale and the starting point (section 3.1), the band by grade (section 3.5) |
| skills | below | below |
| retired skills | `retired_skills`: a skill's ID → `replaced_by`, and `[]` too | transfer (section 3.5) |

A skill has:

| Field | What is in it | Who reads it |
|---|---|---|
| ID | the full ID of a module's skill (concepts 3 spec) | all |
| assessability | `verdict`, `rubric` or `none` | the model — only `verdict` has a state (section 3.4); the engine — where it adapts to the learner |
| `teaches` | the IDs of anchor skills, `core.<name>` | knowing an anchor skill (section 3.4) |
| bands where it is taught | the places of the bands | "mastered" (section 3.1), mastery (section 3.4) |
| prerequisites | the full IDs of the module's skills and the IDs of anchor skills | the engine (concepts 4 spec); the model does not read them |
| tier | `skill_tiers` of the part for the server | the engine; the model does not read it |

**The bands where a skill is taught** are the bands of the tasks where it is a primary skill of the task or of an item (section 3.1). They are derived by a function of the library that builds the graph from the module's skills and the descriptions of its tasks: the rule is one, not a copy each in `api` and `sim`.

**The task description** — one for evidence and for the engine's candidate task (section 3.5):
- the task's full ID, and for a template the template's ID (R60 spec);
- the band — the place and the ID; the difficulty from 1 to 5;
- the primary skills — by full IDs, the supporting ones — by full IDs or `core.<name>`;
- the provenance of the key and of the skill tags (section 3.3);
- the items by place — their own skills, if they have any, and the slots in order: the check kind and, for a choice, the number of options.

Lessons, the pool, lesson runs, the summary of issues and the choice the engine keeps in its own types around the task description, and for `engine`'s rules the graph and the task description cover everything (architecture 3.8 spec; R09 engine; R107 spec). The age category is a filter of `api`, not an input of the engine (R107 spec).

**The input "graphs"** — module → the graph of its current version, or "no current version" (section 3.5).

**Evidence** carries the fields of the tables of section 3.5. **The state** — those of sections 3.1–3.3: for a module, θ, n, the start band, the answers of the trial series with the weight u, whether the last evidence was still in the series, and the traps that have no skills with a state, with the count and the time of the last one (section 3.3); for a skill, δ, m, the number of pieces of evidence, the time of the last practice, τ, the traps with the count and the time of the last one, the mastery band or "none", and the failures in a row. Estimates, weights and the sums of weights n and m are `float64`; the number of pieces of evidence and the time are `int64`; places and difficulty are integers. The grade is a value where "no grade" is a separate case, not 0.

**Operations** (section 3.5):

| Operation | Input | Output |
|---|---|---|
| `BandByGrade` — the band by grade | the grade, the bands of the attempt's version | the place and the ID of the start band |
| `Append` — append | the state, the evidence, the graphs | the state, and the evidence with the start band that the appending took |
| `LevelCorrection` — an adult's correction | the module, the band set, the start band, the time | the evidence of the correction, kind 3 (section 3.5; R45) |
| `Recompute` — recompute | the log in order, the graphs | the state |
| `Read` — read | the state, the graphs, the time t₀ in seconds, the grade, the engine's modules | for a module, whether the trial series has ended and whether the last evidence was still in it (R104 spec), the band the level stands in (R45) and the traps that have no skills with a state; for a skill with practice, the number of pieces of evidence, the time of the last practice, θ + δ, v and the retention at t₀, the traps with the count and the time of the last one; for a mastered one, also the mastery band, the cautious probability with the retention and the review date in seconds or "already"; the starting point of modules without evidence; anchor skills with the state, the reason and the skills that teach them (sections 3.1–3.4) |
| `Readout.Task` — the P of a candidate task | the task description | the P of each item at t₀ (sections 3.1 and 3.2) |
| codec | evidence or state | bytes, and back |

An error is a value, without a panic: a module of evidence with no entry in the input "graphs" (section 3.5), corrupted bytes, NaN or ∞ at encoding (section 3.6).

**The encoding** — protobuf's binary form, proto3; the codec is its own, on the standard library (section 3.6):
- the schema — `.proto` in `proto/novatrail/learnermodel/v1/`, with comments in English at each element. It is also the description of the form for the export: `api` puts its files into the export beside the JSON of the log and the state (R135 spec; section 3.5), so a comment explains the field to the reader of the export, without code;
- each scalar field has `optional`: "not written" differs from 0 and `false` (section 3.5). Without it a scalar field of proto3 has no presence;
- `double` — eight bytes of IEEE in little-endian, as is: −0 stays −0, and the machine's byte order changes nothing;
- the names follow the rules of `spec`'s contracts (R83 spec): the terms of concepts 15 spec in English in snake_case, a repeated field in the plural, and the zero value of an enum is `<NAME>_UNSPECIFIED`, "not specified". An unknown value of an enum is read cautiously, as section 3.5 says;
- the writer writes the fields in ascending order of number, repeated scalars as `packed`, the elements of repeated fields in canonical order: items by place, skills, traps and modules by ID byte by byte — and writes no unknown fields. There are no `map<>` fields in the schema. One value — one set of bytes;
- the reader skips an unknown field number by its wire type — except groups, 3 and 4, which proto3 does not write: they are corruption — and rejects corrupted bytes;
- a field number is neither reused nor removed: the form only grows (section 3.5);
- evidence and state each carry a kind number.

**What buf checks it with** — buf of `spec`'s container, at the version `BUF_VERSION` of `.devcontainer/Dockerfile` (R83 spec), by steps of its own in `ci.yml`:
- `buf lint` — the categories `STANDARD` and `COMMENTS`, with `disallow_comment_ignores: true`, as for `spec`'s contracts (R83 spec): a comment at each element, and a rule cannot be turned off in the file;
- `buf breaking` against the previous tag, the category `FILE`: a removed or renumbered field turns the check red;
- `buf convert` — the golden bytes into JSON by the schema and back: the codec and the schema speak of one thing, and people read the JSON of the same value. The names of the fields in it are as in the schema, `#format=json,use_proto_names=true`: for `spec`'s contracts "A field's name in JSON is its proto name, snake_case" (R83 spec).

An experiment confirmed this: `buf convert` with `use_proto_names=true` turned evidence encoded by hand into JSON and back into the same bytes, and `buf breaking` caught a removed field (measurements 3).

**There is no protobuf runtime.** `google.golang.org/protobuf` is a third-party module (section 3.6), and its deterministic output is promised only within one binary: "It is not guaranteed to remain stable over time" (sources 4.8). Recomputation, though, compares the bytes of the state between builds and versions of Go.

**In code** (R10, R16, R28, R29, R45):
- the types are at the root of the module: the graph `Graph` and `NewGraph`, which derives the bands where a skill is taught; the input "graphs" `Graphs`; the task description `Task`; the evidence `Evidence`, where the task description is `Task`, and what happened with the items in the attempt is `Items` by place; the state `State`; the types of the reading's output. A field that may be absent by its meaning is a value and a flag beside it, without pointers. The kinds are `EvidenceKind = 3` and `StateKind = 2`; what the later kinds added is `CheckKindTime` and `Evidence.LeadingQuestionOnRequest` (R28), `Evidence.LevelSetByAdult`, `ModuleState.FailuresInARow` and `TrialEndedByAdult`, `ModuleReadout.LevelBand` and `HasLevelBand`, and the operation `LevelCorrection(module, band, start, at)`, which gives the evidence of a correction (R45). The rules of R157 spec are `rules.go`, and `startLevel` and `levelBand` are in `model.go`;
- the schema — `common.proto`, `evidence.proto` and `state.proto` in `proto/novatrail/learnermodel/v1/`, `buf.yaml` at the root. What an earlier kind lacks: `CHECK_KIND_TIME = 5` in `common.proto`, which needs no kind (section 3.5); `leading_question_on_request = 20` of evidence of kind 2 and `level_set_by_adult = 21` of kind 3 in `evidence.proto`; `failures_in_a_row = 9` and `trial_ended_by_adult = 10` of `ModuleState` of a state of kind 2 in `state.proto`. Field 8 of `ModuleState`, `last_evidence_in_trial`, is in every kind of state, always written and required;
- the comments of the `.proto` are written for the reader of the export: `api` puts the three files of the version that computed beside the JSON of the log and the state, into `api/proto/novatrail/learnermodel/v1/` (R135 spec; section 3.5). Each field says what it means, in what units, what "not written" is for it and which fields the model does not read; the fields of help say that an answer with them does not set "mastered", and field 21 what a correction is;
- the codec — `EncodeEvidence`, `DecodeEvidence`, `EncodeState` and `DecodeState` over `internal/wire`. The writer and the reader are strict: the writer does not write values without a flag, strings not in UTF-8, or duplicates, and the reader counts as corruption a field of the wrong wire type, a non-repeated field twice, the absence of a field that every kind writes, and elements not in canonical order. A field of a later kind is written from its kind on, required there and refused in an earlier kind, so the bytes of an earlier kind from `api`'s log stay the same bytes (section 3.5). What the reader accepted, the writer gives back as the same bytes;
- in the step, `knownChecks` knows `time`, and an item with it gives an answer, an answer of the trial series and P in `Readout.Task`; `helped` counts a leading question as help too (R28);
- `Read(s, g, t₀, grade, engine modules)` and `Readout.Task(task)` — `readout.go`. For a module the output carries `TrialOver` — n ≥ 5 − ε —, `LastInTrial` from the state, the band the level stands in and the module's traps: which traps are frequent, `engine` decides (section 3.3), and a trap without skills with a state is counted only for the module; for a skill with practice — the number of pieces of evidence, the time of the last practice, θ + δ, v, R at t₀ and the traps, and for a mastered one also the mastery band `BandRef`, R · q at t₀ and the review date, by the functions of section 3.2. `Readout` holds a copy of the state, the graphs, t₀ and the grade for `Task`;
- `Read`, for each module whose current graph it checked, gives once the R at t₀ of each skill of the graph with a verdict — by the step's function `retention` — and the starting point of a module without a state, from its reading. `Task` only reads them, so `Readout` does not change after `Read`, and `Task` can be called from many calls at once (section 3.6, rule 7; R29);
- `Readout.Task` checks the candidate as the codec checks the task description of evidence — by the functions `checkTaskContent` and `checkTaskTags` and the places of the items in order, with the same rejections —, and it also rejects a module without a current version and a module without a state whose current version has no bands: a task of such a module has nothing to be in. The P of an item is the function `chanceOf`, the same as the step's, with transfer by `retired_skills`; for a module without a state — from the starting point by grade, and both the step and the reading take a new module by the function `newModule`, with θ = θ₀: the series computes θ from θ₀ anyway, so this does not change the bits of the state. The provenance of the key and of the tags does not change P. An item without slots or with an unknown check kind is without a verdict and without P. The graphs of the modules that `Read` checked, `Task` does not check again;
- a usual call of `Readout.Task` allocates memory once, for the slice of the answer (R29). Tags and items that are already in order it reads as they are — `checkTaskTags`, `sortedView` —, it takes the module as a substring of the ID, and `transfer` gives a skill of the graph as a slice of the graph itself. Only the step gathers the shares of transfer: `chanceOf` gets a slice for them from the step and `nil` from the reading;
- the release file — `learner-model-proto-<version>.tar.gz`: a folder of the same name with `LICENSE`, `buf.yaml` and all the `.proto` from `proto/`, in a tar that every build writes the same way, and gzip without a name and a time. The recipe `release-files` builds it, and the shared release signs it (section 3.6). It is the description of the form for the export, and `just -f spec/justfile release-verify learner-model <tag> learner-model-proto-<tag>.tar.gz` checks it (R77 spec; R16);
- the checks:
  - the recipe `ci-proto` and the job "Proto": `buf lint`, `buf breaking` against the last tag `v*`, the category `FILE` — skipped, with a word, where the clone has no tag `v*` —, and `buf convert` of the golden bytes `testdata/wire/` — the evidence, of the latest kind with fields 20 and 21, the state and the correction `correction.binpb` — into JSON beside them and back into the same bytes;
  - `ci-release` in the job "Proto" builds the release file twice and compares the bytes, holds its content equal to the `.proto` that git knows, `buf.yaml` and `LICENSE`, and builds the buf schema from the unpacked archive;
  - the golden logs `testdata/fold/` and their reading `testdata/read/` — each log read at three times, with the bits of each number — in all the runs of section 3.6. The logs `anchors`, `mixed`, `olympiad`, `transfer` and `workbook` hold evidence of kind 1, as a version without kind 2 wrote it; `format04` — evidence of kind 1, then of kind 2: the slots `time`, the trap `twelve-hours-off` and a leading question after which the skill is not mastered, while the same page without it masters the skill; `break` and `rules` hold evidence of kind 3: `break` — a break of the module, where the level forgets (section 3.2), `rules` — the three rules of R157 spec. Standalone's SPEC 2.9 starts its learner of grade 3 at 2.5, which is the start of grade 5 here, and SPEC 2.7's failures are put right on a rung, so that the guard keeps out of its numbers;
  - the property "the appended equals the recomputed" runs `time` and field 20 on evidence of kinds 1 and 3;
  - the reading — the example of section 3.2 through `Read` and `Readout.Task`, the P of the reading equals the P of the step bit for bit, in the series and after it, the rejections and the property on random logs;
  - `rules_test.go`: the start below the band of the grade, the band the level stands in at its edges, the correction and its repeat, in the band already, in the series and as the first evidence, what a correction must be and a band above the highest, the guard on single answers, on a choice of two and on pages, its floor, and its state in both kinds of state;
  - the reference `TaskReference` in `task_reference_internal_test.go` — `Task` of the version `v0.1.5`: its checks word for word, and P is `chanceOf` as the step calls it, with R anew and the shares of transfer. `TestTaskIsReference` compares with it P to the bit and the rejections to the word on the prefixes of the golden logs, on the random logs of the property and on the benchmark's pool, with the variants of each task; `TestTaskAllocatesOnce` holds one allocation per call on the pool, and `TestTaskConcurrently` holds calls from several goroutines at once under `-race`. The benchmark is `BenchmarkReadoutTask`, a synthetic module of the form of `olympiad-math` 0.1.1 in four states (measurements 2).

## 4. Sources

The records of external facts are in [sources.md](sources.md), by topic: standalone, learner models and the assessment of answers, forgetting and repetition, Go, CI, GitHub and tools, other NovaTrail repositories; the first line of a record names the subsections it was read for, and a claim of this document cites it, `sources 1.1` (R35). This section stays as a pointer, because the decision log and the archives cite it. The results of measurements are in [measurements.md](measurements.md).

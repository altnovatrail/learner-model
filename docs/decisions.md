# learner-model decision log

> Why learner-model is built the way it is. Each entry: what was decided, why, what was rejected, and in its tail "— *Consequence:*" the files, by their paths, that a change of it edits.
>
> **Rule:** do not propose again what was rejected without new facts. If a decision has to change, ask the author and rewrite its entry in its place, and every entry the change touches: an entry is its text in force — what was decided, why, what was rejected and what a change of it edits —, git keeps the old text, and the commit says why. An entry's number stands in its heading alone: an entry cites no other entry and no section, says in words what it rests on, and carries no date of who decided it and no amendment; a fact that no run of ours checks cites its one record of [sources.md](sources.md), and a measurement its record of [measurements.md](measurements.md). An entry that still has dated amendments is rewritten whole, its amendments folded in, when a task touches it.
>
> The repository was set up on 2026-10-19 from the template of `spec` v0.1.0; a move to a tag of `spec` is its commit, with no entry here; an entry of such a move written before this rule keeps its number as one sentence, since the numbers go without gaps, and holds no decision. How it comes to differ from the template — the tools in the devcontainer, the rules in RUN.md — is recorded here as a decision. Next number — R49.
>
> Open questions that concern several repositories — О-, П-, Т- — are in the registers of `spec`, in `spec/docs/product.md`, `spec/docs/privacy.md` and `spec/docs/repositories.md`. What is still to decide in the component itself is in [architecture.md](architecture.md).

## The tag of `spec`

**R01. The container of `learner-model` is `spec`'s: `.devcontainer/` is what `spec`'s template builds for the kind `go` at the pinned tag, with no difference of its own.** How a move to a tag compares it with the template and raises its versions is in `.claude/rules/ci.md`.

— *Why:*
- a difference from `spec`'s container would have to be recorded here as a decision, with its reason, and none has one;
- with no difference of its own, whatever a move finds between `.devcontainer/` and the recipe's output is a fault to fix, not a choice to weigh.

— *Rejected:*
- carrying the template's changes of the container over selectively, or moving the submodule without them: the container would differ from `spec`'s with no decision holding the difference.

— *Consequence:* `.devcontainer/`; `.claude/rules/ci.md`; CLAUDE.md.

**R02. The move to `spec` v0.6.3 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The learner model

**R03. The learner's level θ is per module, and standalone's correction for a topic belongs to the skill: the estimate of mastery of a skill is θ + δ. The scale of the bands, the trial series, the Elo step and the cautious "mastered" are standalone's; the start is a band's width below the band of the grade, guessing is by the item's answer slots, and the threshold "mastered" is a constant of `learner-model`.** The model — the state, the scale, the start band, P, the step, the uncertainty, "mastered", the operations and what `Append` rejects — is in [architecture.md](architecture.md), its sources in [sources.md](sources.md).

— *Why:*
- **The level per module.** Each module has its own scale of bands, and the format does not compare the hardness of different modules' tasks: "numbers are the model's business, not the format's", as `spec` decided. The band "grade 3" of school maths and the band "3–4" of olympiad maths are not one hardness. `spec` names standalone's model with a level per module as the starting model; the shared taxonomy links the two maths, and the level need not.
- **The correction per skill.** Standalone's topics become skills, and in `spec`'s vocabulary "mastered" and the state belong to the skill; a hierarchical model of adaptive practice corrects "the global skill by the concept skill" (sources 2.4).
- **Standalone's numbers.** Olympiad maths, a choice of five, then computes exactly as standalone, and its test vectors hold; its harness measured its rule "mastered" to give far fewer false "mastered" in fewer answers than its former rule ([standalone@c6cbeb2](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md)). Standalone calls the step 2.5 and the spread of the series a guess, and the harness checks them.
- **The start band in the log, by place.** The grade lies in the profile, not in the log: a recomputation that read it would diverge from the appended state after a change of grade. Each piece of evidence carries the start band, since a withdrawal for `wrong_key` may cut off the first; by place, since the recomputation opens no packages of earlier versions; and the append takes the band from the state, so two first attempts started before and after a change of grade do not set the appended and the recomputed apart.
- **The start a band below, an adult's correction and the guard** are the rules `spec`'s card names in place of a calibration of the start, which waits for real learners' data.
- **The first band without a grade.** The trial series moves the learner by a band either way ([standalone@c6cbeb2](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md)), and a weak learner started from the middle band would more often meet tasks beyond their strength.
- **Guessing by answer slots.** The floor is Birnbaum's: whoever does not know the answer "will guess, and, with probability c_g, will guess correctly" (sources 2.2), and with k options the probability is 1/k + (1 − 1/k) · σ (sources 2.3). An item is correct only when all its answer slots are, so its guessing is their product. A free answer shows no options, and `spec`'s vocabulary calls its guessing "about zero": 0 stands until the harness shows otherwise.
- **The threshold without guessing.** "Mastered" means "solves", not "guesses", and 0.71875 gives olympiad maths exactly standalone's rule.
- **The threshold in `learner-model`.** "Mastered" depends on the history of answers, and `spec`'s vocabulary keeps it in the skill's state; a number of `engine` among the inputs would make the cache depend on what its mark does not name.
- **The item's answer as the unit.** P is computed per item, which has its own answer slots and skills; what a page of many items weighs is the weights of evidence.
- **The zero state and the refusals.** `api` holds the zero state while there is no cache, so the fold starts from it; a value the step relies on that no log could give is the caller's error, which a recomputation must not hide.

— *Rejected:*
- a level per subject area by the code `subject`, or a new module's start from the level in a module of the same subject: both need one scale of hardness across modules and authors, with a term in `spec`'s vocabulary and a field of the next format version;
- the correction at the tier: the tier is computed from the graph and shifts with it, and "mastered" is the skill's anyway; a correction at the tier and the skill together: standalone's harness did not check such a model;
- in the first slice, an initial estimate by age — kept only as a category, and the bands do not name it —, by self-assessment — for language modules, none in the slice — or by a short check — a choice of tasks, `engine`'s rule;
- a step of bands set by the module: a field of the next format version; bands without overlap, with step 5: standalone's scale, its vectors and its harness would no longer hold;
- without a grade, the middle band: above, "The first band without a grade";
- the start band only in the first piece of evidence, the band by its ID instead of its place, the starting point by the grade from the profile at each recomputation: above, "The start band in the log";
- the threshold as an input from `engine`, or "mastered" that `engine` decides: above, "The threshold in `learner-model`";
- the threshold σ ≥ 0.775: stricter than standalone, and its harness does not carry over; P ≥ 0.775 with the guessing of the skill's tasks: "mastered" would come earlier on a choice of two than on a number;
- a small constant of guessing for a free answer, or a number of its own for `predicate`: no data;
- δ̄ as the minimum of the primary skills' δ — it jumps from skill to skill — or a product of probabilities over them — it drops P with the number of skills at the author's same difficulty;
- updating β from answers: the author sets a task's hardness;
- an attempt as one answer for the counters: it would decide the weights of evidence in advance;
- an I that depends on each answer's guessing: the state would accumulate sums, and the harness did not check it; the price is an understated uncertainty on a choice of two;
- a constructor of the empty state, or requiring kind 1 of it: an extra step for every caller, of no use to the fold;
- a module state with only the start band, from evidence without answers: the band would be pinned by evidence the model did not read;
- silently passing the inputs `Append` rejects: the recomputation would hide `api`'s error.

— *Consequence:* [architecture.md](architecture.md); `CLAUDE.md`; `README.md`; `band.go`, `model.go`, `step.go`, `fold.go`, `rules.go` and their tests; `standalone_test.go`; `testdata/fold/`.

**R04. Forgetting is the retention of a skill whose stability grows with spaced practice: P of an item is multiplied by the retention of its primary skills, which falls by a power law with the time since their last practice, and after a pause of the module of 14 days or more its n falls to 20 before the step, so the step of θ widens. The law and the constants are one for all modules. Forgetting does not remove "mastered" but gives a review date; `engine` gets the estimate for the time of the choice and the date; the current time does not enter the cache.** The law, the weight w of the step, the stability, the date and the forgetting of the level are in [architecture.md](architecture.md), the sources in [sources.md](sources.md). The repetition rule — when to give a skill with a date, how many repetitions a day, how they share the queue — is `engine`'s.

— *Why:*
- **Retention with stability** is how spaced repetition works: memory fades until practice, and practice grows its strength. HLR's half-life is "strength in the learner's long-term memory" (sources 3.2); in FSRS "The smaller the value of $R$, the larger the $SInc$ value" ([awesome-fsrs@7bcf7dc](https://github.com/open-spaced-repetition/awesome-fsrs/wiki/The-Algorithm/7bcf7dc343fc81d0fc55db2705c8b8fe85c2daa5)). A failed repetition halves τ, its P is low, and it soon comes up again.
- **Maths is forgotten, and spacing works in it.** By a meta-analysis, "there is substantial skill loss with nonpractice or nonuse", most for "cognitive, artificial, and accuracy-based tasks" (sources 3.6); practice in two sessions "virtually doubled four-week test scores", while more problems per session "had virtually no effect" (sources 3.7). So τ grows from a correct answer after a pause, not from answers in a row, and changes once per skill per piece of evidence, as ACT-R's mechanism "limits long-term benefits from further practice at higher levels of activation" (sources 3.8).
- **A power law** described forgetting best of six laws (sources 3.5), FSRS moved to it from an exponential one in v4 ([awesome-fsrs@7bcf7dc](https://github.com/open-spaced-repetition/awesome-fsrs/wiki/The-Algorithm/7bcf7dc343fc81d0fc55db2705c8b8fe85c2daa5)), and Ebbinghaus found forgetting "very rapid at the beginning of the process and very slow at the end" (sources 3.4).
- **The numbers.** C and F are FSRS-4.5's ([awesome-fsrs@7bcf7dc](https://github.com/open-spaced-repetition/awesome-fsrs/wiki/The-Algorithm/7bcf7dc343fc81d0fc55db2705c8b8fe85c2daa5)); τ₀ = 1 day is SM-2's first interval (sources 3.3) and the Leitner system's first box (sources 3.2); g = 15 grows τ 2.5 times at R = 0.9, SM-2's initial E-Factor (sources 3.3); halving is one Leitner box down. No data stand under them (Д6).
- **One for all modules.** No data would let a module's author choose a speed of forgetting; a field of the format must be known to the engine and the harness before the code; constants by the subject's code put the subject into the core, against the criterion of the first slice. Language may need others: forgetting matters most for "declarative knowledge" (sources 2.3) — Д6.
- **"Mastered" is not removed by time**, so it stays a function of the log: the cache, knowing an anchor skill, the engine's prerequisites and the skill map do not change with the calendar alone. Forgetting shows by a wrong answer — two failures in a row remove "mastered", and it does not move down, as in standalone ([standalone@c6cbeb2](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md)) — and by P, the cautious probability with retention and the date.
- **The date and the estimate together.** The date is by the same threshold as "mastered", so the threshold is one, in one library; P with retention lets the corridor give a first task within the learner's strength after a long break.
- **Time bypasses the cache.** The cache is only what the log, the graphs and the version set, and τ and t_last come from the times of practice in the log, so the appended matches the recomputed; `api` keeps the current time out of the cache and stores a choice's time for its repeat ([api@d896943](https://github.com/altnovatrail/api/blob/d896943/docs/decisions.md)).
- **The time of the last action.** `api` closes an abandoned attempt by a job after 7 days; by the time of closing, the clock would start on the job's day, not the day of the work.
- **τ by the share of correct.** At S̄ = 1 and 0 it gives exactly the rules for all and none correct.
- **The weight w.** S − P is the derivative of the answer's log-likelihood brought to Elo, and w · (S − P) is the same step with retention: a wrong answer after a pause moves θ and δ little — forgetting explains it —, a correct one almost as without forgetting. Without w a correct answer after a long pause would move θ and δ several times more, and θ would inflate after a break, as the example in [architecture.md](architecture.md) shows. The code's form of w avoids dividing small by small, and its branch at R̄ = 1 keeps the bits of evidence without a pause.
- **Whole seconds, saturation and a date without a power.** The model gives the same bits for any state and no NaN or ∞; a power would pass Δt* through three roundings and could land it a second before the threshold.
- **The date no earlier than τ₀.** A repetition on the same day barely grows τ, and a barely mastered skill would come up every few minutes.
- **The trial series without retention.** Its answers lie in the module's state without time, and the skills it catches are almost all new.
- **g = 15 and halving τ stay.** In `sim`'s first judgement g = 10, the multiplier of τ 0.35 and both together were better than the base by the estimate, but none held every "no worse" constraint, and all three caught up more slowly after a break (measurements 12.3).
- **The forgetting of the level.** Catching up after a break was slow because the step with retention met a θ that is not forgotten, with its step at the floor 0.05 past n = 60. N = 20 and d₀ = 14 hold all of `sim`'s constraints, beat the base and catch up after a break faster, at the price of a slightly later "mastered" (measurements 14), and `sim`'s judgement chose them ([sim@7c32016](https://github.com/altnovatrail/sim/blob/7c32016/judgements/model-2/report.md)). They need no field: the module's time is derived from its skills' t_last, which each piece of evidence that moves a skill sets; and unlike a growth of v with time, the rule acts only with evidence.
- **The harness on maths.** The first slice has no language modules, and maths is forgotten too (sources 3.6); standalone's harness knew no time ([standalone@c6cbeb2](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/learners.md)).

— *Rejected:*
- a correction for time only at reading, with the step without forgetting: any answer would reset the clock, a failed repetition would barely move the estimate and the intervals would not grow — a hundred mastered skills would give about fifteen repetitions a day without end; Pelánek's correction is a short "temporal increase in memory activation", on facts of geography, not a law of forgetting (sources 3.9);
- the uncertainty v growing with time, as Glicko's RD (sources 3.10) and Math Garden's 1/30 a day (sources 2.3): θ + δ would not fade, and the corridor would not see forgetting; the skill's part of v is at most 1/(20 · I) ≈ 0.333, so the cautious level falls at most by about 0.58 logits, and a skill mastered with a larger margin would never return;
- FSRS whole — 17 weights in FSRS-4.5, 21 in FSRS-6 ([awesome-fsrs@7bcf7dc](https://github.com/open-spaced-repetition/awesome-fsrs/wiki/The-Algorithm/7bcf7dc343fc81d0fc55db2705c8b8fe85c2daa5)) — and HLR's weights learned from logs: fitted to cards and words, and there are no learners' data for them; ACT-R's activation model (sources 3.8): it needs the skill's whole history of practice with times and a fitted rate of decay;
- an exponential law, as in HLR and FSRS v1 and v3: a power law describes forgetting better (sources 3.5);
- g = 10, the multiplier of τ 0.35 and both together: above, "g = 15 and halving τ stay"; propose them again only with new facts about catching up after a break;
- for the forgetting of the level, N 10 or d₀ 7: they break `M5.caught` on `F-none` of school maths; N 40 or d₀ 28: they hold everything, but more weakly (measurements 14);
- the forgetting of the level together with the numbers `forget-set`, or with one base of halving τ, 0.35 or 0.4: no such set holds all the constraints (measurements 14);
- a separate step of θ after a pause by a new field of the state, leaving n and v alone: state kind 2 and a recomputation of `api`'s cache, while the price of n holds anyway;
- the forgetting of the level at reading, `Read` narrowing n by the time of reading: it is the growth of v with time, above, and was not measured;
- the step without forgetting — P without retention and w = 1: a wrong answer after a pause would lower θ and δ as if nothing were forgotten; S − P with P with retention and no w: a correct answer after a pause would inflate θ and δ — above, "The weight w"; w by its formula as it stands: it divides small by small at small σ;
- a speed of forgetting per module, by a field of the next format version, or by the subject's code: above, "One for all modules";
- "mastered for now" at reading, and a mastery band that falls as forgetting proceeds: "mastered" would change with the calendar, and the band does not move down;
- only the estimate for now for `engine`: a second threshold, "due", would live in the engine beside "mastered"; only the date: the corridor would not see forgetting;
- waiting for language modules, with no forgetting until them: no repetition of the mastered even in maths, and the harness's criterion would come after the first candidate, while a criterion is written before the candidates;
- R̄ as the product or the minimum of the primary skills' R: for the reasons that reject them for δ̄;
- a date without a lower bound, and a date computed with a power: above;
- τ by each item of the page: practice in a row would grow τ, against the experiment on maths problems (sources 3.7);
- rules of τ only for all correct and none correct, the mixed case by a provisional rule: the rule and the golden vectors would be rewritten later;
- all the reading in `Read`, with the retention, the cautious probability and the date computed there: the date would leave the rule of forgetting that sets it;
- a reset of τ to τ₀ after a wrong answer, as SM-2 starts the repetitions "from the beginning" (sources 3.3): one error is not yet forgetting, and only two failures in a row remove "mastered";
- silently raising a τ below τ₀ to τ₀: the recomputation would hide the error of whoever gave the state;
- the time of the attempt's closing as the time of practice, the current time in the cache, and retention in the trial series: above;
- naming the stability S, as FSRS does: S is the outcome of an answer.

— *Consequence:* [architecture.md](architecture.md); `CLAUDE.md`; `README.md`; `model.go`, `step.go`, `readout.go` and their tests; `testdata/fold/`; [run/DEBT.md](../run/DEBT.md).

## The tag of `spec` v0.6.6

**R05. The move to `spec` v0.6.6 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## Weights of evidence

**R06. The model's step takes the first verdict of each item, while help — a rung of the ladder before it, a technique on request, a sketch as a hint, a solution already seen — only keeps "mastered" from being set. Each of the N answers of one attempt weighs 1/(1 + (N − 1)/3), one piece of evidence gives one step, and "mastered" is decided by pieces of evidence. The item's primary skills share the step, after the transfer by `retired_skills`, 1/J each without a major version; the supporting ones get nothing. A key from the author, a solver or an adult, with tags not from AI, weighs 1; any other piece of evidence weighs 0 and moves nothing.** The outcomes, help, the series, the failure, the transfer, provenance, traps and the probe are in [architecture.md](architecture.md), the sources in [sources.md](sources.md). Which traps are frequent, when to take a skill off pause and after how long to give a task again is `engine`'s.

— *Why:*
- **The first verdict.** The guessing floor c and P describe the first try; after a wrong answer and a hint the probability is other, and there are no data on it. Standalone records an answer once ([standalone@89986bc](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/decisions.md)), and KDD Cup 2010 scored models by "Correct First Attempt" (sources 2.5). The log carries the wrong answers before the first verdict, so a new version can read partial credit by recomputation.
- **Help only in "mastered".** In standalone a hint changes what is chosen next, "never the rating", and "mastered" needs an answer without a hint: "The hint would make the answer measure the hint" ([standalone@89986bc](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/SPEC.md)). `spec`'s rules of the outcome leave the request and "explain differently" out of the outcome too. A technique at a step is not help: the new profile goes "by steps", and there "mastered" would never come. A sketch is help for the whole attempt: one at its start would otherwise give "mastered" from the first submission.
- **A repeat of a solved task is help, not weight 0.** The solution reaches the device whenever an attempt closes by a correct answer or by showing it, so this is almost any repeat of a solved task, and in format v0 a workbook page is never a template. At weight 0 a repetition of a mastered skill on an exhausted pool would not refresh τ, and the engine would give the skill again and again; as help it refreshes τ but does not set "mastered", since the learner could remember the key.
- **The discount by attempt.** "Ten fractions in five minutes are not ten independent pieces of evidence", as `spec`'s vocabulary says. Counting one stem's answers as independent "will overestimate the precision with which examinee proficiency is measured", leading to "prematurely ending an examination" (sources 2.6) — here an early "mastered". 1/(1 + (N − 1)ρ) is the share of the effective sample size at an intra-class correlation ρ (sources 2.7); at ρ = 1/3 an attempt weighs at most three answers and a one-item task exactly one, so olympiad maths counts as before.
- **Over the whole attempt, not a set of skills.** The answers depend on the shared stem and session, not on the set of skills; by sets, items with {A} and {A, B} on one page would get no discount, and a page of four sets of five items would move θ like 8.6 answers.
- **One step per piece of evidence.** The order of a page's items changes nothing, as a skill's retention is one per piece of evidence, and one function of the piece of evidence makes the appended state match the recomputed one simply.
- **"Mastered" by pieces of evidence.** Two "not yet" in a row by place on the page would depend on where the author put the items, and scattered errors would never remove "mastered"; for a one-item task a failure is a wrong answer, and two in a row are standalone's ([standalone@c6cbeb2](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/decisions.md)).
- **The share 1/J, and m in full.** P holds the mean δ̄ of the primary skills, and 1/J is the derivative of the log-likelihood by one skill's δ — not a quote: the multidimensional extension of Elo keeps each task in one knowledge component (sources 2.8). m counts the answers where the skill is primary: at 1/J a skill that always stands in a pair would need twice as many tasks for "mastered". Each part of a split skill takes the former share in full because its δ goes exactly as the former skill's, and the item's δ̄ does not change with a major version.
- **Nothing to the supporting ones.** They are not in P, so an answer tells the model nothing about them, and in `spec`'s vocabulary an error does not lower them either.
- **Outcomes by items.** The solution shown is "not solved", as `spec` decided, like "I don't know" in standalone: "a child who has seen the solution has not solved the task" ([standalone@89986bc](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/decisions.md)). What the learner did not answer the model does not judge: "punishing that would punish a stray tap" (same). "Unchecked" and "unreadable" other than blank are no verdict, but the verdicts before "unchecked" came from the program, and at zero weight "the incorrect answers before leaving would be lost to the model", as `spec` reasoned.
- **τ by the share of correct.** At S̄ = 1 and 0 these are the rules of forgetting, one miss in ten does not halve τ, and there are no new constants.
- **Provenance.** A key derived by AI gives no verdict, and a paper task that AI tagged and solved "does not move the rating" (О-11). `spec`'s vocabulary counts an adult's key "when … the skills are not from AI", so tags from a solver or an adult weigh 1; an adult's key also needs a band and a difficulty on the module's scale, or P cannot be computed. A provenance this version does not know may be as unreliable as AI, and a field that changes the weight acts only in a version that knows it. Weight 0 moves nothing at all, or tags from AI would decide whose skill was practised. `spec`'s list of evidence leaves out an AI trap, so its weight is 0.
- **A probe without a lift.** Standalone's lift relied on the pause, a rule of choice, and read a probe off "where the child stands when it is answered" ([standalone@89986bc](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/SPEC.md)); here the pause is `engine`'s, and the model would have to read its state. A correct answer at a low P raises δ strongly anyway.
- **No place for the rubric.** An attempt of a task with a rubric gives no evidence (О-17), and the weight of feedback comes with the answer to О-17.

— *Rejected:*
- partial credit by the number of hints, for example S = 1/2^k: new constants without data, and the log keeps k for later;
- the item's result at the end of the attempt as S, with the outcome as a multiplier of the step: a correct answer after three hints would count as correct, only weaker;
- an answer with help at weight 1/2: one more constant without data; a technique at a step as help: in the mode "by steps" "mastered" would not come; a sketch as help for the items after it alone: above, "Help only in "mastered"";
- a request for a technique as an error of the first attempt, like a hint request in KDD Cup 2010 (sources 2.5): a technique names a way, not an answer, and `spec`'s rules of the outcome do not let the request into it;
- a discount of √N: it grows without bound; a series as one answer with the share of correct ones: a page of twenty items would weigh like one task; a discount by an identical set of skills: above, "Over the whole attempt";
- a step for each answer in the order of the items, with "not yet" in a row at a total weight of 2: the order of places on the page would mean something; counters by pieces of evidence, with the weight only in the step: a page of ten and one task would weigh the same for the uncertainty and "mastered";
- the whole step of δ to each primary skill: at J = 2 δ̄ would shift twice as much as for one answer; m by 1/J: above;
- a step only at a correct answer, τ or a time of practice to the supporting skills: their estimate would rise without measure, and "time to review" would be put off by an answer that says nothing about them;
- not reading "unchecked" at all: the program's verdicts would be lost; the solution shown without "not solved" for empty items: an adult who showed the solution at once would not move the model; "left unfinished" as "not solved", and "unchecked" or "unreadable" other than blank as an error: above, "Outcomes by items";
- a repeat of a solved task at weight 0 or 1/2: above, "A repeat of a solved task is help";
- τ that grows when all are correct and halves otherwise: one miss out of ten would halve it; τ by the majority: a jump at the boundary;
- an adult's key at weight 0 until stage 5: О-11 gives it a verdict; AI traps in the evidence at weight 1/2: `spec`'s list of evidence leaves them out, and it would become a row to `spec`;
- weight 0 that moves only the rating: tags from AI would decide whose skill was practised, and a skill could have practice without m; weight 1 for a provenance "not specified" or unknown: above, "Provenance";
- creating a skill's state for the sake of a trap: a skill has a state only with practice;
- a lift on a probe by a mark "probe" from `engine` in the evidence — a field that standalone rejected in its profile — or by a reachability threshold in the model: above, "A probe without a lift";
- a place in the weights for feedback by rubric at weight 0: there is no evidence, and no place is needed.

— *Consequence:* [architecture.md](architecture.md); `CLAUDE.md`; `README.md`; `step.go`, `model.go`, `rules.go` and their tests; `transfer_test.go`; `testdata/fold/`, `testdata/read/`; [run/DEBT.md](../run/DEBT.md).

## Knowing an anchor skill

**R07. The learner knows an anchor skill when "mastered" is set on at least one skill that teaches it, in any module whose graph came and on any band; "does not know" is given only by the engine's modules, and otherwise the anchor skill is "unknown". Each state comes with a reason and the skills that teach it and is computed at reading, with each module's two marks of the trial series; forgetting does not remove knowing, and skills with `rubric` and `none` have no state.** The rule, the reasons and what `Read` gives are in [architecture.md](architecture.md). What the engine does at "no evidence", when knowing goes away and at a cycle of prerequisites between modules through anchor skills is `engine`'s; knowing by skills with `rubric` waits for О-17.

— *Why:*
- A skill that teaches an anchor skill teaches it whole — "a learner who has the skill can do what the description of each says", as the module format says —, and to know a skill is "mastered" in `spec`'s vocabulary: one threshold, and "or" over the skills and modules.
- The bands of different modules are not comparable, the level is per module, and the format does not say which band a module that rests on the anchor skill expects.
- A module that the engine does not get gives no sessions, so its unmastered skill would hold the dependent skills back forever; and "mastered" derived by evidence stays true after a module is disconnected.
- The reason "no evidence" leaves to `engine` the case where connecting one more module would make it stricter; the grade decides only where the child starts.
- Time does not remove "mastered": otherwise after the holidays the skills of another module would close until a repetition in the first. The retention and review date of the skill that teaches show forgetting.
- `rubric` and `none` give no evidence: an attempt with a rubric gives none (О-17), and `none` is only ever supporting. "Unknown" holds no lesson back, and among skills without a verdict the engine follows the author's order.
- "Does not know" and "unknown" depend on the engine's modules — the connections and the profile —, not on the log, so they cannot be in the cache.
- `engine` needs the marks of the series, and the piece of evidence that ends the series clears its answers, so the mark "last piece of evidence still in the series" is stored; "the series has ended" is n ≥ 5 − ε.
- The format lets a module without tasks leave out `bands`, and `engine` rejects an engine's module without the model's reading: such a module is given with an empty start band, and its skills, taught on no band, cannot be mastered.
- The skill map's contract names the mastery band by ID, and `api` does not look it up itself.
- The P of reading is the P that the step would take for a piece of evidence at t₀ — one function, R̄ = 1 in the trial series, a new module at θ₀ —, and the provenance of the key and the tags does not change P.

— *Rejected:*
- "and" over the skills that teach it — each teaches the anchor skill whole; an own cautious probability with an own threshold — a second threshold beside "mastered";
- "mastered" only on the highest band where the skill is taught — the bands of modules are not comparable, and the format has no field "which band is needed";
- everything only by the engine's modules — after a mastered module is disconnected the anchor skill would become "unknown", and the engine could offer the module again; everything by the graphs that came, without the input of the engine's modules — a disconnected module with an unmastered skill would hold the dependent skills back forever;
- three states without a reason — with one more module the engine would become stricter; "unknown" while the skills that teach it have no answers — the lesson is open at first and closes with the first answers;
- knowing only while the review date of the skill that teaches it has not come — after a break the skills of another module would close until a repetition in the first;
- "does not know" for an anchor skill taught only by skills without a verdict — it would hold the dependent skills back always; a state with a time of practice for `rubric` and `none` — there is no evidence, and a supporting skill gets no practice either;
- knowing in the cache — it depends on what the cache's mark does not hold;
- knowing by answers where the anchor skill is a supporting skill of the task — a supporting skill gets no step, and a shift only upward at a correct answer would pull the estimate up without measure;
- rewriting the mark of the series by any piece of evidence of a module with a state — a piece of evidence at weight 0 or without answers would change the state for the first time; storing "the series has ended" — it is n ≥ 5 − ε already;
- an item's P in the trial series with retention at t₀ — reading and the step would diverge for skills practised long ago;
- `ErrInvalid` for an engine's module whose current version has no bands — one such module would fail reading for the whole learner; skipping it, like a module without a current version — the engine's choice would break; a `NewGraph` that rejects a graph without bands — a module of lessons alone could not be given to the engine; accepting a skill taught on a band that the graph does not have — it would pass for a skill that can be mastered.

— *Consequence:* [architecture.md](architecture.md); `CLAUDE.md`; `README.md`; `readout.go`, `step.go`, `state.go`, `fold.go`, `graph.go` and their tests; `proto/novatrail/learnermodel/v1/state.proto`; `testdata/read/`; [run/DEBT.md](../run/DEBT.md).

## Evidence and the log

**R08. A piece of evidence carries everything that the recomputation reads — the observations of the attempt, not the model's conclusions — and its form only grows: every version reads every form, skipping what it does not know. Beside an attempt the log holds one more fact, an adult's correction of a module's level, as evidence of kind 3. The state is a fold of the log in the order of `api`'s numbers: an append is the same step as the recomputation and gives the same bits on any build of one version of `learner-model`; there is no snapshot. The evidence of a retired skill goes to each of its replacements with the share in full, and a replacement derives "mastered" anew from it. The record of a choice also stores the grade, and a major version of a module keeps the places of the earlier bands.** What evidence carries, the transfer, the start band and the tolerance at the thresholds are in [architecture.md](architecture.md). How `engine` and `api` repeat a choice with an earlier version of `learner-model` is theirs, and the format of the export `spec`'s; the fields of stage 5 — the provenance of the band and the difficulty of a task from paper — wait for О-11.

— *Why:*
- The thresholds — "mastered", n < 5, m ≥ 5, failures in a row — make the state discontinuous: a difference in the last bit sets "mastered" or does not, and the states diverge forever, so a tolerance holds nothing. A fold matches by construction if the step is deterministic and the cache is kept without loss. While a rollout goes on, the earlier and the new image append under one mark, which does not carry the build ([api@d896943](https://github.com/altnovatrail/api/blob/d896943/docs/decisions.md)); the discipline of the code makes the same bits reachable on every build.
- The recomputation is needed only with a change of the cache's mark, when a snapshot under the earlier mark is of no use and the transfer does not apply to it; a snapshot is one more place where the state can diverge from the log; the fold is linear in the log (measurements 11).
- The choice stores the digests of the graphs, by which the engine's modules are restored too, and the time ([api@d896943](https://github.com/altnovatrail/api/blob/d896943/docs/decisions.md)), but not the grade: without it the first choice in a module without evidence, whose starting point goes by grade, will not repeat after a change of grade.
- The log is not rewritten, and the recomputation opens no package of an earlier version: a band's place is all it knows of the band, and while the places of the earlier bands keep, the shift 2.5 · (p − 1) of earlier evidence is the same. The band's ID beside its place is for a table of correspondence, should one ever come: it cannot be appended to the log after the fact.
- The module format says: "the evidence of a retired skill counts for each skill that replaced it". With the share in full the δ of each part of a split skill goes exactly as the former one's, an item's δ̄ and R̄ do not change, and a part's "mastered" comes from the same history. Skills merged in one item are one skill, as if recorded so. "Mastered" is a function of the evidence, so the recomputation derives it anew, and with it knowing the anchor skill the replacement teaches. An item left without skills moves θ: the level belongs to the module, and the task's band and difficulty remain.
- What a piece of evidence does not carry, no later version can read from the old ones. Observations — submissions, verdicts, the numbers of wrong answers — are read by the weights of evidence now and perhaps by partial credit later (Д7), while a conclusion of the weights would be read only by the rule that drew it. `spec`'s list of evidence is the agreement by which `api` builds a piece of evidence; beyond it and what the model reads, it carries only the module's content that later versions may need — the IDs of bands, the supporting skills, the places of the items for traps — and the wrong answers before the first "correct", for partial credit: data are kept to a minimum.
- `api` rolls back only forward (Т-19), so an earlier version recomputes a log that already holds newer forms. Nobody rewrites the log, so reading it while skipping the unknown loses nothing — unlike standalone's profile, which an older instance would rewrite ([standalone@89986bc](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/architecture/04-profile.md)).
- An adult's correction is a fact — the band the adult set and the time —, not a conclusion, and it carries no answer: an earlier version, skipping its field, reads it as evidence with no answer, which moves nothing.

— *Rejected:*
- a match with a tolerance — the thresholds make the state discontinuous; bit for bit only on one build, with the build in the cache's mark — each change of the image would recompute all the logs, and a repeat of a choice on another architecture could diverge;
- a snapshot every N pieces of evidence — it is built anew with each change of the mark and can diverge from the log; a snapshot at each choice — as many copies of the state as there are choices, against keeping data to a minimum;
- a repeat without the grade — the first choice in a new module will not repeat after a change of grade; the start band of the modules without evidence in the record of a choice — a number of the model in `api`'s record instead of the profile's data, one per module;
- a band's place without a rule in `spec` — the comparison of versions would learn of the rule only when some module changes its scale; a band by ID, not by place — the format does not promise that a band's ID will not be taken for another one;
- equal parts between the replacements — the parts would start below the skill they came from and need more answers before "mastered", though the learner learned nothing anew; the share in full but without "mastered" — knowing an anchor skill would go away until new answers; the contribution of each recorded skill separately — for a merged skill the count and practice of one item would double; "as if recorded under the replacements" literally — each part of a split skill would get half;
- walking a chain of retirements recursively — the graph of the current version names the replacements already flat, so the transfer is one lookup;
- a piece of evidence with only what the model reads — the facts of `spec`'s list without weight would have to be removed from `spec`, and partial credit by old evidence could not be introduced; everything the attempt stores — the log would repeat the attempt, with more data; the conclusions of the weights, as "blank" or "not solved", in place of observations — only the rule that drew them could read them; the segment of a lesson — it does not change the weights and has no reader;
- a field missing from an earlier form read as "no" or 0 — it is "not recorded", and the version that added the field decides how to read it;
- a ban on rolling out an earlier version after a new form — `api` would stay on a version with a found error; refusing to compute such a log — the learner would have no model until a new rollout;
- a start band that `api` takes from reading at the start of an attempt, without replacing it at the append — the first two attempts of one module on two devices, with a change of grade between them, would carry different bands;
- a step that reads the graphs of other modules — the cache would depend on modules outside the log;
- the time of practice only as a column of `api` — a piece of evidence carries everything that the recomputation reads.

— *Consequence:* [architecture.md](architecture.md); `CLAUDE.md`; `README.md`; `evidence.go`, `fold.go`, `step.go`, `graph.go`, `types.go` and their tests; `transfer_test.go`; `proto/novatrail/learnermodel/v1/`; `testdata/fold/`, `testdata/wire/`; [run/DEBT.md](../run/DEBT.md).

## The build for the device

**R09. The library gives the same bits on any build of one of its versions by the Go compiler `gc` — on any architecture, GOAMD64 and GOARM64 level, processor and Go version from the one `go.mod` names, WebAssembly too — and is built for the device as it is. The discipline of the code holds this: an explicit `float64()` on the result of each floating-point `*` and `/`, its own exp and log in the public package `detmath`, a port of the pure Go of `math`, only Go's standard library by a list, and the rest of the rules of code; CI checks it from the first code. By Т-10, compute on the server; how offline works is decided by stage 5.** The promise, the rules of code, the lists, the jobs of CI and `detmath` are in [architecture.md](architecture.md), the experiment in [measurements.md](measurements.md). How offline works is stage 5's, by Т-10; the discipline of the code of `engine` and the repeat of a choice on the device are `engine`'s; the targets that CI does not run are Д10.

— *Why:*
- IEEE 754 addition, subtraction, multiplication, division and square root are correctly rounded. Go adds one liberty — "An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements" — and "An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding" ([go_spec@go1.27.3](https://github.com/golang/go/blob/go1.27.3/doc/go_spec.html#L5091-L5096)). With `math`'s own assembly and the conversions that the specification leaves to the implementation, these are the three places where builds diverge, and the rules of code close all three.
- The port with `float64()` gave the same bits in six builds, while `math.Exp` on amd64 diverges from it and changes its bits with the processor even in one binary; one `float64()` left out in the experiment's own harness changed the retention at GOAMD64 v3 — so the rule of `float64()` holds in the tests too (measurements 1.1).
- The rule is for each `*` and `/`, not only for `a*b + c`: a division by a constant with an exact reciprocal becomes a multiplication ([generic.rules@go1.27.3](https://github.com/golang/go/blob/go1.27.3/src/cmd/compile/internal/ssa/_gen/generic.rules#L1331)), and fusion goes across statements and through inlining, also into the code of a consumer that `api` builds with its own flags. Only a check of the library's sources on `go/types` holds it.
- `-d=fmahash=n` turns off every fusion rule ([func.go@go1.27.3](https://github.com/golang/go/blob/go1.27.3/src/cmd/compile/internal/ssa/func.go#L848), [hashdebug.go@go1.27.3](https://github.com/golang/go/blob/go1.27.3/src/cmd/compile/internal/base/hashdebug.go#L64-L70)), so comparing two builds does not depend on the names of instructions and survives new rules and architectures (measurements 1.1).
- `math.Exp`, `math.Log` and `math.Pow` run their own assembly on several architectures, and `math.Exp` on amd64 chooses its path by the processor ([exp_asm.go](https://github.com/golang/go/blob/go1.27.3/src/math/exp_asm.go#L5), [log_asm.go](https://github.com/golang/go/blob/go1.27.3/src/math/log_asm.go#L5), [pow.go](https://github.com/golang/go/blob/go1.27.3/src/math/pow.go#L128), at go1.27.3). The port is the code that `math` runs on wasm and 386, with the same accuracy. Go's BSD-3-Clause and Sun's notice, "provided that this notice is preserved", are compatible with Apache-2.0, and the license check does not see copied code, so the notices stand in the port's file (sources 4.5).
- A third-party module is one more place where the bits or the bytes change with its version; the Go module of `spec` embeds the whole `format/` (`//go:embed all:format`) into every build, WASM too, while the library's types are its own; even `encoding/json` in Go 1.27 is built from v2 ([buildcfg/exp.go@go1.27.3](https://github.com/golang/go/blob/go1.27.3/src/internal/buildcfg/exp.go#L87)), so the bytes of the standard library's encoders are no agreement between Go versions.
- Without CI from the first code, a violation would hide in the last bit until what is appended diverged from what is recomputed during a rollout of `api`. GitHub's arm64 runners are standard and free for a public repository, and `learner-model` is public (sources 5.1).
- Т-10 is shared with `engine`, and offline has a path that Т-10 does not name: the server chooses the queue in advance, and the device only writes evidence. The js/wasm build is useful anyway: one more backend without fusion, where the port is checked against `math`.
- Within the clamps exp, σ and log are finite and normal, so no function of `detmath` gives NaN, ∞ or a subnormal number; outside them `math`'s own assembly on amd64 departs from the correct answer (measurements 1.2). The power's exact branches s = 0 and s = 1 make the rule of τ by the share of correct answers at S̄ = 1 and 0 coincide with the rules for all and none correct, which exp(1 · log a) need not.

— *Rejected:*
- `math.FMA` everywhere instead of separate operations — the same bits too, but on amd64 at v1 and in wasm it is software and slow, and the differential check would no longer tell an intended fusion from an accidental one;
- fixed point — rewriting the whole model, and its own exp and log would be needed anyway; a correctly rounded libm — there is none in the standard library; tables of values — the port of `math`'s pure Go with `float64()` already gives the same bits at `math`'s accuracy;
- `float64()` only at `a*b + c` — a division by a constant and fusion across statements and inlining escape it; a check by a list of fused instructions — it depends on their names, while comparing two builds survives new rules and architectures;
- `-d=fmahash=n` in the consumer's build — a debugging flag of the compiler, and the library does not set the build of `api` and `sim`;
- the build in the cache's mark — not needed: the same bits are reachable;
- third-party modules with a compatible license; the Go module of `spec`; the protobuf runtime — the codec is its own;
- `math.NaN` in the port and a flag of the package `flag` for writing the golden files — neither is on the lists: the port's NaN is the bits of `math.NaN` through `Float64frombits`, and `DETMATH_UPDATE=1` writes the files;
- checks of the device only from stage 5; a set without an arm64 runner; a minimum of the js/wasm build and the list of imports — a violation would hide in the last bit until a rollout of `api`;
- closing Т-10 now — it is shared with `engine`, and offline has a path it does not name; leaving its recommendation without the build for the device — the same bits would go unchecked until stage 5.

— *Consequence:* [architecture.md](architecture.md); `CLAUDE.md`; `README.md`; `go.mod`; `detmath/`; `.github/coderules.json`, `.github/workflows/ci.yml`; `justfile`; [run/DEBT.md](../run/DEBT.md).

## Types and encoding

**R10. The library's types are its own and simple. The graph is per module: the bands in order and the skills with what the model reads, and also with prerequisites and tier for `engine`; the bands where a skill is taught the library derives itself, in `NewGraph`, from the task descriptions. One task description serves evidence and a candidate task. Evidence and state are encoded in protobuf's binary form by a `.proto` schema of its own, by a strict codec of its own on the standard library, and buf checks the schema and the codec in CI.** The packages, the types, the operations, the schema and the rules of the encoding are in [architecture.md](architecture.md). What else the rules of `engine` need in the graph's types `engine` asks of `learner-model`; the format of the export is `spec`'s; the fields of stage 5 — the provenance of the band and the difficulty of a task from paper — wait for О-11.

— *Why:*
- The component's card in `spec` names prerequisites for the engine, and `engine`'s architecture reads a skill's tier, prerequisites and assessability: without them `api` and `sim` would build the graph twice, before `engine`'s detailing and after it. Names, descriptions and lessons the model's rules do not need, and the card keeps the types simple, not a copy of the format.
- The bands where a skill is taught are a rule of the model; derived apart by `api` and `sim`, one module's graph could come in two forms, and `spec` gave the types to one library so that the graph has one form.
- Protobuf's binary form writes a number as its eight IEEE bytes, so −0 and the last bit survive with no rules about text; an unknown field is skipped by its wire type, so an earlier version reads a new form for free; a `.proto` schema with comments is the open description of the form that the export needs; buf is in `spec`'s container — `buf breaking` holds "the form only grows", and `buf convert` gives the JSON of the same value (measurements 3).
- The protobuf runtime for Go is a third-party module, and its deterministic output holds only within one binary: "It is not guaranteed to remain stable over time" ([protobuf@v1.36.12](https://pkg.go.dev/google.golang.org/protobuf@v1.36.12/proto#MarshalOptions)); the recomputation compares the state's bytes between builds.
- Without `optional` a scalar field of proto3 has no presence, and "not written" would read as 0, while evidence tells them apart; for the same reason an item with its own skills is the flag `own_skills`: an empty list of its own primary skills could not otherwise be told from their absence.
- The state is copied by assignment, and an append shares nothing with the state it started from.
- The engine needs to know whether the module's last piece of evidence was still in the trial series, so as not to work through the error of the series' last answer, and the state did not otherwise remember it.
- The log is never rewritten: a value the writer dropped would vanish from it forever, and a repeat or a disorder that the writer never writes can only be corruption.
- `CheckKind` holds only the kinds with a verdict: an attempt of a task with a rubric gives no evidence (О-17), and a new kind is an added value.
- `spec` set the platform's rules for protobuf — the packages `novatrail.<part>.v1`, `STANDARD` and `COMMENTS`, a field's name in JSON as in the schema —; following them starts no second style, and the export's JSON writes the same names, as `buf convert` does with `use_proto_names=true` (measurements 3).

— *Rejected:*
- types from `module-kit` or `spec` — the library imports only Go's standard library;
- only what the model reads in the graph — prerequisites and tier would come in a new version after `engine`'s detailing; the whole server part of the graph — the card keeps the types simple;
- the bands where a skill is taught as a field that `api` and `sim` fill by their own copy of the rule;
- JSON by RFC 8785 — readable without tools, but numbers are text through `strconv`, several times more bytes, and "the form only grows" would need a test of its own; a binary form of its own — neither tools nor a standard, and the export would need its own text about the bytes;
- the protobuf runtime and generated code — a third-party module whose output is not stable between builds;
- proto3 without `optional` — "not written" is indistinguishable from 0; an item's own skills told by a non-empty list — an empty list could not be told from their absence;
- pointers in the Go API — an append would have to copy the state deeply;
- panics — an error is a value of one of three kinds that `errors.Is` tells apart, its message naming what and where;
- a package of its own `learnermodel.v1` and `STANDARD` alone — a second style of protobuf beside `spec`'s contracts;
- the elements of repeated fields in the order of iterating maps, and `map<>` fields — one state would give different bytes;
- a lenient writer that drops a value without its flag — the value would silently vanish from the log; a non-repeated field read twice as protobuf's "the last one wins" — the writer never writes it, so a repeat means corruption; rejecting a varint not in its shortest form or fields not in the order of their numbers — the value does not change by it;
- `CheckKind` values for kinds without a verdict — an attempt with a rubric gives no evidence;
- the last piece of evidence in the series without a new field, by keeping the series' answers until the next piece of evidence — those are answers after the end of the series, which `Append` rejects, and the field of the series' answers would mean two things; storing n up to the last piece of evidence — a number for the sake of one bit;
- `buf breaking` against the latest tag by plain text order — a pre-release `v1.0.0-rc.1` would come after `v1.0.0`; the order is by version with `versionsort.suffix=-`.

— *Consequence:* [architecture.md](architecture.md); `CLAUDE.md`; `README.md`; `types.go`, `evidence.go`, `state.go`, `readout.go`, `errors.go`, `graph.go`, `codec.go`, `codec_check.go`, `codec_read.go` and their tests; `internal/wire/`; `proto/`, `buf.yaml`; `testdata/wire/`; `justfile`; `.github/workflows/ci.yml`.

## The `spec` tag v0.9.3

**R11. The move to `spec` v0.9.3 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## Code before the gate of `sim`

**R12. The model's first version is the baseline, not a candidate: its code, CI and first signed release, `v0.1.0`, came before the first tag of `sim`, checked by its own tests, by the property "appended equals recomputed" and by standalone's test vectors. Every change of the model after it passes the gate of `sim`, a required check of the ruleset of `main`; a later patch with the same code and schema the gate checks like any commit.**

— *Why:*
- `sim` imports `learner-model` with `engine` and takes its version by `go.mod`, while the library's CI takes `sim` by tag: without a version of the model the harness has nothing to run, and "the model after `sim`" is a circle.
- The criterion is written before the first run of a candidate, not before the first line of code: the baseline is what the harness compares the first candidate with.
- standalone's model, the start, passed its own harness; its numbers, Д5–Д7, are what the gate checks first.

— *Rejected:*
- only the scaffold, CI, `detmath`, the types and the codec before the tag of `sim`, the model after: the circle remains, and the work stops halfway;
- the whole first version waits for `sim`: the same circle, longer;
- `sim` keeps a copy of standalone's model as the baseline: the model would live in two repositories, and the harness would measure `learner-model` by someone else's code.

— *Consequence:* `RUN.md`; `CLAUDE.md`; `README.md`; [architecture.md](architecture.md); `.github/workflows/ci.yml`.

## Scaffold and CI

**R13. The library is the Go module `github.com/altnovatrail/learner-model`, whose `go` line is the lowest Go version on which the same bits were checked and which has no `require`. CI calls the shared `go.yml` of `.github` as the job "Go", which runs the recipes `ci-lint`, `ci-test`, `ci-vuln` and `ci-licenses` of the `justfile` and sends the coverage to Codecov; CodeQL reads Go. Every reference to a shared workflow names one release of `.github`. `ci-lint` is `gofmt`, `go vet` and the gopls version check, with no other linter.** The recipes are in the `justfile`, the references in `.github/workflows/`, the threshold of `codecov/patch` in `.github/codecov.yml`, and what CI runs is in [architecture.md](architecture.md). The ruleset `main` requires "Go / lint", "Go / test", "Go / vuln", "Go / licenses", "CodeQL / analyze (go)", "Proto", the four checks of `same-bits` and "Gate", posted by GitHub Actions (`integration_id` 15368), and `codecov/patch`, posted by Codecov (254). A check goes into the ruleset only after its first green run, and `codecov/patch` only once its status has come on a pull request at least once.

— *Why:*
- without the module, the recipes and the job "Go" the code has nowhere to pass its checks, and tests and `codecov/patch` come with every change of code;
- with one release in every reference the files hold one version of the shared checks, and a raise changes one SHA;
- Codecov needs no turning on: its app is installed in the organization and turns the repository on with its first upload;
- it costs nothing: the repository is public, `spec` chose Codecov for public repositories at no cost, and GitHub's standard runners are free for them (sources 5.1).

— *Rejected:*
- raising only some references by hand, as `module-kit` once did: the files would hold different releases, and Dependabot's pull requests would have to be merged or closed one by one;
- staticcheck in `ci-lint`: the rules of code catch the rules of the same bits and SonarCloud the code smells, and one more tool outside Dependabot would be raised by hand with every Go version;
- a `go` line below the version the same bits were checked on, as `go 1.27.0` in `module-kit`: the promise would name versions on which nobody checked it;
- raising the `go` line with the container's Go, as to 1.27.4 for a fix of the standard library: the line is the least version on which the same bits are promised and checked, not the container's; a consumer builds with its own Go, and a vulnerability of the standard library is fixed by the Go of the build, not by the line;
- a badge before its data, as the coverage badge before the first line of code: it shows "unknown";
- requiring `codecov/patch` before its status first came: Codecov posts no status on a pull request without lines of the report, a required check that never comes holds every merge, and the ruleset's bypass list is empty.

— *Consequence:* `go.mod`, `doc.go`, `justfile`, `.gitignore`, `.github/`; README.md; [architecture.md](architecture.md); the ruleset `main`.

**R14. gopls stands in the devcontainer image at an exact version, built from its lock file `.devcontainer/gopls/`; the step and the lock file are those of the kind `go` of `spec`'s template, and `ci-lint` holds `GOPLS_VERSION` of `.devcontainer/Dockerfile` to the version in the lock file, since CI never builds the stage `devcontainer`.** The settings of the Go extension, in `.devcontainer/devcontainer.json`, keep it from updating its tools and from linting on save; an offer of the extension to install a missing tool, the debugger dlv say, comes without an exact version and is declined. Dependabot does not raise the lock file: it goes up with `GOPLS_VERSION` by a move to a newer tag of `spec`.

— *Why:*
- without gopls in the image the Go extension installs the latest gopls by itself, not an exact version, against the rule of exact versions;
- CI never builds the stage `devcontainer`, and `spec`'s `new-repo-check` looks at a copy only when the repository is created: without the check, a hand edit of `GOPLS_VERSION` or of the lock file, or a move that lost the step, would surface only at Rebuild Container; the template's comment of the gopls step says the same;
- gopls is BSD-3-Clause (sources 5.6) and a development tool: it ships nowhere and puts none of its code into what ships.

— *Rejected:*
- leaving gopls to the Go extension: it installs an unpinned version;
- only forbidding the extension to install tools: no language server at all;
- no check of its own, leaving the pair to `spec`'s tests and `new-repo-check`: nothing would hold the lock file and `GOPLS_VERSION` in this repository after its creation.

— *Consequence:* an edit of the `justfile` or `.devcontainer/Dockerfile` releases the next patch with the same code and schema, by the release's `paths`. — `.devcontainer/`, `justfile`, `.github/dependabot.yml`; CLAUDE.md; [architecture.md](architecture.md).

## The `spec` tag v0.9.5

**R15. The move to `spec` v0.9.5 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The first release

**R16. `learner-model` is released by the shared `release.yml` of `.github`, on the line its `RELEASE_LINE` names, when a commit of `main` with green CI changes a path of its `paths`: the code but its tests and `testdata/`, `go.mod`, the schema, the license, the `justfile`, `.devcontainer/Dockerfile` whole and `release.yml` itself. Go consumers take the code as the Go module by the tag; the release's one file is the schema archive `learner-model-proto-<version>.tar.gz`, the same bytes on every build, signed keylessly and checked by `spec`'s `release-verify`. Immutable releases and the ruleset `v-tags` — `refs/tags/v*`, no update, deletion or non-fast-forward, no bypass — are on.** How the archive is built and checked is `release-files` and `ci-release` of the `justfile`; the release in CI is in [architecture.md](architecture.md); the line, the file and the commands to take and verify a release are in [README.md](README.md) of `docs/`, and the repository's README.md has none. An edit of what the file is built with releases a patch with the same code and schema, and no task bypasses it.

— *Why:*
- **A Go module, not a file:** `sim`, `engine` and `api` are Go, `go.mod` with `go.sum` pin the version, and the checksum database keeps the zip's hash for good; `spec`'s module is taken so already, and a second route would split the version in two.
- **The file is the schema:** the shared release needs a file in `dist/`, and the schema is the one thing needed without Go — for the export, to read the log without the library. The signature and the attestation of an immutable release tie it to the tag's commit.
- **A reproducible tar:** the bytes do not depend on the build's time and owner, so anyone can build the file again from the tag and compare, and `ci-release` holds it on every pull request.
- **The settings before the first release:** `spec`'s check of a release asks for them.
- **The Dockerfile whole in `paths`:** the filter sees files, not stages, and a raise of Go in the stage that builds the file may change the archive's bytes.
- **No command in README.md:** a README says only what the component is and where it stands; agents run the commands, and `docs/README.md` holds them.

— *Rejected:*
- an archive of the module's whole source, alone or beside the schema: it repeats the Go module, which consumers take by tag anyway;
- turning the settings on by the author's hands: the executor turns them on through `gh api` and checks them against the API's answer, as it does the ruleset `main`;
- a trial run of the shared release with `publish: false` on a pull request: `ci-release` holds the build of the file, and the shared release checks the signature itself after publishing;
- tests and `testdata/` in `paths`: the module zip carries them, but a release on every edit of a test or a golden file gives consumers versions in which nothing changed for import; an edit of a test comes out with the next release of code, as in `module-kit`;
- removing the Dockerfile from `paths`: a raise of Go in the stage that builds the file would give no release;
- moving the stage `toolchain` into a file of its own: a difference from `spec`'s recipe that every move to a tag would bring in anew.

— *Consequence:* `.github/workflows/release.yml`, `.github/workflows/ci.yml`, `justfile`; README.md; [README.md](README.md); [architecture.md](architecture.md); the ruleset `v-tags`.

## The `spec` tag v0.9.9

**R17. The move to `spec` v0.9.9 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## Shared checks of the same bits

**R18. The checks of the same bits are the shared jobs of `go.yml` of `.github` with the input `same-bits` — "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules" —, and the lists of the rules of code are this repository's data, `.github/coderules.json`. The repository keeps no recipes, jobs, rules test or controls of its own for them, and the ruleset `main` requires the four checks.** What each job runs and why each list holds what it holds is in [architecture.md](architecture.md).

— *Why:*
- `learner-model` and `engine` have the same checks, and two copies of the rules test had already diverged — `engine` had the rule on `replace`, this repository did not — and would diverge further; the rules are one, the lists each repository's own.
- "Go / arm64" runs `ci-test` on arm64 hardware and then its host's targets, so a job of the repository's own on arm64 would repeat it.

— *Rejected:*
- a `rules_test.go` of its own beside "Go / rules": two copies of the same rules, which the shared check is there to end;
- the recipes `test-targets` and `fma-diff` for a run by hand in the container: a third copy of the scenarios of `go.yml`, which would diverge from it;
- the tests' import list of `.github`'s control, with the packages `go/*`: the tests do not need them, and a list wider than needed lets extra through;
- sparing a release by leaving the `justfile` or `release.yml` behind when the checks move: `ci-test` would call removed recipes, and the shared workflows would stand at different releases.

— *Consequence:* `.github/coderules.json`, `.github/workflows/*.yml`, `justfile`; CLAUDE.md, README.md; [architecture.md](architecture.md); `run/DEBT.md`, Д10 and Д11. An edit of the `justfile` or `release.yml` gives a release.

## The `spec` tag v0.9.16

**R19. The move to `spec` v0.9.16 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The gate of `sim`

**R20. The gate of `sim` is the check "Gate" of `.github/workflows/ci.yml`, required by the ruleset of `main` with `integration_id` 15368, the parts no required check: four jobs "Gate / part k", each `sim`'s `just gate-part k 4`, then "Gate", `just gate-check bands`, all in the stage `toolchain` of `sim`'s devcontainer of the same tag, with the pull request's code in `GATE_LEARNER_MODEL`, on every pull request and every commit of `main`, on `ubuntu-24.04-arm`, with `sim`'s cache of Go and 15 minutes a job. `sim` comes from git by the tag and the SHA of its commit, `SIM_TAG` and `SIM_COMMIT` in the workflow's `env`, and each part first checks by `git ls-remote` that the tag names that commit. The gate stands on the latest tag of `sim`, taken only once `just gate exact` on the code of `main` is green to the bit; while a round of candidates is explored and judged, on the tag the round runs on. A candidate that moves the numbers is merged only after `sim` records its bands at the candidate's commit and tags them; then this repository's session raises the pin, and `sim` moves its base to the model's release.** How the jobs and the check of a new tag are built is in [architecture.md](architecture.md); what a pin brings — the base, the modules, the bands — is `sim`'s `gate.json` at its tag.

— *Why:*
- The tag and the SHA in `env`: so `spec`'s rule of the gate and `sim`'s README ask of a caller ([sim@09e5a84](https://github.com/altnovatrail/sim/blob/09e5a84/README.md)); the SHA holds the code, the tag the reader and the check, and a tag moved by mistake stops in seconds.
- `sim`'s recipes: they already write `go.work`, take `engine` from `gate.json` and check the build and the bands, and `GATE_LEARNER_MODEL` is their input for the caller ([sim@09e5a84](https://github.com/altnovatrail/sim/blob/09e5a84/README.md)).
- `sim`'s stage `toolchain`: the gate checks the line `go` and the Go it is built with, and `spec`'s rule of the gate names `sim`'s devcontainer of the same tag, so Go, just and git are the harness's.
- Four parts: the gate was the longest path of CI (measurements 8). The parts give the numbers of `just gate` to the bit, since a learner comes from its seeds alone, and `sim`'s merge holds them together ([sim@0a31313](https://github.com/altnovatrail/sim/blob/0a31313/justfile)).
- arm64: a part splits its share by the runner's processors, so all five jobs run on one kind of runner; arm64 is standard ([github/docs@e0b0c90](https://github.com/github/docs/blob/e0b0c90/data/reusables/actions/supported-github-runners.md)) and free for a public repository (sources 5.1), `sim`'s gate runs faster there than on amd64, and `sim`'s CI holds its numbers to amd64's bits. The cache's steps are `sim`'s own and spare five jobs the build of the bench.
- The latest tag: a candidate's red gate then means a shift of a measure, not new bands.
- Every pull request: `spec`'s rule of the gate asks it; an edit of a document does not move the gate, and a list of the paths that call it would guess what the model reads.
- The budget: 15 minutes a job, as in `sim`'s `gate.yml`, guards against a hung run (measurements 6 and 8); a gate that does not fit is cut by an edit of `sim`, not by the caller's time.
- The session raises the pin, not Dependabot or a shared workflow: the order of releases needs the bands in `sim` before the numbers move.

— *Rejected:*
- a reusable workflow of the gate, `sim`'s `gate.yml` by its tag or one of `.github`: its check is "Gate / <name>", not "Gate", so the ruleset of `main` would change, or a job "Gate" after it would lengthen the path; and `spec`'s rule of the gate makes the job the caller's own;
- recipes or a `go.work` of this repository's own: a copy that drifts from `sim`'s;
- `go mod edit -replace` instead of `go.work`: it edits the harness's `go.mod`, which `spec`'s rule of the gate rejects;
- a toolchain stage of this repository's own: the gate would build on another Go, just and git than the harness's;
- the gate on amd64 or on both: amd64 is slower, and both would repeat what `sim`'s CI holds;
- the gate only on commits that change the code: paths outside the release's `paths` pass the gate anyway, and `spec`'s rule asks for every pull request;
- a tag older than the latest: a candidate's red gate could mean new bands, not a shift of a measure;
- Dependabot or a shared workflow raising `SIM_TAG` and `SIM_COMMIT`: it would not wait for `sim` to record a candidate's bands;
- raising the budget in advance: without a measurement it is a number without data.

— *Consequence:* `.github/workflows/ci.yml`; the ruleset of `main`; README.md, CLAUDE.md; [architecture.md](architecture.md); the times of "Gate" in [measurements.md](measurements.md) and the runners' terms in [sources.md](sources.md). `ci.yml` is not in the release's `paths`, so an edit of it gives no release.

## The leading question of format 0.4

**R21. The leading question of a task of format 0.4 is help, as a technique on request: the answer after it goes into the step, the counters and the stability as it is, but does not set "mastered". Evidence of kind 2 carries it by a field of its own, `leading_question_on_request = 20`, in Go `Evidence.LeadingQuestionOnRequest`, apart from `technique_on_request`, as a fact of the attempt, with no new weight.** How it is read and encoded is in [architecture.md](architecture.md).

— *Why:*
- The question hints at the method more strongly than a technique — that is why the format closed it —, so it weighs no less than one, and help is the strongest weight the rules of evidence give without a new number.
- A field of its own keeps the fact apart in the append-only log: if the harness shows the question must weigh otherwise than a technique, a new version reads it by recomputation.

— *Rejected:*
- the same rule by the field `technique_on_request`: the log would mix the two facts forever, and no recomputation could split them;
- nothing, as a technique at the step: the answer after a question that hinted at the method would set "mastered";
- a smaller weight of the answer: one more constant without data, which the weights of evidence avoid.

— *Consequence:* `proto/novatrail/learnermodel/v1/evidence.proto`, `evidence.go`, `step.go`, the codec's files and `testdata/wire/`; README.md, CLAUDE.md; [architecture.md](architecture.md).

## The plan of E3

**R22. The model's candidates go into `sim`'s judgement as a round: each by its own draft pull request — one constant or a coherent set, with its own group of the estimate —, all judged together, chosen by an exploration of the harness that decides nothing; the first round took up to three, chosen by rough runs on the numbers Д5–Д7. An exploration steps a constant both ways, and an edge — guessing 0, full weight — only inwards. The threshold "mastered" goes only in a pair with the target of `engine`'s corridor, and the boundary of failure ½, which `engine` reads too, only with an edit of `engine` (Д16, Д7). A candidate that changes a number standalone's test vectors pin removes those vectors, with a record in this log: the tag `v0.1.0` holds the baseline with them, the gate and the golden vectors the new one.** The epic's plan is in [plan/E3.md](plan/E3.md).

— *Why:*
- Of those that pass, `sim`'s criterion chooses the simplest, so several candidates in one judgement give more chances than one, and a deciding run is far costlier than a rough one.
- Without an exploration the judgement would more likely reject all.
- A pair with `engine` needs `engine`'s candidate and a paired criterion, and `spec`'s order of releases names no way to release the boundary ½ together with `engine`.
- standalone's vectors are printed for standalone's numbers: a candidate that changes a number is not meant to pass them.

— *Rejected:*
- one candidate a round: if it fails, another round of judgement is needed;
- candidates by reasoning, without runs: faster, but blind;
- the pair of the threshold and the corridor's target, or the boundary ½, in a round of the model alone: they need work in step with `engine`'s session and its candidate;
- the model's constants as an internal set, so that standalone's vectors run on its numbers: more code, tests of code that is not released, and a task of its own first;
- the rule of τ by the share of correct answers in the exploration: it is a rule, not a number;
- C in the exploration: it is hard-wired into a formula and into an exported constant of `detmath`, and changing it is the model's design and the API of `detmath`.

— *Consequence:* RUN.md and `run/` — the epic's tasks, «Когда вернуться» of Д5–Д7, and Д16; `model.go` and the golden vectors of `testdata/`; CLAUDE.md, README.md and [architecture.md](architecture.md).

## The `spec` tag v0.12.0

**R23. The move to `spec` v0.12.0 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The exploration of E3

**R24. Three candidates went into `sim`'s first judgement, `model-1`, all in the group `forgets`, without a target and without a new field: `g10` — `stabilityGrowth` 10 instead of 15; `halve035` — the base of the division of τ 0.35 instead of 0.5; `forget-set` — both. The judgement ended with an exit: all three beat the base by the estimate, but each fails checks "no worse", all three the catching up after a break, `M3.settle` on F-break; so the base stays, g = 15 and the halving of τ, and the numbers stay in Д6. The exploration rejected the other numbers of Д5–Д7, and the harness does not measure the guessing of `predicate`.** The runs are measurements 12.1 and 12.2, the judgement measurements 12.3 and `sim`'s report [judgements/model-1/](https://github.com/altnovatrail/sim/tree/0364b40/judgements/model-1).

— *Why:*
- Both chosen numbers are constants of forgetting taken from systems for flashcards, and the exploration showed the model catches a skill's forgetting better when the stability grows more slowly and a wrong answer cuts it more; both hardly touch the other groups, and the constraints they broke at 300 learners were at the edge of the tolerance (measurements 12.2).
- The set, because both numbers move one group the same way, and the judgement must see whether they add up; `sim` chooses the simplest of those that pass, so the set wins only above the single ones beyond the noise.
- An exit, not a tolerance: the tolerances "no worse" were recorded before the first candidate and are not shifted for one.
- A rough look decides nothing: the judgement's gain in `forgets` was smaller than the rough look's, and failures that 300 learners took for noise were real (measurements 12.3). Why the catching up after a break is slower, the judgement does not read.

— *Rejected:*
- `g10`, `halve035` and `forget-set`, by the judgement: each fails "no worse" — all three `M3.settle` on F-break, `g10` and `forget-set` `M1.rms.200` on olympiad maths too, `forget-set` `E1.corridor` too, `halve035` two more at the edge of the tolerance (measurements 12.3);
- shifting a tolerance "no worse" for a candidate, instead of the exit: the tolerances are fixed before the first candidate;
- `bandStep` 2.0 and 3.0: both ways it breaks the estimate's error `M1` and the corridor `E1` on olympiad maths, dozens of constraints, and 3.0 worsens `forgets` too;
- σ₀ 2.0 and 3.0: 2.0 worsens `off-band`, and 3.0 at 300 learners is no better than the base and breaks 28 constraints. Measured later by fairness by the grade, which the base does not hold, σ₀ alone moves grade 1 only and, at the size of a deciding run, fails "no worse" on the cells of risk; σ₀ comes back only with the step after the trial series and its weight, and with them no set holds every constraint (measurements 5 and 7);
- I 0.12 and 0.202148: in neither direction does a group move beyond the noise;
- τ₀ 0.5 and 2 days: worse both ways, and 2 days worsens `forgets` too;
- g 22.5 and the base of τ's division 0.7: within the noise;
- ρ 1/4: worsens `forgets` and breaks the floor `M4`; ρ 1/2: improves `forgets` at the cost of `learns` and `off-band`, with 20 constraints broken;
- the weight 1/2 of an answer with help: worse `learns` and `pages`. It covers a repeat of a task too — `helped` counts as help the solution of this task already sent to the device —, so a repeat has no variant of its own;
- the guessing of `predicate` in the exploration: no pinned package has a check `predicate`, and the synthetic `S-predicate` runs as a cell without a graph, whose checks are neither met nor failed; the number stays in Д5 until a module with `predicate`;
- a candidate for each group: for σ₀ 3.0 and ρ 1/4 the best groups are paid for by others and by broken constraints, and the criterion judges the overall score;
- more than three candidates: the first round was held to three;
- a run of the three best at 1,000 learners before the judgement: costly (measurements 12.1), and it does not change the choice — the intervals of the overall score of the three were already above zero.

— *Consequence:* nothing in the code, the tests and the documents of `main`: the candidates were draft pull requests, closed with a reason. `run/DEBT.md`, Д5–Д7.

## The plan of E4

**R25. A release that other components wait for goes before a round of the model's candidates. A move to a tag of `spec` and a raise of the gate to a tag of `sim` are each a task and a pull request of their own, on the code of `main`, before the code they serve; if `sim` releases a new tag before such code is merged, the gate is raised to it first by its own pull request, so that a red gate does not confuse new bands with new code.** The epic's plan is in [plan/E4.md](plan/E4.md).

— *Why:*
- `spec`'s order of format 0.4 put the model's release with `time` and the leading question before `sim` and `modules` read it, and `sim`, `modules` and `api` waited for it; a faster `Readout.Task`, a change of the model of its own through the gate, went after it.
- A second round of candidates without an answer to why the model catches up with the truth after a break more slowly would repeat the first, which ended with an exit on it.
- The move goes first because the order of releases and the rules come with the tag, and a move is a task of its own; the gate goes before the code, so that the code passes the gate on the pin on which `sim` reads it.

— *Rejected:*
- the round of candidates first: `sim` and `modules` would wait for the release, and the round would start without an answer on the catching up after a break;
- the code before the move and the gate: the release would pass the gate on the bands of an older `sim` and `engine`, not on the pin where `sim` reads it;
- the move and the gate as one task: each has its own check and its own risk — the move releases a patch with the same code, and the gate may turn red through no fault of the code.

— *Consequence:* RUN.md, `run/EPICS.md`, `run/QUESTIONS.md`, `run/DEBT.md`; CLAUDE.md, README.md.

## The `spec` tag v0.14.0

**R26. The move to `spec` v0.14.0 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The `spec` tag v0.15.0

**R27. The move to `spec` v0.15.0 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## Format 0.4 in the model

**R28. Format 0.4 in the model: the check kind `CHECK_KIND_TIME = 5` is a free answer with guessing 0, read the same in every evidence kind; the leading question given on "Hint" is help, as a technique on request, by the field `leading_question_on_request = 20` of evidence kind 2, which always writes it, while in kind 1 it is "not recorded" and read as "there was no question". No state field was added, the states of the logs written before keep their bits, and the comments of the `.proto` are written for the reader of the export.** How the kinds are written and read, and what the comments say, is in [architecture.md](architecture.md).

— *Why:*
- the rules of the kinds keep "not recorded" apart from "no": `api`'s log holds kind-1 records, where the question could not be recorded, and they must stay distinguishable from records where it was not given — so a new field is a new kind;
- writing field 20 only in kind 2, `false` too, keeps the codec's promise: what the reader accepted, the writer gives back in the same bytes, so kind-1 bytes from the log stay the same bytes;
- before format 0.4 there were no leading questions, so reading "not recorded" as "no question" is how the earlier version read kind 1;
- `time` needs no kind: a new value of a list is an addition the rules of the kinds allow, and no kind-1 record carries it;
- the field's name follows the neighbouring `technique_on_request`, a fact of the attempt; its number is the next free one;
- `api` puts `common.proto`, `evidence.proto` and `state.proto` of the version that computed beside the JSON of the export, so each comment explains its field to a reader without code.

— *Rejected:*
- kind 1 with the field written only when `true`: less code, but "the question was not given" and "`api` did not record it" could no longer be told apart in the log, which the rules of the kinds require;
- field 20 always, in kind 1 too: kind-1 bytes from the log would change after reading and writing, breaking the codec's promise;
- the name `requested_leading_question`, as in `spec`'s contract: there it is a string, the question's text in the conversation between `app` and `api`; here it is a fact of the attempt beside `technique_on_request`;
- `CHECK_KIND_TIME` only in kind 2: a needless condition in the step, with nothing to distinguish — kind 1 with this value does not occur.

— *Consequence:* `proto/novatrail/learnermodel/v1/`, `types.go`, `evidence.go`, `codec.go`, `codec_read.go`, `codec_check.go`, `step.go` and the tests, `testdata/fold/`; README.md; [architecture.md](architecture.md); `run/QUESTIONS.md`, `run/DEBT.md`.

**R29. `Readout.Task` gives the same P to the bit and the same refusals as the reference, `Task` of `v0.1.5`, with one allocation per ordinary call — the slice of the answer — and the same API: what `Read` can compute once for a module it checked, it computes, and `Task` only reads it, so a `Readout` does not change after `Read` and `Task` may be called from many calls at once.** How it is built and tested is in [architecture.md](architecture.md); its time is measurements 2.

— *Why:*
- the time of `Task` within `engine`'s choice was allocation and garbage collection around the computation, not the computation itself (measurements 2);
- a change inside touches neither the API nor `engine`: `engine` and `sim` take it by raising the version, with nothing to wait for, and `sim`'s runs get faster;
- `api` calls the library from many requests, and the library keeps no mutable state they share.

— *Rejected:*
- a new API for many tasks — `Readout.Tasks`, or a task checked once when the pool is loaded: the gain could be larger, but `engine` would need a task and a release of its own, and `sim` would wait for both;
- a lazy cache of retentions in `Readout`, filled by the first `Task`: `Readout` would change after `Read`, and two calls at once would write to one memory;
- dropping `Task`'s checks for the modules `Read` read: its refusals are part of its contract, and without copies the checks cost little;
- changing `sortedCopy` to give the slice itself: `NewGraph` puts its result into the graph, which would then share memory with the caller's input;
- a reference with its own copy of `chanceOf` and `transfer`: SonarCloud took it for duplicated code above its threshold, so the reference holds `Task`'s checks, R and the starting point to `v0.1.5`, and takes P from `chanceOf` as the step calls it.

— *Consequence:* `readout.go`, `step.go`, `fold.go`, `codec_check.go`, `graph.go`, `bench_test.go`, `task_reference_internal_test.go` and the other tests; README.md; [architecture.md](architecture.md); `run/DEBT.md`.

**R30. A round of the model's candidates whose cause is not yet known runs in this order: the analysis of the cause, in parallel with the gate on a tag of `sim` with measured sizes; then the exploration at the size of the deciding run on the judgement's harness; then the drafts of the candidates, `sim`'s judgement and the confirmation on held-out learners, all on one pin of `sim`. Each task is a branch of its own from `main` with its own pull request, which the executor takes through green CI to the merge; the drafts of the candidates are draft pull requests into `main`, and only the chosen one is merged, by the judgement's task, releasing a patch.** How the second round went is the archive [plan/E5.md](plan/E5.md).

— *Why:*
- a judgement shows where the model fails but not why, and an exploration without the cause would again choose the forgetting numbers that delay catching up;
- the analysis does not depend on the harness's sizes, the exploration does: the coarse look of the first round at 300 learners did not see the failures the judgement found;
- a branch per task keeps one release and one recomputation of `api`'s cache: the round's code is only the drafts, of which only the chosen one reaches `main`, and the other tasks change CI and documents, outside the release's `paths`.

— *Rejected:*
- the exploration at once, at the sizes before `sim`'s measurement: its numbers would change with the new tag, and it would have to be repeated;
- waiting for a pin of a later release of `modules`: its date depends on `module-kit`, and the numbers explored do not;
- the pair with `engine` — the threshold "mastered", the corridor's target and the bound ½ (Д7, Д16) — in a round of the model alone: it needs `engine`'s half; it is epic E6 of `run/EPICS.md`, with the conditions on which it opens;
- one branch for the whole epic, each task a commit, merged once after the summary: no fewer releases than a branch per task, since the round's code reaches `main` only with the chosen candidate.

— *Consequence:* [plan/E4.md](plan/E4.md), [plan/E5.md](plan/E5.md); RUN.md, `run/EPICS.md`, `run/QUESTIONS.md`, `run/DEBT.md`; CLAUDE.md; README.md; [architecture.md](architecture.md).

## Catching up after a break

**R31. After a break the model is slow to catch up with a fall of the level because the step with retention puts wrong answers after a pause down to the forgetting of skills while the module's θ is not forgotten — not because of the forgetting numbers alone and not because the model's law diverges from the harness's truth. So forgetting stays in the step, the module's level is forgotten through n after a pause of the module, with no new state field, and numbers that speed forgetting up are explored only together with that rule.** The rule is built in [architecture.md](architecture.md), with forgetting; the runs that found the cause are measurements 13.

— *Why:*
- the step of θ is K_θ · w · (S − P) with P = c + (1 − c) · R̄ · σ: the lower R̄, the less P depends on θ and the smaller w, so an answer after a pause moves θ the more weakly, the more strongly the model believes in forgetting; the forgetting numbers of the first round believe in it more, catch up more slowly and are better calibrated after long pauses — the same shift (measurements 13);
- the base itself catches up slowly because θ is not forgotten and K_θ at n > 60 stands at its floor; a wider step after a pause of the module shortens it at a small price (measurements 13);
- on a truth that forgets by the candidate's own law most of the delay remains, so it comes from the numbers, not from the harness: `F-break` judges what it should;
- the module's last practice is the greatest t_last of its skills, and n is in the state already; the price is the meaning of n — after a pause v is wider and λ lower for every skill of the module, which `engine` reads too, so a candidate is measured on `E1.corridor` and `E2` as well.

— *Rejected:*
- a step without forgetting — P without retention and w = 1: it catches up faster on `F-break`, but on populations without a fall it takes forgetting for a fall of the level and θ goes below the truth, breaking the error of the estimate on four of the five other `F-` of olympiad maths and on `F-C1` of school maths (measurements 13);
- w = 1 for each answer: it catches up even more slowly, since w also damps a correct answer after a pause, which at a low P would pull θ up;
- exploring numbers that speed forgetting up without the forgetting of the level: they always delay catching up;
- a row «в реестр spec» for `sim` about the law of the truth of `F-break`: the delay comes from the numbers themselves;
- a new state field — a separate step of θ after a pause that touches neither n nor v: needed only if the price of n failed "not worse", which the exploration did not show; it would be a new state kind and a recomputation of `api`'s cache, the author's to decide;
- a field "time of the module's last practice": it is derived from the skills;
- small runs of 300 learners: the tolerance of `M3.settle` there is wider than the difference itself (measurements 13);
- tracing the step through the answers after a pause: the runs told the hypotheses apart without it.

— *Consequence:* RUN.md; `run/DEBT.md`, Д6; [architecture.md](architecture.md).

## The exploration of E5

**R32. The forgetting of the module's level — n := min(n, 20) before the step after a pause of the module of 14 days or more, the candidate `theta`, in the group `jumps`, with no target, two constants and no new state field — was the one candidate of the second round; `sim`'s judgement `model-2` chose it, the confirmation on `held-out-model-1` passed, and it is released as `v0.1.8`. None of the sets of the forgetting numbers with it holds all the constraints, and the weight 1/2 for an answer with help is rejected.** The rule and the reasons for N and d₀ are in [architecture.md](architecture.md), with forgetting; the exploration's runs are measurements 14; the judgement's report is [sim@7c32016](https://github.com/altnovatrail/sim/blob/7c32016/judgements/model-2/report.md). `held-out-model-1` is no longer held out: the next round of the model is confirmed on `held-out-model-2`.

— *Why:*
- N 20 and d₀ 14 hold every constraint, are better than the base and shorten catching up on `F-break` outside the noise, the round's target; the rule moves numbers only where a pause of the module of 14 days comes, and shifts the corridor and the loops `E2` nowhere outside the noise, so `engine` pays almost nothing (measurements 14);
- of `th20-14` and `th40-14`, equal by the estimate, only `th20-14` moves catching up on `F-break` of olympiad maths outside the noise; the judgement tries no second candidate on the same held-out learners, so a second one adds no choice;
- no target: the rule moves neither the share that never catches up on G3 nor the lag on G2-half;
- with the rule the forgetting numbers no longer delay catching up, but their failures outside the `F-` populations, where the rule does not act, come from the growth of τ at g 10 itself (measurements 14);
- the exploration ran at the size of the deciding run and on the judgement's harness, as the first round's coarse look did not see the judgement's failures; the price is that the judgement on the working learners is known in advance, and the confirmation on the held-out ones decides.

The numbers of the start (Д5) and ρ and τ₀ (Д6, Д7) were not explored: there were no new facts on them. A set that keeps the gain of the forgetting numbers without failures outside `F-` — a base of dividing τ between 0.4 and 0.5, or a growth of τ that slows down with τ — waits for a candidate with a cause of its own (Д6; [plan/E7.md](plan/E7.md)).

— *Rejected:*
- N = 10 (`th10-14`) and d₀ = 7 (`th20-7`): they fail `M5.caught` on `F-none` of school maths — a wide step there catches up nothing and makes P swing around the truth, and 7 days fire on ordinary pauses of skills;
- N = 40 (`th40-14`) and d₀ = 28 (`th20-28`): they hold everything but are weaker, and 28 days are no cheaper;
- the rule with `forget-set`, g 10 and the base of dividing τ 0.35 (`thfs20-14`): seven failures — the corridor and the error of the estimate outside catching up, and on `F-none` too;
- the rule with the base of dividing τ alone, 0.35 (`thh20-14`) or 0.4 (`thh040-20-14`): one failure each — `M1.rms.200` on G8 of olympiad maths, `M5.caught` on `F-none` of school maths — and the judgement runs the same learners with the same seeds and would repeat it number for number;
- a candidate of forgetting despite a failure by a hair: the judgement would repeat it, and its exit would give nothing but a run;
- a tolerance changed for a candidate: the rule — the larger of the author's tolerance and 4.3 standard errors of this candidate's difference — was set before the first candidate and holds, though the tolerance itself depends on the candidate's noise;
- more runs on the base of dividing τ, 0.45 and nearer 0.5: the exploration's budget was about ten runs, and `M5.caught` on `F-none` of school maths is too noisy to find the boundary without several runs per point;
- an exploration on the `F-` populations alone, then a full run for the finalists: a full run of every variant cost as many hours, and "not worse" was read on every population at once;
- the weight 1/2 for an answer with help (`help05`): 62 failures, the worst on school maths, in the same direction and strength as in the first round without leading questions, so leading questions do not change the conclusion (measurements 14);
- the rule before the check "no answers", as the exploration's variant had it: a piece of evidence without answers after a pause would narrow n, while such evidence must change nothing.

— *Consequence:* `step.go`, `model.go` and the tests, `testdata/`; `run/DEBT.md`, Д5–Д7; [architecture.md](architecture.md).

## The tag v0.16.3 of `spec`

**R33. The move to `spec` v0.16.3 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The shared workflows of `.github` v0.1.24

**R34. A raise of `.github` takes its latest release, for all the shared workflows at once, and goes before the code that the checks it adds are to check. "Go / targets" and "Go / arm64" run the tests in wasip1 too, in `node:wasi` of the Node of the stage `toolchain`, not only build them, so wasip1 is no debt. `.github/dependabot.yml` puts every reference to a shared workflow of `.github` into one group, `shared-workflows`, so that Dependabot raises them in one pull request, to one release; the actions of other owners stay a pull request each.** What the runs are is in [architecture.md](architecture.md).

— *Why:* the first code a new check is for passes it from its first pull request, and one release in every file is one version of the shared checks; the group moves no version, it only gathers the references that go to one release anyway.

— *Rejected:*
- a release older than the latest, as v0.1.22 and v0.1.23, without the run of wasip1: a raise takes the latest release;
- raising only `go.yml`: all the shared workflows stay at one version;
- a Dependabot pull request per reference: a rebase and a round of CI each, and the files would hold different releases until the last one is merged;
- a raise together with a rebuild of the documents and the session settings: it would mix CI with the rebuild;
- a raise with the first task of code that its new check is for: that task's first candidates would pass CI without the check.

— *Consequence:* the rest of Д10 — big-endian targets and the others with only the differential check, 32-bit arm, mips, software floating point — stays open. A raise changes `.github/workflows/release.yml`, one of the release's `paths`, so its merge releases a patch with the same code and schema. — `.github/workflows/*.yml`, `.github/dependabot.yml`; [architecture.md](architecture.md); `run/DEBT.md`, Д10.

## The rebuild of the documents and the session settings

**R35. `learner-model` keeps the way of working and the structure of documents of `spec`'s template: `.claude/settings.json` is the template's output for the kind `go` — the hook "one task, one session", compaction at 400 thousand tokens, claude.ai connectors off and `claudeMdExcludes` —, and `CLAUDE_CODE_PLUGIN_DIRS` in `containerEnv` of `.devcontainer/devcontainer.json` loads the plugin `gopls-lsp` from `spec/claude-plugins/`; CLAUDE.md is thin, with the rules of the parts `go`, `proto` and `ci` in `.claude/rules/`; the rules of RUN.md are the template's, rule 2 naming the rules of the parts; a done task's body and its gist sentence go to the epic's archive at once, and a task outside an epic, when there is no current epic, to the archive of the last one done, section «Задачи вне эпика, сделанные после него»; the facts no run checks are records of `docs/sources.md`, the results of measurements `docs/measurements.md`, and the budgets `docs-budgets.json`. `docs/` is in English, translated by a task of its own.** What each file holds is in CLAUDE.md. A limit of `docs-budgets.json` above `spec`'s block `repository` is the size of its text and about 15 % (measurements 1.9 spec).

— *Why:*
- by `spec`'s shared text of the rebuild, and where it is silent, as `engine` and `module-kit` did;
- a task outside an epic does not wait for the next epic's archive: its body would lie in the plan until then, a difference from the template that spends tokens in every session;
- the plugin comes from the container's environment so that it loads whatever the trust of the folder.

— *Rejected:*
- keeping the bodies of done tasks outside an epic in the plan until the next epic: the template takes a body away at once;
- removing the tasks outside an epic altogether, as `module-kit` did: more than `spec`'s shared text asks;
- raising the limit of a debt row to 6 KB for one long row, Д6: for good and for all rows, while a row's history belongs in the log and the archives;
- the history of the epics in CLAUDE.md, which every call loads: the template asks three to six sentences for the component, and the history is in the archives;
- the plugin by `enabledPlugins` and `extraKnownMarketplaces` of `.claude/settings.json`: it would depend on the trust of the folder, and the marketplace would be recorded in the `~/.claude` the containers share, by one container's path;
- rewriting the links of the log and the archives to the list of sources of [architecture.md](architecture.md): the list is a pointer to the records and keeps them right;
- translating text while moving it, as `sources.md` and `measurements.md` in the rebuild: two checks mixed; a translation is a task of its own, by `spec/tools/translation/README.md` and with `translation-check`.

— *Consequence:* the documents are outside the release's `paths`: the merge releases no tag. — `.claude/settings.json`, `.devcontainer/devcontainer.json`, AGENTS.md, CLAUDE.md, README.md, `.claude/rules/`; RUN.md, `run/`; `docs/`; `docs-budgets.json`.

## A shorter RUN.md

**R36. RUN.md keeps only what is live and no other file says: the executor's rules, the current epic and the bodies of its tasks not yet done. It has no section «Задачи вне эпика»: a task outside an epic is a body of `run/TASKS.md` until done, then goes to the current epic's archive, section «Задачи вне эпика, сделанные за его время», or with no current epic to the last done one's, «… после него». Rule 6 points to CLAUDE.md for what `docs-check-repo` checks instead of repeating it. The rules keep their numbers. The history of the plan goes to the archive of an epic, word for word, not deleted.**

— *Why:* every session reads RUN.md whole at its start, and history or a copy of the shared rules that CLAUDE.md already loads spends tokens in each; `sim` has the same shape, and one shape in both repositories is cheaper to read than two.

— *Rejected:*
- a section «Задачи вне эпика» in RUN.md, empty or made with the first such task: `run/TASKS.md` holds such a task until it is done;
- the template's list of what `docs-check-repo` checks in rule 6: CLAUDE.md says it, and a pointer is shorter than a repeat;
- deleting the history without moving it: the archive of the epic is where it would go at the epic's end, and the log links it;
- renumbering the rules: the plan and CLAUDE.md name them by number.

— *Consequence:* a move to a newer tag of `spec` that brings a new RUN.md of the template reads this entry first. The documents are outside the release's `paths`: the merge releases no tag. — RUN.md; `run/`; [plan/E5.md](plan/E5.md).

## Epic E7

**R37. The round of the forgetting numbers on the base where the level forgets — the base of the division of τ between 0.4 and 0.5 and a growth of τ that slows down with τ, alone and with g, then the drafts of candidates, `sim`'s judgement and a summary — is deferred, `[?]`, until a candidate of forgetting brings a cause of its own: a fault that the gate or a measurement found, a need of the product, a row of `spec`; `sim` judges no round of the model by plan. Returned to, it runs as the rounds before it: the gate on `sim`'s latest tag first, then the exploration at the size of the deciding run on the judgement's tag, with no analysis task; a branch and a pull request per task from `main`, merged on green CI by the executor, or by the orchestrator in its run, and the drafts as branches from `main` that wait for the judgement, which merges the chosen one and releases a patch.** What stands deferred is in [architecture.md](architecture.md), with the model of forgetting, and in Д6.

— *Why:*
- a gate before the exploration makes a red "Gate" on a candidate mean a shift of the measure, not of the tag;
- the exploration's numbers equal the judgement's on the working learners only at the size of the deciding run and on the same tag of `sim`, so the tag is held from the gate to the judgement's choice;
- the cause of `forget-set`'s failures outside `F-` is known — the growth of τ, on which the rule of θ does not act —, so the exploration needs no analysis first;
- a growth of τ that depends on τ needs no new field: τ is in the state.

— *Rejected:*
- an analysis task before the exploration — the cause is known, and an analysis would add no fact the exploration needs;
- one branch for the whole epic — the author chose a branch per task; the code of the round is only the drafts, one of which goes into `main`, so the round is one release either way;
- the exploration on an earlier tag of `sim`, without the gate — the judgement runs on the latest tag, and a later tag would shift the exploration's numbers;
- Д5, ρ, τ₀ and the weight of an answer with help in the round — no new facts about them.

— *Consequence:* the documents are outside the release's `paths`, so a merge releases no tag. — RUN.md; `run/EPICS.md`, `run/QUESTIONS.md`, `run/DEBT.md`, Д6; [architecture.md](architecture.md); [plan/E7.md](plan/E7.md).

## The tag v0.18.3 of `spec`

**R38. The move to `spec` v0.18.3 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The plan's layout of `spec`'s template v0.18.3

**R39. The work plan stands in the layout of `spec`'s template. RUN.md holds its rules and the current epic: the opening, Rules, Epics and the heading of table 1 in English, terse, for the session; the epic's description, its table of tasks and the bodies of tasks in Russian. Tables 2, 3 and 4 are `run/EPICS.md`, `run/QUESTIONS.md` and `run/DEBT.md`, and the open tasks outside an epic are `run/TASKS.md`, their openings and rows in English; no table stays in RUN.md or at the root. A done task's gist sentence goes to the archive with its body. The rules are the template's text, so that the next tag's template compares line by line — rule 4 keeps its items on the browser and AI too —, but for these differences: rule 2 names the rules of the parts; rule 6 points to CLAUDE.md for what `docs-check-repo` checks; rule 7 names the branches of an epic that the executor takes to the merge; the opening says what the component "has deliberately deferred", not "has still to decide", since all that the card left to decide is decided, and has no sentence on the first epic; the language line and rule 10 name the current epic's description and rows in Russian.** `docs-check` reads table 3 from `run/QUESTIONS.md` and table 4 from `run/DEBT.md`.

— *Why:* `spec` asks each repository for the layout with its move to a tag that has it, and the template's text makes the next move a comparison, not a rewrite. The differences kept are what this log decided and the template does not say: the rules of parts are this repository's; a pointer is shorter than a repeat; without the sentence on branches, rule 7 alone would read as "never commit" within an epic whose decision has the executor merge.

— *Rejected:*
- translating an epic's description and the bodies of its tasks: the template keeps them in Russian, since the author writes and reviews them;
- moving the row of a task outside an epic that already stands in table 1 to `run/TASKS.md`: the row is the epic's record of what ran while it went;
- keeping a closed part of a debt row: the template's opening keeps only what is open, and the log and the archive hold the rest;
- editing the words of the log's entries that name `DEBT.md` or "table 3 of RUN.md": RUN.md's opening names each table's file, so such a name still leads to its table; only a link to a file that is gone changes, since `docs-check` refuses it.

— *Consequence:* the documents are outside the release's `paths`: the merge releases no tag. — RUN.md; `run/`; no `DEBT.md` at the root; CLAUDE.md, README.md; [architecture.md](architecture.md); [plan/E5.md](plan/E5.md), [plan/E7.md](plan/E7.md); `docs/decisions.md`.

## Epic E8

**R40. The round of fairness by the grade — candidates that bring `M11.same` of `D-grade5` and `D-grade1` up to `sim`'s bound 0.90 under the model's criterion as is — ended with its exit: no set of the model's numbers of the start and the step holds the criterion, and its drafts and `sim`'s judgement `model-fair-1` are deferred, `[?]`, until there are 100 children, when long runs of fitting constants on `sim` come back. In their place stand the three rules in place of a calibration, built here, with `engine` unchanged, and `spec`'s data protection impact assessment, `spec/docs/dpia.md`, accepts the residual risk. A round returned to keeps the order: an analysis of the cause with a measurement, the gate on the tag of `sim` that holds the judgement's criterion, the exploration at the size of the deciding run on that tag, the drafts, the judgement, and the outcome for the impact assessment as a task of its own.** The cause and the sets are measurements 5 and 7. On a return the judgement waits for `sim`'s judgement of the candidate's commit, by `CANDIDATE_LEARNER_MODEL`, and only its merge for a tag of `sim` with the chosen one's bands; a candidate changes numbers or a rule, and a new field — a new kind of state and a recomputation of `api`'s cache — comes only by the author's answer to a row of table 3 in the planning of the drafts; a start by the grade other than `spec`'s rule is a row «в реестр spec» written before its draft, and the judgement merges it only after `spec`'s answer.

— *Why:*
- `sim`'s evaluation leaves to this repository why the start above the truth is left more slowly; without the cause an exploration guesses which numbers to move, and `sim` no longer judges a candidate without a cause;
- the target is `M11.same` of the two grades, since `fair` holds a candidate no worse than the base but does not reward a better one;
- the gate before the exploration and the exploration on the judgement's tag keep a red "Gate" unambiguous and the exploration's numbers the judgement's;
- the first child waits for the outcome, by the impact assessment's conditions before the first child, and it must come out on an exit too, when the judgement releases nothing.

— *Rejected:*
- another round in place of fairness, or E7's round of forgetting rewritten into it, by this repository's choice — `spec`'s plan gives the round, not a choice;
- the exploration without an analysis — `sim`'s evaluation leaves the asymmetry unexplained, and `sim` no longer judges a candidate without a cause;
- the exploration on another tag of `sim`, without the gate — the judgement goes on the tag of its criterion, and another tag would make a red "Gate" ambiguous;
- the judgement waiting for a release of `sim` — it needs `sim`'s judgement of the candidate's commit alone;
- the outcome inside the judgement — on an exit the judgement has nothing to merge, and the row for the impact assessment would wait on a pull request.

— *Consequence:* the documents are outside the release's `paths`, so a merge releases no tag. — RUN.md; [plan/E7.md](plan/E7.md), [plan/E8.md](plan/E8.md); `docs/decisions.md`; [architecture.md](architecture.md); [run/DEBT.md](../run/DEBT.md), Д5 and Д6; [run/QUESTIONS.md](../run/QUESTIONS.md); [measurements.md](measurements.md); README.md.

## The analysis of E8

**R41. Under grade 5 of olympiad maths, a band above the truth of grade 3, the estimate comes to grade 3's later than under grade 1, a band below it, because of a tail, not because every learner comes slowly, and the tail starts in the trial series, whose tasks move with the start: the engine chooses them by the model's P, so under grade 5 they are a band up too, and the maximum of the series' posterior moves with the start and the tasks alike — a learner who answers them as under grade 3 ends the series as far above as it started. After the series the step of θ, K_θ falling to its floor while each new skill's δ takes much of each pull, closes such a gap slowly. Under grade 1 the module has no band below the first, the series' tasks are grade 3's, and only the prior differs. The cause is the model's to shorten, and `spec`'s answer is the start a band below the band of the grade, in this library.** The runs are measurements 5.

— *Why:*
- in the series θ is the maximum of L(θ) = −(θ − θ₀)² / (2σ₀²) + Σ log-likelihood of the answers; when θ₀ and every βᵢ move by a band and the outcomes stay, the maximum moves with them: the series cannot tell a start a band too high from the truth unless the answers differ, and on three to five tasks one learner in five answers them as under the true band (measurements 5.2);
- a skill is within the engine's reach when a task of it has P ≥ 0.70, the corridor 0.70–0.85 ([`skills.go` of `engine` v0.1.6](https://github.com/altnovatrail/engine/blob/57e0ab8/skills.go)), so the tasks follow the start;
- grade 1's tail is the mirror: a learner far above the start answers the easy series all right, the likelihood of right answers has no maximum, and the prior keeps most of the gap;
- σ₀ moves grade 1 alone; the step after the series and a longer series move grade 5, and no set at the three cells `D-` brings it to 0.90 (measurements 5.3);
- the truth of the cells `D-` is on the middle band, so a start drawn toward the middle of the scale wins `M11` by construction: such a rule is measured on learners whose truth is off the middle band too;
- a set that brings the grades together by moving grade 3 off the truth is no fairer, so beside `M11` an exploration reads the share that never comes, |θ − θ*| of all three cells, and the price on the other cells.

— *Rejected:*
- the engine's choice of the series' tasks in place of the model's start — it would leave the model's start where the asymmetry is;
- the guessing floor as the cause — no item of the series has guessing, and the choices after it are alike under the grades;
- forgetting in the step as the cause — without it the tail barely shrinks (measurements 5.3);
- the bands of the tasks after the 50th answer as the cause — `M11.band` holds in `sim`'s evaluation;
- σ₀ alone as a set — it does not reach the cause under grade 5, and a narrower or a wider σ₀ is worse under both grades, the wider one because θ in the series swings far enough for the engine to give grade 1 other tasks than grade 3;
- the step without forgetting and K₀ of δ 0 as candidates — the first barely moves `M11` and leaves the tail, the second removes the skills' corrections, which the model keeps;
- runs of the other cells in the analysis — the cause is in the three cells `D-`, and the price of a set is the exploration's, at the size of the deciding run;
- a trace of the engine's choice — the tasks the model's trace shows told the causes apart without it.

— *Consequence:* nothing changes in the code and the tests. — [measurements.md](measurements.md); `docs/decisions.md`; [architecture.md](architecture.md); RUN.md; [sources.md](sources.md); [run/DEBT.md](../run/DEBT.md), Д5; [run/QUESTIONS.md](../run/QUESTIONS.md); [plan/E8.md](plan/E8.md).

## The exploration of E8

**R42. No set of the model's numbers of the start and the step — K₀ of θ and of δ, the floor and the decay of K_θ, the weight of the trial series and σ₀, alone and together — is a candidate for `sim`'s judgement `model-fair-1`: each one checked against "no worse" lifts `M11.same` of grades 1 and 5 and fails "no worse" elsewhere, and none brings grade 5 to the bound 0.90 (measurements 7.2, 7.3). The bound is not reached by fitting constants, and the base stays as it is. Long runs of fitting constants on `sim` stop, by the platform's rule in `spec`'s decision log, until there are 100 children; `sim`'s gate in CI, one-off analyses of how the model works and the rules in place of a calibration go on. The levers left are those rules and the engine's choice of the series' tasks.** How the sets are named and run is measurements 7.1.

— *Why:*
- the sets not checked lie between those checked or beyond the strongest, so none is likely to hold where its neighbours do not;
- the gap after the series is closed between n = 20 and n = 50: a K_θ that decays faster moves nothing or less, and the floor of K_θ moves the share that never comes, not `M11.same`, which reads the first 50 answers;
- a wider step of θ meets the full weight of a repeat on `W-repeat` (Д7; measurements 7.5) and breaks the error of the estimate on school maths; moving the pull from δ to θ holds `W-repeat`, but a δ that learns more slowly delays "mastered", and a smaller or a larger share breaks other cells (measurements 7.3, 7.4);
- a copy of the base's decisive file with a variant's lines on the cells of risk screens "no worse" there cheaply and exactly, but no more: the other cells stay the base's, and a set that failed four there failed eighteen on its decisive run (measurements 7.3, 7.4).

— *Rejected:*
- carrying the round to the bound 0.90 — the stop of long runs;
- more decisive runs — the stop, and the stopping rules too: the sets left are milder than `k04d02` or stronger than `kb06n50w10`, and the first lift less while the second fail more;
- a faster decay of K_θ, or a wider K₀ of θ below n = 50 alone — the first moves nothing or less, and the price of the second is not in the step after the 50th answer (measurements 7.2, 7.3);
- a start by the grade of the model's own choosing, as toward the middle of the scale — not measured: it wins `M11` on the cells `D-` by construction, and where a learner of a grade starts is `spec`'s rule to write first;
- a candidate with failures, as the best that can be had — the judgement's constraints are written down before the first candidate ([sim@61ed988](https://github.com/altnovatrail/sim/blob/61ed988/docs/decisions.md)), and the judgement would repeat each failure number for number on the same learners;
- a rule for the weight of a repeat in this round — the weights' own question, with its own cause, and its first variant, the weight 1/2 of an answer with help, failed far more checks (Д7).

— *Consequence:* nothing changes in the code and the tests, and the documents are outside the release's `paths`, so a merge releases no tag. — [measurements.md](measurements.md); `docs/decisions.md`; RUN.md; [architecture.md](architecture.md); [run/DEBT.md](../run/DEBT.md), Д5 and Д7; [run/QUESTIONS.md](../run/QUESTIONS.md); [plan/E8.md](plan/E8.md).

## The speed of CI and the tag v0.18.4 of `spec`

**R43. CI is fast, and what its checks prove stays the same. `ci-test` runs the tests twice, each a check of its own: `go test -race`, then the coverage with `-covermode=set`. The gate is four jobs "Gate / part k", each `sim`'s `just gate-part k 4`, and after them the check "Gate", `just gate-check bands`, all on `ubuntu-24.04-arm`, with `sim`'s cache of Go by the steps of `sim`'s `gate.yml` word for word. "Proto" takes the image of the stage `toolchain` by the restore and the script of the image steps of `.github`'s composite action `toolchain` word for word, with `SAVE` `"false"`, and no cache of Go. `.devcontainer/` is the output of `spec`'s template, with the stage `toolchain-core`, in which the release builds and which `toolchain` adds nothing to.** How the jobs run is in `.github/workflows/ci.yml`, the `justfile` and [architecture.md](architecture.md).

— *Why:*
- the gate was the longest path of CI, and `-race` with `-covermode=atomic` in one run took four times the two runs apart (measurements 3.1 spec; measurements 8);
- Go takes `atomic` by itself under `-race`, "significantly more expensive", while `set` says whether a statement ran ([test.go@go1.27.4](https://github.com/golang/go/blob/go1.27.4/src/cmd/go/internal/test/test.go#L232-L240)), all that `codecov/patch` reads, and both profiles mark the same statements (measurements 8);
- what the gate proves does not change: a learner comes from its seeds alone, `sim` merges the parts by its own check, and `sim`'s CI holds amd64 and arm64 to the same bits ([sim@0a31313](https://github.com/altnovatrail/sim/blob/0a31313/docs/decisions.md)); `spec`'s plan of the speed of CI moves the gate to arm64 and leaves such a change to the repository;
- the cache of Go, rejected as giving "almost nothing" when the gate was one job ([sim@09e5a84](https://github.com/altnovatrail/sim/blob/09e5a84/RUN.md)), spares five jobs the build of the bench, and its steps are `sim`'s own, not a key of this repository's.

— *Rejected:*
- one run with `-race` and the coverage, `-covermode=atomic`: four times the two runs apart;
- `-race` on part of the packages: all three together take about 13 s with it (measurements 8), so there is nothing to win;
- the two runs of `ci-test` at once: both use the runner's processors, and two outputs at once would mix in the log;
- `sim`'s `gate.yml` called by its tag: its check is "Gate / <name>", not "Gate", so the ruleset of `main`, a setting outside the repository, would change, or a job "Gate" after it would add a job to the path; the gate stays a job of the caller's own, as `engine` keeps it;
- the gate on amd64 or on both: amd64 is slower by about a minute, and both would repeat what `sim`'s CI holds;
- the composite action `toolchain` itself in "Proto": it is for a job that runs Go, and would take the caches of Go and save the image on `main`;
- a save of the image by "Proto": `go.yml`'s job `toolchain` saves the same key on a push to `main`, and a second save would race it for nothing;
- "Proto" in `toolchain-core`: the two stages hold the same here, and `go.yml` keeps one entry of `toolchain` for every job, which "Proto" shares.

— *Consequence:* `.devcontainer/Dockerfile`, the `justfile` and `.github/workflows/release.yml` are in the release's `paths`, so a change of them releases a patch with the same code and schema. The parts of the gate are no required check, and "Gate" keeps its name in the ruleset of `main`. — `justfile`; `.github/workflows/`; `.devcontainer/`; `spec/`; CLAUDE.md, README.md, RUN.md; [architecture.md](architecture.md); [.claude/rules/ci.md](../.claude/rules/ci.md); [sources.md](sources.md); [measurements.md](measurements.md); `docs/decisions.md`; [plan/E8.md](plan/E8.md).

## The tag v0.18.6 of `spec`

**R44. The move to `spec` v0.18.6 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The three rules in place of a calibration

**R45. The model holds the three rules in place of a calibration. A module's level starts a band's width below the band of the learner's grade, below the lowest band too, while the start band of the evidence and the state stays the grade's. An adult's correction, evidence of kind 3 that `LevelCorrection` gives, sets the level at the middle of the band chosen, in one step, unless it stands in that band already, ends the trial series and starts the guard's count anew. Three failures of the module in a row past the trial series lower the level by a band, once a run and not below −2.5: a failure is a piece of evidence whose answers the model expected, above the guessing floor, to be at least half known, and fewer than half of whose items with an answer were solved at the end of the attempt. `Read` gives the band the level stands in, and the state is of kind 2.** How each is computed is in [architecture.md](architecture.md), with the learner's state and with the evidence. `sim` judged the rules on their pull request's commit, and the gate stands on the tag of `sim` that recorded its bands, `SIM_TAG` of `.github/workflows/ci.yml`. The judgement held the target and the guard's false firings and failed one check of `fair` — the kind of task under grade 1, whose start is now below the lowest band — and the correction against `D-grade1`, while the correction held against the truth; `spec/docs/dpia.md` holds the numbers and accepts the failure of `fair` in the residual risk.

— *Why:*
- the start band stays the grade's: every log starts its modules by the new rule on recomputation, the start band keeps one meaning, no place 0 enters the evidence, whose places start at 1, and `spec` counts the start from the grade's band;
- the failure is the model's own expectation and the end of the attempt: in the corridor a learner at its level fails about one first answer in four, so three first answers wrong in a row would fire on most such learners within 200 answers, against the at most 5 % of false firings by which `sim` judges the guard (measurements 9.1); "stuck at the end" alone fired mostly on tasks the model already gave a P near 0, which say nothing of the level, and broke the error of the estimate (measurements 9.2); a task the model gave as fitting that the learner could not solve even on the ladder says the level stands too high, which is what the guard is for;
- above the guessing floor: P is never below c, so by P a choice of two would be fitting however far above the level, while R̄ · σ is the chance the learner knows the answer, and with no guessing it is P;
- one lowering a run: without it a learner stuck on a run of tasks that stay hard sinks a band each three to the floor;
- the middle of the band and the nearest band: a band's middle task, difficulty 3, stands at its offset, so the correction puts the level where the band's tasks are in the corridor and the band it stands in is the one set; the nearest middle makes "Easier" from a level near the next band's middle a step of about a band, not of almost two;
- the correction ends the series: the series finds the start and the adult gave it; left on, its next answer would compute θ anew from θ₀ and undo the correction;
- the form of the record: an earlier version skips field 21 and refuses evidence whose task names no module, so the record names the module in its task ID and holds nothing an earlier version reads as an answer, and it moves nothing there;
- `sim`'s judgement, not a rough look: a rough look of 300 learners a cell did not see the failure of `fair` (measurements 9.3).

— *Rejected:*
- a failure by a skill's rule or by the first verdict, fewer than half right at it — it fires on learners at their level (measurements 9.1);
- one "stuck at the end" without the model's expectation — it fires on tasks the model already took for too hard (measurements 9.2);
- the expectation by P with its guessing floor — by it a choice of two is never too hard;
- a lowering at each third failure of a run — a stuck learner would sink to the floor;
- a guard with no floor — a learner stuck on everything would sink without bound;
- a start band a band lower in the evidence, `BandByGrade` giving the band below, or a place 0 below the lowest — the start band would mean two things, a place 0 changes the checks of the evidence, and the logs before kind 3 would keep the former start on recomputation;
- a correction that keeps the series, its centre moved to the band set — a new field of the state, and five more answers would move the level a band either way, against "in one step";
- a correction to a band the current version does not have refused — a log written on a version with more bands could no longer be recomputed;
- the adjacency of the band checked in the library — `api`'s method checks it, and a log that `wrong_key` cut may put a correction after another level than the one it was made at, which the recomputation must still read;
- a record with a field of its own for the module and no task — an earlier version refuses it, and an earlier version must read a new form.

— *Consequence:* what the model computes changes, so "Gate" stands on the bands `sim` recorded for this change. — `rules.go`, `model.go`, `evidence.go`, `state.go`, `readout.go`, the codec; `proto/`; `testdata/` — the golden files, the golden log `rules` and the golden correction `testdata/wire/correction.binpb`; `.github/workflows/ci.yml`; [architecture.md](architecture.md); [measurements.md](measurements.md); CLAUDE.md, README.md; `docs/decisions.md`; [run/DEBT.md](../run/DEBT.md); [run/QUESTIONS.md](../run/QUESTIONS.md); [plan/E8.md](plan/E8.md).

## The tag v0.18.14 of `spec`

**R46. The move to `spec` v0.18.14 holds no decision of its own in force; what it brought is its commit and `spec`'s release.**

## The context of E8's sessions

**R47. E8's sessions carried more context a call than E7's because they were longer, not because of the documents: over each session's first 80 main calls E8 carried less than E7, what entered from the documents was E7's in size, and three long sessions are the whole excess over the reference (measurements 4, 10). So the length of a session is no reason to change `docs-budgets.json`, the order of reading that RUN.md sets or the documents. How to wait for a long run and the figure an analysis by `spec`'s threshold begins with are `spec`'s: a run longer than a few minutes goes to the background, as the shared rules say, and `tokens-report` prints the context of each session's first 80 main calls.**

— *Why:*
- a lower limit would cut what a task has to read, not the length that cost;
- RUN.md's reads come at the start of a session, where E8 is below E7;
- the length of a task is for RUN.md's rule on a task larger than one review to judge: the analysis and the exploration were long by their runs, a split of an exploration into sessions brings each the documents again, and nothing measures that it would be cheaper;
- a sentence of the shared rules reaches every repository whose task runs long.

— *Rejected:*
- raising or lowering `docs-budgets.json`, lower limits for the entries of `docs/decisions.md` or the sections of [architecture.md](architecture.md), or a split of the log — none of them brought E8 more than E7;
- another order of RUN.md's reading, or a split of an exploration into sessions — above;
- a sentence on waiting in RUN.md or in a rule of `.claude/rules/` here — no part's `paths` matches the temporary clone where an exploration runs `sim`, and RUN.md keeps the template's rules but for this repository's layout of the plan;
- E8 as a new reference — the reference is this repository's first row of at least 300 main calls after the translation of `docs/` into English, E7's (measurements 1.9 spec), and a threshold met by length is no reason to move it.

— *Consequence:* nothing changes in the code, the tests, CI, the container or the budgets. — [measurements.md](measurements.md); `docs/decisions.md`; `docs-budgets.json`; RUN.md; [run/QUESTIONS.md](../run/QUESTIONS.md); [plan/E8.md](plan/E8.md).

## The tag v0.18.27 of `spec`

**R48. The move to `spec` v0.18.27 and to `.github` v0.1.33 holds no decision of its own in force; what it brought is its commit and the two releases.**

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

A file of Go's specification or sources, of `golang/go` at the tag of its version, of the Go Wiki in `golang/wiki` at its commit or of a Go module at its version is a link beside the claim; here are only the licenses.

### 4.5. The licenses of Go's standard library

- [`LICENSE`](https://github.com/golang/go/blob/go1.27.3/LICENSE) [`src/math/exp.go`](https://github.com/golang/go/blob/go1.27.3/src/math/exp.go#L28-L34) and [`src/math/log.go`](https://github.com/golang/go/blob/go1.27.3/src/math/log.go#L17-L23) of `golang/go` at the tag go1.27.3, read 2026-11-24: Go's license is BSD-3-Clause, "Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met"; the pure Go `exp`, `expmulti` and `log` of `math/exp.go` and `math/log.go` carry the Sun Microsystems notice from FreeBSD msun: "Permission to use, copy, modify, and distribute this software is freely granted, provided that this notice is preserved".

## 5. CI, GitHub and tools

A file of `altnovatrail/.github` is git, a link to it at the commit of its release beside the claim; the versions of the tools are set in the `justfile` and `.devcontainer/`, and Codecov's state is the check `codecov/patch` of every pull request. Here are only GitHub's terms for its runners and the licenses of the tools.

### 5.1. GitHub's runners

- [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), read 2026-11-24: "Use of the standard GitHub-hosted runners is free and unlimited on public repositories" — its source, `data/reusables/actions/supported-github-runners.md` of `github/docs` at e0b0c90. `altnovatrail/learner-model` is public (`gh repo view --json visibility`, 2026-11-24).

### 5.6. The licenses of the tools of the checks and gopls

- GitHub's API, `repos/<repository>/license`, read 2026-11-24: govulncheck and gopls — BSD-3-Clause (`golang/vuln` and `golang/tools`), go-licenses — Apache-2.0 (`google/go-licenses`).

## 6. Other NovaTrail repositories

No records: a file of another NovaTrail repository is git, a link to the repository at its commit beside the claim.

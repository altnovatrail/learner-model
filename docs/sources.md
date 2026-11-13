# learner-model — источники

> Внешние факты, на которых стоят документы `docs/` репозитория: модель и код standalone, исследования, документация Go, protobuf, GitHub и Codecov, документы, код и релизы других репозиториев NovaTrail. Здесь их записи — адрес, версия или дата страницы, дата прочтения и цитата на языке источника, — а [architecture.md](architecture.md) держит утверждение и ссылку на запись по её номеру, как на раздел: `sources 1.1` (RULES.md, The bar; R143 spec).
>
> Запись — подраздел «N.M» темы: один источник или одно чтение нескольких источников для одного решения. Номер записи не меняется: новая встаёт в конец своей темы, новая тема — в конец файла. Источник, прочитанный снова, — пункт с его датой в той же записи, а прежний текст не правится. Образец — sources spec. Итоги замеров — не здесь, а в [measurements.md](measurements.md).
>
> Записи 1.1–6.10 перенесены сюда дословно 2026-11-13 из [architecture.md](architecture.md), раздела 4 «Источники», задачей T49 (R35): заголовок записи — имя прежнего пункта, а первая строка называет подразделы architecture и решения, для которых пункт прочитан. «Раздел N» внутри записи — раздел того текста, на который она ссылается. В architecture на месте прежнего раздела остался указатель: на него ссылаются журнал решений и архивы, а их не переписывают.

## 1. MathTrail: standalone

### 1.1. Модель ученика standalone

Для architecture 2.1–2.5.

- [mathtrail/mathtrail-standalone, `SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md), ветка `main`, коммит c6cbeb2, прочитано 2026-10-19:
  - раздел 2.1 — «a Rasch model with a guessing floor, updated Elo-style after every answer» и «P = 0.2 + 0.8 · σ(θ + δ_topic − β)»;
  - раздел 2.2 — «Only correctness enters the formula» и «The hint, the pace and the run of failures are recorded (04-profile) and change what is chosen next (section 3) — never the rating»; об исключении — «One kind of task is read for more than its correctness: a probe»;
  - раздел 2.5 — «освоено» ставится, если ученик на осторожном уровне решал бы среднюю задачу уровня «with **P ≥ 0.775**»; оно держится на уровне — «The topic counts as mastered only while its recommended point (2.3) stays at that level or below»; а при смене правила «The masteries of the earlier rule are cleared»;
  - раздел 2.8 — тест-векторы пакета рейтингов: «The implementation computes in `float64` and must match the four decimals printed here to within 1e-4».

### 1.2. Профиль standalone

Для architecture 2.4 и 3.5.

- [`docs/architecture/04-profile.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/architecture/04-profile.md), тот же коммит, прочитано 2026-10-19:
  - о сводке по темам — «This is the part that is never pruned — it is the summary that lets the history window stay short»; об окне — «`recent` — the last **20** entries, answers and skipped tasks alike»;
  - раздел Versions and migrations — «Reading an older version runs the migration chain in memory — pure functions, one per step, each with its golden file in `testdata/`».

### 1.3. Стенд standalone

Для architecture 3.2.

- [`docs/learners.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/learners.md), тот же коммит, прочитано 2026-10-19 и 2026-11-02: ход прогона — «A task is answered two minutes after it is issued, the next is issued at once, and every tenth answer starts a new day»; таблица генераторов, G3-drop — «their level falls by 1.0 at once, at an answer from the 50th to the 150th, as after a long break».

### 1.4. Код модели standalone

Для architecture 2.1, 2.5.

- [`internal/domain/rating/`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/rating), тот же коммит, прочитано 2026-10-19: в пакете рейтингов лежит и коридор, `corridor.go`; файлы без `_test` импортируют только `cmp`, `math` и `slices`.

### 1.5. Модель ученика standalone, подробно

Для architecture 3.1.

- [`SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md), коммит c6cbeb2, прочитано 2026-11-02. Сверено с `main` 89986bc от 2026-11-01: пакет `internal/domain/rating` с c6cbeb2 не менялся, а в SPEC 2 изменились только абзацы о пробах (R311 standalone) и о задачах с картинки (R322 standalone), без формул:
  - раздел 2.1:
    - «β = the shift of its level + (difficulty − 3)»;
    - «Each level stands 2.5 above the one before it, so the two hardest difficulties of a level overlap the two easiest of the next and no two points coincide» и «2.5 is a guess with no data behind it»;
    - о старте — «It is written into `ratings.start` when the profile is created, and θ starts there. The grade decides nothing else: changed later, it moves neither θ, nor θ₀, nor the trial series of 2.2.1» (О-56 standalone);
  - раздел 2.2:
    - «θ += K_θ · (S − P)»;
    - «K_δ = K₀ / (1 + a · n)» и «K_θ = max(K₀ / (1 + a · n), 0.05)» при K₀ 0.2 для θ и 0.4 для δ и a = 0.05;
    - «β is not updated in v1»;
  - раздел 2.2.1:
    - «the first five answers are a trial series»;
    - «a normal prior centred on the start times the chance of each answer given; its spread is σ₀ = 2.5»;
    - «the corrections of the topics do not move»;
    - «A grid rather than Newton's method from θ₀, because L can have two peaks»;
  - раздел 2.5:
    - «the topic becomes mastered at the highest level it is taught at, at or below the level of the task, whose middle task — difficulty 3, β = the level's shift — the child would answer with P ≥ 0.775 at the cautious level λ = θ + δ − √v»;
    - «v = 1 / (1/σ₀² + I · n) + 1 / (I · (20 + m))» и «0.9383 | ln(0.71875 / 0.28125)»;
    - «mastery never moves down»;
    - «On the learners' bench 3.7 % of these masteries are false, where the earlier rule's were 63.5 %, and they come in 4.6 answers against 9.9» (R187 standalone).

### 1.6. Код модели standalone, подробно

Для architecture 3.1.

- [`internal/domain/rating/`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/rating), коммит c6cbeb2, прочитано 2026-11-02:
  - `ladder.go` — `levelShift = 2.5`, а над `Start` — «A grade outside 1 to 6 starts from the nearest level's shift rather than from no number at all»;
  - `trial.go` — `TrialAnswers = 5`, `startSpread = levelShift`, `gridStep = 0.05`;
  - `mastery.go` — над `MasteredAt`: «the highest of the levels it is taught at, at or below the level of the task answered, whose middle task the child at the cautious level would answer at the corridor's middle or better».

### 1.7. Решения standalone о модели

Для architecture 3.1.

- [`docs/decisions.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/docs/decisions.md), коммит c6cbeb2, прочитано 2026-11-02:
  - R79 standalone — «The grade picks the child's starting point θ₀, the shift of its level, which is written into `ratings.start` when the profile is created and never moved again»;
  - R187 standalone — «Mastery is declared only above the level held, and lost after two wrong answers in a row, as before».

### 1.8. Ответы, подсказки, пропуски и ловушки в standalone

Для architecture 3.3.

- [`SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/SPEC.md), [`docs/decisions.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/decisions.md) и [`docs/architecture/04-profile.md`](https://github.com/mathtrail/mathtrail-standalone/blob/89986bc/docs/architecture/04-profile.md), ветка `main`, коммит 89986bc, прочитано 2026-11-02:
  - SPEC 2.2 — «Only correctness enters the formula» и «The hint, the pace and the run of failures are recorded (04-profile) and change what is chosen next (section 3) — never the rating, which stays an honest estimate of what the child knows»; о пробе — «A right answer to a probe lifts its topic» (R299 standalone) и «Whether a task is a probe is read off where the child stands when it is answered»;
  - SPEC 2.5 — «The hint would make the answer measure the hint, and five answers keep a lucky start from mastering a topic»;
  - R93 standalone — «"I don't know" is an answer, and a wrong one» и «a child who has seen the solution has not solved the task; counting it as wrong keeps the rating an honest estimate»;
  - R98 standalone — о пропущенной задаче: «Nothing that reads answers reads it — not the rating, the trial series, the misconception map, the streaks or the rule» и «punishing that would punish a stray tap»;
  - R101 standalone — «The answer is recorded once»;
  - `04-profile.md`, поле `traps` сводки темы — «How often this child fell for each trap in this topic, ever. The rule picks the two most frequent».

### 1.9. Журнал и версии у standalone

Для architecture 3.5.

- [mathtrail/mathtrail-standalone](https://github.com/mathtrail/mathtrail-standalone/tree/89986bc), коммит 89986bc на `main`, прочитано 2026-11-02 через `gh api`:
  - `docs/decisions.md`, R165 standalone, среди отвергнутого — «the levels worked back from the answers, which breaks the day the step of the rating changes»;
  - `docs/architecture/04-profile.md`, Versions and migrations, — «Reading a newer version is refused», потому что «two revisions are live during every Cloud Run rollout, and an old instance rewriting a new file would quietly drop whatever it did not understand», и «Unknown fields at the current version are ignored on read and are not written back».

## 2. Модели ученика и оценка ответов

### 2.1. Модель Раша

Для architecture 3.1.

- G. Rasch, [On General Laws and the Meaning of Measurement in Psychology](https://digitalassets.lib.berkeley.edu/math/ucb/text/math_s4_v4_article-22.pdf), Proceedings of the Fourth Berkeley Symposium on Mathematical Statistics and Probability, vol. 4, 1961, pp. 321–333, прочитано 2026-11-02 по текстовому слою скана:
  - раздел 3 — «For the responses to the items of an intelligence test the following model is suggested… The probability that person number v gives a correct answer to item number i is» — дальше формула (3.1). Формулу текстовый слой передаёт с ошибками, поэтому здесь она не повторяется: что у модели Раша и у Elo одна форма, говорит Pelánek 2016, ниже;
  - книги Rasch 1960 открытой копии нет, и она не прочитана. Standalone по той же причине цитирует статью 1961 года ([`research/literature/learner-models.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/research/literature/learner-models.md), коммит c6cbeb2).

### 2.2. Пол угадывания

Для architecture 3.1.

- A. Birnbaum, [Some Latent Trait Models and Their Use in Inferring an Examinee's Ability](https://faculty.ucmerced.edu/jvevea/classes/290_21/readings/week%209/Birnbaum.pdf), в F. M. Lord, M. R. Novick, Statistical Theories of Mental Test Scores, 1968, гл. 17, прочитано 2026-11-02 по текстовому слою скана:
  - раздел 17.3 — «it further assumes that if he does not know it he will guess, and, with probability c_g, will guess correctly».

### 2.3. Elo в обучении

Для architecture 3.1 и 3.2.

- R. Pelánek, [Applications of the Elo Rating System in Adaptive Educational Systems](https://www.fi.muni.cz/~xpelanek/publications/CAE-elo.pdf), Computers & Education 98 (2016), 169–179; авторская версия, прочитано 2026-11-02:
  - раздел 2.2:
    - обновление — «θs := θs + K(correct_si − P(correct_si = 1))»;
    - «In this formulation the Elo rating system is closely related to the Rasch model (one parameter model) used in IRT. These two approaches share the same functional form of the model for predicting the probability of a correct answer, they differ only in the parameter estimation procedure»;
    - об угадывании — «for question with k options the probability of a correct answer becomes» 1/k + (1 − 1/k) / (1 + e^(−(θs − di)));
  - раздел 2.3 — «An intuitive improvement used in most Elo extensions is to use an "uncertainty function" instead of a constant K»; о Math Garden — «Klinkenberg et al. (2011) use uncertainty initialized as U = 1 and updated after each answer as follows:» U := U − 1/40 + D/30, «where D is a time from a previous attempt (in days)». Статья Klinkenberg et al. 2011 не прочитана;
  - раздел 4.3 — «It is also useful to incorporate forgetting, particularly when we are modeling declarative knowledge (in contrast to procedural skills, which are the typical focus of student modeling)»;
  - книги Elo 1978 открытой копии нет, и она не прочитана. Standalone её тоже не цитирует.

### 2.4. Общий уровень и поправка

Для architecture 3.1.

- J. Nižnan, R. Pelánek, J. Řihák, [Student Models for Prior Knowledge Estimation](https://files.eric.ed.gov/fulltext/ED560514.pdf), Proceedings of the 8th International Conference on Educational Data Mining, 2015, pp. 109–116, прочитано 2026-11-02:
  - раздел 2.1 — «We use an uncertainty function U(n) = α/(1 + βn), where n is the number of previous updates to the estimated parameters»;
  - раздел 2.3 — «In addition to the global skill θs the model now uses also the concept skill θsc … we just correct the global skill by the concept skill».

### 2.5. Первая попытка

Для architecture 3.3.

- [KDD Cup 2010: Educational Data Mining Challenge, Data Format](https://pslcdatashop.web.cmu.edu/KDDCup/rules_data_format.jsp), прочитано 2026-11-02: «Correct First Attempt: the tutor's evaluation of the student's first attempt on the step—1 if correct, 0 if an error»; ошибкой первой попытки там считается и просьба о подсказке — «the first attempt for the step was an error (incorrect attempt or hint request)».

### 2.6. Тестлеты

Для architecture 3.3.

- E. T. Bradlow, H. Wainer, X. Wang, A Bayesian Random Effects Model for Testlets, Psychometrika 64 (1999), 153–168, [doi:10.1007/BF02294533](https://doi.org/10.1007/BF02294533). Статья в закрытом доступе; прочитана аннотация, через OpenAlex, 2026-11-02: «all items given to an examinee are unlikely to be conditionally independent (given examinee proficiency). Models that assume conditional independence will overestimate the precision with which examinee proficiency is measured. Overstatement of precision may lead to inaccurate inferences such as prematurely ending an examination in which the stopping rule is based on the estimated standard error of examinee proficiency».

### 2.7. Эффект плана и эффективный размер выборки

Для architecture 3.3.

- G. Kalton, J. M. Brick, T. Lê, [Estimating components of design effects for use in sample design](https://unstats.un.org/unsd/hhsurveys/pdf/chapter_6.pdf), глава VI в Household Sample Surveys in Developing and Transition Countries, ООН, ST/ESA/STAT/SER.F/96, март 2005, прочитано 2026-11-02 по текстовому слою:
  - «The effective sample size is the size of an unrestricted sample that would yield the same level of precision for the survey estimate as that attained by the complex design»;
  - эффект кластера — формула (18), 1 + (b − 1)ρ, — «depends on two factors: the subsample size within selected PSUs (b) and the intra-class correlation (ρ)». Формулы и букву ρ текстовый слой передаёт с ошибками, поэтому формула здесь записана по смыслу;
  - термин «design effect» глава возводит к Kish 1965 — книга не прочитана, и из неё ничего не утверждается.

### 2.8. Elo с несколькими навыками

Для architecture 3.3.

- Pelánek 2016, выше, раздел 2.4, прочитано 2026-11-02: многомерное расширение берёт случай, когда «each item belongs to a single knowledge component», и после ответа обновляет все навыки j с корреляцией c_ij: θ_sj := θ_sj + c_ij · K · (correct − P). Деления шага между основными навыками одного пункта там нет, и 1/J в разделе 3.3 выведено из производной, а не взято оттуда.

## 3. Забывание и повторение

### 3.1. Интервальное повторение: FSRS

Для architecture 3.2.

- [The Algorithm](https://github.com/open-spaced-repetition/awesome-fsrs/wiki/The-Algorithm), вики `open-spaced-repetition/awesome-fsrs`, коммит 7bcf7dc, прочитано 2026-11-02:
  - обозначения — «$S$: Stability (interval when $R=90%$)»;
  - FSRS-4.5 — R(t, S) = (1 + FACTOR · t / S)^DECAY, «where $R(t,S)=0.9$ when $t=S$», и «In FSRS-4.5, $DECAY=-0.5$ and $FACTOR=\cfrac{19}{81}$»; «This version uses 17 parameters», а у FSRS-6 — «This version uses 21 parameters»;
  - FSRS v1 и v3 — показательный закон, R(t, S) = 0,9^(t/S);
  - FSRS v4, рост стабильности после успешного повторения — «The larger the value of $S$, the smaller the $SInc$ value. This means that the higher the stability of the memory, the harder it becomes to make the memory even more stable» и «The smaller the value of $R$, the larger the $SInc$ value. This means that the spacing effect accumulates over time».

### 3.2. HLR и система Лейтнера

Для architecture 3.2.

- B. Settles, B. Meeder, [A Trainable Spaced Repetition Model for Language Learning](https://aclanthology.org/P16-1174.pdf), Proceedings of the 54th Annual Meeting of the Association for Computational Linguistics, 2016, pp. 1848–1858, прочитано 2026-11-02:
  - раздел 3.2 — ящики Лейтнера с интервалами в 1, 2, 4 суток и дальше: «All cards start out in the 1-day box», а при ошибке карточка переходит в ящик с интервалом короче;
  - раздел 3.3 — «This posits that memory decays exponentially over time»: p = 2^(−Δ/h), где h — «the half-life or measure of strength in the learner's long-term memory»; и «Under the assumption that half-life should increase exponentially with each repeated exposure (a common practice in spacing and lag effect research)» полупериод — функция признаков, веса которых учатся по данным Duolingo.

### 3.3. SM-2

Для architecture 3.2.

- P. Wozniak, [SuperMemo 2: Algorithm](https://super-memory.com/english/ol/sm2.htm), из магистерской работы 1990 года, на сайте с 1998-05-10, прочитано 2026-11-02: «With all items associate an E-Factor equal to 2.5»; «I(1):=1», «I(2):=6», «for n>2: I(n):=I(n-1)*EF», интервал — «(in days)»; «If the quality response was lower than 3 then start repetitions for the item from the beginning without changing the E-Factor».

### 3.4. Кривая забывания

Для architecture 3.2.

- H. Ebbinghaus, [Memory: A Contribution to Experimental Psychology, глава VII](https://psychclassics.yorku.ca/Ebbinghaus/memory7.htm), 1885, перевод 1913, прочитано 2026-11-02, § 29, Discussion of Results: «forgetting would be very rapid at the beginning of the process and very slow at the end» и «the quotients of the amounts retained and the amounts forgotten were inversely as the logarithms of the times».

### 3.5. Форма забывания

Для architecture 3.2.

- J. T. Wixted, E. B. Ebbesen, On the Form of Forgetting, Psychological Science 2 (1991), 409–415, [doi:10.1111/j.1467-9280.1991.tb00175.x](https://doi.org/10.1111/j.1467-9280.1991.tb00175.x). Сама статья не прочитана: открытая копия — скан без текстового слоя. Прочитана аннотация, через OpenAlex, 2026-11-02: «In all cases, the course of forgetting was best described by a simple power function of time relative to five other reasonable alternatives (linear, exponential, exponential-power, hyperbolic, and logarithmic)».

### 3.6. Забывание навыков

Для architecture 3.2.

- W. Arthur Jr., W. Bennett Jr., P. L. Stanush, T. L. McNelly, Factors That Influence Skill Decay and Retention: A Quantitative Review and Analysis, Human Performance 11 (1998), 57–101, [doi:10.1207/s15327043hup1101_3](https://doi.org/10.1207/s15327043hup1101_3). Статья в закрытом доступе; прочитана аннотация, через OpenAlex, 2026-11-02: «Results indicate that there is substantial skill loss with nonpractice or nonuse, with the amount of skill loss ranging from an effect size (d) of -0.01 immediately after training to a d of -1.4 after more than 365 days of nonuse» и «Physical, natural, and speed-based tasks were less susceptible to skill loss than cognitive, artificial, and accuracy-based tasks».

### 3.7. Интервалы в математике

Для architecture 3.2.

- D. Rohrer, K. Taylor, [The Effects of Overlearning and Distributed Practice on the Retention of Mathematics Knowledge](https://files.eric.ed.gov/fulltext/ED505642.pdf), Applied Cognitive Psychology 20 (2006), 1209–1224; рукопись авторов в ERIC ED505642, прочитано 2026-11-02 по текстовому слою. Испытуемые — студенты университета, «116 undergraduates» в первом опыте. General Discussion — «In Experiment 1, distributing 10 practice problems across two sessions instead of massing all 10 problems in the same session had no effect on one-week test scores but virtually doubled four-week test scores. In Experiment 2, increasing the number of problems solved in a single session from three to nine had virtually no effect on test scores on either the one- or four-week test».

### 3.8. Модель активации ACT-R

Для architecture 3.2.

- P. I. Pavlik Jr., J. R. Anderson, Practice and forgetting effects on vocabulary memory: an activation-based model of the spacing effect, Cognitive Science 29 (2005), 559–586, [PubMed 21702785](https://pubmed.ncbi.nlm.nih.gov/21702785/). Сама статья не прочитана: копия авторов не открылась. Прочитана аннотация, 2026-11-02: «each time an item is practiced it receives an increment of strength but that these increments decay as a power function of time» и «This mechanism limits long-term benefits from further practice at higher levels of activation and produces the spacing effect».

### 3.9. Поправка на время у Пеланека

Для architecture 3.2.

- R. Pelánek, [Modeling Students' Memory for Application in Adaptive Educational Systems](https://www.fi.muni.cz/~xpelanek/publications/memory-geography.pdf), Proceedings of the 8th International Conference on Educational Data Mining, 2015, pp. 480–483; авторская версия, прочитано 2026-11-02 по текстовому слою:
  - раздел 2.2 — «instead of P(m) to use P(m + f(t)), where t is the time (in seconds) from the last attempt», а «the value f(t) corresponds to temporal increase in memory activation due to (short) time from previous exposure of an item». Простой функции f данные не дали, и она ступенчатая, подобранная по данным, а за последней ступенью — «0 otherwise»;
  - данные — места на карте, slepemapy.cz, раздел 3.

### 3.10. Glicko

Для architecture 3.2.

- M. E. Glickman, [The Glicko system](http://www.glicko.net/glicko/glicko.pdf), без даты, прочитано 2026-11-02: «In the Glicko system, a player's rating changes only from game outcomes, but his/her RD changes both from game outcomes and also from the passage of time when not playing»; после t периодов без игр — RD = min(√(RD²_old + c² · t); 350).

## 4. Go

### 4.1. Операции с плавающей точкой в Go

Для architecture 2.5.

- [The Go Programming Language Specification](https://go.dev/ref/spec#Floating_point_operators), раздел Floating-point operators, Language version go1.27 (May 26, 2026), прочитано 2026-10-19: «An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements, and produce a result that differs from the value obtained by executing and rounding the instructions individually».

### 4.2. Явное преобразование против слияния в Go

Для architecture 2.5 и 3.5.

- [The Go Programming Language Specification](https://go.dev/ref/spec#Floating_point_operators), раздел Floating-point operators, Language version go1.27 (May 26, 2026), прочитано 2026-11-02: «An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding» и пример «r = float64(x*y) + z» среди тех, где FMA запрещено.

### 4.3. Исходники Go 1.27.3

Для architecture 2.5 и 3.5.

- каталог `src` Go 1.27.3 в контейнере (`go env GOROOT`), прочитано 2026-11-02:
  - `math/exp_asm.go` — `//go:build amd64 || arm64 || loong64 || riscv64 || s390x` и `const haveArchExp = true`: у `math.Exp` на этих архитектурах своя реализация на ассемблере, а на остальных, WASM тоже, — на Go (`math/exp_noasm.go`);
  - `math/exp_amd64.go`, строка 11, — `var useFMA = cpu.X86.HasAVX && cpu.X86.HasFMA`: путь в `math/exp_amd64.s` выбирается по процессору во время работы. Выключить его даёт `internal/cpu/cpu_x86.go` — `option{Name: "fma", Feature: &X86.HasFMA}`, то есть `GODEBUG=cpu.fma=off`;
  - `math/log_asm.go` — `//go:build amd64 || s390x`;
  - `math/pow.go`, строка 128, — `a1 = Exp(yf * Log(x))`;
  - правила компилятора: в `cmd/compile/internal/ssa/_gen/ARM64.rules` — `(FADDD a (FMULD  x y)) && a.Block.Func.useFMA(v) => (FMADDD  a x y)`, а в `AMD64.rules` — `(ADDS(S|D) (MULS(S|D) x y) z) && buildcfg.GOAMD64 >= 3 && z.Block.Func.useFMA(v) => (VFMADD231S(S|D) z x y)`.

### 4.4. Слияние, преобразования и порядок в Go

Для architecture 3.6.

- [The Go Programming Language Specification](https://go.dev/ref/spec), Language version go1.27 (May 26, 2026), копия `doc/go_spec.html` Go 1.27.3 в контейнере, прочитано 2026-11-02:
  - Floating-point operators — «An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements» и «An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding»; среди примеров, где слияние разрешено, — `t = x*y; r = t + z`;
  - Conversions — «In all non-constant conversions involving floating-point or complex values, if the result type cannot represent the value the conversion succeeds but the result value is implementation-dependent»;
  - For statements with range clause — «The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next».

### 4.5. Компилятор и стандартная библиотека Go 1.27.3

Для architecture 3.6 и 3.7.

- каталог `src` в контейнере (`/usr/local/go`), прочитано 2026-11-02:
  - правила слияния, у каждого условие `useFMA`: `cmd/compile/internal/ssa/_gen/ARM64.rules`, строки 1798–1809 — `FMADDD`, `FMSUBD`, `FNMSUBD`, `FNMADDD`; `PPC64.rules`, 1013–1014; `S390X.rules`, 1266–1267; `RISCV64.rules`, 808–810; `LOONG64.rules`, 770–781; `AMD64.rules`, 1501 — только `ADDS(S|D)` и только `buildcfg.GOAMD64 >= 3`. В `Wasm.rules`, `386.rules`, `MIPS.rules` и `MIPS64.rules` таких правил нет;
  - `generic.rules`, строка 1331, — `(Div64F x (Const64F <t> [c])) && reciprocalExact64(c) => (Mul64F x (Const64F <t> [1/c]))`: деление на константу с точной обратной становится умножением;
  - `cmd/compile/internal/ssagen/ssa.go`, строка 2709, — преобразование `float64` в `float64` даёт `OpRound64F`, и правила слияния через него не проходят;
  - `cmd/compile/internal/ssa/func.go`, `useFMA`, — «useFMA allows targeted debugging w/ GOFMAHASH», а `base/hashdebug.go`, строки 64–70, — значение «n» не выбирает ничего: `-d=fmahash=n` выключает все правила слияния;
  - `cmd/compile/internal/base/flag.go`, строки 268–273, — «quietly disable the convert hash changes», и `AMD64.rules`, строки 168–180, — другое преобразование чисел в целые на amd64 за этим выключателем;
  - `math/sqrt.go` — «Return correctly rounded sqrt»; `Wasm.rules`, строка 349, — `(Sqrt ...) => (F64Sqrt ...)`;
  - `math/exp.go` и `math/log.go` — чистый Go `exp`, `expmulti` и `log`, с уведомлением Sun Microsystems из FreeBSD msun: «Permission to use, copy, modify, and distribute this software is freely granted, provided that this notice is preserved»; лицензия Go — `LICENSE`, BSD-3-Clause;
  - `internal/cpu/cpu_x86.go`, строки 95–102, — `fma` среди выключаемых `GODEBUG` только при GOAMD64 ниже v3, а `runtime/extern.go` — «cpu.all=off disables the use of all optional instruction set extensions»;
  - `runtime/asm_amd64.s`, строка 69, — сборка v3 на процессоре без v3 не запускается: «This program can only be run on AMD64 processors with v3 microarchitecture support»;
  - `encoding/json` собран из файлов v2 — в `internal/buildcfg/exp.go` у базовых экспериментов `JSONv2: true`;
  - `slices/sort.go` — «This sort is not guaranteed to be stable»; `unicode/tables.go` — `const Version = "17.0.0"`;
  - `lib/wasm/go_js_wasm_exec` — тесты js/wasm идут в Node, `exec node --stack-size=8192 "$DIR/wasm_exec_node.js"`; `go tool dist list -json` — у `js/wasm` и `wasip1/wasm` `"FirstClass": false`.

### 4.6. Ограничения сборки по имени файла в Go 1.27.3

Для architecture 3.6; поправка 2026-11-07 к R09.

- каталог `src` в контейнере, прочитано 2026-11-07: `go/build/build.go`, `goodOSArchFile`, — имя файла до первой точки, без первого слова и без `_test`, кончается на `_GOOS`, `_GOARCH` или `_GOOS_GOARCH`; списки систем и архитектур — `KnownOS` и `KnownArch` в `internal/syslist/syslist.go`. Сборка с `-gcflags='-S -d=fmahash=n'` печатает на каждом месте, где правило слияния сработало бы, строку «fmahash triggered [DISABLED]» с хешем, но без имени файла и строки — где именно, показывает разность ассемблера.

### 4.7. Go и WebAssembly

Для architecture 3.6.

- [Go Wiki: WebAssembly](https://go.dev/wiki/WebAssembly), прочитано 2026-11-02: «Set `GOOS=js` and `GOARCH=wasm` environment variables to compile for WebAssembly»; «Go 1.21 introduced WASI as a supported platform. To build for WASI, use the `wasip1` port»; «At present, Go generates large Wasm files, with the smallest possible size being around ~2MB» и «Using `gz` compression reduces the ~2MB (minimum file size) example WASM file down to around 500kB»; «TinyGo supports a subset of the Go language».

### 4.8. protobuf для Go

Для architecture 3.7.

- [`google.golang.org/protobuf/proto`, `MarshalOptions`](https://pkg.go.dev/google.golang.org/protobuf/proto#MarshalOptions), v1.36.12, прочитано 2026-11-02: `Deterministic` — «the same message will always be serialized to the same bytes within the same binary», а «the deterministic serialization is NOT canonical across languages. It is not guaranteed to remain stable over time».

## 5. CI, GitHub и инструменты

### 5.1. Раннеры GitHub

Для architecture 3.6.

- [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), прочитано 2026-11-02: метки `ubuntu-24.04-arm`, `ubuntu-22.04-arm` и `ubuntu-26.04-arm` — в таблице стандартных раннеров для публичных репозиториев, и «Use of the standard GitHub-hosted runners is free and unlimited on public repositories». Процессоры раннеров x64 страница не называет. `altnovatrail/learner-model` — публичный (`gh repo view`, 2026-11-02).

### 5.2. Общая обёртка Go

Для architecture 3.6.

- [`.github/workflows/go.yml`](https://github.com/altnovatrail/.github/blob/af40f949676fb3a01abb0a1e6ffbc99ce824660b/.github/workflows/go.yml), `altnovatrail/.github`, `main` af40f94, прочитано 2026-11-02 через `gh api`: «by four recipes of that repository's own — ci-lint, ci-test, ci-vuln and ci-licenses — one job each, in the toolchain stage of the caller's .devcontainer/Dockerfile, where the caller pins its Go»; job идут на `ubuntu-24.04`, без arm64 и без WASM.

### 5.3. Общие workflow `.github` v0.1.15

Для architecture 3.6; R13.

- релиз [v0.1.15](https://github.com/altnovatrail/.github/releases/tag/v0.1.15), коммит 89d3cee, прочитано 2026-11-06; `gh release verify v0.1.15 -R altnovatrail/.github` — «Release v0.1.15 verified!»:
  - `go.yml` — «ci-test writes the coverage of the tests as one Go coverage profile, coverage/go.out, and the job coverage sends it to Codecov by OIDC, with no secret»; `ci-licenses` получает `ALLOWED_LICENSES` — «the ASF's Category A as SPDX identifiers»;
  - `codeql.yml` — «Besides actions, javascript-typescript and go are taken», а Go читается из сборки Go стадии `toolchain`; у v0.1.0 вход `languages` брал только `javascript-typescript`;
  - `secrets.yml`, `dco.yml`, `documents.yml` и `scorecard.yml` между v0.1.0 и v0.1.15 не менялись (`git diff --stat v0.1.0 v0.1.15`);
  - тексты вызывающих файлов, `codecov.yml` и порядок перехода на покрытие — [`docs/architecture.md`](https://github.com/altnovatrail/.github/blob/main/docs/architecture.md) `.github`, раздел 3.3, `main`, прочитано 2026-11-06.

### 5.4. Общие workflow `.github` v0.1.21

Для architecture 3.6; R18.

- релиз [v0.1.21](https://github.com/altnovatrail/.github/releases/tag/v0.1.21), коммит 706f406, прочитано 2026-11-08 в клоне во временной папке; `gh release verify v0.1.21 -R altnovatrail/.github` — «Release v0.1.21 verified!». С ним прежнее «на `ubuntu-24.04`, без arm64 и без WASM» выше верно только без входа `same-bits`:
  - `go.yml` — «With same-bits, four more jobs hold the caller's bits the same on every build»: «"Go / targets" runs the tests on the builds where the bits are most at risk, and builds and vets wasip1; "Go / fma-diff" compares the assembler of each package with fusion and without it, with a control of its own here; "Go / arm64" runs ci-test and the targets of its host on an arm64 runner; "Go / rules" holds the code to the rules of code by the tool tools/coderules of .github, of this workflow's own version, and the caller's lists, .github/coderules.json»; «Without the input the workflow is what it was»;
  - `tools/coderules/lists.go` — поля списков `code-imports`, `test-imports`, `test-math`, `requires` и `nested-modules`, «A field it does not know is a mistake»;
  - между v0.1.15 и v0.1.21 `codeql.yml` и `scorecard.yml` подняли только версии действий, `documents.yml` зовёт `spec/tools/toolchain/toolchain.sh`, а `release.yml` поменял только ветку npm (`git diff v0.1.15 v0.1.21`).

### 5.5. Общие workflow `.github` v0.1.24

Для architecture 3.6; R34.

- релиз [v0.1.24](https://github.com/altnovatrail/.github/releases/tag/v0.1.24), коммит fab102d, прочитано 2026-11-13 в клоне во временной папке; `gh release verify v0.1.24 -R altnovatrail/.github` — «Release v0.1.24 verified!»; последний релиз (`gh release list -R altnovatrail/.github`, 2026-11-13). Между v0.1.21 и v0.1.24 (`git diff v0.1.21 v0.1.24 -- .github/workflows/`):
  - `go.yml` — «"Go / targets" runs the tests on the builds where the bits are most at risk, wasip1 among them»: `GOOS=wasip1 GOARCH=wasm go test -count=1` с `-exec` — `node --stack-size=8192` и исполнитель на `node:wasi`, `new WASI({ version: "preview1", args, env: process.env, preopens: { "/": "/" }, returnOnExit: true })`; так же — в «Go / arm64». Ещё каждый job в стадии `toolchain` берёт её образ из кэша GitHub Actions по ключу ОС, архитектуры и хеша Dockerfile и строит его только при промахе, «so that an outage of what the Dockerfile downloads from does not fail the checks while the Dockerfile stays the same». Входы и права вызова те же;
  - `documents.yml` — тот же кэш стадии для `spec/tools/toolchain/toolchain.sh` через `TOOLCHAIN_IMAGE_FILE`, который скрипт читает с `spec` v0.13.0;
  - `release.yml` и `scorecard.yml` — только ссылки в комментариях; `codeql.yml`, `secrets.yml` и `dco.yml` не менялись.

### 5.6. Инструменты проверок и gopls

Для architecture 3.6; R13, R14, R19.

- прокси модулей Go, `https://proxy.golang.org/<модуль>/@latest`, прочитано 2026-11-06: `golang.org/x/vuln` — v1.9.0 от 2026-10-11, `github.com/google/go-licenses/v2` — v2.0.1 от 2025-09-08, `golang.org/x/tools/gopls` — v0.24.0 от 2026-10-14. Лицензии — по API GitHub, `repos/<репозиторий>/license`, 2026-11-06: govulncheck и gopls — BSD-3-Clause (`golang/vuln` и `golang/tools`), go-licenses — Apache-2.0 (`google/go-licenses`). С `spec` v0.9.16 lock-файл gopls — заготовки `spec`: `.devcontainer/` из вывода рецепта `new-repo` для `learner-model` вида `go` отличается от прежнего только комментариями, а `go.sum` — тот же байт в байт (`diff -r`, 2026-11-09).

### 5.7. Codecov

Для architecture 3.6; R13.

- прочитано 2026-11-06: приложение Codecov стоит у организации на все репозитории — `gh api orgs/altnovatrail/installations`, `"app_slug":"codecov","repository_selection":"all"`; `learner-model` в Codecov уже есть, но не включён — `https://api.codecov.io/api/v2/github/altnovatrail/repos/learner-model/`, `"active":false,"activated":false`.

## 6. Другие репозитории NovaTrail

### 6.1. Журнал и кэш у `api`

Для architecture 3.5.

- [altnovatrail/api](https://github.com/altnovatrail/api/tree/d896943), коммит d896943 на `main`, прочитано 2026-11-02 через `gh api`:
  - R06 api, «Что решено», — свидетельство `api` хранит «его тип как есть и столбцы для себя: ID профиля, ID попытки, дайджест пакета, модуль, время и порядковый номер. Содержания в свидетельстве нет», а кэш — «строка на профиль: состояние, версия `learner-model`, закреплённый релиз `modules` и номер последнего учтённого свидетельства»; выгрузка — «журнал, текущее состояние и выборы с кодами причин»;
  - поправки 2026-11-02 к R06 api — «Пока идёт выкладка, запросы может обслуживать и прежний образ: пометку новее своей — по релизу или окну — он не переписывает, а состояние для ответа считает, не записывая его» и «выбор движка хранит и своё время»;
  - R08 api — «Свидетельство о навыке без пары не идёт ни одному навыку, а о навыке, которого граф текущей версии не знает, — ждёт в журнале, пока следующая версия его не назовёт»;
  - architecture 3.7 api, «Свидетельство», — «Нужен ли весу свидетельства факт «в уроке» — попытка знает своё прохождение, — решает `learner-model`».

### 6.2. Повтор выбора у `engine`

Для architecture 3.5.

- [altnovatrail/engine](https://github.com/altnovatrail/engine/tree/08bc48a), коммит 08bc48a на `main`, прочитано 2026-11-02 через `gh api`: в RUN.md T03 engine решает, «что входит в запись выбора и что из журнала свидетельств нужно, чтобы повторить выбор по сохранённым версиям», а T04 engine поля, которых нет в типах `learner-model`, шлёт «строкой таблицы 3 «в реестр spec», к T05 learner-model». Детализация `engine` не начата.

### 6.3. Стенд `sim` после суда `model-2`

Для architecture 1 и 3.6; поправка 2026-11-13 к R20; поправка 2026-11-13 к R32.

- [altnovatrail/sim на `v0.1.13`](https://github.com/altnovatrail/sim/tree/v0.1.13), прочитано 2026-11-13: релиз неизменяемый, от 2026-11-13, последний, тег лёгкий, на коммите `7c3201622165070519a489f7dfd2066dbf06d189`, мерже pull request #58 `sim` (`gh release view v0.1.13 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — «Release v0.1.13 verified!»). `gate.json` — `learner-model` на a063835 без версии, `engine` v0.1.6 на 57e0ab8, `modules` v0.2.5, `spec` v0.16.3, 69 чисел, полосы записаны заново — `"T47: the judgement model-2 chose theta, learner-model a063835, and the held-out learners confirmed it (judgements/model-2/report.md; R46)"`; `go.mod` — `learner-model` `v0.1.7-0.20261113040121-a063835087c7`. Отчёт суда — `judgements/model-2/report.md`: «Chosen and confirmed: theta.» Между `v0.1.10` и `v0.1.13` рецепт `just gate` и `internal/gate/` не изменились, `Learners` — 250, а в README «Gate» и «Taking a release» — только «bench» стало «simulator» (`git diff v0.1.10 v0.1.13`).

### 6.4. Стенд `sim` с размерами после `learner-model` v0.1.6

Для architecture 1 и 3.6; поправка 2026-11-12 к R20.

- [altnovatrail/sim на `v0.1.10`](https://github.com/altnovatrail/sim/tree/v0.1.10), прочитано 2026-11-12: релиз неизменяемый, от 2026-11-12, последний, тег лёгкий, на коммите `ff921be033917605b7be9d8d194a3c41a90c280c` (`gh release view v0.1.10 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — «Release v0.1.10 verified!»). `gate.json` — `learner-model` v0.1.6 на 7d796e2, `engine` v0.1.6 на 57e0ab8, `modules` v0.2.5 — `olympiad-math` и `grade-3-math` 0.3.0, `english-essays` 0.2.0, — `spec` v0.15.0, 69 чисел с полосами, записанными заново на 250 учениках, — «T46: the gate takes 250 learners a cell, not 150, after learner-model v0.1.6 made its readout faster». `internal/gate/gate.go` — `const Learners = 250`: «250 keep the slowest run of the job on 150, scaled, at about ten minutes at every caller, whose gate keeps no cache of Go, within the twelve that 300 would only just meet»; рецепт `just decisive` — `-learners 2000`. Между `v0.1.8` и `v0.1.10` — два тега: `v0.1.9` — `time` и наводящий вопрос формата 0.4 (T42, R41 `sim`), пакеты `modules` v0.2.5 с ними (T43, R42 `sim`) и база на `learner-model` v0.1.6 (T45, R43 `sim`), а `v0.1.10` — размеры (T46; поправки 2026-11-12 к R31 и R33 `sim`). Раздел README «Gate», вход `GATE_LEARNER_MODEL` и раздел «Taking a release» — те же (`git diff d957fc7 ff921be -- README.md justfile`): в README изменились только «Status» и размер решающего прогона в «Judging a candidate».

### 6.5. Стенд `sim` на `spec` v0.15.0

Для architecture 1 и 3.6; вторая поправка 2026-11-11 к R20.

- [altnovatrail/sim на `v0.1.8`](https://github.com/altnovatrail/sim/tree/v0.1.8), прочитано 2026-11-11: релиз неизменяемый, от 2026-11-11, последний, тег лёгкий, на коммите `d957fc7e9ab411fb1f2e9442b99e2a888d8d6775` (`gh release view v0.1.8 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — «Release v0.1.8 verified!»). Между `v0.1.6` и `v0.1.8` — переход `sim` на `spec` v0.15.0 (T38, R39 `sim`) и суд, который держит кандидата одной библиотеки ко второй библиотеке базы (T40, R40 `sim`); в `gate.json` изменилась только зависимость `spec`, а библиотеки, пакеты, полосы и рецепт `just gate` — те же (`git diff v0.1.6 v0.1.8`).

### 6.6. Стенд `sim` на `engine` v0.1.6

Для architecture 1 и 3.6; поправка 2026-11-11 к R20.

- [altnovatrail/sim на `v0.1.6`](https://github.com/altnovatrail/sim/tree/v0.1.6), прочитано 2026-11-11: релиз неизменяемый, от 2026-11-11, последний, тег — на коммите `a97be529249042595e544da6192555fa783ff9ec` (`gh release view v0.1.6 -R altnovatrail/sim`, `gh release list`, `git ls-remote`, `gh release verify` — «Release v0.1.6 verified!»). `gate.json` — `engine` v0.1.6 на 57e0ab8, `learner-model` v0.1.4 на acfc4f2, `modules` v0.2.0, `spec` v0.12.0, 150 учеников на ячейку (T41 `sim`; поправки 2026-11-11 к R26 и R31 `sim`). Рецепт `just gate` пишет рабочее пространство скрытым рецептом `workspace`, а `GATE_LEARNER_MODEL` берёт папку, как прежде, или коммит полным SHA; README, «Gate» — «A caller puts its own code in place of a library with `GATE_LEARNER_MODEL` or `GATE_ENGINE`: a folder, or a commit by its full SHA» (`git diff v0.1.2 v0.1.6`). Замер `sim` по самому медленному раннеру — «Gate» `learner-model` на 150 учениках около 11,2 минуты при бюджете 15 (поправка 2026-11-11 к R31 `sim`).

### 6.7. Стенд `sim` на новом закреплении

Для architecture 1 и 3.6; поправка 2026-11-09 к R20.

- [altnovatrail/sim на `v0.1.2`](https://github.com/altnovatrail/sim/tree/v0.1.2), прочитано 2026-11-09: релиз неизменяемый, от 2026-11-09, тег — на коммите `e61f90aa1bd894d1ce8f35547e0757d38f28c0b2` (`gh release view v0.1.2 -R altnovatrail/sim`, `git ls-remote`, `gh release verify` — «Release v0.1.2 verified!»). Рецепт `just gate` и вход `GATE_LEARNER_MODEL` — те же, что у `v0.1.0`, а «Taking a release» README не изменился (`git diff v0.1.0 v0.1.2`). `gate.json` — `modules` `v0.1.4`, `olympiad-math` 0.1.1, ячейка `A-linked` на двух математиках, по 100 учеников на ячейку; библиотеки в нём по-прежнему `learner-model` `v0.1.0` на d2620a1 и `engine` `v0.1.0` на c6c5bb9, как и в `go.mod`. Свой CI `sim` на pull request #40 `sim` — «Gate / amd64» 11 минут, «Gate / arm64» 6 минут (RUN.md `sim`, T30, «Сделано»). «Gate» здесь на pull request #41 — 7 мин 10 с, из них сборка стадии `toolchain` — 23 с (RUN.md, T29).

### 6.8. Стенд `sim` и его гейт

Для architecture 1 и 3.6; R20.

- [altnovatrail/sim на `v0.1.0`](https://github.com/altnovatrail/sim/tree/v0.1.0), прочитано 2026-11-09: релиз неизменяемый, от 2026-11-09, тег — на коммите `09e5a84f8a26fedcf5c8996bf81bf29c616917a2` (`gh release view v0.1.0 -R altnovatrail/sim`, `git ls-remote`). README, «Gate» — «A caller puts its own code in place of a library with `GATE_LEARNER_MODEL=<folder>` or `GATE_ENGINE=<folder>`», а «Taking a release» — «write both the tag and the SHA of its commit in its workflow… check out sim at that SHA, then run `just gate`… in the toolchain stage of its devcontainer»; код выхода — 0, когда гейт зелёный, 1 — когда красный или вход отвергнут, 2 — на вызов, которого рецепт не знает. `gate.json` называет `learner-model` `v0.1.0` на d2620a1 и `engine` `v0.1.0` на c6c5bb9; многоразового workflow нет намеренно (R05 `sim`, «Отвергнуто»). Замер `sim` на раннере — около четырёх минут на `ubuntu-24.04`, «кеш Go почти ничего не даёт» (RUN.md `sim`, T24).

### 6.9. Модуль Go `spec`

Для architecture 3.6.

- [`spec.go`](../spec/spec.go) в подмодуле на v0.6.6: `//go:embed all:format`.

### 6.10. Соседи о сборке

Для architecture 3.6.

- прочитано 2026-11-02 через `gh api`:
  - `engine`, `main` 08bc48a: T05 engine решает, «что нужно движку, чтобы однажды выбирать и на устройстве», и спрашивает, «проверять ли сборку для устройства в CI с первого кода»; T04 engine поля, которых нет в типах `learner-model`, шлёт «строкой таблицы 3 «в реестр spec», к T05 learner-model»; в разделе 2.4 — «у навыка — ID, название, описание, ярус, пререквизиты и оцениваемость»;
  - `sim`, `main` 067f624, docs/architecture.md, раздел 2.3, — «полосы стража в standalone записаны на amd64 при `GOAMD64` v1 и на процессоре с FMA. На arm64, на amd64 при `GOAMD64` v3 и выше и на процессоре без FMA последние цифры чисел могут разойтись, хотя и в пределах полос»;
  - `app`, `main` 375a09b, docs/architecture.md, раздел 3.1, — «Приложение — веб, для планшета и компьютера, без офлайна и нативных приложений», на React 19, Vite 8 и TypeScript 7 (R03 app).

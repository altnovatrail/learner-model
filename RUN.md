# RUN.md — work plan of learner-model

What learner-model does, does not do and has deliberately deferred — [docs/architecture.md](docs/architecture.md); why — [docs/decisions.md](docs/decisions.md). The platform is in `spec/`, a submodule at tag v0.18.3: the component's card — architecture 3.8 spec; the order of repositories — repositories 6 spec; learners' data and law — privacy spec; the platform's decisions — spec's log. Small tasks, one at a time, each reviewed by the author; only the current epic is in detail (R19 spec). Other tables, each a file opening with its rules (R153 spec): 2, epics for later — [run/EPICS.md](run/EPICS.md); 3, questions to the author — [run/QUESTIONS.md](run/QUESTIONS.md); 4, debt — [run/DEBT.md](run/DEBT.md). Open tasks outside an epic — [run/TASKS.md](run/TASKS.md) (R155 spec).

Language (R154 spec): rules and tables in English, terse, for the session; the current epic — its description, table 1's rows and the bodies of its tasks — and the bodies of `run/TASKS.md` in Russian: the author writes and reviews them.

## Rules

A task starts with «Выполни задачу T01 из RUN.md» in a new session in the devcontainer. One task — one session (R143 spec): the hook `UserPromptSubmit` of [.claude/settings.json](.claude/settings.json) — `spec/tools/session/` — stops «Выполни задачу T…» or «планируй T…» of another task in the same session and names `/clear`; the same task passes.

1. Do only the named task; touch other tasks' files only where it requires.
2. Read first: [CLAUDE.md](CLAUDE.md) with the rules it imports (RULES.md of `spec/`), the task — in table 1, or only its own body of `run/TASKS.md` (its opening says how) —, the sections it cites — its own and spec's —, [docs/decisions.md](docs/decisions.md).
   - Conflict with the documents or something missing — stop and ask; never guess.
   - A gap in the component itself — a subsection of [docs/architecture.md](docs/architecture.md), section 2. A gap across repositories — a row of [run/QUESTIONS.md](run/QUESTIONS.md) whose question opens with the mark «В реестр spec»: spec's session finds it on `main`, and it gets its О-, П- or Т- number once the author decides it there (R81 spec; RULES.md, Across repositories).
   - From `run/DEBT.md` and `run/QUESTIONS.md` take only the rows the task and its sections name and those `grep -n` finds by the paths and parts it changes (R143 spec).
   - The checks of a part the task changes — its rule in `.claude/rules/` (CLAUDE.md, Environment; R35).
3. **Devcontainer only.** If `NOVATRAIL_DEVCONTAINER` is not `1`, stop before the first edit and ask to reopen in the container; "only documents" is no reason to go on. If `spec/` is empty, ask for `git submodule update --init`. Exception: a cloud session the author started for the task runs the checks it can reach, and its report names each one that did not run.
4. Outside this repository — editing `spec` or other repositories, their settings, anything that costs money — only after a warning (what, where, cost) and the author's confirmation. `spec/` is read only here; a task of its own moves it to a new tag.
   - Exception, no confirmation: the task's checks by R146 spec — browser: Playwright in the devcontainer, custom-token sign-in of a stand-in user on the dev project; AI: the `tutor` stub in `api` dev mode. A real AI call only where the task checks the model's words: dev key under the author's spend limit at the provider, within the task's ceiling, or with confirmation where it has none; «Сделано» names the provider, the calls and the cost.
   - Nothing on devices: engines are Chromium, WebKit and Firefox in CI; what only a device shows — pen, device voices, VoiceOver and TalkBack, the home screen — stays a row of [run/DEBT.md](run/DEBT.md) until the author decides otherwise.
5. Facts on AI providers, laws, prices and other platforms — from the primary source, with quote and date read (RULES.md). What is not confirmed is not stated: it becomes an open question.
6. End with the checks of what the task changed. Documents: `just -f spec/justfile docs-check-repo` (R32 spec; what it checks — CLAUDE.md, Environment; R36). Then `/code-review high`, one round: fix the findings, do not rerun. A small change may skip the review even if «Как проверить» names it.
7. Do not commit: the author commits after review. A run of spec's orchestrator, «Сделай эпик E…», is the author's word otherwise: its executor works by the contract `spec/tools/orchestrator/executor.md` — commits, pushes and opens the pull request, never merges (R156 spec). Where the current epic's decision says so, the executor takes each task through its own branch and pull request to the merge (R30, decision 5; R37, decision 3).
8. End with a short report in Russian: what was done, files, how to check, what the review found or why there was none, what is open and what goes «в реестр spec».
9. When done: table 1 `[ ]` → `[x]`; append to the body «Сделано (дата)» — what was done, how checked — and «Открыто». A question for the author from «Открыто» is also a row of [run/QUESTIONS.md](run/QUESTIONS.md); what is put off, a row of [run/DEBT.md](run/DEBT.md).
   - Then move the body verbatim, relative links fixed, to the epic archive `docs/plan/<epic>.md`, with a gist sentence — an unheaded list item at the start of the receiving section: «T01 — сделана <дата>: <итог одной фразой>». Neither stays here, and a search by heading finds one body, in the archive (R153 spec).
   - A task outside an epic leaves `run/TASKS.md` for the current epic's archive, section «Задачи вне эпика, сделанные за его время»; with no current epic, to the last done epic's archive, section «Задачи вне эпика, сделанные после него» (R19 spec and its 2026-11-13 amendment; R35).
10. Language — RULES.md, Hard rules (R154 spec): this plan and `run/` in English; the current epic's description and rows, the bodies of tasks, the archives `docs/plan/` and the reports in Russian; `docs/` in English (R36).
11. A task larger than one review (≈ 30–60 minutes of reading): propose a split before starting.
12. Before a new task, close the findings on the previous one, once SonarCloud, CodeQL and Scorecard run here (R30 spec); an allowed departure from rule 1. SonarCloud sees merged code: its verdict on the previous task is in [New Code](https://sonarcloud.io/summary/new_code?id=altnovatrail_learner-model); CodeQL and Scorecard — the Security tab, `gh api 'repos/altnovatrail/learner-model/code-scanning/alerts?state=open'`. Fix in code first; "false positive" is the last resort, with the reason in words; what a change cannot close stays open — shared reasons in R33 spec, own ones in the log, like R51 canvas-api. Never skip silently. Report these fixes in a separate paragraph: the author commits them apart from the task.

**Statuses:** `[ ]` not started; `[x]` done, the executor marks it; `[-]` not done — the author's decision, reason in the task; `[?]` deferred — the author's decision, reason and when to return in the task. A task that fails the author's review returns to `[ ]` with what to fix.

**Who:** «Claude» — the executor does it, the author reviews; «Claude, проверяю я» — texts the author reads line by line; «я + Claude» — the author does some steps: decisions, sign-ins, GitHub settings, confirmations by rule 4.

**Plan:** «да» — plan mode first: decisions missing from the documents, or larger than one review (rule 11); «—» — run at once.

**Task body** (Russian): «Цель», «Сделать», «Как проверить», «Разделы»; as needed «Не делать», «Решить при планировании», «Мой шаг», «Стоимость», «Риски»; from planning «Почему», «Решено автором (дата)», «Выяснилось при планировании»; when done «Сделано (дата)», «Открыто» (rule 9).

**Citations** — RULES.md, Citations. A task — «T03», from another repository «T03 learner-model»; one not here is in `run/TASKS.md` or an archive: `grep -n '^#### T03\.' RUN.md run/TASKS.md docs/plan/*.md`. Task numbers are global and never reused: take the one free when writing it down.

## Epics

- **Table 1** is the current epic: its tasks in detail, run one at a time.
- **Next epic.** Once table 1 is all checked, the author picks the next epic from [run/EPICS.md](run/EPICS.md) and splits it with Claude in plan mode into tasks of one review each, into table 1 and this file.
- **A done epic** goes to `docs/plan/<epic>.md` (R59 canvas-api): bodies and gist sentences are there already (rule 9); move its description and table of tasks verbatim; questions and debt stay in `run/`; until the next epic, table 1 reads «Current epic — none» and names it as the last done epic; measure the epic's sessions — `just -f spec/justfile tokens-report --since <first day of the epic> --row` — as a row of the table of measurements of `docs/measurements.md`, which the first measurement starts after measurements 1.9 spec; a new analysis by the threshold of R145 spec. Never read an archive whole: take the task cited.

## Table 1. Current epic — E8. Круг справедливости по классу

Развёрнут 2026-11-16 задачей T58 по решению автора 2026-11-15 (R156 spec, «Row 4.17»; R40): база не держит критерий справедливости по классу (R58 sim), а его ждёт первый ребёнок (dpia 8 spec). E7 — в архиве [docs/plan/E7.md](docs/plan/E7.md): T51 сделана, T52–T55 — `[?]`.

**Результат** (строка 4.17 COORDINATION.md spec):
- записано решением журнала, почему со старта на 2,5 выше правды модель сходит к ней дольше — за 83 ответа в среднем, — чем поднимается со старта на 2,5 ниже — за 33 (R58 sim), и что из этого идёт в кандидаты;
- суд E10 `sim` — по коммиту кандидата, с `CANDIDATE_LEARNER_MODEL` (строка 4.18 COORDINATION.md spec) — прошёл по кандидатам, что поднимают `M11.same` у `D-grade5` и `D-grade1` к 0,90 при критерии модели как есть. Каждый кандидат — черновой pull request; кандидат с иным стартом по классу, чем полоса класса (concepts 11 spec; R51 spec), — ещё и строка «в реестр spec». Выбранный вышел релизом по R103 spec, или выход записан;
- итог круга — строкой «в реестр spec» для dpia 10.2 spec и dpia 8 spec: `M11` выбранного, или выход с числами, по которым автор решает об остаточном риске в dpia 7 spec (поправка 2026-11-15 к R147 spec).

| Статус | ID | Задача | Кто | План | Зависит от |
|---|---|---|---|---|---|
| [ ] | T59 | Разбор: почему со старта выше правды оценка сходит дольше | Claude | да | — |
| [ ] | T60 | Гейт на теге `sim` с критерием суда круга | Claude | — | критерий суда E10 записан в `sim` (строка 4.18 COORDINATION.md spec) |
| [ ] | T61 | Разведка круга | я + Claude | да | T59, T60; В51 |
| [ ] | T62 | Кандидаты — черновые pull request | Claude | да | T61 |
| [ ] | T63 | Суд и выпуск выбранного | Claude | — | T62; суд E10 `sim` по коммиту кандидата (строка 4.18 COORDINATION.md spec); кандидату с иным стартом по классу — ответ `spec` на его строку «в реестр spec» |
| [ ] | T64 | Итог круга и архив эпика | Claude | — | T63 |

**Почему такой порядок.**
- Разбор — первым и параллельно с гейтом, как T42 в E5: «почему» R58 sim оставил `learner-model` — «Why the start above the truth is left more slowly is not measured here; it is a question for `learner-model`», — и без причины разведка выбирала бы числа наугад. Тега суда разбору не нужно: он гоняет ячейки `D-` на последнем теге `sim`, а `v0.1.21` даёт те же биты, что `v0.1.20`, чей `paths.json` называет файл базы `fairness-1` (R58, R60 sim).
- Гейт — до разведки, как T51 в E7: разведка и суд идут на одном теге, а красный «Gate» кандидата должен значить сдвиг меры, а не новый тег. Тег — тот, на котором `sim` записал критерий E10, до первого кандидата (R103 spec; R49 sim).
- Разведка — после обоих, на размере решающего прогона (поправка 2026-11-10 к R24).
- Кандидаты и суд — как T53 и T54 в E7. Суд ждёт не релиза `sim`, а его суда по коммиту кандидата: `sim` гонит черновик из его коммита (строка 4.18). Тег `sim` с полосами выбранного нужен только мержу выбранного (R103 spec).
- Итог — своей задачей: строку для dpia ждёт первый ребёнок, и она выходит и при выходе круга, когда T63 ничего не выпускает.

**Решено автором (2026-11-15)** (R156 spec, «Row 4.17»; RUN.md spec, E14): круг справедливости — новый текущий эпик; T52–T55 E7 — `[?]` по R59 sim, E7 — в архив с T51. Отвергнуто: переписать E7 под круг справедливости; оставить выбор сессии `learner-model`. План эпика — исполнителя: в прогоне оркестратора автор планов не утверждает, а вопросы плана — строки таблицы 3 (R156 spec; R40).

## Epic E8. Круг справедливости по классу

**Почему эпик такой.** Оценка `fairness-1` `sim` на базе — `learner-model` v0.1.8 и `engine` v0.1.6 — нашла: у учеников одного навыка под классами 1 и 5 олимпиадной математики, где модель стартует на полосу ниже и выше правды, оценка к 50-му ответу сходится с оценкой под классом 3 у 0,812 [0,795; 0,829] и 0,5525 [0,53; 0,5735] учеников при границе 0,90; в среднем — за 32,5 и 83,3 ответа; полосы задач сходятся (R58 sim; dpia 10.2 spec, item 10). Ученик, чей класс выше того, что он умеет, дольше всех читается слишком высоко. `fair` держит кандидата «не хуже базы», но не награждает за лучшее (R55 sim), так что круг идёт со своей целью — судом E10 `sim` (R59 sim; строка `sim` В22, вариант (а)).

**Выяснилось при записи** (2026-11-16; `sim` `main` c5e9a4c) — R40, «What was found»:
- E10 `sim` «Суд круга справедливости модели» в его таблице 2 ждёт, пока `learner-model` возьмёт круг. Критерий — критерий модели как есть, цель — `M11.same` у `D-grade1` и `D-grade5`, своя папка и свой отложенный набор, как у `engine-lesson-1` (R49, R52 sim); `sim` записывает его до первого кандидата. Последний тег `sim` — `v0.1.21`, а «Gate» здесь — на `v0.1.14`;
- старт модели — θ₀, смещение полосы класса, 2,5 · (p − 1); пока вес ответов меньше 5, θ — максимум апостериорной вероятности серии проб при нормальном априорном с центром θ₀ и σ₀ = 2,5; затем шаг Elo с K_θ = max(0,2 / (1 + 0,05 · n); 0,05) (R03; `trial.go`, `step.go`, `model.go`). Это числа без нового поля: σ₀, вес конца серии проб, K_θ. Старт в полосе класса — правило R51 spec, «A learner of a grade starts in the band of their grade», и иной старт меняет `spec`;
- σ₀ 3,0 уже пробовала разведка T30: на 300 учениках не лучше базы по ошибке оценки (R24; Д5); `M11` тогда не было.

**Не входит:**
- `engine`: его выбор задачи тоже двигает оценку, и `fair` держит обе библиотеки (R55 sim), но этот круг — модели; кандидат `engine` идёт своей причиной в `engine`;
- числа забывания — T52–T55 E7, `[?]` (R40); Д6, Д7 и в Д5 всё, кроме чисел старта;
- код `sim`: критерий и суд — E10 `sim`.

**Ветки задач** (как в E7: R37, решение 3):
- задача — своя ветка от `main` и свой pull request в `main`. В прогоне оркестратора `spec` мержит оркестратор на зелёных проверках, а исполнитель не мержит (R156 spec); вне прогона исполнитель доводит pull request до мержа сам, с «Gate» и `codecov/patch` (R30, решение 5). Правки по правилу 12 — отдельным коммитом;
- черновик кандидата T62 — своя ветка от `main` и черновой pull request; суд `sim` читает его коммит. Выбранного T63 переводит из черновика в готовый, и его мерж merge-коммитом выпускает патч линии `0.1`; остальные закрываются с причиной;
- что пришло в `main` до мержа кандидата, вливается в его ветку merge-коммитом, а не перебазированием: SHA коммита кандидата называют журнал и суд.

**Общее для задач:**
- перед задачей — правило 12; документы — `just -f spec/justfile docs-check-repo`, зелёный;
- кандидат меняет числа или правило без нового поля; новое поле — новый вид состояния и пересчёт кэша `api` (R08; R80 spec) — только решением автора, строкой таблицы 3 при планировании T62. Чем меньше констант, тем вероятнее выбор (R04 `sim`);
- что меняет кандидат, кроме числа: тесты, где число стоит руками; золотые файлы — `LEARNERMODEL_UPDATE=1 go test .`; затронутые векторы standalone снимаются с записью (R22);
- «Gate» на кандидате краснеет, как только число гейта выходит из полосы, — к лучшему или к худшему (R20). Зелёный «Gate» на старых полосах — не повод сливать: мерж ждёт суда и полос `sim` (R103 spec), и это сказано в описании pull request;
- код — с тестами, `codecov/patch` — не меньше 90 % (R48 spec).

#### T59. Разбор: почему со старта выше правды оценка сходит дольше

- **Цель:** записано решением журнала, с замером, почему модель со старта на 2,5 выше правды сходит к ней за 83 ответа в среднем, а со старта на 2,5 ниже поднимается за 33 (R58 sim), и какие числа или правила из этого идут в разведку T61.
- **Почему:** `sim` это не замерял и оставил вопрос `learner-model` (R58 sim). Так T42 в [docs/plan/E5.md](docs/plan/E5.md) нашёл причину провала первого суда до разведки E5 (R31).
- **Решить при планировании:**
  - какие объяснения проверить замером: серия проб — максимум апостериорной вероятности с σ₀ = 2,5 и угадыванием на выборе (`trial.go`); шаг после неё — K_θ и S − P, когда правда даёт P много ниже или выше оценки; доля, которую забирают поправки навыков δ, — `M11` читает только θ; задачи, которые `engine` даёт по завышенной оценке: на трудной задаче ответ правды близок к угадыванию;
  - на чём: ячейки `D-grade1`, `D-grade3` и `D-grade5` на последнем теге `sim`, база — по его `go.mod`, след θ после каждого ответа (R56 sim); варианты — в клоне в папке сессии, в репозиторий не идут; размер — сколько нужно, чтобы различить объяснения, не больше решающего прогона;
  - если причина — в выборе `engine`, а не в модели: строка «в реестр spec» для `engine` и `sim`, а круг идёт тем, что может модель.
- **Сделать:** замеры; решение журнала — причина, замер, как его повторить, с SHA-256 файлов, и что идёт в T61 — числа, правила и меры; раздел 3.1 [docs/architecture.md](docs/architecture.md), если он говорит иначе; строка Д5.
- **Не делать:** чисел модели в `main` не менять.
- **Как проверить:** решение журнала называет прогоны так, что их можно повторить; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R03, R24, R31; R55, R56 и R58 `sim`; dpia 10.2 spec; Д5.

#### T60. Гейт на теге `sim` с критерием суда круга

- **Цель:** job «Gate» стоит на теге `sim`, на котором записан критерий суда E10, — последнем к началу задачи, — и на коде `main` он зелёный.
- **Почему:** разведка T61 и суд E10 идут на одном теге, а «Gate» — на теге, на котором `sim` судит (R20; T51 в [docs/plan/E7.md](docs/plan/E7.md)).
- **Сделать** (как T51): `SIM_TAG` и `SIM_COMMIT` в `.github/workflows/ci.yml` — тег и SHA его коммита, проверенные `git ls-remote`, `gh release view` и `gh release verify`; что изменилось у вызывающего — рецепт `just gate`, `gate.json`, «Taking a release» README `sim`; `just gate exact` из клона `sim` на коде `main`; поправка к R20; разделы 1 и 3.6 [docs/architecture.md](docs/architecture.md), README.md; запись источника в `docs/sources.md`, как sources 6.11.
- **Не делать:** чисел модели не трогать.
- **Как проверить:** «Gate» на pull request — зелёный, числа — в полосах; `just -f spec/justfile docs-check-repo`.
- **Разделы:** R20; R49, R52 и R60 `sim`; строка 4.18 COORDINATION.md spec; sources 6.11.

#### T61. Разведка круга

- **Цель:** выбраны до трёх кандидатов суда E10 — числа или правило старта и шага, цель — `M11.same` у `D-grade5` и `D-grade1`, число изменённых констант у каждого, — и выбор записан решением журнала с тем, что показала разведка.
- **Почему:** как T43 в [docs/plan/E5.md](docs/plan/E5.md) — на размере решающего прогона и на теге суда, — но с причиной T59. Суд держит кандидата всеми ограничениями критерия модели как есть (строка 4.18 COORDINATION.md spec): кандидат, который поднимает `M11.same`, но срывает «не хуже» где-то ещё, суда не пройдёт.
- **Решить при планировании:**
  - наборы — из причины T59: σ₀, вес конца серии проб, K_θ и их пары; правила остановки;
  - иной старт по классу — например, ниже полосы класса: R51 spec и concepts 11 spec велят начинать в полосе класса, так что такой кандидат — ещё и строка «в реестр spec» с числами разведки, до его черновика;
  - меры — `M11.same`, `M11.answers` и `M11.band` у `D-grade1` и `D-grade5` против цели E10, и все «не хуже», полы и `fair` критерия модели; где сорвалась σ₀ 3,0 в T30 (R24).
- **Мой шаг:** бюджет прогонов и время машины — В51.
- **Сделать:** прогоны на теге гейта T60, база — без кандидата, по `go.mod` `sim`; варианты — коммиты поверх `main` в клоне в папке сессии, в репозиторий не идут; разбор — `sim judge` критерием E10 и `sim compare` против базы; решение журнала — кандидаты, их меры, почему и как повторить прогоны, с SHA-256 файлов; строка Д5; раздел 2.1 [docs/architecture.md](docs/architecture.md), «Deferred».
- **Как проверить:** решение журнала называет прогоны так, что их можно повторить; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Риски:** если ни один набор не держит всех ограничений, кандидатов нет — это выход круга: исполнитель предлагает его строкой таблицы 3, T62 и T63 становятся `[-]` решением автора, с записью, а T64 несёт в dpia числа разведки.
- **Разделы:** R03, R22, R24, R32; R49, R52, R55, R58 и R59 `sim`; строка 4.18 COORDINATION.md spec; R51 spec; concepts 11 spec; Д5.

#### T62. Кандидаты — черновые pull request

- **Цель:** у каждого кандидата T61 — свой черновой pull request, готовый к суду E10 `sim`.
- **Почему:** суд `sim` гонит кандидата из его коммита, `CANDIDATE_LEARNER_MODEL` (строка 4.18 COORDINATION.md spec; T44 в [docs/plan/E5.md](docs/plan/E5.md)).
- **Решить при планировании:** какие тесты стоят руками на изменённых числах; если кандидату всё же нужно новое поле — вид состояния 2 и пересчёт кэша `api` (R08; R80 spec) — строкой таблицы 3.
- **Сделать** — в каждом pull request: числа или правило кандидата; тесты, где они стоят руками; золотые файлы — `LEARNERMODEL_UPDATE=1`; затронутые векторы standalone сняты, с записью (R22); раздел 3.1 [docs/architecture.md](docs/architecture.md), поправка к R03 и строка Д5; прогон `just decisive` на коммите кандидата сверен с прогоном разведки; описание называет для суда имя кандидата, цель, число изменённых констант и нет ли нового поля (R25 `sim`) и говорит, что мерж ждёт суда и полос `sim`, каким бы ни был «Gate».
- **Не делать:** не сливать; базу — `main` — не трогать.
- **Как проверить:** в каждом pull request `just ci-lint` и `ci-test` — зелёные, кроме «Gate»; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг на каждый.
- **Разделы:** R03, R12, R20, R22; R25 `sim`; строка 4.18 COORDINATION.md spec.

#### T63. Суд и выпуск выбранного

- **Цель:** суд E10 `sim` по кандидатам прошёл, и его итог стоит в коде и документах.
- **Почему:** R103 spec — порядок релизов: полосы в `sim` записываются на последнем коммите кандидата раньше мержа здесь. Задача ждёт суда `sim` по коммиту кандидата, а не тега `sim` под круг: суд гонит черновик из его коммита (строка 4.18 COORDINATION.md spec).
- **Сделать** (как T45 в [docs/plan/E5.md](docs/plan/E5.md)): прочитать отчёт суда и подтверждение на отложенном наборе E10; если выбран кандидат и подтверждён — на его pull request `SIM_TAG` и `SIM_COMMIT` на тег `sim` с его полосами, «Gate» зелёный, pull request — из черновика в готовый, и его мерж выпускает патч линии `0.1`; кандидат с иным стартом по классу — только после ответа `spec` на его строку; остальные закрываются с причиной; если нет — все закрываются, причина и что показал суд записываются; поправки к R03 и к решению T61.
- **Как проверить:** CI pull request выбранного, с «Gate», — зелёный; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Риски:** с релизом выбранного `api` пересчитывает кэш — его пометка — версия модели (R80 spec), — и оценки учеников, чей класс не совпадает с их уровнем, меняются.
- **Разделы:** R03, R12, R20, R22; R25, R49 и R52 `sim`; R103 spec; строка 4.18 COORDINATION.md spec.

#### T64. Итог круга и архив эпика

- **Цель:** итог круга — у `spec` строкой «в реестр spec» для dpia 10.2 spec и dpia 8 spec, раздел 3 [docs/architecture.md](docs/architecture.md) говорит то, что построено, а E8 — в архиве.
- **Почему:** первый ребёнок ждёт, пока база не держит критерий или пока автор не примет остаточный риск в dpia 7 spec с числами этого круга (поправка 2026-11-15 к R147 spec; dpia 8 spec).
- **Сделать:** строка «в реестр spec» — выбранный и его `M11` у `D-grade1` и `D-grade5` с интервалами из отчёта суда и держит ли он критерий R55 `sim`, или выход с числами разведки и суда, — для item 10 dpia 10.2 spec и условия dpia 8 spec; что меняет карточку architecture 3.8 spec или concepts 11 spec — тоже строкой «в реестр spec»; раздел 3 — сверить с кодом; долги [run/DEBT.md](run/DEBT.md) — Д5 по итогу; строки таблицы 3 эпика — что с ними; T52–T55 E7 — по-прежнему `[?]` или вернуть, если причина пришла; замер токенов сессий эпика — строкой [measurements 4](docs/measurements.md); перенести E8 в архив `docs/plan/E8.md`.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R03; R55 и R58 `sim`; поправка 2026-11-15 к R147 spec; dpia 7 spec, dpia 8 spec и dpia 10.2 spec; architecture 3.8 spec; concepts 11 spec; таблицы 2–4.

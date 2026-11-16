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

## Table 1. Current epic — E7. Третий круг кандидатов модели

Развёрнут 2026-11-13, когда E5 кончился релизом `theta`, `v0.1.8` (R30; R37). E5 — в архиве [docs/plan/E5.md](docs/plan/E5.md).

**Результат:**
- «Gate» стоит на последнем теге `sim`, где база — релиз `v0.1.8`, и разведка и суд идут на том же теге;
- третий суд `sim` по кандидатам модели прошёл. Кандидатов до трёх, это числа забывания R04 на базе с `theta`: основание деления τ между 0,4 и 0,5 или рост τ, который замедляется с самой τ, как у FSRS (sources 3.1). Цель — сохранить выигрыш `forget-set` в группе `forgets` без провалов вне `F-` популяций (R32). Каждый кандидат — черновой pull request; затем суд `model-3` и подтверждение выбранного на нетронутом `held-out-model-2` (R46 `sim`). Выбранный вышел релизом, или выход записан;
- Д6 говорит, что показал круг.

| Статус | ID | Задача | Кто | План | Зависит от |
|---|---|---|---|---|---|
| [x] | T51 | Гейт на теге `sim` с базой на `v0.1.8` | Claude | — | — |
| [ ] | T52 | Разведка третьего круга | я + Claude | да | T51 |
| [ ] | T53 | Кандидаты — черновые pull request | Claude, проверяю я | да | T52 |
| [ ] | T54 | Суд `model-3` и выпуск выбранного | я + Claude | — | T53; суд `model-3` RUN.md `sim` (В50; ответ — поправка 2026-11-16 к R37) |
| [ ] | T55 | Итог эпика | Claude, проверяю я | — | T54 |
| [x] | T56 | Переход на тег `spec` v0.18.3 — вне эпика | Claude | — | — |
| [x] | T57 | Раскладка плана по R153–R155 spec — вне эпика | Claude | — | T56 |

**Почему такой порядок.**
- Гейт — до разведки, как T36 в E4 и T41 в E5: разведка и суд идут на тех размерах и полосах, на которых `sim` будет судить, а красный «Gate» кандидата должен значить сдвиг меры, а не новый тег.
- Разбора, как T42 в E5, нет: причину провалов `forget-set` вне `F-` уже назвал R32 — рост τ при g 10 (R37, решение 2).
- Кандидаты, суд и итог — как T44–T46 в E5.
- От T51 до выбора суда тег `sim` не меняется: числа разведки на этом теге — те же, что даст суд на рабочих учениках (R32). Что круг `engine` «Закрепление в уроке» ждёт своего тега `sim` и суда в то же время, — строка В50 «в реестр spec»; ответ `spec` — поправка 2026-11-16 к R37.

**Решено автором (2026-11-13), при развёртывании E7** (R37):
1. E7 — следующий эпик: «разворачивай задачи E7»; состав и порядок — по строке E7 таблицы 2.
2. Ветки — по задаче, как в E4 и E5 с T44 (раздел «Ветки задач» ниже).

## Epic E7. Третий круг кандидатов модели

**Почему эпик такой.** С правилом θ числа забывания больше не задерживают догонку после перерыва, но ни один набор разведки T43 не держит всех ограничений. `forget-set` выигрывает в группе `forgets` +0,277 и даёт семь провалов: пять вне `F-`, где правило θ не действует, — `E1.corridor` на G2, G3 и S-adult, `M1.rms.200` на G3-drop и S-scale1.5, — и два на `F-none`: `M1.rms.200` олимпиадной и `M5.caught` школьной. Основание деления τ само по себе — 0,35 и 0,4 — срывает по одной проверке (R32). Числа и их источники — Д6.

**Выяснилось при развёртывании** — R37, «What was found»: тег `sim` с базой на релизе — `v0.1.14`, а «Gate» здесь — на `v0.1.13`; задачи под суд `model-3` у `sim` нет, и `engine` ждёт от `sim` своего тега и суда (В50); рост τ, зависящий от самой τ, обходится без нового поля. Решающий прогон `sim` — 2 000 рабочих учеников на ячейку, рецептом `just decisive engine <файл>` с `CANDIDATE_LEARNER_MODEL` — папкой варианта, около 38 минут на свободной машине в 20 ядер, до 62 на занятой, по одному (R32). Зерно `held-out-model-2` `sim` записывает поправкой к своему R04 до первого прогона кандидата (R46 `sim`).

**Не входит:**
- пара с `engine` — порог «освоено», цель коридора и граница ½ (Д7, Д16) — эпик E6 в таблице 2;
- Д5, ρ, τ₀ и вес ответа с помощью: новых фактов о них нет (R24, R32);
- код `sim` и `engine`: суд — задача `sim`, закрепление в уроке — E7 `engine`.

**Ветки задач** (R37, решение 3; как в E5 с T44):
- задача — своя ветка от `main`, свой pull request в `main` и мерж, который исполнитель доводит до зелёного CI сам, с «Gate» и `codecov/patch` (R30, решение 5). Правки по правилу 12 — отдельным коммитом;
- черновик кандидата T53 — своя ветка от `main` и черновой pull request в `main`. Суд `sim` читает коммит кандидата. Выбранного T54 сливает в `main` merge-коммитом — этот мерж выпускает патч линии `0.1`, — а остальные закрываются с причиной;
- если в `main` до мержа кандидата пришло что-то ещё, оно вливается в ветку кандидата merge-коммитом, а не перебазированием: SHA коммита кандидата называют журнал и суд.

**Общее для задач:**
- перед задачей — правило 12; документы — `just -f spec/justfile docs-check-repo`, зелёный;
- кандидат меняет числа или правило без нового поля; новое поле — новый вид состояния и пересчёт кэша `api` (R08; R80 spec) — только с решением автора при планировании T53. Чем меньше констант, тем вероятнее выбор (R04 `sim`);
- что меняет кандидат, кроме числа: тесты, где число стоит руками; золотые файлы — `LEARNERMODEL_UPDATE=1 go test .`; затронутые векторы standalone снимаются с записью (R22);
- «Gate» на кандидате краснеет, как только число гейта выходит из полосы, — к лучшему или к худшему (R20). Зелёный «Gate» на старых полосах — не повод сливать: мерж ждёт суда и полос `sim` (R103 spec), и это сказано в описании pull request;
- код — с тестами, `codecov/patch` — не меньше 90 % (R48 spec).

#### T52. Разведка третьего круга

- **Цель:** выбраны до трёх кандидатов третьего суда — числа или правило, группа оценки, цель и число изменённых констант у каждого, — и выбор записан решением журнала с тем, что показала разведка.
- **Почему:** как T43 в [docs/plan/E5.md](docs/plan/E5.md) — на размере решающего прогона и на том же теге, что суд, — но на базе с `theta` и с числами забывания, которые R32 оставил следующему кругу.
- **Решить при планировании:**
  - основание деления τ между 0,4 и 0,5: какие точки (0,45 — первой) и правила остановки. `M5.caught` на `F-none` школьной шумная — стандартная ошибка разности 0,6–0,9 ответа, допуск суда 2,58–3,75 (R32), — так что граница между 0,4 и 0,5 по одной точке не видна;
  - рост τ, который замедляется с самой τ, как у FSRS (sources 3.1): вид — например, g · (1 − R) · (τ / τ₀)^(−k) вместо g · (1 − R) — и его константы; с ним и g меньше 15 — набор, который ускоряет забывание у малых τ, не раздувая больших;
  - что не повторять без новых фактов: 0,35 и 0,4 сами по себе, g 10 без правила роста, ρ, τ₀, вес ответа с помощью и Д5 (R24, R32);
  - меры — группа `forgets` и то, что сорвалось в T43: `M5.caught` на `F-none` школьной, `M1.rms.200` на G8, G3-drop, S-scale1.5 и `F-none` олимпиадной, `E1.corridor` на G2, G3 и S-adult; и `M3.settle` на `F-break` — правило θ не должно потерять догонку;
  - сколько прогонов: в T43 автор дал около десяти (R32).
- **Мой шаг:** бюджет прогонов и время машины.
- **Сделать:** до первого прогона — проверить, что `sim` записал зерно `held-out-model-2` (R46 `sim`), а если нет — остановиться и спросить; прогоны на теге `sim` гейта T51, база — без кандидата, по `go.mod` `sim`; варианты — коммиты поверх `v0.1.8` в клоне в папке сессии, в репозиторий не идут; разбор — `sim judge` критерием модели и `sim compare` против базы; решение журнала — кандидаты, их меры, почему и как повторить прогоны; строка Д6, раздел 2.2 [docs/architecture.md](docs/architecture.md), «Deferred».
- **Как проверить:** решение журнала называет прогоны так, что их можно повторить, с SHA-256 файлов; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Риски:** если ни один набор не держит всех ограничений, кандидатов нет — это выход круга: исполнитель предлагает его, а T53 и T54 становятся `[-]` моим решением, с записью; Д6 тогда говорит, что показала разведка.
- **Разделы:** R04, R22, R24, R31, R32; R35, R37 и R46 `sim`; 3.6 `sim`; Д6; sources 3.1.
- **Выяснилось при планировании** (2026-11-14; `sim` `main` 9c8f782 = `v0.1.14`):
  - правило 12 по T51: CI `main` d2d7ccd зелёный; SonarCloud — ворота OK, находок и hotspot нет; CodeQL — пусто; Scorecard — те же шесть из общего списка R33 spec. Править нечего;
  - зерна `held-out-model-2` в `sim` нет: R04 `sim` (поправка 2026-11-13 «the base on `theta`») пишет его поправкой «before the first run of its candidate», а R12 `sim` знает только `held-out-model-1` и `held-out-engine-1`;
  - истина всех популяций вне `F-` забывает по базовому закону модели — g = 15 и деление пополам (architecture 3.1 `sim`); закон меняют только `F-exp`, `F-C1`, `F-growth5`, `F-none` и `F-break`. Поэтому кандидат с другими g или делением расходится с истиной вне `F-`, а рост τ, замедляющийся с τ, при g = 15 совпадает с ней при τ = τ₀ и расходится лишь у больших τ, которые набираются в 180-дневных прогонах `F-`;
  - суд `model-3` сравнивает с базой `theta`, а T43 — с `v0.1.6`: на `F-` разность `thh040-20-14` несла цену самого `theta` — `M5.caught` на `F-none` школьной +1,26 из +3,07, — и против новой базы основание 0,4 может держать. Срыв 0,35 — `M1.rms.200` на G8, вне `F-`, где `theta` не действует, — останется;
  - у FSRS рост стабильности умножается на S^(−w9), w9 = 0,1666 по умолчанию (py-fsrs, e40a5a7); `detmath.Pow` берёт отрицательную степень, так что k варианту не нужен новый код `detmath`;
  - база разведки должна дать SHA-256 `base` из `paths.json` `sim` `v0.1.14` — `7d3cb52e72cf4c55…bf16`; файлов прогонов T43 не осталось.
- **Решено автором (2026-11-14):**
  1. зерно `held-out-model-2` — ждать `sim`: его записывает сессия `sim`, а прогоны идут, только когда оно в `main` `sim`; подготовка без прогонов — клоны и коммиты вариантов — до того;
  2. бюджет — около десяти прогонов, база — первый из них;
  3. этап 1 — основание деления τ 0,4 против новой базы и 0,45, рост τ, замедляющийся с τ, — g · (1 − R) · (τ/τ₀)^(−k) при g 15 — с k 1/6 и 1/3; 0,35 не повторять. Этап 2 — до пяти прогонов по правилам остановки: 0,425; k 1/2, 1/4 или 1/12; лучший k с g 12,5; набор деления и k;
  4. после тега `v0.1.15` `sim` (ниже) — первым базу `sim` двигает суд закрепления в уроке `engine-lesson-1`, как в таблице 2 RUN.md `sim`, а T52 ждёт его итога: выбран вариант — сначала «Gate» на новый тег `sim`, как T51, и разведка на новой базе; выход — разведка на `v0.1.15`. Зерно `held-out-model-2` по-прежнему нужно до первого прогона.
- **Выяснилось при исполнении** (2026-11-14; `sim` `main` 6310470):
  - `sim` выпустил `v0.1.15`, 52a3658: шаг закрепления в уроке (R48 `sim`) и критерий суда `engine-lesson-1` с зерном `held-out-engine-lesson-1` 20261113 (R49 `sim`). На базе шаг даёт те же биты: `engine` v0.1.6 в конце сегмента сначала продолжает открытый урок, что бы ни оценила модель. Файл базы отличается ключами справки `E6` — `base` в `paths.json` теперь `2e53e18c08df5016…9aad`, — пути, `gate.json` и тело рецепта `decisive` те же;
  - зерна `held-out-model-2` нет и на `main` 6310470; зерно `sim` — дата записи (раздел 3.8 `sim`). Эпика под суд `model-3` в таблице 2 `sim` нет: там E7 — суд `engine-lesson-1` — и E6, `engine-2`, а по R35 `sim` второй суд идёт на базе с выбором первого. Кандидата закрепления в уроке у `engine` ещё нет — открытых pull request нет;
  - варианты этапа 1 собраны в клоне в папке сессии поверх `v0.1.8` — `h040` cb8fa1d, `h045` d8475c3, `k6` 9d73309, `k3` 40b63a3; прогонов не было.

#### T53. Кандидаты — черновые pull request

- **Цель:** у каждого кандидата T52 — свой черновой pull request, готовый к суду `model-3` `sim`.
- **Почему:** суд ждёт pull request кандидата в библиотеке (R35 `sim`; T44 в [docs/plan/E5.md](docs/plan/E5.md)).
- **Решить при планировании:** как правило роста τ стоит в `stabilityAfter` и какие тесты стоят руками на нём; если кандидату всё же нужно новое поле — вид состояния 2 и пересчёт кэша `api` (R08; R80 spec).
- **Сделать** — в каждом pull request: числа или правило кандидата; тесты, где они стоят руками; золотые файлы — `LEARNERMODEL_UPDATE=1`; затронутые векторы standalone сняты, с записью (R22); раздел 3 [docs/architecture.md](docs/architecture.md), поправка к R04 и строка Д6; прогон `just decisive` на коммите кандидата сверен с прогоном разведки; описание называет для суда имя кандидата, группу, цель, число изменённых констант и нет ли нового поля (R25 `sim`) и говорит, что мерж ждёт суда и полос `sim`, каким бы ни был «Gate».
- **Не делать:** не сливать; базу — `main` — не трогать.
- **Как проверить:** в каждом pull request `just ci-lint` и `ci-test` — зелёные, кроме «Gate»; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг на каждый.
- **Разделы:** R04, R12, R20, R22, R32; R25 и R35 `sim`.

#### T54. Суд `model-3` и выпуск выбранного

- **Цель:** третий суд `sim` по кандидатам модели прошёл, и его итог стоит в коде и документах.
- **Почему:** R103 spec — порядок релизов: полосы в `sim` записываются на последнем коммите кандидата раньше мержа здесь.
- **Мой шаг:** в `sim` — задача суда `model-3` (В50; ответ `spec` — поправка 2026-11-16 к R37), чтение его отчёта и подтверждения на `held-out-model-2`, мерж pull request с полосами.
- **Сделать** (как T45 в [docs/plan/E5.md](docs/plan/E5.md)): если выбран кандидат и подтверждён — на его pull request `SIM_TAG` и `SIM_COMMIT` на тег `sim` с его полосами, «Gate» зелёный, мерж в `main` — он выпускает патч линии `0.1`, — остальные закрываются с причиной; если нет — все закрываются, причина и что показал суд записываются, Д6 остаётся; поправка к R04 и к решению T52; строка «в реестр spec» — третий круг модели и его судьба, для `sim`, `engine` и `api`.
- **Как проверить:** CI pull request выбранного, с «Gate», — зелёный; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Риски:** релиз выбранного выходит мержем T54 в `main`, а с ним `api` пересчитывает кэш: пометка кэша — версия модели (R80 spec).
- **Разделы:** R12, R20, R22, R37; R35 и R46 `sim`; R103 spec.

#### T55. Итог эпика

- **Цель:** раздел 3 [docs/architecture.md](docs/architecture.md) говорит то, что построено после третьего суда, следующий эпик предложен, а E7 — в архиве.
- **Сделать:** сверить раздел 3 с кодом; что меняет карточку architecture 3.8 spec — строкой «в реестр spec»; долги [run/DEBT.md](run/DEBT.md) — закрыть, что закрыл суд; строки таблицы 3 эпика — назвать, что с ними; пересмотреть таблицу 2 — E6 и что ещё предложить; замер токенов сессий эпика — строкой в `docs/measurements.md`; перенести E7 в архив `docs/plan/E7.md` — своим pull request.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** раздел 3 [docs/architecture.md](docs/architecture.md); architecture 3.8 spec; таблицы 2–4; R37.

# TASKS.md — tasks outside an epic of `learner-model`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T01\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T65. Анализ контекста сессий E8 по порогу R145 spec

- **Цель:** понять, почему у сессий E8 контекст главного вызова на 35 % выше опорного, и решить, что из этого убрать.
- **Почему:** строка E8 в [measurements 4](../docs/measurements.md) — 173,3 тысячи на главный вызов против опорных 128,2 у E7, на 545 главных вызовах: выше 25 % порога R145 spec, а анализ по порогу — задача вне эпика (R155 spec). Записала T64 при архиве E8.
- **Сделать:**
  - `just -f spec/justfile tokens-report --since 2026-11-16T01:48:00Z --until 2026-11-16T16:17:00Z` — по сессиям, по инструментам, по файлам и по программам Bash — сверить с тем же отчётом по окну опорной строки E7;
  - назвать, что даёт рост: вывод Bash — 47,4 М из 100,2 М, перенесённых в контекст; чтения `docs/decisions.md` через Bash — 87 раз, 11,2 М — и RUN.md — 73 раза, 4,3 М; длинные сессии T61 и T59 — 187,7 и 165,6 тысячи на вызов, но вместе с вызовами субагентов, так что их главные вызовы выше: число по главным вызовам сессии — первым делом;
  - решить, что из этого — размер документов и их бюджеты в `docs-budgets.json`, порядок чтения в правилах или то, что задачи E8 читали прогоны `sim`, — и что делать;
  - итог — замером в `docs/measurements.md`; правка бюджетов или правил — решением журнала; причина в общих правилах или инструментах `spec` — строкой «в реестр spec» в [QUESTIONS.md](QUESTIONS.md).
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** R145 spec, R155 spec; measurements 1.9 spec; measurements 4.

#### T70. Переход на тег `spec` v0.18.14

- **Цель:** `learner-model` стоит на последнем теге `spec`, v0.18.14, — его правила, план и документы по этому тегу, — и остальные задачи эпика прогона COORDINATION.md (слово автора 2026-11-17) `spec` идут здесь на нём.
- **Почему:** R156 spec — параллельная работа волны 4 (слово автора 2026-11-17): пока `app` идёт эпиком E5, остальные репозитории переходят на последний тег `spec` и закрывают строки таблицы 4 `spec`, которые ждут этого перехода; переход — задача и решение журнала (RULES.md spec, Across repositories). Задача вне эпика или в текущем эпике — по правилам этого RUN.md. Записана по тексту прогона оркестратора (правило 6 COORDINATION.md spec).
- **Сделать:**
  1. `gh release verify v0.18.14 -R altnovatrail/spec`; подмодуль `spec/` — на v0.18.14; где репозиторий берёт модуль Go или пакет npm `spec` — та же версия, пакет npm — через `release-verify-npm` (R77 spec); тег — везде, где этот репозиторий его называет: CLAUDE.md, README, начало RUN.md, статус `docs/architecture.md`.
  2. Что теги после v0.18.6 до v0.18.14 несут этому репозиторию: `git -C spec log --oneline v0.18.6..v0.18.14`; строки таблицы 4 `spec` с именем репозитория в обратных кавычках — `grep -n '`learner-model`' spec/run/DEBT.md`; карточка этого репозитория в architecture 3 spec; RULES.md spec. Что из них берёт эта задача, а что — задача своя, решает сессия; отложенное — строкой таблицы 4.
  3. Заготовка `spec/repo-template/` (repositories 2.2 spec): `just new-repo-draft learner-model <вид> <карточка> v0.18.14` (вид — по `repo-template/kinds.json`, карточка — номер раздела architecture 3 spec) на клоне `spec` v0.18.14 во временной папке, как T56 `sim`. `.devcontainer/`, `.claude/`, CLAUDE.md, правила RUN.md и раскладка таблиц плана — по ней, где она изменилась, с отличиями, которые держат решения этого журнала.
  4. Решение журнала о переходе — что тег принёс и что взято —, поправки к записям, которые он меняет; строки таблиц 3 и 4, которые тег закрывает.
  5. Переход больше одного ревью — разбить по правилу 11: эта задача — первая часть, остальные — следующими задачами, и назвать их в отчёте.
- **Как проверить:** `just -f spec/justfile docs-check-repo` на новом теге — зелёный; проверки частей, которые задача меняет, — правило 6 этого RUN.md; `/code-review high`, один круг.
- **Мой шаг:** если задача меняет `.devcontainer/` — назвать это в отчёте: после мержа автор пересобирает контейнер, и оркестратор ждёт этого перед следующей задачей.
- **Разделы:** R156 spec, R76 spec; RULES.md spec; repositories 2.2 spec; карточка этого репозитория в architecture 3 spec.

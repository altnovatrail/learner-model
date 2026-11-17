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

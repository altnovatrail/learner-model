# TASKS.md — tasks outside an epic of `learner-model`

Open tasks outside the current epic, each as its body, in Russian (R154 spec, R155 spec).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T01\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T76. Тексты без истории (R160 spec): разделы 3.1–3.4 архитектуры

- **Цель:** разделы 3.1–3.4 `docs/architecture.md` описывают текущее устройство модели, без хроник (R160 spec).
- **Почему:** вторая часть документов из разделения T74 по правилу 11 RUN.md.
- **Сделать:** как в T75 — без хроник версий и задач, решение — номером, без пересказа; записи журнала — только те, что задача меняет, переписанные на месте целиком. Больше одного ревью — разделить по правилу 11 при планировании.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`.
- **Разделы:** R160 spec; [docs/architecture.md](../docs/architecture.md), разделы 3.1–3.4.
- **Зависит от:** T74.

#### T77. Тексты без истории (R160 spec): разделы 3.5–3.7 архитектуры

- **Цель:** разделы 3.5–3.7 `docs/architecture.md` — лог и пересчёт, сборка и те же биты, типы и кодирование — описывают текущее устройство, без хроник: версии `.github`, теги и релизы — файлом, где они заданы (R160 spec).
- **Почему:** третья часть документов из разделения T74 по правилу 11 RUN.md.
- **Сделать:** как в T75; цепочки вроде «from T26 v0.1.21, from T48 v0.1.24 …» — ссылкой на `.github/workflows/`; записи журнала — только те, что задача меняет, переписанные на месте целиком. Больше одного ревью — разделить по правилу 11 при планировании.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`.
- **Разделы:** R160 spec; [docs/architecture.md](../docs/architecture.md), разделы 3.5–3.7.
- **Зависит от:** T74.

#### T78. Комментарии кода по R160 spec

- **Цель:** комментарии Go, `.proto` вне описаний схемы, YAML, `justfile` и Dockerfile говорят только то, чего не скажет код: почему, ограничение, ловушка, контракт экспортируемого API; номер решения — якорем в конце, `(R25)` (R160 spec, «Comments»).
- **Почему:** четвёртая часть из разделения T74 по правилу 11 RUN.md; `docs-check` с `"comments": true` проверяет номера, но не содержание.
- **Сделать:**
  - без пересказа кода, истории, номеров задач и разделов документов (`docs/architecture.md, section 3.6` и т. п.); `R10 in docs/decisions.md` — `(R10)`;
  - описания полей `.proto` — для читателя выгрузки (`.claude/rules/proto.md`), их R160 spec не трогает, кроме номеров;
  - уведомления лицензий в `detmath/port.go` — как есть;
  - изменение `*.go` вне тестов, `proto/` и `buf.yaml` выпускает следующий патч (R16): одна задача — один выпуск. Больше одного ревью — разделить по правилу 11 при планировании.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; проверки частей `go`, `proto` и `ci` по `.claude/rules/`; `/code-review high`.
- **Разделы:** R160 spec, R32 spec; RULES.md spec — The bar; `.claude/rules/`.
- **Зависит от:** T74.

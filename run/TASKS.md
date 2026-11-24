# TASKS.md — tasks outside an epic of `learner-model`

Open tasks outside the current epic, each as its body, in Russian.

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T01\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations). When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T99. `docs/sources.md` по правилу записей, часть 3: Go, CI, GitHub и инструменты

- **Цель:** утверждения разделов 4 и 5 держит ссылка на версию или коммит — спецификация и исходники Go на теге версии, `google.golang.org/protobuf` на версии, файлы `altnovatrail/.github` на теге и коммите, — тест, проверка CI или `measurements`; записями остаются только условия, цены, лицензии и страницы без версии, которых прогон не проверит.
- **Почему:** третья часть T97 по правилу 11 этого RUN.md; строка 7.65 COORDINATION.md `spec`, по тексту T291 `spec`; R171 `spec`.
- **Сделать:**
  - 4.1–4.6 и 4.8 — ссылки на `doc/go_spec.html` и файлы `src` в `golang/go` на теге версии, которую читала запись, и на `pkg.go.dev` на версии; 4.7, Go Wiki, — есть ли у страницы версия, решить по правилу записей;
  - 5.2–5.5, 5.8, 5.9 — ссылки на файлы `.github` на коммите релиза; ссылку на запись 13.8 `.github` сверить с записью `.github` на его `main`;
  - 5.1 — запись держит только условие «free and unlimited» для публичных репозиториев; ярлыки раннеров — ссылка на страницу на коммите `github/docs` или задание CI, которое на них идёт; 5.6 — версии держат `.devcontainer/Dockerfile` и лок-файл gopls, лицензии — по правилу записей; 5.7 — состояние Codecov держит проверка `codecov/patch`, запись удалить;
  - около 40 ссылок в `docs/decisions.md`, `docs/architecture.md` и `run/DEBT.md`; смысл утверждений не менять; номера не переиспользовать; у темы, что опустела, — строка, где теперь её факты; архивы `docs/plan/` не трогать.
- **Как проверить:** `just -f spec/justfile docs-check-repo` — зелёный; `git grep -nE 'sources [45]\.[0-9]+' -- docs run ':!docs/plan'` — только оставшиеся записи; размер `docs/sources.md` до и после — в отчёт; `/code-review high`, один круг. Строка 7.65 `spec` закрыта, когда сделаны T97–T99.
- **Разделы:** `spec/RULES.md`, «The bar» и «Citations»; R171 `spec`; строка 7.65 COORDINATION.md `spec`; T97.

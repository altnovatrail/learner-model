---
paths:
  - "**/*.go"
  - "go.mod"
  - "go.sum"
  - "testdata/**"
  - "detmath/**"
  - "internal/**"
  - ".github/coderules.json"
---

# The Go code of learner-model

The library is one Go module, `github.com/altnovatrail/learner-model`, at the root, with the package `detmath` — the port of `math.Exp` and `math.Log` that gives the same bits on every build — and the codec of protobuf's wire form in `internal/wire` (R09, R10). It requires nothing: its code takes Go's standard library alone (R09). `testdata/` holds the golden logs and their states, `testdata/fold/`, the golden reads, `testdata/read/`, and the golden bytes of the codec, `testdata/wire/` — the last are also the schema's, `.claude/rules/proto.md`. A change of code comes with its tests, and Codecov waits for `codecov/patch`, at least 90% of the changed lines covered (R48 spec; R13).

- `just ci-lint`, `just ci-test`, `just ci-vuln` and `just ci-licenses` — `gofmt` and `go vet`, and gopls's lock file against `GOPLS_VERSION`; the tests with `-race` and their coverage in `coverage/go.out`; govulncheck; go-licenses with `ALLOWED_LICENSES`, which CI takes from `.github`'s `go.yml`. A task that touches Go code ends with them green.
- **The same bits** on every build are held by the eight rules of code of architecture 3.6 — `float64()` round each product and quotient, from `math` only functions with an exact result and exp and log from `detmath`, no NaN or ∞, checked conversions, a canonical order, no global mutable state, imports from its list — and checked by the jobs of `.github`'s shared `go.yml` with `same-bits`: "Go / targets" runs the tests on the builds where the bits are most at risk, wasip1 too, "Go / fma-diff" compares the assembler of the fused targets and "Go / rules" the rules of code, by the lists of `.github/coderules.json` (R09, R18, R34). A task that touches the arithmetic or the imports ends with them green on its pull request.
- **Golden files** are written anew only where a change of what the model computes is meant: `LEARNERMODEL_UPDATE=1 go test .`, and the JSON of the codec's bytes then with `just proto-json`. A state with a new field is a new kind of state, and `api` recomputes its cache (R08; R80 spec).
- **The gate.** Every change after the first version passes the gate of `sim`, the job "Gate" (R12, R20): a change that keeps what the model computes keeps all the gate's numbers to the bit; one that changes it is judged by `sim` before it is released (R103 spec), `.claude/rules/ci.md`.
- A session gets gopls's diagnostics after each edit of a Go file, through the plugin `gopls-lsp` of `.claude/settings.json` (R35); they do not stand in for the checks above.

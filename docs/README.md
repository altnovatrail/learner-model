# The documents of learner-model

What the repository's README.md leaves out: the index of `docs/`, the detail of the releases and how to work on the repository.

- [architecture.md](architecture.md) — the component: what it does and does not do, its data and links, what was deliberately deferred, and how it is built.
- [decisions.md](decisions.md) — the decision log: what was decided, why, and what was rejected.
- [sources.md](sources.md) — the external facts the documents rest on, each with its source, its date and its quote.
- [measurements.md](measurements.md) — the results of measurements.
- `plan/` — the archives of the done epics and their tasks.

## Releases

Releases come out by themselves: each commit of `main` whose CI has passed and which changes the code, the schema or what builds the release's file gets the next patch of the line `0.1`, with its tag, by .github's shared release workflow ([R16](decisions.md)). Nobody tags by hand, and a release is immutable.

- **Go code** takes the library as a Go module, by the tag: `go get github.com/altnovatrail/learner-model@<tag>`. `go.sum` and the Go checksum database hold the version from then on.
- **The release's one file** is the schema of the evidence and the state, for what reads them without Go: `learner-model-proto-<tag>.tar.gz`, with `LICENSE`, `buf.yaml` and the `.proto` files of `proto/`, signed without a key; `gh release download <tag> -R altnovatrail/learner-model` takes it with its bundle. Its `.proto` files are the ones `api` puts into a learner's export, under `api/proto/novatrail/learnermodel/v1/`, beside the JSON of the evidence log and the state they describe, so their comments say what each field is without the code (R135 spec). Verify it where it enters NovaTrail with spec's recipe, `just -f spec/justfile release-verify learner-model <tag> learner-model-proto-<tag>.tar.gz`, with its bundle `<file>.sigstore.json` beside it. `just release-files <tag>`, in the toolchain stage of the tag's `.devcontainer/Dockerfile`, builds the same bytes from the tag's tree ([architecture.md](architecture.md), section 3.7).

## Contributing

The working rules, for people and coding agents alike, are in [CLAUDE.md](../CLAUDE.md) and [AGENTS.md](../AGENTS.md), and the work plan is [RUN.md](../RUN.md), with its other tables in `run/`. All work happens in the devcontainer of `.devcontainer/`. The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec), checked out here as a git submodule: clone this repository with `git clone --recurse-submodules`.

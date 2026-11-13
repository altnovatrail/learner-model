---
paths:
  - "proto/**"
  - "buf.yaml"
  - "testdata/wire/**"
---

# The schema of the evidence and the state

`proto/novatrail/learnermodel/v1/` is the schema by which the library encodes an evidence record and a state in protobuf's binary form, and `buf.yaml` its module (R10; architecture 3.7). The codec of `internal/wire` writes and reads the same bytes by hand, and the golden bytes `testdata/wire/*.binpb`, with their JSON beside them, hold the two to each other.

- `just ci-proto` — `buf lint` by `buf.yaml`; `buf breaking`, category FILE, against the latest tag `v*`; and `buf convert` of each golden run of bytes to JSON by the schema, which must be the JSON beside it, and back to the same bytes. The job "Proto" runs it, and `just ci-release` after it. A task that touches the schema or the codec ends with it green.
- `just proto-json` writes the JSON of the golden bytes anew, after `LEARNERMODEL_UPDATE=1 go test .`.
- **A record only grows.** A new observation is a new field with the next free number; a field is never renumbered, retyped or reused, and recompute reads every record of the log as it was written (R08; R80 spec). A new field of the state is a new kind of state.
- **The comments are for the reader of the export.** `api` puts the `.proto` files of the version that computed into a learner's export, beside the JSON of the log and the state: each field says what it means, its units, what "not recorded" is, and which fields the model does not read (R135 spec; Д181 spec).
- **The release's one file** is the schema's archive, `learner-model-proto-<version>.tar.gz`, built by `just release-files` and held to the tree by `just ci-release` (R16): a change of `proto/` or `buf.yaml` releases the next patch, `.claude/rules/ci.md`.

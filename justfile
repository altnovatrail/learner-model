# learner-model's own recipes. The ci- ones are what .github's shared go.yml runs, one job each, in
# the toolchain stage of .devcontainer/Dockerfile (R16 .github; R13 in docs/decisions.md): their
# names are the same in every Go repository of NovaTrail, and what they run is this repository's
# own.

# The tools only the checks need, at exact versions. Each is a Go program run straight from its
# module, as module-kit's justfile runs them, so nothing is installed for it and go.mod gains no
# require: the Go module proxy serves that version and the checksum database vouches for it.
# Dependabot does not see them, so they are raised by hand.
GOVULNCHECK := "golang.org/x/vuln/cmd/govulncheck@v1.9.0"
GO_LICENSES := "github.com/google/go-licenses/v2@v2.0.1"

# List the recipes.
default:
    @just --list

# Fails on code gofmt would rewrite, on what go vet reports, and on a gopls lock file that names
# another version than GOPLS_VERSION in the Dockerfile.
# Check formatting, go vet and the version of gopls.
ci-lint:
    #!/usr/bin/env bash
    set -euo pipefail
    # The Go files git tracks or would, outside every testdata/, the root's too, which glob magic
    # takes and a bare **/ would not: spec/ is a submodule, another repository's, and git lists
    # none of its files. A file deleted from the working tree but still in the index is skipped,
    # so that gofmt does not fail on it with no word of formatting.
    unformatted="$(git ls-files -z --cached --others --exclude-standard -- '*.go' ':(exclude,glob)**/testdata/**' \
        | while IFS= read -r -d '' file; do if [ -e "$file" ]; then printf '%s\0' "$file"; fi; done \
        | xargs -0 -r gofmt -l)"
    if [ -n "$unformatted" ]; then
        echo "ci-lint: gofmt would rewrite these files; run gofmt -w on them:" >&2
        echo "$unformatted" >&2
        exit 1
    fi
    go vet ./...
    # gopls comes into the devcontainer image from its lock file, .devcontainer/gopls/, which the go
    # kind of spec's template gives, and the image's build stops when the gopls it built does not
    # name GOPLS_VERSION (R19). CI never builds that stage, so the two are held to each other here:
    # a pull request that raises one alone, or a move to a new tag of spec that loses the step or
    # the lock file, fails here.
    lock="$(sed -n 's/^[[:space:]]*\(require[[:space:]]\{1,\}\)\{0,1\}golang\.org\/x\/tools\/gopls[[:space:]]\{1,\}v\([^[:space:]]\{1,\}\).*$/\2/p' \
        .devcontainer/gopls/go.mod)"
    arg="$(sed -n 's/^ARG GOPLS_VERSION=\(.*\)$/\1/p' .devcontainer/Dockerfile)"
    if [ -z "$lock" ] || [ "$arg" != "$lock" ]; then
        echo "ci-lint: .devcontainer/gopls/go.mod names gopls ${lock:-at no version} and" \
            ".devcontainer/Dockerfile has GOPLS_VERSION ${arg:-unset}; raise them together (R19)" >&2
        exit 1
    fi

# The coverage goes into coverage/go.out, one Go profile, which the job coverage of .github's shared
# go.yml sends to Codecov (R48 spec; R13). Go counts statements and leaves the tests themselves out;
# a package that no test runs is in it too, its statements at zero. Every test counts for every
# package of the module, -coverpkg=./...: a test of golden vectors runs the code of other packages,
# and its lines count where they are. The checks of the same bits on every build — the tests on the
# targets, the comparison of the assembler and the rules of code — are go.yml's, by its input
# same-bits, with the lists of .github/coderules.json (R09; R18), and its job arm64 runs this recipe
# too, on arm64 hardware: what it runs must run on every host.
# Run the tests, with the race detector, and write their coverage.
ci-test:
    mkdir -p coverage
    go test -race -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage/go.out ./...

# The schema of the evidence and the state, by buf of the container (R10; docs/architecture.md,
# section 3.7): buf lint by buf.yaml; buf breaking, category FILE, against the latest tag v* of
# this repository, and before the first tag nothing to compare with; and buf convert of each
# golden run of bytes of the codec, testdata/wire/*.binpb, to JSON by the schema, which must be
# the JSON beside it, and of that JSON back to the same bytes. So the codec and the schema say the
# same. The JSON is written anew with `just proto-json`, after LEARNERMODEL_UPDATE=1 go test.
# Check the schema with buf, and the golden bytes of the codec against it.
ci-proto:
    #!/usr/bin/env bash
    set -euo pipefail
    buf lint
    latest="$(git -c versionsort.suffix=- tag --list 'v*' --sort=-v:refname | head -n 1)"
    if [ -n "$latest" ]; then
        buf breaking --against ".git#tag=$latest"
    else
        echo "ci-proto: no tag v* yet: buf breaking has nothing to compare with" >&2
    fi
    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp"' EXIT
    for pair in evidence:Evidence state:State; do
        name="${pair%%:*}" type="novatrail.learnermodel.v1.${pair##*:}"
        buf convert proto --type "$type" --from "testdata/wire/$name.binpb#format=binpb" \
            --to "$tmp/$name.json#format=json,use_proto_names=true"
        if ! cmp -s "$tmp/$name.json" "testdata/wire/$name.json"; then
            echo "ci-proto: buf reads testdata/wire/$name.binpb otherwise than testdata/wire/$name.json says:" >&2
            diff <(cat "$tmp/$name.json"; echo) <(cat "testdata/wire/$name.json"; echo) >&2 || true
            exit 1
        fi
        buf convert proto --type "$type" --from "testdata/wire/$name.json#format=json" --to "$tmp/$name.binpb#format=binpb"
        if ! cmp -s "$tmp/$name.binpb" "testdata/wire/$name.binpb"; then
            echo "ci-proto: buf writes testdata/wire/$name.json as other bytes than the codec's" >&2
            exit 1
        fi
        echo "ci-proto: $name: the codec and the schema agree" >&2
    done

# Write the JSON of the golden bytes of the codec anew, by the schema.
proto-json:
    #!/usr/bin/env bash
    set -euo pipefail
    for pair in evidence:Evidence state:State; do
        name="${pair%%:*}"
        buf convert proto --type "novatrail.learnermodel.v1.${pair##*:}" --from "testdata/wire/$name.binpb#format=binpb" \
            --to "testdata/wire/$name.json#format=json,use_proto_names=true"
    done

# Fail on a known vulnerability in code this module reaches.
ci-vuln:
    go run {{ GOVULNCHECK }} ./...

# Fails if a package this module is built from carries a license outside ALLOWED_LICENSES, SPDX
# identifiers separated by commas, which .github's go.yml sets.
# Check the licenses of what the module is built from against ALLOWED_LICENSES.
ci-licenses:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "${ALLOWED_LICENSES:-}" ]; then
        echo "ci-licenses: ALLOWED_LICENSES is not set. CI takes it from .github's go.yml; to run this here, set it to the list there." >&2
        exit 1
    fi
    go run {{ GO_LICENSES }} check ./... --allowed_licenses="$ALLOWED_LICENSES"

# What .github's shared release.yml runs to build a release, in the toolchain stage of the commit
# it releases, from that commit's tree with no git in it, and then signs every file (R16 in
# docs/decisions.md; R18 .github). Go code takes the library as a Go module, by the tag, through
# the module proxy and the checksum database; the release's one file is the schema of the evidence
# and the state, for what reads them without Go: the export of a learner's data, say (section 3.7
# of docs/architecture.md). dist/ is made afresh and holds learner-model-proto-<version>.tar.gz:
# LICENSE, buf.yaml and every .proto file of proto/ in a folder of that name, in a tar every build
# writes alike — names sorted, no times or owners, modes 0644 — and gzip without a name or a time.
# Build a release's files into dist/, for a version such as v0.1.0.
release-files version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{ quote(version) }}
    if ! [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
        echo "release-files: the version is v<major>.<minor>.<patch>, as v0.1.0, not ${version}" >&2
        exit 1
    fi
    name="learner-model-proto-${version}"
    rm -rf dist
    mkdir dist
    (printf '%s\n' LICENSE buf.yaml; find proto -type f -name '*.proto') | LC_ALL=C sort |
        tar --create --verbatim-files-from --files-from=- --format=posix \
            --pax-option=exthdr.name=%d/PaxHeaders/%f,delete=atime,delete=ctime \
            --transform="s,^,${name}/," --mtime=@0 --owner=0 --group=0 --numeric-owner --mode=0644 |
        gzip -9n > "dist/${name}.tar.gz"

# Builds the release's files twice, a second apart, and fails unless the two are the same bytes;
# unless the archive holds LICENSE, buf.yaml and every .proto file git tracks in proto/, as they are
# in the tree, and nothing else; and unless buf builds the schema from the unpacked archive alone.
# So the file a release signs is the schema of its commit, and anyone can build it again and
# compare. It refuses to run beside a dist/ of one's own, which release-files would replace.
# Check that the release's file is the schema, and the same bytes on every build.
ci-release:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -e dist ]; then
        echo "ci-release: dist/ is there already, and release-files would replace it: move it away first" >&2
        exit 1
    fi
    version=v0.0.0-ci
    name="learner-model-proto-${version}"
    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp" dist' EXIT
    just release-files "$version"
    mv dist "$tmp/first"
    sleep 1
    just release-files "$version"
    if [ "$(ls dist)" != "${name}.tar.gz" ] || ! cmp -s "$tmp/first/${name}.tar.gz" "dist/${name}.tar.gz"; then
        echo "ci-release: two builds of the release's files differ, or dist/ holds another file:" >&2
        ls -l "$tmp/first" dist >&2
        exit 1
    fi
    mkdir "$tmp/unpacked"
    tar -xzf "dist/${name}.tar.gz" -C "$tmp/unpacked"
    got="$(cd "$tmp/unpacked/$name" && find . -type f | sed 's|^\./||' | LC_ALL=C sort)"
    want="$( (printf '%s\n' LICENSE buf.yaml; git ls-files -- 'proto/*.proto') | LC_ALL=C sort)"
    if [ "$got" != "$want" ]; then
        echo "ci-release: the archive holds" >&2
        echo "$got" >&2
        echo "and git tracks" >&2
        echo "$want" >&2
        exit 1
    fi
    while IFS= read -r file; do
        if ! cmp -s "$file" "$tmp/unpacked/$name/$file"; then
            echo "ci-release: ${file} in the archive is not the tree's" >&2
            exit 1
        fi
    done <<< "$want"
    (cd "$tmp/unpacked/$name" && buf build -o /dev/null)
    echo "ci-release: ${name}.tar.gz is the schema, the same bytes on two builds" >&2

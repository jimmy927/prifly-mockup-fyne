#!/usr/bin/env bash
# needs: go
# The smallest proof that the recipe still works: the Quillnote fixture compiles to a Go WebAssembly module
# with Fyne's browser driver (GOOS=js GOARCH=wasm), the same compile `fyne package -os wasm` runs. Needs the
# module cache to hold Fyne already (or network once); packing and posting are the recipe's later steps.
set -euo pipefail
here="$(cd "$(dirname "$0")" && pwd)"
export PATH="$HOME/.local/go/bin:$PATH"
out="$(mktemp -d)"
trap 'rm -rf "$out"' EXIT
cd "$here/app"
GOOS=js GOARCH=wasm go build -ldflags "-s -w" -tags release -o "$out/Quillnote.wasm" .
test -s "$out/Quillnote.wasm"
test -f "$(go env GOROOT)/misc/wasm/wasm_exec.js" || test -f "$(go env GOROOT)/lib/wasm/wasm_exec.js"
echo "built $(wc -c <"$out/Quillnote.wasm") bytes"

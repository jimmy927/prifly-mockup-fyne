# prifly-mockup-fyne

A prifly mockup-stack recipe for **Fyne**, the Go native GUI toolkit. It tells a
prifly session how to render a mockup with the real Fyne widgets: Fyne's own WebAssembly target, packed into one
self-contained HTML page, instead of an HTML imitation of a Fyne window.

## What is in it

- `mockup-stacks/fyne.md`: the recipe (frontmatter that prifly's loader reads, then Steps, Gotchas, States).
- `mockup-stacks/fixtures/fyne/`: a small proof project, "Quillnote" (tabs, a list, a form, a confirm dialog, settings), sources
  only. `build.sh` compiles it to wasm; prifly's checker runs it when `go` is on the PATH.
- `scripts/check.ts`: a Bun-only validator of the recipe frontmatter, run by CI.

## Use it in prifly

This repo is a prifly extension. In prifly, open **Extensions**, install from the git URL
`https://github.com/jimmy927/prifly-mockup-fyne.git`, then enable it. Sessions then get the recipe from
`mcp__prifly__mockup_recipe` (layer order: project `.prifly/mockup-stacks/` > user folder > enabled extensions >
built-in). A project can also copy `mockup-stacks/fyne.md` into its own `.prifly/mockup-stacks/`.

## The recipe in one paragraph

Build the app with `fyne package -os wasm --release` (Go 1.23, Fyne 2.8.1: about 13 s, a 35 MB module), pack the `.wasm` and
`wasm_exec.js` with the kit's `pack-wasm.py generic` into one 14.7 MB page, and start it from a small `boot.js` that works around
three things: Fyne 2.8 panics when IndexedDB is unavailable in a file page, the kit's ready signal fires before Fyne's first frame,
and a mockup-control strip above the canvas must not shift it. States are `#hash` links in that strip. The Steps and Gotchas in
`mockup-stacks/fyne.md` have the commands, timings and sizes.

## Screenshots

The recipe's card in prifly:

![The Fyne recipe card in prifly](docs/prifly-card.png)

A mockup rendered by it:

![A Fyne mockup open in prifly](docs/prifly-open.png)

## Checking

```
bun scripts/check.ts
```

validates the frontmatter against the loader's schema. From a prifly checkout,
`bun run mockup-stacks:check --dir <path-to-this-repo>/mockup-stacks` runs the full checker on this folder and also builds the fixture.

## License

MIT, see `LICENSE`.

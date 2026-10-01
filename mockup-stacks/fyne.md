---
id: fyne
name: "Fyne (Go native GUI toolkit, WebAssembly)"
aliases: [fyne, fyne.io, go gui, golang gui, fyne wasm]
detect:
  deps: ["fyne.io/fyne/v2", "fyne.io/fyne"]
  files: ["FyneApp.toml"]
route: wasm
kind: desktop
where: local
toolchain: "Go 1.23.4 (~/.local/go) + fyne CLI v1.7.3 (go install fyne.io/tools/cmd/fyne@latest, GOBIN ~/.cache/mockup-toolchains/fyne/bin), Fyne library v2.8.1, kit pack-wasm.py generic"
versions: "Go 1.23.4, fyne CLI 1.7.3, fyne.io/fyne/v2 2.8.1"
renderedBy: ["~/.cache/mockup-toolchains/fyne", "go"]
verified: 2026-10-01
verifiedBy: "verified by prifly sessions on 2026-10-01"
fixture: fixtures/fyne
---

## Steps
Fyne has its own official browser port: its GLFW layer compiles to WebAssembly (`GOOS=js GOARCH=wasm`) and draws the
real widgets with WebGL into a `<canvas>`. That is the route; no HTML imitation, no screenshots of a desktop window.
1. Toolchain once: `GOBIN=~/.cache/mockup-toolchains/fyne/bin go install fyne.io/tools/cmd/fyne@latest` (12 s). Go itself is
   `~/.local/go/bin/go`. A native `go build`/`go vet` of a Fyne app needs cgo, GL and Wayland headers and failed here
   (`WL_MARSHAL_FLAG_DESTROY undeclared`); do not try it, check with `GOOS=js GOARCH=wasm go vet ./...`.
2. Put the app in `<scratchpad>/app`: `go mod init`, `go get fyne.io/fyne/v2@latest`, `go mod tidy`; `Icon.png` next to `main.go`
   (`fyne package` refuses to run without it). Real widgets, invented fixed data. Read the state from `window.__state`
   in a `//go:build js` file (`syscall/js`, see the fixture's `state_js.go`) and keep a `!js` twin so the package still
   type-checks natively; and call a one-line `announceFirstFrame()` (also in `state_js.go`) before `w.ShowAndRun()`.
3. Build under the lock: `K=$(mockup-kit path); PATH=$HOME/.local/go/bin:$HOME/.cache/mockup-toolchains/fyne/bin:$PATH;
   cd app; nice -n 10 flock ~/.cache/mockup-build.lock fyne package -os wasm --release --app-id com.example.app --name Quillnote`
   (12.7 s wall, 61 s CPU, with a warm module cache; 2-3 s when only your code changed). Output `app/wasm/`: `Quillnote.wasm`
   (35 MB, 11 MB gzipped), `wasm_exec.js`, and an `index.html` + gifs that you do NOT use. `--release` strips symbols; without it the
   module is 36.6 MB.
4. Copy the fixture's `boot.js` and adapt the wasm name (it is the file name under `app/wasm/`). Pack with the generic preset:
   `python3 $K/scripts/pack-wasm.py generic --wasm wasm/Quillnote.wasm --script wasm/wasm_exec.js --boot boot.js --title Quillnote --height 800 --out quillnote.html`
   (3.6 s, 14.7 MB, within the 16 MB limit) then `mockup-kit check quillnote.html`. `boot.js` is the part that makes it work, see Gotchas.
5. One more pack per state for the picture, `--default-state editor|settings|dialog --quiet --out chk-<state>.html`, because the tool
   refuses `file#state`. `mcp__prifly__mockup` with `check: true` on each (1658 wide; the load takes under 8 s here), then post `quillnote.html` once.
   `renderedBy`: `["~/.cache/mockup-toolchains/fyne", "go"]` written as absolute paths.

## Gotchas
- **Blank page, no error from the tool.** Go panics are only on the console. While building, add a debug script that mirrors
  `console.log/error` into a fixed `<pre>` (a 10-line `--script debug.js`) and read the picture; drop it before posting.
- **Fyne 2.8 panics at start in a file page**: `NewIndexDBRepository` -> `(*OpenDBRequest).Await` nil pointer, because
  `indexedDB.open` throws at an opaque origin and the Go library does not check it. `boot.js` replaces `window.indexedDB` with an
  object whose `open()` returns an `Object.create(IDBOpenDBRequest.prototype)` that fails asynchronously (go-indexeddb checks
  `instanceof IDBRequest`, and the members must be own properties). Fyne then logs "failed to create repository", "mkdir ... not
  implemented on js" and carries on; those lines are harmless.
- **The kit's `__mockupReady` fires too early**: it waits for a canvas plus 1.2 s, and Fyne's canvas exists before anything is drawn.
  `boot.js` replaces `window.__mockupReady` with its own promise that the Go side resolves (`__goReady`) 1.5 s after the window is
  shown, 7.5 s at the latest. The kit's probe workaround (forced `preserveDrawingBuffer`) is what lets the tool copy the WebGL canvas.
- **A mockup control must not shift Fyne's canvas.** Fyne sizes its canvas buffer once from `window.innerHeight` and sets the canvas
  CSS to `100dvw x 100dvh`. A bar above it needs: `body{padding-top:34px; box-sizing:border-box}`, the canvas CSS height fixed to
  `800-34` px with `!important`, and `innerHeight` reporting `766` only until the canvas exists. Leave the override in place and the tool,
  which sizes its picture from `innerHeight`, crops the bottom 34 px; use `position:absolute` for the canvas and the drawing shifts up by 34 px.
- Window size: the browser window is the Fyne window, `w.Resize(...)` is ignored; design for the viewport (1658 x 800).
- Container order trap: `container.NewBorder(top, bottom, left, right, center)` keeps `Objects` as centre, top, bottom, left, right when you
  index into list-item templates (the wrong index panics with an interface conversion and blanks the page).
- Dialogs (`dialog.ShowConfirm`, `ShowForm`) are drawn inside the canvas as overlays, so they are part of the picture. The `dialog` state
  shows one by calling it from `apply()` before `ShowAndRun`; it rendered correctly.
- `theme.LightTheme()` pins the look; otherwise Fyne follows `prefers-color-scheme` of the tool's browser.
- No emoji or web fonts: Fyne ships its own Noto-based font and text is drawn by it, nothing is fetched. Icons come from `theme.*Icon()`.
- Hash handling is a reload (the kit reloads on `hashchange`), so each state starts the 35 MB module again: expect 3-6 s per state change.

## States
`#notes` (default; first note selected), `#editor` (third note selected), `#settings`, `#about`, `#dialog` (delete-confirmation overlay).
`boot.js` draws a 34 px dark strip on top labelled "Mockup control (not part of the app)" with one pill per state, linking to the hashes;
the Fyne canvas sits below it at full width. The Go side maps the hash with `apply(state)`.

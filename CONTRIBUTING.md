# Contributing

Fixes and new recipes are welcome. A recipe records something that was built and posted, not something that should work.

- One file per stack in `mockup-stacks/<id>.md`; the file name is the `id`.
- Frontmatter keys: `id`, `name`, `aliases`, `detect` (`deps`, `files`, `globs`), `route`, `kind`, `where`, `toolchain`, `versions`,
  `renderedBy`, `verified` (YYYY-MM-DD), `verifiedBy`, `fixture`. Then `## Steps`, `## Gotchas`, `## States`.
- Name only versions that were actually run, and give timings and sizes for the slow steps.
- Keep it private-free: no home paths (write `~` or `<scratchpad>`), machine names, addresses, session ids or company names.
- A fixture is sources only (no binaries), builds offline in under two minutes, and has a `build.sh` whose second line is
  `# needs: <programs>`.
- Run `bun scripts/check.ts` before opening a pull request; CI runs the same.
- Change the recipe only after re-running the steps, then update `verified` and `versions`.

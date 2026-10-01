/**
 * `bun scripts/check.ts`: validates every recipe in mockup-stacks/ against the schema prifly's loader uses
 * (apps/desktop-host/src/mockup-stacks/recipe.ts), without needing prifly or any dependency. Exit 1 on a problem.
 */
import { existsSync, readdirSync, readFileSync, statSync } from "node:fs";
import { basename, isAbsolute, join, normalize } from "node:path";

const root = join(import.meta.dir, "..");
const dir = join(root, "mockup-stacks");
const NOT_RECIPES = new Set(["readme.md", "coverage.md"]);
const FRONT = /^﻿?---[ \t]*\r?\n([\s\S]*?)\r?\n---[ \t]*(?:\r?\n|$)/;
const KNOWN = new Set([
  "id", "name", "aliases", "detect", "route", "verified", "versions", "fixture",
  "priority", "kind", "where", "toolchain", "renderedBy", "verifiedBy",
]);

const scalar = (v: unknown) => typeof v === "string" || typeof v === "number";

function names(value: unknown, label: string, problems: string[]): void {
  if (value === undefined) return;
  if (!Array.isArray(value) || !value.every((v) => scalar(v) && String(v).length > 0)) {
    problems.push(`${label}: a list of non-empty names`);
  }
}

function validate(file: string, text: string): string[] {
  const problems: string[] = [];
  const match = FRONT.exec(text);
  if (match === null) return ["no frontmatter block (--- ... --- at the top)"];
  let raw: unknown;
  try {
    raw = Bun.YAML.parse(match[1] ?? "");
  } catch (caught) {
    return [`frontmatter is not valid YAML: ${caught instanceof Error ? caught.message : String(caught)}`];
  }
  if (typeof raw !== "object" || raw === null || Array.isArray(raw)) {
    return ["frontmatter is not a mapping of keys to values"];
  }
  const fm = raw as Record<string, unknown>;
  for (const key of Object.keys(fm)) if (!KNOWN.has(key)) problems.push(`unknown key "${key}" (the loader drops it)`);

  if (typeof fm.id !== "string" || !/^[a-z0-9][a-z0-9-]*$/.test(fm.id)) problems.push("id: lower-case words joined by hyphens");
  else if (fm.id !== basename(file, ".md")) problems.push(`id "${fm.id}" does not match the file name`);
  if (!scalar(fm.name) || String(fm.name).length === 0) problems.push("name: required");
  if (!scalar(fm.route) || String(fm.route).length === 0) problems.push("route: required");
  if (!(scalar(fm.verified) || fm.verified instanceof Date) || !/^\d{4}-\d{2}-\d{2}/.test(fm.verified instanceof Date ? fm.verified.toISOString() : String(fm.verified))) {
    problems.push("verified: a date, YYYY-MM-DD");
  }
  names(fm.aliases, "aliases", problems);
  names(fm.renderedBy, "renderedBy", problems);
  if (fm.detect !== undefined) {
    if (typeof fm.detect !== "object" || fm.detect === null || Array.isArray(fm.detect)) problems.push("detect: a mapping");
    else for (const key of ["deps", "files", "globs"]) names((fm.detect as Record<string, unknown>)[key], `detect.${key}`, problems);
  }
  for (const key of ["versions", "kind", "where", "toolchain", "verifiedBy"]) {
    if (fm[key] !== undefined && !scalar(fm[key])) problems.push(`${key}: text`);
  }
  if (fm.priority !== undefined && !Number.isInteger(fm.priority)) problems.push("priority: an integer");
  if (fm.fixture !== undefined) {
    const path = String(fm.fixture);
    if (path === "" || isAbsolute(path) || normalize(path).startsWith("..")) problems.push("fixture: a path inside the recipe's folder");
    else if (!existsSync(join(dir, path))) problems.push(`fixture ${path} does not exist`);
    else {
      const target = join(dir, path);
      const script = statSync(target).isDirectory() ? join(target, "build.sh") : target;
      if (!existsSync(script)) problems.push(`no build.sh in ${path}`);
    }
  }
  const body = text.slice(match[0].length);
  for (const heading of ["## Steps", "## Gotchas", "## States"]) {
    if (!body.includes(`\n${heading}\n`) && !body.startsWith(`${heading}\n`)) problems.push(`missing section "${heading}"`);
  }
  if (/\/home\/[a-z]/i.test(text) || /\/Users\/[a-z]/i.test(text)) problems.push("contains a home-directory path (use ~ or <scratchpad>)");
  return problems;
}

let failed = false;
const files = readdirSync(dir).filter((f) => f.endsWith(".md") && !NOT_RECIPES.has(f.toLowerCase())).sort();
if (files.length === 0) {
  console.log("FAIL  no recipes in mockup-stacks/");
  failed = true;
}
for (const file of files) {
  const problems = validate(file, readFileSync(join(dir, file), "utf8"));
  if (problems.length === 0) console.log(`ok    ${file}`);
  else {
    failed = true;
    for (const p of problems) console.log(`FAIL  ${file}: ${p}`);
  }
}
process.exit(failed ? 1 : 0);

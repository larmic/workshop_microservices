// Guard-Gate: syntax check for the inline <script> blocks of the static pages.
// The pages have no build step, a typo in a script only shows up in the
// browser. Each inline block is checked with `node --check`; it is padded
// with empty lines so reported line numbers match the HTML file.
//
// Usage: node scripts/check-inline-js.mjs [file.html ...]   (from the repo root)
import { spawnSync } from "node:child_process";
import { mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

const defaults = [
  "services/dashboard/static/index.html",
  "services/slides/index.html",
  "pages/landing/index.html",
  "pages/vorbereitung/index.html",
  "pages/feedback/index.html",
  "pages/material/index.html",
];
const files = process.argv.length > 2 ? process.argv.slice(2) : defaults;
const scriptTag = /<script\b([^>]*)>([\s\S]*?)<\/script\s*>/gi;
const classic = new Set(["", "text/javascript", "application/javascript"]);

const dir = mkdtempSync(join(tmpdir(), "inline-js-"));
let blocks = 0;
let failed = 0;
try {
  for (const file of files) {
    const html = readFileSync(file, "utf8");
    for (const m of html.matchAll(scriptTag)) {
      const attrs = m[1];
      if (/\bsrc\s*=/i.test(attrs)) continue;
      const type = (attrs.match(/\btype\s*=\s*["']?([^"'\s>]+)/i)?.[1] ?? "").toLowerCase();
      if (type !== "module" && !classic.has(type)) continue; // JSON, templates, …
      const startLine = html.slice(0, m.index + m[0].indexOf(">") + 1).split("\n").length;
      const source = "\n".repeat(startLine - 1) + m[2];
      const tmp = join(dir, `block-${blocks}.${type === "module" ? "mjs" : "cjs"}`);
      writeFileSync(tmp, source);
      blocks++;
      const res = spawnSync(process.execPath, ["--check", tmp], { encoding: "utf8" });
      if (res.status !== 0) {
        failed++;
        console.error(`${file} (script ab Zeile ${startLine}):\n${res.stderr.replaceAll(tmp, file)}`);
      }
    }
  }
} finally {
  rmSync(dir, { recursive: true, force: true });
}

if (failed > 0) {
  console.error(`Inline-JavaScript: ${failed} von ${blocks} Blöcken mit Syntaxfehler`);
  process.exit(1);
}
console.log(`Inline-JavaScript: ${blocks} Blöcke in ${files.length} Dateien syntaktisch ok`);

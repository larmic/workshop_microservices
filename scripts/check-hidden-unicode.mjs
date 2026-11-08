// Guard-Gate: unsichtbare Unicode-Zeichen (Kategorie Cf: Zero-Width, Bidi-Steuerzeichen,
// bedingter Trennstrich, BOM …) in versionierten Textdateien. Solche Zeichen können
// Code und Text verschleiern ("Trojan Source"); Renovate warnt ebenfalls davor.
// Gewollte Zeichen gehören als Escape in den Quelltext, z. B. "\u200B".
import { execFileSync } from "node:child_process";
import { readFileSync } from "node:fs";

const files = execFileSync("git", ["ls-files", "-z"], { encoding: "utf8" }).split("\0").filter(Boolean);
const decoder = new TextDecoder("utf-8", { fatal: true });
const hidden = /\p{Cf}/gu;
let findings = 0;

for (const file of files) {
  let text;
  try {
    text = decoder.decode(readFileSync(file));
  } catch {
    continue; // Binärdatei (Bilder, Schriften)
  }
  text.split("\n").forEach((line, i) => {
    for (const m of line.matchAll(hidden)) {
      const cp = m[0].codePointAt(0).toString(16).toUpperCase().padStart(4, "0");
      console.error(`${file}:${i + 1}:${m.index + 1}: unsichtbares Zeichen U+${cp}`);
      findings++;
    }
  });
}

if (findings > 0) {
  console.error(`Guard-Gate unicode: ${findings} unsichtbare Zeichen gefunden (als Escape schreiben oder entfernen)`);
  process.exit(1);
}
console.log(`Unsichtbare Unicode-Zeichen: keine in ${files.length} Dateien`);

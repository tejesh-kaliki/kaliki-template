#!/usr/bin/env node
// Auto-discovers every services/*.yaml and joins them into openapi.yaml.
// No need to maintain a service list as you add domains — just drop a new
// services/<domain>.yaml and rebuild.
//
// Usage: node scripts/build-spec.mjs [config/servers-<env>.yaml] [-o openapi.yaml]
import { readdirSync } from "node:fs";
import { execFileSync } from "node:child_process";
import { join } from "node:path";

const serverConfig = process.argv[2] || "config/servers-local.yaml";
const output = process.argv[4] || "openapi.yaml";
const servicesDir = "services";

// Shared component specs (if present) must be joined before the specs that
// $ref them. Everything else is alphabetical for deterministic output.
const priority = ["global.yaml", "shared.yaml"];
const rank = (f) => {
  const i = priority.indexOf(f);
  return i === -1 ? priority.length : i;
};

const files = readdirSync(servicesDir)
  .filter((f) => f.endsWith(".yaml") || f.endsWith(".yml"))
  .sort((a, b) => rank(a) - rank(b) || a.localeCompare(b))
  .map((f) => join(servicesDir, f));

if (files.length === 0) {
  console.error(`No service specs found in ${servicesDir}/`);
  process.exit(1);
}

console.log(`Joining ${files.length} spec(s) -> ${output}`);
execFileSync(
  "npx",
  ["redocly", "join", serverConfig, ...files, "-o", output],
  { stdio: "inherit" },
);

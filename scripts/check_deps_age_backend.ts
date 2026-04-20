#!/usr/bin/env node
/**
 * Verify that every Go dependency (direct + transitive) was released at least
 * MIN_DAYS ago. Rewrite of check_deps_age.py in TypeScript.
 *
 * Cache: .deps_age_cache stores { "module@version": "<ISO timestamp>" }.
 * Release dates are fetched once per module@version and reused across runs.
 * Age is re-evaluated against the current date on every run.
 *
 * On violation: identifies the direct dependency responsible via `go mod graph`
 * and suggests a downgrade command.
 *
 * Usage: node scripts/check_deps_age_backend.ts
 * Tests: node --test scripts/check_deps_age_backend.ts
 */

import { execSync } from 'node:child_process';
import { readFileSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const MIN_DAYS = 7;
const HEADERS = { 'User-Agent': 'beanmemo-dep-checker/1.0' };

const PASS = '\x1b[32m✓\x1b[0m';
const FAIL = '\x1b[31m✗\x1b[0m';
const WARN = '\x1b[33m!\x1b[0m';

// ─── Types ────────────────────────────────────────────────────────────────────

interface ModuleInfo { module: string; version: string; }
interface ViolationInfo { module: string; version: string; releasedAt: Date; }

// ─── Pure functions ───────────────────────────────────────────────────────────

/**
 * Escape a module path for use in proxy.golang.org URLs.
 * Per the Go module proxy protocol, each uppercase letter X is replaced with !x.
 */
export function proxyEscape(module: string): string {
  return module.replace(/[A-Z]/g, c => `!${c.toLowerCase()}`);
}

/**
 * Parse the output of `go list -m all` into module/version pairs.
 * The first line (the root module itself, no version field) is skipped.
 */
export function parseGoListOutput(output: string): ModuleInfo[] {
  return output
    .split('\n')
    .map(line => line.trim())
    .filter(Boolean)
    .map(line => line.split(/\s+/))
    .filter(parts => parts.length >= 2)
    .map(parts => ({ module: parts[0], version: parts[1] }));
}

/**
 * Parse `go mod graph` output and return all entries that directly require
 * module@version. The root module appears without @version in the graph.
 */
export function findRequirers(graphOutput: string, module: string, version: string): string[] {
  const target = `${module}@${version}`;
  return graphOutput
    .split('\n')
    .map(line => line.trim())
    .filter(Boolean)
    .map(line => line.split(/\s+/))
    .filter(parts => parts.length === 2 && parts[1] === target)
    .map(parts => parts[0]);
}

/**
 * Given a map of version -> ISO timestamp, return the latest version published
 * strictly before `before`, or null if none exists.
 */
export function findLatestVersionBefore(
  timeMap: Record<string, string>,
  before: Date,
): string | null {
  const candidates = Object.entries(timeMap)
    .filter(([v]) => v !== 'created' && v !== 'modified')
    .filter(([, t]) => new Date(t) < before)
    .sort(([, a], [, b]) => new Date(b).getTime() - new Date(a).getTime());
  return candidates.length > 0 ? candidates[0][0] : null;
}

// ─── Tests ───────────────────────────────────────────────────────────────────
// Run with: node --test scripts/check_deps_age_backend.ts

if (process.env.NODE_TEST_CONTEXT) {
  const { test } = await import('node:test');
  const assert = (await import('node:assert/strict')).default;

  // proxyEscape
  test('proxyEscape: escapes uppercase letters with !lowercase', () => {
    assert.equal(proxyEscape('github.com/Foo/Bar'), 'github.com/!foo/!bar');
  });
  test('proxyEscape: leaves all-lowercase paths unchanged', () => {
    assert.equal(proxyEscape('github.com/foo/bar'), 'github.com/foo/bar');
  });

  // parseGoListOutput
  test('parseGoListOutput: skips root module line (no version)', () => {
    const output = 'github.com/my/module\ngithub.com/dep/a v1.2.3\ngithub.com/dep/b v2.0.0';
    assert.deepEqual(parseGoListOutput(output), [
      { module: 'github.com/dep/a', version: 'v1.2.3' },
      { module: 'github.com/dep/b', version: 'v2.0.0' },
    ]);
  });
  test('parseGoListOutput: handles empty output', () => {
    assert.deepEqual(parseGoListOutput(''), []);
  });

  // findRequirers
  const graphOutput = [
    'github.com/my/module github.com/dep/a@v1.2.3',
    'github.com/dep/b@v2.0.0 github.com/dep/a@v1.2.3',
    'github.com/dep/c@v3.0.0 github.com/dep/b@v2.0.0',
  ].join('\n');

  test('findRequirers: returns root module and transitive requirers', () => {
    const result = findRequirers(graphOutput, 'github.com/dep/a', 'v1.2.3');
    assert.deepEqual(result, ['github.com/my/module', 'github.com/dep/b@v2.0.0']);
  });
  test('findRequirers: returns empty array when version not in graph', () => {
    assert.deepEqual(findRequirers(graphOutput, 'github.com/dep/a', 'v9.9.9'), []);
  });

  // findLatestVersionBefore
  const before = new Date('2026-03-30T00:00:00Z');
  const timeMap: Record<string, string> = {
    created:  '2020-01-01T00:00:00Z',
    modified: '2026-04-01T00:00:00Z',
    'v1.0.0': '2026-03-28T00:00:00Z', // before cutoff → candidate
    'v1.0.1': '2026-03-29T00:00:00Z', // before cutoff → candidate (newest)
    'v1.0.2': '2026-04-01T00:00:00Z', // after cutoff  → excluded
  };

  test('findLatestVersionBefore: returns the latest compliant version', () => {
    assert.equal(findLatestVersionBefore(timeMap, before), 'v1.0.1');
  });
  test('findLatestVersionBefore: returns null when all versions are after cutoff', () => {
    assert.equal(findLatestVersionBefore({ 'v1.0.0': '2026-04-01T00:00:00Z' }, before), null);
  });
  test('findLatestVersionBefore: ignores created/modified keys', () => {
    const onlyMeta = { created: '2020-01-01T00:00:00Z', modified: '2020-01-01T00:00:00Z' };
    assert.equal(findLatestVersionBefore(onlyMeta, before), null);
  });
}

// ─── Main ─────────────────────────────────────────────────────────────────────

if (!process.env.NODE_TEST_CONTEXT) {
  const root = resolve(dirname(fileURLToPath(import.meta.url)), '..', 'backend');
  const cacheFile = resolve(root, '.deps_age_cache');

  const now = new Date();
  const threshold = new Date(now.getTime() - MIN_DAYS * 24 * 60 * 60 * 1000);

  // Cache helpers
  function loadCache(): Record<string, string> {
    try { return JSON.parse(readFileSync(cacheFile, 'utf8')); } catch { return {}; }
  }
  function saveCache(cache: Record<string, string>): void {
    writeFileSync(cacheFile, JSON.stringify(cache, null, 2));
  }

  // HTTP helpers
  async function fetchJson<T>(url: string): Promise<T | null> {
    try {
      const res = await fetch(url, { headers: HEADERS });
      if (res.status === 404) return null;
      if (!res.ok) { console.error(`  ${WARN} HTTP ${res.status} for ${url}`); return null; }
      return res.json() as Promise<T>;
    } catch (e) {
      console.error(`  ${WARN} fetch error for ${url}: ${e}`);
      return null;
    }
  }
  async function fetchText(url: string): Promise<string | null> {
    try {
      const res = await fetch(url, { headers: HEADERS });
      if (res.status === 404 || !res.ok) return null;
      return res.text();
    } catch { return null; }
  }

  async function getModuleReleaseTime(
    module: string,
    version: string,
    cache: Record<string, string>,
  ): Promise<Date | null> {
    const key = `${module}@${version}`;
    if (cache[key]) { try { return new Date(cache[key]); } catch { /* fall through */ } }
    const data = await fetchJson<{ Time: string }>(
      `https://proxy.golang.org/${proxyEscape(module)}/@v/${version}.info`,
    );
    if (!data) return null;
    cache[key] = data.Time;
    return new Date(data.Time);
  }

  // Check all modules against the age threshold
  async function checkGoModules(): Promise<ViolationInfo[]> {
    console.log('Checking Go dependencies (direct + transitive via `go list -m all`) …');
    const cache = loadCache();

    let goListOutput: string;
    try {
      goListOutput = execSync('go list -m all', { cwd: root }).toString();
    } catch (e) {
      console.error(`  ${WARN} \`go list -m all\` failed: ${e}`);
      return [];
    }

    const results = await Promise.all(
      parseGoListOutput(goListOutput).map(async ({ module, version }) => {
        const releasedAt = await getModuleReleaseTime(module, version, cache);
        if (!releasedAt) {
          console.log(`  ${WARN} ${module}@${version}: skipped (info unavailable)`);
          return null;
        }
        const ageDays = Math.floor((now.getTime() - releasedAt.getTime()) / 86400000);
        if (releasedAt > threshold) {
          console.log(`  ${FAIL} [Go] ${module}@${version}: released ${ageDays} day(s) ago (minimum ${MIN_DAYS})`);
          return { module, version, releasedAt };
        }
        console.log(`  ${PASS} ${module}@${version}: ${ageDays}d old`);
        return null;
      }),
    );

    saveCache(cache);
    return results.filter((r): r is ViolationInfo => r !== null);
  }

  // Find the latest version of dep published before `before`
  async function suggestDowngrade(dep: string, before: Date): Promise<string | null> {
    const escaped = proxyEscape(dep);
    const text = await fetchText(`https://proxy.golang.org/${escaped}/@v/list`);
    if (!text) return null;
    const versions = text.split('\n').map(v => v.trim()).filter(Boolean);

    const entries = await Promise.all(
      versions.map(async ver => {
        const data = await fetchJson<{ Time: string }>(
          `https://proxy.golang.org/${escaped}/@v/${ver}.info`,
        );
        if (!data) return null;
        const released = new Date(data.Time);
        return released < before ? ([released.toISOString(), ver] as const) : null;
      }),
    );

    const candidates = entries
      .filter((e): e is [string, string] => e !== null)
      .sort(([a], [b]) => new Date(b).getTime() - new Date(a).getTime());
    return candidates.length > 0 ? candidates[0][1] : null;
  }

  async function adviseViolation(
    module: string,
    version: string,
    releasedAt: Date,
  ): Promise<void> {
    console.log(`\n  Investigating fix for ${module}@${version} …`);

    let graphOutput: string;
    try {
      graphOutput = execSync('go mod graph', { cwd: root }).toString();
    } catch {
      console.log(`  ${WARN} Could not run \`go mod graph\`.`);
      return;
    }

    const requirers = findRequirers(graphOutput, module, version);
    if (requirers.length === 0) {
      console.log(`  ${WARN} Could not determine version requirers via \`go mod graph\`.`);
      return;
    }

    const rootIsRequirer = requirers.some(r => !r.includes('@'));
    const otherRequirers = requirers.filter(r => r.includes('@'));

    if (rootIsRequirer) {
      if (otherRequirers.length > 0) console.log(`  Also required by: ${otherRequirers.join(', ')}`);
      console.log('  Version is pinned as indirect in go.mod (set by a previous `go mod tidy`).');
      const best = await suggestDowngrade(module, releasedAt);
      if (best) {
        console.log('  Suggested fix:');
        console.log(`    go get ${module}@${best} && go mod tidy`);
      } else {
        console.log(`  ${WARN} No older version of ${module} found — wait ${MIN_DAYS} days or pin manually.`);
      }
      return;
    }

    for (const requirer of otherRequirers) {
      const reqModule = requirer.split('@')[0];
      console.log(`  Version selected because ${requirer} requires ${module}@${version}.`);
      const best = await suggestDowngrade(reqModule, releasedAt);
      if (best) {
        console.log(`  Suggested fix (downgrade ${reqModule}):`);
        console.log(`    go get ${reqModule}@${best} && go mod tidy`);
        console.log(`  Note: verify that ${best} does not itself require ${module}@${version}.`);
      } else {
        console.log(`  ${WARN} No older version of ${reqModule} found — wait ${MIN_DAYS} days or pin manually.`);
      }
    }
  }

  // Entry point
  const violations = await checkGoModules();
  console.log();

  if (violations.length > 0) {
    console.log(`\x1b[31m❌ ${violations.length} dependency/ies are too new (< ${MIN_DAYS} days old):\x1b[0m`);
    for (const { module, version } of violations) console.log(`  • [Go] ${module}@${version}`);
    console.log('\n\x1b[33mSuggested fixes:\x1b[0m');
    await Promise.all(violations.map(({ module, version, releasedAt }) =>
      adviseViolation(module, version, releasedAt),
    ));
    process.exit(1);
  } else {
    console.log(`\x1b[32m✅ All Go dependencies pass the ${MIN_DAYS}-day age check.\x1b[0m`);
  }
}

#!/usr/bin/env node
/**
 * Syncs frontend dependencies by running pnpm install, and detects packages
 * that violate the minimumReleaseAge constraint.
 *
 * Default (no flags): reports the offending package with a suggested fix, then
 * exits with code 1. Does NOT modify package.json. Safe for automated use.
 *
 * --pin: automatically pins violating packages to an older compliant version via
 * pnpm.overrides in package.json and retries until all violations are resolved.
 * Modifies package.json, so must only be run with explicit human intent
 * (via `pnpm run sync:pin`).
 *
 * Usage: node scripts/sync-frontend.ts [--pin]
 * Tests: node --test scripts/sync-frontend.ts
 */

import { execSync } from 'node:child_process';
import { readFileSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

// ─── Types ────────────────────────────────────────────────────────────────────

interface PackageJson {
  pnpm?: { overrides?: Record<string, string> };
  [key: string]: unknown;
}

// ─── Pure functions ───────────────────────────────────────────────────────────

/** Parse minimumReleaseAge (minutes) from pnpm-workspace.yaml content. */
export function parseMinimumReleaseAge(yaml: string): number | null {
  const match = yaml.match(/^minimumReleaseAge:\s*(\d+)/m);
  return match ? parseInt(match[1]) : null;
}

/**
 * Given the npm registry time map for a package and a cutoff date,
 * return the latest version published before the cutoff, or null if none exists.
 */
export function findCompliantVersion(
  times: Record<string, string>,
  cutoff: Date,
): string | null {
  const candidates = Object.entries(times)
    .filter(([v]) => v !== 'created' && v !== 'modified')
    .filter(([, t]) => new Date(t) < cutoff)
    .sort(([, a], [, b]) => new Date(b).getTime() - new Date(a).getTime());
  return candidates.length > 0 ? candidates[0][0] : null;
}

/** Return a new PackageJson with the given pnpm override applied (non-mutating). */
export function applyOverride(
  pkgJson: PackageJson,
  pkg: string,
  version: string,
): PackageJson {
  return {
    ...pkgJson,
    pnpm: {
      ...pkgJson.pnpm,
      overrides: { ...pkgJson.pnpm?.overrides, [pkg]: version },
    },
  };
}

// ─── Tests ───────────────────────────────────────────────────────────────────
// Run with: node --test scripts/sync-frontend.ts

if (process.env.NODE_TEST_CONTEXT) {
  const { test } = await import('node:test');
  const assert = (await import('node:assert/strict')).default;

  // parseMinimumReleaseAge
  test('parseMinimumReleaseAge: returns minutes when present', () => {
    assert.equal(parseMinimumReleaseAge('minimumReleaseAge: 10080\n'), 10080);
  });
  test('parseMinimumReleaseAge: returns null when absent', () => {
    assert.equal(parseMinimumReleaseAge('allowBuilds:\n  esbuild: true\n'), null);
  });

  // findCompliantVersion
  const cutoff = new Date('2026-03-30T00:00:00Z');
  const times: Record<string, string> = {
    created:  '2020-01-01T00:00:00Z',
    modified: '2026-04-01T00:00:00Z',
    '1.0.0':  '2026-03-28T00:00:00Z', // before cutoff → compliant
    '1.0.1':  '2026-03-29T00:00:00Z', // before cutoff → compliant (newest)
    '1.0.2':  '2026-04-01T00:00:00Z', // after cutoff  → violation
  };

  test('findCompliantVersion: returns latest compliant version', () => {
    assert.equal(findCompliantVersion(times, cutoff), '1.0.1');
  });
  test('findCompliantVersion: returns null when no compliant version exists', () => {
    assert.equal(findCompliantVersion({ '1.0.0': '2026-04-01T00:00:00Z' }, cutoff), null);
  });
  test('findCompliantVersion: ignores created/modified keys', () => {
    const onlyMeta = { created: '2020-01-01T00:00:00Z', modified: '2020-01-01T00:00:00Z' };
    assert.equal(findCompliantVersion(onlyMeta, cutoff), null);
  });

  // applyOverride
  test('applyOverride: adds override to empty pnpm field', () => {
    const result = applyOverride({}, '@emnapi/core', '1.9.1');
    assert.deepEqual(result.pnpm?.overrides, { '@emnapi/core': '1.9.1' });
  });
  test('applyOverride: merges with existing overrides', () => {
    const input: PackageJson = { pnpm: { overrides: { browserslist: '4.28.1' } } };
    const result = applyOverride(input, '@emnapi/core', '1.9.1');
    assert.deepEqual(result.pnpm?.overrides, {
      browserslist: '4.28.1',
      '@emnapi/core': '1.9.1',
    });
  });
  test('applyOverride: does not mutate input', () => {
    const input: PackageJson = { pnpm: { overrides: { browserslist: '4.28.1' } } };
    applyOverride(input, '@emnapi/core', '1.9.1');
    assert.deepEqual(input.pnpm?.overrides, { browserslist: '4.28.1' });
  });

}

// ─── Main ─────────────────────────────────────────────────────────────────────

if (!process.env.NODE_TEST_CONTEXT) {
  const pin = process.argv.includes('--pin');
  const rootDir = resolve(dirname(fileURLToPath(import.meta.url)), '..', 'frontend');

  const workspaceYaml = readFileSync(resolve(rootDir, 'pnpm-workspace.yaml'), 'utf8');
  const minAgeMinutes = parseMinimumReleaseAge(workspaceYaml);
  if (!minAgeMinutes) {
    console.log('minimumReleaseAge not set in pnpm-workspace.yaml, nothing to do.');
    process.exit(0);
  }
  const cutoff = new Date(Date.now() - minAgeMinutes * 60 * 1000);
  console.log(`minimumReleaseAge: ${minAgeMinutes} min  (cutoff: ${cutoff.toISOString()})`);

  let iterations = 0;
  while (true) {
    iterations++;
    let output: string;
    try {
      output = execSync('pnpm install', { cwd: rootDir, stdio: 'pipe' }).toString();
      console.log(output);
      console.log('pnpm install succeeded.');
      break;
    } catch (e) {
      const err = e as { stdout?: Buffer; stderr?: Buffer };
      output = (err.stdout?.toString() ?? '') + (err.stderr?.toString() ?? '');
    }

    const match = output.match(/ERR_PNPM_NO_MATURE_MATCHING_VERSION.*?of (\S+) does not meet/);
    if (!match) {
      process.stderr.write(output);
      process.exit(1);
    }

    const pkg = match[1];
    console.error(`\nViolation detected: ${pkg}`);

    const raw = execSync(`pnpm view ${pkg} time --json`, { cwd: rootDir }).toString();
    const version = findCompliantVersion(JSON.parse(raw) as Record<string, string>, cutoff);
    if (!version) {
      console.error(`No compliant version found for ${pkg} (cutoff: ${cutoff.toISOString()})`);
      process.exit(1);
    }

    if (pin) {
      console.log(`Pinning to: ${pkg}@${version}`);
      const pkgJsonPath = resolve(rootDir, 'package.json');
      const pkgJson = JSON.parse(readFileSync(pkgJsonPath, 'utf8')) as PackageJson;
      writeFileSync(pkgJsonPath, JSON.stringify(applyOverride(pkgJson, pkg, version), null, 2) + '\n');
    } else {
      console.error(`Suggested fix: add the following to frontend/package.json under pnpm.overrides:`);
      console.error(`  "${pkg}": "${version}"`);
      console.error(`Or run: pnpm run sync:pin`);
      process.exit(1);
    }
  }

  if (pin) console.log(`\nDone in ${iterations} iteration(s).`);
}

/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
/* eslint-disable no-console */

const fs = require('fs');
const path = require('path');
const { execSync, spawnSync } = require('child_process');

const ROOT_DIR = path.join(__dirname, '..');
const PKG_PATH = path.join(ROOT_DIR, 'package.json');
const LOCK_PATH = path.join(ROOT_DIR, 'pnpm-lock.yaml');
const LIST_MAX_BUFFER = 1024 * 1024 * 100;
const LIST_DEPTHS = ['Infinity', 8, 6, 4, 2];

/**
 * Validates that a package name follows npm naming conventions.
 * Prevents command injection by ensuring only safe characters are present.
 * @param {string} name - The package name to validate
 * @returns {boolean} - True if valid, false otherwise
 */
function isValidPackageName(name) {
  if (!name || typeof name !== 'string') return false;

  // npm package names can contain:
  // - lowercase letters, numbers, hyphens, underscores, dots
  // - @ symbol (for scoped packages)
  // - / (for scoped packages)
  // Must not contain shell metacharacters: $, `, ;, |, &, <, >, (, ), etc.
  const validPattern = /^[@a-z0-9\-_./]+$/i;

  // Additional checks
  const maxLength = 214; // npm package name limit
  if (name.length > maxLength) return false;

  return validPattern.test(name);
}

/**
 * Simple exact-version comparison.
 * Returns true if the resolved version is older than the target override.
 */
function isOlder(found, target) {
  if (!found || !target) return false;
  const f = found
    .replace(/[^0-9.]/g, '')
    .split('.')
    .map((n) => parseInt(n || 0));
  const t = target
    .replace(/[^0-9.]/g, '')
    .split('.')
    .map((n) => parseInt(n || 0));
  for (let i = 0; i < Math.max(f.length, t.length); i++) {
    if ((f[i] || 0) < (t[i] || 0)) return true;
    if ((f[i] || 0) > (t[i] || 0)) return false;
  }
  return false;
}

function writeReport(report, error) {
  const outPath = path.join(ROOT_DIR, 'DEP_OVERRIDE_REPORT.md');
  fs.writeFileSync(outPath, report);
  const message =
    error || '✅ Dependency override audit complete! Project has been restored to its original state.';
  console.log(message);
  console.log('📄 See DEP_OVERRIDE_REPORT.md for details.');
}

function getCommandOutput(error) {
  const stdout = error.stdout ? error.stdout.toString().trim() : '';
  const stderr = error.stderr ? error.stderr.toString().trim() : '';
  return { stdout, stderr, errorMsg: stderr || stdout || error.message };
}

function shouldRetryWithReducedDepth(error) {
  const { errorMsg } = getCommandOutput(error);
  return error.code === 'ENOBUFS' || errorMsg.includes('Invalid string length');
}

function getDependencyTree(packageName) {
  // Validate package name to prevent command injection
  if (!isValidPackageName(packageName)) {
    throw new Error(
      `Invalid package name: "${packageName}". Package names must contain only alphanumeric characters, hyphens, underscores, dots, @, and /.`
    );
  }

  let lastError;

  for (const depth of LIST_DEPTHS) {
    try {
      // Use spawnSync with array arguments to prevent shell injection
      const result = spawnSync(
        'pnpm',
        ['list', packageName, '--recursive', '--depth', String(depth), '--json'],
        {
          cwd: ROOT_DIR,
          maxBuffer: LIST_MAX_BUFFER,
          encoding: 'utf8',
        }
      );

      // Check for errors
      if (result.error) {
        throw result.error;
      }

      const stdout = result.stdout || '';
      const stderr = result.stderr || '';

      // pnpm exits with code 1 and an empty JSON array when the package is absent
      if (stdout.trim() === '[]') {
        return { orphaned: true };
      }

      // If there's an error status, create an error object similar to execSync
      if (result.status !== 0) {
        const error = new Error(`Command failed with exit code ${result.status}`);
        error.code = result.status === null ? 'ENOBUFS' : undefined;
        error.stdout = Buffer.from(stdout);
        error.stderr = Buffer.from(stderr);
        throw error;
      }

      return {
        data: JSON.parse(stdout),
        usedFallbackDepth: depth !== 'Infinity',
        depth,
      };
    } catch (error) {
      const { stdout } = getCommandOutput(error);

      // pnpm exits with code 1 and an empty JSON array when the package is absent.
      if (stdout === '[]') {
        return { orphaned: true };
      }

      if (!shouldRetryWithReducedDepth(error)) {
        throw error;
      }

      lastError = error;
    }
  }

  throw lastError;
}

function genOverrideReport() {
  console.log('🚀 Starting Dependency Override Audit. Backing up package.json...');

  let report = '# 🛡️ PNPM Override Audit Report\n\n';
  report += `Generated on: ${new Date().toLocaleString()}\n\n`;

  if (!fs.existsSync(PKG_PATH)) {
    report += `❌ **Error:** package.json not found at ${PKG_PATH}\n`;
    writeReport(report, '❌ package.json not found.');
    return;
  }

  // 1. Read and backup the original state
  const originalPkgStr = fs.readFileSync(PKG_PATH, 'utf8');
  const originalLockStr = fs.existsSync(LOCK_PATH) ? fs.readFileSync(LOCK_PATH, 'utf8') : null;

  const pkgJson = JSON.parse(originalPkgStr);
  const overrides = pkgJson.pnpm?.overrides || {};

  if (Object.keys(overrides).length === 0) {
    report += '✅ **No overrides found in package.json.**\n';
    writeReport(report);
    return;
  }

  try {
    // 2. Strip overrides and save the modified package.json
    const tempPkgJson = JSON.parse(originalPkgStr);
    delete tempPkgJson.pnpm.overrides;
    fs.writeFileSync(PKG_PATH, JSON.stringify(tempPkgJson, null, 2));

    // 3. Reinstall dependencies to recalculate lockfile AND physically update node_modules
    console.log('⏳ Relinking node_modules to their natural state (this may take a minute)...');
    execSync('pnpm install --no-frozen-lockfile --ignore-scripts', {
      cwd: ROOT_DIR,
      stdio: 'ignore',
    });

    // 4. Audit each removed override using pnpm list
    for (const [overrideName, targetVersion] of Object.entries(overrides)) {
      console.log(`🔎 Auditing natural resolution for ${overrideName}...`);
      const culprits = new Map();

      let treeResult;
      try {
        treeResult = getDependencyTree(overrideName);
      } catch (e) {
        console.error(`└──⚠️ Could not fetch tree for ${overrideName}.`);
        const { errorMsg } = getCommandOutput(e);

        report += `## \`${overrideName}\`\n**Target Override:** \`${targetVersion}\`\n\n`;

        report += `❓ **UNKNOWN (Error)**\n\n`;
        report += `> The script encountered an error resolving this package:\n> \`${errorMsg}\`\n`;

        report += `\n---\n`;
        continue; // Immediately jump to the next override in the loop
      }

      report += `## \`${overrideName}\`\n**Target Override:** \`${targetVersion}\`\n\n`;

      if (treeResult.orphaned) {
        report += `✅ **SAFE TO REMOVE (Orphaned)**\n\n`;
        report += `> This package does not exist anywhere in the naturally resolved dependency tree. It was likely removed by an upstream dependency update.\n`;
        report += `\n---\n`;
        continue;
      }

      const data = treeResult.data;

      const scanTree = (parentName, parentVersion, depsObject) => {
        if (!depsObject) return;

        for (const [depName, depInfo] of Object.entries(depsObject)) {
          if (!depInfo) continue;

          if (depName === overrideName && depInfo.version) {
            const resolvedVersion = depInfo.version;

            if (isOlder(resolvedVersion, targetVersion)) {
              culprits.set(`${parentName}@${parentVersion}`, resolvedVersion);
            }
          }

          // If this dependency has its own dependencies, it becomes the new parent
          if (depInfo.dependencies) {
            scanTree(depName, depInfo.version, depInfo.dependencies);
          }
        }
      };

      // Start the scan from the top-level workspaces/projects
      data.forEach((project) => {
        const projectName = project.name || 'Root Project';
        const projectVersion = project.version || 'unknown';
        const allDeps = {
          ...project.dependencies,
          ...project.devDependencies,
          ...project.optionalDependencies,
        };

        scanTree(projectName, projectVersion, allDeps);
      });

      // Generate markdown segment
      if (culprits.size > 0) {
        if (treeResult.usedFallbackDepth) {
          report += `> pnpm could not serialize the full recursive tree for this package, so this result is based on a bounded recursive scan to depth ${treeResult.depth}. Older versions were found within that scanned portion of the tree.\n\n`;
        }

        report += `⚠️ **REQUIRED**\n\n`;
        report += `> These packages will continue to receive the overridden version until they are updated to naturally resolve to >= ${targetVersion}.\n\n`;
        report += `| Parent Package | Naturally Resolved Version |\n| :--- | :--- |\n`;
        const sortedCulprits = Array.from(culprits.entries()).sort();
        for (const [parent, resolved] of sortedCulprits) {
          report += `| \`${parent}\` | \`${resolved}\` |\n`;
        }
      } else {
        if (treeResult.usedFallbackDepth) {
          report += `> pnpm could not serialize the full recursive tree for this package, so this result is based on a bounded recursive scan to depth ${treeResult.depth}. No older versions were found within that scanned portion of the tree, but deeper paths were not inspected.\n\n`;
        }

        report += `✅ **SAFE TO REMOVE**\n\n`;
        report += `> All packages naturally resolve to >= ${targetVersion} without the override.\n`;
      }
      report += `\n---\n`;
    }
  } finally {
    // 5. Restore original package.json, lockfile, and node_modules
    console.log('🧹 Cleaning up: Restoring package.json, lockfile, and node_modules...');

    // Put the files back
    fs.writeFileSync(PKG_PATH, originalPkgStr);
    if (originalLockStr) {
      fs.writeFileSync(LOCK_PATH, originalLockStr);
    }

    // Run a full install again to force pnpm to re-apply overrides to node_modules
    try {
      execSync('pnpm install --no-frozen-lockfile --ignore-scripts', {
        cwd: ROOT_DIR,
        stdio: 'ignore',
      });
    } catch (e) {
      console.error("⚠️ Cleanup failed, you may need to run 'pnpm install' manually.");
    }

    // Write the report
    writeReport(report);
  }
}

genOverrideReport();

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, findAll, triggerKeyEvent, visit } from '@ember/test-helpers';

// REPL selectors
const REPL = {
  toggle: '[data-test-console-toggle]',
  consoleInput: '[data-test-component="console/command-input"] input',
  logOutputItems: '[data-test-component="console/output-log"] > div',
};

/**
 * Helper functions to run common commands in the consoleComponent during tests.
 * Note: A user must be logged in during the test context for the commands to run.
 *
 * Example:
 * import { v4 as uuidv4 } from 'uuid';
 * import { runCmd, mountEngineCmd } from 'vault/tests/helpers/commands';
 *
 * async function mountEngine() {
 *   const backend = `pki-${uuidv4()}`;
 *   await runCmd(mountEngineCmd('pki', backend));
 *   return backend;
 * }
 */

// Command execution helpers
export const runCmd = async (commands, throwErrors = true) => {
  if (!commands) throw new Error('runCmd requires commands array passed in');
  if (!Array.isArray(commands)) commands = [commands];

  await click(REPL.toggle);
  await enterCommands(commands);
  const lastOutput = await lastLogOutput();
  await click(REPL.toggle);

  if (throwErrors && lastOutput.includes('Error')) {
    throw new Error(`Error occurred while running commands: "${commands.join('; ')}" - ${lastOutput}`);
  }

  return lastOutput;
};

export const enterCommands = async (commands) => {
  const toExecute = Array.isArray(commands) ? commands : [commands];
  for (const command of toExecute) {
    await fillIn(REPL.consoleInput, command);
    await triggerKeyEvent(REPL.consoleInput, 'keyup', 'Enter');
  }
};

export const lastLogOutput = async () => {
  const items = findAll(REPL.logOutputItems);
  if (!items.length) return '';
  return items[items.length - 1].innerText;
};

// Command builders
export const mountEngineCmd = (type, customName = '') => {
  const name = customName || type;
  return type === 'kv-v2'
    ? `write sys/mounts/${name} type=kv options=version=2`
    : `write sys/mounts/${name} type=${type}`;
};

export const deleteEngineCmd = (name) => `delete sys/mounts/${name}`;

export const mountAuthCmd = (type, customName = '') => {
  const name = customName || type;
  return `write sys/auth/${name} type=${type}`;
};

export const deleteAuthCmd = (name) => `delete sys/auth/${name}`;

export const createPolicyCmd = (name, contents) => {
  const policyContent = window.btoa(contents);
  return `write sys/policies/acl/${name} policy=${policyContent}`;
};

export const createTokenCmd = (policyName = 'default') =>
  `write -field=client_token auth/token/create policies=${policyName} ttl=1h`;

export const tokenWithPolicyCmd = (name, policy) => [createPolicyCmd(name, policy), createTokenCmd(name)];

export const createNS = (namespace) => `write sys/namespaces/${namespace} -f`;

export const deleteNS = (namespace) => `delete sys/namespaces/${namespace} -f`;

/**
 * @description
 * Iterates over an array of namespace paths and ensures each nested level is created.
 * It visits the root namespace before attempting to create the next segment.
 *
 * @example input: ['foo/bar', 'baz/qux/quux']
 * This will create: foo, foo/bar, baz, baz/qux, baz/qux/quux
 *
 * @param {string[]} namespaces - Array of strings of namespace paths (containing backslashes)
 */
export const createNSFromPaths = async (namespaces) => {
  for (const ns of namespaces) {
    const parts = ns.split('/');
    let currentPath = '';

    for (const part of parts) {
      const url = `/vault/dashboard${currentPath && `?namespace=${currentPath.replaceAll('/', '%2F')}`}`;
      await visit(url);

      currentPath = currentPath ? `${currentPath}/${part}` : part;
      await runCmd(createNS(part), false);
    }

    // Reset to root namespace after creating each path
    await visit('/vault/dashboard');
  }
};

/**
 * @description
 * Deletes namespaces by removing each segment of the path from deepest to top.
 *
 * @example input: ['foo/bar', 'baz/qux']
 * This will delete: foo, baz
 *
 * @param {string[]} namespaces - Array of strings of namespace paths (containing backslashes)
 */
export const deleteNSFromPaths = async (namespaces) => {
  for (const ns of namespaces) {
    const parts = ns.split('/');
    // Work from deepest child up to the top-level namespace
    for (let i = parts.length - 1; i >= 0; i--) {
      const parentPath = parts.slice(0, i).join('/');
      const toDelete = parts[i];
      // Build the URL for the parent namespace (or root if none)
      const url = parentPath
        ? `/vault/dashboard?namespace=${parentPath.replaceAll('/', '%2F')}`
        : '/vault/dashboard';
      await visit(url);
      await runCmd(deleteNS(toDelete), false);
    }
    // Reset to root namespace after deleting each path
    await visit('/vault/dashboard');
  }
};

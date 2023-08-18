import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { create } from 'ember-cli-page-object';
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Helper functions to run common commands in the consoleComponent during tests.
 * Please note that a user must be logged in during the test context for the commands to run.
 * By default runCmd throws an error if the last log includes "Error". To override this,
 * pass boolean false to run the commands and not throw errors
 *
 * Example:
 *
 * import { v4 as uuidv4 } from 'uuid';
 * import { runCmd, mountEngineCmd } from 'vault/tests/helpers/commands';
 *
 *
 * async function mountEngineExitOnError() {
 *    const backend = `pki-${uuidv4()}`;
 *    await runCmd(mountEngineCmd('pki', backend));
 *    return backend;
 * }
 *
 * async function mountEngineSquashErrors() {
 *    const backend = `pki-${uuidv4()}`;
 *    await runCmd(mountEngineCmd('pki', backend), false);
 *    return backend;
 * }
 */

const cc = create(consoleClass);

/**
 * runCmd is used to run commands and throw an error if the output includes "Error"
 * @param {string || string[]} commands array of commands that should run
 * @param {boolean} throwErrors
 * @returns the last log output. Throws an error if it includes an error
 */
export async function runCmd(commands, throwErrors = true) {
  if (!commands) {
    throw new Error('runCmd requires commands array passed in');
  }
  if (!Array.isArray(commands)) {
    commands = [commands];
  }
  await cc.runCommands(commands);
  const lastOutput = cc.lastLogOutput;
  if (throwErrors && lastOutput.includes('Error')) {
    throw new Error(`Error occurred while running commands: "${commands.join('; ')}" - ${lastOutput}`);
  }
  return lastOutput;
}

// Common commands
export function mountEngineCmd(type, customName = '') {
  const name = customName || type;
  if (type === 'kv-v2') {
    return `write sys/mounts/${name} type=kv options=version=2`;
  }
  return `write sys/mounts/${name} type=${type}`;
}

export function deleteEngineCmd(name) {
  return `delete sys/mounts/${name}`;
}

export function mountAuthCmd(type, customName = '') {
  const name = customName || type;
  return `write sys/auth/${name} type=${type}`;
}

export function deleteAuthCmd(name) {
  return `delete sys/auth/${name}`;
}

export function createPolicyCmd(name, contents) {
  const policyContent = window.btoa(contents);
  return `write sys/policies/acl/${name} policy=${policyContent}`;
}

export function createTokenCmd(policyName = 'default') {
  return `write -field=client_token auth/token/create policies=${policyName} ttl=1h`;
}

export const tokenWithPolicyCmd = function (name, policy) {
  return [createPolicyCmd(name, policy), createTokenCmd(name)];
};

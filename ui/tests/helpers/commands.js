/**
 * Helper functions to run common commands in the consoleComponent during tests.
 * Please note that a user must be logged in during the test context for the commands to run
 *
 * Example:
 *
 * import { v4 as uuidv4 } from 'uuid';
 * import { create } from 'ember-cli-page-object';
 * import consoleClass from 'vault/tests/pages/components/console/ui-panel';
 * import { runCmd, mountEngineCmd } from 'vault/tests/helpers/commands';
 *
 * const consoleComponent = create(consoleClass);
 *
 * async function mountEngineExitOnError() {
 *    const backend = `pki-${uuidv4()}`;
 *    await runCmd(consoleComponent, mountEngineCmd('pki', backend));
 *    return backend;
 * }
 *
 * async function mountEngineSquashErrors() {
 *    const backend = `pki-${uuidv4()}`;
 *    await consoleComponent.runCommands(mountEngineCmd('pki', backend));
 *    return backend;
 * }
 */

export function mountEngineCmd(type, customName = '') {
  const name = customName || type;
  return [`write sys/mounts/${name} type=${type}`];
}

export function deleteEngineCmd(name) {
  return [`delete sys/mounts/${name}`];
}

export function createPolicyCmd(name, contents) {
  const policyContent = window.btoa(contents);
  return [`write sys/policies/acl/${name} policy=${policyContent}`];
}

export function tokenWithPolicyCmd(policyName = 'default') {
  return [`write -field=client_token auth/token/create policies=${policyName} ttl=1h`];
}

/**
 * runCmd is used to run commands and throw an error if the output includes "Error"
 * @param {Component} console instance
 * @param {Array<string>} commands array of commands that should be run
 * @returns the last log output. Throws an error if it includes an error
 */
export async function runCmd(component, commands) {
  if (!component || !commands) {
    throw new Error('runCmd requires consoleComponent and commands passed in');
  }
  await component.runCommands(commands);
  const lastOutput = component.lastLogOutput;
  if (lastOutput.includes('Error')) {
    throw new Error(`Error occurred while running commands: "${commands.join('; ')}" - ${lastOutput}`);
  }
  return lastOutput;
}

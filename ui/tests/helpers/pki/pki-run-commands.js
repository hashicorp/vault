/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { create } from 'ember-cli-page-object';

const consoleComponent = create(consoleClass);

export const tokenWithPolicy = async function (name, policy) {
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write -field=client_token auth/token/create policies=${name}`,
  ]);
  return consoleComponent.lastLogOutput;
};

export const runCommands = async function (commands) {
  try {
    await consoleComponent.runCommands(commands);
    const res = consoleComponent.lastLogOutput;
    if (res.includes('Error')) {
      throw new Error(res);
    }
    return res;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error(
      `The following occurred when trying to run the command(s):\n ${commands.join('\n')} \n\n ${
        consoleComponent.lastLogOutput
      }`
    );
    throw error;
  }
};

// Clears pki-related data and capabilities so that admin
// capabilities from setup don't rollover
export function clearRecords(store) {
  store.unloadAll('pki/action');
  store.unloadAll('pki/issuer');
  store.unloadAll('pki/key');
  store.unloadAll('pki/role');
  store.unloadAll('pki/sign-intermediate');
  store.unloadAll('pki/tidy');
  store.unloadAll('pki/config/urls');
  store.unloadAll('pki/config/crl');
  store.unloadAll('pki/config/cluster');
  store.unloadAll('pki/config/acme');
  store.unloadAll('pki/certificate/generate');
  store.unloadAll('pki/certificate/sign');
  store.unloadAll('capabilities');
}

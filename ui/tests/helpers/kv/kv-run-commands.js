/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { create } from 'ember-cli-page-object';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import listPage from 'vault/tests/pages/secrets/backend/list';

// TODO update writeSecret with create flow from KV ember engine when built
export const writeSecret = async function (backend, path, key, val) {
  await listPage.visitRoot({ backend });
  await listPage.create();
  return editPage.createSecret(path, key, val);
};

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

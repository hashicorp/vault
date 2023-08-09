/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-env node */
/* eslint-disable no-console */

const testHelper = require('./test-helper');

(async function () {
  try {
    let unsealKeys = process.env.VAULT_UNSEAL_KEYS;
    if (!unsealKeys) {
      console.error(
        'Cannot run ember tests without unseal keys, please make sure to export the keys, in an env ' +
          'var named: VAULT_UNSEAL_KEYS'
      );
      process.exit(1);
    } else {
      unsealKeys = JSON.parse(unsealKeys);
    }

    const rootToken = process.env.VAULT_TOKEN;
    if (!rootToken) {
      console.error(
        'Cannot run ember tests without root token, please make sure to export the root token, in an env ' +
          'var named: VAULT_TOKEN'
      );
      process.exit(1);
    }

    testHelper.writeKeysFile(unsealKeys, rootToken);
  } catch (error) {
    console.log(error);
    process.exit(1);
  }

  const vaultAddr = process.env.VAULT_ADDR;
  if (!vaultAddr) {
    console.error(
      'Cannot run ember tests without the Vault Address, please make sure to export the vault address, in an env ' +
        'var named: VAULT_ADDR'
    );
    process.exit(1);
  }

  console.log('VAULT_ADDR=' + vaultAddr);

  try {
    const testArgs = ['test', '-c', 'testem.enos.js'];

    if (process.env.TEST_FILTER && process.env.TEST_FILTER.length > 0) {
      testArgs.push('-f=' + process.env.TEST_FILTER);
    }

    await testHelper.run('ember', [...testArgs, ...process.argv.slice(2)], false);
  } catch (error) {
    console.log(error);
    process.exit(1);
  } finally {
    process.exit(0);
  }
})();

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

    if (process.env.TEST_FILTERS) {
      const filters = JSON.parse(process.env.TEST_FILTERS).map((filter) => '-f=' + filter);
      testArgs.push(...filters);
    }

    await testHelper.run('ember', testArgs, false);
  } catch (error) {
    console.log(error);
    process.exit(1);
  } finally {
    process.exit(0);
  }
})();

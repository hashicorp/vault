/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
/* eslint-disable no-console */
/* eslint-disable no-process-exit */
/* eslint-disable n/no-extraneous-require */

const testHelper = require('./test-helper');

(async function () {
  // ignore first 2 args (node and path) and extract flags to pass to test/exam command
  const args = process.argv.slice(2);
  // in CI use local vault binary, otherwise assume vault is in PATH
  const vaultCommand = process.env.CI ? '../bin/vault' : 'vault';
  try {
    testHelper.run(
      vaultCommand,
      [
        'server',
        '-dev',
        '-dev-ha',
        '-dev-transactional',
        '-dev-root-token-id=root',
        '-dev-listen-address=127.0.0.1:9200',
      ],
      false
    );
    try {
      const withServer = args.includes('--server') || args.includes('-s');
      // current issue with headless Chrome where an event listener in Hds::Modal is not triggered resulting in a pending test waiter and timeout
      // the workaround for now is to run the tests in headless firefox for local runs
      if (!withServer && !process.env.CI) {
        args.push('--launch=Firefox');
      }
      await testHelper.run('ember', ['exam', ...args]);
    } catch (error) {
      console.log(error);
      process.exit(1);
    } finally {
      process.exit(0);
    }
  } catch (error) {
    console.log(error);
    process.exit(1);
  }
})();

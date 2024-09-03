/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
/* eslint-disable no-console */
/* eslint-disable no-process-exit */
/* eslint-disable n/no-extraneous-require */

var readline = require('readline');
const testHelper = require('./test-helper');

var output = '';
var unseal, root, written, initError;

async function processLines(input, eachLine = () => {}) {
  const rl = readline.createInterface({
    input,
    terminal: true,
  });
  for await (const line of rl) {
    eachLine(line);
  }
}

(async function () {
  try {
    const vault = testHelper.run(
      'vault',
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
    processLines(vault.stdout, function (line) {
      if (written) {
        output = null;
        return;
      }
      output = output + line;
      var unsealMatch = output.match(/Unseal Key: (.+)$/m);
      if (unsealMatch && !unseal) {
        unseal = [unsealMatch[1]];
      }
      var rootMatch = output.match(/Root Token: (.+)$/m);
      if (rootMatch && !root) {
        root = rootMatch[1];
      }
      var errorMatch = output.match(/Error initializing core: (.*)$/m);
      if (errorMatch) {
        initError = errorMatch[1];
      }
      if (root && unseal && !written) {
        testHelper.writeKeysFile(unseal, root);
        written = true;
        console.log('VAULT SERVER READY');
      } else if (initError) {
        console.log('VAULT SERVER START FAILED');
        console.log(
          'If this is happening, run `export VAULT_LICENSE_PATH=/Users/username/license.hclic` to your valid local vault license filepath, or use OSS Vault'
        );
        process.exit(1);
      }
    });
    try {
      // only the test:filter command specifies --server by default
      const verb = process.argv[2] === '--server' ? 'test' : 'exam';
      await testHelper.run('ember', [verb, ...process.argv.slice(2)]);
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

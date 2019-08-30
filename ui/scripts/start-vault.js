#!/usr/bin/env node
/* eslint-env node */
/* eslint-disable no-console */

var fs = require('fs');
var path = require('path');
var readline = require('readline');
var execa = require('execa');
var chalk = require('chalk');

function run(command, args = [], shareStd = true) {
  console.log(chalk.dim('$ ' + command + ' ' + args.join(' ')));
  // cleanup means that execa will handle stopping the vault subprocess
  // inherit all of the stdin/out/err so that testem still works as if you were running it directly
  if (shareStd) {
    return execa(command, args, { cleanup: true, stdin: 'inherit', stdout: 'inherit', stderr: 'inherit' });
  }
  let p = execa(command, args, { cleanup: true });
  p.stdout.pipe(process.stdout);
  p.stderr.pipe(process.stderr);
  return p;
}

var output = '';
var unseal, root, written;

async function processLines(input, eachLine = () => {}) {
  const rl = readline.createInterface({
    input,
    terminal: true,
  });
  for await (const line of rl) {
    eachLine(line);
  }
}

(async function() {
  try {
    let vault = run(
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

    processLines(vault.stdout, function(line) {
      if (written) {
        output = null;
        return;
      }
      output = output + line;
      var unsealMatch = output.match(/Unseal Key: (.+)$/m);
      if (unsealMatch && !unseal) {
        unseal = unsealMatch[1];
      }
      var rootMatch = output.match(/Root Token: (.+)$/m);
      if (rootMatch && !root) {
        root = rootMatch[1];
      }
      if (root && unseal && !written) {
        fs.writeFile(
          path.join(process.cwd(), 'tests/helpers/vault-keys.js'),
          `export default ${JSON.stringify({ unseal, root }, null, 2)}`,
          err => {
            if (err) throw err;
          }
        );
        written = true;
        console.log('VAULT SERVER READY');
      }
    });
    try {
      if (process.argv[2] === '--browserstack') {
        await run('ember', ['browserstack:connect']);
        try {
          await run('ember', ['test', '-f=secrets/secret/create', '-c', 'testem.browserstack.js']);

          console.log('success');
          process.exit(0);
        } finally {
          if (process.env.CI === 'true') {
            await run('ember', ['browserstack:results']);
          }
          await run('ember', ['browserstack:disconnect']);
        }
      } else {
        await run('ember', ['test', ...process.argv.slice(2)]);
      }
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

#!/usr/bin/env node
/* eslint-env node */
/* eslint-disable no-console */

const execa = require('execa');
const chalk = require('chalk');

function run(command, args = []) {
  console.log(chalk.dim('$ ' + command + ' ' + args.join(' ')));

  let p = execa(command, args);
  p.stdout.pipe(process.stdout);
  p.stderr.pipe(process.stderr);

  return p;
}

(async function() {
  try {
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
  } catch (error) {
    console.log('error');
    console.log(error);
    process.exit(1);
  }
})();

#!/usr/bin/env node
/* eslint-disable */

if (process.argv[2]) {
  process.kill(process.argv[2], 'SIGINT');
  process.exit(0);
}

process.env.TERM = 'dumb';
var fs = require('fs');
var path = require('path');
var readline = require('readline');
var spawn = require('child_process').spawn;
var vault = spawn('vault', [
  'server',
  '-dev',
  '-dev-ha',
  '-dev-transactional',
  '-dev-root-token-id=root',
  '-dev-listen-address=127.0.0.1:9200',
]);

var output = '';
var unseal, root;

readline
  .createInterface({
    input: vault.stdout,
    terminal: false,
  })
  .on('line', function(line) {
    output = output + line;
    console.log(line);
    var unsealMatch = output.match(/Unseal Key: (.+)$/m);
    if (unsealMatch && !unseal) {
      unseal = unsealMatch[1];
    }
    var rootMatch = output.match(/Root Token: (.+)$/m);
    if (rootMatch && !root) {
      root = rootMatch[1];
    }
    if (root && unseal) {
      fs.writeFile(
        path.join(process.cwd(), 'tests/helpers/vault-keys.js'),
        `export default ${JSON.stringify({ unseal, root }, null, 2)}`,
        err => {
          if (err) throw err;
        }
      );

      console.log('VAULT SERVER READY');
    }
  });

vault.stderr.on('data', function(data) {
  console.log(data.toString());
});

vault.on('close', function(code) {
  console.log(`child process exited with code ${code}`);
  process.exit();
});
vault.on('error', function(error) {
  console.log(`child process errored: ${error}`);
  process.exit();
});

var pidFile = 'vault-ui-integration-server.pid';
process.on('SIGINT', function() {
  vault.kill('SIGINT');
  process.exit();
});
process.on('exit', function() {
  vault.kill('SIGINT');
});

fs.writeFile(pidFile, process.pid, err => {
  if (err) throw err;
  console.log('The file has been saved!');
});

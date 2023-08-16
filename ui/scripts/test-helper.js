/* eslint-env node */
/* eslint-disable no-console */

const fs = require('fs');
const path = require('path');
const chalk = require('chalk');
const execa = require('execa');

/**
 * Writes a vault keys file that can be imported in other scripts, that includes the unseal keys and the root token.
 * @param unsealKeys an array of unseal keys, must contain at least one key
 * @param rootToken the root token
 * @param filePath optional file path, if not provided the default path of <cwd>/tests/helpers/vault-keys.js
 * will be used.
 */
function writeKeysFile(unsealKeys, rootToken, filePath) {
  if (filePath === undefined) {
    filePath = path.join(process.cwd(), 'tests/helpers/vault-keys.js');
  }
  let keys = {};
  keys.unsealKeys = unsealKeys;
  keys.rootToken = rootToken;

  fs.writeFile(filePath, `export default ${JSON.stringify(keys, null, 2)}`, (err) => {
    if (err) throw err;
  });
}

/**
 * Runs the provided command and pipes the processes stdout and stderr to the terminal. Upon completion with
 * success or error the child process will be cleaned up.
 * @param command some command to run
 * @param args some arguments for the command to run
 * @param shareStd if true the sub process created by the command will share the stdout and stderr of the parent
 * process
 * @returns {*} The child_process for the executed command which is also a Promise.
 */
function run(command, args = [], shareStd = true) {
  console.log(chalk.dim('$ ' + command + ' ' + args.join(' ')));
  // cleanup means that execa will handle stopping the subprocess
  // inherit all of the stdin/out/err so that testem still works as if you were running it directly
  if (shareStd) {
    return execa(command, args, { cleanup: true, stdin: 'inherit', stdout: 'inherit', stderr: 'inherit' });
  }
  let p = execa(command, args, { cleanup: true });
  p.stdout.pipe(process.stdout);
  p.stderr.pipe(process.stderr);
  return p;
}

module.exports = {
  writeKeysFile: writeKeysFile,
  run: run,
};

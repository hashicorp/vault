/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { REPL } from './selectors';
import { click, fillIn, findAll, triggerKeyEvent } from '@ember/test-helpers';

export const runCmd = async (commands, throwErrors = true) => {
  if (!commands) {
    throw new Error('runCmd requires commands array passed in');
  }
  if (!Array.isArray(commands)) {
    commands = [commands];
  }

  await click(REPL.toggle);
  await enterCommands(commands);
  const lastOutput = await lastLogOutput();
  await click(REPL.toggle);
  if (throwErrors && lastOutput.includes('Error')) {
    throw new Error(`Error occurred while running commands: "${commands.join('; ')}" - ${lastOutput}`);
  }
  return lastOutput;
};

export const enterCommands = async (commands) => {
  const toExecute = Array.isArray(commands) ? commands : [commands];
  for (const command of toExecute) {
    await fillIn(REPL.consoleInput, command);
    await triggerKeyEvent(REPL.consoleInput, 'keyup', 'Enter');
  }
};

export const lastLogOutput = async () => {
  const items = await findAll(REPL.logOutputItems);
  const count = items.length;
  if (count === 0) {
    // If no logOutput items are found, we can assume the response is empty
    return '';
  }
  const outputItemText = items[count - 1].innerText;
  return outputItemText;
};

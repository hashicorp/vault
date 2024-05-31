/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { text, triggerable, clickable, collection, fillable, value, isPresent } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';
import { settled } from '@ember/test-helpers';

import keys from 'core/utils/key-codes';

export default {
  toggle: clickable('[data-test-console-toggle]'),
  dismissConsole: clickable(['data-test-dismiss-console-button']),
  consoleInput: fillable('[data-test-component="console/command-input"] input'),
  consoleInputValue: value('[data-test-component="console/command-input"] input'),
  logOutput: text('[data-test-component="console/output-log"]'),
  logOutputItems: collection('[data-test-component="console/output-log"] > div', {
    text: text(),
  }),
  lastLogOutput: getter(function () {
    const count = this.logOutputItems.length;
    if (count === 0) {
      // If no logOutput items are found, we can assume the response is empty
      return '';
    }
    const outputItemText = this.logOutputItems[count - 1].text;
    return outputItemText;
  }),
  logTextItems: collection('[data-test-component="console/log-text"]', {
    text: text(),
  }),
  lastTextOutput: getter(function () {
    const count = this.logTextItems.length;
    if (count === 0) {
      // If no logOutput items are found, we can assume the response is empty
      return '';
    }
    return this.logTextItems.objectAt(count - 1).text;
  }),
  logJSONItems: collection('[data-test-component="console/log-json"]', {
    text: text(),
  }),
  lastJSONOutput: getter(function () {
    const count = this.logJSONItems.length;
    return this.logJSONItems.objectAt(count - 1).text;
  }),
  up: triggerable('keyup', '[data-test-component="console/command-input"] input', {
    eventProperties: { keyCode: keys.UP },
  }),
  down: triggerable('keyup', '[data-test-component="console/command-input"] input', {
    eventProperties: { keyCode: keys.DOWN },
  }),
  enter: triggerable('keyup', '[data-test-component="console/command-input"] input', {
    eventProperties: { keyCode: keys.ENTER },
  }),
  hasInput: isPresent('[data-test-component="console/command-input"] input'),
  runCommands: async function (commands, shouldToggle = true) {
    const toExecute = Array.isArray(commands) ? commands : [commands];
    if (shouldToggle) {
      await this.toggle(); // toggle the console open
    }
    for (const command of toExecute) {
      await this.consoleInput(command);
      await this.enter();
      await settled();
    }
    if (shouldToggle) {
      await this.toggle(); // toggle it closed
    }
  },
};

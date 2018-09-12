import { text, triggerable, clickable, collection, fillable, value, isPresent } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';

import keys from 'vault/lib/keycodes';

export default {
  toggle: clickable('[data-test-console-toggle]'),
  consoleInput: fillable('[data-test-component="console/command-input"] input'),
  consoleInputValue: value('[data-test-component="console/command-input"] input'),
  logOutput: text('[data-test-component="console/output-log"]'),
  logOutputItems: collection('[data-test-component="console/output-log"] > div', {
    text: text(),
  }),
  lastLogOutput: getter(function() {
    let count = this.logOutputItems.length;
    return this.logOutputItems.objectAt(count - 1).text;
  }),
  logTextItems: collection('[data-test-component="console/log-text"]', {
    text: text(),
  }),
  lastTextOutput: getter(function() {
    let count = this.logTextItems.length;
    return this.logTextItems.objectAt(count - 1).text;
  }),
  logJSONItems: collection('[data-test-component="console/log-json"]', {
    text: text(),
  }),
  lastJSONOutput: getter(function() {
    let count = this.logJSONItems.length;
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
  runCommands: async function(commands) {
    let toExecute = Array.isArray(commands) ? commands : [commands];
    for (let command of toExecute) {
      await this.consoleInput(command);
      await this.enter();
    }
  },
};

import { text, triggerable, collection, fillable, value, isPresent } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';

import keys from 'vault/lib/keycodes';

export default {
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
  runCommands(commands) {
    let toExecute = Array.isArray(commands) ? commands : [commands];
    return toExecute.forEach(command => {
      this.consoleInput(command);
      this.enter();
    });
  }
};

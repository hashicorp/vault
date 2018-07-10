import { text, triggerable, fillable, value, isPresent } from 'ember-cli-page-object';
import keys from 'vault/lib/keycodes';

export default {
  consoleInput: fillable('[data-test-component="console/command-input"] input'),
  consoleInputValue: value('[data-test-component="console/command-input"] input'),
  logOutput: text('[data-test-component="console/output-log"]'),
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
};

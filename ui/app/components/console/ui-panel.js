import Ember from 'ember';
import {
  parseCommand,
  shiftCommandIndex,
  extractDataAndFlags,
  logFromResponse,
  logFromError,
  logErrorFromInput,
} from 'vault/lib/console-helpers';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  classNames: 'console-ui-panel-scoller',
  classNameBindings: ['isFullscreen:fullscreen'],
  isFullscreen: false,
  console: inject.service(),
  inputValue: null,
  commandHistory: computed('log.[]', function() {
    return this.get('log').filterBy('type', 'command');
  }),
  log: computed(function() {
    return [];
  }),
  commandIndex: null,

  clearLog() {
    let history = this.get('commandHistory').slice();
    history.setEach('hidden', true);
    let log = this.get('log');
    log.clear();
    log.addObjects(history);
  },

  logAndOutput(command, logContent) {
    let log = this.get('log');
    this.set('inputValue', '');
    log.pushObject({ type: 'command', content: command });
    this.set('commandIndex', null);
    if (logContent) {
      log.pushObject(logContent);
    }
  },

  executeCommand(command, shouldThrow = false) {
    let serviceArgs;
    if (command === 'clear') {
      this.logAndOutput(command);
      this.clearLog();
      return;
    }
    if (command === 'fullscreen') {
      this.toggleProperty('isFullscreen');
      this.logAndOutput(command);
      return;
    }
    // parse to verify it's valid
    try {
      serviceArgs = parseCommand(command, shouldThrow);
    } catch (e) {
      this.logAndOutput(command, { type: 'help' });
      return;
    }
    // we have a invalid command but don't want to throw
    if (serviceArgs === false) {
      return;
    }

    let [method, flagArray, path, dataArray] = serviceArgs;

    if (dataArray || flagArray) {
      var { data, flags } = extractDataAndFlags(dataArray, flagArray);
    }

    let inputError = logErrorFromInput(path, method, flags, dataArray);
    if (inputError) {
      this.logAndOutput(command, inputError);
    }
    let service = this.get('console');
    let serviceFn = service[method];

    serviceFn.call(service, path, data, flags.wrapTTL)
      .then(resp => {
        this.logAndOutput(command, logFromResponse(resp, path, method, flags));
      })
      .catch(error => {
        this.logAndOutput(command, logFromError(error, path, method));
      });
  },

  shiftCommandIndex(keyCode) {
    let [index, newInputValue] = shiftCommandIndex(
      keyCode,
      this.get('commandHistory'),
      this.get('commandIndex')
    );
    this.set('commandIndex', index);
    this.set('inputValue', newInputValue);
  },

  actions: {
    setValue(val) {
      this.set('inputValue', val);
    },
    toggleFullscreen() {
      this.toggleProperty('isFullscreen');
    },
    executeCommand(val) {
      this.executeCommand(val, true);
    },
    shiftCommandIndex(direction) {
      this.shiftCommandIndex(direction);
    },
  },
});

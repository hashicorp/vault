import Ember from 'ember';
import {
  parseCommand,
  extractDataAndFlags,
  logFromResponse,
  logFromError,
  logErrorFromInput,
  executeUICommand,
} from 'vault/lib/console-helpers';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  classNames: 'console-ui-panel-scroller',
  classNameBindings: ['isFullscreen:fullscreen'],
  isFullscreen: false,
  console: inject.service(),
  inputValue: null,
  log: computed.alias('console.log'),

  logAndOutput(command, logContent) {
    this.set('inputValue', '');
    this.get('console').logAndOutput(command, logContent);
  },

  executeCommand(command, shouldThrow = false) {
    let service = this.get('console');
    let serviceArgs;
    
    if(executeUICommand(command, (args) => this.logAndOutput(args), (args) => service.clearLog(args), () => this.toggleProperty('isFullscreen'))){
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
      return;
    }
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
    this.get('console').shiftCommandIndex(keyCode, (val) => {
      this.set('inputValue', val);
    });
  },

  actions: {
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

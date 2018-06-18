import Ember from 'ember';
import { task } from 'ember-concurrency';
import {
  parseCommand,
  extractDataAndFlags,
  logFromResponse,
  logFromError,
  logErrorFromInput,
  executeUICommand,
} from 'vault/lib/console-helpers';

const { inject, computed, getOwner, run } = Ember;

export default Ember.Component.extend({
  classNames: 'console-ui-panel-scroller',
  classNameBindings: ['isFullscreen:fullscreen'],
  isFullscreen: false,
  console: inject.service(),
  router: inject.service(),
  inputValue: null,
  log: computed.alias('console.log'),

  didRender() {
    this._super(...arguments);
    this.scrollToBottom();
  },

  logAndOutput(command, logContent) {
    this.get('console').logAndOutput(command, logContent);
    run.schedule('afterRender', () => this.scrollToBottom());
  },

  isRunning: computed.or('executeCommand.isRunning', 'refreshRoute.isRunning'),

  executeCommand: task(function*(command, shouldThrow = false) {
    this.set('inputValue', '');
    let service = this.get('console');
    let serviceArgs;

    if (
      executeUICommand(
        command,
        args => this.logAndOutput(args),
        args => service.clearLog(args),
        () => this.toggleProperty('isFullscreen'),
        () => this.get('refreshRoute').perform()
      )
    ) {
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
    try {
      let resp = yield service[method].call(service, path, data, flags.wrapTTL);
      this.logAndOutput(command, logFromResponse(resp, path, method, flags));
    } catch (error) {
      this.logAndOutput(command, logFromError(error, path, method));
    }
  }),

  refreshRoute: task(function*() {
    let owner = getOwner(this);
    let routeName = this.get('router.currentRouteName');
    let route = owner.lookup(`route:${routeName}`);

    try {
      yield route.refresh();
      this.logAndOutput(null, { type: 'success', content: 'The current screen has been refreshed!' });
    } catch (error) {
      this.logAndOutput(null, { type: 'error', content: 'The was a problem refreshing the current screen.' });
    }
  }),

  shiftCommandIndex(keyCode) {
    this.get('console').shiftCommandIndex(keyCode, val => {
      this.set('inputValue', val);
    });
  },

  scrollToBottom() {
    this.element.scrollTop = this.element.scrollHeight;
  },

  actions: {
    toggleFullscreen() {
      this.toggleProperty('isFullscreen');
    },
    executeCommand(val) {
      this.get('executeCommand').perform(val, true);
    },
    shiftCommandIndex(direction) {
      this.shiftCommandIndex(direction);
    },
  },
});

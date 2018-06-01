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

const { inject, computed, getOwner } = Ember;

export default Ember.Component.extend({
  classNames: 'console-ui-panel-scroller',
  classNameBindings: ['isFullscreen:fullscreen'],
  isFullscreen: false,
  console: inject.service(),
  router: inject.service(),
  inputValue: null,
  log: computed.alias('console.log'),

  didReceiveAttrs() {
    let val = this.get('inputValue');
    let oldVal = this.get('oldInputValue');
    this.set('valChanged', val !== oldVal);
    this.set('oldInputValue', val);
  },

  didRender() {
    if (this.get('valChanged')) {
      // make sure we're scrolled to the bottom;
      this.scrollToBottom();
    }
  },

  logAndOutput(command, logContent) {
    this.get('console').logAndOutput(command, logContent);
  },

  executeCommand(command, shouldThrow = false) {
    this.set('inputValue', '');
    let service = this.get('console');
    let serviceArgs;

    if (
      executeUICommand(
        command,
        args => this.logAndOutput(args),
        args => service.clearLog(args),
        () => this.toggleProperty('isFullscreen'),
        () => this.refreshRoute()
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
    this.get('runner').perform(
      () => service[method].call(service, path, data, flags.wrapTTL),
      resp => this.logAndOutput(command, logFromResponse(resp, path, method, flags)),
      error => this.logAndOutput(command, logFromError(error, path, method))
    );
  },

  refreshRoute() {
    let owner = getOwner(this);
    let routeName = this.get('router.currentRouteName');
    let route = owner.lookup(`route:${routeName}`);

    this.get('runner').perform(
      () => route.refresh(),
      () => this.logAndOutput(null, { type: 'success', content: 'The current screen has been refreshed!' }),
      () =>
        this.logAndOutput(null, {
          type: 'error',
          content: 'The was a problem refreshing the current screen.',
        })
    );
  },

  runner: task(function*(fn, success, err) {
    yield fn().then(success, err);
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
      this.executeCommand(val, true);
    },
    shiftCommandIndex(direction) {
      this.shiftCommandIndex(direction);
    },
  },
});

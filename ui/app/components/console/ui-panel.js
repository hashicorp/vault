import { inject as service } from '@ember/service';
import { alias, or } from '@ember/object/computed';
import Component from '@ember/component';
import { getOwner } from '@ember/application';
import { run } from '@ember/runloop';
import { task } from 'ember-concurrency';
import ControlGroupError from 'vault/lib/control-group-error';
import {
  parseCommand,
  extractDataAndFlags,
  logFromResponse,
  logFromError,
  logErrorFromInput,
  executeUICommand,
} from 'vault/lib/console-helpers';

export default Component.extend({
  console: service(),
  router: service(),
  controlGroup: service(),
  store: service(),
  'data-test-component': 'console/ui-panel',

  classNames: 'console-ui-panel',
  classNameBindings: ['isFullscreen:fullscreen'],
  isFullscreen: false,
  inputValue: null,
  log: alias('console.log'),

  didRender() {
    this._super(...arguments);
    this.scrollToBottom();
  },

  logAndOutput(command, logContent) {
    this.get('console').logAndOutput(command, logContent);
    run.schedule('afterRender', () => this.scrollToBottom());
  },

  isRunning: or('executeCommand.isRunning', 'refreshRoute.isRunning'),

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
      if (error instanceof ControlGroupError) {
        return this.logAndOutput(command, this.get('controlGroup').logFromError(error));
      }
      this.logAndOutput(command, logFromError(error, path, method));
    }
  }),

  refreshRoute: task(function*() {
    let owner = getOwner(this);
    let routeName = this.get('router.currentRouteName');
    let route = owner.lookup(`route:${routeName}`);

    try {
      this.get('store').clearAllDatasets();
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
    closeConsole() {
      this.set('console.isOpen', false);
    },
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

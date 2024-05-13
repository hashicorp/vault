/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { alias, or } from '@ember/object/computed';
import Component from '@ember/component';
import { getOwner } from '@ember/application';
import { schedule } from '@ember/runloop';
import { camelize } from '@ember/string';
import { task } from 'ember-concurrency';
import { buildWaiter } from '@ember/test-waiters';
import ControlGroupError from 'vault/lib/control-group-error';
import {
  parseCommand,
  logFromResponse,
  logFromError,
  formattedErrorFromInput,
  executeUICommand,
  extractFlagsFromStrings,
  extractDataFromStrings,
} from 'vault/lib/console-helpers';

const waiter = buildWaiter('web-repl');

export default Component.extend({
  console: service(),
  router: service(),
  controlGroup: service(),
  store: service(),
  'data-test-component': 'console/ui-panel',
  attributeBindings: ['data-test-component'],

  classNames: 'console-ui-panel',
  classNameBindings: ['isFullscreen:fullscreen'],
  isFullscreen: false,
  inputValue: null,
  cliLog: alias('console.log'),

  didRender() {
    this._super(...arguments);
    this.scrollToBottom();
  },

  logAndOutput(command, logContent) {
    this.console.logAndOutput(command, logContent);
    schedule('afterRender', () => this.scrollToBottom());
  },

  isRunning: or('executeCommand.isRunning', 'refreshRoute.isRunning'),

  executeCommand: task(function* (command, shouldThrow = false) {
    this.set('inputValue', '');
    const waiterToken = waiter.beginAsync();
    const service = this.console;
    let serviceArgs;

    if (
      executeUICommand(command, (args) => this.logAndOutput(args), {
        api: () => this.routeToExplore.perform(command),
        clearall: () => service.clearLog(true),
        clear: () => service.clearLog(),
        fullscreen: () => this.toggleProperty('isFullscreen'),
        refresh: () => this.refreshRoute.perform(),
      })
    ) {
      waiter.endAsync(waiterToken);
      return;
    }

    // parse to verify it's valid
    try {
      serviceArgs = parseCommand(command);
    } catch (e) {
      if (shouldThrow) {
        this.logAndOutput(command, { type: 'help' });
      }
      waiter.endAsync(waiterToken);
      return;
    }

    const { method, flagArray, path, dataArray } = serviceArgs;
    const flags = extractFlagsFromStrings(flagArray, method);
    const data = extractDataFromStrings(dataArray);

    const inputError = formattedErrorFromInput(path, method, flags, dataArray);
    if (inputError) {
      this.logAndOutput(command, inputError);
      waiter.endAsync(waiterToken);
      return;
    }
    try {
      const resp = yield service[camelize(method)].call(service, path, data, flags);
      this.logAndOutput(command, logFromResponse(resp, path, method, flags));
    } catch (error) {
      if (error instanceof ControlGroupError) {
        waiter.endAsync(waiterToken);
        return this.logAndOutput(command, this.controlGroup.logFromError(error));
      }
      this.logAndOutput(command, logFromError(error, path, method));
    }
    waiter.endAsync(waiterToken);
  }),

  refreshRoute: task(function* () {
    const owner = getOwner(this);
    const currentRoute = owner.lookup(`router:main`).currentRouteName;

    try {
      this.store.clearAllDatasets();
      yield this.router.transitionTo(currentRoute);
      this.logAndOutput(null, { type: 'success', content: 'The current screen has been refreshed!' });
    } catch (error) {
      this.logAndOutput(null, { type: 'error', content: 'The was a problem refreshing the current screen.' });
    }
  }),

  routeToExplore: task(function* (command) {
    const filter = command.replace('api', '').trim();
    let content =
      'Welcome to the Vault API explorer! \nYou can search for endpoints, see what parameters they accept, and even execute requests with your current token.';
    if (filter) {
      content = `Welcome to the Vault API explorer! \nWe've filtered the list of endpoints for '${filter}'.`;
    }
    try {
      yield this.router.transitionTo('vault.cluster.tools.open-api-explorer', {
        queryParams: { filter },
      });
      this.logAndOutput(null, {
        type: 'success',
        content,
      });
    } catch (error) {
      if (error.message === 'TransitionAborted') {
        this.logAndOutput(null, {
          type: 'success',
          content,
        });
      } else {
        this.logAndOutput(null, {
          type: 'error',
          content: 'There was a problem navigating to the api explorer.',
        });
      }
    }
  }),

  shiftCommandIndex(keyCode) {
    this.console.shiftCommandIndex(keyCode, (val) => {
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
      this.executeCommand.perform(val, true);
    },
    shiftCommandIndex(direction) {
      this.shiftCommandIndex(direction);
    },
  },
});

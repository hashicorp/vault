/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import { alias, or } from '@ember/object/computed';
import { schedule } from '@ember/runloop';
import { camelize } from '@ember/string';
import { task } from 'ember-concurrency';
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

/**
 * @module UiPanel
 * UiPanel is the console window provided so users can run a limited set of CLI commands on the GUI.
 *
 * @example
 * ```js
 <UiPanel todo/>
 * ```
 * @param {string} [mode=null] - todo
 */

export default class UiPanel extends Component {
  @service console;
  @service router;
  @service controlGroup;
  @service store;

  @tracked isFullScreen = false;
  @tracked inputValue = null;
  @tracked element = null;

  @alias('console.log') cliLog;

  constructor() {
    super(...arguments);
  }

  scrollToBottom(element) {
    // We do not have access to element after entering a command. Save the original element var from the executeCommand task to use for this situation.
    const container = !element ? this.element : element;
    container.scrollTop = container.scrollHeight;
  }

  logAndOutput(command, logContent) {
    this.console.logAndOutput(command, logContent);
    schedule('afterRender', () => this.scrollToBottom());
  }

  @or('executeCommand.isRunning', 'refreshRoute.isRunning') isRunning;

  @task
  *executeCommand(element, shouldThrow = true) {
    this.element = element;
    const command = element.value;
    this.inputValue = '';
    const service = this.console;
    let serviceArgs;

    if (
      executeUICommand(command, (args) => this.logAndOutput(args), {
        api: () => this.routeToExplore.perform(command),
        clearall: () => service.clearLog(true),
        clear: () => service.clearLog(),
        fullscreen: () => (this.isFullscreen = !this.isFullScreen),
        refresh: () => this.refreshRoute.perform(),
      })
    ) {
      return;
    }

    // parse to verify it's valid
    try {
      serviceArgs = parseCommand(command);
    } catch (e) {
      if (shouldThrow) {
        this.logAndOutput(command, { type: 'help' });
      }
      return;
    }

    const { method, flagArray, path, dataArray } = serviceArgs;
    const flags = extractFlagsFromStrings(flagArray, method);
    const data = extractDataFromStrings(dataArray);

    const inputError = formattedErrorFromInput(path, method, flags, dataArray);
    if (inputError) {
      this.logAndOutput(command, inputError);
      return;
    }
    try {
      const resp = yield service[camelize(method)].call(service, path, data, flags);
      this.logAndOutput(command, logFromResponse(resp, path, method, flags));
    } catch (error) {
      if (error instanceof ControlGroupError) {
        return this.logAndOutput(command, this.controlGroup.logFromError(error));
      }
      this.logAndOutput(command, logFromError(error, path, method));
    }
  }

  @task
  *refreshRoute() {
    const currentRoute = this.router.currentRouteName;
    try {
      this.store.clearAllDatasets();
      yield this.router.transitionTo(currentRoute);
      this.logAndOutput(null, { type: 'success', content: 'The current screen has been refreshed!' });
    } catch (error) {
      this.logAndOutput(null, { type: 'error', content: 'The was a problem refreshing the current screen.' });
    }
  }

  @task
  *routeToExplore(command) {
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
  }

  @action
  closeConsole() {
    this.console.isOpen = false;
  }
  @action
  toggleFullscreen() {
    this.isFullScreen = !this.isFullScreen;
  }
  @action
  shiftCommandIndex(keyCode) {
    this.console.shiftCommandIndex(keyCode, (val) => {
      this.inputValue = val;
    });
  }
}

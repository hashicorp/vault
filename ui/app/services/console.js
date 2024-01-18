/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import { filterBy } from '@ember/object/computed';

import Service, { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { shiftCommandIndex } from 'vault/lib/console-helpers';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { sanitizePath, ensureTrailingSlash } from 'core/utils/sanitize-path';

const VERBS = {
  read: 'GET',
  list: 'GET',
  write: 'POST',
  delete: 'DELETE',
};

export default class ConsoleService extends Service {
  @tracked isOpen = false;
  @tracked log = [];
  @tracked commandIndex = null;

  @service store;

  @filterBy('log', 'type', 'command') commandHistory;

  /* eslint ember/no-computed-properties-in-native-classes: 'warn' */
  shiftCommandIndex(keyCode, setCommandFn = () => {}) {
    const [newIndex, newCommand] = shiftCommandIndex(keyCode, this.commandHistory, this.commandIndex);
    if (newCommand !== undefined && newIndex !== undefined) {
      this.commandIndex = newIndex;
      setCommandFn(newCommand);
    }
  }

  clearLog(clearAll = false) {
    let history;
    if (!clearAll) {
      history = this.commandHistory.slice();
      history.setEach('hidden', true);
    }
    this.log.clear();
    if (history) {
      this.log.addObjects(history);
    }
  }

  logAndOutput(command, logContent) {
    if (command) {
      this.log.pushObject({ type: 'command', content: command });
      this.commandIndex = null;
    }
    if (logContent) {
      this.log.pushObject(logContent);
    }
  }

  ajax(operation, path, options = {}) {
    const verb = VERBS[operation];
    const adapter = this.store.adapterFor('console');
    const url = adapter.buildURL(encodePath(path));
    const { data, wrapTTL } = options;
    return adapter.ajax(url, verb, {
      data,
      wrapTTL,
    });
  }

  kvGet(path, data, flags = {}) {
    const { wrapTTL, metadata } = flags;
    // Split on first / to find backend and secret path
    const pathSegment = metadata ? 'metadata' : 'data';
    const [backend, secretPath] = path.split(/\/(.+)?/);
    const kvPath = `${backend}/${pathSegment}/${secretPath}`;
    return this.ajax('read', sanitizePath(kvPath), { wrapTTL });
  }

  async read(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    return await this.ajax('read', sanitizePath(path), { wrapTTL });
  }

  async write(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    return await this.ajax('write', sanitizePath(path), { data, wrapTTL });
  }

  async delete(path) {
    return await this.ajax('delete', sanitizePath(path));
  }

  async list(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    const listPath = ensureTrailingSlash(sanitizePath(path));
    return await this.ajax('list', listPath, {
      data: {
        list: true,
      },
      wrapTTL,
    });
  }
}

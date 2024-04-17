/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import Service from '@ember/service';
import { getOwner } from '@ember/application';
import { shiftCommandIndex } from 'vault/lib/console-helpers';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { sanitizePath, ensureTrailingSlash } from 'core/utils/sanitize-path';
import { tracked } from '@glimmer/tracking';
import { addManyToArray } from 'vault/helpers/add-to-array';

const VERBS = {
  read: 'GET',
  list: 'GET',
  write: 'POST',
  delete: 'DELETE',
};

export default class ConsoleService extends Service {
  @tracked isOpen = false;
  @tracked commandIndex = null;
  @tracked log = [];

  get commandHistory() {
    return this.log.filter((log) => log.type === 'command');
  }

  // Not a getter so it can be stubbed in tests
  adapter() {
    return getOwner(this).lookup('adapter:console');
  }

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
    this.log = [];
    if (history) {
      this.log = addManyToArray(this.log, history);
    }
  }

  logAndOutput(command, logContent) {
    const log = this.log.slice();
    if (command) {
      log.push({ type: 'command', content: command });
      this.commandIndex = null;
    }
    if (logContent) {
      log.push(logContent);
    }
    this.log = log;
  }

  ajax(operation, path, options = {}) {
    const verb = VERBS[operation];
    const adapter = this.adapter();
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

  read(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    return this.ajax('read', sanitizePath(path), { wrapTTL });
  }

  write(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    return this.ajax('write', sanitizePath(path), { data, wrapTTL });
  }

  delete(path) {
    return this.ajax('delete', sanitizePath(path));
  }

  list(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    const listPath = ensureTrailingSlash(sanitizePath(path));
    return this.ajax('list', listPath, {
      data: {
        list: true,
      },
      wrapTTL,
    });
  }
}

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import { filterBy } from '@ember/object/computed';

import Service from '@ember/service';

import { getOwner } from '@ember/application';
import { computed } from '@ember/object';
import { shiftCommandIndex } from 'vault/lib/console-helpers';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { sanitizePath, ensureTrailingSlash } from 'core/utils/sanitize-path';

const VERBS = {
  read: 'GET',
  list: 'GET',
  write: 'POST',
  delete: 'DELETE',
};

export default Service.extend({
  isOpen: false,

  adapter() {
    return getOwner(this).lookup('adapter:console');
  },
  commandHistory: filterBy('log', 'type', 'command'),
  log: computed(function () {
    return [];
  }),
  commandIndex: null,

  shiftCommandIndex(keyCode, setCommandFn = () => {}) {
    const [newIndex, newCommand] = shiftCommandIndex(keyCode, this.commandHistory, this.commandIndex);
    if (newCommand !== undefined && newIndex !== undefined) {
      this.set('commandIndex', newIndex);
      setCommandFn(newCommand);
    }
  },

  clearLog(clearAll = false) {
    const log = this.log;
    let history;
    if (!clearAll) {
      history = this.commandHistory.slice();
      history.setEach('hidden', true);
    }
    log.clear();
    if (history) {
      log.addObjects(history);
    }
  },

  logAndOutput(command, logContent) {
    const log = this.log;
    if (command) {
      log.pushObject({ type: 'command', content: command });
      this.set('commandIndex', null);
    }
    if (logContent) {
      log.pushObject(logContent);
    }
  },

  ajax(operation, path, options = {}) {
    const verb = VERBS[operation];
    const adapter = this.adapter();
    const url = adapter.buildURL(encodePath(path));
    const { data, wrapTTL } = options;
    return adapter.ajax(url, verb, {
      data,
      wrapTTL,
    });
  },

  kvGet(path, data, flags = {}) {
    const { wrapTTL, metadata } = flags;
    // Split on first / to find backend and secret path
    const pathSegment = metadata ? 'metadata' : 'data';
    const [backend, secretPath] = path.split(/\/(.+)?/);
    const kvPath = `${backend}/${pathSegment}/${secretPath}`;
    return this.ajax('read', sanitizePath(kvPath), { wrapTTL });
  },

  read(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    return this.ajax('read', sanitizePath(path), { wrapTTL });
  },

  write(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    return this.ajax('write', sanitizePath(path), { data, wrapTTL });
  },

  delete(path) {
    return this.ajax('delete', sanitizePath(path));
  },

  list(path, data, flags) {
    const wrapTTL = flags?.wrapTTL;
    const listPath = ensureTrailingSlash(sanitizePath(path));
    return this.ajax('list', listPath, {
      data: {
        list: true,
      },
      wrapTTL,
    });
  },
});

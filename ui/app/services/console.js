// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import Service from '@ember/service';

import { getOwner } from '@ember/application';
import { computed } from '@ember/object';
import { shiftCommandIndex } from 'vault/lib/console-helpers';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export function sanitizePath(path) {
  //remove whitespace + remove trailing and leading slashes
  return path.trim().replace(/^\/+|\/+$/g, '');
}
export function ensureTrailingSlash(path) {
  return path.replace(/(\w+[^/]$)/g, '$1/');
}

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
  commandHistory: computed('log.[]', function() {
    return this.get('log').filterBy('type', 'command');
  }),
  log: computed(function() {
    return [];
  }),
  commandIndex: null,

  shiftCommandIndex(keyCode, setCommandFn = () => {}) {
    let [newIndex, newCommand] = shiftCommandIndex(
      keyCode,
      this.get('commandHistory'),
      this.get('commandIndex')
    );
    if (newCommand !== undefined && newIndex !== undefined) {
      this.set('commandIndex', newIndex);
      setCommandFn(newCommand);
    }
  },

  clearLog(clearAll = false) {
    let log = this.get('log');
    let history;
    if (!clearAll) {
      history = this.get('commandHistory').slice();
      history.setEach('hidden', true);
    }
    log.clear();
    if (history) {
      log.addObjects(history);
    }
  },

  logAndOutput(command, logContent) {
    let log = this.get('log');
    if (command) {
      log.pushObject({ type: 'command', content: command });
      this.set('commandIndex', null);
    }
    if (logContent) {
      log.pushObject(logContent);
    }
  },

  ajax(operation, path, options = {}) {
    let verb = VERBS[operation];
    let adapter = this.adapter();
    let url = adapter.buildURL(encodePath(path));
    let { data, wrapTTL } = options;
    return adapter.ajax(url, verb, {
      data,
      wrapTTL,
    });
  },

  read(path, data, wrapTTL) {
    return this.ajax('read', sanitizePath(path), { wrapTTL });
  },

  write(path, data, wrapTTL) {
    return this.ajax('write', sanitizePath(path), { data, wrapTTL });
  },

  delete(path) {
    return this.ajax('delete', sanitizePath(path));
  },

  list(path, data, wrapTTL) {
    let listPath = ensureTrailingSlash(sanitizePath(path));
    return this.ajax('list', listPath, {
      data: {
        list: true,
      },
      wrapTTL,
    });
  },
});

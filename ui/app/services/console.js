// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import Ember from 'ember';

const { Service, getOwner } = Ember;

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
  adapter() {
    return getOwner(this).lookup('adapter:console');
  },

  ajax(operation, path, options) {
    let verb = VERBS[operation];
    let url = `${this.urlPrefix()}/${path}`;
    let { data, wrapTTL } = options;
    return this.adapter().ajax(url, verb, {
      data,
      wrapTTL,
    });
  },

  read(path) {
    return this.ajax('read', sanitizePath(path));
  },

  write(path, data) {
    return this.ajax('write', sanitizePath(path), { data });
  },

  delete(path) {
    return this.ajax('delete', sanitizePath(path));
  },

  list(path) {
    let listPath = ensureTrailingSlash(sanitizePath(path));
    return this.ajax('list', listPath, {
      data: {
        list: true,
      },
    });
  },
});

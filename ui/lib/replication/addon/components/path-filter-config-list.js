import Component from '@ember/component';
import { set, get, computed } from '@ember/object';
import { inject as service } from '@ember/service';
import { readOnly } from '@ember/object/computed';
import { task } from 'ember-concurrency';

export default Component.extend({
  namespace: service(),
  store: service(),
  config: null,
  possiblePaths: null,
  namespaces: readOnly('namespace.accessibleNamespaces'),
  autoCompleteOptions: null,

  fetchMounts() {
    return hash({
      mounts: this.store.findAll('secret-engine'),
      auth: this.store.findAll('auth-method'),
    }).then(({ mounts, auth }) => {
      return resolve(mounts.toArray().concat(auth.toArray()));
    });
  },

  init() {
    this._super(...arguments);
    this.setAutoCompleteOptions.perform();
  },

  fetchMountsForNamespace: task(function*(ns) {
    let adapter = this.store.adapterFor('application');
    let mounts = ns
      ? yield adapter.ajax('/v1/sys/internal/ui/mounts', 'GET', { namespace: ns })
      : yield adapter.ajax('/v1/sys/internal/ui/mounts', 'GET');
    [['secret', 'secret-engine'], ['auth', 'auth-method']].forEach(([key, modelType]) => {
      for (let [id, info] of Object.entries(mounts.data[key])) {
        let longId = ns ? `${ns}/${id}` : id;
        info.path = longId;

        // don't add singleton mounts
        if (!this.singletonMountTypes.includes(info.type)) {
          this.store.push({
            data: {
              type: modelType,
              id: longId.replace(/\/$/, ''),
              attributes: {
                ...info,
              },
            },
          });
        }
      }
    });
    return {
      secret: this.store.peekAll('secret-engine'),
      auth: this.store.peekAll('auth-method'),
    };
  }),

  setAutoCompleteOptions: task(function*() {
    let { namespaces, singletonMountTypes: singletons, autoCompleteOptions } = this;
    // fetch auth and secret methods from sys/internal/ui/mounts
    // for any namespaces that are already autocompleted
    let { secret, auth } = yield this.fetchMountsForNamespace.perform();
    let options = namespaces.concat(secret.toArray(), auth.toArray()).map(val => {
      // namespaces are just strings, mounts are objects
      let path = val.id ? val.path : val;
      return {
        id: path,
        name: path,
        searchText: path,
      };
    });

    this.set('autoCompleteOptions', options);
  }).keepLatest(),

  // singleton mounts are not eligible for per-mount-filtering
  singletonMountTypes: computed(function() {
    return ['cubbyhole', 'system', 'token', 'identity', 'ns_system', 'ns_identity'];
  }),

  willDestroyElement() {
    this.store.unloadAll('secret-engine');
    this.store.unloadAll('auth-method');
    this._super(...arguments);
  },

  actions: {
    pathsChanged(paths) {
      console.table(paths);
      this.set('config.paths', paths);
    },

    inputChanged(val, powerSelect) {
      console.log(val, powerSelect);
    },

    addOrRemovePath(path, e) {
      let config = get(this, 'config') || [];
      let paths = get(config, 'paths').slice();

      if (e.target.checked) {
        paths.addObject(path);
      } else {
        paths.removeObject(path);
      }

      set(config, 'paths', paths);
    },
  },
});

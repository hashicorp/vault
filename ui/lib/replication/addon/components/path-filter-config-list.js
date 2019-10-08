import Component from '@ember/component';
import { set, get, computed } from '@ember/object';
import { inject as service } from '@ember/service';
import { readOnly } from '@ember/object/computed';
import { task, timeout } from 'ember-concurrency';

export default Component.extend({
  namespace: service(),
  store: service(),
  config: null,
  possiblePaths: null,
  namespaces: readOnly('namespace.accessibleNamespaces'),
  autoCompleteOptions: null,
  namespacesFetched: null,

  init() {
    this._super(...arguments);
    this.setAutoCompleteOptions.perform();
  },

  fetchMountsForNamespace: task(function*(ns) {
    let adapter = this.store.adapterFor('application');
    let secret = [];
    let auth = [];
    let mounts = ns
      ? yield adapter.ajax('/v1/sys/internal/ui/mounts', 'GET', { namespace: ns })
      : yield adapter.ajax('/v1/sys/internal/ui/mounts', 'GET');

    ['secret', 'auth'].forEach(key => {
      for (let [id, info] of Object.entries(mounts.data[key])) {
        let longId;
        if (key === 'auth') {
          longId = ns ? `${ns}/auth/${id}` : `auth/${id}`;
        } else {
          longId = ns ? `${ns}/${id}` : id;
        }
        info.path = longId;
        longId = longId.replace(/\/$/, '');

        // don't add singleton mounts
        if (!this.singletonMountTypes.includes(info.type)) {
          (key === 'secret' ? secret : auth).push({
            id: longId,
            name: longId,
            searchText: `${longId} ${info.type} ${info.accessor}`,
          });
        }
      }
    });
    return {
      secret,
      auth,
    };
  }),

  setAutoCompleteOptions: task(function*(term, removeOptions) {
    let { namespaces, autoCompleteOptions } = this;
    if (removeOptions) {
      return;
    }
    // fetch auth and secret methods from sys/internal/ui/mounts for the given namespace
    let { secret: secretList, auth: authList } = yield this.fetchMountsForNamespace.perform(namespaceToFetch);
    if (autoCompleteOptions) {
      var currentSecrets = autoCompleteOptions.findBy('groupName', 'Secret Engines');
      var currentAuths = autoCompleteOptions.findBy('groupName', 'Auth Methods');
    }
    let formattedNamespaces = namespaces.map(val => {
      return {
        id: val,
        name: val,
        searchText: val,
      };
    });
    let options = [];
    if (formattedNamespaces.length) {
      options.push({ groupName: 'Namespaces', options: formattedNamespaces });
    }
    if (secretList.length) {
      let secretOptions = currentSecrets ? [...currentSecrets.options, ...secretList] : secretList;
      options.push({ groupName: 'Secret Engines', options: secretOptions.uniqBy('id') });
    }
    if (authList.length) {
      let authOptions = currentAuths ? [...currentAuths.options, ...authList] : authList;
      options.push({ groupName: 'Auth Methods', options: authOptions.uniqBy('id') });
    }
    this.set('autoCompleteOptions', options);
    return options;
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
      // need to check the current paths and see if there's to remove then call
      this.set('config.paths', paths);
    },

    inputChanged(val, powerSelect) {
      console.log(val, powerSelect);

      if (this.namespaces.includes(val)) {
        this.setAutoCompleteOptions.perform(val);
      }
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

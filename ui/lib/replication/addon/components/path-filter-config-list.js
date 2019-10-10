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
  lastOptions: null,
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

  filterOptions(list, term) {
    return list
      .map(({ groupName, options }) => {
        let trimmedOptions = term ? options.filter(op => op.searchText.includes(term)) : options;
        return trimmedOptions.length ? { groupName, options: trimmedOptions } : null;
      })
      .compact();
  },

  setAutoCompleteOptions: task(function*(term, powerSelect) {
    let { namespaces, lastOptions } = this;
    let namespaceToFetch = namespaces.find(ns => ns === term);
    let secretList = [];
    let authList = [];
    let options = [];
    if (!term || (term && namespaceToFetch)) {
      // fetch auth and secret methods from sys/internal/ui/mounts for the given namespace
      let result = yield this.fetchMountsForNamespace.perform(namespaceToFetch);
      secretList = result.secret;
      authList = result.auth;
    }
    var currentSecrets = lastOptions && lastOptions.findBy('groupName', 'Secret Engines');
    var currentAuths = lastOptions && lastOptions.findBy('groupName', 'Auth Methods');
    let formattedNamespaces = namespaces.map(val => {
      return {
        id: val,
        name: val,
        searchText: val,
      };
    });

    options.push({ groupName: 'Namespaces', options: formattedNamespaces });
    let secretOptions = currentSecrets ? [...currentSecrets.options, ...secretList] : secretList;

    options.push({ groupName: 'Secret Engines', options: secretOptions.uniqBy('id') });
    let authOptions = currentAuths ? [...currentAuths.options, ...authList] : authList;

    options.push({ groupName: 'Auth Methods', options: authOptions.uniqBy('id') });
    if (!term) {
      this.set('autoCompleteOptions', this.filterOptions(options));
      return;
    }
    let filtered = this.filterOptions(options, term);
    this.set('lastOptions', filtered);
    return filtered;
  }),

  // singleton mounts are not eligible for per-mount-filtering
  singletonMountTypes: computed(function() {
    return ['cubbyhole', 'system', 'token', 'identity', 'ns_system', 'ns_identity'];
  }),

  willDestroyElement() {
    this._super(...arguments);
  },

  actions: {
    pathsChanged(paths) {
      set(this.config, 'paths', paths);
    },
  },
});

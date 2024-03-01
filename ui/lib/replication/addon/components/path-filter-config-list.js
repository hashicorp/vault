/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { set, computed } from '@ember/object';
import { service } from '@ember/service';
import { readOnly } from '@ember/object/computed';
import { task, timeout } from 'ember-concurrency';

export default Component.extend({
  'data-test-component': 'path-filter-config',
  attributeBindings: ['data-test-component'],
  namespace: service(),
  store: service(),
  config: null,
  namespaces: readOnly('namespace.accessibleNamespaces'),
  lastOptions: null,
  autoCompleteOptions: null,
  namespacesFetched: null,
  startedWithMode: false,

  init() {
    this._super(...arguments);
    this.setAutoCompleteOptions.perform();
    if (this.config.mode) {
      this.set('startedWithMode', true);
    }
  },

  fetchMountsForNamespace: task(function* (ns) {
    const adapter = this.store.adapterFor('application');
    const secret = [];
    const auth = [];
    const mounts = ns
      ? yield adapter.ajax('/v1/sys/internal/ui/mounts', 'GET', { namespace: ns })
      : yield adapter.ajax('/v1/sys/internal/ui/mounts', 'GET');

    ['secret', 'auth'].forEach((key) => {
      for (const [id, info] of Object.entries(mounts.data[key])) {
        let longId;
        if (key === 'auth') {
          longId = ns ? `${ns}/auth/${id}` : `auth/${id}`;
        } else {
          longId = ns ? `${ns}/${id}` : id;
        }
        info.path = longId;

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
    const paths = this.config.paths;
    return list
      .map(({ groupName, options }) => {
        const trimmedOptions = options.filter((op) => {
          if (term) {
            return op.searchText.includes(term) && !paths.includes(op.id);
          }
          return !paths.includes(op.id);
        });
        return trimmedOptions.length ? { groupName, options: trimmedOptions } : null;
      })
      .compact();
  },

  setAutoCompleteOptions: task(function* (term) {
    const { namespaces, lastOptions } = this;
    const namespaceToFetch = namespaces.find((ns) => ns === term);
    let secretList = [];
    let authList = [];
    const options = [];
    if (term) {
      yield timeout(200);
    }
    if (!term || (term && namespaceToFetch)) {
      // fetch auth and secret methods from sys/internal/ui/mounts for the given namespace
      const result = yield this.fetchMountsForNamespace.perform(namespaceToFetch);
      secretList = result.secret;
      authList = result.auth;
    }
    var currentSecrets = lastOptions && lastOptions.find((opt) => opt.groupName === 'Secret Engines');
    var currentAuths = lastOptions && lastOptions.find((opt) => opt.groupName === 'Auth Methods');
    const formattedNamespaces = namespaces.map((val) => {
      return {
        id: val,
        name: val,
        searchText: val,
      };
    });

    options.push({ groupName: 'Namespaces', options: formattedNamespaces });
    const secretOptions = currentSecrets ? [...currentSecrets.options, ...secretList] : secretList;

    options.push({ groupName: 'Secret Engines', options: secretOptions.uniqBy('id') });
    const authOptions = currentAuths ? [...currentAuths.options, ...authList] : authList;

    options.push({ groupName: 'Auth Methods', options: authOptions.uniqBy('id') });
    const filtered = term ? this.filterOptions(options, term) : this.filterOptions(options);
    if (!term) {
      this.set('autoCompleteOptions', filtered);
    }
    this.set('lastOptions', filtered);
    return filtered;
  }),

  // singleton mounts are not eligible for per-mount-filtering
  singletonMountTypes: computed(function () {
    return ['cubbyhole', 'system', 'token', 'identity', 'ns_system', 'ns_identity', 'ns_token'];
  }),

  actions: {
    async pathsChanged(paths) {
      // set paths on the model
      set(this.config, 'paths', paths);

      // if dropdown is empty or has options without groupName, re-fetch
      if (
        this.autoCompleteOptions.length === 0 ||
        this.autoCompleteOptions.any((option) => !Object.keys(option).includes('groupName'))
      ) {
        await this.setAutoCompleteOptions.perform();
      }
      if (paths.length) {
        // remove the selected item from the default list of options
        const filtered = this.filterOptions(this.autoCompleteOptions);
        this.set('autoCompleteOptions', filtered);
      } else {
        // if there's no paths, we need to re-fetch like on init
        this.setAutoCompleteOptions.perform();
      }
    },
  },
});

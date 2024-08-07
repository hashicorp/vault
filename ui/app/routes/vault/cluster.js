/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { computed } from '@ember/object';
import { reject } from 'rsvp';
import Route from '@ember/routing/route';
import { task, timeout } from 'ember-concurrency';
import Ember from 'ember';
import getStorage from '../../lib/token-storage';
import localStorage from 'vault/lib/local-storage';
import ClusterRoute from 'vault/mixins/cluster-route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';
import { assert } from '@ember/debug';

const POLL_INTERVAL_MS = 10000;

export const getManagedNamespace = (nsParam, root) => {
  if (!nsParam || nsParam.replaceAll('/', '') === root) return root;
  // Check if param starts with root and /
  if (nsParam.startsWith(`${root}/`)) {
    return nsParam;
  }
  // Otherwise prepend the given param with the root
  return `${root}/${nsParam}`;
};

export default Route.extend(ModelBoundaryRoute, ClusterRoute, {
  auth: service(),
  currentCluster: service(),
  customMessages: service(),
  flagsService: service('flags'),
  namespaceService: service('namespace'),
  permissions: service(),
  router: service(),
  store: service(),
  version: service(),
  modelTypes: computed(function () {
    return ['node', 'secret', 'secret-engine'];
  }),

  queryParams: {
    namespaceQueryParam: {
      refreshModel: true,
    },
  },

  getClusterId(params) {
    const { cluster_name } = params;
    const records = this.store.peekAll('cluster');
    const cluster = records.find((record) => record.name === cluster_name);
    return cluster?.id ?? null;
  },

  async beforeModel() {
    const params = this.paramsFor(this.routeName);
    let namespace = params.namespaceQueryParam;
    const currentTokenName = this.auth.currentTokenName;
    const managedRoot = this.flagsService.hvdManagedNamespaceRoot;
    assert(
      'Cannot use VAULT_CLOUD_ADMIN_NAMESPACE flag with non-enterprise Vault version',
      !(managedRoot && this.version.isCommunity)
    );

    // activatedFlags are called this high in routing to return a response used to show/hide Secrets sync on sidebar nav.
    await this.flagsService.fetchActivatedFlags();

    if (!namespace && currentTokenName && !Ember.testing) {
      // if no namespace queryParam and user authenticated,
      // use user's root namespace to redirect to properly param'd url
      const storage = getStorage().getItem(currentTokenName);
      namespace = storage?.userRootNamespace;
      // only redirect if something other than nothing
      if (namespace) {
        this.router.transitionTo({ queryParams: { namespace } });
      }
    } else if (managedRoot !== null) {
      const managed = getManagedNamespace(namespace, managedRoot);
      if (managed !== namespace) {
        this.router.transitionTo({ queryParams: { namespace: managed } });
      }
    }
    this.namespaceService.setNamespace(namespace);
    const id = this.getClusterId(params);
    if (id) {
      this.auth.setCluster(id);
      if (this.auth.currentToken) {
        this.version.fetchVersion();
        await this.permissions.getPaths.perform();
      }
      return this.version.fetchFeatures();
    } else {
      return reject({ httpStatus: 404, message: 'not found', path: params.cluster_name });
    }
  },

  model(params) {
    // if a user's browser settings block localStorage they will be unable to use Vault. The method will throw the error and the rest of the application will not load.
    localStorage.isLocalStorageSupported();

    const id = this.getClusterId(params);
    return this.store.findRecord('cluster', id);
  },

  poll: task(function* () {
    while (true) {
      // when testing, the polling loop causes promises to never settle so acceptance tests hang
      // to get around that, we just disable the poll in tests
      if (Ember.testing) {
        return;
      }
      yield timeout(POLL_INTERVAL_MS);
      try {
        /* eslint-disable-next-line ember/no-controller-access-in-routes */
        yield this.controller.model.reload();
        yield this.transitionToTargetRoute();
      } catch (e) {
        // we want to keep polling here
      }
    }
  })
    .cancelOn('deactivate')
    .keepLatest(),

  afterModel(model, transition) {
    this._super(...arguments);
    this.currentCluster.setCluster(model);
    if (model.needsInit && this.auth.currentToken) {
      // clear token to prevent infinite load state
      this.auth.deleteCurrentToken();
    }

    // Check that namespaces is enabled and if not,
    // clear the namespace by transition to this route w/o it
    if (this.namespaceService.path && !this.version.hasNamespaces) {
      return this.router.transitionTo(this.routeName, { queryParams: { namespace: '' } });
    }
    return this.transitionToTargetRoute(transition);
  },

  setupController() {
    this._super(...arguments);
    this.poll.perform();
  },

  actions: {
    error(e) {
      if (e.httpStatus === 503 && e.errors[0] === 'Vault is sealed') {
        this.refresh();
      }
      return true;
    },
  },
});

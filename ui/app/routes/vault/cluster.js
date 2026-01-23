/**
 * Copyright IBM Corp. 2016, 2025
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

import { v4 as uuidv4 } from 'uuid';

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
  api: service(),
  analytics: service(),
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
      // In test mode, polling causes acceptance tests to hang due to never-settling promises.
      // To avoid this, polling is disabled during tests.
      // If your test depends on cluster status changes (e.g., replication mode),
      // manually trigger polling using pollCluster from 'vault/tests/helpers/poll-cluster'.
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

  // Note: do not make this afterModel hook async, it will break the DR secondary flow.
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
    // Skip analytics initialization if the cluster is a DR secondary:
    // 1. There is little value in collecting analytics in this state.
    // 2. The analytics service requires resolving async setup (e.g. await),
    //   which delays the afterModel hook resolution and breaks the DR secondary flow.
    if (model.dr?.isSecondary) {
      return this.transitionToTargetRoute(transition);
    }

    this.addAnalyticsService(model);

    return this.transitionToTargetRoute(transition);
  },

  async addAnalyticsService(model) {
    // identify user for analytics service
    if (this.analytics.activated) {
      let licenseId = '';

      try {
        const licenseStatus = await this.api.sys.systemReadLicenseStatus();
        licenseId = licenseStatus?.data?.autoloaded?.licenseId;
      } catch (e) {
        // license is not retrievable
        licenseId = '';
      }

      try {
        const entity_id = this.auth.authData?.entityId;
        const entity = entity_id ? entity_id : `root_${uuidv4()}`;

        this.analytics.identifyUser(entity, {
          licenseId: licenseId,
          licenseState: model.license?.state || 'community',
          version: model.version.version,
          storageType: model.storageType,
          replicationMode: model.replicationMode,
          isEnterprise: Boolean(model.license),
        });
      } catch (e) {
        console.error('unable to start analytics', e);
      }
    }
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

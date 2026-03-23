/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { assert } from '@ember/debug';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { filterEnginesByMountCategory, isAddonEngine } from 'core/utils/all-engines-metadata';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';
import { hash } from 'rsvp';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { getEnginePathParam } from 'vault/utils/backend-route-helpers';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';
import { getModelTypeForEngine } from 'vault/utils/model-helpers/secret-engine-helpers';
import { normalizePath } from 'vault/utils/path-encoding-helpers';

const SUPPORTED_BACKENDS = supportedSecretBackends();

function getValidPage(pageParam) {
  if (typeof pageParam === 'number') {
    return pageParam;
  }
  if (typeof pageParam === 'string') {
    try {
      return parseInt(pageParam, 10) || 1;
    } catch (e) {
      return 1;
    }
  }
  return 1;
}

export default Route.extend({
  pagination: service(),
  store: service(),
  templateName: 'vault/cluster/secrets/backend/list',
  pathHelp: service('path-help'),
  router: service(),

  queryParams: {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
    tab: {
      refreshModel: true,
    },
  },

  secretParam() {
    const { secret } = this.paramsFor(this.routeName);
    return secret ? normalizePath(secret) : '';
  },

  beforeModel() {
    const secret = this.secretParam();
    const backend = getEnginePathParam(this);
    const { tab } = this.paramsFor('vault.cluster.secrets.backend.list-root');
    const secretEngine = this.modelFor('vault.cluster.secrets.backend');
    const type = secretEngine?.engineType;
    const effectiveType = getEffectiveEngineType(type);
    assert('secretEngine.engineType is not defined', !!type);
    // if configuration only, redirect to configuration route
    if (engineDisplayData(effectiveType)?.isOnlyMountable) {
      return this.router.transitionTo('vault.cluster.secrets.backend.configuration', backend);
    }

    const engineRoute = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: true }).find(
      (engine) => engine.type === effectiveType
    )?.engineRoute;
    if (!type || !SUPPORTED_BACKENDS.includes(effectiveType)) {
      return this.router.transitionTo('vault.cluster.secrets');
    }
    if (this.routeName === 'vault.cluster.secrets.backend.list' && !secret.endsWith('/')) {
      return this.router.replaceWith('vault.cluster.secrets.backend.list', secret + '/');
    }
    if (isAddonEngine(effectiveType, secretEngine.version)) {
      if (engineRoute === 'kv.list' && pathIsDirectory(secret)) {
        return this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', backend, secret);
      }
      return this.router.transitionTo(`vault.cluster.secrets.backend.${engineRoute}`, backend);
    } else if (secretEngine.isV2KV) {
      // if it's KV v2 but not registered as an addon, it's type generic
      return this.router.transitionTo('vault.cluster.secrets.backend.kv.list', backend);
    }
    const modelType = this.getModelType(effectiveType, tab);
    return this.pathHelp.hydrateModel(modelType, backend).then(() => {
      this.store.unloadAll('capabilities');
    });
  },

  getModelType(type, tab) {
    return getModelTypeForEngine(type, { tab });
  },

  async model(params) {
    const secret = this.secretParam() || '';
    const backend = getEnginePathParam(this);
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    const effectiveType = getEffectiveEngineType(backendModel.engineType);
    const modelType = this.getModelType(effectiveType, params.tab);

    return hash({
      secret,
      secrets: this.pagination
        .lazyPaginatedQuery(modelType, {
          id: secret,
          backend,
          responsePath: 'data.keys',
          page: getValidPage(params.page),
          pageFilter: params.pageFilter,
        })
        .then((model) => {
          this.set('has404', false);
          return model;
        })
        .catch((err) => {
          if (backendModel && err.httpStatus === 404) {
            return [];
          } else {
            // else we're throwing and dealing with this in the error action
            throw err;
          }
        }),
    });
  },

  setupController(controller, resolvedModel) {
    const secretParams = this.paramsFor(this.routeName);
    const secret = resolvedModel.secret;
    const model = resolvedModel.secrets;
    const backend = getEnginePathParam(this);
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    const has404 = this.has404;
    // only clear store cache if this is a new model
    if (secret !== controller?.baseKey?.id) {
      this.pagination.clearDataset();
    }
    controller.set('hasModel', true);
    controller.setProperties({
      model,
      has404,
      backend,
      backendModel,
      baseKey: { id: secret },
      backendType: getEffectiveEngineType(backendModel.engineType),
    });
    if (!has404) {
      const pageFilter = secretParams.pageFilter;
      let filter;
      if (secret) {
        filter = secret + (pageFilter || '');
      } else if (pageFilter) {
        filter = pageFilter;
      }
      controller.setProperties({
        filter: filter || '',
        page: model.meta?.currentPage || 1,
      });
    }
  },

  resetController(controller, isExiting) {
    this._super(...arguments);
    if (isExiting) {
      controller.set('pageFilter', null);
      controller.set('filter', null);
    }
  },

  actions: {
    error(error, transition) {
      const secret = this.secretParam();
      const backend = getEnginePathParam(this);
      const is404 = error.httpStatus === 404;
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      const hasModel = this.controllerFor(this.routeName).hasModel;

      set(error, 'secret', secret);
      set(error, 'isRoot', true);
      set(error, 'backend', backend);
      // only swallow the error if we have a previous model
      if (hasModel && is404) {
        this.set('has404', true);
        transition.abort();
        return false;
      }
      return true;
    },

    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.pagination.clearDataset();
      }
      return true;
    },
    reload() {
      this.pagination.clearDataset();
      this.refresh();
    },
  },
});

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { set } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { allEngines, isAddonEngine } from 'vault/helpers/mountable-secret-engines';
import { service } from '@ember/service';
import { normalizePath } from 'vault/utils/path-encoding-helpers';
import { assert } from '@ember/debug';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';

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

  modelTypeForTransform(tab) {
    let modelType;
    switch (tab) {
      case 'role':
        modelType = 'transform/role';
        break;
      case 'template':
        modelType = 'transform/template';
        break;
      case 'alphabet':
        modelType = 'transform/alphabet';
        break;
      default: // CBS TODO: transform/transformation
        modelType = 'transform';
        break;
    }
    return modelType;
  },

  secretParam() {
    const { secret } = this.paramsFor(this.routeName);
    return secret ? normalizePath(secret) : '';
  },

  enginePathParam() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },

  beforeModel() {
    const secret = this.secretParam();
    const backend = this.enginePathParam();
    const { tab } = this.paramsFor('vault.cluster.secrets.backend.list-root');
    const secretEngine = this.store.peekRecord('secret-engine', backend);
    const type = secretEngine?.engineType;
    assert('secretEngine.engineType is not defined', !!type);
    const engineRoute = allEngines().find((engine) => engine.type === type)?.engineRoute;

    if (!type || !SUPPORTED_BACKENDS.includes(type)) {
      return this.router.transitionTo('vault.cluster.secrets');
    }
    if (this.routeName === 'vault.cluster.secrets.backend.list' && !secret.endsWith('/')) {
      return this.router.replaceWith('vault.cluster.secrets.backend.list', secret + '/');
    }
    if (isAddonEngine(type, secretEngine.version)) {
      if (engineRoute === 'kv.list' && pathIsDirectory(secret)) {
        return this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', backend, secret);
      }
      return this.router.transitionTo(`vault.cluster.secrets.backend.${engineRoute}`, backend);
    } else if (secretEngine.isV2KV) {
      // if it's KV v2 but not registered as an addon, it's type generic
      return this.router.transitionTo('vault.cluster.secrets.backend.kv.list', backend);
    }
    const modelType = this.getModelType(backend, tab);
    return this.pathHelp.getNewModel(modelType, backend).then(() => {
      this.store.unloadAll('capabilities');
    });
  },

  getModelType(backend, tab) {
    const secretEngine = this.store.peekRecord('secret-engine', backend);
    const type = secretEngine.engineType;
    const types = {
      database: tab === 'role' ? 'database/role' : 'database/connection',
      transit: 'transit-key',
      ssh: 'role-ssh',
      transform: this.modelTypeForTransform(tab),
      aws: 'role-aws',
      cubbyhole: 'secret',
      kv: 'secret',
      keymgmt: `keymgmt/${tab || 'key'}`,
      generic: 'secret',
    };
    return types[type];
  },

  async model(params) {
    const secret = this.secretParam() || '';
    const backend = this.enginePathParam();
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    const modelType = this.getModelType(backend, params.tab);

    return hash({
      secret,
      secrets: this.store
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
    const backend = this.enginePathParam();
    const backendModel = this.store.peekRecord('secret-engine', backend);
    const has404 = this.has404;
    // only clear store cache if this is a new model
    if (secret !== controller?.baseKey?.id) {
      this.store.clearAllDatasets();
    }
    controller.set('hasModel', true);
    controller.setProperties({
      model,
      has404,
      backend,
      backendModel,
      baseKey: { id: secret },
      backendType: backendModel.engineType,
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
      const backend = this.enginePathParam();
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
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
});

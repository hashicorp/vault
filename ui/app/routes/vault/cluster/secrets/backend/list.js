/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { set } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { allEngines, isAddonEngine } from 'vault/helpers/mountable-secret-engines';
import { inject as service } from '@ember/service';
import { normalizePath } from 'vault/utils/path-encoding-helpers';
import { assert } from '@ember/debug';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default Route.extend({
  store: service(),
  templateName: 'vault/cluster/secrets/backend/list',
  pathHelp: service('path-help'),
  router: service(),

  // By default assume user doesn't have permissions
  noMetadataPermissions: true,

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
    const engineRoute = allEngines().findBy('type', type)?.engineRoute;

    if (!type || !SUPPORTED_BACKENDS.includes(type)) {
      return this.router.transitionTo('vault.cluster.secrets');
    }
    if (this.routeName === 'vault.cluster.secrets.backend.list' && !secret.endsWith('/')) {
      return this.router.replaceWith('vault.cluster.secrets.backend.list', secret + '/');
    }
    if (isAddonEngine(type, secretEngine.version)) {
      return this.router.transitionTo(`vault.cluster.secrets.backend.${engineRoute}`, backend);
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
      // secret or secret-v2
      cubbyhole: 'secret',
      kv: secretEngine.modelTypeForKV,
      keymgmt: `keymgmt/${tab || 'key'}`,
      generic: secretEngine.modelTypeForKV,
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
          page: params.page || 1,
          pageFilter: params.pageFilter,
        })
        .then((model) => {
          this.set('noMetadataPermissions', false);
          this.set('has404', false);
          return model;
        })
        .catch((err) => {
          // if we're at the root we don't want to throw
          if (backendModel && err.httpStatus === 404 && secret === '') {
            this.set('noMetadataPermissions', false);
            return [];
          } else if (err.httpStatus === 403 && backendModel.isV2KV) {
            this.set('noMetadataPermissions', true);
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
    const noMetadataPermissions = this.noMetadataPermissions;
    // only clear store cache if this is a new model
    if (secret !== controller.get('baseKey.id')) {
      this.store.clearAllDatasets();
    }
    controller.set('hasModel', true);
    controller.setProperties({
      model,
      has404,
      noMetadataPermissions,
      backend,
      backendModel,
      baseKey: { id: secret },
      backendType: backendModel.get('engineType'),
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
      const hasModel = this.controllerFor(this.routeName).get('hasModel');

      // this will occur if we've deleted something,
      // and navigate to its parent and the parent doesn't exist -
      // this if often the case with nested keys in kv-like engines
      if (transition.data.isDeletion && is404) {
        throw error;
      }
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

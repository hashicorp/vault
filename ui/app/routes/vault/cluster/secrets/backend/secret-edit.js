/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { set } from '@ember/object';
import Ember from 'ember';
import { resolve } from 'rsvp';
import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { encodePath, normalizePath } from 'vault/utils/path-encoding-helpers';
import { keyIsFolder, parentKeyForKey } from 'core/utils/key-utils';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';
import { getModelTypeForEngine } from 'vault/utils/model-helpers/secret-engine-helpers';
import { getBackendEffectiveType, getEnginePathParam } from 'vault/utils/backend-route-helpers';
import { isValidProvider } from 'vault/utils/keymgmt-provider-utils';
import KeymgmtKeyForm from 'vault/forms/keymgmt/key';
import KeymgmtProviderForm from 'vault/forms/keymgmt/provider';
import {
  SecretsApiKeyManagementListKmsProvidersForKeyListEnum,
  SecretsApiTransformListRolesListEnum,
} from '@hashicorp/vault-client-typescript';

/**
 * @type Class
 */
export default Route.extend({
  store: service(),
  router: service(),
  pathHelp: service('path-help'),
  api: service(),
  capabilitiesService: service('capabilities'),

  secretParam() {
    const { secret } = this.paramsFor(this.routeName);
    return secret ? normalizePath(secret) : '';
  },

  capabilities(secret, modelType) {
    const backend = getEnginePathParam(this);
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    const backendType = getEffectiveEngineType(backendModel.engineType);
    let path;
    if (backendType === 'transit') {
      path = backend + '/keys/' + secret;
    } else if (backendType === 'ssh' || backendType === 'aws') {
      path = backend + '/roles/' + secret;
    } else if (modelType.startsWith('transform/')) {
      path = this.buildTransformPath(backend, secret, modelType);
    } else {
      path = backend + '/' + secret;
    }
    return this.store.findRecord('capabilities', path);
  },

  buildTransformPath(backend, secret, modelType) {
    const noun = modelType.split('/')[1];
    return `${backend}/${noun}/${secret}`;
  },

  transformSecretName(secret, modelType) {
    const noun = modelType.split('/')[1];
    return secret.replace(`${noun}/`, '');
  },

  async fetchTransformRoles(backend) {
    try {
      const { keys } = await this.api.secrets.transformListRoles(
        backend,
        SecretsApiTransformListRolesListEnum.TRUE
      );
      return keys || [];
    } catch (error) {
      return [];
    }
  },

  backendType() {
    return getBackendEffectiveType(this);
  },

  templateName: 'vault/cluster/secrets/backend/secretEditLayout',

  beforeModel({ to: { queryParams } }) {
    const secret = this.secretParam();
    const secretEngine = this.modelFor('vault.cluster.secrets.backend');
    return this.buildModel(secret, queryParams).then(() => {
      const parentKey = parentKeyForKey(secret);
      const mode = this.routeName.split('.').pop();
      // for kv v2, redirect users from the old url to the new engine url (1.15.0 +)
      if (secretEngine.type === 'kv' && secretEngine.version === 2) {
        let route, params;
        switch (true) {
          case !secret:
            // if no secret param redirect to the create route
            route = 'vault.cluster.secrets.backend.kv.create';
            params = [secretEngine.id];
            break;
          case this.routeName === 'vault.cluster.secrets.backend.show':
            route = 'vault.cluster.secrets.backend.kv.secret.index';
            params = [secretEngine.id, secret];
            break;
          case this.routeName === 'vault.cluster.secrets.backend.edit':
            route = 'vault.cluster.secrets.backend.kv.secret.details.edit';
            params = [secretEngine.id, secret];
            break;
          default:
            route = 'vault.cluster.secrets.backend.kv.secret.index';
            params = [secretEngine.id, secret];
            break;
        }
        this.router.transitionTo(route, ...params);
        return;
      }
      if (mode === 'edit' && keyIsFolder(secret)) {
        if (parentKey) {
          this.router.transitionTo('vault.cluster.secrets.backend.list', encodePath(parentKey));
        } else {
          this.router.transitionTo('vault.cluster.secrets.backend.list-root');
        }
        return;
      }
    });
  },

  buildModel(secret, queryParams) {
    const backend = getEnginePathParam(this);
    const modelType = this.modelType(backend, secret, { queryParams });
    // Keymgmt resources are loaded through API-backed forms, so Ember Data hydration is unnecessary.
    if (modelType === 'secret' || modelType.startsWith('keymgmt/')) {
      return resolve();
    }
    return this.pathHelp.hydrateModel(modelType, backend);
  },

  modelType(backend, secret, options = {}) {
    const backendModel = this.modelFor('vault.cluster.secrets.backend', backend);
    const { engineType } = backendModel;
    const effectiveType = getEffectiveEngineType(engineType);

    return getModelTypeForEngine(effectiveType, {
      secret,
      itemType: options.queryParams?.itemType,
    });
  },

  async fetchKeyProvider(name, backend) {
    try {
      const providerResp = await this.api.secrets.keyManagementListKmsProvidersForKey(
        name,
        backend,
        SecretsApiKeyManagementListKmsProvidersForKeyListEnum.TRUE
      );
      return providerResp.keys?.[0] || null;
    } catch (e) {
      const { status } = await this.api.parseError(e);
      if (status === 403) {
        return { permissionsError: true };
      }
      if (status === 404) {
        return [];
      }
      throw e;
    }
  },

  async fetchKeyDistribution(name, provider, backend) {
    try {
      const distResp = await this.api.secrets.keyManagementReadKeyInKmsProvider(name, provider, backend);
      return {
        ...distResp.data,
        purposeArray: distResp.data?.purpose?.split(',') || [],
      };
    } catch (e) {
      const { status } = await this.api.parseError(e);
      // Return null for 403 - distribution is optional, no need for permissionsError like provider
      if (status === 403) {
        return null;
      }
      throw e;
    }
  },

  buildKeyVersionsData(versions) {
    let created = null;
    let last_rotated = null;
    let versionsArray = [];

    if (versions) {
      const versionKeys = Object.keys(versions);
      if (versionKeys.length > 0) {
        // This computes value of "created" from first version
        const firstKey = versionKeys[0];
        const firstVersion = versions[firstKey];
        if (firstVersion?.creation_time) {
          created = new Date(firstVersion.creation_time);
        }

        // This computes value of "last_rotated" from last version (if more than one)
        if (versionKeys.length > 1) {
          const lastKey = versionKeys[versionKeys.length - 1];
          const lastVersion = versions[lastKey];
          if (lastVersion?.creation_time) {
            last_rotated = new Date(lastVersion.creation_time);
          }
        }

        versionsArray = versionKeys
          .map((key) => {
            const version = versions[key];
            if (!version?.creation_time) return null;
            return {
              ...version,
              id: parseInt(key, 10),
            };
          })
          .filter((v) => v !== null);
      }
    }

    return { created, last_rotated, versions: versionsArray };
  },

  async fetchKeymgmtKey(backend, name) {
    const { data } = await this.api.secrets.keyManagementReadKey(name, backend);

    const provider = await this.fetchKeyProvider(name, backend);

    let distribution = null;
    if (isValidProvider(provider)) {
      distribution = await this.fetchKeyDistribution(name, provider, backend);
    }

    // This builds computed version data.
    const { created, last_rotated, versions } = this.buildKeyVersionsData(data.versions);

    const form = new KeymgmtKeyForm(
      {
        ...data,
        name,
        backend,
        provider,
        distribution,
        created,
        last_rotated,
        versions,
      },
      { isNew: false }
    );
    return form;
  },

  async fetchKeymgmtKeyCapabilities(backend, name) {
    const keyPath = this.capabilitiesService.pathFor('keymgmtKey', { backend, name });
    const keysPath = this.capabilitiesService.pathFor('keymgmtKeys', { backend });
    const keyProvidersPath = this.capabilitiesService.pathFor('keymgmtKeyProviders', { backend, name });

    const capabilities = await this.capabilitiesService.fetch([keyPath, keysPath, keyProvidersPath]);

    return {
      canDelete: capabilities[keyPath]?.canDelete,
      canUpdate: capabilities[keyPath]?.canUpdate,
      canEdit: capabilities[keyPath]?.canUpdate,
      canRead: capabilities[keyPath]?.canRead,
      canList: capabilities[keysPath]?.canList,
      canListProviders: capabilities[keyProvidersPath]?.canList,
    };
  },

  async fetchKeymgmtProvider(backend, name) {
    const { data } = await this.api.secrets.keyManagementReadKmsProvider(name, backend);

    const form = new KeymgmtProviderForm(
      {
        ...data,
        name,
        backend,
        keys: [],
      },
      { isNew: false }
    );

    return form;
  },

  async fetchKeymgmtProviderCapabilities(backend, name) {
    const providerPath = this.capabilitiesService.pathFor('keymgmtProvider', { backend, id: name });
    const providersPath = this.capabilitiesService.pathFor('keymgmtProviders', { backend });
    const providerKeysPath = this.capabilitiesService.pathFor('keymgmtProviderKeys', { backend, id: name });

    const capabilities = await this.capabilitiesService.fetch([
      providerPath,
      providersPath,
      providerKeysPath,
    ]);

    return {
      canDelete: capabilities[providerPath]?.canDelete,
      canUpdate: capabilities[providerPath]?.canUpdate,
      canEdit: capabilities[providerPath]?.canUpdate,
      canRead: capabilities[providerPath]?.canRead,
      canList: capabilities[providersPath]?.canList,
      canListKeys: capabilities[providerKeysPath]?.canList,
      canCreateKeys: capabilities[providerKeysPath]?.canCreate,
    };
  },

  async handleSecretModelError(capabilitiesPromise, secretId, modelType, error) {
    // capabilities is a promise proxy, not a real object
    // to work around this we explicitly assign it to a const and await it
    const capabilities = await capabilitiesPromise;
    // can't read the path and don't have update capability, so re-throw
    if (!capabilities.canUpdate && modelType === 'secret') {
      throw error;
    }
    this.store.push({
      data: {
        id: secretId,
        type: modelType,
        attributes: {
          failedServerRead: true,
        },
      },
    });
    const secretModel = this.store.peekRecord(modelType, secretId);
    return secretModel;
  },

  async model(params, { to: { queryParams } }) {
    let secret = this.secretParam();
    const backend = getEnginePathParam(this);
    const modelType = this.modelType(backend, secret, { queryParams });
    const type = params.type || '';
    if (!secret) {
      secret = '\u0020';
    }
    if (modelType.startsWith('transform/')) {
      secret = this.transformSecretName(secret, modelType);
    }
    if (modelType === 'database/role') {
      secret = secret.replace('role/', '');
    }
    let secretModel;
    let transformRoles;
    let capabilities;

    // Handle keymgmt resources with API service
    if (modelType === 'keymgmt/key') {
      secretModel = await this.fetchKeymgmtKey(backend, secret);
      capabilities = await this.fetchKeymgmtKeyCapabilities(backend, secret);
    } else if (modelType === 'keymgmt/provider') {
      secretModel = await this.fetchKeymgmtProvider(backend, secret);
      capabilities = await this.fetchKeymgmtProviderCapabilities(backend, secret);
    } else {
      capabilities = await this.capabilities(secret, modelType);
      try {
        secretModel = await this.store.queryRecord(modelType, { id: secret, backend, type });
      } catch (err) {
        // we've failed the read request, but if it's a kv-v1 type backend, we want to
        // do additional checks of the capabilities
        if (err.httpStatus === 403 && modelType === 'secret') {
          secretModel = await this.handleSecretModelError(capabilities, secret, modelType, err);
        } else {
          throw err;
        }
      }
    }

    // fetch roles for transform type to display in detail view
    if (modelType === 'transform') {
      transformRoles = await this.fetchTransformRoles(backend);
    }

    return {
      secret: secretModel,
      capabilities,
      transformRoles,
    };
  },

  setupController(controller, model) {
    this._super(...arguments);
    const secret = this.secretParam();
    const backend = getEnginePathParam(this);
    const preferAdvancedEdit =
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      this.controllerFor('vault.cluster.secrets.backend').preferAdvancedEdit || false;
    const backendType = this.backendType();
    // mode will be 'show', 'edit', 'create'
    const mode = this.routeName.split('.').pop().replace('-root', '');

    // Handle keymgmt forms differently - Resource or Form doesn't have setProperties
    const modelType = this.modelType(backend, secret);
    if (!['keymgmt/key', 'keymgmt/provider'].includes(modelType)) {
      model.secret.setProperties({ backend });
    }

    controller.setProperties({
      model: model.secret,
      form: ['keymgmt/key', 'keymgmt/provider'].includes(modelType) ? model.secret : null,
      capabilities: model.capabilities,
      baseKey: { id: secret },
      mode,
      backend,
      preferAdvancedEdit,
      backendType,
      transformRoles: model.transformRoles,
    });
  },

  resetController(controller) {
    if (controller.reset && typeof controller.reset === 'function') {
      controller.reset();
    }
  },

  actions: {
    error(error) {
      const secret = this.secretParam();
      const backend = getEnginePathParam(this);
      set(error, 'keyId', backend + '/' + secret);
      set(error, 'backend', backend);
      return true;
    },

    refreshModel() {
      this.refresh();
    },

    willTransition(transition) {
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      const { mode, model } = this.controller;

      // If model is clean or deleted, continue
      if (!model.hasDirtyAttributes || model.isDeleted) {
        return true;
      }

      const changed = model.changedAttributes();
      const changedKeys = Object.keys(changed);

      // until we have time to move `backend` on a v1 model to a relationship,
      // it's going to dirty the model state, so we need to look for it
      // and explicity ignore it here
      if (mode !== 'show' && changedKeys.length && changedKeys[0] !== 'backend') {
        if (
          Ember.testing ||
          window.confirm(
            'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
          )
        ) {
          model && model.rollbackAttributes();
          return true;
        } else {
          transition.abort();
          return false;
        }
      }
      return this._super(...arguments);
    },
  },
});

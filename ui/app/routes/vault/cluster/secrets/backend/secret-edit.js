/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { set } from '@ember/object';
import Ember from 'ember';
import { resolve } from 'rsvp';
import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { encodePath, normalizePath } from 'vault/utils/path-encoding-helpers';
import { keyIsFolder, parentKeyForKey } from 'core/utils/key-utils';

export default Route.extend({
  store: service(),
  router: service(),
  pathHelp: service('path-help'),

  secretParam() {
    const { secret } = this.paramsFor(this.routeName);
    return secret ? normalizePath(secret) : '';
  },

  enginePathParam() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },

  capabilities(secret, modelType) {
    const backend = this.enginePathParam();
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    const backendType = backendModel.engineType;
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

  modelTypeForTransform(secretName) {
    if (!secretName) return 'transform';
    if (secretName.startsWith('role/')) {
      return 'transform/role';
    }
    if (secretName.startsWith('template/')) {
      return 'transform/template';
    }
    if (secretName.startsWith('alphabet/')) {
      return 'transform/alphabet';
    }
    return 'transform'; // TODO: transform/transformation;
  },

  transformSecretName(secret, modelType) {
    const noun = modelType.split('/')[1];
    return secret.replace(`${noun}/`, '');
  },

  backendType() {
    return this.modelFor('vault.cluster.secrets.backend').engineType;
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
        // if no secret param redirect to the create route
        // if secret param they are either viewing or editing secret so navigate to the details route
        if (!secret) {
          this.router.transitionTo('vault.cluster.secrets.backend.kv.create', secretEngine.id);
        } else {
          this.router.transitionTo(
            'vault.cluster.secrets.backend.kv.secret.details',
            secretEngine.id,
            secret
          );
        }
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
    const backend = this.enginePathParam();
    const modelType = this.modelType(backend, secret, { queryParams });
    if (modelType === 'secret') {
      return resolve();
    }
    return this.pathHelp.getNewModel(modelType, backend);
  },

  modelType(backend, secret, options = {}) {
    const backendModel = this.modelFor('vault.cluster.secrets.backend', backend);
    const { engineType } = backendModel;
    const types = {
      database: secret && secret.startsWith('role/') ? 'database/role' : 'database/connection',
      transit: 'transit-key',
      ssh: 'role-ssh',
      transform: this.modelTypeForTransform(secret),
      aws: 'role-aws',
      cubbyhole: 'secret',
      kv: 'secret',
      keymgmt: `keymgmt/${options.queryParams?.itemType || 'key'}`,
      generic: 'secret',
    };
    return types[engineType];
  },

  handleSecretModelError(capabilities, secretId, modelType, error) {
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
    const backend = this.enginePathParam();
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

    const capabilities = this.capabilities(secret, modelType);
    try {
      secretModel = await this.store.queryRecord(modelType, { id: secret, backend, type });
    } catch (err) {
      // we've failed the read request, but if it's a kv-v1 type backend, we want to
      // do additional checks of the capabilities
      if (err.httpStatus === 403 && modelType === 'secret') {
        await capabilities;
        secretModel = this.handleSecretModelError(capabilities, secret, modelType, err);
      } else {
        throw err;
      }
    }
    await capabilities;

    return {
      secret: secretModel,
      capabilities,
    };
  },

  setupController(controller, model) {
    this._super(...arguments);
    const secret = this.secretParam();
    const backend = this.enginePathParam();
    const preferAdvancedEdit =
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      this.controllerFor('vault.cluster.secrets.backend').preferAdvancedEdit || false;
    const backendType = this.backendType();
    model.secret.setProperties({ backend });
    controller.setProperties({
      model: model.secret,
      capabilities: model.capabilities,
      baseKey: { id: secret },
      // mode will be 'show', 'edit', 'create'
      mode: this.routeName.split('.').pop().replace('-root', ''),
      backend,
      preferAdvancedEdit,
      backendType,
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
      const backend = this.enginePathParam();
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

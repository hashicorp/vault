import { set } from '@ember/object';
import { hash, resolve } from 'rsvp';
import { inject as service } from '@ember/service';
import DS from 'ember-data';
import Route from '@ember/routing/route';
import utils from 'vault/lib/key-utils';
import { getOwner } from '@ember/application';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import { encodePath, normalizePath } from 'vault/utils/path-encoding-helpers';

export default Route.extend(UnloadModelRoute, {
  pathHelp: service('path-help'),
  secretParam() {
    let { secret } = this.paramsFor(this.routeName);
    return secret ? normalizePath(secret) : '';
  },
  enginePathParam() {
    let { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },
  capabilities(secret) {
    const backend = this.enginePathParam();
    let backendModel = this.modelFor('vault.cluster.secrets.backend');
    let backendType = backendModel.get('engineType');
    if (backendType === 'kv' || backendType === 'cubbyhole' || backendType === 'generic') {
      return resolve({});
    }
    let path;
    if (backendType === 'transit') {
      path = backend + '/keys/' + secret;
    } else if (backendType === 'ssh' || backendType === 'aws') {
      path = backend + '/roles/' + secret;
    } else {
      path = backend + '/' + secret;
    }
    return this.store.findRecord('capabilities', path);
  },

  backendType() {
    return this.modelFor('vault.cluster.secrets.backend').get('engineType');
  },

  templateName: 'vault/cluster/secrets/backend/secretEditLayout',

  beforeModel() {
    // currently there is no recursive delete for folders in vault, so there's no need to 'edit folders'
    // perhaps in the future we could recurse _for_ users, but for now, just kick them
    // back to the list
    let secret = this.secretParam();
    return this.buildModel(secret).then(() => {
      const parentKey = utils.parentKeyForKey(secret);
      const mode = this.routeName.split('.').pop();
      if (mode === 'edit' && utils.keyIsFolder(secret)) {
        if (parentKey) {
          return this.transitionTo('vault.cluster.secrets.backend.list', encodePath(parentKey));
        } else {
          return this.transitionTo('vault.cluster.secrets.backend.list-root');
        }
      }
    });
  },

  buildModel(secret) {
    const backend = this.enginePathParam();

    let modelType = this.modelType(backend, secret);
    if (['secret', 'secret-v2'].includes(modelType)) {
      return resolve();
    }
    let owner = getOwner(this);
    return this.pathHelp.getNewModel(modelType, owner, backend);
  },

  modelType(backend, secret) {
    let backendModel = this.modelFor('vault.cluster.secrets.backend', backend);
    let type = backendModel.get('engineType');
    let types = {
      transit: 'transit-key',
      ssh: 'role-ssh',
      aws: 'role-aws',
      pki: secret && secret.startsWith('cert/') ? 'pki-certificate' : 'role-pki',
      cubbyhole: 'secret',
      kv: backendModel.get('modelTypeForKV'),
      generic: backendModel.get('modelTypeForKV'),
    };
    return types[type];
  },

  model(params) {
    let secret = this.secretParam();
    let backend = this.enginePathParam();
    let backendModel = this.modelFor('vault.cluster.secrets.backend', backend);
    let modelType = this.modelType(backend, secret);

    if (!secret) {
      secret = '\u0020';
    }
    if (modelType === 'pki-certificate') {
      secret = secret.replace('cert/', '');
    }
    return hash({
      secret: this.store
        .queryRecord(modelType, { id: secret, backend })
        .then(secretModel => {
          if (modelType === 'secret-v2') {
            let targetVersion = parseInt(params.version || secretModel.currentVersion, 10);
            let version = secretModel.versions.findBy('version', targetVersion);
            // 404 if there's no version
            if (!version) {
              let error = new DS.AdapterError();
              set(error, 'httpStatus', 404);
              throw error;
            }
            secretModel.set('engine', backendModel);

            return version.reload().then(() => {
              secretModel.set('selectedVersion', version);
              return secretModel;
            });
          }
          return secretModel;
        })
        .catch(err => {
          //don't have access to the metadata, so we'll make
          //a stub metadata model and try to load the version
          if (modelType === 'secret-v2' && err.httpStatus === 403) {
            let secretModel = this.store.createRecord('secret-v2');
            secretModel.setProperties({
              engine: backendModel,
              id: secret,
              // so we know it's a stub model and won't be saving it
              // because we don't have access to that endpoint
              isStub: true,
            });
            let targetVersion = params.version ? parseInt(params.version, 10) : null;
            let versionId = targetVersion ? [backend, secret, targetVersion] : [backend, secret];
            return this.store
              .findRecord('secret-v2-version', JSON.stringify(versionId), { reload: true })
              .then(versionModel => {
                secretModel.set('selectedVersion', versionModel);
                return secretModel;
              });
          }
          throw err;
        }),
      capabilities: this.capabilities(secret),
    });
  },

  setupController(controller, model) {
    this._super(...arguments);
    let secret = this.secretParam();
    let backend = this.enginePathParam();
    const preferAdvancedEdit =
      this.controllerFor('vault.cluster.secrets.backend').get('preferAdvancedEdit') || false;
    const backendType = this.backendType();
    model.secret.setProperties({ backend });
    controller.setProperties({
      model: model.secret,
      capabilities: model.capabilities,
      baseKey: { id: secret },
      // mode will be 'show', 'edit', 'create'
      mode: this.routeName
        .split('.')
        .pop()
        .replace('-root', ''),
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
      let secret = this.secretParam();
      let backend = this.enginePathParam();
      set(error, 'keyId', backend + '/' + secret);
      set(error, 'backend', backend);
      return true;
    },

    refreshModel() {
      this.refresh();
    },

    willTransition(transition) {
      let { mode, model } = this.controller;
      let version = model.get('selectedVersion');
      let changed = model.changedAttributes();
      let changedKeys = Object.keys(changed);
      // until we have time to move `backend` on a v1 model to a relationship,
      // it's going to dirty the model state, so we need to look for it
      // and explicity ignore it here
      if (
        (mode !== 'show' && (changedKeys.length && changedKeys[0] !== 'backend')) ||
        (mode !== 'show' && version && Object.keys(version.changedAttributes()).length)
      ) {
        if (
          window.confirm(
            'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
          )
        ) {
          version && version.rollbackAttributes();
          this.unloadModel();
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

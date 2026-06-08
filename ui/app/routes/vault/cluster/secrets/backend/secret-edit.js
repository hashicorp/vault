/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { set, get } from '@ember/object';
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
import TotpKeyForm from 'vault/forms/totp/key';
import clamp from 'vault/utils/clamp';
import TransitKeyForm from 'vault/forms/transit/key';

import SshRoleForm from 'vault/forms/ssh/role';
import AlphabetForm from 'vault/forms/transform/alphabet';
import TemplateForm from 'vault/forms/transform/template';
import RoleForm from 'vault/forms/transform/role';
import TransformationForm from 'vault/forms/transform/transformation';
import Form from 'vault/forms/form';
import {
  SecretsApiKeyManagementListKmsProvidersForKeyListEnum,
  SecretsApiTransformListRolesListEnum,
} from '@hashicorp/vault-client-typescript';

// TODO: Move this into a util file or class
const ACTION_VALUES = {
  encrypt: {
    isSupported: 'supports_encryption',
    description: 'Looks up wrapping properties for the given token.',
    glyph: 'lock-fill',
  },
  decrypt: {
    isSupported: 'supports_decryption',
    description: 'Decrypts the provided ciphertext using this key.',
    glyph: 'mail-open',
  },
  datakey: {
    isSupported: 'supports_encryption',
    description: 'Generates a new key and value encrypted with this key.',
    glyph: 'key',
  },
  rewrap: {
    isSupported: 'supports_encryption',
    description: 'Rewraps the ciphertext using the latest version of the named key.',
    glyph: 'reload',
  },
  sign: {
    isSupported: 'supports_signing',
    description: 'Get the cryptographic signature of the given data.',
    glyph: 'pencil-tool',
  },
  hmac: {
    isSupported: true,
    description: 'Generate a data digest using a hash algorithm.',
    glyph: 'shuffle',
  },
  verify: {
    isSupported: true,
    description: 'Validate the provided signature for the given data.',
    glyph: 'check-circle',
  },
  export: {
    isSupported: 'exportable',
    description: 'Get the named key.',
    glyph: 'external-link',
  },
};

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
    // TODO: Remove buildModel once all remaining engine types (e.g. database, transit) are migrated
    // to API-backed Form instances — none will need hydrateModel and this function becomes redundant.
    const skipHydration =
      modelType === 'secret' ||
      modelType.startsWith('keymgmt/') ||
      modelType.startsWith('transform') ||
      modelType === 'totp-key' ||
      modelType === 'role-ssh';
    if (skipHydration) {
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

  async fetchTotpKey(backend, name) {
    const resp = await this.api.secrets.totpReadKey(name, backend);
    const data = resp.data || {};
    return new TotpKeyForm(
      {
        ...data,
        name,
        backend,
      },
      { isNew: false }
    );
  },

  async fetchTotpKeyCapabilities(backend, name) {
    const keyPath = this.capabilitiesService.pathFor('totpKey', { backend, name });
    const keysPath = this.capabilitiesService.pathFor('totpKeys', { backend });

    const capabilities = await this.capabilitiesService.fetch([keyPath, keysPath]);

    return {
      canDelete: capabilities[keyPath]?.canDelete,
      canRead: capabilities[keyPath]?.canRead,
      canList: capabilities[keysPath]?.canList,
    };
  },

  async fetchTransitKey(name, backend) {
    const res = await this.api.secrets.transitReadKey(name, backend);

    const transitModel = {
      backend,
      id: name,
      ...res,
      ...this.transitEncryptionKeyVersions(
        res.keys,
        res.min_decryption_version,
        res.min_encryption_version,
        res.latest_version,
        res.supports_signing,
        res.supports_encryption
      ),
    };

    return transitModel;
  },

  async fetchTransitKeyCapabilities(backend, name, secretModel) {
    const rotatePath = this.capabilitiesService.pathFor('transitKeyRotate', { backend, name });
    const keyPath = this.capabilitiesService.pathFor('transitKey', { backend, name });

    const capabilities = await this.capabilitiesService.fetch([rotatePath, keyPath]);

    return {
      canRotate: capabilities[rotatePath]?.canUpdate,
      canRead: capabilities[keyPath]?.canRead,
      canUpdate: capabilities[keyPath]?.canUpdate,
      canDelete: capabilities[keyPath]?.canDelete !== false && secretModel.deletion_allowed !== false, // check deletion_allowed from the model in addition to the capability since both are required to allow deletion
    };
  },

  // TODO: Move these into separate classes or utils
  transitEncryptionKeyVersions(
    keys,
    minDecryptionVersion,
    minEncryptionVersion,
    latestVersion,
    supportsSigning,
    supportsEncryption
  ) {
    const validKeys = Object.keys(keys);
    const keyVersions = [];
    const keysForEncryption = [];

    // get keyVersions
    let maxVersion = Math.max(...validKeys);
    while (maxVersion > 0) {
      keyVersions.unshift(maxVersion);
      maxVersion--;
    }

    // get encryptionKeyVersions using keyVersions
    const encryptionKeyVersions = keyVersions
      .filter((version) => {
        return version >= minDecryptionVersion;
      })
      .reverse();

    // get keysForEncryption
    const minVersion = clamp(minEncryptionVersion - 1, 0, latestVersion);
    while (latestVersion > minVersion) {
      keysForEncryption.push(latestVersion);
      latestVersion--;
    }

    // get exportKeyTypes
    const exportKeyTypes = ['hmac'];
    if (supportsSigning) {
      exportKeyTypes.unshift('signing');
    }
    if (supportsEncryption) {
      exportKeyTypes.unshift('encryption');
    }

    return {
      encryptionKeyVersions,
      keyVersions,
      keysForEncryption,
      exportKeyTypes,
      validKeyVersions: Object.keys(keys),
    };
  },

  transitSupportedActions(secretModel) {
    return Object.keys(ACTION_VALUES)
      .filter((name) => {
        const { isSupported } = ACTION_VALUES[name];
        return typeof isSupported === 'boolean' || get(secretModel, isSupported);
      })
      .map((name) => {
        const { description, glyph } = ACTION_VALUES[name];
        return { name, description, glyph };
      });
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

  async fetchSshRole(backend, name) {
    try {
      const { data } = await this.api.secrets.sshReadRole(name, backend);
      return new SshRoleForm({ ...data, name, backend }, { isNew: false });
    } catch (error) {
      const { message } = await this.api.parseError(error);
      throw new Error(message);
    }
  },

  async fetchSshRoleCapabilities(backend, name) {
    try {
      const rolePath = this.capabilitiesService.pathFor('sshRole', { backend, id: name });
      const credentialsPath = this.capabilitiesService.pathFor('sshCredentials', { backend, id: name });
      const signPath = this.capabilitiesService.pathFor('sshSign', { backend, id: name });

      const capabilities = await this.capabilitiesService.fetch([rolePath, credentialsPath, signPath]);

      return {
        canDelete: capabilities[rolePath]?.canDelete,
        canUpdate: capabilities[rolePath]?.canUpdate,
        canEdit: capabilities[rolePath]?.canUpdate,
        canRead: capabilities[rolePath]?.canRead,
        canGenerate: capabilities[credentialsPath]?.canUpdate,
        canSign: capabilities[signPath]?.canUpdate,
      };
    } catch (error) {
      const { message } = await this.api.parseError(error);
      throw new Error(message);
    }
  },

  async fetchTransformAlphabet(backend, name) {
    const resp = await this.api.secrets.transformReadAlphabet(name, backend);
    const data = resp.data || {};
    return new AlphabetForm({ ...data, name, backend }, { isNew: false });
  },

  async fetchTransformAlphabetCapabilities(backend, name) {
    const alphabetPath = this.capabilitiesService.pathFor('transformAlphabet', { backend, name });
    const alphabetsPath = this.capabilitiesService.pathFor('transformAlphabets', { backend });

    const capabilities = await this.capabilitiesService.fetch([alphabetPath, alphabetsPath]);

    return {
      canDelete: capabilities[alphabetPath]?.canDelete,
      canUpdate: capabilities[alphabetPath]?.canUpdate,
      canRead: capabilities[alphabetPath]?.canRead,
      canList: capabilities[alphabetsPath]?.canList,
    };
  },

  async fetchTransformTemplate(backend, name) {
    const resp = await this.api.secrets.transformReadTemplate(name, backend);
    const data = resp.data || {};
    return new TemplateForm(
      {
        name,
        backend,
        type: data.type,
        pattern: data.pattern,
        alphabet: data.alphabet ? [data.alphabet] : [],
        encode_format: data.encode_format,
        decode_formats: data.decode_formats,
      },
      { isNew: false }
    );
  },

  async fetchTransformTemplateCapabilities(backend, name) {
    const templatePath = this.capabilitiesService.pathFor('transformTemplate', { backend, name });
    const templatesPath = this.capabilitiesService.pathFor('transformTemplates', { backend });

    const capabilities = await this.capabilitiesService.fetch([templatePath, templatesPath]);

    return {
      canDelete: capabilities[templatePath]?.canDelete,
      canUpdate: capabilities[templatePath]?.canUpdate,
      canRead: capabilities[templatePath]?.canRead,
      canList: capabilities[templatesPath]?.canList,
    };
  },

  async fetchTransformRole(backend, name) {
    const resp = await this.api.secrets.transformReadRole(name, backend);
    const data = resp.data || {};
    return new RoleForm({ name, backend, transformations: data.transformations || [] }, { isNew: false });
  },

  async fetchTransformRoleCapabilities(backend, name) {
    const rolePath = this.capabilitiesService.pathFor('transformRole', { backend, name });
    const rolesPath = this.capabilitiesService.pathFor('transformRoles', { backend });

    const capabilities = await this.capabilitiesService.fetch([rolePath, rolesPath]);

    return {
      canDelete: capabilities[rolePath]?.canDelete,
      canUpdate: capabilities[rolePath]?.canUpdate,
      canRead: capabilities[rolePath]?.canRead,
      canList: capabilities[rolesPath]?.canList,
    };
  },

  async fetchTransformTransformation(backend, name) {
    const resp = await this.api.secrets.transformReadTransformation(name, backend);
    const data = resp.data || resp || {};
    return new TransformationForm(
      {
        name,
        backend,
        type: data.type || 'fpe',
        tweak_source: data.tweak_source,
        masking_character: data.masking_character,
        template: data.template ? [data.template] : [],
        allowed_roles: data.allowed_roles || [],
        deletion_allowed: data.deletion_allowed,
        mapping_mode: data.mapping_mode,
        convergent: data.convergent,
        max_ttl: data.max_ttl,
        stores: data.stores || [],
      },
      { isNew: false }
    );
  },

  async fetchTransformTransformationCapabilities(backend, name) {
    const transformationPath = this.capabilitiesService.pathFor('transformTransformation', { backend, name });
    const transformationsPath = this.capabilitiesService.pathFor('transformTransformations', { backend });

    const capabilities = await this.capabilitiesService.fetch([transformationPath, transformationsPath]);

    return {
      canDelete: capabilities[transformationPath]?.canDelete,
      canUpdate: capabilities[transformationPath]?.canUpdate,
      canRead: capabilities[transformationPath]?.canRead,
      canList: capabilities[transformationsPath]?.canList,
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
    let capabilities;

    if (modelType === 'totp-key') {
      secretModel = await this.fetchTotpKey(backend, secret);
      capabilities = await this.fetchTotpKeyCapabilities(backend, secret);
    } else if (modelType === 'transit-key') {
      secretModel = await this.fetchTransitKey(secret, backend);
      capabilities = await this.fetchTransitKeyCapabilities(backend, secret, secretModel);
      secretModel.supportedActions = this.transitSupportedActions(secretModel);
      // replace secretModel with form
      secretModel = new TransitKeyForm(secretModel, { isNew: false });
    } else if (modelType === 'keymgmt/key') {
      secretModel = await this.fetchKeymgmtKey(backend, secret);
      capabilities = await this.fetchKeymgmtKeyCapabilities(backend, secret);
    } else if (modelType === 'keymgmt/provider') {
      secretModel = await this.fetchKeymgmtProvider(backend, secret);
      capabilities = await this.fetchKeymgmtProviderCapabilities(backend, secret);
    } else if (modelType === 'role-ssh') {
      secretModel = await this.fetchSshRole(backend, secret);
      capabilities = await this.fetchSshRoleCapabilities(backend, secret);
    } else if (modelType === 'transform/alphabet') {
      secretModel = await this.fetchTransformAlphabet(backend, secret);
      capabilities = await this.fetchTransformAlphabetCapabilities(backend, secret);
    } else if (modelType === 'transform/template') {
      secretModel = await this.fetchTransformTemplate(backend, secret);
      capabilities = await this.fetchTransformTemplateCapabilities(backend, secret);
    } else if (modelType === 'transform/role') {
      secretModel = await this.fetchTransformRole(backend, secret);
      capabilities = await this.fetchTransformRoleCapabilities(backend, secret);
    } else if (modelType === 'transform') {
      secretModel = await this.fetchTransformTransformation(backend, secret);
      capabilities = await this.fetchTransformTransformationCapabilities(backend, secret);
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

    return {
      secret: secretModel,
      capabilities,
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

    // Form-based models (keymgmt, TOTP, SSH, transform sub-types) use Form class instances which
    // don't have the Ember Data setProperties method.
    const isFormModel = model.secret instanceof Form;
    if (!isFormModel) {
      model.secret.setProperties({ backend });
    }

    controller.setProperties({
      model: model.secret,
      form: isFormModel ? model.secret : null,
      capabilities: model.capabilities,
      baseKey: { id: secret },
      mode,
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

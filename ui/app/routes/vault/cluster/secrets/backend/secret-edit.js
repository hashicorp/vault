/**
 * Copyright IBM Corp. 2016, 2026
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
import TotpKeyForm from 'vault/forms/totp/key';
import TransitKeyForm from 'vault/forms/transit/key';
import TransitKeyResource from 'vault/resources/transit/key';
import DatabaseConnectionForm, { normalizeDatabaseConnection } from 'vault/forms/database/connection';
import DatabaseRoleForm, { normalizeDatabaseRole } from 'vault/forms/database/role';

import SshRoleForm from 'vault/forms/ssh/role';
import RoleAwsForm from 'vault/forms/aws-role';
import AlphabetForm from 'vault/forms/transform/alphabet';
import TemplateForm from 'vault/forms/transform/template';
import RoleForm from 'vault/forms/transform/role';
import TransformationForm from 'vault/forms/transform/transformation';
import Form from 'vault/forms/form';
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
  api: service(),
  capabilitiesService: service('capabilities'),
  version: service(),

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

  beforeModel() {
    const secret = this.secretParam();
    const secretEngine = this.modelFor('vault.cluster.secrets.backend');
    return resolve().then(() => {
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

    return new TransitKeyResource({
      backend,
      id: name,
      ...res,
    });
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

  async fetchDatabaseConnection(backend, name) {
    try {
      const { data } = await this.api.secrets.databaseReadConnectionConfiguration(name, backend);
      const connection = normalizeDatabaseConnection(data, name, backend);
      // the edit view changes the connection so it gets a form, the show view is read-only plain data
      if (this.routeName.startsWith('vault.cluster.secrets.backend.edit')) {
        return new DatabaseConnectionForm(connection, {
          isNew: false,
          isEnterprise: this.version.isEnterprise,
        });
      }
      return connection;
    } catch (error) {
      const { message, status, response } = await this.api.parseError(error);
      // the application route sends control group errors to the approval page
      if (response?.isControlGroupError) throw response;
      // preserve the status so the error substate can render the not found / forbidden views
      throw Object.assign(new Error(message), { httpStatus: status, errors: [message] });
    }
  },

  async fetchDatabaseConnectionCapabilities(backend, name) {
    const connectionPath = this.capabilitiesService.pathFor('databaseConnection', { backend, id: name });
    const resetPath = this.capabilitiesService.pathFor('databaseReset', { backend, id: name });
    const rotateRootPath = this.capabilitiesService.pathFor('databaseRotateRoot', { backend, id: name });
    const rolePath = this.capabilitiesService.pathFor('databaseRoleWildcard', { backend });
    const staticRolePath = this.capabilitiesService.pathFor('databaseStaticRoleWildcard', { backend });

    const capabilities = await this.capabilitiesService.fetch([
      connectionPath,
      resetPath,
      rotateRootPath,
      rolePath,
      staticRolePath,
    ]);
    return {
      canRead: capabilities[connectionPath]?.canRead,
      canEdit: capabilities[connectionPath]?.canUpdate,
      canDelete: capabilities[connectionPath]?.canDelete,
      canReset: capabilities[resetPath]?.canUpdate || capabilities[resetPath]?.canCreate,
      canRotateRoot: capabilities[rotateRootPath]?.canUpdate || capabilities[rotateRootPath]?.canCreate,
      canAddRole: capabilities[rolePath]?.canCreate || capabilities[staticRolePath]?.canCreate,
    };
  },

  // Reads a database role. A role may be static or dynamic and the route does not always know which,
  // so both are requested unless the type is given. The edit view changes the role so it gets a form,
  // the show view is read-only plain data.
  async fetchDatabaseRole(backend, name, type) {
    const readStatic = () => this.api.secrets.databaseReadStaticRole(name, backend);
    const readDynamic = () => this.api.secrets.databaseReadRole(name, backend);
    let staticResp, dynamicResp;
    if (type === 'static') {
      [staticResp] = await Promise.allSettled([readStatic()]);
    } else if (type === 'dynamic') {
      [dynamicResp] = await Promise.allSettled([readDynamic()]);
    } else {
      [staticResp, dynamicResp] = await Promise.allSettled([readStatic(), readDynamic()]);
    }

    const succeeded = [staticResp, dynamicResp].find((resp) => resp?.status === 'fulfilled');
    if (!succeeded) {
      const failures = [staticResp, dynamicResp].filter(Boolean);
      const parsed = await Promise.all(failures.map((resp) => this.api.parseError(resp.reason)));
      // the application route sends control group errors to the approval page
      const controlGroup = parsed.find(({ response }) => response?.isControlGroupError);
      if (controlGroup) {
        throw controlGroup.response;
      }
      // prefer the error with the higher status, a 404 on one type should not hide a 403 on the other
      const worst = parsed.reduce((a, b) => ((b.status ?? 0) > (a.status ?? 0) ? b : a));
      // preserve the status so the error substate can render the not found / forbidden views
      throw Object.assign(new Error(worst.message), { httpStatus: worst.status, errors: [worst.message] });
    }

    const roleType = succeeded === staticResp ? 'static' : 'dynamic';
    const role = normalizeDatabaseRole(succeeded.value?.data, name, backend, roleType);
    if (this.routeName.startsWith('vault.cluster.secrets.backend.edit')) {
      return new DatabaseRoleForm(role, { isNew: false });
    }
    return role;
  },

  async fetchDatabaseRoleCapabilities(backend, name, type) {
    const params = { backend, id: name };
    const rolePath = this.capabilitiesService.pathFor(
      type === 'static' ? 'databaseStaticRole' : 'databaseRole',
      params
    );
    const credsPath = this.capabilitiesService.pathFor('databaseCreds', params);
    const staticCredsPath = this.capabilitiesService.pathFor('databaseStaticCreds', params);
    const rotatePath = this.capabilitiesService.pathFor('databaseRotateRole', params);
    const dynamicWildcard = this.capabilitiesService.pathFor('databaseRoleWildcard', { backend });
    const staticWildcard = this.capabilitiesService.pathFor('databaseStaticRoleWildcard', { backend });

    const capabilities = await this.capabilitiesService.fetch([
      rolePath,
      credsPath,
      staticCredsPath,
      rotatePath,
      dynamicWildcard,
      staticWildcard,
    ]);
    return {
      canEditRole: capabilities[rolePath]?.canUpdate,
      canDelete: capabilities[rolePath]?.canDelete,
      canGenerateCredentials: capabilities[credsPath]?.canRead,
      canGetCredentials: capabilities[staticCredsPath]?.canRead,
      canRotateRoleCredentials: capabilities[rotatePath]?.canUpdate,
      canCreateDynamic: capabilities[dynamicWildcard]?.canCreate,
      canCreateStatic: capabilities[staticWildcard]?.canCreate,
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

  async fetchAwsRole(backend, name) {
    const data = await this.api.secrets.awsReadRole(name, backend);
    return new RoleAwsForm({ ...data, name, backend }, { isNew: false });
  },

  async fetchAwsRoleCapabilities(backend, name) {
    const rolePath = `${backend}/roles/${name}`;
    const credentialsPath = `${backend}/creds/${name}`;
    const capabilities = await this.capabilitiesService.fetch([rolePath, credentialsPath]);
    return {
      canRead: capabilities[rolePath]?.canRead,
      canEdit: capabilities[rolePath]?.canUpdate,
      canDelete: capabilities[rolePath]?.canDelete,
      canGenerate: capabilities[credentialsPath]?.canUpdate,
    };
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
      if (this.routeName.startsWith('vault.cluster.secrets.backend.edit')) {
        secretModel = new TransitKeyForm(secretModel, { isNew: false });
      }
    } else if (modelType === 'keymgmt/key') {
      secretModel = await this.fetchKeymgmtKey(backend, secret);
      capabilities = await this.fetchKeymgmtKeyCapabilities(backend, secret);
    } else if (modelType === 'keymgmt/provider') {
      secretModel = await this.fetchKeymgmtProvider(backend, secret);
      capabilities = await this.fetchKeymgmtProviderCapabilities(backend, secret);
    } else if (modelType === 'database/role') {
      secretModel = await this.fetchDatabaseRole(backend, secret, type);
      const roleType = secretModel.type;
      capabilities = await this.fetchDatabaseRoleCapabilities(backend, secret, roleType);
    } else if (modelType === 'database/connection') {
      secretModel = await this.fetchDatabaseConnection(backend, secret);
      capabilities = await this.fetchDatabaseConnectionCapabilities(backend, secret);
    } else if (modelType === 'role-ssh') {
      secretModel = await this.fetchSshRole(backend, secret);
      capabilities = await this.fetchSshRoleCapabilities(backend, secret);
    } else if (modelType === 'role-aws') {
      secretModel = await this.fetchAwsRole(backend, secret);
      capabilities = await this.fetchAwsRoleCapabilities(backend, secret);
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

    // Form-based models and plain data do not have Ember Data's setProperties method.
    const isFormModel = model.secret instanceof Form;
    if (
      !isFormModel &&
      !(model.secret instanceof TransitKeyResource) &&
      typeof model.secret.setProperties === 'function'
    ) {
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

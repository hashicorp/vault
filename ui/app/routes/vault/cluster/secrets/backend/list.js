/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { assert } from '@ember/debug';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { filterEnginesByMountCategory, isAddonEngine } from 'core/utils/all-engines-metadata';
import { paginate } from 'core/utils/paginate-list';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { getEnginePathParam } from 'vault/utils/backend-route-helpers';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';
import { getKeymgmtProviderIcon } from 'vault/utils/keymgmt-provider-utils';
import { getModelTypeForEngine } from 'vault/utils/model-helpers/secret-engine-helpers';
import { normalizePath } from 'vault/utils/path-encoding-helpers';
import { resolve } from 'rsvp';
import {
  SecretsApiKeyManagementListKeysListEnum,
  SecretsApiKeyManagementListKmsProvidersListEnum,
  SecretsApiTotpListKeysListEnum,
  SecretsApiSshListRolesListEnum,
} from '@hashicorp/vault-client-typescript';

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
  api: service(),
  capabilitiesService: service('capabilities'),
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
    // Keymgmt, TOTP, and SSH routes use API-backed forms instead of Ember Data models, so skip model hydration.
    if (effectiveType === 'keymgmt' || effectiveType === 'totp' || effectiveType === 'ssh') {
      return resolve();
    }

    return this.pathHelp.hydrateModel(modelType, backend).then(() => {
      this.store.unloadAll('capabilities');
    });
  },

  getModelType(type, tab) {
    return getModelTypeForEngine(type, { tab });
  },

  async fetchTotpKeys(backend, page, pageFilter) {
    try {
      const resp = await this.api.secrets.totpListKeys(backend, SecretsApiTotpListKeysListEnum.TRUE);
      const keys = resp.keys || [];

      const pathsToFetch = keys.map((name) => this.capabilitiesService.pathFor('totpKey', { backend, name }));
      const capabilities = pathsToFetch.length ? await this.capabilitiesService.fetch(pathsToFetch) : {};

      const items = keys.map((name) => {
        const keyPath = this.capabilitiesService.pathFor('totpKey', { backend, name });
        return {
          id: name,
          name,
          backend,
          canRead: capabilities[keyPath]?.canRead || false,
          canDelete: capabilities[keyPath]?.canDelete || false,
        };
      });

      return paginate(items, { page, filter: pageFilter });
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return [];
      }
      throw error;
    }
  },

  async fetchKeysWithCapabilities(backend) {
    const { keys } = await this.api.secrets.keyManagementListKeys(
      backend,
      SecretsApiKeyManagementListKeysListEnum.TRUE
    );

    // Fetch capabilities for all keys
    const pathsToFetch = (keys || []).flatMap((keyName) => {
      const keyPath = this.capabilitiesService.pathFor('keymgmtKey', { backend, name: keyName });
      return [keyPath];
    });

    const capabilities = await this.capabilitiesService.fetch(pathsToFetch);

    // Transform string array into objects for list display with capabilities
    const keysList = (keys || []).map((keyName) => {
      const keyPath = this.capabilitiesService.pathFor('keymgmtKey', { backend, name: keyName });
      return {
        id: keyName,
        name: keyName,
        backend,
        icon: 'key',
        type: 'key',
        canRead: capabilities[keyPath]?.canRead || false,
        canEdit: capabilities[keyPath]?.canUpdate || false,
        canDelete: capabilities[keyPath]?.canDelete || false,
      };
    });

    return { keysList, capabilities };
  },

  async fetchKeymgmtKeys(backend, page, pageFilter) {
    try {
      const { keysList } = await this.fetchKeysWithCapabilities(backend);
      return paginate(keysList, { page, filter: pageFilter });
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return [];
      }
      throw error;
    }
  },

  async fetchProvidersWithCapabilities(backend) {
    const { keys: providerNames } = await this.api.secrets.keyManagementListKmsProviders(
      backend,
      SecretsApiKeyManagementListKmsProvidersListEnum.TRUE
    );

    const providersWithData = await Promise.all(
      (providerNames || []).map(async (providerName) => {
        const { data } = await this.api.secrets.keyManagementReadKmsProvider(providerName, backend);
        return {
          ...data,
          id: providerName,
          name: providerName,
          backend,
          type: 'provider',
          icon: getKeymgmtProviderIcon(data.provider),
        };
      })
    );

    const pathsToFetch = providersWithData.map((provider) =>
      this.capabilitiesService.pathFor('keymgmtProvider', { backend, id: provider.id })
    );
    const capabilities = await this.capabilitiesService.fetch(pathsToFetch);

    const providersList = providersWithData.map((provider) => {
      const providerPath = this.capabilitiesService.pathFor('keymgmtProvider', {
        backend,
        id: provider.id,
      });
      return {
        ...provider,
        canRead: capabilities[providerPath]?.canRead || false,
        canEdit: capabilities[providerPath]?.canUpdate || false,
        canDelete: capabilities[providerPath]?.canDelete || false,
      };
    });

    return { providersList, capabilities };
  },

  async fetchKeymgmtProviders(backend, page, pageFilter) {
    try {
      const { providersList } = await this.fetchProvidersWithCapabilities(backend);
      return paginate(providersList, { page, filter: pageFilter });
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return [];
      }
      throw error;
    }
  },

  async fetchSshRolesWithCapabilities(backend) {
    // Fetch roles and zero-address config in parallel (zero-address may 404 if never configured)
    const [listResponse, zeroAddressResult] = await Promise.allSettled([
      this.api.secrets.sshListRoles(backend, SecretsApiSshListRolesListEnum.TRUE),
      this.api.secrets.sshReadZeroAddressConfigurationRaw({ ssh_mount_path: backend }),
    ]);

    if (listResponse.status === 'rejected') throw listResponse.reason;

    const roles = this.api.keyInfoToArray(listResponse.value);

    // Build set of zero-address role names from the config endpoint
    let zeroAddressRoles = new Set();
    if (zeroAddressResult.status === 'fulfilled') {
      const body = await zeroAddressResult.value.raw.json();
      const names = body?.data?.roles;
      if (Array.isArray(names)) {
        zeroAddressRoles = new Set(names);
      }
    }

    // Build all capability paths for batch fetch
    const zeroAddressPath = this.capabilitiesService.pathFor('sshZeroAddress', { backend });
    const rolePaths = roles.map((role) => ({
      role: this.capabilitiesService.pathFor('sshRole', { backend, id: role.id }),
      credentials: this.capabilitiesService.pathFor('sshCredentials', { backend, id: role.id }),
      sign: this.capabilitiesService.pathFor('sshSign', { backend, id: role.id }),
    }));

    // Fetch all capabilities in a single request
    const allPaths = [
      ...rolePaths.flatMap((paths) => [paths.role, paths.credentials, paths.sign]),
      zeroAddressPath,
    ];
    const capabilities = await this.capabilitiesService.fetch(allPaths);

    // Merge role data with capabilities
    return roles.map((role, index) => {
      const paths = rolePaths[index];
      return {
        ...role,
        backend,
        zero_address: zeroAddressRoles.has(role.id),
        canRead: capabilities[paths.role]?.canRead || false,
        canEdit: capabilities[paths.role]?.canUpdate || false,
        canDelete: capabilities[paths.role]?.canDelete || false,
        canGenerate: capabilities[paths.credentials]?.canUpdate || false,
        canSign: capabilities[paths.sign]?.canUpdate || false,
        canEditZeroAddress: capabilities[zeroAddressPath]?.canUpdate || false,
      };
    });
  },

  async fetchSshRoles(backend, page, pageFilter) {
    try {
      const roles = await this.fetchSshRolesWithCapabilities(backend);
      return paginate(roles, { page, filter: pageFilter });
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return [];
      }
      throw error;
    }
  },

  async model(params) {
    const secret = this.secretParam() || '';
    const backend = getEnginePathParam(this);
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    const effectiveType = getEffectiveEngineType(backendModel.engineType);
    const modelType = this.getModelType(effectiveType, params.tab);

    // Handle TOTP, keymgmt and ssh resources with API service
    let secrets;
    if (effectiveType === 'totp') {
      const page = getValidPage(params.page);
      secrets = await this.fetchTotpKeys(backend, page, params.pageFilter);
      this.set('has404', false);
    } else if (effectiveType === 'keymgmt') {
      const page = getValidPage(params.page);
      const filter = params.pageFilter;
      secrets =
        params.tab === 'provider'
          ? await this.fetchKeymgmtProviders(backend, page, filter)
          : await this.fetchKeymgmtKeys(backend, page, filter);

      this.set('has404', false);
    } else if (effectiveType === 'ssh') {
      const page = getValidPage(params.page);
      const filter = params.pageFilter;
      secrets = await this.fetchSshRoles(backend, page, filter);
      this.set('has404', false);
    } else {
      try {
        secrets = await this.pagination.lazyPaginatedQuery(modelType, {
          id: secret,
          backend,
          responsePath: 'data.keys',
          page: getValidPage(params.page),
          pageFilter: params.pageFilter,
        });
        this.set('has404', false);
      } catch (err) {
        if (backendModel && err.httpStatus === 404) {
          secrets = [];
        } else {
          throw err;
        }
      }
    }

    return {
      secret,
      secrets,
    };
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

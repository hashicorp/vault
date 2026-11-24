/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import SecretsEngineForm from 'vault/forms/secrets/engine';
import { getExternalPluginNameFromBuiltin } from 'vault/utils/external-plugin-helpers';
import { getAllVersionsForEngineType, type EngineVersionInfo } from 'vault/utils/plugin-catalog-helpers';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';
import { IdentityApiOidcListKeysListEnum } from '@hashicorp/vault-client-typescript/dist/apis/IdentityApi';

import type ApiService from 'vault/services/api';
import type PluginCatalogService from 'vault/services/plugin-catalog';

export default class VaultClusterSecretsEnableCreateRoute extends Route {
  @service('plugin-catalog') declare readonly pluginCatalog: PluginCatalogService;
  @service declare api: ApiService;

  async model(params: { mount_type: string }) {
    const { mount_type } = params;

    const defaults = {
      path: mount_type, // Default path to match the engine type
      config: { listing_visibility: false },
      kv_config: {
        max_versions: 0,
        cas_required: false,
        delete_version_after: undefined,
      },
      options: { version: 2 },
    };

    const form = new SecretsEngineForm(defaults, { isNew: true });
    // Explicitly set the type on the form after creation
    form.type = mount_type;
    // Apply type-specific defaults (e.g., PKI max lease TTL)
    form.applyTypeSpecificDefaults();

    // Fetch plugin catalog data to get available versions for this engine type
    const { availableVersions, hasUnversionedPlugins } = await this.fetchPluginVersions(form);
    // Get pinned version for this plugin type
    const pinnedVersion = await this.fetchPinnedVersion(mount_type, availableVersions);

    // WIF engines require fetching OIDC keys to populate an additional form field
    const oidcKeys = await this.fetchOidcKeys(mount_type);

    return {
      form,
      availableVersions,
      hasUnversionedPlugins,
      pinnedVersion,
      oidcKeys,
    };
  }

  async fetchPluginVersions(form: SecretsEngineForm) {
    const pluginCatalogResponse = await this.pluginCatalog.fetchPluginCatalog();
    let availableVersions: EngineVersionInfo[] = [];
    let hasUnversionedPlugins = false;

    if (pluginCatalogResponse.data?.detailed) {
      const versionResult = getAllVersionsForEngineType(
        pluginCatalogResponse.data.detailed,
        form.type,
        'secret'
      );

      availableVersions = versionResult.versions;
      hasUnversionedPlugins = versionResult.hasUnversionedPlugins;

      // Set up the plugin version field with available versions
      form.setupPluginVersionField(availableVersions);
    }
    return { availableVersions, hasUnversionedPlugins };
  }

  async fetchPinnedVersion(mountType: string, availableVersions: EngineVersionInfo[]) {
    let pinnedVersion: string | null = null;

    // Only fetch external pinned version if there are external versions available
    const hasExternalVersions = availableVersions.some((version) => !version.isBuiltin);
    if (hasExternalVersions) {
      try {
        // Convert builtin type to external plugin name for API call
        const externalPluginName = getExternalPluginNameFromBuiltin(mountType);
        if (externalPluginName) {
          const response = await this.api.sys.pluginsCatalogPinsReadPinnedVersion(
            externalPluginName,
            'secret'
          );
          pinnedVersion = response?.version || null;
        }
      } catch (error) {
        // Silently handle errors - pins are optional
        pinnedVersion = null;
      }
    }
    return pinnedVersion;
  }

  async fetchOidcKeys(type: string) {
    // an additional form field for identity_token_key is displayed for WIF engines
    // fetch oidc keys to populate the search select component
    let oidcKeys: { id: string }[] = [];
    const isWIF = !!ALL_ENGINES.find((engine) => engine.type === type && engine.isWIF);
    if (isWIF) {
      try {
        const { keys } = await this.api.identity.oidcListKeys(IdentityApiOidcListKeysListEnum.TRUE);
        // SearchSelect requires options to be objects
        oidcKeys = keys?.map((key) => ({ id: key })) || [];
      } catch (e) {
        // swallow fetch error and fallback component will be used
      }
    }
    return oidcKeys;
  }
}

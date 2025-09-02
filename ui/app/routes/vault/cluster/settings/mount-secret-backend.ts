/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import SecretsEngineForm from 'vault/forms/secrets/engine';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type { ModelFrom } from 'vault/vault/route';
import type AuthService from 'vault/services/auth';
import type NamespaceService from 'vault/services/namespace';
import type PluginCatalogService from 'vault/services/plugin-catalog';

export type MountSecretBackendModel = ModelFrom<VaultClusterSettingsMountSecretBackendRoute>;

export default class VaultClusterSettingsMountSecretBackendRoute extends Route {
  @service declare readonly auth: AuthService;
  @service declare readonly namespace: NamespaceService;
  @service('plugin-catalog') declare readonly pluginCatalog: PluginCatalogService;

  async model() {
    const defaults = {
      config: { listing_visibility: false },
      kv_config: {
        max_versions: 0,
        cas_required: false,
        delete_version_after: undefined,
      },
      options: { version: 2 },
    };

    const secretEngineForm = new SecretsEngineForm(defaults, { isNew: true });

    // Fetch plugin catalog data to enhance the secret engines list
    const pluginCatalogResponse = await this.pluginCatalog.fetchPluginCatalog();

    return hash({
      form: secretEngineForm,
      pluginCatalogData: pluginCatalogResponse.data,
      pluginCatalogError: pluginCatalogResponse.error,
    });
  }
}

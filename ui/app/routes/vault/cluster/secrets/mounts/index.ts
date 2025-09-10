/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import SecretsEngineForm from 'vault/forms/secrets/engine';
import Router from 'vault/router';
import type PluginCatalogService from 'vault/services/plugin-catalog';

import type { ModelFrom } from 'vault/vault/route';

export type MountSecretBackendModel = ModelFrom<VaultClusterSecretsMountsIndexRouter>;

export default class VaultClusterSecretsMountsIndexRouter extends Route {
  @service declare router: Router;
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

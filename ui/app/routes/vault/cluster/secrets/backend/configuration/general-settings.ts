/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { getPluginVersionsFromEngineType } from 'vault/utils/plugin-catalog-helpers';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type PluginCatalogService from 'vault/services/plugin-catalog';

export default class SecretsBackendConfigurationGeneralSettingsRoute extends Route {
  @service declare readonly api: ApiService;
  @service('plugin-catalog') declare readonly pluginCatalog: PluginCatalogService;

  async model() {
    const secretsEngine = this.modelFor('vault.cluster.secrets.backend') as SecretsEngineResource;
    const { data } = await this.pluginCatalog.getRawPluginCatalogData();

    const versions = getPluginVersionsFromEngineType(data?.secret, secretsEngine.type);

    return { secretsEngine, versions };
  }
}

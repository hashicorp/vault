/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

import type SecretsEngineResource from 'vault/resources/secrets/engine';

export default class SecretsBackendConfigurationGeneralSettingsRoute extends Route {
  async model() {
    const secretsEngine = this.modelFor('vault.cluster.secrets.backend') as SecretsEngineResource;
    // TODO: get list of versions using the sys/plugins/catalog endpoint.
    return { secretsEngine };
  }
}

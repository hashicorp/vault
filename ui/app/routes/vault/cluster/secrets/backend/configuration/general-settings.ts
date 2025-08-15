/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import engineDisplayData from 'vault/helpers/engines-display-data';

import type SecretsEngineResource from 'vault/resources/secrets/engine';

export default class SecretsBackendConfigurationGeneralSettingsRoute extends Route {
  async model() {
    const secretsEngine = this.modelFor('vault.cluster.secrets.backend') as SecretsEngineResource;
    const engineData = engineDisplayData(secretsEngine.type);

    return {
      secretsEngine,
      engineData,
    };
  }
}

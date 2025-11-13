/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import engineDisplayData from 'vault/helpers/engines-display-data';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type RouterService from '@ember/routing/router-service';

/**
 * @module ConfigureTabs handles the display of the ttl picker for the lease duration card in general settings.
 * 
 * @example
 * <SecretEngine::ConfigureTabs
    @model={{this.model}}
    />
 *
 * @param {object} secretsEngine - secrets engine resource.
 * @param {boolean} config - config model for the secret engine.
 */

interface Args {
  secretsEngine: SecretsEngineResource;
  config: Record<string, unknown>;
}

export default class ConfigureTabs extends Component<Args> {
  @service declare readonly router: RouterService;

  get routeName() {
    if (this.router.currentRouteName === 'vault.cluster.secrets.backend.configuration.edit') {
      return 'vault.cluster.secrets.backend.configuration.edit';
    }
    return 'vault.cluster.secrets.backend.configuration.plugin-settings';
  }

  get engineRoute() {
    const engineData = engineDisplayData(this.args?.secretsEngine?.type);
    const baseUrl = 'vault.cluster.secrets.backend.';
    if (this.args?.config && Object.keys(this.args?.config).length > 0)
      return `${baseUrl}${this.args?.secretsEngine?.type}.configuration`;
    return `${baseUrl}${engineData?.engineConfigureRoute}`;
  }

  get pluginSettingsRoute() {
    const engineData = engineDisplayData(this.args?.secretsEngine?.type);
    if (engineData?.engineConfigureRoute) {
      return this.engineRoute;
    }

    return this.routeName;
  }
}

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module DashboardVaultConfigurationCard
 * DashboardVaultConfigurationCard component are used to display vault configuration.
 *
 * @example
 * ```js
 * <DashboardVaultConfigurationCard @vaultConfiguration={{@model.vaultConfiguration}} />
 * ```
 * @param {object} vaultConfiguration - object of vault configuration key/values
 */

export default class DashboardSecretsEnginesCard extends Component {
  get tlsDisabled() {
    const tlsDisableConfig = this.args.vaultConfiguration?.listeners.find((listener) => {
      if (listener.config && listener.config.tls_disable) return listener.config.tls_disable;
    });

    return tlsDisableConfig?.config.tls_disable ? 'Enabled' : 'Disabled';
  }
}

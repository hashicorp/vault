/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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
  get tls() {
    // since the default for tlsDisable is false it may not be in the config
    // consider tls enabled if tlsDisable is undefined or false AND both tlsCertFile and tlsKeyFile are defined
    const tlsListener = this.args.vaultConfiguration?.listeners.find((listener) => {
      const { tlsDisable, tlsCertFile, tlsKeyFile } = listener.config || {};
      return !tlsDisable && tlsCertFile && tlsKeyFile;
    });

    return tlsListener ? 'Enabled' : 'Disabled';
  }
}

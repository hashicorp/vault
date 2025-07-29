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
    // since the default for tls_disable is false it may not be in the config
    // consider tls enabled if tls_disable is undefined or false AND both tls_cert_file and tls_key_file are defined
    const tlsListener = this.args.vaultConfiguration?.listeners.find((listener) => {
      const { tls_disable, tls_cert_file, tls_key_file } = listener.config || {};
      return !tls_disable && tls_cert_file && tls_key_file;
    });

    return tlsListener ? 'Enabled' : 'Disabled';
  }
}

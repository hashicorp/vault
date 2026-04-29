/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface VaultListenerConfig {
  tls_disable?: boolean;
  tls_cert_file?: string;
  tls_key_file?: string;
}

interface VaultListener {
  config?: VaultListenerConfig;
}

interface VaultStorage {
  type?: string;
}

interface VaultConfiguration {
  api_addr?: string;
  default_lease_ttl?: number | string;
  max_lease_ttl?: number | string;
  log_format?: string;
  log_level?: string;
  storage?: VaultStorage;
  listeners: VaultListener[];
}

export interface Args {
  vaultConfiguration?: VaultConfiguration;
}

/**
 * @module Dashboard::Widgets::ClusterConfiguration
 * Dashboard widget component to display vault cluster configuration details.
 *
 * @example
 * ```js
 * <Dashboard::Widgets::ClusterConfiguration @vaultConfiguration={{@model.vaultConfiguration}} />
 * ```
 * @param {VaultConfiguration} vaultConfiguration - Vault configuration object with listeners, storage, etc.
 */

export default class DashboardWidgetsClusterConfiguration extends Component<Args> {
  get tls() {
    // since the default for tls_disable is false it may not be in the config
    // consider tls enabled if tls_disable is undefined or false AND both tls_cert_file and tls_key_file are defined
    const tlsListener = this.args.vaultConfiguration?.listeners?.find((listener) => {
      const { tls_disable, tls_cert_file, tls_key_file } = listener.config || {};
      return !tls_disable && tls_cert_file && tls_key_file;
    });

    return tlsListener ? 'Enabled' : 'Disabled';
  }
}

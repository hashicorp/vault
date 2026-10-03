/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type SecretMountPath from 'vault/services/secret-mount-path';

export default class KmipConfigurationPageComponent extends Component {
  @service declare readonly secretMountPath: SecretMountPath;

  displayFields = [
    { key: 'listen_addrs', label: 'Listen addresses' },
    { key: 'default_tls_client_key_bits', label: 'Default TLS client key bits' },
    { key: 'default_tls_client_key_type', label: 'Default TLS client key type' },
    { key: 'default_tls_client_ttl', label: 'Default TLS client TTL' },
    { key: 'server_hostnames', label: 'Server hostnames' },
    { key: 'server_ips', label: 'Server IPs' },
    { key: 'tls_ca_key_bits', label: 'TLS CA key bits' },
    { key: 'tls_ca_key_type', label: 'TLS CA key type' },
    { key: 'tls_min_version', label: 'Minimum TLS version' },
  ];
}

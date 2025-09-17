/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const keyMgmtMockModel = {
  secretsEngine: {
    accessor: 'keymgmt_accessor',
    config: {
      default_lease_ttl: 2073600,
      force_no_cache: false,
      listing_visibility: 'hidden',
      max_lease_ttl: 4320000,
    },
    description: 'hello',
    external_entropy_access: false,
    local: true,
    options: {},
    path: 'keymgmt/',
    plugin_version: '',
    running_plugin_version: 'v0.17.1+builtin',
    running_sha256: '',
    seal_wrap: false,
    type: 'keymgmt',
    uuid: '4ea92618-5b52-f89a-9cbe-b65dc7e65689',
    id: 'keymgmt',
    backendConfigurationLink: `vault.cluster.secrets.backend.configuration`,
  },
  versions: ['v0.17.1+builtin'],
};

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import apiPath from 'vault/utils/api-path';

export const SUDO_PATHS = [
  'sys/seal',
  'sys/replication/performance/primary/secondary-token',
  'sys/replication/dr/primary/secondary-token',
  'sys/replication/reindex',
  'sys/leases/lookup/',
];

export const SUDO_PATH_PREFIXES = ['sys/leases/revoke-prefix', 'sys/leases/revoke-force'];

export const PATH_MAP = {
  customLogin: apiPath`sys/config/ui/login/default-auth/${'id'}`,
  customMessages: apiPath`sys/config/ui/custom-messages/${'id'}`,
  syncActivate: apiPath`sys/activation-flags/secrets-sync/activate`,
  syncDestination: apiPath`sys/sync/destinations/${'type'}/${'name'}`,
  syncSetAssociation: apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/set`,
  syncRemoveAssociation: apiPath`sys/sync/destinations/${'type'}/${'name'}/associations/remove`,
  kvConfig: apiPath`${'path'}/config`,
};

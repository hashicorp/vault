/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// This model represents the capabilities on a given `path`
// `path` is also the primaryId
// https://developer.hashicorp.com/vault/docs/concepts/policies#capabilities

import Model, { attr } from '@ember-data/model';

import { computed } from '@ember/object';

const SUDO_PATHS = [
  'sys/seal',
  'sys/replication/performance/primary/secondary-token',
  'sys/replication/dr/primary/secondary-token',
  'sys/replication/reindex',
  'sys/leases/lookup/',
];

const SUDO_PATH_PREFIXES = ['sys/leases/revoke-prefix', 'sys/leases/revoke-force'];

export { SUDO_PATHS, SUDO_PATH_PREFIXES };

const computedCapability = function (capability) {
  return computed('path', 'capabilities', 'capabilities.[]', function () {
    const capabilities = this.capabilities;
    const path = this.path;
    if (!capabilities) {
      return false;
    }
    if (capabilities.includes('root')) {
      return true;
    }
    if (capabilities.includes('deny')) {
      return false;
    }
    // if the path is sudo protected, they'll need sudo + the appropriate capability
    if (SUDO_PATHS.includes(path) || SUDO_PATH_PREFIXES.find((item) => path.startsWith(item))) {
      return capabilities.includes('sudo') && capabilities.includes(capability);
    }
    return capabilities.includes(capability);
  });
};

export default Model.extend({
  path: attr('string'),
  capabilities: attr('array'),
  allowedParameters: attr(),
  deniedParameters: attr(),
  canCreate: computedCapability('create'),
  canDelete: computedCapability('delete'),
  canList: computedCapability('list'),
  canPatch: computedCapability('patch'),
  canRead: computedCapability('read'),
  canSudo: computedCapability('sudo'),
  canUpdate: computedCapability('update'),
});
